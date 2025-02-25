package streamsubscription

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/equinix/terraform-provider-equinix/internal/fabric"
	int_fw "github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceByIDsModel struct {
	ID             types.String `tfsdk:"id"`
	StreamID       types.String `tfsdk:"stream_id"`
	SubscriptionID types.String `tfsdk:"subscription_id"`
	baseStreamSubscriptionModel
}

type dataSourceAll struct {
	ID         types.String                                                 `tfsdk:"id"`
	StreamID   types.String                                                 `tfsdk:"stream_id"`
	Pagination fwtypes.ObjectValueOf[paginationModel]                       `tfsdk:"pagination"`
	Data       fwtypes.ListNestedObjectValueOf[baseStreamSubscriptionModel] `tfsdk:"data"`
}

type paginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type resourceModel struct {
	StreamID types.String   `tfsdk:"stream_id"`
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	baseStreamSubscriptionModel
}

type baseStreamSubscriptionModel struct {
	Type           types.String                                 `tfsdk:"type"`
	Name           types.String                                 `tfsdk:"name"`
	Description    types.String                                 `tfsdk:"description"`
	Enabled        types.Bool                                   `tfsdk:"enabled"`
	Filters        fwtypes.ListNestedObjectValueOf[filterModel] `tfsdk:"filters"`         // List of filters
	MetricSelector fwtypes.ObjectValueOf[selectorModel]         `tfsdk:"metric_selector"` // Object of MetricSelectorModel
	EventSelector  fwtypes.ObjectValueOf[selectorModel]         `tfsdk:"event_selector"`  // Object of EventSelectorModel
	Sink           fwtypes.ObjectValueOf[sinkModel]             `tfsdk:"sink"`            // Object of SinkModel
	Href           types.String                                 `tfsdk:"href"`
	UUID           types.String                                 `tfsdk:"uuid"`
	State          types.String                                 `tfsdk:"state"`
	ChangeLog      fwtypes.ObjectValueOf[changeLogModel]        `tfsdk:"change_log"` // Object of ChangeLogModel
}

type filterModel struct {
	Property types.String                      `tfsdk:"property"`
	Operator types.String                      `tfsdk:"operator"`
	Values   fwtypes.ListValueOf[types.String] `tfsdk:"values"`
	Or       types.Bool                        `tfsdk:"or"`
}

type selectorModel struct {
	Include fwtypes.ListValueOf[types.String] `tfsdk:"include"`
	Except  fwtypes.ListValueOf[types.String] `tfsdk:"except"`
}

type sinkModel struct {
	URI              types.String                               `tfsdk:"uri"`
	Type             types.String                               `tfsdk:"type"`
	BatchEnabled     types.Bool                                 `tfsdk:"batch_enabled"`
	BatchSizeMax     types.Int32                                `tfsdk:"batch_size_max"`
	BatchWaitTimeMax types.Int32                                `tfsdk:"batch_wait_time_max"`
	Host             types.String                               `tfsdk:"host"`
	Credential       fwtypes.ObjectValueOf[sinkCredentialModel] `tfsdk:"credential"` // Object of CredentialModel
	Settings         fwtypes.ObjectValueOf[sinkSettingsModel]   `tfsdk:"settings"`   // Object of SinkSettingsModel
}

type sinkCredentialModel struct {
	Type           types.String `tfsdk:"type"`
	AccessToken    types.String `tfsdk:"access_token"`
	IntegrationKey types.String `tfsdk:"integration_key"`
	APIKey         types.String `tfsdk:"api_key"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
}

type sinkSettingsModel struct {
	EventIndex     types.String `tfsdk:"event_index"`
	MetricIndex    types.String `tfsdk:"metric_index"`
	Source         types.String `tfsdk:"source"`
	ApplicationKey types.String `tfsdk:"application_key"`
	EventURI       types.String `tfsdk:"event_uri"`
	MetricURI      types.String `tfsdk:"metric_uri"`
}

type changeLogModel struct {
	CreatedBy         types.String `tfsdk:"created_by"`
	CreatedByFullName types.String `tfsdk:"created_by_full_name"`
	CreatedByEmail    types.String `tfsdk:"created_by_email"`
	CreatedDateTime   types.String `tfsdk:"created_date_time"`
	UpdatedBy         types.String `tfsdk:"updated_by"`
	UpdatedByFullName types.String `tfsdk:"updated_by_full_name"`
	UpdatedByEmail    types.String `tfsdk:"updated_by_email"`
	UpdatedDateTime   types.String `tfsdk:"updated_date_time"`
	DeletedBy         types.String `tfsdk:"deleted_by"`
	DeletedByFullName types.String `tfsdk:"deleted_by_full_name"`
	DeletedByEmail    types.String `tfsdk:"deleted_by_email"`
	DeletedDateTime   types.String `tfsdk:"deleted_date_time"`
}

func (m *dataSourceByIDsModel) parse(ctx context.Context, streamSubscription *fabricv4.StreamSubscription) diag.Diagnostics {
	m.StreamID = types.StringValue(streamSubscription.GetUuid())
	m.SubscriptionID = types.StringValue(streamSubscription.GetUuid())
	m.ID = types.StringValue(streamSubscription.GetUuid())

	diags := m.baseStreamSubscriptionModel.parse(ctx, streamSubscription)
	if diags.HasError() {
		return diags
	}

	return diags
}

func (m *dataSourceAll) parse(ctx context.Context, streamSubscriptionsResponse *fabricv4.GetAllStreamSubscriptionResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(streamSubscriptionsResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by stream subscriptions data source",
			"either the account does not have any stream subscription data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]baseStreamSubscriptionModel, len(streamSubscriptionsResponse.GetData()))
	streamSubscriptions := streamSubscriptionsResponse.GetData()
	for index, streamSubscription := range streamSubscriptions {
		var streamSubscriptionModel baseStreamSubscriptionModel
		diags = streamSubscriptionModel.parse(ctx, &streamSubscription)
		if diags.HasError() {
			return diags
		}
		data[index] = streamSubscriptionModel
	}
	responsePagination := streamSubscriptionsResponse.GetPagination()
	pagination := paginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}

	m.ID = types.StringValue(data[0].UUID.ValueString())
	m.StreamID = types.StringValue(data[0].UUID.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[paginationModel](ctx, &pagination)
	m.Data = fwtypes.NewListNestedObjectValueOfValueSlice[baseStreamSubscriptionModel](ctx, data)

	return diags
}

func (m *resourceModel) parse(ctx context.Context, streamSubscription *fabricv4.StreamSubscription) diag.Diagnostics {
	m.ID = types.StringValue(streamSubscription.GetUuid())

	diags := m.baseStreamSubscriptionModel.parse(ctx, streamSubscription)
	if diags.HasError() {
		return diags
	}

	return diags
}

func (m *baseStreamSubscriptionModel) parse(ctx context.Context, streamSubscription *fabricv4.StreamSubscription) diag.Diagnostics {

	var mDiags diag.Diagnostics

	m.Type = types.StringValue(string(streamSubscription.GetType()))
	m.Name = types.StringValue(streamSubscription.GetName())
	m.Description = types.StringValue(streamSubscription.GetDescription())
	m.Href = types.StringValue(streamSubscription.GetHref())
	m.UUID = types.StringValue(streamSubscription.GetUuid())
	m.State = types.StringValue(string(streamSubscription.GetState()))
	m.Enabled = types.BoolValue(streamSubscription.GetEnabled())

	// Parse filters
	streamSubscriptionFilters := streamSubscription.GetFilters()
	filterModels := make([]filterModel, len(streamSubscriptionFilters.GetAnd()))
	for i, filter := range streamSubscriptionFilters.GetAnd() {
		if len(filter.StreamFilterOrFilter.GetOr()) > 0 {
			for j, orFilter := range filter.StreamFilterOrFilter.GetOr() {
				orFilterModel, diags := parseSimpleExpression(ctx, &orFilter, true)
				if diags.HasError() {
					mDiags.Append(diags...)
					return mDiags
				}
				// If the first OrGroup selector assign it to the space made in the slice
				// Else append it to the end.
				// We do this because we can't do the exact representation of the API model
				// and this will be a longer list with orGroup boolean instead of a sub list
				if j == 0 {
					filterModels[i] = orFilterModel
				} else {
					filterModels = append(filterModels, orFilterModel)
				}
			}
		} else {
			// The unmarshal for this will always put the values in the additional properties for the
			// StreamFilterOrFilter because it checks for that embedded struct first and
			// the unmarshal allows it to do so without error; so it will never proceed to the
			// StreamFilterSimpleExpression struct. So if GetOr doesn't have any values we check the additional
			// properties map of StreamFilterOrFilter instead. Something to address at API Spec level
			// before code generation of equinix-sdk-go/fabricv4 for long term fix
			values := int_fw.StringSliceToAttrValue(converters.IfArrToStringArr(filter.StreamFilterOrFilter.AdditionalProperties["values"].([]interface{})))
			fwValues, diags := fwtypes.NewListValueOf[types.String](ctx, values)
			if diags.HasError() {
				mDiags.Append(diags...)
				return mDiags
			}
			filterModels[i] = filterModel{
				Property: types.StringValue(filter.StreamFilterOrFilter.AdditionalProperties["property"].(string)),
				Operator: types.StringValue(filter.StreamFilterOrFilter.AdditionalProperties["operator"].(string)),
				Values:   fwValues,
				Or:       types.BoolValue(false),
			}
		}
	}
	m.Filters = fwtypes.NewListNestedObjectValueOfValueSlice[filterModel](ctx, filterModels)

	// Parse MetricSelector
	metricSelectorObject, diags := parseSelectorModel(ctx, streamSubscription.GetMetricSelector())
	if diags.HasError() {
		mDiags.Append(diags...)
		return mDiags
	}
	m.MetricSelector = metricSelectorObject

	// Parse EventSelector
	eventSelectorObject, diags := parseSelectorModel(ctx, streamSubscription.GetEventSelector())
	if diags.HasError() {
		mDiags.Append(diags...)
		return mDiags
	}
	m.EventSelector = eventSelectorObject

	planSinkModel := sinkModel{}
	if !m.Sink.IsNull() && !m.Sink.IsUnknown() {
		diags = m.Sink.As(ctx, &planSinkModel, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return mDiags
		}
	}

	// Parse Sink
	streamSubSink := streamSubscription.GetSink()
	sink := sinkModel{
		URI:              types.StringValue(streamSubSink.GetUri()),
		Type:             types.StringValue(string(streamSubSink.GetType())),
		BatchEnabled:     types.BoolValue(streamSubSink.GetBatchEnabled()),
		BatchSizeMax:     types.Int32Value(streamSubSink.GetBatchSizeMax()),
		BatchWaitTimeMax: types.Int32Value(streamSubSink.GetBatchWaitTimeMax()),
		Host:             types.StringValue(streamSubSink.GetHost()),
	}

	if planSinkModel.URI.ValueString() != "" {
		sink.URI = types.StringValue(planSinkModel.URI.ValueString())
	}

	sinkCredential := streamSubSink.GetCredential()
	credentialModel := sinkCredentialModel{
		Type:           types.StringValue(string(sinkCredential.GetType())),
		AccessToken:    types.StringValue(sinkCredential.GetAccessToken()),
		IntegrationKey: types.StringValue(sinkCredential.GetIntegrationKey()),
		APIKey:         types.StringValue(sinkCredential.GetApiKey()),
		Username:       types.StringValue(sinkCredential.GetUsername()),
		Password:       types.StringValue(sinkCredential.GetPassword()),
	}

	if !planSinkModel.Credential.IsNull() && !planSinkModel.Credential.IsUnknown() {
		planCredentialModel := sinkCredentialModel{}
		diags = planSinkModel.Credential.As(ctx, &planCredentialModel, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return mDiags
		}
		switch fabricv4.StreamSubscriptionSinkCredentialType(planCredentialModel.Type.ValueString()) {
		case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_ACCESS_TOKEN:
			credentialModel.AccessToken = types.StringValue(planCredentialModel.AccessToken.ValueString())
		case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_API_KEY:
			credentialModel.APIKey = types.StringValue(planCredentialModel.APIKey.ValueString())
		case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_INTEGRATION_KEY:
			credentialModel.IntegrationKey = types.StringValue(planCredentialModel.IntegrationKey.ValueString())
		case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_USERNAME_PASSWORD:
			credentialModel.Username = types.StringValue(planCredentialModel.Username.ValueString())
			credentialModel.Password = types.StringValue(planCredentialModel.Password.ValueString())
		}
	}

	sink.Credential = fwtypes.NewObjectValueOf[sinkCredentialModel](ctx, &credentialModel)

	streamSubSinkSettings := streamSubSink.GetSettings()
	sinkSettings := sinkSettingsModel{
		EventIndex:     types.StringValue(streamSubSinkSettings.GetEventIndex()),
		MetricIndex:    types.StringValue(streamSubSinkSettings.GetMetricIndex()),
		Source:         types.StringValue(streamSubSinkSettings.GetSource()),
		ApplicationKey: types.StringValue(streamSubSinkSettings.GetApplicationKey()),
		EventURI:       types.StringValue(streamSubSinkSettings.GetEventUri()),
		MetricURI:      types.StringValue(streamSubSinkSettings.GetMetricUri()),
	}

	if !planSinkModel.Settings.IsNull() && !planSinkModel.Settings.IsUnknown() {
		planSettingsModel := sinkSettingsModel{}
		diags = planSinkModel.Settings.As(ctx, &planSettingsModel, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return mDiags
		}
		if planSettingsModel.ApplicationKey.ValueString() != "" {
			sinkSettings.ApplicationKey = types.StringValue(planSettingsModel.ApplicationKey.ValueString())
		}
	}

	sink.Settings = fwtypes.NewObjectValueOf[sinkSettingsModel](ctx, &sinkSettings)

	m.Sink = fwtypes.NewObjectValueOf[sinkModel](ctx, &sink)

	// Parse ChangeLog
	streamSubscriptionChangeLog := streamSubscription.GetChangeLog()
	changeLog := changeLogModel{
		CreatedBy:         types.StringValue(streamSubscriptionChangeLog.GetCreatedBy()),
		CreatedByFullName: types.StringValue(streamSubscriptionChangeLog.GetCreatedByFullName()),
		CreatedByEmail:    types.StringValue(streamSubscriptionChangeLog.GetCreatedByEmail()),
		CreatedDateTime:   types.StringValue(streamSubscriptionChangeLog.GetCreatedDateTime().Format(fabric.TimeFormat)),
		UpdatedBy:         types.StringValue(streamSubscriptionChangeLog.GetUpdatedBy()),
		UpdatedByFullName: types.StringValue(streamSubscriptionChangeLog.GetUpdatedByFullName()),
		UpdatedByEmail:    types.StringValue(streamSubscriptionChangeLog.GetUpdatedByEmail()),
		UpdatedDateTime:   types.StringValue(streamSubscriptionChangeLog.GetUpdatedDateTime().Format(fabric.TimeFormat)),
		DeletedBy:         types.StringValue(streamSubscriptionChangeLog.GetDeletedBy()),
		DeletedByFullName: types.StringValue(streamSubscriptionChangeLog.GetDeletedByFullName()),
		DeletedByEmail:    types.StringValue(streamSubscriptionChangeLog.GetDeletedByEmail()),
		DeletedDateTime:   types.StringValue(streamSubscriptionChangeLog.GetDeletedDateTime().Format(fabric.TimeFormat)),
	}
	m.ChangeLog = fwtypes.NewObjectValueOf[changeLogModel](ctx, &changeLog)

	return mDiags
}

func parseSimpleExpression(ctx context.Context, expression *fabricv4.StreamFilterSimpleExpression, orGroup bool) (filterModel, diag.Diagnostics) {
	values := int_fw.StringSliceToAttrValue(expression.GetValues())
	fwValues, diags := fwtypes.NewListValueOf[types.String](ctx, values)
	if diags.HasError() {
		return filterModel{}, diags
	}
	return filterModel{
		Property: types.StringValue(expression.GetProperty()),
		Operator: types.StringValue(expression.GetOperator()),
		Values:   fwValues,
		Or:       types.BoolValue(orGroup),
	}, nil
}

func parseSelectorModel(ctx context.Context, streamSubSelector fabricv4.StreamSubscriptionSelector) (fwtypes.ObjectValueOf[selectorModel], diag.Diagnostics) {
	var diags diag.Diagnostics
	inclusions, diags := fwtypes.NewListValueOf[types.String](ctx, int_fw.StringSliceToAttrValue(streamSubSelector.GetInclude()))
	if diags.HasError() {
		return fwtypes.NewObjectValueOfNull[selectorModel](ctx), diags
	}
	exclusions, diags := fwtypes.NewListValueOf[types.String](ctx, int_fw.StringSliceToAttrValue(streamSubSelector.GetExcept()))
	if diags.HasError() {
		return fwtypes.NewObjectValueOfNull[selectorModel](ctx), diags
	}
	selector := selectorModel{
		Include: inclusions,
		Except:  exclusions,
	}
	return fwtypes.NewObjectValueOf[selectorModel](ctx, &selector), diags
}
