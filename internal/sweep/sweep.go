package sweep

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/env"
)

const (
	// duplicated from equinix_sweeoer_test.go
	testResourcePrefix        = "tfacc"
	cannotConvertTimeoutToInt = "cannot convert value of '%s' env variable to int"
	missingFabricSecrets      = "missing fabric clientId - %s, and clientSecret - %s"
	missingMetalToken         = "to run sweepers of Equinix Metal Resources, you must set %s"
)

var (
	FabricTestResourceSuffixes = []string{"_PFCR", "_PNFV", "_PPDS"}
)

func IsSweepableTestResource(namePrefix string) bool {
	return strings.HasPrefix(namePrefix, testResourcePrefix)
}

func IsSweepableFabricTestResource(resourceName string) bool {
	for _, suffix := range FabricTestResourceSuffixes {
		if strings.HasSuffix(resourceName, suffix) {
			return true
		}
	}
	return false
}

func GetConfigForFabric() (*config.Config, error) {
	endpoint := env.GetWithDefault(config.EndpointEnvVar, config.DefaultBaseURL)
	clientId := env.GetWithDefault(config.ClientIDEnvVar, "")
	clientSecret := env.GetWithDefault(config.ClientSecretEnvVar, "")
	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf(missingFabricSecrets, config.ClientIDEnvVar, config.ClientSecretEnvVar)
	}

	clientTimeout := env.GetWithDefault(config.ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf(cannotConvertTimeoutToInt, config.ClientTimeoutEnvVar)
	}

	return &config.Config{
		BaseURL:        endpoint,
		ClientID:       clientId,
		ClientSecret:   clientSecret,
		RequestTimeout: time.Duration(clientTimeoutInt) * time.Second,
	}, nil
}

func GetConfigForMetal() (*config.Config, error) {
	endpoint := env.GetWithDefault(config.EndpointEnvVar, config.DefaultBaseURL)
	clientTimeout := env.GetWithDefault(config.ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf(cannotConvertTimeoutToInt, config.ClientTimeoutEnvVar)
	}
	metalAuthToken := env.GetWithDefault(config.MetalAuthTokenEnvVar, "")

	if metalAuthToken == "" {
		return nil, fmt.Errorf(missingMetalToken, config.MetalAuthTokenEnvVar)
	}

	return &config.Config{
		AuthToken:      metalAuthToken,
		BaseURL:        endpoint,
		RequestTimeout: time.Duration(clientTimeoutInt) * time.Second,
	}, nil
}
