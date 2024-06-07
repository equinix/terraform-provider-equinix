package equinix

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/equinix/ecx-go/v2"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	fabric_network "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/network"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vrf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Provider returns Equinix terraform *schema.Provider
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(config.EndpointEnvVar, config.DefaultBaseURL),
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				Description:  fmt.Sprintf("The Equinix API base URL to point out desired environment. Defaults to %s", config.DefaultBaseURL),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(config.ClientIDEnvVar, ""),
				Description: "API Consumer Key available under My Apps section in developer portal",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(config.ClientSecretEnvVar, ""),
				Description: "API Consumer secret available under My Apps section in developer portal",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(config.ClientTokenEnvVar, ""),
				Description: "API token from the developer sandbox",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(config.MetalAuthTokenEnvVar, ""),
				Description: "The Equinix Metal API auth key for API operations",
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(config.ClientTimeoutEnvVar, config.DefaultTimeout),
				ValidateFunc: validation.IntAtLeast(1),
				Description:  fmt.Sprintf("The duration of time, in seconds, that the Equinix Platform API Client should wait before canceling an API request.  Defaults to %d", config.DefaultTimeout),
			},
			"response_max_page_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(100),
				Description:  "The maximum number of records in a single response for REST queries that produce paginated responses",
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of retries.",
			},
			"max_retry_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Maximum number of seconds to wait before retrying a request.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"equinix_ecx_port":                   dataSourceECXPort(),
			"equinix_ecx_l2_sellerprofile":       dataSourceECXL2SellerProfile(),
			"equinix_ecx_l2_sellerprofiles":      dataSourceECXL2SellerProfiles(),
			"equinix_fabric_routing_protocol":    dataSourceRoutingProtocol(),
			"equinix_fabric_connection":          dataSourceFabricConnection(),
			"equinix_fabric_cloud_router":        dataSourceFabricCloudRouter(),
			"equinix_fabric_cloud_routers":       dataSourceFabricGetCloudRouters(),
			"equinix_fabric_network":             fabric_network.DataSource(),
			"equinix_fabric_networks":            fabric_network.DataSourceSearch(),
			"equinix_fabric_port":                dataSourceFabricPort(),
			"equinix_fabric_ports":               dataSourceFabricGetPortsByName(),
			"equinix_fabric_service_profile":     dataSourceFabricServiceProfileReadByUuid(),
			"equinix_fabric_service_profiles":    dataSourceFabricSearchServiceProfilesByName(),
			"equinix_network_account":            dataSourceNetworkAccount(),
			"equinix_network_device":             dataSourceNetworkDevice(),
			"equinix_network_device_type":        dataSourceNetworkDeviceType(),
			"equinix_network_device_software":    dataSourceNetworkDeviceSoftware(),
			"equinix_network_device_platform":    dataSourceNetworkDevicePlatform(),
			"equinix_metal_hardware_reservation": dataSourceMetalHardwareReservation(),
			"equinix_metal_metro":                dataSourceMetalMetro(),
			"equinix_metal_facility":             dataSourceMetalFacility(),
			"equinix_metal_ip_block_ranges":      dataSourceMetalIPBlockRanges(),
			"equinix_metal_precreated_ip_block":  dataSourceMetalPreCreatedIPBlock(),
			"equinix_metal_operating_system":     dataSourceOperatingSystem(),
			"equinix_metal_spot_market_price":    dataSourceSpotMarketPrice(),
			"equinix_metal_device":               dataSourceMetalDevice(),
			"equinix_metal_devices":              dataSourceMetalDevices(),
			"equinix_metal_device_bgp_neighbors": dataSourceMetalDeviceBGPNeighbors(),
			"equinix_metal_plans":                dataSourceMetalPlans(),
			"equinix_metal_port":                 dataSourceMetalPort(),
			"equinix_metal_reserved_ip_block":    dataSourceMetalReservedIPBlock(),
			"equinix_metal_spot_market_request":  dataSourceMetalSpotMarketRequest(),
			"equinix_metal_virtual_circuit":      dataSourceMetalVirtualCircuit(),
			"equinix_metal_vrf":                  vrf.DataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"equinix_ecx_l2_connection":          resourceECXL2Connection(),
			"equinix_ecx_l2_connection_accepter": resourceECXL2ConnectionAccepter(),
			"equinix_ecx_l2_serviceprofile":      resourceECXL2ServiceProfile(),
			"equinix_fabric_network":             fabric_network.Resource(),
			"equinix_fabric_cloud_router":        resourceFabricCloudRouter(),
			"equinix_fabric_connection":          resourceFabricConnection(),
			"equinix_fabric_routing_protocol":    resourceFabricRoutingProtocol(),
			"equinix_fabric_service_profile":     resourceFabricServiceProfile(),
			"equinix_network_device":             resourceNetworkDevice(),
			"equinix_network_ssh_user":           resourceNetworkSSHUser(),
			"equinix_network_bgp":                resourceNetworkBGP(),
			"equinix_network_ssh_key":            resourceNetworkSSHKey(),
			"equinix_network_acl_template":       resourceNetworkACLTemplate(),
			"equinix_network_device_link":        resourceNetworkDeviceLink(),
			"equinix_network_file":               resourceNetworkFile(),
			"equinix_metal_user_api_key":         resourceMetalUserAPIKey(),
			"equinix_metal_project_api_key":      resourceMetalProjectAPIKey(),
			"equinix_metal_device":               resourceMetalDevice(),
			"equinix_metal_device_network_type":  resourceMetalDeviceNetworkType(),
			"equinix_metal_port":                 resourceMetalPort(),
			"equinix_metal_reserved_ip_block":    resourceMetalReservedIPBlock(),
			"equinix_metal_ip_attachment":        resourceMetalIPAttachment(),
			"equinix_metal_spot_market_request":  resourceMetalSpotMarketRequest(),
			"equinix_metal_virtual_circuit":      resourceMetalVirtualCircuit(),
			"equinix_metal_vrf":                  vrf.Resource(),
			"equinix_metal_bgp_session":          resourceMetalBGPSession(),
			"equinix_metal_port_vlan_attachment": resourceMetalPortVlanAttachment(),
		},
		ProviderMetaSchema: map[string]*schema.Schema{
			"module_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return configureProvider(ctx, d, provider)
	}
	return provider
}

type providerMeta struct {
	ModuleName string `cty:"module_name"`
}

func configureProvider(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	mrws := d.Get("max_retry_wait_seconds").(int)
	rt := d.Get("request_timeout").(int)

	config := config.Config{
		AuthToken:      d.Get("auth_token").(string),
		BaseURL:        d.Get("endpoint").(string),
		ClientID:       d.Get("client_id").(string),
		ClientSecret:   d.Get("client_secret").(string),
		Token:          d.Get("token").(string),
		RequestTimeout: time.Duration(rt) * time.Second,
		PageSize:       d.Get("response_max_page_size").(int),
		MaxRetries:     d.Get("max_retries").(int),
		MaxRetryWait:   time.Duration(mrws) * time.Second,
	}
	meta := providerMeta{}

	if err := d.GetProviderMeta(&meta); err != nil {
		return nil, diag.FromErr(err)
	}
	config.TerraformVersion = p.TerraformVersion
	if config.TerraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		config.TerraformVersion = "0.11+compatible"
	}

	stopCtx, ok := schema.StopContext(ctx)
	if !ok {
		stopCtx = ctx
	}
	if err := config.Load(stopCtx); err != nil {
		return nil, diag.FromErr(err)
	}
	return &config, nil
}

func stringsFound(source []string, target []string) bool {
	for i := range source {
		if !slices.Contains(target, source[i]) {
			return false
		}
	}
	return true
}

func atLeastOneStringFound(source []string, target []string) bool {
	for i := range source {
		if slices.Contains(target, source[i]) {
			return true
		}
	}
	return false
}

func isEmpty(v interface{}) bool {
	switch v := v.(type) {
	case int:
		return v == 0
	case *int:
		return ecx.IntValue(v) == 0
	case string:
		return v == ""
	case *string:
		return ecx.StringValue(v) == ""
	case nil:
		return true
	default:
		return false
	}
}

func slicesMatch(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	visited := make([]bool, len(s1))
	for i := 0; i < len(s1); i++ {
		found := false
		for j := 0; j < len(s2); j++ {
			if visited[j] {
				continue
			}
			if s1[i] == s2[j] {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func slicesMatchCaseInsensitive(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	visited := make([]bool, len(s1))
	for i := 0; i < len(s1); i++ {
		found := false
		for j := 0; j < len(s2); j++ {
			if visited[j] {
				continue
			}
			if strings.EqualFold(s1[i], s2[j]) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
