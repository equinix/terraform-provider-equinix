package equinix

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/equinix/terraform-provider-equinix/version"

	"github.com/equinix/ecx-go/v2"
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

type mockECXClient struct {
	GetUserPortsFn func() ([]ecx.Port, error)

	GetL2OutgoingConnectionsFn     func(statuses []string) ([]ecx.L2Connection, error)
	GetL2ConnectionFn              func(uuid string) (*ecx.L2Connection, error)
	CreateL2ConnectionFn           func(conn ecx.L2Connection) (*string, error)
	CreateL2RedundantConnectionFn  func(priConn, secConn ecx.L2Connection) (*string, *string, error)
	NewL2ConnectionUpdateRequestFn func(uuid string) ecx.L2ConnectionUpdateRequest
	DeleteL2ConnectionFn           func(uuid string) error
	ConfirmL2ConnectionFn          func(uuid string, confirmConn ecx.L2ConnectionToConfirm) (*ecx.L2ConnectionConfirmation, error)

	GetL2SellerProfilesFn    func() ([]ecx.L2ServiceProfile, error)
	GetL2ServiceProfileFn    func(uuid string) (*ecx.L2ServiceProfile, error)
	CreateL2ServiceProfileFn func(sp ecx.L2ServiceProfile) (*string, error)
	UpdateL2ServiceProfileFn func(sp ecx.L2ServiceProfile) error
	DeleteL2ServiceProfileFn func(uuid string) error
}

func (m *mockECXClient) GetUserPorts() ([]ecx.Port, error) {
	return m.GetUserPortsFn()
}

func (m *mockECXClient) GetL2OutgoingConnections(statuses []string) ([]ecx.L2Connection, error) {
	return m.GetL2OutgoingConnectionsFn(statuses)
}

func (m *mockECXClient) GetL2Connection(uuid string) (*ecx.L2Connection, error) {
	return m.GetL2ConnectionFn(uuid)
}

func (m *mockECXClient) CreateL2Connection(conn ecx.L2Connection) (*string, error) {
	return m.CreateL2ConnectionFn(conn)
}

func (m *mockECXClient) CreateL2RedundantConnection(priConn, secConn ecx.L2Connection) (*string, *string, error) {
	return m.CreateL2RedundantConnectionFn(priConn, secConn)
}

func (m *mockECXClient) NewL2ConnectionUpdateRequest(uuid string) ecx.L2ConnectionUpdateRequest {
	return m.NewL2ConnectionUpdateRequestFn(uuid)
}

func (m *mockECXClient) DeleteL2Connection(uuid string) error {
	return m.DeleteL2ConnectionFn(uuid)
}

func (m *mockECXClient) ConfirmL2Connection(uuid string, confirmConn ecx.L2ConnectionToConfirm) (*ecx.L2ConnectionConfirmation, error) {
	return m.ConfirmL2ConnectionFn(uuid, confirmConn)
}

func (m *mockECXClient) GetL2SellerProfiles() ([]ecx.L2ServiceProfile, error) {
	return m.GetL2SellerProfilesFn()
}

func (m *mockECXClient) GetL2ServiceProfile(uuid string) (*ecx.L2ServiceProfile, error) {
	return m.GetL2ServiceProfileFn(uuid)
}

func (m *mockECXClient) CreateL2ServiceProfile(sp ecx.L2ServiceProfile) (*string, error) {
	return m.CreateL2ServiceProfileFn(sp)
}

func (m *mockECXClient) UpdateL2ServiceProfile(sp ecx.L2ServiceProfile) error {
	return m.UpdateL2ServiceProfileFn(sp)
}

func (m *mockECXClient) DeleteL2ServiceProfile(uuid string) error {
	return m.DeleteL2ServiceProfileFn(uuid)
}

var _ ecx.Client = (*mockECXClient)(nil)

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

func TestProvider_atLeastOneStringFound(t *testing.T) {
	// given
	needles := []string{"key4", "key2"}
	hay := []string{"key1", "key2"}
	// when
	result := atLeastOneStringFound(needles, hay)
	// then
	assert.True(t, result, "Given strings were found")
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
	if !comparisons.IsEmpty(value) {
		return d.Set(key, value)
	}
	return nil
}
