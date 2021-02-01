package equinix

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/equinix/ecx-go/v2"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	endpointEnvVar      = "EQUINIX_API_ENDPOINT"
	clientIDEnvVar      = "EQUINIX_API_CLIENTID"
	clientSecretEnvVar  = "EQUINIX_API_CLIENTSECRET"
	clientTimeoutEnvVar = "EQUINIX_API_TIMEOUT"
)

//resourceDataProvider provies interface to schema.ResourceData
//for convenient mocking purposes
type resourceDataProvider interface {
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
	HasChange(key string) bool
	GetChange(key string) (interface{}, interface{})
}

//Provider returns Equinix terraform *schema.Provider
func Provider() *schema.Provider {
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
				DefaultFunc:  schema.EnvDefaultFunc(clientTimeoutEnvVar, 30),
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
			"equinix_network_ssh_key":            resourceNetworkSSHKey(),
			"equinix_network_acl_template":       resourceNetworkACLTemplate(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return configureProvider(ctx, d, provider)
	}
	return provider
}

func configureProvider(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
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
	stopCtx, ok := schema.StopContext(ctx)
	if !ok {
		stopCtx = ctx
	}
	if err := config.Load(stopCtx); err != nil {
		return nil, diag.FromErr(err)
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
	switch v.(type) {
	case int:
		return v.(int) == 0
	case *int:
		return ecx.IntValue(v.(*int)) == 0
	case string:
		return v.(string) == ""
	case *string:
		return ecx.StringValue(v.(*string)) == ""
	case nil:
		return true
	default:
		return false
	}
}

func isEmptyOld(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	if val.IsZero() {
		return true
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.IsZero()
}

func setSchemaValueIfNotEmpty(key string, value interface{}, d *schema.ResourceData) error {
	if !isEmpty(value) {
		return d.Set(key, value)
	}
	return nil
}
