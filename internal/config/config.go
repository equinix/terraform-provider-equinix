// Package config provides configuration and client initialization for the Equinix
// Terraform provider. It handles authentication to Equinix services (Fabric,
// Network Edge, and Metal) via OAuth, tokens, and Workload Identity Federation.
// The package manages HTTP client configuration including retries and timeouts,
// and provides utility functions for generating service-specific API clients
// compatible with both Terraform SDK v2 and Framework interfaces.
package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/ne-go"
	"github.com/equinix/oauth2-go"
	"github.com/equinix/terraform-provider-equinix/version"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
	xoauth2 "golang.org/x/oauth2"
)

const (
	// EndpointEnvVar is the environment variable name used to override the default Equinix API endpoint.
	EndpointEnvVar = "EQUINIX_API_ENDPOINT"
	// ClientIDEnvVar is the environment variable name for the Equinix API OAuth client ID.
	ClientIDEnvVar       = "EQUINIX_API_CLIENTID"
	ClientSecretEnvVar   = "EQUINIX_API_CLIENTSECRET"
	ClientTokenEnvVar    = "EQUINIX_API_TOKEN"
	ClientTimeoutEnvVar  = "EQUINIX_API_TIMEOUT"
	MetalAuthTokenEnvVar = "METAL_AUTH_TOKEN"
	AuthScopeEnvVar      = "EQUINIX_STS_AUTH_SCOPE"
	StsSourceTokenEnvVar = "EQUINIX_STS_SOURCE_TOKEN"
	StsEndpointEnvVar    = "EQUINIX_STS_ENDPOINT"
)

// ProviderMeta contains metadata about the Terraform module using this provider.
// It's primarily used to track module names for user agent identification in API requests.
type ProviderMeta struct {
	ModuleName string `cty:"module_name"`
}

const (
	consumerToken = "aZ9GmqHTPtxevvFq9SK3Pi2yr9YCbRzduCSXF2SNem5sjB91mDq7Th3ZwTtRqMWZ"
	metalBasePath = "/metal/v1/"
	uaEnvVar      = "TF_APPEND_USER_AGENT"
)

var (
	// DefaultBaseURL is the standard production API endpoint for Equinix services.
	DefaultBaseURL = "https://api.equinix.com"
	// DefaultStsBaseURL is the default Security Token Service (STS) endpoint
	DefaultStsBaseURL = "https://sts.eqix.equinix.com"
	DefaultTimeout    = 30
)

// Config is the configuration structure used to instantiate the Equinix
// provider.
type Config struct {
	BaseURL        string
	AuthToken      string
	ClientID       string
	ClientSecret   string
	MaxRetries     int
	MaxRetryWait   time.Duration
	RequestTimeout time.Duration
	PageSize       int
	Token          string
	StsAuthScope   string
	StsBaseURL     string
	StsSourceToken string

	authClient *http.Client

	Ne    ne.Client
	Metal *packngo.Client

	neUserAgent    string
	metalUserAgent string

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

	authClient := c.createAuthClient(ctx)
	authClient.Timeout = c.requestTimeout()
	//nolint:staticcheck // We should move to subsystem loggers, but that is a much bigger change
	authClient.Transport = logging.NewTransport("Equinix", authClient.Transport)
	c.authClient = authClient
	neClient := ne.NewClient(ctx, c.BaseURL, authClient)

	if c.PageSize > 0 {
		neClient.SetPageSize(c.PageSize)
	}
	c.neUserAgent = c.tfSdkUserAgent("equinix/ne-go")
	neClient.SetHeaders(map[string]string{
		"User-agent": c.neUserAgent,
	})

	c.Ne = neClient
	c.Metal = c.NewMetalClient()
	return nil
}

// NewFabricClientForSDK returns a terraform sdkv2 plugin compatible
// equinix-sdk-go/fabricv4 client to be used to access Fabric's V4 APIs
func (c *Config) NewFabricClientForSDK(ctx context.Context, d *schema.ResourceData) *fabricv4.APIClient {
	client := c.newFabricClient(ctx)

	baseUserAgent := c.tfSdkUserAgent(client.GetConfig().UserAgent)
	client.GetConfig().UserAgent = generateModuleUserAgentString(d, baseUserAgent)

	return client
}

// NewFabricClientForTesting creates a Fabric client configured specifically for acceptance tests.
// Deprecated: when the acceptance package starts to contain API clients for testing/cleanup this will move with them
func (c *Config) NewFabricClientForTesting(ctx context.Context) *fabricv4.APIClient {
	client := c.newFabricClient(ctx)

	client.GetConfig().UserAgent = fmt.Sprintf("tf-acceptance-tests %v", client.GetConfig().UserAgent)

	return client
}

// NewFabricClientForFramework creates a Fabric client compatible with the Terraform Plugin Framework.
func (c *Config) NewFabricClientForFramework(ctx context.Context, meta tfsdk.Config) *fabricv4.APIClient {
	client := c.newFabricClient(ctx)

	baseUserAgent := c.tfFrameworkUserAgent(client.GetConfig().UserAgent)
	client.GetConfig().UserAgent = generateFwModuleUserAgentString(ctx, meta, baseUserAgent)

	return client
}

// newFabricClient returns the base fabricv4 client that is then used for either the sdkv2 or framework
// implementations of the Terraform Provider with exported Methods
func (c *Config) newFabricClient(ctx context.Context) *fabricv4.APIClient {
	authClient := c.createAuthClient(ctx)

	// Configure HTTP client with retries and logging
	httpClient := c.configureHTTPClient(authClient)

	return c.createFabricClient(httpClient)
}

func (c *Config) createAuthClient(ctx context.Context) *http.Client {
	if c.StsAuthScope != "" && c.StsSourceToken != "" {
		// If the StsAuthScope and StsSourceToken are set, use STS based authentication
		return c.createStsIdentityClient(ctx, c.StsSourceToken)
	}

	return c.createOAuthClient(ctx)
}

func (c *Config) createStsIdentityClient(ctx context.Context, stsSourceToken string) *http.Client {
	httpClient := &http.Client{Timeout: c.requestTimeout()}

	refreshSource := &stsRefreshTokenSource{
		ctx:            ctx,
		authScope:      c.StsAuthScope,
		stsSourceToken: stsSourceToken,
		client:         httpClient,
		stsBaseURL:     c.StsBaseURL,
	}

	tokenSource := xoauth2.ReuseTokenSource(nil, refreshSource)

	return &http.Client{
		Transport: &xoauth2.Transport{
			Source: tokenSource,
		},
	}
}

func (c *Config) createOAuthClient(ctx context.Context) *http.Client {
	if c.Token != "" {
		tokenSource := xoauth2.StaticTokenSource(&xoauth2.Token{AccessToken: c.Token})
		oauthTransport := &xoauth2.Transport{
			Source: tokenSource,
		}
		var authClient = &http.Client{
			Transport: oauthTransport,
		}
		return authClient
	}
	authConfig := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		BaseURL:      c.BaseURL,
	}
	return authConfig.New(ctx)
}

func (c *Config) configureHTTPClient(authClient *http.Client) *http.Client {
	transport := logging.NewTransport("Equinix Fabric (fabricv4)", authClient.Transport)

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

// NewMetalClient returns a new packngo client for accessing Equinix Metal's API.
// Deprecated: migrate to NewMetalClientForSdk or NewMetalClientForFramework instead
func (c *Config) NewMetalClient() *packngo.Client {
	transport := http.DefaultTransport
	//nolint:staticcheck // We should move to subsystem loggers, but that is a much bigger change
	transport = logging.NewTransport("Equinix Metal (packngo)", transport)
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = transport
	retryClient.RetryMax = c.MaxRetries
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryWaitMax = c.MaxRetryWait
	standardClient := retryClient.StandardClient()
	baseURL, _ := url.Parse(c.BaseURL)
	baseURL.Path = path.Join(baseURL.Path, metalBasePath) + "/"
	client, _ := packngo.NewClientWithBaseURL(consumerToken, c.AuthToken, standardClient, baseURL.String())
	client.UserAgent = c.tfSdkUserAgent(client.UserAgent)
	c.metalUserAgent = client.UserAgent
	return client
}

// NewMetalClientForSDK returns a Metal client compatible with the Terraform Plugin SDK v2.
func (c *Config) NewMetalClientForSDK(d *schema.ResourceData) *metalv1.APIClient {
	client := c.newMetalClient()

	baseUserAgent := c.tfSdkUserAgent(client.GetConfig().UserAgent)
	client.GetConfig().UserAgent = generateModuleUserAgentString(d, baseUserAgent)

	return client
}

// NewMetalClientForFramework returns a Metal client compatible with the Terraform Plugin Framework.
func (c *Config) NewMetalClientForFramework(ctx context.Context, meta tfsdk.Config) *metalv1.APIClient {
	client := c.newMetalClient()

	baseUserAgent := c.tfFrameworkUserAgent(client.GetConfig().UserAgent)
	client.GetConfig().UserAgent = generateFwModuleUserAgentString(ctx, meta, baseUserAgent)

	return client
}

// NewMetalClientForTesting is a short-term shim to allow tests to continue to have a client for cleanup and validation
// code that is outside of the resource or datasource under test
// Deprecated: when possible, API clients for test cleanup/validation should be moved to the acceptance package
func (c *Config) NewMetalClientForTesting() *metalv1.APIClient {
	client := c.newMetalClient()

	client.GetConfig().UserAgent = fmt.Sprintf("tf-acceptance-tests %v", client.GetConfig().UserAgent)

	return client
}

func (c *Config) newMetalClient() *metalv1.APIClient {
	transport := http.DefaultTransport
	//nolint:staticcheck // We should move to subsystem loggers, but that is a much bigger change
	transport = logging.NewTransport("Equinix Metal (metal-go)", transport)
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = transport
	retryClient.RetryMax = c.MaxRetries
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryWaitMax = c.MaxRetryWait
	standardClient := retryClient.StandardClient()

	baseURL, _ := url.Parse(c.BaseURL)
	baseURL.Path = path.Join(baseURL.Path, metalBasePath)

	configuration := metalv1.NewConfiguration()
	configuration.Servers = metalv1.ServerConfigurations{
		metalv1.ServerConfiguration{
			URL: baseURL.String(),
		},
	}
	configuration.HTTPClient = standardClient
	configuration.AddDefaultHeader("X-Auth-Token", c.AuthToken)
	client := metalv1.NewAPIClient(configuration)
	return client
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

// AddModuleToNEUserAgent updates the Network Edge client's User-Agent header to include Terraform module information.
func (c *Config) AddModuleToNEUserAgent(client *ne.Client, d *schema.ResourceData) {
	cli := *client
	rc := cli.(*ne.RestClient)
	rc.SetHeader("User-agent", generateModuleUserAgentString(d, c.neUserAgent))
	*client = rc
}

// AddFwModuleToMetalUserAgent TODO (ocobleseqx) - known issue, Metal services are initialized using the metal client pointer
// if two or more modules in same project interact with metal resources they will override
// the UserAgent resulting in swapped UserAgent.
// This can be fixed by letting the headers be overwritten on the initialized Packngo ServiceOp
// clients on a query-by-query basis.
func (c *Config) AddFwModuleToMetalUserAgent(ctx context.Context, meta tfsdk.Config) {
	c.Metal.UserAgent = generateFwModuleUserAgentString(ctx, meta, c.metalUserAgent)
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

// AddModuleToMetalUserAgent updates the Metal client's User-Agent string to include Terraform module information.
func (c *Config) AddModuleToMetalUserAgent(d *schema.ResourceData) {
	c.Metal.UserAgent = generateModuleUserAgentString(d, c.metalUserAgent)
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

// Implement the refresh token source
type stsRefreshTokenSource struct {
	ctx            context.Context
	authScope      string
	stsSourceToken string
	client         *http.Client
	stsBaseURL     string
}

// Token implements the oauth2.TokenSource interface
func (s *stsRefreshTokenSource) Token() (*xoauth2.Token, error) {
	accessToken, expiry, err := oidcTokenExchange(
		s.ctx,
		s.authScope,
		s.stsSourceToken,
		s.client,
		s.stsBaseURL,
	)
	if err != nil {
		return nil, err
	}

	return &xoauth2.Token{
		AccessToken: accessToken,
		Expiry:      time.Now().Add(time.Duration(expiry) * time.Second),
	}, nil
}

func oidcTokenExchange(ctx context.Context, authScope string, stsSourceToken string, client *http.Client, stsBaseURL string) (string, int, error) {
	if authScope == "" {
		return "", 0, fmt.Errorf("authorization scope cannot be empty for OIDC token exchange")
	}

	if stsSourceToken == "" {
		return "", 0, fmt.Errorf("sts source token cannot be empty for OIDC token exchange")
	}

	baseURL, err := url.Parse(stsBaseURL)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse STS base URL: %w", err)
	}
	baseURL.Path = path.Join(baseURL.Path, "/use/token")
	tokenURL := baseURL.String()

	form := url.Values{
		"grant_type":         {"urn:ietf:params:oauth:grant-type:token-exchange"},
		"scope":              {authScope},
		"subject_token_type": {"urn:ietf:params:oauth:token-type:id_token"},
		"subject_token":      {stsSourceToken},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create token exchange request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer func() {
		if cerr := httpResp.Body.Close(); cerr != nil {
			log.Printf("[WARN] error closing token exchange response body: %v", cerr)
		}
	}()

	body, readErr := io.ReadAll(httpResp.Body)
	if readErr != nil {
		return "", 0, fmt.Errorf("error reading token exchange response body: %w", readErr)
	}

	if httpResp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("token exchange failed with status %d: %s", httpResp.StatusCode, string(body))
	}

	var response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", 0, fmt.Errorf("failed to parse token exchange response: %w", err)
	}
	return response.AccessToken, response.ExpiresIn, nil
}
