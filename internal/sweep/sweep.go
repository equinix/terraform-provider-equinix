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
	testResourcePrefix = "tfacc"
	missingMetalToken  = "to run sweepers of Equinix Metal Resources, you must set %s"
)

func IsSweepableTestResource(namePrefix string) bool {
	return strings.HasPrefix(namePrefix, testResourcePrefix)
}

func GetConfigForMetal() (*config.Config, error) {
	endpoint := env.GetWithDefault(config.EndpointEnvVar, config.DefaultBaseURL)
	clientTimeout := env.GetWithDefault(config.ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value of '%s' env variable to int", config.ClientTimeoutEnvVar)
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
