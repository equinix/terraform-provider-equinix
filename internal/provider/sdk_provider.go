package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/metal_connection"

	"github.com/equinix/terraform-provider-equinix/internal/config"
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
			"equinix_ecx_port":                   equinix.DataSourceECXPort(),
			"equinix_ecx_l2_sellerprofile":       equinix.DataSourceECXL2SellerProfile(),
			"equinix_ecx_l2_sellerprofiles":      equinix.DataSourceECXL2SellerProfiles(),
			"equinix_fabric_routing_protocol":    equinix.DataSourceRoutingProtocol(),
			"equinix_fabric_connection":          equinix.DataSourceFabricConnection(),
			"equinix_fabric_cloud_router":        equinix.DataSourceCloudRouter(),
			"equinix_fabric_port":                equinix.DataSourceFabricPort(),
			"equinix_fabric_ports":               equinix.DataSourceFabricGetPortsByName(),
			"equinix_fabric_service_profile":     equinix.DataSourceFabricServiceProfileReadByUuid(),
			"equinix_fabric_service_profiles":    equinix.DataSourceFabricSearchServiceProfilesByName(),
			"equinix_network_account":            equinix.DataSourceNetworkAccount(),
			"equinix_network_device":             equinix.DataSourceNetworkDevice(),
			"equinix_network_device_type":        equinix.DataSourceNetworkDeviceType(),
			"equinix_network_device_software":    equinix.DataSourceNetworkDeviceSoftware(),
			"equinix_network_device_platform":    equinix.DataSourceNetworkDevicePlatform(),
			"equinix_metal_hardware_reservation": equinix.DataSourceMetalHardwareReservation(),
			"equinix_metal_metro":                equinix.DataSourceMetalMetro(),
			"equinix_metal_facility":             equinix.DataSourceMetalFacility(),
			"equinix_metal_connection":           metal_connection.DataSource(),
			"equinix_metal_gateway":              equinix.DataSourceMetalGateway(),
			"equinix_metal_ip_block_ranges":      equinix.DataSourceMetalIPBlockRanges(),
			"equinix_metal_precreated_ip_block":  equinix.DataSourceMetalPreCreatedIPBlock(),
			"equinix_metal_operating_system":     equinix.DataSourceOperatingSystem(),
			"equinix_metal_organization":         equinix.DataSourceMetalOrganization(),
			"equinix_metal_spot_market_price":    equinix.DataSourceSpotMarketPrice(),
			"equinix_metal_device":               equinix.DataSourceMetalDevice(),
			"equinix_metal_devices":              equinix.DataSourceMetalDevices(),
			"equinix_metal_device_bgp_neighbors": equinix.DataSourceMetalDeviceBGPNeighbors(),
			"equinix_metal_plans":                equinix.DataSourceMetalPlans(),
			"equinix_metal_port":                 equinix.DataSourceMetalPort(),
			"equinix_metal_project":              equinix.DataSourceMetalProject(),
			"equinix_metal_reserved_ip_block":    equinix.DataSourceMetalReservedIPBlock(),
			"equinix_metal_spot_market_request":  equinix.DataSourceMetalSpotMarketRequest(),
			"equinix_metal_virtual_circuit":      equinix.DataSourceMetalVirtualCircuit(),
			"equinix_metal_vlan":                 equinix.DataSourceMetalVlan(),
			"equinix_metal_vrf":                  equinix.DataSourceMetalVRF(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"equinix_ecx_l2_connection":          equinix.ResourceECXL2Connection(),
			"equinix_ecx_l2_connection_accepter": equinix.ResourceECXL2ConnectionAccepter(),
			"equinix_ecx_l2_serviceprofile":      equinix.ResourceECXL2ServiceProfile(),
			"equinix_fabric_cloud_router":        equinix.ResourceCloudRouter(),
			"equinix_fabric_connection":          equinix.ResourceFabricConnection(),
			"equinix_fabric_routing_protocol":    equinix.ResourceFabricRoutingProtocol(),
			"equinix_fabric_service_profile":     equinix.ResourceFabricServiceProfile(),
			"equinix_network_device":             equinix.ResourceNetworkDevice(),
			"equinix_network_ssh_user":           equinix.ResourceNetworkSSHUser(),
			"equinix_network_bgp":                equinix.ResourceNetworkBGP(),
			"equinix_network_ssh_key":            equinix.ResourceNetworkSSHKey(),
			"equinix_network_acl_template":       equinix.ResourceNetworkACLTemplate(),
			"equinix_network_device_link":        equinix.ResourceNetworkDeviceLink(),
			"equinix_network_file":               equinix.ResourceNetworkFile(),
			"equinix_metal_user_api_key":         equinix.ResourceMetalUserAPIKey(),
			"equinix_metal_project_api_key":      equinix.ResourceMetalProjectAPIKey(),
			"equinix_metal_connection":           metal_connection.Resource(),
			"equinix_metal_device":               equinix.ResourceMetalDevice(),
			"equinix_metal_device_network_type":  equinix.ResourceMetalDeviceNetworkType(),
			"equinix_metal_organization_member":  equinix.ResourceMetalOrganizationMember(),
			"equinix_metal_port":                 equinix.ResourceMetalPort(),
			"equinix_metal_project":              equinix.ResourceMetalProject(),
			"equinix_metal_organization":         equinix.ResourceMetalOrganization(),
			"equinix_metal_reserved_ip_block":    equinix.ResourceMetalReservedIPBlock(),
			"equinix_metal_ip_attachment":        equinix.ResourceMetalIPAttachment(),
			"equinix_metal_spot_market_request":  equinix.ResourceMetalSpotMarketRequest(),
			"equinix_metal_vlan":                 equinix.ResourceMetalVlan(),
			"equinix_metal_virtual_circuit":      equinix.ResourceMetalVirtualCircuit(),
			"equinix_metal_vrf":                  equinix.ResourceMetalVRF(),
			"equinix_metal_bgp_session":          equinix.ResourceMetalBGPSession(),
			"equinix_metal_port_vlan_attachment": equinix.ResourceMetalPortVlanAttachment(),
			"equinix_metal_gateway":              equinix.ResourceMetalGateway(),
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
