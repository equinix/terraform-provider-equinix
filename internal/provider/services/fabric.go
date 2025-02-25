package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection_route_aggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/metro"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_aggregation_rule"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/routeaggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"
	streamattachment "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream_attachment"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{
		connection_route_aggregation.NewResource,
		routeaggregation.NewResource,
		route_aggregation_rule.NewResource,
		stream.NewResource,
		streamattachment.NewResource,
	}
}

func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		connection_route_aggregation.NewDataSourceByConnectionRouteAggregationID,
		connection_route_aggregation.NewDataSourceAllConnectionRouteAggregations,
		metro.NewDataSourceMetroCode,
		metro.NewDataSourceMetros,
		routeaggregation.NewDataSourceByRouteAggregationID,
		routeaggregation.NewDataSourceAllRouteAggregation,
		route_aggregation_rule.NewDataSourceByRouteAggregationRuleID,
		route_aggregation_rule.NewDataSourceAllRouteAggregationRule,
		stream.NewDataSourceByStreamID,
		stream.NewDataSourceAllStreams,
		streamattachment.NewDataSourceAllStreamAttachments,
		streamattachment.NewDataSourceByIDs,
	}
}
