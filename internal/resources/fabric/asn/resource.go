package asn

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	typePublic    = "PUBLIC_ASN"
	signaturePath = "/ownershipChallengeDetails/signature"
	deletedMarker = "tf-marker-for-deleted-asn"
)

var (
	waitDelay       = 5 * time.Second
	waitMinInterval = 5 * time.Second
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_asn",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema(ctx)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest, err := buildCreateRequest(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("failed building ASN create request", err.Error())
		return
	}

	created, _, apiErr := client.ASNsApi.CreateAsn(ctx).CreateAsnRequest(createRequest).Execute()
	if apiErr != nil {
		resp.Diagnostics.AddError("api error creating asn", equinix_errors.FormatFabricError(apiErr).Error())
		return
	}
	id := created.GetUuid()

	createTimeout, diags := plan.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	signature := signatureFromPlan(ctx, &plan)
	waitForSignature := plan.Type.ValueString() == typePublic && signature == ""
	if signature != "" {

		if _, waitErr := getWaiter(ctx, client, id, createTimeout, []string{string(fabricv4.ASNSTATE_REQUESTED)}, []string{string(fabricv4.ASNSTATE_PENDING)}).WaitForStateContext(ctx); waitErr != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("failed waiting for ASN %s to accept the ownership signature", id), waitErr.Error())
			resp.State.SetAttribute(ctx, path.Root("id"), id)
			return
		}
		if patchErr := patchSignature(ctx, client, id, signature); patchErr != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("api error submitting ownership signature for ASN %s", id), patchErr.Error())
			resp.State.SetAttribute(ctx, path.Root("id"), id)
			return
		}
	}

	checked, waitErr := getWaiter(ctx, client, id, createTimeout, createPendingStates(), createTargetStates(waitForSignature)).WaitForStateContext(ctx)
	if waitErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed waiting for ASN %s to be provisioned", id), waitErr.Error())
		resp.State.SetAttribute(ctx, path.Root("id"), id)
		return
	}

	resp.Diagnostics.Append(plan.parse(ctx, checked.(*fabricv4.Asn))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	asn, httpResp, err := client.ASNsApi.GetAsnById(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, httpResp.StatusCode) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed retrieving ASN %s", state.ID.ValueString()),
			equinix_errors.FormatFabricError(err).Error(),
		)
		return
	}

	if asn.GetState() == fabricv4.ASNSTATE_DEPROVISIONED {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(state.parse(ctx, asn)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	var plan, state ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	updateTimeout, diags := plan.Timeouts.Update(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	signature := signatureFromPlan(ctx, &plan)
	if signature != "" && signature != signatureFromPlan(ctx, &state) {
		if err := patchSignature(ctx, client, id, signature); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("api error submitting ownership signature for ASN %s", id), err.Error())
			return
		}
	}

	checked, waitErr := getWaiter(ctx, client, id, updateTimeout, createPendingStates(), createTargetStates(false)).WaitForStateContext(ctx)
	if waitErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed waiting for ASN %s to be provisioned", id), waitErr.Error())
		return
	}

	resp.Diagnostics.Append(plan.parse(ctx, checked.(*fabricv4.Asn))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, deleteResp, err := client.ASNsApi.DeleteAsnById(ctx, id).Execute()
	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed deleting ASN %s", id),
				equinix_errors.FormatFabricError(err).Error(),
			)
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, waitErr := getDeleteWaiter(ctx, client, id, deleteTimeout).WaitForStateContext(ctx); waitErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed waiting for ASN %s to be deprovisioned", id), waitErr.Error())
	}
}

func buildCreateRequest(ctx context.Context, plan *ResourceModel) (fabricv4.CreateAsnRequest, error) {
	project, diags := plan.Project.ToPtr(ctx)
	if diags.HasError() || project == nil {
		return fabricv4.CreateAsnRequest{}, fmt.Errorf("failed reading project from plan")
	}

	req := fabricv4.CreateAsnRequest{
		Project: fabricv4.Project{ProjectId: project.ProjectID.ValueString()},

		AdditionalProperties: map[string]interface{}{typeKey: plan.Type.ValueString()},
	}

	if !plan.Number.IsNull() && !plan.Number.IsUnknown() {
		req.SetNumber(plan.Number.ValueInt64())
	}
	if !plan.RirRegion.IsNull() && !plan.RirRegion.IsUnknown() {
		req.SetRirRegion(fabricv4.AsnRirRegion(plan.RirRegion.ValueString()))
	}
	if !plan.Range.IsNull() && !plan.Range.IsUnknown() {
		req.SetRange(fabricv4.AsnRange(plan.Range.ValueString()))
	}
	if !plan.Location.IsNull() && !plan.Location.IsUnknown() {
		loc, d := plan.Location.ToPtr(ctx)
		if d.HasError() || loc == nil {
			return fabricv4.CreateAsnRequest{}, fmt.Errorf("failed reading location from plan")
		}
		req.SetLocation(fabricv4.SimplifiedLocation{MetroCode: loc.MetroCode.ValueStringPointer()})
	}
	if !plan.OwnershipChallengeDetails.IsNull() && !plan.OwnershipChallengeDetails.IsUnknown() {
		challenge, d := plan.OwnershipChallengeDetails.ToPtr(ctx)
		if d.HasError() || challenge == nil {
			return fabricv4.CreateAsnRequest{}, fmt.Errorf("failed reading ownership challenge details from plan")
		}
		req.SetOwnershipChallengeDetails(fabricv4.AsnOwnershipChallengeDetails{
			Type: fabricv4.AsnOwnershipChallengeType(challenge.Type.ValueString()),
		})
	}

	return req, nil
}

func signatureFromPlan(ctx context.Context, m *ResourceModel) string {
	if m.OwnershipChallengeDetails.IsNull() || m.OwnershipChallengeDetails.IsUnknown() {
		return ""
	}
	challenge, diags := m.OwnershipChallengeDetails.ToPtr(ctx)
	if diags.HasError() || challenge == nil || challenge.Signature.IsNull() || challenge.Signature.IsUnknown() {
		return ""
	}
	return challenge.Signature.ValueString()
}

func patchSignature(ctx context.Context, client *fabricv4.APIClient, id, signature string) error {
	_, _, err := client.ASNsApi.PatchAsnById(ctx, id).AsnPatchOperation([]fabricv4.AsnPatchOperation{{
		Op:    fabricv4.ASNPATCHOPERATIONOP_REPLACE,
		Path:  signaturePath,
		Value: signature,
	}}).Execute()
	if err != nil {
		return equinix_errors.FormatFabricError(err)
	}
	return nil
}

func createPendingStates() []string {
	return []string{
		string(fabricv4.ASNSTATE_REQUESTED),
		string(fabricv4.ASNSTATE_PROVISIONING),
	}
}

func createTargetStates(allowPending bool) []string {
	states := []string{
		string(fabricv4.ASNSTATE_PROVISIONED),
		string(fabricv4.ASNSTATE_ATTACHED),
	}
	if allowPending {
		states = append(states, string(fabricv4.ASNSTATE_PENDING))
	}
	return states
}

func getWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration, pending, target []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (any, string, error) {
			asn, _, err := client.ASNsApi.GetAsnById(ctx, id).Execute()
			if err != nil {
				return nil, "", err
			}
			return asn, string(asn.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      waitDelay,
		MinTimeout: waitMinInterval,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{

		Pending: []string{
			string(fabricv4.ASNSTATE_REQUESTED),
			string(fabricv4.ASNSTATE_PROVISIONING),
			string(fabricv4.ASNSTATE_PENDING),
			string(fabricv4.ASNSTATE_PROVISIONED),
			string(fabricv4.ASNSTATE_ATTACHED),
			string(fabricv4.ASNSTATE_DEPROVISIONING),
		},
		Target: []string{deletedMarker, string(fabricv4.ASNSTATE_DEPROVISIONED)},
		Refresh: func() (any, string, error) {
			asn, resp, err := client.ASNsApi.GetAsnById(ctx, id).Execute()
			if err != nil {
				if resp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return asn, deletedMarker, nil
				}
				return nil, "", err
			}
			return asn, string(asn.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      waitDelay,
		MinTimeout: waitMinInterval,
	}
}
