package connectionrouteaggregation

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSourceByConnectionRouteAggregationID() datasource.DataSource {
	return &DataSourceByConnectionRouteAggregationID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_connection_route_aggregation",
			},
		),
	}
}

type DataSourceByConnectionRouteAggregationID struct {
	framework.BaseDataSource
}

func (r *DataSourceByConnectionRouteAggregationID) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleConnectionRouteAggregationSchema(ctx)
}

func (r *DataSourceByConnectionRouteAggregationID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data dataSourceByIDModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	routeAggregationID := data.RouteAggregationID.ValueString()
	connectionID := data.ConnectionID.ValueString()

	routeAggregation, _, err := client.RouteAggregationsApi.GetConnectionRouteAggregationByUuid(ctx, routeAggregationID, connectionID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving connection route aggregation data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, routeAggregation)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)

}
