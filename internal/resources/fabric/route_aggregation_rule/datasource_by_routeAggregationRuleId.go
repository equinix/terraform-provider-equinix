package route_aggregation_rule

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSourceByRouteAggregationRuleID() datasource.DataSource {
	return &DataSourceByRouteAggregationRuleID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_route_aggregation_rule",
			},
		),
	}
}

type DataSourceByRouteAggregationRuleID struct {
	framework.BaseDataSource
}

func (r *DataSourceByRouteAggregationRuleID) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleRouteAggregationRuleSchema(ctx)
}

func (r *DataSourceByRouteAggregationRuleID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceByIdModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	routeAggregationRuleID := data.RouteAggregationRuleId.ValueString()
	routeAggregationId := data.RouteAggregationID.ValueString()

	routeAggregation, _, err := client.RouteAggregationRulesApi.GetRouteAggregationRuleByUuid(ctx, routeAggregationId, routeAggregationRuleID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving route aggregation rule data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, routeAggregation)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)

}
