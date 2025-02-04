package route_aggregation

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSourceByRouteAggregationID() datasource.DataSource {
	return &DataSourceByRouteAggregationID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_route_aggregation",
			},
		),
	}
}

type DataSourceByRouteAggregationID struct {
	framework.BaseDataSource
}

func (r *DataSourceByRouteAggregationID) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleRouteAggregationSchema(ctx)
}

func (r *DataSourceByRouteAggregationID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceByIdModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	routeAggregationID := data.RouteAggregationId.ValueString()

	routeAggregation, _, err := client.RouteAggregationsApi.GetRouteAggregationByUuid(ctx, routeAggregationID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving route aggregation data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, routeAggregation)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)

}
