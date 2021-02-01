package equinix

import (
	"context"
	"fmt"
	"time"

	"github.com/equinix/ecx-go/v2"
	"github.com/equinix/ne-go"
	"github.com/equinix/oauth2-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

//Config is the configuration structure used to instantiate the Equinix
//provider.
type Config struct {
	BaseURL        string
	ClientID       string
	ClientSecret   string
	RequestTimeout time.Duration

	ecx ecx.Client
	ne  ne.Client
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
	authClient := authConfig.New(ctx)
	authClient.Timeout = c.requestTimeout()
	authClient.Transport = logging.NewTransport("Equinix", authClient.Transport)
	c.ecx = ecx.NewClient(ctx, c.BaseURL, authClient)
	c.ne = ne.NewClient(ctx, c.BaseURL, authClient)
	return nil
}

func (c *Config) requestTimeout() time.Duration {
	if c.RequestTimeout == 0 {
		return 5 * time.Second
	}
	return c.RequestTimeout
}
