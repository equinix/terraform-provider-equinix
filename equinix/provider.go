package equinix

import (
	"fmt"
	"regexp"
	"time"

	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	endpointEnvVar      = "EQUINIX_API_ENDPOINT"
	clientIDEnvVar      = "EQUINIX_API_CLIENTID"
	clientSecretEnvVar  = "EQUINIX_API_CLIENTSECRET"
	clientTimeoutEnvVar = "EQUINIX_API_TIMEOUT"
)

//Provider returns Equinix terraform ResourceProvider
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(endpointEnvVar, "https://api.equinix.com"),
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"client_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(clientIDEnvVar, nil),
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"client_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(clientSecretEnvVar, nil),
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				DefaultFunc:  schema.EnvDefaultFunc(clientTimeoutEnvVar, nil),
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"equinix_ecx_port":                dataSourceECXPort(),
			"equinix_ecx_l2_sellerprofile":    dataSourceECXL2SellerProfile(),
			"equinix_network_account":         dataSourceNetworkAccount(),
			"equinix_network_device_type":     dataSourceNetworkDeviceType(),
			"equinix_network_device_software": dataSourceNetworkDeviceSoftware(),
			"equinix_network_device_platform": dataSourceNetworkDevicePlatform(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"equinix_ecx_l2_connection":          resourceECXL2Connection(),
			"equinix_ecx_l2_connection_accepter": resourceECXL2ConnectionAccepter(),
			"equinix_ecx_l2_serviceprofile":      resourceECXL2ServiceProfile(),
			"equinix_network_device":             resourceNetworkDevice(),
			"equinix_network_ssh_user":           resourceNetworkSSHUser(),
			"equinix_network_bgp":                resourceNetworkBGP(),
			"equinix_network_acl_template":       resourceNetworkACLTemplate(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return configureProvider(d, provider)
	}
	return provider
}

func configureProvider(d *schema.ResourceData, p *schema.Provider) (interface{}, error) {
	config := Config{}
	if v, ok := d.GetOk("endpoint"); ok {
		config.BaseURL = v.(string)
	}
	if v, ok := d.GetOk("client_id"); ok {
		config.ClientID = v.(string)
	}
	if v, ok := d.GetOk("client_secret"); ok {
		config.ClientSecret = v.(string)
	}
	if v, ok := d.GetOk("request_timeout"); ok {
		config.RequestTimeout = time.Duration(v.(int)) * time.Second
	}
	if err := config.Load(p.StopContext()); err != nil {
		return nil, err
	}
	return &config, nil
}

func expandListToStringList(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = fmt.Sprint(v)
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

func stringsFound(source []string, target []string) bool {
	for i := range source {
		if !isStringInSlice(source[i], target) {
			return false
		}
	}
	return true
}

func isStringInSlice(needle string, hay []string) bool {
	for i := range hay {
		if needle == hay[i] {
			return true
		}
	}
	return false
}
