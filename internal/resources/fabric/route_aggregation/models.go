package route_aggregation

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
	RouteAggregationId types.String `tfsdk:"route_aggregation_id"`
	ID                 types.String `tfsdk:"id"`
	BaseRouteAggregationModel
}

type DatsSourceAllRouteAggregationsModel struct {
	ID         types.String                                               `tfsdk:"id"`
	Data       fwtypes.ListNestedObjectValueOf[BaseRouteAggregationModel] `tfsdk:"data"`
	Filter     types.Object                                               `tfsdk:"filter"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]                     `tfsdk:"pagination"`
	Sort       fwtypes.ObjectValueOf[SortModel]                           `tfsdk:"sort"`
}

type FilterModel struct {
	Property types.String   `tfsdk:"property"`
	Operator types.String   `tfsdk:"operator"`
	Values   []types.String `tfsdk:"values"`
}

type PaginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type SortModel struct {
	Direction types.String `tfsdk:"direction"`
	Property  types.String `tfsdk:"property"`
}

type ResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	BaseRouteAggregationModel
}

type BaseRouteAggregationModel struct {
	Type             types.String                          `tfsdk:"type"`
	Name             types.String                          `tfsdk:"name"`
	Description      types.String                          `tfsdk:"description"`
	Href             types.String                          `tfsdk:"href"`
	Uuid             types.String                          `tfsdk:"uuid"`
	State            types.String                          `tfsdk:"state"`
	ConnectionsCount types.Int32                           `tfsdk:"connections_count"`
	RulesCount       types.Int32                           `tfsdk:"rules_count"`
	Project          fwtypes.ObjectValueOf[ProjectModel]   `tfsdk:"project"`
	Change           fwtypes.ObjectValueOf[ChangeModel]    `tfsdk:"change"`
	ChangeLog        fwtypes.ObjectValueOf[ChangeLogModel] `tfsdk:"change_log"`
}

type ProjectModel struct {
	ProjectId types.String `tfsdk:"project_id"`
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

func (m *DataSourceByIdModel) parse(ctx context.Context, routeAggregation *fabricv4.RouteAggregationsData) diag.Diagnostics {
	m.RouteAggregationId = types.StringValue(routeAggregation.GetUuid())
	m.ID = types.StringValue(routeAggregation.GetUuid())

	diags := parseRouteAggregation(ctx, routeAggregation,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.ConnectionsCount,
		&m.RulesCount,
		&m.Project,
		&m.Change,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}
	return diags
}

func (m *DatsSourceAllRouteAggregationsModel) parse(ctx context.Context, routeAggregationsResponse *fabricv4.RouteAggregationsSearchResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(routeAggregationsResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by streams data source", "either the account does not have any streams data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]BaseRouteAggregationModel, len(routeAggregationsResponse.GetData()))
	routeAggregations := routeAggregationsResponse.GetData()
	for index, routeAggregation := range routeAggregations {
		var routeAggregationModel BaseRouteAggregationModel
		diags = routeAggregationModel.parse(ctx, &routeAggregation)
		if diags.HasError() {
			return diags
		}
		data[index] = routeAggregationModel
	}
	responsePagination := routeAggregationsResponse.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}
	m.ID = types.StringValue(data[0].Uuid.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)

	dataPtr := make([]*BaseRouteAggregationModel, len(data))
	for i := range data {
		dataPtr[i] = &data[i]
	}
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[BaseRouteAggregationModel](ctx, dataPtr)

	return diags
}
func (m *ResourceModel) parse(ctx context.Context, routeAggregation *fabricv4.RouteAggregationsData) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(routeAggregation.GetUuid())

	diags = parseRouteAggregation(ctx, routeAggregation,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.ConnectionsCount,
		&m.RulesCount,
		&m.Project,
		&m.Change,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}

	return diags

}

func (m *BaseRouteAggregationModel) parse(ctx context.Context, routeAggregation *fabricv4.RouteAggregationsData) diag.Diagnostics {
	var diags diag.Diagnostics = parseRouteAggregation(ctx, routeAggregation,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.ConnectionsCount,
		&m.RulesCount,
		&m.Project,
		&m.Change,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}
	return diags
}

func parseRouteAggregation(ctx context.Context, routeAggregation *fabricv4.RouteAggregationsData,
	type_, name, description, href, uuid, state *basetypes.StringValue,
	connectionsCount, rulesCount *basetypes.Int32Value,
	project *fwtypes.ObjectValueOf[ProjectModel],
	change *fwtypes.ObjectValueOf[ChangeModel],
	changeLog *fwtypes.ObjectValueOf[ChangeLogModel]) diag.Diagnostics {
	var diag diag.Diagnostics

	*type_ = types.StringValue(string(routeAggregation.GetType()))
	*name = types.StringValue(routeAggregation.GetName())
	*description = types.StringValue(routeAggregation.GetDescription())
	*href = types.StringValue(routeAggregation.GetHref())
	*uuid = types.StringValue(routeAggregation.GetUuid())
	*state = types.StringValue(string(routeAggregation.GetState()))
	*connectionsCount = types.Int32Value(routeAggregation.GetConnectionsCount())
	*rulesCount = types.Int32Value(routeAggregation.GetRulesCount())

	routeAggregationProject := routeAggregation.GetProject()
	projectModel := ProjectModel{
		ProjectId: types.StringValue(routeAggregationProject.GetProjectId()),
	}
	*project = fwtypes.NewObjectValueOf[ProjectModel](ctx, &projectModel)

	routeAggregationChange := routeAggregation.GetChange()
	changeModel := ChangeModel{
		Uuid: types.StringValue(routeAggregationChange.GetUuid()),
		Type: types.StringValue(string(routeAggregationChange.GetType())),
		Href: types.StringValue(routeAggregationChange.GetHref()),
	}
	*change = fwtypes.NewObjectValueOf[ChangeModel](ctx, &changeModel)

	const TIMEFORMAT = "2008-02-02T14:02:02.000Z"
	routeAggregationChangeLog := routeAggregation.GetChangeLog()
	changeLogModel := ChangeLogModel{
		CreatedBy:         types.StringValue(routeAggregationChangeLog.GetCreatedBy()),
		CreatedByFullName: types.StringValue(routeAggregationChangeLog.GetCreatedByFullName()),
		CreatedByEmail:    types.StringValue(routeAggregationChangeLog.GetCreatedByEmail()),
		CreatedDateTime:   types.StringValue(routeAggregationChangeLog.GetCreatedDateTime().Format(TIMEFORMAT)),
		UpdatedBy:         types.StringValue(routeAggregationChangeLog.GetUpdatedBy()),
		UpdatedByFullName: types.StringValue(routeAggregationChangeLog.GetUpdatedByFullName()),
		UpdatedByEmail:    types.StringValue(routeAggregationChangeLog.GetUpdatedByEmail()),
		UpdatedDateTime:   types.StringValue(routeAggregationChangeLog.GetUpdatedDateTime().Format(TIMEFORMAT)),
		DeletedBy:         types.StringValue(routeAggregationChangeLog.GetDeletedBy()),
		DeletedByFullName: types.StringValue(routeAggregationChangeLog.GetDeletedByFullName()),
		DeletedByEmail:    types.StringValue(routeAggregationChangeLog.GetDeletedByEmail()),
		DeletedDateTime:   types.StringValue(routeAggregationChangeLog.GetDeletedDateTime().Format(TIMEFORMAT)),
	}
	*changeLog = fwtypes.NewObjectValueOf[ChangeLogModel](ctx, &changeLogModel)
	return diag
}
