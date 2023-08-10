package equinix

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/equinix/ecx-go/v2"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	metalMutexKV         = NewMutexKV()
	DeviceNetworkTypes   = []string{"layer3", "hybrid", "layer2-individual", "layer2-bonded"}
	DeviceNetworkTypesHB = []string{"layer3", "hybrid", "hybrid-bonded", "layer2-individual", "layer2-bonded"}
	NetworkTypeList      = strings.Join(DeviceNetworkTypes, ", ")
	NetworkTypeListHB    = strings.Join(DeviceNetworkTypesHB, ", ")
)

const (
	endpointEnvVar       = "EQUINIX_API_ENDPOINT"
	clientIDEnvVar       = "EQUINIX_API_CLIENTID"
	clientSecretEnvVar   = "EQUINIX_API_CLIENTSECRET"
	clientTokenEnvVar    = "EQUINIX_API_TOKEN"
	clientTimeoutEnvVar  = "EQUINIX_API_TIMEOUT"
	metalAuthTokenEnvVar = "METAL_AUTH_TOKEN"
)

// resourceDataProvider provies interface to schema.ResourceData
// for convenient mocking purposes
type resourceDataProvider interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
	HasChange(key string) bool
	GetChange(key string) (interface{}, interface{})
}

// Provider returns Equinix terraform *schema.Provider
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(endpointEnvVar, DefaultBaseURL),
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				Description:  fmt.Sprintf("The Equinix API base URL to point out desired environment. Defaults to %s", DefaultBaseURL),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(clientIDEnvVar, ""),
				Description: "API Consumer Key available under My Apps section in developer portal",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(clientSecretEnvVar, ""),
				Description: "API Consumer secret available under My Apps section in developer portal",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(clientTokenEnvVar, ""),
				Description: "API token from the developer sandbox",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(metalAuthTokenEnvVar, ""),
				Description: "The Equinix Metal API auth key for API operations",
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(clientTimeoutEnvVar, DefaultTimeout),
				ValidateFunc: validation.IntAtLeast(1),
				Description:  fmt.Sprintf("The duration of time, in seconds, that the Equinix Platform API Client should wait before canceling an API request.  Defaults to %d", DefaultTimeout),
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
			"equinix_ecx_port":                   dataSourceECXPort(),
			"equinix_ecx_l2_sellerprofile":       dataSourceECXL2SellerProfile(),
			"equinix_ecx_l2_sellerprofiles":      dataSourceECXL2SellerProfiles(),
			"equinix_fabric_routingprotocol":     dataSourceRoutingProtocol(),
			"equinix_fabric_connection":          dataSourceFabricConnection(),
			"equinix_fabric_gateway":             dataSourceFabricGateway(),
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
			"equinix_metal_connection":           dataSourceMetalConnection(),
			"equinix_metal_gateway":              dataSourceMetalGateway(),
			"equinix_metal_ip_block_ranges":      dataSourceMetalIPBlockRanges(),
			"equinix_metal_precreated_ip_block":  dataSourceMetalPreCreatedIPBlock(),
			"equinix_metal_operating_system":     dataSourceOperatingSystem(),
			"equinix_metal_organization":         dataSourceMetalOrganization(),
			"equinix_metal_spot_market_price":    dataSourceSpotMarketPrice(),
			"equinix_metal_device":               dataSourceMetalDevice(),
			"equinix_metal_device_bgp_neighbors": dataSourceMetalDeviceBGPNeighbors(),
			"equinix_metal_plans":                dataSourceMetalPlans(),
			"equinix_metal_port":                 dataSourceMetalPort(),
			"equinix_metal_project":              dataSourceMetalProject(),
			"equinix_metal_project_ssh_key":      dataSourceMetalProjectSSHKey(),
			"equinix_metal_reserved_ip_block":    dataSourceMetalReservedIPBlock(),
			"equinix_metal_spot_market_request":  dataSourceMetalSpotMarketRequest(),
			"equinix_metal_virtual_circuit":      dataSourceMetalVirtualCircuit(),
			"equinix_metal_vlan":                 dataSourceMetalVlan(),
			"equinix_metal_vrf":                  dataSourceMetalVRF(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"equinix_ecx_l2_connection":          resourceECXL2Connection(),
			"equinix_ecx_l2_connection_accepter": resourceECXL2ConnectionAccepter(),
			"equinix_ecx_l2_serviceprofile":      resourceECXL2ServiceProfile(),
			"equinix_fabric_gateway":             resourceFabricGateway(),
			"equinix_fabric_connection":          resourceFabricConnection(),
			"equinix_fabric_routingprotocol":     resourceFabricRoutingProtocol(),
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
			"equinix_metal_connection":           resourceMetalConnection(),
			"equinix_metal_device":               resourceMetalDevice(),
			"equinix_metal_device_network_type":  resourceMetalDeviceNetworkType(),
			"equinix_metal_ssh_key":              resourceMetalSSHKey(),
			"equinix_metal_organization_member":  resourceMetalOrganizationMember(),
			"equinix_metal_port":                 resourceMetalPort(),
			"equinix_metal_project_ssh_key":      resourceMetalProjectSSHKey(),
			"equinix_metal_project":              resourceMetalProject(),
			"equinix_metal_organization":         resourceMetalOrganization(),
			"equinix_metal_reserved_ip_block":    resourceMetalReservedIPBlock(),
			"equinix_metal_ip_attachment":        resourceMetalIPAttachment(),
			"equinix_metal_spot_market_request":  resourceMetalSpotMarketRequest(),
			"equinix_metal_vlan":                 resourceMetalVlan(),
			"equinix_metal_virtual_circuit":      resourceMetalVirtualCircuit(),
			"equinix_metal_vrf":                  resourceMetalVRF(),
			"equinix_metal_bgp_session":          resourceMetalBGPSession(),
			"equinix_metal_port_vlan_attachment": resourceMetalPortVlanAttachment(),
			"equinix_metal_gateway":              resourceMetalGateway(),
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

	config := Config{
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
	config.terraformVersion = p.TerraformVersion
	if config.terraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		config.terraformVersion = "0.11+compatible"
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

var resourceDefaultTimeouts = &schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}

func expandListToStringList(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = fmt.Sprint(v)
	}
	return result
}

func expandListToInt32List(list []interface{}) []int32 {
	result := make([]int32, len(list))
	for i, v := range list {
		result[i] = int32(v.(int))
	}
	return result
}

func expandSetToStringList(set *schema.Set) []string {
	list := set.List()
	return expandListToStringList(list)
}

func expandInterfaceMapToStringMap(mapIn map[string]interface{}) map[string]string {
	mapOut := make(map[string]string)
	for k, v := range mapIn {
		mapOut[k] = fmt.Sprintf("%v", v)
	}
	return mapOut
}

func hasApplicationErrorCode(errors []rest.ApplicationError, code string) bool {
	for _, err := range errors {
		if err.Code == code {
			return true
		}
	}
	return false
}

func hasModelErrorCode(errors []v4.ModelError, code string) bool {
	for _, err := range errors {
		if err.ErrorCode == code {
			return true
		}
	}
	return false
}

func stringIsMetroCode() schema.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile("^[A-Z]{2}$"), "MetroCode must consist of two capital letters")
}

func stringIsEmailAddress() schema.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile("^[^ @]+@[^ @]+$"), "not valid email address")
}

func stringIsPortDefinition() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^(([0-9]+(,[0-9]+){0,9})|([0-9]+-[0-9]+)|(any))$"),
		"port definition has to be: up to 10 comma sepparated numbers (22,23), range (20-23) or word 'any'")
}

func stringIsSpeedBand() schema.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile("^[0-9]+(MB|GB)$"), "SpeedBand should consist of digit followed by MB or GB")
}

func stringsFound(source []string, target []string) bool {
	for i := range source {
		if !isStringInSlice(source[i], target) {
			return false
		}
	}
	return true
}

func atLeastOneStringFound(source []string, target []string) bool {
	for i := range source {
		if isStringInSlice(source[i], target) {
			return true
		}
	}
	return false
}

func isStringInSlice(needle string, hay []string) bool {
	for i := range hay {
		if needle == hay[i] {
			return true
		}
	}
	return false
}

func getResourceDataChangedKeys(keys []string, d resourceDataProvider) map[string]interface{} {
	changed := make(map[string]interface{})
	for _, key := range keys {
		if v := d.Get(key); v != nil && d.HasChange(key) {
			changed[key] = v
		}
	}
	return changed
}

func getResourceDataListElementChanges(keys []string, listKeyName string, listIndex int, d resourceDataProvider) map[string]interface{} {
	changed := make(map[string]interface{})
	if !d.HasChange(listKeyName) {
		return changed
	}
	old, new := d.GetChange(listKeyName)
	oldList := old.([]interface{})
	newList := new.([]interface{})
	if len(oldList) < listIndex || len(newList) < listIndex {
		return changed
	}
	return getMapChangedKeys(keys, oldList[listIndex].(map[string]interface{}), newList[listIndex].(map[string]interface{}))
}

func getMapChangedKeys(keys []string, old, new map[string]interface{}) map[string]interface{} {
	changed := make(map[string]interface{})
	for _, key := range keys {
		if !reflect.DeepEqual(old[key], new[key]) {
			changed[key] = new[key]
		}
	}
	return changed
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

func isRestNotFoundError(err error) bool {
	if restErr, ok := err.(rest.Error); ok {
		if restErr.HTTPCode == http.StatusNotFound {
			return true
		}
	}
	return false
}

func schemaSetToMap(set *schema.Set) map[int]interface{} {
	transformed := make(map[int]interface{})
	if set != nil {
		list := set.List()
		for i := range list {
			transformed[set.F(list[i])] = list[i]
		}
	}
	return transformed
}
