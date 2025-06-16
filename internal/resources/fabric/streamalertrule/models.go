package streamalertrule

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/fabric"
	int_fw "github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type selectorModel struct {
	Include fwtypes.ListValueOf[types.String] `tfsdk:"include"`
}

type streamAlertRuleResourceModel struct {
	StreamID types.String   `tfsdk:"stream_id"`
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	baseStreamAlertRulesModel
}

type baseStreamAlertRulesModel struct {
	Type              types.String                         `tfsdk:"type"`
	Name              types.String                         `tfsdk:"name"`
	Description       types.String                         `tfsdk:"description"`
	Enabled           types.Bool                           `tfsdk:"enabled"`
	MetricName        types.String                         `tfsdk:"metric_name"`
	ResourceSelector  fwtypes.ObjectValueOf[selectorModel] `tfsdk:"resource_selector"` // Object of ResourceSelectorModel
	WindowSize        types.String                         `tfsdk:"window_size"`
	WarningThreshold  types.String                         `tfsdk:"warning_threshold"`
	CriticalThreshold types.String                         `tfsdk:"critical_threshold"`
	Operand           types.String                         `tfsdk:"operand"`
	Href              types.String                         `tfsdk:"href"`
	UUID              types.String                         `tfsdk:"uuid"`
	State             types.String                         `tfsdk:"state"`
	ChangeLog         fwtypes.ObjectValueOf[LogModel]      `tfsdk:"change_log"` // Object of LogModel
}

type LogModel struct {
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

type dataSourceStreamAlertRuleByIDsModel struct {
	ID          types.String `tfsdk:"id"`
	StreamID    types.String `tfsdk:"stream_id"`
	AlertRuleID types.String `tfsdk:"alert_rule_id"`
	baseStreamAlertRulesModel
}

type paginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type dataSourceAllStreamAlertRulesModel struct {
	ID         types.String                                               `tfsdk:"id"`
	StreamID   types.String                                               `tfsdk:"stream_id"`
	Pagination fwtypes.ObjectValueOf[paginationModel]                     `tfsdk:"pagination"`
	Data       fwtypes.ListNestedObjectValueOf[baseStreamAlertRulesModel] `tfsdk:"data"`
}

func (m *baseStreamAlertRulesModel) parse(ctx context.Context, streamAlertRule *fabricv4.StreamAlertRule) diag.Diagnostics {
	diags := parseAlertRule(ctx, streamAlertRule, &m.Type, &m.Name, &m.Description,
		&m.Href, &m.UUID, &m.State, &m.WindowSize, &m.WarningThreshold,
		&m.CriticalThreshold, &m.Operand, &m.MetricName, &m.Enabled, &m.ResourceSelector, &m.ChangeLog)
	return diags
}

func (m *streamAlertRuleResourceModel) parse(ctx context.Context, streamAlertRule *fabricv4.StreamAlertRule) diag.Diagnostics {
	m.ID = types.StringValue(streamAlertRule.GetUuid())
	diags := m.baseStreamAlertRulesModel.parse(ctx, streamAlertRule)
	if diags.HasError() {
		return diags
	}
	return diags
}

func parseAlertRule(ctx context.Context, streamAlertRule *fabricv4.StreamAlertRule,
	streamType, name, description, href, uuid, state, windowSize, warningThreshold,
	criticalThreshold, operand, metricName *basetypes.StringValue,
	enabled *basetypes.BoolValue,
	resourceSelector *fwtypes.ObjectValueOf[selectorModel],
	changeLog *fwtypes.ObjectValueOf[LogModel]) diag.Diagnostics {

	var mDiags diag.Diagnostics

	*streamType = types.StringValue(string(streamAlertRule.GetType()))
	*name = types.StringValue(streamAlertRule.GetName())
	*description = types.StringValue(streamAlertRule.GetDescription())
	*href = types.StringValue(streamAlertRule.GetHref())
	*uuid = types.StringValue(streamAlertRule.GetUuid())
	*state = types.StringValue(string(streamAlertRule.GetState()))
	*enabled = types.BoolValue(streamAlertRule.GetEnabled())
	*windowSize = types.StringValue(streamAlertRule.GetWindowSize())
	*warningThreshold = types.StringValue(streamAlertRule.GetWarningThreshold())
	*criticalThreshold = types.StringValue(streamAlertRule.GetCriticalThreshold())
	*operand = types.StringValue(string(streamAlertRule.GetOperand()))
	*metricName = types.StringValue(string(streamAlertRule.GetMetricName()))

	// Parse ResourceSelector
	getResourceSelector := streamAlertRule.GetResourceSelector()
	inclusions, diags := fwtypes.NewListValueOf[types.String](ctx, int_fw.StringSliceToAttrValue(getResourceSelector.GetInclude()))
	if diags.HasError() {
		return diags
	}
	selector := selectorModel{
		Include: inclusions,
	}
	*resourceSelector = fwtypes.NewObjectValueOf[selectorModel](ctx, &selector)

	// Parse ChangeLog
	streamAlertRuleChangeLog := streamAlertRule.GetChangeLog()
	changeLogModel := LogModel{
		CreatedBy:         types.StringValue(streamAlertRuleChangeLog.GetCreatedBy()),
		CreatedByFullName: types.StringValue(streamAlertRuleChangeLog.GetCreatedByFullName()),
		CreatedByEmail:    types.StringValue(streamAlertRuleChangeLog.GetCreatedByEmail()),
		CreatedDateTime:   types.StringValue(streamAlertRuleChangeLog.GetCreatedDateTime().Format(fabric.TimeFormat)),
		UpdatedBy:         types.StringValue(streamAlertRuleChangeLog.GetUpdatedBy()),
		UpdatedByFullName: types.StringValue(streamAlertRuleChangeLog.GetUpdatedByFullName()),
		UpdatedByEmail:    types.StringValue(streamAlertRuleChangeLog.GetUpdatedByEmail()),
		UpdatedDateTime:   types.StringValue(streamAlertRuleChangeLog.GetUpdatedDateTime().Format(fabric.TimeFormat)),
		DeletedBy:         types.StringValue(streamAlertRuleChangeLog.GetDeletedBy()),
		DeletedByFullName: types.StringValue(streamAlertRuleChangeLog.GetDeletedByFullName()),
		DeletedByEmail:    types.StringValue(streamAlertRuleChangeLog.GetDeletedByEmail()),
		DeletedDateTime:   types.StringValue(streamAlertRuleChangeLog.GetDeletedDateTime().Format(fabric.TimeFormat)),
	}
	*changeLog = fwtypes.NewObjectValueOf[LogModel](ctx, &changeLogModel)

	return mDiags
}

func (m *dataSourceAllStreamAlertRulesModel) parse(ctx context.Context, streamAlertRulesResponse *fabricv4.GetAllStreamAlertRuleResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(streamAlertRulesResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by stream alert rule data source",
			"either the account does not have any stream alert rule data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]baseStreamAlertRulesModel, len(streamAlertRulesResponse.GetData()))
	streamAlertRules := streamAlertRulesResponse.GetData()
	for index, streamAlertRule := range streamAlertRules {
		var streamAlertRuleModel baseStreamAlertRulesModel
		diags = streamAlertRuleModel.parse(ctx, &streamAlertRule)
		if diags.HasError() {
			return diags
		}
		data[index] = streamAlertRuleModel
	}
	responsePagination := streamAlertRulesResponse.GetPagination()
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
	m.Data = fwtypes.NewListNestedObjectValueOfValueSlice[baseStreamAlertRulesModel](ctx, data)

	return diags
}

func (m *dataSourceStreamAlertRuleByIDsModel) parse(ctx context.Context, streamAlertRule *fabricv4.StreamAlertRule) diag.Diagnostics {
	m.StreamID = types.StringValue(streamAlertRule.GetUuid())
	m.AlertRuleID = types.StringValue(streamAlertRule.GetUuid())
	m.ID = types.StringValue(streamAlertRule.GetUuid())

	diags := m.baseStreamAlertRulesModel.parse(ctx, streamAlertRule)
	if diags.HasError() {
		return diags
	}

	return diags
}
