package equinix

import (
	"fmt"
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
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"equinix_ecx_port":             dataSourceECXPort(),
			"equinix_ecx_l2_sellerprofile": dataSourceECXL2SellerProfile(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"equinix_ecx_l2_connection":          resourceECXL2Connection(),
			"equinix_ecx_l2_connection_accepter": resourceECXL2ConnectionAccepter(),
			"equinix_ecx_l2_serviceprofile":      resourceECXL2ServiceProfile(),
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

func expandSetToStringList(set *schema.Set) []string {
	list := set.List()
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = fmt.Sprint(v)
	}
	return result
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

func isStringInSlice(needle string, hay []string) bool {
	for i := range hay {
		if needle == hay[i] {
			return true
		}
	}
	return false
}
