package stream_subscription

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_stream_subscription",
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

	stream, _, err := client.StreamSubscriptionsApi.CreateStreamSubscriptions(ctx, plan.StreamID.ValueString()).StreamSubscriptionPostRequest(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed creating stream subscription", equinix_errors.FormatFabricError(err).Error())
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
			fmt.Sprintf("failed creating stream subscription %s", stream.GetUuid()), err.Error())
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
			fmt.Sprintf("failed retrieving stream subscription %s", id), equinix_errors.FormatFabricError(err).Error())
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
		resp.Diagnostics.AddWarning("no updatable fields have changed",
			"terraform detected a config change, but it is for a field that isn't updatable for the stream subscription resource. please revert to prior config")
		return
	}

	updateRequest.SetName(plan.Name.ValueString())
	updateRequest.SetDescription(plan.Description.ValueString())

	_, _, err := client.StreamsApi.UpdateStreamByUuid(ctx, id).StreamPutRequest(updateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed updating stream subscription %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updateWaiter := getCreateUpdateWaiter(ctx, client, id, updateTimeout)
	streamSubscriptionChecked, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed updating stream subscription %s", id), err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, streamSubscriptionChecked.(*fabricv4.StreamSubscription))...)
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

func buildCreateRequest(ctx context.Context, plan ResourceModel) (fabricv4.StreamSubscriptionPostRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.StreamSubscriptionPostRequest{}

	request.SetName(plan.Name.ValueString())
	request.SetType(fabricv4.StreamSubscriptionPostRequestType(plan.Type.ValueString()))
	request.SetDescription(plan.Description.ValueString())
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		request.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.MetricSelector.IsNull() && !plan.MetricSelector.IsUnknown() {
		// Build MetricSelector
	}

	if !plan.EventSelector.IsNull() && !plan.EventSelector.IsUnknown() {
		// Build EventSelector
	}

	if !plan.Sink.IsNull() && !plan.Sink.IsUnknown() {
		// Update sink request
	}

	return request, diags
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, streamID, streamSubscriptionID string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			streamSubscription, _, err := client.StreamSubscriptionsApi.GetStreamSubscriptionByUuid(ctx, streamID, streamSubscriptionID).Execute()
			if err != nil {
				return 0, "", err
			}
			return streamSubscription, string(streamSubscription.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, streamID, streamSubscriptionID string, timeout time.Duration) *retry.StateChangeConf {
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
			streamSubscription, resp, err := client.StreamSubscriptionsApi.GetStreamSubscriptionByUuid(ctx, streamID, streamSubscriptionID).Execute()
			if err != nil {
				if resp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return streamSubscription, deletedMarker, nil
				}
				return 0, "", err
			}
			return streamSubscription, string(streamSubscription.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
