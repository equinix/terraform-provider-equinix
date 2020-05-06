package equinix

import (
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
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					endpointEnvVar,
				}, nil),
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					clientIDEnvVar,
				}, nil),
			},
			"client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					clientSecretEnvVar,
				}, nil),
			},
			"request_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"equinix_ecx_l2_connection": resourceECXL2Connection(),
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
