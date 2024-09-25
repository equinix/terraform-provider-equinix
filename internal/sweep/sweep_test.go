package sweep_test

import (
	"testing"

	fabric_cloud_router "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/cloud_router"
	fabric_connection "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/connection"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/device"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/ssh_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/user_api_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/virtual_circuit"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vrf"
	ne_acl_template "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/acl_template"
	ne_device "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/device"
	ne_device_link "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/device_link"
	ne_ssh_key "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/ssh_key"
	ne_ssh_user "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/ssh_user"

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
	ne_device.AddTestSweeper()
	ne_device_link.AddTestSweeper()
	ne_acl_template.AddTestSweeper()
	ne_ssh_key.AddTestSweeper()
	ne_ssh_user.AddTestSweeper()
	organization.AddTestSweeper()
	project.AddTestSweeper()
	ssh_key.AddTestSweeper()
	user_api_key.AddTestSweeper()
	virtual_circuit.AddTestSweeper()
	vlan.AddTestSweeper()
	vrf.AddTestSweeper()
}
