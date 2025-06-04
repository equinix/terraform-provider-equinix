package stream_alert_rule

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"log"
	"net/http"
	"slices"
	"time"
)

// NewResource creates ba new stream alert rule
func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_stream_alert_rule",
			},
		),
	}
}

// Resource represents the stream alert rule
type Resource struct {
	framework.BaseResource
}

// Schema returns the resource schema
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

	var plan resourceModel
	diags := req.Plan.Get(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)
	alertRulePostRequest, diags := buildCreateRequest(ctx, plan)
	log.Println("deep1" + alertRulePostRequest.GetName())

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	log.Println("deep2" + plan.StreamID.ValueString())

	streamAlertRule, _, err := client.StreamAlertRulesApi.CreateStreamAlertRules(ctx, plan.StreamID.ValueString()).AlertRulePostRequest(alertRulePostRequest).Execute()
	alertRuleUuid := streamAlertRule.GetUuid()
	plan.ID = types.StringValue(alertRuleUuid)
	if err != nil {
		resp.Diagnostics.AddError("failed creating stream alert rule", equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, plan.StreamID.ValueString(), streamAlertRule.GetUuid(), createTimeout)
	alertRuleChecked, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed creating stream alert rule %s", streamAlertRule.GetUuid()), err.Error())
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, alertRuleChecked.(*fabricv4.StreamAlertRule))...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, alertRuleID, streamAlertRuleID string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.STREAMALERTRULESTATE_INACTIVE),
			"INACTIVE",
		},
		Target: []string{
			string(fabricv4.STREAMALERTRULESTATE_ACTIVE),
		},
		Refresh: func() (interface{}, string, error) {
			streamAlertRule, _, err := client.StreamAlertRulesApi.GetStreamAlertRuleByUuid(ctx, alertRuleID, streamAlertRuleID).Execute()
			if err != nil {
				return 0, "", err
			}
			return streamAlertRule, string(streamAlertRule.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}

func buildCreateRequest(ctx context.Context, plan resourceModel) (fabricv4.AlertRulePostRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.AlertRulePostRequest{}

	request.SetName(plan.Name.ValueString())
	request.SetType(fabricv4.AlertRulePostRequestType(plan.Type.ValueString()))
	request.SetDescription(plan.Description.ValueString())
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		request.SetEnabled(plan.Enabled.ValueBool())
	}
	request.SetWarningThreshold(plan.WarningThreshold.ValueString())
	request.SetCriticalThreshold(plan.CriticalThreshold.ValueString())
	request.SetMetricName(fabricv4.StreamAlertRuleMetricName(plan.MetricName.ValueString()))
	request.SetWindowSize(plan.WindowSize.ValueString())
	request.SetOperand(fabricv4.StreamAlertRuleOperand(plan.Operand.ValueString()))

	if !plan.ResourceSelector.IsNull() && !plan.ResourceSelector.IsUnknown() {
		// Build ResourceSelector
		var resourceSelector fabricv4.ResourceSelector
		resourceSelector, diags = buildStreamAlertRuleSelector(ctx, plan.ResourceSelector)
		if diags.HasError() {
			return fabricv4.AlertRulePostRequest{}, diags
		}
		request.SetResourceSelector(resourceSelector)
	}
	return request, diags
}

func buildStreamAlertRuleSelector(ctx context.Context, selector fwtypes.ObjectValueOf[selectorModel]) (fabricv4.ResourceSelector, diag.Diagnostics) {
	var selectorValue selectorModel
	diags := selector.As(ctx, &selectorValue, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.ResourceSelector{}, diags
	}

	var resourceSelector fabricv4.ResourceSelector
	if !selectorValue.Include.IsNull() && !selectorValue.Include.IsUnknown() {
		include := []string{}
		diags = selectorValue.Include.ElementsAs(ctx, &include, false)
		if diags.HasError() {
			return fabricv4.ResourceSelector{}, diags
		}
		resourceSelector.SetInclude(include)
	}
	return resourceSelector, diags
}

// Read retrieves a new stream alert rule
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

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()
	streamAlertRuleID := state.StreamID.ValueString()

	streamAlertRule, _, err := client.StreamAlertRulesApi.GetStreamAlertRuleByUuid(ctx, streamAlertRuleID, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed retrieving stream alert rule %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, streamAlertRule)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies an existing stream alert rule
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Retrieve values from plan
	var state, plan resourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	plan.StreamID = state.StreamID
	id := state.ID.ValueString()
	streamID := state.StreamID.ValueString()
	updateRequest, diags := buildUpdateRequest(ctx, plan)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, _, err := client.StreamAlertRulesApi.UpdateStreamAlertRuleByUuid(ctx, streamID, id).AlertRulePutRequest(updateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed updating stream alert rule %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updateWaiter := getCreateUpdateWaiter(ctx, client, streamID, id, updateTimeout)
	streamAlertRuleChecked, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed updating stream subscription %s", id), err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, streamAlertRuleChecked.(*fabricv4.StreamAlertRule))...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueString() == "" {
		plan.ID = types.StringValue(id)
	}

	// Set the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func buildUpdateRequest(ctx context.Context, plan resourceModel) (fabricv4.AlertRulePutRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.AlertRulePutRequest{}

	postRequest, diags := buildCreateRequest(ctx, plan)

	request.SetName(postRequest.GetName())
	request.SetDescription(postRequest.GetDescription())
	request.SetWarningThreshold(postRequest.GetWarningThreshold())
	request.SetCriticalThreshold(postRequest.GetCriticalThreshold())
	request.SetWindowSize(postRequest.GetWindowSize())
	request.SetMetricName(postRequest.GetMetricName())
	request.SetOperand(postRequest.GetOperand())

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		request.SetEnabled(postRequest.GetEnabled())
	}

	if !plan.ResourceSelector.IsNull() && !plan.ResourceSelector.IsUnknown() {
		request.SetResourceSelector(postRequest.GetResourceSelector())
	}

	return request, diags
}

// Delete removes the stream alert rule
func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve the API client
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Retrieve the current state
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	streamAlertRuleID := state.StreamID.ValueString()

	_, deleteResp, err := client.StreamAlertRulesApi.DeleteStreamAlertRuleByUuid(ctx, streamAlertRuleID, id).Execute()
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
	deleteWaiter := getDeleteWaiter(ctx, client, streamAlertRuleID, id, deleteTimeout)
	_, err = deleteWaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed deleting Stream %s", id), err.Error())
		return
	}

}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, streamID, streamAlertRuleID string, timeout time.Duration) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the connection appears to be deleted successfully based on
	// status code
	deletedMarker := "tf-marker-for-deleted-stream-alert-rule"
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.STREAMALERTRULESTATE_INACTIVE),
		},
		Target: []string{
			deletedMarker,
			string(fabricv4.STREAMALERTRULESTATE_ACTIVE),
		},
		Refresh: func() (interface{}, string, error) {
			streamAlertRule, resp, err := client.StreamAlertRulesApi.GetStreamAlertRuleByUuid(ctx, streamID, streamAlertRuleID).Execute()
			if err != nil {
				if resp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return streamAlertRule, deletedMarker, nil
				}
				return 0, "", err
			}
			return streamAlertRule, string(streamAlertRule.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
