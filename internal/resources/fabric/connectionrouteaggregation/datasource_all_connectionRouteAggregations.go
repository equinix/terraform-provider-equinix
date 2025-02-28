package connectionrouteaggregation

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSourceAllConnectionRouteAggregations() datasource.DataSource {
	return &DataSourceAllConnectionRouteAggregations{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_connection_route_aggregations",
			},
		),
	}
}

type DataSourceAllConnectionRouteAggregations struct {
	framework.BaseDataSource
}

func (r *DataSourceAllConnectionRouteAggregations) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllConnectionRouteAggregationSchema(ctx)
}

func (r *DataSourceAllConnectionRouteAggregations) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data datsSourceAllConnectionRouteAggregationModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	connectionID := data.ConnectionID.ValueString()

	connectionRouteAggregations, _, err := client.RouteAggregationsApi.GetConnectionRouteAggregations(ctx, connectionID).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving connection route aggregations data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, connectionRouteAggregations)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
