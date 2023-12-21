package acceptance

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/version"

	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	// duplicated from equinix_sweeoer_test.go
	tstResourcePrefix = "tfacc"
	missingMetalToken = "To run acceptance tests of Equinix Metal Resources, you must set %s"
)

var (
	TestAccProvider          *schema.Provider
	TestAccProviders         map[string]*schema.Provider
	TestAccProviderFactories map[string]func() (*schema.Provider, error)
	TestExternalProviders    map[string]resource.ExternalProvider
	TestAccFrameworkProvider *provider.FrameworkProvider
)

func init() {
	TestAccProvider = provider.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"equinix": TestAccProvider,
	}
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"equinix": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
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

	if _, err = GetFromEnv(config.ClientTokenEnvVar); err != nil {
		_, err = GetFromEnv(config.ClientIDEnvVar)
		if err == nil {
			_, err = GetFromEnv(config.ClientSecretEnvVar)
		}
	}

	if err == nil {
		_, err = GetFromEnv(config.MetalAuthTokenEnvVar)
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

func IsSweepableTestResource(namePrefix string) bool {
	return strings.HasPrefix(namePrefix, tstResourcePrefix)
}

func GetConfigForNonStandardMetalTest() (*config.Config, error) {
	endpoint := GetFromEnvDefault(config.EndpointEnvVar, config.DefaultBaseURL)
	clientTimeout := GetFromEnvDefault(config.ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value of '%s' env variable to int", config.ClientTimeoutEnvVar)
	}
	metalAuthToken := GetFromEnvDefault(config.MetalAuthTokenEnvVar, "")

	if metalAuthToken == "" {
		return nil, fmt.Errorf(missingMetalToken, config.MetalAuthTokenEnvVar)
	}

	return &config.Config{
		AuthToken:      metalAuthToken,
		BaseURL:        endpoint,
		RequestTimeout: time.Duration(clientTimeoutInt) * time.Second,
	}, nil
}
