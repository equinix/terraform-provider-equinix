package equinix

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const tstResourcePrefix = "tf-tst"

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (*Config, error) {
	endpoint, err := getFromEnv(endpointEnvVar)
	if err != nil {
		return nil, err
	}
	clientID, err := getFromEnv(clientIDEnvVar)
	if err != nil {
		return nil, err
	}
	clientSecret, err := getFromEnv(clientSecretEnvVar)
	if err != nil {
		return nil, err
	}
	clientTimeout, err := getFromEnv(clientTimeoutEnvVar)
	if err != nil {
		return nil, err
	}
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value of %sa env variable to int", clientTimeoutEnvVar)
	}
	return &Config{
		BaseURL:        endpoint,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		RequestTimeout: time.Duration(clientTimeoutInt) * time.Second,
	}, nil
}

func isSweepableTestResource(name string) bool {
	return strings.HasPrefix(name, tstResourcePrefix)
}
