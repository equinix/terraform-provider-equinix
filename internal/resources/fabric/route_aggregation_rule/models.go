package route_aggregation_rule

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DataSourceByIdModel struct {
	RouteAggregationRuleId types.String `tfsdk:"route_aggregation_rule_id"`
	ID                     types.String `tfsdk:"id"`
	BaseRouteAggregationRuleModel
}

type DatsSourceAllRouteAggregationRulesModel struct {
	ID                 types.String                                                   `tfsdk:"id"`
	Data               fwtypes.ListNestedObjectValueOf[BaseRouteAggregationRuleModel] `tfsdk:"data"`
	Pagination         fwtypes.ObjectValueOf[PaginationModel]                         `tfsdk:"pagination"`
	RouteAggregationId types.String                                                   `tfsdk:"route_aggregation_id"`
}

type PaginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type ResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	BaseRouteAggregationRuleModel
}

type BaseRouteAggregationRuleModel struct {
	RouteAggregationID types.String                          `tfsdk:"route_aggregation_id"`
	Name               types.String                          `tfsdk:"name"`
	Description        types.String                          `tfsdk:"description"`
	Prefix             types.String                          `tfsdk:"prefix"`
	Href               types.String                          `tfsdk:"href"`
	Type               types.String                          `tfsdk:"type"`
	Uuid               types.String                          `tfsdk:"uuid"`
	State              types.String                          `tfsdk:"state"`
	Change             fwtypes.ObjectValueOf[ChangeModel]    `tfsdk:"change"`
	ChangeLog          fwtypes.ObjectValueOf[ChangeLogModel] `tfsdk:"change_log"`
}

type ChangeModel struct {
	Uuid types.String `tfsdk:"uuid"`
	Type types.String `tfsdk:"type"`
	Href types.String `tfsdk:"href"`
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

func (m *DataSourceByIdModel) parse(ctx context.Context, routeAggregationRule *fabricv4.RouteAggregationRulesData) diag.Diagnostics {
	m.RouteAggregationRuleId = types.StringValue(routeAggregationRule.GetUuid())
	m.ID = types.StringValue(routeAggregationRule.GetUuid())

	diags := parseRouteAggregationRule(ctx, routeAggregationRule,
		&m.Name,
		&m.Description,
		&m.Prefix,
		&m.Href,
		&m.Type,
		&m.Uuid,
		&m.State,
		&m.Change,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}
	return diags
}

func (m *DatsSourceAllRouteAggregationRulesModel) parse(ctx context.Context, routeAggregationRulesResponse *fabricv4.GetRouteAggregationRulesResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(routeAggregationRulesResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by route aggregation rules data source", "either the account does not have any route aggregation rules data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]BaseRouteAggregationRuleModel, len(routeAggregationRulesResponse.GetData()))
	routeAggregationRules := routeAggregationRulesResponse.GetData()
	for index, routeAggregationRule := range routeAggregationRules {
		var routeAggregationRuleModel BaseRouteAggregationRuleModel
		diags = routeAggregationRuleModel.parse(ctx, &routeAggregationRule)
		if diags.HasError() {
			return diags
		}
		data[index] = routeAggregationRuleModel
	}
	responsePagination := routeAggregationRulesResponse.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}
	m.ID = types.StringValue(data[0].Uuid.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)

	dataPtr := make([]*BaseRouteAggregationRuleModel, len(data))
	for i := range data {
		dataPtr[i] = &data[i]
	}
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[BaseRouteAggregationRuleModel](ctx, dataPtr)

	return diags
}

func (m *ResourceModel) parse(ctx context.Context, routeAggregationRule *fabricv4.RouteAggregationRulesData) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(routeAggregationRule.GetUuid())

	diags = parseRouteAggregationRule(ctx, routeAggregationRule,
		&m.Name,
		&m.Description,
		&m.Prefix,
		&m.Href,
		&m.Type,
		&m.Uuid,
		&m.State,
		&m.Change,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}
	return diags
}

func (m *BaseRouteAggregationRuleModel) parse(ctx context.Context, routeAggregationRule *fabricv4.RouteAggregationRulesData) diag.Diagnostics {
	var diags diag.Diagnostics = parseRouteAggregationRule(ctx, routeAggregationRule,
		&m.Name,
		&m.Description,
		&m.Prefix,
		&m.Href,
		&m.Type,
		&m.Uuid,
		&m.State,
		&m.Change,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}
	return diags
}

func parseRouteAggregationRule(ctx context.Context, routeAggregationRule *fabricv4.RouteAggregationRulesData,
	name, description, prefix, href, type_, uuid, state *basetypes.StringValue,
	change *fwtypes.ObjectValueOf[ChangeModel],
	changeLog *fwtypes.ObjectValueOf[ChangeLogModel]) diag.Diagnostics {
	var diag diag.Diagnostics

	*name = types.StringValue(routeAggregationRule.GetName())
	*description = types.StringValue(routeAggregationRule.GetDescription())
	*prefix = types.StringValue(routeAggregationRule.GetPrefix())
	*href = types.StringValue(routeAggregationRule.GetHref())
	*type_ = types.StringValue(string(routeAggregationRule.GetType()))
	*uuid = types.StringValue(routeAggregationRule.GetUuid())
	*state = types.StringValue(string(routeAggregationRule.GetState()))

	routeAggregationRuleChange := routeAggregationRule.GetChange()
	changeModel := ChangeModel{
		Uuid: types.StringValue(routeAggregationRuleChange.GetUuid()),
		Type: types.StringValue(string(routeAggregationRuleChange.GetType())),
		Href: types.StringValue(routeAggregationRuleChange.GetHref()),
	}
	*change = fwtypes.NewObjectValueOf[ChangeModel](ctx, &changeModel)

	const TIMEFORMAT = "2008-02-02T14:02:02.000Z"
	routeAggregationRuleChangeLog := routeAggregationRule.GetChangeLog()
	changeLogModel := ChangeLogModel{
		CreatedBy:         types.StringValue(routeAggregationRuleChangeLog.GetCreatedBy()),
		CreatedByFullName: types.StringValue(routeAggregationRuleChangeLog.GetCreatedByFullName()),
		CreatedByEmail:    types.StringValue(routeAggregationRuleChangeLog.GetCreatedByEmail()),
		CreatedDateTime:   types.StringValue(routeAggregationRuleChangeLog.GetCreatedDateTime().Format(TIMEFORMAT)),
		UpdatedBy:         types.StringValue(routeAggregationRuleChangeLog.GetUpdatedBy()),
		UpdatedByFullName: types.StringValue(routeAggregationRuleChangeLog.GetUpdatedByFullName()),
		UpdatedByEmail:    types.StringValue(routeAggregationRuleChangeLog.GetUpdatedByEmail()),
		UpdatedDateTime:   types.StringValue(routeAggregationRuleChangeLog.GetUpdatedDateTime().Format(TIMEFORMAT)),
		DeletedBy:         types.StringValue(routeAggregationRuleChangeLog.GetDeletedBy()),
		DeletedByFullName: types.StringValue(routeAggregationRuleChangeLog.GetDeletedByFullName()),
		DeletedByEmail:    types.StringValue(routeAggregationRuleChangeLog.GetDeletedByEmail()),
		DeletedDateTime:   types.StringValue(routeAggregationRuleChangeLog.GetDeletedDateTime().Format(TIMEFORMAT)),
	}
	*changeLog = fwtypes.NewObjectValueOf[ChangeLogModel](ctx, &changeLogModel)
	return diag
}
