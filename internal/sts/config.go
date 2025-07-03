package sts

import (
	"github.com/equinix/equinix-sdk-go/services/stsv1alpha"
	"sync"
)

// Config describes oauth2 client credentials flow
type Config struct {
	// ClientID is the application's ID.
	StsAuthScope string
	// ClientSecret is the application's secret.
	StsSourceToken string
	// StsBaseURL is the base endpoint of a server that  token endpoint
	StsBaseURL string
}

// StsTokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided context and the
// client ID and client secret.
func (c *Config) StsTokenSource() *StsContextAwareTokenSource {
	config := stsv1alpha.NewConfiguration()
	config.Servers = stsv1alpha.ServerConfigurations{
		stsv1alpha.ServerConfiguration{
			URL: c.StsBaseURL,
		},
	}
	restClient := stsv1alpha.NewAPIClient(config)
	source := StsContextAwareTokenSource{
		c,
		restClient,
		sync.Mutex{},
		nil,
	}
	return &source
}
