package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestFrameworkProviderConfig_toOldStyleConfig(t *testing.T) {
	tests := []struct {
		name                          string
		config                        FrameworkProviderConfig
		expectedTokenExchangeScope    string
		expectedTokenExchangeSubToken string
	}{
		{
			name: "sts_auth_scope takes precedence over token_exchange_scope",
			config: FrameworkProviderConfig{
				StsAuthScope:       types.StringValue("sts-scope"),
				TokenExchangeScope: types.StringValue("token-exchange-scope"),
			},
			expectedTokenExchangeScope: "sts-scope",
		},
		{
			name: "token_exchange_scope used when sts_auth_scope is not set",
			config: FrameworkProviderConfig{
				StsAuthScope:       types.StringNull(),
				TokenExchangeScope: types.StringValue("token-exchange-scope"),
			},
			expectedTokenExchangeScope: "token-exchange-scope",
		},
		{
			name: "sts_source_token takes precedence over token_exchange_subject_token",
			config: FrameworkProviderConfig{
				StsSourceToken:            types.StringValue("sts-token"),
				TokenExchangeSubjectToken: types.StringValue("token-exchange-token"),
			},
			expectedTokenExchangeSubToken: "sts-token",
		},
		{
			name: "token_exchange_subject_token used when sts_source_token is not set",
			config: FrameworkProviderConfig{
				StsSourceToken:            types.StringNull(),
				TokenExchangeSubjectToken: types.StringValue("token-exchange-token"),
			},
			expectedTokenExchangeSubToken: "token-exchange-token",
		},
		{
			name: "both sts fields take precedence",
			config: FrameworkProviderConfig{
				StsAuthScope:              types.StringValue("sts-scope"),
				TokenExchangeScope:        types.StringValue("token-exchange-scope"),
				StsSourceToken:            types.StringValue("sts-token"),
				TokenExchangeSubjectToken: types.StringValue("token-exchange-token"),
			},
			expectedTokenExchangeScope:    "sts-scope",
			expectedTokenExchangeSubToken: "sts-token",
		},
		{
			name: "empty string sts_auth_scope does not override token_exchange_scope",
			config: FrameworkProviderConfig{
				StsAuthScope:       types.StringValue(""),
				TokenExchangeScope: types.StringValue("token-exchange-scope"),
			},
			expectedTokenExchangeScope: "token-exchange-scope",
		},
		{
			name: "empty string sts_source_token does not override token_exchange_subject_token",
			config: FrameworkProviderConfig{
				StsSourceToken:            types.StringValue(""),
				TokenExchangeSubjectToken: types.StringValue("token-exchange-token"),
			},
			expectedTokenExchangeSubToken: "token-exchange-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldConfig := tt.config.toOldStyleConfig()

			if tt.expectedTokenExchangeScope != "" {
				assert.Equal(t, tt.expectedTokenExchangeScope, oldConfig.TokenExchangeScope)
			}
			if tt.expectedTokenExchangeSubToken != "" {
				assert.Equal(t, tt.expectedTokenExchangeSubToken, oldConfig.TokenExchangeSubjectToken)
			}
		})
	}
}
