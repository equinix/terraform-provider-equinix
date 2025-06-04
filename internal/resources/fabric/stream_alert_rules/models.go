package stream_alert_rules

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/fabric"
	int_fw "github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	_ "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type selectorModel struct {
	Include fwtypes.ListValueOf[types.String] `tfsdk:"include"`
}

type resourceModel struct {
	StreamID types.String   `tfsdk:"stream_id"`
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	baseStreamAlertRulesModel
}

type baseStreamAlertRulesModel struct {
	Type              types.String                          `tfsdk:"type"`
	Name              types.String                          `tfsdk:"name"`
	Description       types.String                          `tfsdk:"description"`
	Enabled           types.Bool                            `tfsdk:"enabled"`
	MetricName        types.String                          `tfsdk:"metric_name"`
	ResourceSelector  fwtypes.ObjectValueOf[selectorModel]  `tfsdk:"resource_selector"` // Object of ResourceSelectorModel
	WindowSize        types.String                          `tfsdk:"window_size"`
	WarningThreshold  types.String                          `tfsdk:"warning_threshold"`
	CriticalThreshold types.String                          `tfsdk:"critical_threshold"`
	Operand           types.String                          `tfsdk:"operand"`
	Href              types.String                          `tfsdk:"href"`
	UUID              types.String                          `tfsdk:"uuid"`
	State             types.String                          `tfsdk:"state"`
	ChangeLog         fwtypes.ObjectValueOf[changeLogModel] `tfsdk:"change_log"` // Object of ChangeLogModel
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

func (m *baseStreamAlertRulesModel) parse(ctx context.Context, streamAlertRule *fabricv4.StreamAlertRule) diag.Diagnostics {

	var mDiags diag.Diagnostics

	m.Type = types.StringValue(string(streamAlertRule.GetType()))
	m.Name = types.StringValue(streamAlertRule.GetName())
	m.Description = types.StringValue(streamAlertRule.GetDescription())
	m.Href = types.StringValue(streamAlertRule.GetHref())
	m.UUID = types.StringValue(streamAlertRule.GetUuid())
	m.State = types.StringValue(string(streamAlertRule.GetState()))
	m.Enabled = types.BoolValue(streamAlertRule.GetEnabled())
	m.WindowSize = types.StringValue(streamAlertRule.GetWindowSize())
	m.WarningThreshold = types.StringValue(streamAlertRule.GetWarningThreshold())
	m.CriticalThreshold = types.StringValue(streamAlertRule.GetCriticalThreshold())
	m.Operand = types.StringValue(string(streamAlertRule.GetOperand()))

	// Parse ResourceSelector
	resourceSelectorObject, diags := parseSelectorModel(ctx, streamAlertRule.GetResourceSelector())
	if diags.HasError() {
		mDiags.Append(diags...)
		return mDiags
	}
	m.ResourceSelector = resourceSelectorObject

	// Parse ChangeLog
	streamSubscriptionChangeLog := streamAlertRule.GetChangeLog()
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

func parseSelectorModel(ctx context.Context, alertRuleSubSelector fabricv4.ResourceSelector) (fwtypes.ObjectValueOf[selectorModel], diag.Diagnostics) {
	var diags diag.Diagnostics
	inclusions, diags := fwtypes.NewListValueOf[types.String](ctx, int_fw.StringSliceToAttrValue(alertRuleSubSelector.GetInclude()))
	if diags.HasError() {
		return fwtypes.NewObjectValueOfNull[selectorModel](ctx), diags
	}
	selector := selectorModel{
		Include: inclusions,
	}
	return fwtypes.NewObjectValueOf[selectorModel](ctx, &selector), diags
}
