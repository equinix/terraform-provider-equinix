package equinix

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const tstResourcePrefix = "tfacc"

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

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

	if clientToken == "" && (clientID == "" || clientSecret == "") {
		return nil, fmt.Errorf("To run acceptance tests sweeper, one of '%s' or pair '%s' - '%s' must be set for Equinix Fabric and Network Edge",
			config.ClientTokenEnvVar, config.ClientIDEnvVar, config.ClientSecretEnvVar)
	}

	return &config.Config{
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
