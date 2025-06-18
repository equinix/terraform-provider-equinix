// Package services for registering service-specific test sweepers
package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/streamalertrule"

	fabric_cloud_router "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/cloud_router"
	fabric_connection "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection"
	fabric_route_filter "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter"
	fabric_route_aggregation "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/routeaggregation"
	fabric_stream "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"

	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/network"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/serviceprofile"
)

// AddFabricTestSweepers registers test sweepers for Fabric resources
func AddFabricTestSweepers() {
	fabric_cloud_router.AddTestSweeper()
	fabric_connection.AddTestSweeper()
	fabric_route_filter.AddTestSweeper()
	fabric_route_aggregation.AddTestSweeper()
	fabric_stream.AddTestSweeper()
	network.AddTestSweeper()
	serviceprofile.AddTestSweeper()
	streamalertrule.AddTestSweeper()
}
