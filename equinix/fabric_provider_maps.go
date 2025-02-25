package equinix

import (
	fabric_connection "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection"
	fabric_connection_route_filter "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection_route_filter"
	fabric_market_place_subscription "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/marketplace"
	fabric_network "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/network"
	fabric_route_filter "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter"
	fabric_route_filter_rule "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter_rule"
	fabric_service_token "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/servicetoken"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func fabricDatasources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_fabric_routing_protocol":          dataSourceRoutingProtocol(),
		"equinix_fabric_connection":                fabric_connection.DataSource(),
		"equinix_fabric_connections":               fabric_connection.DataSourceSearch(),
		"equinix_fabric_connection_route_filter":   fabric_connection_route_filter.DataSource(),
		"equinix_fabric_connection_route_filters":  fabric_connection_route_filter.DataSourceGetAllRules(),
		"equinix_fabric_cloud_router":              dataSourceFabricCloudRouter(),
		"equinix_fabric_cloud_routers":             dataSourceFabricGetCloudRouters(),
		"equinix_fabric_market_place_subscription": fabric_market_place_subscription.DataSourceFabricMarketplaceSubscription(),
		"equinix_fabric_network":                   fabric_network.DataSource(),
		"equinix_fabric_networks":                  fabric_network.DataSourceSearch(),
		"equinix_fabric_port":                      dataSourceFabricPort(),
		"equinix_fabric_ports":                     dataSourceFabricGetPortsByName(),
		"equinix_fabric_route_filter":              fabric_route_filter.DataSource(),
		"equinix_fabric_route_filters":             fabric_route_filter.DataSourceSearch(),
		"equinix_fabric_route_filter_rule":         fabric_route_filter_rule.DataSource(),
		"equinix_fabric_route_filter_rules":        fabric_route_filter_rule.DataSourceGetAllRules(),
		"equinix_fabric_service_profile":           dataSourceFabricServiceProfileReadByUuid(),
		"equinix_fabric_service_profiles":          dataSourceFabricSearchServiceProfilesByName(),
		"equinix_fabric_service_token":             fabric_service_token.DataSource(),
		"equinix_fabric_service_tokens":            fabric_service_token.DataSourceSearch(),
	}
}

func fabricResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_fabric_network":                 fabric_network.Resource(),
		"equinix_fabric_cloud_router":            resourceFabricCloudRouter(),
		"equinix_fabric_connection":              fabric_connection.Resource(),
		"equinix_fabric_connection_route_filter": fabric_connection_route_filter.Resource(),
		"equinix_fabric_route_filter":            fabric_route_filter.Resource(),
		"equinix_fabric_route_filter_rule":       fabric_route_filter_rule.Resource(),
		"equinix_fabric_routing_protocol":        resourceFabricRoutingProtocol(),
		"equinix_fabric_service_profile":         resourceFabricServiceProfile(),
		"equinix_fabric_service_token":           fabric_service_token.Resource(),
	}
}
