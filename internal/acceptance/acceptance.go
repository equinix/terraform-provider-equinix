// Package acceptance provides Utilities and test framework setup for running
// acceptance tests for the Equinix Terraform provider. It handles provider
// configuration, authentication verification, and prerequisite checks for
// testing against Equinix Fabric, Network Edge, and Metal services.
package acceptance

import (
	"os"
	"sync"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/env"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/equinix/terraform-provider-equinix/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	terraformsdk "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	missingMetalToken = "To run acceptance tests of Equinix Metal Resources, you must set %s"
)

var (
	// TestAccProvider is the Equinix provider instance used for acceptance testing
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

func init() {
	TestAccProvider = equinix.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"equinix": TestAccProvider,
	}
	TestExternalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source: "hashicorp/random",
		},
	}
	TestAccFrameworkProvider = provider.CreateFrameworkProvider(version.ProviderVersion).(*provider.FrameworkProvider)
}

// TestAccPreCheck verifies that the required environment variables are set
// for running acceptance tests. It checks for authentication credentials for
// Equinix Fabric, Network Edge, and Metal services.
func TestAccPreCheck(t *testing.T) {
	var err error

	if _, err = env.Get(config.ClientTokenEnvVar); err != nil {
		_, err = env.Get(config.ClientIDEnvVar)
		if err == nil {
			_, err = env.Get(config.ClientSecretEnvVar)
		}

		// If neither token nor client ID/secret are configured, check for STS source token
		if err != nil {
			_, authScopeErr := env.Get(config.AuthScopeEnvVar)
			_, stsTokenErr := env.Get(config.StsSourceTokenEnvVar)

			if authScopeErr == nil && stsTokenErr == nil {
				err = nil
			}
		}
	}

	if err == nil {
		_, err = env.Get(config.MetalAuthTokenEnvVar)
	}

	if err != nil {
		t.Fatalf("To run acceptance tests, one of '%s', pair '%s' - '%s', or pair '%s' - '%s' must be set for Equinix Fabric and Network Edge, and '%s' for Equinix Metal",
			config.ClientTokenEnvVar, config.ClientIDEnvVar, config.ClientSecretEnvVar,
			config.AuthScopeEnvVar, config.StsSourceTokenEnvVar, config.MetalAuthTokenEnvVar)
	}
}

// TestAccPreCheckMetal specifically verifies that the Equinix Metal authentication token
// environment variable is set for running Metal-specific acceptance tests.
func TestAccPreCheckMetal(t *testing.T) {
	if os.Getenv(config.MetalAuthTokenEnvVar) == "" {
		t.Fatalf(missingMetalToken, config.MetalAuthTokenEnvVar)
	}
}

// TestAccPreCheckProviderConfigured ensures the provider is properly configured
// before running tests. It uses sync.Once to guarantee the provider is
// configured exactly once across all test executions.
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
