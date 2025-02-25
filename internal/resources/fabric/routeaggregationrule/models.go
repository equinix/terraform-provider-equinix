package routeaggregationrule

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceByIDModel struct {
	RouteAggregationRuleID types.String `tfsdk:"route_aggregation_rule_id"`
	ID                     types.String `tfsdk:"id"`
	baseRouteAggregationRuleModel
}

type datsSourceAllRouteAggregationRulesModel struct {
	ID                 types.String                                                   `tfsdk:"id"`
	Data               fwtypes.ListNestedObjectValueOf[baseRouteAggregationRuleModel] `tfsdk:"data"`
	Pagination         fwtypes.ObjectValueOf[paginationModel]                         `tfsdk:"pagination"`
	RouteAggregationID types.String                                                   `tfsdk:"route_aggregation_id"`
}

type paginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type resourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	baseRouteAggregationRuleModel
}

type baseRouteAggregationRuleModel struct {
	RouteAggregationID types.String                          `tfsdk:"route_aggregation_id"`
	Name               types.String                          `tfsdk:"name"`
	Description        types.String                          `tfsdk:"description"`
	Prefix             types.String                          `tfsdk:"prefix"`
	Href               types.String                          `tfsdk:"href"`
	Type               types.String                          `tfsdk:"type"`
	UUID               types.String                          `tfsdk:"uuid"`
	State              types.String                          `tfsdk:"state"`
	Change             fwtypes.ObjectValueOf[changeModel]    `tfsdk:"change"`
	ChangeLog          fwtypes.ObjectValueOf[changeLogModel] `tfsdk:"change_log"`
}

type changeModel struct {
	UUID types.String `tfsdk:"uuid"`
	Type types.String `tfsdk:"type"`
	Href types.String `tfsdk:"href"`
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

func (m *dataSourceByIDModel) parse(ctx context.Context, routeAggregationRule *fabricv4.RouteAggregationRulesData) diag.Diagnostics {
	m.RouteAggregationRuleID = types.StringValue(routeAggregationRule.GetUuid())
	m.ID = types.StringValue(routeAggregationRule.GetUuid())
	diags := m.baseRouteAggregationRuleModel.parse(ctx, routeAggregationRule)
	return diags
}

func (m *datsSourceAllRouteAggregationRulesModel) parse(ctx context.Context, routeAggregationRulesResponse *fabricv4.GetRouteAggregationRulesResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(routeAggregationRulesResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by route aggregation rules data source", "either the account does not have any route aggregation rules data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]baseRouteAggregationRuleModel, len(routeAggregationRulesResponse.GetData()))
	routeAggregationRules := routeAggregationRulesResponse.GetData()
	for index, routeAggregationRule := range routeAggregationRules {
		var routeAggregationRuleModel baseRouteAggregationRuleModel
		diags = routeAggregationRuleModel.parse(ctx, &routeAggregationRule)
		if diags.HasError() {
			return diags
		}
		data[index] = routeAggregationRuleModel
	}
	responsePagination := routeAggregationRulesResponse.GetPagination()
	pagination := paginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}
	m.ID = types.StringValue(data[0].UUID.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[paginationModel](ctx, &pagination)

	dataPtr := make([]*baseRouteAggregationRuleModel, len(data))
	for i := range data {
		dataPtr[i] = &data[i]
	}
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[baseRouteAggregationRuleModel](ctx, dataPtr)

	return diags
}

func (m *resourceModel) parse(ctx context.Context, routeAggregationRule *fabricv4.RouteAggregationRulesData) diag.Diagnostics {

	m.ID = types.StringValue(routeAggregationRule.GetUuid())

	diags := m.baseRouteAggregationRuleModel.parse(ctx, routeAggregationRule)
	return diags
}

func (m *baseRouteAggregationRuleModel) parse(ctx context.Context, routeAggregationRule *fabricv4.RouteAggregationRulesData) diag.Diagnostics {

	var diag diag.Diagnostics
	m.Name = types.StringValue(routeAggregationRule.GetName())
	m.Description = types.StringValue(routeAggregationRule.GetDescription())
	m.Prefix = types.StringValue(routeAggregationRule.GetPrefix())
	m.Href = types.StringValue(routeAggregationRule.GetHref())
	m.Type = types.StringValue(string(routeAggregationRule.GetType()))
	m.UUID = types.StringValue(routeAggregationRule.GetUuid())
	m.State = types.StringValue(string(routeAggregationRule.GetState()))
	routeAggregationRuleChange := routeAggregationRule.GetChange()
	changemodel := changeModel{
		UUID: types.StringValue(routeAggregationRuleChange.GetUuid()),
		Type: types.StringValue(string(routeAggregationRuleChange.GetType())),
		Href: types.StringValue(routeAggregationRuleChange.GetHref()),
	}
	m.Change = fwtypes.NewObjectValueOf[changeModel](ctx, &changemodel)

	const TIMEFORMAT = "2008-02-02T14:02:02.000Z"
	routeAggregationRuleChangeLog := routeAggregationRule.GetChangeLog()
	changelogModel := changeLogModel{
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
	m.ChangeLog = fwtypes.NewObjectValueOf[changeLogModel](ctx, &changelogModel)

	return diag
}
