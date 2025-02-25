package routeaggregationrule

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewDataSourceAllRouteAggregationRule() datasource.DataSource {
	return &DataSourceAllRouteAggregationRules{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_route_aggregation_rules",
			},
		),
	}
}

type DataSourceAllRouteAggregationRules struct {
	framework.BaseDataSource
}

func (r *DataSourceAllRouteAggregationRules) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllRouteAggregationRulesSchema(ctx)
}

func (r *DataSourceAllRouteAggregationRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data datsSourceAllRouteAggregationRulesModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	routeAggregationID := data.RouteAggregationID.ValueString()
	var tfpagination paginationModel
	if !data.Pagination.IsNull() && !data.Pagination.IsUnknown() {
		diags := data.Pagination.As(ctx, &tfpagination, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			response.Diagnostics.Append(diags...)
			return
		}
	}

	offset := tfpagination.Offset.ValueInt32()
	limit := tfpagination.Limit.ValueInt32()
	if limit == 0 {
		limit = 20
	}
	routeAggregationRequest := client.RouteAggregationRulesApi.GetRouteAggregationRules(ctx, routeAggregationID)
	if !tfpagination.Limit.IsNull() {
		routeAggregationRequest = routeAggregationRequest.Limit(limit)
	}
	if !tfpagination.Offset.IsNull() {
		routeAggregationRequest = routeAggregationRequest.Offset(offset)
	}
	routeAggregations, _, err := routeAggregationRequest.Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving route aggregation rules data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, routeAggregations)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
