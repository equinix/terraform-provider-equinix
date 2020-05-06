package equinix

import (
	"context"
	ecx "ecx-go-client/v3"
	"ecx-go-client/v3/oauth2"
	"fmt"
	"time"
)

//Config is the configuration structure used to instantiate the Equinix
//provider.
type Config struct {
	BaseURL        string
	ClientID       string
	ClientSecret   string
	RequestTimeout time.Duration

	ecx ecx.Client
}

//Load function validates configuration structure fields and configures
//all required API clients.
func (c *Config) Load(ctx context.Context) error {
	if c.BaseURL == "" {
		return fmt.Errorf("baseURL cannot be empty")
	}
	if c.ClientID == "" {
		return fmt.Errorf("clientId cannot be empty")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("clientSecret cannot be empty")
	}
	authConfig := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		BaseURL:      c.BaseURL}
	authClient := authConfig.Client(ctx)
	authClient.Timeout = c.requestTimeout()
	c.ecx = ecx.NewClient(c.BaseURL, authClient)
	return nil
}

func (c *Config) requestTimeout() time.Duration {
	if c.RequestTimeout == 0 {
		return 5 * time.Second
	}
	return c.RequestTimeout
}
