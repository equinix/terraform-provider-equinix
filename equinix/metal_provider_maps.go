package equinix

import (
	metal_device "github.com/equinix/terraform-provider-equinix/internal/resources/metal/device"
	metal_port "github.com/equinix/terraform-provider-equinix/internal/resources/metal/port"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/virtual_circuit"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vrf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func metalDatasources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_metal_hardware_reservation": dataSourceMetalHardwareReservation(),
		"equinix_metal_metro":                dataSourceMetalMetro(),
		"equinix_metal_facility":             dataSourceMetalFacility(),
		"equinix_metal_ip_block_ranges":      dataSourceMetalIPBlockRanges(),
		"equinix_metal_precreated_ip_block":  dataSourceMetalPreCreatedIPBlock(),
		"equinix_metal_operating_system":     dataSourceOperatingSystem(),
		"equinix_metal_spot_market_price":    dataSourceSpotMarketPrice(),
		"equinix_metal_device":               metal_device.DataSource(),
		"equinix_metal_devices":              metal_device.ListDataSource(),
		"equinix_metal_device_bgp_neighbors": dataSourceMetalDeviceBGPNeighbors(),
		"equinix_metal_plans":                dataSourceMetalPlans(),
		"equinix_metal_port":                 metal_port.DataSource(),
		"equinix_metal_reserved_ip_block":    dataSourceMetalReservedIPBlock(),
		"equinix_metal_spot_market_request":  dataSourceMetalSpotMarketRequest(),
		"equinix_metal_virtual_circuit":      virtual_circuit.DataSource(),
		"equinix_metal_vrf":                  vrf.DataSource(),
	}
}

func metalResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_metal_user_api_key":         resourceMetalUserAPIKey(),
		"equinix_metal_project_api_key":      resourceMetalProjectAPIKey(),
		"equinix_metal_device":               metal_device.Resource(),
		"equinix_metal_device_network_type":  resourceMetalDeviceNetworkType(),
		"equinix_metal_port":                 metal_port.Resource(),
		"equinix_metal_reserved_ip_block":    resourceMetalReservedIPBlock(),
		"equinix_metal_ip_attachment":        resourceMetalIPAttachment(),
		"equinix_metal_spot_market_request":  resourceMetalSpotMarketRequest(),
		"equinix_metal_virtual_circuit":      virtual_circuit.Resource(),
		"equinix_metal_vrf":                  vrf.Resource(),
		"equinix_metal_bgp_session":          resourceMetalBGPSession(),
		"equinix_metal_port_vlan_attachment": resourceMetalPortVlanAttachment(),
	}
}
