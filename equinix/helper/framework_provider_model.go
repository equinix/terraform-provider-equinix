package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
)

type FrameworkProviderModel struct {
	Endpoint            types.String `tfsdk:"endpoint,omitempty"`
	ClientID            types.String `tfsdk:"client_id,omitempty"`
	ClientSecret        types.String `tfsdk:"client_secret,omitempty"`
	Token               types.String `tfsdk:"token,omitempty"`
	AuthToken           types.String `tfsdk:"auth_token,omitempty"`
	RequestTimeout      types.Int64  `tfsdk:"request_timeout,omitempty"`
	ResponseMaxPageSize types.Int64  `tfsdk:"response_max_page_size,omitempty"`
	MaxRetries          types.Int64  `tfsdk:"max_retries,omitempty"`
	MaxRetryWaitSeconds types.Int64  `tfsdk:"max_retry_wait_seconds,omitempty"`
}

type FrameworkProviderMeta struct {
	Client *packngo.Client
	Config *FrameworkProviderModel
}
