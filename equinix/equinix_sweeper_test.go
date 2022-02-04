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
	clientID := ""
	clientSecret := ""
	clientToken, err := getFromEnv(clientTokenEnvVar)
	if err != nil {
		clientID, err = getFromEnv(clientIDEnvVar)
		if err != nil {
			return nil, fmt.Errorf("one of '%s' or pair '%s' - '%s' must be set for acceptance tests", clientTokenEnvVar, clientIDEnvVar, clientSecretEnvVar)
		}
		clientSecret, err = getFromEnv(clientSecretEnvVar)
		if err != nil {
			return nil, fmt.Errorf("one of '%s' or pair '%s' - '%s' must be set for acceptance tests", clientTokenEnvVar, clientIDEnvVar, clientSecretEnvVar)
		}
	}
	clientTimeout := getFromEnvDefault(clientTimeoutEnvVar, strconv.Itoa(DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value of '%s' env variable to int", clientTimeoutEnvVar)
	}
	metalAuthToken, err := getFromEnv(metalAuthTokenEnvVar)
	if err != nil {
		return nil, err
	}
	return &Config{
		AuthToken:      metalAuthToken,
		BaseURL:        endpoint,
		Token:			clientToken,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		RequestTimeout: time.Duration(clientTimeoutInt) * time.Second,
	}, nil
}

func isSweepableTestResource(namePrefix string) bool {
	return strings.HasPrefix(namePrefix, tstResourcePrefix)
}
