package equinix

import (
	"ecx-go/v3"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	endpointEnvVar     = "EQUINIX_API_ENDPOINT"
	clientIDEnvVar     = "EQUINIX_API_CLIENTID"
	clientSecretEnvVar = "EQUINIX_API_CLIENTSECRET"
)

//Provider returns Equinix terraform ResourceProvider
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					endpointEnvVar,
				}, nil),
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					clientIDEnvVar,
				}, nil),
			},
			"client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					clientSecretEnvVar,
				}, nil),
			},
			"request_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"equinix_ecx_l2_connection":     resourceECXL2Connection(),
			"equinix_ecx_l2_serviceprofile": resourceECXL2ServiceProfile(),
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

func hasECXErrorCode(errors []ecx.Error, code string) bool {
	for _, err := range errors {
		if err.ErrorCode == code {
			return true
		}
	}
	return false
}
