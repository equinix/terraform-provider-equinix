package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func networkEdgeDatasources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_network_account":         dataSourceNetworkAccount(),
		"equinix_network_device":          dataSourceNetworkDevice(),
		"equinix_network_device_type":     dataSourceNetworkDeviceType(),
		"equinix_network_device_software": dataSourceNetworkDeviceSoftware(),
		"equinix_network_device_platform": dataSourceNetworkDevicePlatform(),
	}
}

func networkEdgeResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_network_device":       resourceNetworkDevice(),
		"equinix_network_ssh_user":     resourceNetworkSSHUser(),
		"equinix_network_bgp":          resourceNetworkBGP(),
		"equinix_network_ssh_key":      resourceNetworkSSHKey(),
		"equinix_network_acl_template": resourceNetworkACLTemplate(),
		"equinix_network_device_link":  resourceNetworkDeviceLink(),
		"equinix_network_file":         resourceNetworkFile(),
	}
}
