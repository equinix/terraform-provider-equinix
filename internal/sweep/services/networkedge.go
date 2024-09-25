// Package services for registering service-specific test sweepers
package services

import (
	ne_acl_template "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/acl_template"
	ne_device "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/device"
	ne_device_link "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/device_link"
	ne_ssh_key "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/ssh_key"
	ne_ssh_user "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/ssh_user"
)

func AddNetworkEdgeTestSweepers() {
	ne_device.AddTestSweeper()
	ne_device_link.AddTestSweeper()
	ne_acl_template.AddTestSweeper()
	ne_ssh_key.AddTestSweeper()
	ne_ssh_user.AddTestSweeper()
}
