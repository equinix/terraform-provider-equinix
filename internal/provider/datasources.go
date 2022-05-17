package provider

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/ecx/l2sellerprofile"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ecx/l2sellerprofiles"
	ecxport "github.com/equinix/terraform-provider-equinix/internal/resources/ecx/port"

	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/connection"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/device"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/devicebgpneighbors"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/facility"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/gateway"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/hardwarereservation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/ipblockranges"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/metro"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/operatingsystem"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/plans"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/port"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/precreatedipblock"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/projectsshkey"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/reservedipblock"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/spotmarketprice"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/spotmarketrequest"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/virtualcircuit"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vrf"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ne/account"
	nedevice "github.com/equinix/terraform-provider-equinix/internal/resources/ne/device"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ne/devicesoftware"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ne/devicetype"
	platform "github.com/equinix/terraform-provider-equinix/internal/resources/ne/platform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Datasources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_ecx_port":                   ecxport.DataSource(),
		"equinix_ecx_l2_sellerprofile":       l2sellerprofile.DataSource(),
		"equinix_ecx_l2_sellerprofiles":      l2sellerprofiles.DataSource(),
		"equinix_network_account":            account.DataSource(),
		"equinix_network_device":             nedevice.DataSource(),
		"equinix_network_device_type":        devicetype.DataSource(),
		"equinix_network_device_software":    devicesoftware.DataSource(),
		"equinix_network_device_platform":    platform.DataSource(),
		"equinix_metal_hardware_reservation": hardwarereservation.DataSource(),
		"equinix_metal_metro":                metro.DataSource(),
		"equinix_metal_facility":             facility.DataSource(),
		"equinix_metal_connection":           connection.DataSource(),
		"equinix_metal_gateway":              gateway.DataSource(),
		"equinix_metal_ip_block_ranges":      ipblockranges.DataSource(),
		"equinix_metal_precreated_ip_block":  precreatedipblock.DataSource(),
		"equinix_metal_operating_system":     operatingsystem.DataSource(),
		"equinix_metal_organization":         organization.DataSource(),
		"equinix_metal_spot_market_price":    spotmarketprice.DataSource(),
		"equinix_metal_device":               device.DataSource(),
		"equinix_metal_device_bgp_neighbors": devicebgpneighbors.DataSource(),
		"equinix_metal_plans":                plans.DataSource(),
		"equinix_metal_port":                 port.DataSource(),
		"equinix_metal_project":              project.DataSource(),
		"equinix_metal_project_ssh_key":      projectsshkey.DataSource(),
		"equinix_metal_reserved_ip_block":    reservedipblock.DataSource(),
		"equinix_metal_spot_market_request":  spotmarketrequest.DataSource(),
		"equinix_metal_virtual_circuit":      virtualcircuit.DataSource(),
		"equinix_metal_vlan":                 vlan.DataSource(),
		"equinix_metal_vrf":                  vrf.DataSource(),
	}
}
