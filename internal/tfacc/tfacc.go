package tfacc

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	// internal/provider can not be imported to prevent cyclic import

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const TestResourcePrefix = "tfacc"

const (
	PriPortEnvVar              = "TF_ACC_FABRIC_PRI_PORT_NAME"
	SecPortEnvVar              = "TF_ACC_FABRIC_SEC_PORT_NAME"
	AwsSpEnvVar                = "TF_ACC_FABRIC_L2_AWS_SP_NAME"
	AwsAuthKeyEnvVar           = "TF_ACC_FABRIC_L2_AWS_ACCOUNT_ID"
	AzureSpEnvVar              = "TF_ACC_FABRIC_L2_AZURE_SP_NAME"
	AzureXRServiceKeyEnvVar    = "TF_ACC_FABRIC_L2_AZURE_XROUTE_SERVICE_KEY"
	GcpOneSpEnvVar             = "TF_ACC_FABRIC_L2_GCP1_SP_NAME"
	GcpOneConnServiceKeyEnvVar = "TF_ACC_FABRIC_L2_GCP1_INTERCONN_SERVICE_KEY"
	GcpTwoSpEnvVar             = "TF_ACC_FABRIC_L2_GCP2_SP_NAME"
	GcpTwoConnServiceKeyEnvVar = "TF_ACC_FABRIC_L2_GCP2_INTERCONN_SERVICE_KEY"

	NEDeviceAccountNameEnvVar          = "TF_ACC_NETWORK_DEVICE_BILLING_ACCOUNT_NAME"
	NEDeviceSecondaryAccountNameEnvVar = "TF_ACC_NETWORK_DEVICE_SECONDARY_BILLING_ACCOUNT_NAME"
	NEDeviceMetroEnvVar                = "TF_ACC_NETWORK_DEVICE_METRO"
	NEDeviceSecondaryMetroEnvVar       = "TF_ACC_NETWORK_DEVICE_SECONDARY_METRO"
	NEDeviceCSRSDWANLicenseFileEnvVar  = "TF_ACC_NETWORK_DEVICE_CSRSDWAN_LICENSE_FILE"
	NEDeviceVSRXLicenseFileEnvVar      = "TF_ACC_NETWORK_DEVICE_VSRX_LICENSE_FILE"
	NEDeviceVersaController1EnvVar     = "TF_ACC_NETWORK_DEVICE_VERSA_CONTROLLER1"
	NEDeviceVersaController2EnvVar     = "TF_ACC_NETWORK_DEVICE_VERSA_CONTROLLER2"
	NEDeviceVersaLocalIDEnvVar         = "TF_ACC_NETWORK_DEVICE_VERSA_LOCALID"
	NEDeviceVersaRemoteIDEnvVar        = "TF_ACC_NETWORK_DEVICE_VERSA_REMOTEID"
	NEDeviceVersaSerialNumberEnvVar    = "TF_ACC_NETWORK_DEVICE_VERSA_SERIAL"
	NEDeviceCGENIXLicenseKeyEnvVar     = "TF_ACC_NETWORK_DEVICE_CGENIX_LICENSE_KEY"
	NEDeviceCGENIXLicenseSecretEnvVar  = "TF_ACC_NETWORK_DEVICE_CGENIX_LICENSE_SECRET"
	NEDevicePANWLicenseTokenEnvVar     = "TF_ACC_NETWORK_DEVICE_PANW_LICENSE_TOKEN"
)

var (
	AccProviders map[string]*schema.Provider
	AccProvider  *schema.Provider
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Test helper functions
//_______________________________________________________________________

func PreCheck(t *testing.T) {
	var err error

	if _, err = getFromEnv(config.ClientTokenEnvVar); err != nil {
		_, err = getFromEnv(config.ClientIDEnvVar)
		if err == nil {
			_, err = getFromEnv(config.ClientSecretEnvVar)
		}
	}

	if err == nil {
		_, err = getFromEnv(config.MetalAuthTokenEnvVar)
	}

	if err != nil {
		t.Fatalf("To run acceptance tests, one of '%s' or pair '%s' - '%s' must be set for Equinix Fabric and Network Edge, and '%s' for Equinix Metal",
			config.ClientTokenEnvVar, config.ClientIDEnvVar, config.ClientSecretEnvVar, config.MetalAuthTokenEnvVar)
	}
}

type AccConfigFn func(*TestAccConfig)

func NewTestAccConfig(ctx map[string]interface{}, fns ...AccConfigFn) *TestAccConfig {
	conf := &TestAccConfig{
		ctx:    ctx,
		config: "",
	}
	for _, fn := range fns {
		fn(conf)
	}
	return conf
}

func (t *TestAccConfig) Context() map[string]interface{} {
	return t.ctx
}

func (t *TestAccConfig) Use(fn func(map[string]interface{}) string) {
	t.Append(fn(t.ctx))
}

func (t *TestAccConfig) Append(s string) {
	t.config += s
}

func (t *TestAccConfig) Build() string {
	return t.config
}

func NPrintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		var strVal string
		switch val.(type) {
		case []string:
			r := regexp.MustCompile(`" "`)
			strVal = r.ReplaceAllString(fmt.Sprintf("%q", val), `", "`)
		default:
			strVal = fmt.Sprintf("%v", val)
		}
		format = strings.Replace(format, "%{"+key+"}", strVal, -1)
	}
	return format
}

func getFromEnv(varName string) (string, error) {
	if v := os.Getenv(varName); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("environmental variable '%s' is not set", varName)
}

func getFromEnvDefault(varName string, defaultValue string) string {
	if v := os.Getenv(varName); v != "" {
		return v
	}
	return defaultValue
}

type TestAccConfig struct {
	ctx    map[string]interface{}
	config string
}

func init() {
	AccProvider = provider.Provider()
	AccProviders = map[string]*schema.Provider{
		"equinix": AccProvider,
	}
}

func SharedConfigForRegion(region string) (*config.Config, error) {
	endpoint := getFromEnvDefault(config.EndpointEnvVar, config.DefaultBaseURL)
	clientToken := getFromEnvDefault(config.ClientTokenEnvVar, "")
	clientID := getFromEnvDefault(config.ClientIDEnvVar, "")
	clientSecret := getFromEnvDefault(config.ClientSecretEnvVar, "")
	clientTimeout := getFromEnvDefault(config.ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value of '%s' env variable to int", config.ClientTimeoutEnvVar)
	}
	metalAuthToken := getFromEnvDefault(config.MetalAuthTokenEnvVar, "")

	if clientToken == "" && (clientID == "" || clientSecret == "") && metalAuthToken == "" {
		return nil, fmt.Errorf("To run acceptance tests sweeper, one of '%s' or pair '%s' - '%s' must be set for Equinix Fabric and Network Edge, and '%s' for Equinix Metal",
			config.ClientTokenEnvVar, config.ClientIDEnvVar, config.ClientSecretEnvVar, config.MetalAuthTokenEnvVar)
	}

	return &config.Config{
		AuthToken:      metalAuthToken,
		BaseURL:        endpoint,
		Token:          clientToken,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		RequestTimeout: time.Duration(clientTimeoutInt) * time.Second,
	}, nil
}

func IsSweepableTestResource(namePrefix string) bool {
	return strings.HasPrefix(namePrefix, TestResourcePrefix)
}
