package equinix

import (
	"context"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns Equinix terraform *schema.Provider
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema:         provider.Schema(),
		DataSourcesMap: provider.Datasources(),
		ResourcesMap:   provider.Resources(),
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return configureProvider(ctx, d, provider)
	}
	return provider
}

func configureProvider(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	mrws := d.Get("max_retry_wait_seconds").(int)
	rt := d.Get("request_timeout").(int)

	config := config.Config{
		AuthToken:      d.Get("auth_token").(string),
		BaseURL:        d.Get("endpoint").(string),
		ClientID:       d.Get("client_id").(string),
		ClientSecret:   d.Get("client_secret").(string),
		Token:          d.Get("token").(string),
		RequestTimeout: time.Duration(rt) * time.Second,
		PageSize:       d.Get("response_max_page_size").(int),
		MaxRetries:     d.Get("max_retries").(int),
		MaxRetryWait:   time.Duration(mrws) * time.Second,
	}

	config.TerraformVersion = p.TerraformVersion
	if config.TerraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		config.TerraformVersion = "0.11+compatible"
	}

	stopCtx, ok := schema.StopContext(ctx)
	if !ok {
		stopCtx = ctx
	}
	if err := config.Load(stopCtx); err != nil {
		return nil, diag.FromErr(err)
	}
	return &config, nil
}
