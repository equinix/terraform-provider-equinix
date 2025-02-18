package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection_route_aggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_aggregation_rule"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{
		stream.NewResource,
		route_aggregation_rule.NewResource,
		connection_route_aggregation.NewResource,
	}
}

func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		stream.NewDataSourceByStreamID,
		stream.NewDataSourceAllStreams,
		route_aggregation_rule.NewDataSourceByRouteAggregationRuleID,
		route_aggregation_rule.NewDataSourceAllRouteAggregationRule,
		connection_route_aggregation.NewDataSourceByConnectionRouteAggregationID,
		connection_route_aggregation.NewDataSourceAllConnectionRouteAggregations,
	}
}
