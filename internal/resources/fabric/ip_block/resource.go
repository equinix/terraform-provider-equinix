package ipblock

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_ip_block",
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

	createRequest, err := buildCreateRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("failed building IP block create request", err.Error())
		return
	}

	ipBlock, _, apiErr := client.IPBlocksApi.SubmitIpBlock(ctx).SubmitIpBlockRequestBody(createRequest).Execute()
	if apiErr != nil {
		resp.Diagnostics.AddError("api error creating ip block", equinix_errors.FormatFabricError(apiErr).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkedBlock, waitErr := getCreateWaiter(ctx, client, ipBlock.GetUuid(), createTimeout).WaitForStateContext(ctx)
	if waitErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed waiting for IP block %s to become active", ipBlock.GetUuid()), waitErr.Error())
		return
	}

	resp.Diagnostics.Append(plan.parse(ctx, checkedBlock.(*fabricv4.IpBlock))...)
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

	ipBlock, _, err := client.IPBlocksApi.GetIpBlock(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed retrieving IP block %s", state.ID.ValueString()),
			equinix_errors.FormatFabricError(err).Error(),
		)
		return
	}

	resp.Diagnostics.Append(state.parse(ctx, ipBlock)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not reachable since all user-facing input fields use RequiresReplace.
// It is implemented defensively to surface an actionable message if the framework
// ever routes a change here.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"IP block is immutable",
		"No fields on equinix_fabric_ip_block support in-place updates. Terraform should have triggered a destroy/create cycle. Please file an issue if you see this warning.",
	)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, deleteResp, err := client.IPBlocksApi.DeleteIpBlockById(ctx, id).Execute()
	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed deleting IP block %s", id),
				equinix_errors.FormatFabricError(err).Error(),
			)
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, waitErr := getDeleteWaiter(ctx, client, id, deleteTimeout).WaitForStateContext(ctx)
	if waitErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed waiting for IP block %s to be deleted", id), waitErr.Error())
	}
}

func buildCreateRequest(ctx context.Context, plan ResourceModel) (fabricv4.SubmitIpBlockRequestBody, error) {
	var project ResourceProjectModel
	if diags := plan.Project.As(ctx, &project, basetypes.ObjectAsOptions{}); diags.HasError() {
		return fabricv4.SubmitIpBlockRequestBody{}, fmt.Errorf("failed reading project from plan")
	}

	req := fabricv4.SubmitIpBlockRequestBody{
		Type:    fabricv4.TypeOfIpBlockProduct(plan.Type.ValueString()),
		Project: fabricv4.IpBlockProjectRequest{ProjectId: project.ProjectID.ValueString()},
	}

	if !plan.Location.IsNull() && !plan.Location.IsUnknown() {
		var loc LocationModel
		if diags := plan.Location.As(ctx, &loc, basetypes.ObjectAsOptions{}); diags.HasError() {
			return fabricv4.SubmitIpBlockRequestBody{}, fmt.Errorf("failed reading location from plan")
		}
		req.SetLocation(fabricv4.IpBlockLocation{MetroCode: loc.MetroCode.ValueString()})
	}

	if !plan.PrefixLength.IsNull() && !plan.PrefixLength.IsUnknown() {
		v := plan.PrefixLength.ValueInt32()
		req.SetPrefixLength(v)
	}

	if !plan.Prefix.IsNull() && !plan.Prefix.IsUnknown() {
		v := plan.Prefix.ValueString()
		req.SetPrefix(v)
	}

	if !plan.Account.IsNull() && !plan.Account.IsUnknown() {
		var acc AccountModel
		if diags := plan.Account.As(ctx, &acc, basetypes.ObjectAsOptions{}); diags.HasError() {
			return fabricv4.SubmitIpBlockRequestBody{}, fmt.Errorf("failed reading account from plan")
		}
		req.SetAccount(fabricv4.IpBlockAccount{AccountNumber: acc.AccountNumber.ValueString()})
	}

	return req, nil
}

func getCreateWaiter(ctx context.Context, client *fabricv4.APIClient, uuid string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{string(fabricv4.IPBLOCKSTATE_PENDING)},
		Target:  []string{string(fabricv4.IPBLOCKSTATE_ACTIVE)},
		Refresh: func() (any, string, error) {
			ipBlock, _, err := client.IPBlocksApi.GetIpBlock(ctx, uuid).Execute()
			if err != nil {
				return nil, "", err
			}
			return ipBlock, string(ipBlock.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, uuid string, timeout time.Duration) *retry.StateChangeConf {
	deletedMarker := "tf-marker-for-deleted-ip-block"
	return &retry.StateChangeConf{
		Pending: []string{string(fabricv4.IPBLOCKSTATE_DELETING)},
		Target:  []string{deletedMarker, string(fabricv4.IPBLOCKSTATE_DELETED)},
		Refresh: func() (any, string, error) {
			ipBlock, resp, err := client.IPBlocksApi.GetIpBlock(ctx, uuid).Execute()
			if err != nil {
				if resp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return ipBlock, deletedMarker, nil
				}
				return nil, "", err
			}
			return ipBlock, string(ipBlock.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
