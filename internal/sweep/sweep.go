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
	tokenExchangeScope := env.GetWithDefault(config.TokenExchangeScopeEnvVar, "")
	envVar := env.GetWithDefault(config.TokenExchangeSubjectTokenEnvVarEnvVar, config.DefaultTokenExchangeSubjectTokenEnvVar)
	token := env.GetWithDefault(envVar, "")
	if (clientId == "" || clientSecret == "") && (tokenExchangeScope == "" || token == "") {
		return nil, fmt.Errorf("missing Fabric credentials: either %s and %s, or %s and a subject token (from %s, default %s)", config.ClientIDEnvVar, config.ClientSecretEnvVar, config.TokenExchangeScopeEnvVar, config.TokenExchangeSubjectTokenEnvVarEnvVar, config.DefaultTokenExchangeSubjectTokenEnvVar)
	}

	clientTimeout := env.GetWithDefault(config.ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf(cannotConvertTimeoutToInt, config.ClientTimeoutEnvVar)
	}

	return &config.Config{
		BaseURL:                         endpoint,
		ClientID:                        clientId,
		ClientSecret:                    clientSecret,
		RequestTimeout:                  time.Duration(clientTimeoutInt) * time.Second,
		TokenExchangeScope:              tokenExchangeScope,
		TokenExchangeSubjectTokenEnvVar: envVar,
		TokenExchangeSubjectToken:       token,
		StsBaseURL:                      env.GetWithDefault(config.StsEndpointEnvVar, ""),
	}, nil
}
