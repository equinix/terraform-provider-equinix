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
	missingSecrets            = "missing Equinix credentials: set one of '%s', the pair '%s' and '%s', or '%s' and a subject token in '%s'"
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

// GetConfig returns a provider configuration for sweepers that call the
// Equinix API. Credentials are not product-specific: any of an API token, a
// client ID and secret pair, or an STS token exchange scope and subject token
// will do.
func GetConfig() (*config.Config, error) {
	endpoint := env.GetWithDefault(config.EndpointEnvVar, config.DefaultBaseURL)
	clientToken := env.GetWithDefault(config.ClientTokenEnvVar, "")
	clientID := env.GetWithDefault(config.ClientIDEnvVar, "")
	clientSecret := env.GetWithDefault(config.ClientSecretEnvVar, "")
	tokenExchangeScope := env.GetWithDefault(config.TokenExchangeScopeEnvVar, "")
	subjectTokenEnvVar := env.GetWithDefault(config.TokenExchangeSubjectTokenEnvVarEnvVar, config.DefaultTokenExchangeSubjectTokenEnvVar)
	subjectToken := env.GetWithDefault(subjectTokenEnvVar, "")

	if clientToken == "" && (clientID == "" || clientSecret == "") && (tokenExchangeScope == "" || subjectToken == "") {
		return nil, fmt.Errorf(missingSecrets,
			config.ClientTokenEnvVar, config.ClientIDEnvVar, config.ClientSecretEnvVar,
			config.TokenExchangeScopeEnvVar, subjectTokenEnvVar)
	}

	clientTimeout := env.GetWithDefault(config.ClientTimeoutEnvVar, strconv.Itoa(config.DefaultTimeout))
	clientTimeoutInt, err := strconv.Atoi(clientTimeout)
	if err != nil {
		return nil, fmt.Errorf(cannotConvertTimeoutToInt, config.ClientTimeoutEnvVar)
	}

	return &config.Config{
		BaseURL:                         endpoint,
		Token:                           clientToken,
		ClientID:                        clientID,
		ClientSecret:                    clientSecret,
		RequestTimeout:                  time.Duration(clientTimeoutInt) * time.Second,
		TokenExchangeScope:              tokenExchangeScope,
		TokenExchangeSubjectTokenEnvVar: subjectTokenEnvVar,
		TokenExchangeSubjectToken:       subjectToken,
		StsBaseURL:                      env.GetWithDefault(config.StsEndpointEnvVar, config.DefaultStsBaseURL),
	}, nil
}
