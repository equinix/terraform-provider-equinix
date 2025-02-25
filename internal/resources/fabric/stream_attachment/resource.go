package streamattachment

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_stream_attachment",
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

	assetID, asset, streamID := plan.AssetID.ValueString(), plan.Asset.ValueString(), plan.StreamID.ValueString()

	putRequest := fabricv4.StreamAssetPutRequest{}
	putRequest.SetMetricsEnabled(plan.MetricsEnabled.ValueBool())
	_, _, err := client.StreamsApi.UpdateStreamAssetByUuid(ctx, assetID, fabricv4.Asset(asset), streamID).StreamAssetPutRequest(putRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed creating stream attachment", equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, assetID, asset, streamID, createTimeout)
	attachment, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed creating stream attachment %s", attachment.(*fabricv4.StreamAsset).GetUuid()), err.Error())
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, attachment.(*fabricv4.StreamAsset))...)
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

	assetID, asset, streamID := state.AssetID.ValueString(), state.Asset.ValueString(), state.StreamID.ValueString()

	attachment, _, err := client.StreamsApi.GetStreamAssetByUuid(ctx, assetID, fabricv4.Asset(asset), streamID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed retrieving stream attachment %s", attachment.GetUuid()), equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, attachment)...)
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

	if plan.MetricsEnabled.ValueBool() == state.MetricsEnabled.ValueBool() {
		resp.Diagnostics.AddWarning("no updatable fields have changed",
			"terraform detected a config change, but it is for a field that isn't updatable for the stream attachment resource. please revert to prior config")
		return
	}

	id, assetID, asset, streamID := state.ID.ValueString(), plan.AssetID.ValueString(), plan.Asset.ValueString(), plan.StreamID.ValueString()

	putRequest := fabricv4.StreamAssetPutRequest{}
	putRequest.SetMetricsEnabled(plan.MetricsEnabled.ValueBool())
	_, _, err := client.StreamsApi.UpdateStreamAssetByUuid(ctx, assetID, fabricv4.Asset(asset), streamID).StreamAssetPutRequest(putRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed updating stream attachment %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updateWaiter := getCreateUpdateWaiter(ctx, client, assetID, asset, streamID, updateTimeout)
	attachment, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed updating stream attachment %s", id), err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, attachment.(*fabricv4.StreamAsset))...)
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

	id, assetID, asset, streamID := state.ID.ValueString(), state.AssetID.ValueString(), state.Asset.ValueString(), state.StreamID.ValueString()

	_, deleteResp, err := client.StreamsApi.DeleteStreamAssetByUuid(ctx, assetID, fabricv4.Asset(asset), streamID).Execute()
	if err != nil {

		//Design decision from API team was to return 400 for all errors instead of 404 for not found
		if deleteResp == nil || !slices.Contains([]int{http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed deleting stream attachment %s", id), equinix_errors.FormatFabricError(err).Error())
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	deleteWaiter := getDeleteWaiter(ctx, client, assetID, asset, streamID, deleteTimeout)
	_, err = deleteWaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed deleting stream attachment %s", id), err.Error())
		return
	}

}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, assetID, asset, streamID string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.STREAMASSETATTACHMENTSTATUS_ATTACHING),
		},
		Target: []string{
			string(fabricv4.STREAMASSETATTACHMENTSTATUS_ATTACHED),
		},
		Refresh: func() (interface{}, string, error) {
			stream, _, err := client.StreamsApi.GetStreamAssetByUuid(ctx, assetID, fabricv4.Asset(asset), streamID).Execute()
			if err != nil {
				return 0, "", err
			}
			return stream, string(stream.GetAttachmentStatus()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, assetID, asset, streamID string, timeout time.Duration) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the connection appears to be deleted successfully based on
	// status code
	deletedMarker := "tf-marker-for-deleted-stream-attachment"
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.STREAMASSETATTACHMENTSTATUS_ATTACHED),
			string(fabricv4.STREAMASSETATTACHMENTSTATUS_DETACHING),
		},
		Target: []string{
			deletedMarker,
			string(fabricv4.STREAMASSETATTACHMENTSTATUS_DETACHED),
		},
		Refresh: func() (interface{}, string, error) {
			stream, resp, err := client.StreamsApi.GetStreamAssetByUuid(ctx, assetID, fabricv4.Asset(asset), streamID).Execute()
			if err != nil {
				//Design decision from API team was to return 400 for all errors instead of 404 for not found
				if resp != nil && slices.Contains([]int{http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return stream, deletedMarker, nil
				}
				return 0, "", err
			}
			return stream, string(stream.GetAttachmentStatus()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
