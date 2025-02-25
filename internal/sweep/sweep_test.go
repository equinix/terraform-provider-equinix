package sweep_test

import (
	"testing"

	fabric_cloud_router "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/cloud_router"
	fabric_connection "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection"
	fabric_route_filter "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter"
	fabric_route_aggregation "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/routeaggregation"
	fabric_stream "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"

	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/network"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/connection"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/device"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/ssh_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/user_api_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/virtual_circuit"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vrf"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	// Register sweepers for individual resource packages
	addTestSweepers()

	resource.TestMain(m)
}

func addTestSweepers() {
	connection.AddTestSweeper()
	device.AddTestSweeper()
	fabric_cloud_router.AddTestSweeper()
	fabric_connection.AddTestSweeper()
	fabric_route_filter.AddTestSweeper()
	fabric_route_aggregation.AddTestSweeper()
	fabric_stream.AddTestSweeper()
	network.AddTestSweeper()
	organization.AddTestSweeper()
	project.AddTestSweeper()
	ssh_key.AddTestSweeper()
	user_api_key.AddTestSweeper()
	virtual_circuit.AddTestSweeper()
	vlan.AddTestSweeper()
	vrf.AddTestSweeper()
}
