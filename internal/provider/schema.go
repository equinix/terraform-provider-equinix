package provider

import (
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoint": {
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc(config.EndpointEnvVar, config.DefaultBaseURL),
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Description:  fmt.Sprintf("The Equinix API base URL to point out desired environment. Defaults to %s", config.DefaultBaseURL),
		},
		"client_id": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc(config.ClientIDEnvVar, ""),
			Description: "API Consumer Key available under My Apps section in developer portal",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc(config.ClientSecretEnvVar, ""),
			Description: "API Consumer secret available under My Apps section in developer portal",
		},
		"token": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc(config.ClientTokenEnvVar, ""),
			Description: "API token from the developer sandbox",
		},
		"auth_token": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc(config.MetalAuthTokenEnvVar, ""),
			Description: "The Equinix Metal API auth key for API operations",
		},
		"request_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc(config.ClientTimeoutEnvVar, config.DefaultTimeout),
			ValidateFunc: validation.IntAtLeast(1),
			Description:  fmt.Sprintf("The duration of time, in seconds, that the Equinix Platform API Client should wait before canceling an API request.  Defaults to %d", config.DefaultTimeout),
		},
		"response_max_page_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(100),
			Description:  "The maximum number of records in a single response for REST queries that produce paginated responses",
		},
		"max_retries": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  10,
		},
		"max_retry_wait_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  30,
		},
	}
}
