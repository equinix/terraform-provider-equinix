package equinix

import (
	neaccount "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/account"
	neacltemplate "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/acl_template"
	nebgp "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/bgp"
	nedevice "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/device"
	nedevicelink "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/device_link"
	nedevicesoftware "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/device_software"
	nedevicetype "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/device_type"
	nefile "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/file"
	nedeviceplatform "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/platform"
	nesshkey "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/ssh_key"
	nesshuser "github.com/equinix/terraform-provider-equinix/internal/resources/networkedge/ssh_user"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func networkEdgeDatasources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_network_account":         neaccount.DataSource(),
		"equinix_network_device":          nedevice.DataSource(),
		"equinix_network_device_type":     nedevicetype.DataSource(),
		"equinix_network_device_software": nedevicesoftware.DataSource(),
		"equinix_network_device_platform": nedeviceplatform.DataSource(),
	}
}

func networkEdgeResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"equinix_network_device":       nedevice.Resource(),
		"equinix_network_ssh_user":     nesshuser.Resource(),
		"equinix_network_bgp":          nebgp.Resource(),
		"equinix_network_ssh_key":      nesshkey.Resource(),
		"equinix_network_acl_template": neacltemplate.Resource(),
		"equinix_network_device_link":  nedevicelink.Resource(),
		"equinix_network_file":         nefile.Resource(),
	}
}
