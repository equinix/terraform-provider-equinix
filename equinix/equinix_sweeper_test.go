package equinix

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const tstResourcePrefix = "tfacc"

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (*Config, error) {
	endpoint := getFromEnvDefault(endpointEnvVar, DefaultBaseURL)
	clientToken := getFromEnvDefault(clientTokenEnvVar, "")
	clientID := getFromEnvDefault(clientIDEnvVar, "")
	clientSecret := getFromEnvDefault(clientSecretEnvVar, "")
	clientTimeout := getFromEnvDefault(clientTimeoutEnvVar, strconv.Itoa(DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value of '%s' env variable to int", clientTimeoutEnvVar)
	}
	metalAuthToken := getFromEnvDefault(metalAuthTokenEnvVar, "")

	if clientToken == "" && (clientID == "" || clientSecret == "") && metalAuthToken == "" {
		return nil, fmt.Errorf("To run acceptance tests sweeper, one of '%s' or pair '%s' - '%s' must be set for Equinix Fabric and Network Edge, and '%s' for Equinix Metal",
			clientTokenEnvVar, clientIDEnvVar, clientSecretEnvVar, metalAuthTokenEnvVar)
	}

	return &Config{
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
