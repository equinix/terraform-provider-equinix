package streamsubscription

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

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

	var plan resourceModel
	diags := req.Plan.Get(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	createRequest, diags := buildCreateRequest(ctx, plan)
	if diags.HasError() {
		return
	}

	streamSubscription, _, err := client.StreamSubscriptionsApi.CreateStreamSubscriptions(ctx, plan.StreamID.ValueString()).StreamSubscriptionPostRequest(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed creating stream subscription", equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, plan.StreamID.ValueString(), streamSubscription.GetUuid(), createTimeout)
	streamChecked, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed creating stream subscription %s", streamSubscription.GetUuid()), err.Error())
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, streamChecked.(*fabricv4.StreamSubscription))...)
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
	streamID := state.StreamID.ValueString()

	streamSubscription, _, err := client.StreamSubscriptionsApi.GetStreamSubscriptionByUuid(ctx, streamID, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed retrieving stream subscription %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, streamSubscription)...)
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
	var state, plan resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	streamID := state.StreamID.ValueString()

	updateRequest, diags := buildUpdateRequest(ctx, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, _, err := client.StreamSubscriptionsApi.UpdateStreamSubscriptionByUuid(ctx, streamID, id).StreamSubscriptionPutRequest(updateRequest).Execute()
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

	updateWaiter := getCreateUpdateWaiter(ctx, client, streamID, id, updateTimeout)
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
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	streamID := state.StreamID.ValueString()

	_, deleteResp, err := client.StreamSubscriptionsApi.DeleteStreamSubscriptionByUuid(ctx, streamID, id).Execute()
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
	deleteWaiter := getDeleteWaiter(ctx, client, streamID, id, deleteTimeout)
	_, err = deleteWaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed deleting Stream %s", id), err.Error())
		return
	}

}

func buildUpdateRequest(ctx context.Context, plan resourceModel) (fabricv4.StreamSubscriptionPutRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.StreamSubscriptionPutRequest{}

	postRequest, diags := buildCreateRequest(ctx, plan)

	request.SetName(postRequest.GetName())
	request.SetDescription(postRequest.GetDescription())

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		request.SetEnabled(postRequest.GetEnabled())
	}

	if !plan.Filters.IsNull() && !plan.Filters.IsUnknown() {
		request.SetFilters(postRequest.GetFilters())
	}

	if !plan.MetricSelector.IsNull() && !plan.MetricSelector.IsUnknown() {
		request.SetMetricSelector(postRequest.GetMetricSelector())
	}

	if !plan.EventSelector.IsNull() && !plan.EventSelector.IsUnknown() {
		request.SetEventSelector(postRequest.GetEventSelector())
	}

	if !plan.Sink.IsNull() && !plan.Sink.IsUnknown() {
		request.SetSink(postRequest.GetSink())
	}

	return request, diags
}

func buildCreateRequest(ctx context.Context, plan resourceModel) (fabricv4.StreamSubscriptionPostRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.StreamSubscriptionPostRequest{}

	request.SetName(plan.Name.ValueString())
	request.SetType(fabricv4.StreamSubscriptionPostRequestType(plan.Type.ValueString()))
	request.SetDescription(plan.Description.ValueString())
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		request.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.Filters.IsNull() && !plan.Filters.IsUnknown() {
		filterModels := make([]filterModel, len(plan.Filters.Elements()))
		diags = plan.Filters.ElementsAs(ctx, &filterModels, false)
		if diags.HasError() {
			return fabricv4.StreamSubscriptionPostRequest{}, diags
		}
		var streamSubscriptionFilter fabricv4.StreamSubscriptionFilter
		var filters []fabricv4.StreamFilter
		var orFilter fabricv4.StreamFilterOrFilter
		for _, filter := range filterModels {
			var expression fabricv4.StreamFilterSimpleExpression
			expression.SetOperator(filter.Operator.ValueString())
			expression.SetProperty(filter.Property.ValueString())
			var values []string
			diags = filter.Values.ElementsAs(ctx, &values, false)
			if diags.HasError() {
				return fabricv4.StreamSubscriptionPostRequest{}, diags
			}
			expression.SetValues(values)
			if filter.Or.ValueBool() {
				orFilter.SetOr(append(orFilter.GetOr(), expression))
			} else {
				filters = append(filters, fabricv4.StreamFilter{
					StreamFilterSimpleExpression: &expression,
				})
			}
		}

		if len(orFilter.GetOr()) > 0 {
			filters = append(filters, fabricv4.StreamFilter{
				StreamFilterOrFilter: &orFilter,
			})
		}
		streamSubscriptionFilter.SetAnd(filters)
		request.SetFilters(streamSubscriptionFilter)
	}

	if !plan.MetricSelector.IsNull() && !plan.MetricSelector.IsUnknown() {
		// Build MetricSelector
		var metricSelector fabricv4.StreamSubscriptionSelector
		metricSelector, diags = buildStreamSubscriptionSelector(ctx, plan.MetricSelector)
		if diags.HasError() {
			return fabricv4.StreamSubscriptionPostRequest{}, diags
		}
		request.SetMetricSelector(metricSelector)
	}

	if !plan.EventSelector.IsNull() && !plan.EventSelector.IsUnknown() {
		// Build EventSelector
		var eventSelector fabricv4.StreamSubscriptionSelector
		eventSelector, diags = buildStreamSubscriptionSelector(ctx, plan.EventSelector)
		if diags.HasError() {
			return fabricv4.StreamSubscriptionPostRequest{}, diags
		}
		request.SetEventSelector(eventSelector)
	}

	if !plan.Sink.IsNull() && !plan.Sink.IsUnknown() {
		// Update sink request
		var sink fabricv4.StreamSubscriptionSink
		sink, diags = buildStreamSubscriptionSink(ctx, plan.Sink)
		if diags.HasError() {
			return fabricv4.StreamSubscriptionPostRequest{}, diags
		}
		request.SetSink(sink)
	}

	return request, diags
}

func buildStreamSubscriptionSelector(ctx context.Context, selector fwtypes.ObjectValueOf[selectorModel]) (fabricv4.StreamSubscriptionSelector, diag.Diagnostics) {
	var selectorValue selectorModel
	diags := selector.As(ctx, &selectorValue, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.StreamSubscriptionSelector{}, diags
	}

	var streamSubscriptionSelector fabricv4.StreamSubscriptionSelector
	if !selectorValue.Include.IsNull() && !selectorValue.Include.IsUnknown() {
		include := []string{}
		diags = selectorValue.Include.ElementsAs(ctx, &include, false)
		if diags.HasError() {
			return fabricv4.StreamSubscriptionSelector{}, diags
		}
		streamSubscriptionSelector.SetInclude(include)
	}

	if !selectorValue.Except.IsNull() && !selectorValue.Except.IsUnknown() {
		except := []string{}
		diags = selectorValue.Except.ElementsAs(ctx, &except, false)
		if diags.HasError() {
			return fabricv4.StreamSubscriptionSelector{}, diags
		}
		streamSubscriptionSelector.SetExcept(except)
	}

	return streamSubscriptionSelector, diags
}

func buildStreamSubscriptionSink(ctx context.Context, sinkObject fwtypes.ObjectValueOf[sinkModel]) (fabricv4.StreamSubscriptionSink, diag.Diagnostics) {
	var sinkValue sinkModel
	diags := sinkObject.As(ctx, &sinkValue, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.StreamSubscriptionSink{}, diags
	}

	var sink fabricv4.StreamSubscriptionSink

	sink.SetType(fabricv4.StreamSubscriptionSinkType(sinkValue.Type.ValueString()))
	if !sinkValue.URI.IsNull() && !sinkValue.URI.IsUnknown() {
		sink.SetUri(sinkValue.URI.ValueString())
	}

	if !sinkValue.BatchEnabled.IsNull() && !sinkValue.BatchEnabled.IsUnknown() {
		sink.SetBatchEnabled(sinkValue.BatchEnabled.ValueBool())
	}

	if !sinkValue.BatchSizeMax.IsNull() && !sinkValue.BatchSizeMax.IsUnknown() {
		sink.SetBatchSizeMax(sinkValue.BatchSizeMax.ValueInt32())
	}

	if !sinkValue.BatchWaitTimeMax.IsNull() && !sinkValue.BatchWaitTimeMax.IsUnknown() {
		sink.SetBatchWaitTimeMax(sinkValue.BatchWaitTimeMax.ValueInt32())
	}

	if !sinkValue.Host.IsNull() && !sinkValue.Host.IsUnknown() {
		sink.SetHost(sinkValue.Host.ValueString())
	}

	if !sinkValue.Credential.IsNull() && !sinkValue.Credential.IsUnknown() {
		// Build credential
		var credentialModel sinkCredentialModel
		diags = sinkValue.Credential.As(ctx, &credentialModel, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fabricv4.StreamSubscriptionSink{}, diags
		}

		var credential fabricv4.StreamSubscriptionSinkCredential
		credential.SetType(fabricv4.StreamSubscriptionSinkCredentialType(credentialModel.Type.ValueString()))

		switch credential.GetType() {
		case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_ACCESS_TOKEN:
			credential.SetAccessToken(credentialModel.AccessToken.ValueString())
		case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_INTEGRATION_KEY:
			credential.SetIntegrationKey(credentialModel.IntegrationKey.ValueString())
		case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_API_KEY:
			credential.SetApiKey(credentialModel.APIKey.ValueString())
		case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_USERNAME_PASSWORD:
			credential.SetUsername(credentialModel.Username.ValueString())
			credential.SetPassword(credentialModel.Password.ValueString())
		default:
			diags.AddError("sink credential type is invalid", fmt.Sprintf("update sink credential type to one of the following %v", fabricv4.AllowedStreamSubscriptionSinkCredentialTypeEnumValues))
			return fabricv4.StreamSubscriptionSink{}, diags
		}

		sink.SetCredential(credential)
	}

	if !sinkValue.Settings.IsNull() && !sinkValue.Settings.IsUnknown() {
		// Build settings
		var settingsModel sinkSettingsModel
		diags = sinkValue.Settings.As(ctx, &settingsModel, basetypes.ObjectAsOptions{})

		var settings fabricv4.StreamSubscriptionSinkSetting
		if !settingsModel.EventIndex.IsNull() && !settingsModel.EventIndex.IsUnknown() {
			settings.SetEventIndex(settingsModel.EventIndex.ValueString())
		}
		if !settingsModel.MetricIndex.IsNull() && !settingsModel.MetricIndex.IsUnknown() {
			settings.SetMetricIndex(settingsModel.MetricIndex.ValueString())
		}
		if !settingsModel.Source.IsNull() && !settingsModel.Source.IsUnknown() {
			settings.SetSource(settingsModel.Source.ValueString())
		}
		if !settingsModel.ApplicationKey.IsNull() && !settingsModel.ApplicationKey.IsUnknown() {
			settings.SetApplicationKey(settingsModel.ApplicationKey.ValueString())
		}
		if !settingsModel.EventURI.IsNull() && !settingsModel.EventURI.IsUnknown() {
			settings.SetEventUri(settingsModel.EventURI.ValueString())
		}
		if !settingsModel.MetricURI.IsNull() && !settingsModel.MetricURI.IsUnknown() {
			settings.SetMetricUri(settingsModel.MetricURI.ValueString())
		}
		sink.SetSettings(settings)
	}

	return sink, diags
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, streamID, streamSubscriptionID string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONING),
			"REPROVISIONING",
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
	deletedMarker := "tf-marker-for-deleted-stream-subscription"
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
