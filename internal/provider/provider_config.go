package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// FrameworkProviderConfig holds the configuration for the Equinix provider.
type FrameworkProviderConfig struct {
	BaseURL                         types.String `tfsdk:"endpoint"`
	ClientID                        types.String `tfsdk:"client_id"`
	ClientSecret                    types.String `tfsdk:"client_secret"`
	Token                           types.String `tfsdk:"token"`
	AuthToken                       types.String `tfsdk:"auth_token"`
	RequestTimeout                  types.Int64  `tfsdk:"request_timeout"`
	PageSize                        types.Int64  `tfsdk:"response_max_page_size"`
	MaxRetries                      types.Int64  `tfsdk:"max_retries"`
	MaxRetryWaitSeconds             types.Int64  `tfsdk:"max_retry_wait_seconds"`
	TokenExchangeScope              types.String `tfsdk:"token_exchange_scope"`
	StsBaseURL                      types.String `tfsdk:"sts_endpoint"`
	StsSourceToken                  types.String `tfsdk:"sts_source_token"`
	StsAuthScope                    types.String `tfsdk:"sts_auth_scope"`
	TokenExchangeSubjectToken       types.String `tfsdk:"token_exchange_subject_token"`
	TokenExchangeSubjectTokenEnvVar types.String `tfsdk:"token_exchange_subject_token_env_var"`
}

func (c *FrameworkProviderConfig) toOldStyleConfig() *config.Config {
	// this immitates func configureProvider in proivder.go
	
	// Determine which scope to use: sts_auth_scope takes precedence over token_exchange_scope
	scope := c.TokenExchangeScope.ValueString()
	if !c.StsAuthScope.IsNull() && c.StsAuthScope.ValueString() != "" {
		scope = c.StsAuthScope.ValueString()
	}
	
	// Determine which source token to use: sts_source_token takes precedence over token_exchange_subject_token
	sourceToken := c.TokenExchangeSubjectToken.ValueString()
	if !c.StsSourceToken.IsNull() && c.StsSourceToken.ValueString() != "" {
		sourceToken = c.StsSourceToken.ValueString()
	}
	
	return &config.Config{
		AuthToken:                       c.AuthToken.ValueString(),
		BaseURL:                         c.BaseURL.ValueString(),
		ClientID:                        c.ClientID.ValueString(),
		ClientSecret:                    c.ClientSecret.ValueString(),
		Token:                           c.Token.ValueString(),
		RequestTimeout:                  time.Duration(c.RequestTimeout.ValueInt64()) * time.Second,
		PageSize:                        int(c.PageSize.ValueInt64()),
		MaxRetries:                      int(c.MaxRetries.ValueInt64()),
		MaxRetryWait:                    time.Duration(c.MaxRetryWaitSeconds.ValueInt64()) * time.Second,
		TokenExchangeScope:              scope,
		StsBaseURL:                      c.StsBaseURL.ValueString(),
		TokenExchangeSubjectToken:       sourceToken,
		TokenExchangeSubjectTokenEnvVar: c.TokenExchangeSubjectTokenEnvVar.ValueString(),
	}
}

// Configure initializes the provider configuration, reading values from the provider block
func (fp *FrameworkProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var fwconfig FrameworkProviderConfig

	// This call reads the configuration from the provider block in the
	// Terraform configuration to the FrameworkProviderConfig struct (config)
	resp.Diagnostics.Append(req.Config.Get(ctx, &fwconfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// We need to supply values from envvar and defaults, because framework
	// provider does not support loading from envvar and defaults :/.
	// (it can validate though)

	// this immitates func Provider() *schema.Provider from provider.go

	fwconfig.BaseURL = determineStrConfValue(
		fwconfig.BaseURL, config.EndpointEnvVar, config.DefaultBaseURL)

	fwconfig.ClientID = determineStrConfValue(
		fwconfig.ClientID, config.ClientIDEnvVar, "")

	fwconfig.ClientSecret = determineStrConfValue(
		fwconfig.ClientSecret, config.ClientSecretEnvVar, "")

	fwconfig.Token = determineStrConfValue(
		fwconfig.Token, config.ClientTokenEnvVar, "")

	fwconfig.AuthToken = determineStrConfValue(
		fwconfig.AuthToken, config.MetalAuthTokenEnvVar, "")

	fwconfig.RequestTimeout = determineIntConfValue(
		fwconfig.RequestTimeout, config.ClientTimeoutEnvVar, int64(config.DefaultTimeout), &resp.Diagnostics)

	fwconfig.MaxRetries = determineIntConfValue(
		fwconfig.MaxRetries, "", 10, &resp.Diagnostics)

	fwconfig.MaxRetryWaitSeconds = determineIntConfValue(
		fwconfig.MaxRetryWaitSeconds, "", 30, &resp.Diagnostics)

	fwconfig.TokenExchangeScope = determineStrConfValue(
		fwconfig.TokenExchangeScope, config.TokenExchangeScopeEnvVar, "")

	fwconfig.StsBaseURL = determineStrConfValue(
		fwconfig.StsBaseURL, config.StsEndpointEnvVar, config.DefaultStsBaseURL)

	fwconfig.StsSourceToken = determineStrConfValue(
		fwconfig.StsSourceToken, config.StsSourceTokenEnvVar, "")

	fwconfig.StsAuthScope = determineStrConfValue(
		fwconfig.StsAuthScope, config.StsAuthScopeEnvVar, "")

	fwconfig.TokenExchangeSubjectToken = determineStrConfValue(
		fwconfig.TokenExchangeSubjectToken, "", "")

	fwconfig.TokenExchangeSubjectTokenEnvVar = determineStrConfValue(
		fwconfig.TokenExchangeSubjectTokenEnvVar, config.TokenExchangeSubjectTokenEnvVarEnvVar, config.DefaultTokenExchangeSubjectTokenEnvVar)

	if resp.Diagnostics.HasError() {
		return
	}

	oldStyleConfig := fwconfig.toOldStyleConfig()
	err := oldStyleConfig.Load(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to load provider configuration",
			err.Error(),
		)
		return
	}
	resp.ResourceData = oldStyleConfig
	resp.DataSourceData = oldStyleConfig

	fp.Meta = oldStyleConfig
}

// GetIntFromEnv retrieves an integer value from the environment variable specified by key.
func GetIntFromEnv(
	key string,
	defaultValue int64,
	diags *diag.Diagnostics,
) int64 {
	if key == "" {
		return defaultValue
	}
	envVarVal := os.Getenv(key)
	if envVarVal == "" {
		return defaultValue
	}

	intVal, err := strconv.ParseInt(envVarVal, 10, 64)
	if err != nil {
		diags.AddWarning(
			fmt.Sprintf(
				"Failed to parse the environment variable %v "+
					"to an integer. Will use default value: %d instead",
				key,
				defaultValue,
			),
			err.Error(),
		)
		return defaultValue
	}

	return intVal
}

func determineIntConfValue(v basetypes.Int64Value, envVar string, defaultValue int64, diags *diag.Diagnostics) basetypes.Int64Value {
	if !v.IsNull() {
		return v
	}
	return types.Int64Value(GetIntFromEnv(envVar, defaultValue, diags))
}

func determineStrConfValue(v basetypes.StringValue, envVar, defaultValue string) basetypes.StringValue {
	if !v.IsNull() {
		return v
	}
	returnVal := os.Getenv(envVar)

	if returnVal == "" {
		returnVal = defaultValue
	}

	return types.StringValue(returnVal)
}
