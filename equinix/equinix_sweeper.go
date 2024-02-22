package equinix

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
)

const tstResourcePrefix = "tfacc"

func sharedConfigForRegion(region string) (*config.Config, error) {
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

func isSweepableTestResource(namePrefix string) bool {
	return strings.HasPrefix(namePrefix, tstResourcePrefix)
}

func getFromEnvDefault(varName string, defaultValue string) string {
	if v := os.Getenv(varName); v != "" {
		return v
	}
	return defaultValue
}

// Deprecated: this function is a shim to allow us to
// continue to run all Metal sweepers while we migrate resources
// and datasources out of the equinix package and into
// isolated packages.  This should be removed after all
// resources have been migrated out of package equinix.
func AddMetalTestSweepers() {
	addMetalDeviceSweeper()
	addMetalOrganizationSweeper()
	addMetalUserAPIKeySweeper()
	addMetalVirtualCircuitSweeper()
	addMetalVlanSweeper()
}
