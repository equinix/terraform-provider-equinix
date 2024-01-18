package acceptance

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/config"

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
)

func init() {
	TestAccProvider = equinix.Provider()
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
}

func TestAccPreCheckMetal(t *testing.T) {
	if os.Getenv(config.MetalAuthTokenEnvVar) == "" {
		t.Fatalf(missingMetalToken, config.MetalAuthTokenEnvVar)
	}
}

func IsSweepableTestResource(namePrefix string) bool {
	return strings.HasPrefix(namePrefix, tstResourcePrefix)
}

func getFromEnvDefault(varName string, defaultValue string) string {
	if v := os.Getenv(varName); v != "" {
		return v
	}
	return defaultValue
}

func GetConfigForNonStandardMetalTest() (*config.Config, error) {
	endpoint := getFromEnvDefault(config.EndpointEnvVar, config.DefaultBaseURL)
	clientTimeout := getFromEnvDefault(config.ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value of '%s' env variable to int", config.ClientTimeoutEnvVar)
	}
	metalAuthToken := getFromEnvDefault(config.MetalAuthTokenEnvVar, "")

	if metalAuthToken == "" {
		return nil, fmt.Errorf(missingMetalToken, config.MetalAuthTokenEnvVar)
	}

	return &config.Config{
		AuthToken:      metalAuthToken,
		BaseURL:        endpoint,
		RequestTimeout: time.Duration(clientTimeoutInt) * time.Second,
	}, nil
}
