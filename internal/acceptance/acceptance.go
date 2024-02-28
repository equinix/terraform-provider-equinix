package acceptance

import (
	"os"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/env"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/equinix/terraform-provider-equinix/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
