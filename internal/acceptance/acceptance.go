package acceptance

import (
	"os"
	"sync"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/env"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	terraformsdk "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	missingMetalToken = "To run acceptance tests of Equinix Metal Resources, you must set %s"
)

var (
	TestAccProvider          *schema.Provider
	TestAccProviders         map[string]*schema.Provider
	TestExternalProviders    map[string]resource.ExternalProvider
	TestAccFrameworkProvider *provider.FrameworkProvider
	// testAccProviderConfigure ensures Provider is only configured once
	//
	// The PreCheck(t) function is invoked for every test and this prevents
	// extraneous reconfiguration to the same values each time. However, this does
	// not prevent reconfiguration that may happen should the address of
	// Provider be errantly reused in ProviderFactories.
	testAccProviderConfigure sync.Once
)

func TestAccPreCheck(t *testing.T) {
	var err error

	if _, err = env.Get(config.ClientTokenEnvVar); err != nil {
		_, err = env.Get(config.ClientIDEnvVar)
		if err == nil {
			_, err = env.Get(config.ClientSecretEnvVar)
		}
	}

	if err == nil {
		_, err = env.Get(config.MetalAuthTokenEnvVar)
	}

	if err != nil {
		t.Fatalf("To run acceptance tests, one of '%s' or pair '%s' - '%s' must be set for Equinix Fabric and Network Edge, and '%s' for Equinix Metal",
			config.ClientTokenEnvVar, config.ClientIDEnvVar, config.ClientSecretEnvVar, config.MetalAuthTokenEnvVar)
	}
}

func TestAccPreCheckMetal(t *testing.T) {
	if os.Getenv(config.MetalAuthTokenEnvVar) == "" {
		t.Fatalf(missingMetalToken, config.MetalAuthTokenEnvVar)
	}
}

func TestAccPreCheckProviderConfigured(t *testing.T) {
	// Since we are outside the scope of the Terraform configuration we must
	// call Configure() to properly initialize the provider configuration.
	testAccProviderConfigure.Do(func() {
		diags := TestAccProvider.Configure(Context(t), terraformsdk.NewResourceConfigRaw(nil))
		if diags.HasError() {
			t.Fatalf("configuring provider")
		}
	})
}
