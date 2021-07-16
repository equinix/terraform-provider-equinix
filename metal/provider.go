package metal

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	metalMutexKV         = NewMutexKV()
	DeviceNetworkTypes   = []string{"layer3", "hybrid", "layer2-individual", "layer2-bonded"}
	DeviceNetworkTypesHB = []string{"layer3", "hybrid", "hybrid-bonded", "layer2-individual", "layer2-bonded"}
	NetworkTypeList      = strings.Join(DeviceNetworkTypes, ", ")
	NetworkTypeListHB    = strings.Join(DeviceNetworkTypesHB, ", ")
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"METAL_AUTH_TOKEN",
					"PACKET_AUTH_TOKEN",
				}, nil),
				Description: "The API auth key for API operations.",
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"max_retry_wait_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"metal_hardware_reservation": dataSourceMetalHardwareReservation(),
			"metal_metro":                dataSourceMetalMetro(),
			"metal_facility":             dataSourceMetalFacility(),
			"metal_connection":           dataSourceMetalConnection(),
			"metal_ip_block_ranges":      dataSourceMetalIPBlockRanges(),
			"metal_precreated_ip_block":  dataSourceMetalPreCreatedIPBlock(),
			"metal_operating_system":     dataSourceOperatingSystem(),
			"metal_organization":         dataSourceMetalOrganization(),
			"metal_spot_market_price":    dataSourceSpotMarketPrice(),
			"metal_device":               dataSourceMetalDevice(),
			"metal_device_bgp_neighbors": dataSourceMetalDeviceBGPNeighbors(),
			"metal_port":                 dataSourceMetalPort(),
			"metal_project":              dataSourceMetalProject(),
			"metal_project_ssh_key":      dataSourceMetalProjectSSHKey(),
			"metal_reserved_ip_block":    dataSourceMetalReservedIPBlock(),
			"metal_spot_market_request":  dataSourceMetalSpotMarketRequest(),
			"metal_volume":               dataSourceMetalVolume(),
			"metal_virtual_circuit":      dataSourceMetalVirtualCircuit(),
			"metal_vlan":                 dataSourceMetalVlan(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"metal_user_api_key":         resourceMetalUserAPIKey(),
			"metal_project_api_key":      resourceMetalProjectAPIKey(),
			"metal_connection":           resourceMetalConnection(),
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
			"metal_virtual_circuit":      resourceMetalVirtualCircuit(),
			"metal_bgp_session":          resourceMetalBGPSession(),
			"metal_port_vlan_attachment": resourceMetalPortVlanAttachment(),
		},
	}

	provider.ConfigureFunc = providerConfigure(provider.TerraformVersion)

	return provider
}

func providerConfigure(tfVersion string) func(d *schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		mrws := d.Get("max_retry_wait_seconds").(int)
		config := Config{
			AuthToken:    d.Get("auth_token").(string),
			MaxRetries:   d.Get("max_retries").(int),
			MaxRetryWait: time.Duration(mrws) * time.Second,
		}
		config.terraformVersion = tfVersion
		if config.terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			config.terraformVersion = "0.11+compatible"
		}
		return config.Client(), nil
	}
}

var resourceDefaultTimeouts = &schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}
