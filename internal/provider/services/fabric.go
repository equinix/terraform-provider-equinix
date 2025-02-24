package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/metro"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_aggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{
		stream.NewResource,
		route_aggregation.NewResource,
	}
}

func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		metro.NewDataSourceMetroCode,
		metro.NewDataSourceMetros,
		route_aggregation.NewDataSourceByRouteAggregationID,
		route_aggregation.NewDataSourceAllRouteAggregation,
		stream.NewDataSourceByStreamID,
		stream.NewDataSourceAllStreams,
	}
}
