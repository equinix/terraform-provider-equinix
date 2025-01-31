package stream

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_stream",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = resourceSchema(ctx)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {

	var plan ResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	createRequest, diags := buildCreateRequest(ctx, plan)
	if diags.HasError() {
		return
	}

	stream, _, err := client.StreamsApi.CreateStreams(ctx).StreamPostRequest(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed creating Stream", equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, stream.GetUuid(), createTimeout)
	streamChecked, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed creating Stream %s", stream.GetUuid()), err.Error())
		return
	}

	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, streamChecked.(*fabricv4.Stream))...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	stream, _, err := client.StreamsApi.GetStreamByUuid(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed retrieving Stream %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, stream)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Retrieve values from plan
	var state, plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	updateRequest := fabricv4.StreamPutRequest{}

	needsUpdate := plan.Name.ValueString() != state.Name.ValueString() ||
		plan.Description.ValueString() != state.Description.ValueString()

	if !needsUpdate {
		resp.Diagnostics.AddWarning("No updatable fields have changed",
			"Terraform detected a config change, but it is for a field that isn't updatable for the stream resource. Please revert to prior config")
		return
	}

	updateRequest.SetName(plan.Name.ValueString())
	updateRequest.SetDescription(plan.Description.ValueString())

	_, _, err := client.StreamsApi.UpdateStreamByUuid(ctx, id).StreamPutRequest(updateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed updating Stream %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updateWaiter := getCreateUpdateWaiter(ctx, client, id, updateTimeout)
	streamChecked, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed updating Stream %s", id), err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, streamChecked.(*fabricv4.Stream))...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve the API client
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Retrieve the current state
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, deleteResp, err := client.StreamsApi.DeleteStreamByUuid(ctx, id).Execute()
	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed deleting Stream %s", id), equinix_errors.FormatFabricError(err).Error())
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	deleteWaiter := getDeleteWaiter(ctx, client, id, deleteTimeout)
	_, err = deleteWaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed deleting Stream %s", id), err.Error())
		return
	}

}

func buildCreateRequest(ctx context.Context, plan ResourceModel) (fabricv4.StreamPostRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.StreamPostRequest{}

	request.SetName(plan.Name.ValueString())
	request.SetType(fabricv4.StreamPostRequestType(plan.Type.ValueString()))
	request.SetDescription(plan.Description.ValueString())

	var project ProjectModel
	if !plan.Project.IsNull() && !plan.Project.IsUnknown() {
		diags = plan.Project.As(ctx, &project, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fabricv4.StreamPostRequest{}, diags
		}
		request.SetProject(fabricv4.Project{ProjectId: project.ProjectID.ValueString()})
	}

	return request, diags
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			stream, _, err := client.StreamsApi.GetStreamByUuid(ctx, id).Execute()
			if err != nil {
				return 0, "", err
			}
			return stream, stream.GetState(), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the connection appears to be deleted successfully based on
	// status code
	deletedMarker := "tf-marker-for-deleted-connection"
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.STREAMSUBSCRIPTIONSTATE_DEPROVISIONING),
		},
		Target: []string{
			deletedMarker,
			string(fabricv4.STREAMSUBSCRIPTIONSTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			stream, resp, err := client.StreamsApi.GetStreamByUuid(ctx, id).Execute()
			if err != nil {
				if resp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return stream, deletedMarker, nil
				}
				return 0, "", err
			}
			return stream, stream.GetState(), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
