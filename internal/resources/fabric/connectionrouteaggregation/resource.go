package connectionrouteaggregation

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_connection_route_aggregation",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r Resource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	//No Update Method Supported by Connection Route Aggregation
}

func (r Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema(ctx)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan resourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	routeAggregationID := plan.RouteAggregationID.ValueString()
	connectionID := plan.ConnectionID.ValueString()

	connectionRouteAggregation, _, err := client.RouteAggregationsApi.AttachConnectionRouteAggregation(ctx, routeAggregationID, connectionID).Execute()

	if err != nil {
		resp.Diagnostics.AddError("Failed attaching connection to route aggregation", equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, routeAggregationID, connectionID, createTimeout)
	connectionRouteAggregationChecked, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed attaching Route Aggregation %s", connectionRouteAggregation.GetUuid()), err.Error())
		return
	}

	resp.Diagnostics.Append(plan.parse(ctx, connectionRouteAggregationChecked.(*fabricv4.ConnectionRouteAggregationData))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state resourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	id := state.ID.ValueString()
	routeAggregationID := state.RouteAggregationID.ValueString()
	connectionID := state.ConnectionID.ValueString()

	connectionRouteAggregation, _, err := client.RouteAggregationsApi.GetConnectionRouteAggregationByUuid(ctx, routeAggregationID, connectionID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed retrieving Connection Route Aggregation Attachment %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, connectionRouteAggregation)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	//Retrieve the API client
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	//Retrieve the current state
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	routeAggregationID := state.RouteAggregationID.ValueString()
	connectionID := state.ConnectionID.ValueString()

	_, deleteResp, err := client.RouteAggregationsApi.DetachConnectionRouteAggregation(ctx, routeAggregationID, connectionID).Execute()

	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed detaching Connection Route Aggregation %s", id), equinix_errors.FormatFabricError(err).Error())
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	deletewaiter := getDeleteWaiter(ctx, client, routeAggregationID, connectionID, deleteTimeout)
	_, err = deletewaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed detaching Connection Route Aggregation %s", id), err.Error())
		return
	}
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, routeAggregationID string, connectionID string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHED),
		},
		Refresh: func() (interface{}, string, error) {
			connectionRouteAggregation, _, err := client.RouteAggregationsApi.GetConnectionRouteAggregationByUuid(ctx, routeAggregationID, connectionID).Execute()
			if err != nil {
				return 0, "", err
			}
			return connectionRouteAggregation, string(connectionRouteAggregation.GetAttachmentStatus()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, routeAggregationID string, connectionID string, timeout time.Duration) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the resource appears to be deleted successfully based on
	// status code or specific error code
	deletedMarker := "tf-marker-for-deleted-route-aggregation-rule"
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_DETACHING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_DETACHED),
			deletedMarker,
		},
		Refresh: func() (interface{}, string, error) {
			routeAggregationRule, resp, err := client.RouteAggregationsApi.GetConnectionRouteAggregationByUuid(ctx, routeAggregationID, connectionID).Execute()
			if err != nil {
				if resp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return routeAggregationRule, deletedMarker, nil
				}
				return 0, "", err
			}
			return routeAggregationRule, string(routeAggregationRule.GetAttachmentStatus()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
