package equinix

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const tstResourcePrefix = "tfacc"

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (*config.Config, error) {
	endpoint := getFromEnvDefault(EndpointEnvVar, config.DefaultBaseURL)
	clientToken := getFromEnvDefault(ClientTokenEnvVar, "")
	clientID := getFromEnvDefault(ClientIDEnvVar, "")
	clientSecret := getFromEnvDefault(ClientSecretEnvVar, "")
	clientTimeout := getFromEnvDefault(ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value of '%s' env variable to int", ClientTimeoutEnvVar)
	}
	metalAuthToken := getFromEnvDefault(MetalAuthTokenEnvVar, "")

	if clientToken == "" && (clientID == "" || clientSecret == "") && metalAuthToken == "" {
		return nil, fmt.Errorf("To run acceptance tests sweeper, one of '%s' or pair '%s' - '%s' must be set for Equinix Fabric and Network Edge, and '%s' for Equinix Metal",
			ClientTokenEnvVar, ClientIDEnvVar, ClientSecretEnvVar, MetalAuthTokenEnvVar)
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
