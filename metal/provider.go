package metal

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var metalMutexKV = mutexkv.NewMutexKV()

func Provider() terraform.ResourceProvider {

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PACKET_AUTH_TOKEN", nil),
				Description: "The API auth key for API operations.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"metal_ip_block_ranges":      dataSourceMetalIPBlockRanges(),
			"metal_precreated_ip_block":  dataSourceMetalPreCreatedIPBlock(),
			"metal_operating_system":     dataSourceOperatingSystem(),
			"metal_organization":         dataSourceMetalOrganization(),
			"metal_spot_market_price":    dataSourceSpotMarketPrice(),
			"metal_device":               dataSourceMetalDevice(),
			"metal_device_bgp_neighbors": dataSourceMetalDeviceBGPNeighbors(),
			"metal_project":              dataSourceMetalProject(),
			"metal_project_ssh_key":      dataSourceMetalProjectSSHKey(),
			"metal_spot_market_request":  dataSourceMetalSpotMarketRequest(),
			"metal_volume":               dataSourceMetalVolume(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"metal_device":               resourceMetalDevice(),
			"metal_device_network_type":  resourceMetalDeviceNetworkType(),
			"metal_ssh_key":              resourceMetalSSHKey(),
			"metal_project_ssh_key":      resourceMetalProjectSSHKey(),
			"metal_project":              resourceMetalProject(),
			"metal_organization":         resourceMetalOrganization(),
			"metal_volume":               resourceMetalVolume(),
			"metal_volume_attachment":    resourceMetalVolumeAttachment(),
			"metal_reserved_ip_block":    resourceMetalReservedIPBlock(),
			"metal_ip_attachment":        resourceMetalIPAttachment(),
			"metal_spot_market_request":  resourceMetalSpotMarketRequest(),
			"metal_vlan":                 resourceMetalVlan(),
			"metal_bgp_session":          resourceMetalBGPSession(),
			"metal_port_vlan_attachment": resourceMetalPortVlanAttachment(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AuthToken: d.Get("auth_token").(string),
	}
	return config.Client(), nil
}

var resourceDefaultTimeouts = &schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}
