package routeaggregationrule

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_route_aggregation_rule",
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

	createRequest, diags := buildCreateRequest(ctx, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	routeAggregationRule, _, err := client.RouteAggregationRulesApi.CreateRouteAggregationRule(ctx, routeAggregationID).RouteAggregationRulesBase(createRequest).Execute()

	if err != nil {
		resp.Diagnostics.AddError("Failed creating route aggregation rule", equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, routeAggregationID, routeAggregationRule.GetUuid(), createTimeout)
	routeAggregationRuleChecked, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed creating Route Aggregation %s", routeAggregationRule.GetUuid()), err.Error())
		return
	}

	resp.Diagnostics.Append(plan.parse(ctx, routeAggregationRuleChecked.(*fabricv4.RouteAggregationRulesData))...)
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

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()
	routeAggregationID := state.RouteAggregationID.ValueString()

	routeAggregationRule, _, err := client.RouteAggregationRulesApi.GetRouteAggregationRuleByUuid(ctx, routeAggregationID, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed retrieving Route Aggregation Rule %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, routeAggregationRule)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	//Retrieve values from plan
	var state, plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := state.ID.ValueString()
	routeAggregationID := state.RouteAggregationID.ValueString()

	newPrefix, oldPrefix := plan.Prefix.ValueString(), state.Prefix.ValueString()

	if newPrefix == oldPrefix {
		resp.Diagnostics.AddWarning("No updatable fields have changed", "Terraform detected a config change, but it is for a field that isn't updatable for the route aggregation rule resource. Please revert to prior config")
		return
	}

	updateRequest := []fabricv4.RouteAggregationRulesPatchRequestItem{{
		Op:    "replace",
		Path:  "/prefix",
		Value: newPrefix,
	}}

	_, _, err := client.RouteAggregationRulesApi.PatchRouteAggregationRuleByUuid(ctx, routeAggregationID, id).RouteAggregationRulesPatchRequestItem(updateRequest).Execute()

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed updating Route Aggregation %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updateWaiter := getCreateUpdateWaiter(ctx, client, routeAggregationID, id, updateTimeout)
	routeAggregationRuleChecked, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed updating Route Aggregation Rule%s", id), err.Error())
		return
	}

	resp.Diagnostics.Append(plan.parse(ctx, routeAggregationRuleChecked.(*fabricv4.RouteAggregationRulesData))...)
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
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	routeAggregationID := state.RouteAggregationID.ValueString()
	_, deleteResp, err := client.RouteAggregationRulesApi.DeleteRouteAggregationRuleByUuid(ctx, routeAggregationID, id).Execute()

	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed deleting Route Aggregation Rule %s", id), equinix_errors.FormatFabricError(err).Error())
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	deletewaiter := getDeleteWaiter(ctx, client, routeAggregationID, id, deleteTimeout)
	_, err = deletewaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed deleting Route Aggregation Rule %s", id), err.Error())
		return
	}
}

func buildCreateRequest(_ context.Context, plan resourceModel) (fabricv4.RouteAggregationRulesBase, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.RouteAggregationRulesBase{}

	request.SetName(plan.Name.ValueString())
	request.SetDescription(plan.Description.ValueString())
	request.SetPrefix(plan.Prefix.ValueString())

	return request, diags
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, routeAggregationID string, routeAggregationRuleID string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.ROUTEAGGREGATIONRULESTATE_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.ROUTEAGGREGATIONRULESTATE_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			routeAggregationRule, _, err := client.RouteAggregationRulesApi.GetRouteAggregationRuleByUuid(ctx, routeAggregationID, routeAggregationRuleID).Execute()
			if err != nil {
				return 0, "", err
			}
			return routeAggregationRule, string(routeAggregationRule.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, routeAggregationID string, id string, timeout time.Duration) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the resource appears to be deleted successfully based on
	// status code or specific error code
	deletedMarker := "tf-marker-for-deleted-route-aggregation-rule"
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.ROUTEAGGREGATIONRULESTATE_DEPROVISIONING),
		},
		Target: []string{
			deletedMarker,
		},
		Refresh: func() (interface{}, string, error) {
			routeAggregationRule, resp, err := client.RouteAggregationRulesApi.GetRouteAggregationRuleByUuid(ctx, routeAggregationID, id).Execute()
			if err != nil {
				if resp != nil {
					if slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
						return routeAggregationRule, deletedMarker, nil
					}
					apiError, ok := err.(*fabricv4.GenericOpenAPIError)
					if ok {
						errorBody := string(apiError.Body())
						if strings.Contains(errorBody, "EQ-3044402") {
							return routeAggregationRule, deletedMarker, nil
						}
					}
				}
				return 0, "", err
			}
			return routeAggregationRule, string(routeAggregationRule.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
