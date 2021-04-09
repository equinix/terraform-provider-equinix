package equinix

import (
	"context"
	"fmt"
	"net/http"
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
				Description:  "The Equinix API base URL to point out desired environment. Defaults to https://api.equinix.com",
			},
			"client_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(clientIDEnvVar, nil),
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "API Consumer Key available under My Apps section in developer portal",
			},
			"client_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(clientSecretEnvVar, nil),
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "API Consumer secret available under My Apps section in developer portal",
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(clientTimeoutEnvVar, 30),
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The duration of time, in seconds, that the Equinix Platform API Client should wait before canceling an API request",
			},
			"response_max_page_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(100),
				Description:  "The maximum number of records in a single response for REST queries that produce paginated responses",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"equinix_ecx_port":                dataSourceECXPort(),
			"equinix_ecx_l2_sellerprofile":    dataSourceECXL2SellerProfile(),
			"equinix_ecx_l2_sellerprofiles":   dataSourceECXL2SellerProfiles(),
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
			"equinix_network_device_link":        resourceNetworkDeviceLink(),
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
	if v, ok := d.GetOk("response_max_page_size"); ok {
		config.PageSize = v.(int)
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
