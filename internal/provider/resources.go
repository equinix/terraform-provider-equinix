package provider

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/ecx/l2connection"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ecx/l2connectionaccepter"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ecx/l2serviceprofile"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/bgpsession"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/connection"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/device"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/devicenetworktype"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/gateway"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/ipattachment"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/port"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/portvlanattachment"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/projectapikey"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/projectsshkey"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/reservedipblock"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/spotmarketrequest"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/sshkey"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/userapikey"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/virtualcircuit"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vrf"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ne/acltemplate"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ne/bgp"
	nedevice "github.com/equinix/terraform-provider-equinix/internal/resources/ne/device"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ne/devicelink"
	nesshkey "github.com/equinix/terraform-provider-equinix/internal/resources/ne/sshkey"
	"github.com/equinix/terraform-provider-equinix/internal/resources/ne/sshuser"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_ecx_l2_connection":          l2connection.Resource(),
		"equinix_ecx_l2_connection_accepter": l2connectionaccepter.Resource(),
		"equinix_ecx_l2_serviceprofile":      l2serviceprofile.Resource(),
		"equinix_network_device":             nedevice.Resource(),
		"equinix_network_ssh_user":           sshuser.Resource(),
		"equinix_network_bgp":                bgp.Resource(),
		"equinix_network_ssh_key":            nesshkey.Resource(),
		"equinix_network_acl_template":       acltemplate.Resource(),
		"equinix_network_device_link":        devicelink.Resource(),
		"equinix_metal_user_api_key":         userapikey.Resource(),
		"equinix_metal_project_api_key":      projectapikey.Resource(),
		"equinix_metal_connection":           connection.Resource(),
		"equinix_metal_device":               device.Resource(),
		"equinix_metal_device_network_type":  devicenetworktype.Resource(),
		"equinix_metal_ssh_key":              sshkey.Resource(),
		"equinix_metal_port":                 port.Resource(),
		"equinix_metal_project_ssh_key":      projectsshkey.Resource(),
		"equinix_metal_project":              project.Resource(),
		"equinix_metal_organization":         organization.Resource(),
		"equinix_metal_reserved_ip_block":    reservedipblock.Resource(),
		"equinix_metal_ip_attachment":        ipattachment.Resource(),
		"equinix_metal_spot_market_request":  spotmarketrequest.Resource(),
		"equinix_metal_vlan":                 vlan.Resource(),
		"equinix_metal_virtual_circuit":      virtualcircuit.Resource(),
		"equinix_metal_vrf":                  vrf.Resource(),
		"equinix_metal_bgp_session":          bgpsession.Resource(),
		"equinix_metal_port_vlan_attachment": portvlanattachment.Resource(),
		"equinix_metal_gateway":              gateway.Resource(),
	}
}
