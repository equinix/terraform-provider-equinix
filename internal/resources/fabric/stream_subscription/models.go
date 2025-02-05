package stream_subscription

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/fabric"
	int_fw "github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DataSourceByIDsModel struct {
	ID             types.String `tfsdk:"id"`
	StreamID       types.String `tfsdk:"stream_id"`
	SubscriptionID types.String `tfsdk:"subscription_id"`
	BaseStreamSubscriptionModel
}

type DataSourceAll struct {
	ID         types.String                                                 `tfsdk:"id"`
	StreamID   types.String                                                 `tfsdk:"stream_id"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]                       `tfsdk:"pagination"`
	Data       fwtypes.ListNestedObjectValueOf[BaseStreamSubscriptionModel] `tfsdk:"data"`
}

type PaginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type ResourceModel struct {
	StreamID types.String   `tfsdk:"stream_id"`
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	BaseStreamSubscriptionModel
}

type BaseStreamSubscriptionModel struct {
	Type           types.String                                 `tfsdk:"type"`
	Name           types.String                                 `tfsdk:"name"`
	Description    types.String                                 `tfsdk:"description"`
	Enabled        types.Bool                                   `tfsdk:"enabled"`
	Filters        fwtypes.ListNestedObjectValueOf[FilterModel] `tfsdk:"filters"`         // List of filters
	MetricSelector fwtypes.ObjectValueOf[SelectorModel]         `tfsdk:"metric_selector"` // Object of MetricSelectorModel
	EventSelector  fwtypes.ObjectValueOf[SelectorModel]         `tfsdk:"event_selector"`  // Object of EventSelectorModel
	Sink           fwtypes.ObjectValueOf[SinkModel]             `tfsdk:"sink"`            // Object of SinkModel
	Href           types.String                                 `tfsdk:"href"`
	Uuid           types.String                                 `tfsdk:"uuid"`
	State          types.String                                 `tfsdk:"state"`
	ChangeLog      fwtypes.ObjectValueOf[ChangeLogModel]        `tfsdk:"change_log"` // Object of ChangeLogModel
}

type FilterModel struct {
	Property types.String                      `tfsdk:"property"`
	Operator types.String                      `tfsdk:"operator"`
	Values   fwtypes.ListValueOf[types.String] `tfsdk:"values"`
	Or       types.Bool                        `tfsdk:"or"`
}

type SelectorModel struct {
	Include fwtypes.ListValueOf[types.String] `tfsdk:"include"`
	Except  fwtypes.ListValueOf[types.String] `tfsdk:"except"`
}

type SinkModel struct {
	Uri              types.String                               `tfsdk:"uri"`
	Type             types.String                               `tfsdk:"type"`
	BatchEnabled     types.Bool                                 `tfsdk:"batch_enabled"`
	BatchSizeMax     types.Int32                                `tfsdk:"batch_size_max"`
	BatchWaitTimeMax types.Int32                                `tfsdk:"batch_wait_time_max"`
	Host             types.String                               `tfsdk:"host"`
	Credential       fwtypes.ObjectValueOf[SinkCredentialModel] `tfsdk:"credential"` // Object of CredentialModel
	Settings         fwtypes.ObjectValueOf[SinkSettingsModel]   `tfsdk:"settings"`   // Object of SinkSettingsModel
}

type SinkCredentialModel struct {
	Type           types.String `tfsdk:"type"`
	AccessToken    types.String `tfsdk:"access_token"`
	IntegrationKey types.String `tfsdk:"integration_key"`
	ApiKey         types.String `tfsdk:"api_key"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
}

type SinkSettingsModel struct {
	EventIndex      types.String `tfsdk:"event_index"`
	MetricIndex     types.String `tfsdk:"metric_index"`
	Source          types.String `tfsdk:"source"`
	ApplicationKey  types.String `tfsdk:"application_key"`
	EventUri        types.String `tfsdk:"event_uri"`
	MetricUri       types.String `tfsdk:"metric_uri"`
	TransformAlerts types.Bool   `tfsdk:"transform_alerts"`
}

type ChangeLogModel struct {
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

func (m *DataSourceByIDsModel) parse(ctx context.Context, streamSubscription *fabricv4.StreamSubscription) diag.Diagnostics {
	m.StreamID = types.StringValue(streamSubscription.GetUuid())
	m.SubscriptionID = types.StringValue(streamSubscription.GetUuid())
	m.ID = types.StringValue(streamSubscription.GetUuid())

	diags := parseStreamSubscription(ctx, streamSubscription,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Enabled,
		&m.Filters,
		&m.MetricSelector,
		&m.EventSelector,
		&m.Sink,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}

	return diags
}

func (m *DataSourceAll) parse(ctx context.Context, streamSubscriptionsResponse *fabricv4.GetAllStreamSubscriptionResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(streamSubscriptionsResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by streamSubscription subscriptions data source",
			"either the account does not have any streamSubscription subscription data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]BaseStreamSubscriptionModel, len(streamSubscriptionsResponse.GetData()))
	streamSubscriptions := streamSubscriptionsResponse.GetData()
	for index, streamSubscription := range streamSubscriptions {
		var streamSubscriptionModel BaseStreamSubscriptionModel
		diags = streamSubscriptionModel.parse(ctx, &streamSubscription)
		if diags.HasError() {
			return diags
		}
		data[index] = streamSubscriptionModel
	}
	responsePagination := streamSubscriptionsResponse.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}

	m.ID = types.StringValue(data[0].Uuid.ValueString())
	m.StreamID = types.StringValue(data[0].Uuid.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)
	m.Data = fwtypes.NewListNestedObjectValueOfValueSlice[BaseStreamSubscriptionModel](ctx, data)

	return diags
}

func (m *ResourceModel) parse(ctx context.Context, streamSubscription *fabricv4.StreamSubscription) diag.Diagnostics {
	m.ID = types.StringValue(streamSubscription.GetUuid())

	diags := parseStreamSubscription(ctx, streamSubscription,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Enabled,
		&m.Filters,
		&m.MetricSelector,
		&m.EventSelector,
		&m.Sink,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}

	return diags
}

func (m *BaseStreamSubscriptionModel) parse(ctx context.Context, streamSubscription *fabricv4.StreamSubscription) diag.Diagnostics {
	diags := parseStreamSubscription(ctx, streamSubscription,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Enabled,
		&m.Filters,
		&m.MetricSelector,
		&m.EventSelector,
		&m.Sink,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}

	return diags
}

func parseStreamSubscription(ctx context.Context, streamSubscription *fabricv4.StreamSubscription,
	streamSubscriptionType, name, description *basetypes.StringValue,
	enabled *basetypes.BoolValue,
	filters *fwtypes.ListNestedObjectValueOf[FilterModel],
	metricSelector, eventSelector *fwtypes.ObjectValueOf[SelectorModel],
	sink *fwtypes.ObjectValueOf[SinkModel],
	href, uuid, state *basetypes.StringValue,
	changeLog *fwtypes.ObjectValueOf[ChangeLogModel]) diag.Diagnostics {

	var diags diag.Diagnostics

	*streamSubscriptionType = types.StringValue(string(streamSubscription.GetType()))
	*name = types.StringValue(streamSubscription.GetName())
	*description = types.StringValue(streamSubscription.GetDescription())
	*href = types.StringValue(streamSubscription.GetHref())
	*uuid = types.StringValue(streamSubscription.GetUuid())
	*state = types.StringValue(string(streamSubscription.GetState()))
	*enabled = types.BoolValue(streamSubscription.GetEnabled())

	// Parse filters
	streamSubscriptionFilters := streamSubscription.GetFilters()
	filterModels := make([]FilterModel, len(streamSubscriptionFilters.GetAnd()))
	for i, filter := range streamSubscriptionFilters.GetAnd() {
		if len(filter.StreamFilterOrFilter.GetOr()) > 0 {
			for j, orFilter := range filter.StreamFilterOrFilter.GetOr() {
				orFilterModel, diags := parseSimpleExpression(ctx, &orFilter, true)
				if diags.HasError() {
					return diags
				}
				if j == 0 {
					filterModels[i] = orFilterModel
				} else {
					filterModels = append(filterModels, orFilterModel)
				}
			}
		} else {
			filterModels[i], diags = parseSimpleExpression(ctx, filter.StreamFilterSimpleExpression, false)
			if diags.HasError() {
				return diags
			}
		}
	}
	*filters = fwtypes.NewListNestedObjectValueOfValueSlice[FilterModel](ctx, filterModels)

	// Parse MetricSelector
	metricSelectorObject, diags := parseSelectorModel(ctx, streamSubscription.GetMetricSelector())
	if diags.HasError() {
		return diags
	}
	*metricSelector = metricSelectorObject

	// Parse EventSelector
	eventSelectorObject, diags := parseSelectorModel(ctx, streamSubscription.GetEventSelector())
	if diags.HasError() {
		return diags
	}
	*eventSelector = eventSelectorObject

	// Parse Sink
	streamSubSink := streamSubscription.GetSink()
	sinkModel := SinkModel{
		Uri:              types.StringValue(streamSubSink.GetUri()),
		Type:             types.StringValue(string(streamSubSink.GetType())),
		BatchEnabled:     types.BoolValue(streamSubSink.GetBatchEnabled()),
		BatchSizeMax:     types.Int32Value(streamSubSink.GetBatchSizeMax()),
		BatchWaitTimeMax: types.Int32Value(streamSubSink.GetBatchWaitTimeMax()),
		Host:             types.StringValue(streamSubSink.GetHost()),
	}
	sinkCredential := streamSubSink.GetCredential()
	credentialModel := SinkCredentialModel{}
	credentialModel.Type = types.StringValue(string(sinkCredential.GetType()))

	switch sinkCredential.GetType() {
	case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_ACCESS_TOKEN:
		credentialModel.AccessToken = types.StringValue(sinkCredential.GetAccessToken())
	case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_INTEGRATION_KEY:
		credentialModel.IntegrationKey = types.StringValue(sinkCredential.GetIntegrationKey())
	case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_API_KEY:
		credentialModel.ApiKey = types.StringValue(sinkCredential.GetApiKey())
		sinkCredential.SetApiKey(credentialModel.ApiKey.ValueString())
	case fabricv4.STREAMSUBSCRIPTIONSINKCREDENTIALTYPE_USERNAME_PASSWORD:
		credentialModel.Username = types.StringValue(sinkCredential.GetUsername())
		credentialModel.Password = types.StringValue(sinkCredential.GetPassword())
	}

	sinkModel.Credential = fwtypes.NewObjectValueOf[SinkCredentialModel](ctx, &credentialModel)

	sinkSettings := streamSubSink.GetSettings()
	sinkSettingsModel := SinkSettingsModel{
		EventIndex:      types.StringValue(sinkSettings.GetEventIndex()),
		MetricIndex:     types.StringValue(sinkSettings.GetMetricIndex()),
		Source:          types.StringValue(sinkSettings.GetSource()),
		ApplicationKey:  types.StringValue(sinkSettings.GetApplicationKey()),
		EventUri:        types.StringValue(sinkSettings.GetEventUri()),
		MetricUri:       types.StringValue(sinkSettings.GetMetricUri()),
		TransformAlerts: types.BoolValue(sinkSettings.GetTransformAlerts()),
	}

	sinkModel.Settings = fwtypes.NewObjectValueOf[SinkSettingsModel](ctx, &sinkSettingsModel)

	*sink = fwtypes.NewObjectValueOf[SinkModel](ctx, &sinkModel)

	// Parse ChangeLog
	streamSubscriptionChangeLog := streamSubscription.GetChangeLog()
	changeLogModel := ChangeLogModel{
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
	*changeLog = fwtypes.NewObjectValueOf[ChangeLogModel](ctx, &changeLogModel)

	return diags
}

func parseSimpleExpression(ctx context.Context, expression *fabricv4.StreamFilterSimpleExpression, orGroup bool) (FilterModel, diag.Diagnostics) {
	values := int_fw.StringSliceToAttrValue(expression.GetValues())
	fwValues, diags := fwtypes.NewListValueOf[types.String](ctx, values)
	if diags.HasError() {
		return FilterModel{}, diags
	}
	return FilterModel{
		Property: types.StringValue(expression.GetProperty()),
		Operator: types.StringValue(expression.GetOperator()),
		Values:   fwValues,
		Or:       types.BoolValue(orGroup),
	}, nil
}

func parseSelectorModel(ctx context.Context, selector fabricv4.StreamSubscriptionSelector) (fwtypes.ObjectValueOf[SelectorModel], diag.Diagnostics) {
	var diags diag.Diagnostics
	inclusions, diags := fwtypes.NewListValueOf[types.String](ctx, int_fw.StringSliceToAttrValue(selector.GetInclude()))
	if diags.HasError() {
		return fwtypes.NewObjectValueOfNull[SelectorModel](ctx), diags
	}
	exclusions, diags := fwtypes.NewListValueOf[types.String](ctx, int_fw.StringSliceToAttrValue(selector.GetExcept()))
	if diags.HasError() {
		return fwtypes.NewObjectValueOfNull[SelectorModel](ctx), diags
	}
	selectorModel := SelectorModel{
		Include: inclusions,
		Except:  exclusions,
	}
	return fwtypes.NewObjectValueOf[SelectorModel](ctx, &selectorModel), diags
}
