// Package config uses environment variables to configure API clients for the provider
package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/extensions/equinixoauth2"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/internal/sts"
	"github.com/equinix/terraform-provider-equinix/version"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2"
)

// These constants track environment variable names
// that are relevant to the provider
const (
	EndpointEnvVar                         = "EQUINIX_API_ENDPOINT"
	ClientIDEnvVar                         = "EQUINIX_API_CLIENTID"
	ClientSecretEnvVar                     = "EQUINIX_API_CLIENTSECRET"
	ClientTokenEnvVar                      = "EQUINIX_API_TOKEN"
	ClientTimeoutEnvVar                    = "EQUINIX_API_TIMEOUT"
	TokenExchangeScopeEnvVar               = "EQUINIX_TOKEN_EXCHANGE_SCOPE"
	TokenExchangeSubjectTokenEnvVarEnvVar  = "EQUINIX_TOKEN_EXCHANGE_SUBJECT_TOKEN_ENV_VAR"
	StsEndpointEnvVar                      = "EQUINIX_STS_ENDPOINT"
	DefaultTokenExchangeSubjectTokenEnvVar = "EQUINIX_TOKEN_EXCHANGE_SUBJECT_TOKEN"
)

// ProviderMeta allows passing additional metadata
// specific to our provider in provider config blocks
// See https://developer.hashicorp.com/terraform/internals/provider-meta
type ProviderMeta struct {
	ModuleName string `cty:"module_name"`
}

const (
	uaEnvVar = "TF_APPEND_USER_AGENT"
)

var (
	// DefaultBaseURL is the default URL to use for API requests
	DefaultBaseURL = "https://api.equinix.com"
	// DefaultStsBaseURL is the default Security Token Service (STS) endpoint
	DefaultStsBaseURL = "https://sts.eqix.equinix.com"
	// DefaultTimeout is the default request timeout to use for API requests
	DefaultTimeout = 30
)

// Config is the configuration structure used to instantiate the Equinix
// provider.
type Config struct {
	BaseURL                         string
	AuthToken                       string
	ClientID                        string
	ClientSecret                    string
	MaxRetries                      int
	MaxRetryWait                    time.Duration
	RequestTimeout                  time.Duration
	PageSize                        int
	Token                           string
	TokenExchangeScope              string
	StsBaseURL                      string
	TokenExchangeSubjectToken       string
	TokenExchangeSubjectTokenEnvVar string

	authClient *http.Client

	Ne ne.Client

	neUserAgent string

	TerraformVersion string
}

// Load function validates configuration structure fields and configures
// all required API clients.
func (c *Config) Load(ctx context.Context) error {
	if c.BaseURL == "" {
		return fmt.Errorf("'baseURL' cannot be empty")
	}

	// Validate BaseURL
	_, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	// If StsBaseURL is set, validate it too
	if c.StsBaseURL != "" {
		_, err := url.Parse(c.StsBaseURL)
		if err != nil {
			return fmt.Errorf("invalid STS base URL: %w", err)
		}
	}
	c.authClient = c.newAuthClient()

	neClient := ne.NewClient(ctx, c.BaseURL, c.authClient)

	if c.PageSize > 0 {
		neClient.SetPageSize(c.PageSize)
	}
	c.neUserAgent = c.tfSdkUserAgent("equinix/ne-go")
	neClient.SetHeaders(map[string]string{
		"User-agent": c.neUserAgent,
	})

	c.Ne = neClient
	return nil
}

func (c *Config) newAuthClient() *http.Client {
	var authTransport http.RoundTripper
	if c.TokenExchangeScope != "" {
		sourceToken := c.resolveSourceToken()
		if sourceToken != "" {
			authConfig := sts.Config{
				StsAuthScope:   c.TokenExchangeScope,
				StsSourceToken: sourceToken,
				StsBaseURL:     c.StsBaseURL,
			}
			authTransport = authConfig.New()
		}
	}

	// If no STS auth, fall back to existing logic
	if authTransport == nil {
		if c.Token != "" {
			tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.Token})
			oauthTransport := &oauth2.Transport{
				Source: tokenSource,
			}
			authTransport = oauthTransport
		} else {
			authConfig := equinixoauth2.Config{
				ClientID:     c.ClientID,
				ClientSecret: c.ClientSecret,
				BaseURL:      c.BaseURL,
			}
			authTransport = authConfig.New()
		}
	}

	authClient := http.Client{
		Timeout: c.requestTimeout(),
		//nolint:staticcheck
		Transport: logging.NewTransport("Equinix", authTransport),
	}
	return &authClient
}

func (c *Config) resolveSourceToken() string {
	// First priority: explicitly configured token
	if c.TokenExchangeSubjectToken != "" {
		return c.TokenExchangeSubjectToken
	}

	// Second priority: token from environment variable
	if c.TokenExchangeSubjectTokenEnvVar != "" {
		return os.Getenv(c.TokenExchangeSubjectTokenEnvVar)
	}

	return ""
}

// NewFabricClientForSDK returns a terraform sdkv2 plugin compatible
// equinix-sdk-go/fabricv4 client to be used to access Fabric's V4 APIs
func (c *Config) NewFabricClientForSDK(_ context.Context, d *schema.ResourceData) *fabricv4.APIClient {
	client := c.newFabricClient()

	baseUserAgent := c.tfSdkUserAgent(client.GetConfig().UserAgent)
	client.GetConfig().UserAgent = generateModuleUserAgentString(d, baseUserAgent)

	return client
}

// NewFabricClientForTesting is a shim for Fabric tests.
// Deprecated: when the acceptance package starts to contain API clients for testing/cleanup this will move with them
func (c *Config) NewFabricClientForTesting(_ context.Context) *fabricv4.APIClient {
	client := c.newFabricClient()

	client.GetConfig().UserAgent = fmt.Sprintf("tf-acceptance-tests %v", client.GetConfig().UserAgent)

	return client
}

// NewFabricClientForFramework returns a terraform framework compatible
// equinix-sdk-go/fabricv4 client to be used to access Fabric's V4 APIs
func (c *Config) NewFabricClientForFramework(ctx context.Context, meta tfsdk.Config) *fabricv4.APIClient {
	client := c.newFabricClient()

	baseUserAgent := c.tfFrameworkUserAgent(client.GetConfig().UserAgent)
	client.GetConfig().UserAgent = generateFwModuleUserAgentString(ctx, meta, baseUserAgent)

	return client
}

// newFabricClient returns the base fabricv4 client that is then used for either the sdkv2 or framework
// implementations of the Terraform Provider with exported Methods
func (c *Config) newFabricClient() *fabricv4.APIClient {
	authClient := c.newAuthClient()

	// Configure HTTP client with retries and logging
	httpClient := c.configureHTTPClient(authClient)

	return c.createFabricClient(httpClient)
}

func (c *Config) configureHTTPClient(client *http.Client) *http.Client {
	//nolint:staticcheck // We should move to subsystem loggers, but that is a much bigger change
	transport := logging.NewTransport("Equinix Fabric (fabricv4)", client.Transport)

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = transport
	retryClient.HTTPClient.Timeout = c.requestTimeout()
	retryClient.RetryMax = c.MaxRetries
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryWaitMax = c.MaxRetryWait

	return retryClient.StandardClient()
}

func (c *Config) createFabricClient(httpClient *http.Client) *fabricv4.APIClient {
	baseURL, _ := url.Parse(c.BaseURL)

	configuration := fabricv4.NewConfiguration()
	configuration.Servers = fabricv4.ServerConfigurations{
		fabricv4.ServerConfiguration{
			URL: baseURL.String(),
		},
	}
	configuration.HTTPClient = httpClient
	configuration.AddDefaultHeader("X-SOURCE", "API")
	configuration.AddDefaultHeader("X-CORRELATION-ID", correlationId(25))

	return fabricv4.NewAPIClient(configuration)
}

func (c *Config) requestTimeout() time.Duration {
	if c.RequestTimeout == 0 {
		return 5 * time.Second
	}
	return c.RequestTimeout
}

func appendUserAgentFromEnv(ua string) string {
	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

// AddModuleToNEUserAgent injects the ModuleName into SDK resource metadata for analytics
func (c *Config) AddModuleToNEUserAgent(client *ne.Client, d *schema.ResourceData) {
	cli := *client
	rc := cli.(*ne.RestClient)
	rc.SetHeader("User-agent", generateModuleUserAgentString(d, c.neUserAgent))
	*client = rc
}

func generateFwModuleUserAgentString(ctx context.Context, meta tfsdk.Config, baseUserAgent string) string {
	var m ProviderMeta
	diags := meta.Get(ctx, &m)
	if diags.HasError() {
		log.Printf("[WARN] error retrieving provider_meta")
		return baseUserAgent
	}
	if m.ModuleName != "" {
		return strings.Join([]string{m.ModuleName, baseUserAgent}, " ")
	}
	return baseUserAgent
}

func generateModuleUserAgentString(d *schema.ResourceData, baseUserAgent string) string {
	var m ProviderMeta
	err := d.GetProviderMeta(&m)
	if err != nil {
		log.Printf("[WARN] error retrieving provider_meta")
		return baseUserAgent
	}

	if m.ModuleName != "" {
		return strings.Join([]string{m.ModuleName, baseUserAgent}, " ")
	}
	return baseUserAgent
}

func (c *Config) tfSdkUserAgent(suffix string) string {
	sdkModulePath := "github.com/hashicorp/terraform-plugin-sdk/v2"
	baseUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s",
		c.TerraformVersion, moduleVersionFromBuild(sdkModulePath))
	baseUserAgent = appendUserAgentFromEnv(baseUserAgent)
	userAgent := fmt.Sprintf("%s terraform-provider-equinix/%s %s", baseUserAgent, version.ProviderVersion, suffix)
	return strings.TrimSpace(userAgent)
}

func (c *Config) tfFrameworkUserAgent(suffix string) string {
	frameworkModulePath := "github.com/hashicorp/terraform-plugin-framework"
	baseUserAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin Framework/%s",
		c.TerraformVersion, moduleVersionFromBuild(frameworkModulePath))
	baseUserAgent = appendUserAgentFromEnv(baseUserAgent)
	userAgent := fmt.Sprintf("%s terraform-provider-equinix/%s %s", baseUserAgent, version.ProviderVersion, suffix)
	return strings.TrimSpace(userAgent)
}

func moduleVersionFromBuild(modulePath string) string {
	buildInfo, ok := debug.ReadBuildInfo()

	if !ok {
		return "buildinfo-failed"
	}

	for _, dependency := range buildInfo.Deps {
		if dependency.Path == modulePath {
			return dependency.Version
		}
	}

	return "unknown-version"
}
