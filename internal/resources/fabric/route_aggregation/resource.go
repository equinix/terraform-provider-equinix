package route_aggregation

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"net/http"
	"slices"
	"strings"
	"time"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_route_aggregation",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	createRequest, diags := buildCreateRequest(ctx, plan)
	if diags.HasError() {
		return
	}

	routeAggregation, _, err := client.RouteAggregationsApi.CreateRouteAggregation(ctx).RouteAggregationsBase(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed creating route aggregation"), equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, routeAggregation.GetUuid(), createTimeout)
	routeAggregationChecked, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed creating Route Aggregation %s", routeAggregation.GetUuid()), err.Error())
		return
	}

	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.parse(ctx, routeAggregationChecked.(*fabricv4.RouteAggregationsData))...)
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
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	routeAggregation, _, err := client.RouteAggregationsApi.GetRouteAggregationByUuid(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed retrieving Route Aggregation %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, routeAggregation)...)
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

	//Retrieve values from plan
	var state, plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	newName, oldName := plan.Name.ValueString(), plan.Name.ValueString()

	if newName != oldName {
		resp.Diagnostics.AddWarning("No updatable fields have changed", "Terraform detected a config change, but it is for a field that isn't updatable for the route aggregation resource. Please revert to prior config")
		return
	}

	updateRequest := []fabricv4.RouteAggregationsPatchRequestItem{{
		Op:    "replace",
		Path:  "/name",
		Value: map[string]interface{}{"": newName},
	},
	}

	_, _, err := client.RouteAggregationsApi.PatchRouteAggregationByUuid(ctx, id).RouteAggregationsPatchRequestItem(updateRequest).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed updating Route Aggregation %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updateWaiter := getCreateUpdateWaiter(ctx, client, id, updateTimeout)
	routeAggregationChecked, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed updating Route Aggregation %s", id), err.Error())
		return
	}

	//set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, routeAggregationChecked.(*fabricv4.RouteAggregationsData))...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Set the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	//Retrieve the API client
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	//Retrieve the current state
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, deleteResp, err := client.RouteAggregationsApi.DeleteRouteAggregationByUuid(ctx, id).Execute()
	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed deleting Stream %s", id), equinix_errors.FormatFabricError(err).Error())
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	deletewaiter := getDeleteWaiter(ctx, client, id, deleteTimeout)
	_, err = deletewaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed deleting Stream %s", id), err.Error())
		return
	}
}

func buildCreateRequest(ctx context.Context, plan ResourceModel) (fabricv4.RouteAggregationsBase, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.RouteAggregationsBase{}

	request.SetType(fabricv4.RouteAggregationsBaseType(plan.Type.ValueString()))
	request.SetName(plan.Name.ValueString())
	request.SetDescription(plan.Description.ValueString())

	var project ProjectModel
	if !plan.Project.IsNull() && !plan.Project.IsUnknown() {
		diags = plan.Project.As(ctx, &project, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fabricv4.RouteAggregationsBase{}, diags
		}
		request.SetProject(fabricv4.Project{ProjectId: project.ProjectId.ValueString()})
	}
	return request, diags
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.ROUTEAGGREGATIONSTATE_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.ROUTEAGGREGATIONSTATE_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			routeAggregation, _, err := client.RouteAggregationsApi.GetRouteAggregationByUuid(ctx, id).Execute()
			if err != nil {
				return 0, "", err
			}
			return routeAggregation, string(routeAggregation.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the resource appears to be deleted successfully based on
	// status code or specific error code
	deletedMarker := "tf-marker-for-deleted-route-aggregation"
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.ROUTEAGGREGATIONSTATE_DEPROVISIONING),
		},
		Target: []string{
			deletedMarker,
		},
		Refresh: func() (interface{}, string, error) {
			routeAggregation, resp, err := client.RouteAggregationsApi.GetRouteAggregationByUuid(ctx, id).Execute()
			if err != nil {
				if resp != nil {
					if slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
						return routeAggregation, deletedMarker, nil
					}
					apiError, ok := err.(*fabricv4.GenericOpenAPIError)
					if ok {
						errorBody := string(apiError.Body())
						if strings.Contains(errorBody, "EQ-3044301") {
							return routeAggregation, deletedMarker, nil
						}
					}
				}
				return 0, "", err
			}
			return routeAggregation, string(routeAggregation.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
