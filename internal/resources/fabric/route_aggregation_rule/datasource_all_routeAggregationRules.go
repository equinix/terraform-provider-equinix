package route_aggregation_rule

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
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllRouteAggregationRulesSchema(ctx)
}

func (r *DataSourceAllRouteAggregationRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DatsSourceAllRouteAggregationRulesModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	routeAggregationId := data.RouteAggregationId.ValueString()

	var tfpagination PaginationModel
	diags := data.Pagination.As(ctx, &tfpagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return
	}
	offset := tfpagination.Offset.ValueInt32()
	limit := tfpagination.Limit.ValueInt32()
	if limit == 0 {
		limit = 20
	}

	routeAggregations, _, err := client.RouteAggregationRulesApi.GetRouteAggregationRules(ctx, routeAggregationId).Limit(limit).Offset(offset).Execute()

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
