package equinix

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/equinix/terraform-provider-equinix/version"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

var (
	testAccProviders         map[string]*schema.Provider
	testAccProvider          *schema.Provider
	testExternalProviders    map[string]resource.ExternalProvider
	testAccFrameworkProvider *provider.FrameworkProvider

	testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		"equinix": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				testAccProviders["equinix"].GRPCProvider,
				providerserver.NewProtocol5(
					testAccFrameworkProvider,
				),
			}
			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}
			return muxServer.ProviderServer(), nil
		},
	}
)

type testAccConfig struct {
	ctx    map[string]interface{}
	config string
}

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"equinix": testAccProvider,
	}
	testExternalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source: "hashicorp/random",
		},
	}
	// during framework migration, it is required to duplicate this (TestAccFrameworkProvider declared in internal package)
	// for e2e tests that need already migrated resources. Importing from internal produces and import cycle error
	testAccFrameworkProvider = provider.CreateFrameworkProvider(version.ProviderVersion).(*provider.FrameworkProvider)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_setSchemaValueIfNotEmpty(t *testing.T) {
	// given
	key := "test"
	s := map[string]*schema.Schema{
		key: {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	var b *int = nil
	d := schema.TestResourceDataRaw(t, s, make(map[string]interface{}))
	// when
	setSchemaValueIfNotEmpty(key, b, d)
	// then
	_, ok := d.GetOk(key)
	assert.False(t, ok, "Key was not set")
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Test helper functions
//_______________________________________________________________________

func testAccPreCheck(t *testing.T) {
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

func newTestAccConfig(ctx map[string]interface{}) *testAccConfig {
	return &testAccConfig{
		ctx:    ctx,
		config: "",
	}
}

func (t *testAccConfig) build() string {
	return t.config
}

// nprintf returns a string with all the placeholders replaced by the values from the params map
//
// Deprecated: nprintf is shared between NE resource tests and has been
// centralized ahead of those NE resources moving to separate packages.
// Use github.com/equinix/terraform-provider-equinix/internal/nprintf.NPrintf instead
func nprintf(format string, params map[string]interface{}) string {
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

func copyMap(source map[string]interface{}) map[string]interface{} {
	target := make(map[string]interface{})
	for k, v := range source {
		target[k] = v
	}
	return target
}

func setSchemaValueIfNotEmpty(key string, value interface{}, d *schema.ResourceData) error {
	if !isEmpty(value) {
		return d.Set(key, value)
	}
	return nil
}
