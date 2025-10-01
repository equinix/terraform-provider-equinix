// Package services for registering service-specific test sweepers
package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/connection"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/device"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/ssh_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/user_api_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/virtualcircuit"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vrf"
)

// AddMetalTestSweepers registers test sweepers for Fabric resources
func AddMetalTestSweepers() {
	connection.AddTestSweeper()
	device.AddTestSweeper()
	organization.AddTestSweeper()
	project.AddTestSweeper()
	ssh_key.AddTestSweeper()
	user_api_key.AddTestSweeper()
	virtualcircuit.AddTestSweeper()
	vlan.AddTestSweeper()
	vrf.AddTestSweeper()
}
