// Package config uses environment variables to configure API clients for the provider
package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/extensions/equinixoauth2"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/version"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
	"golang.org/x/oauth2"
)

// These constants track environment variable names
// that are relevant to the provider
const (
	EndpointEnvVar       = "EQUINIX_API_ENDPOINT"
	ClientIDEnvVar       = "EQUINIX_API_CLIENTID"
	ClientSecretEnvVar   = "EQUINIX_API_CLIENTSECRET"
	ClientTokenEnvVar    = "EQUINIX_API_TOKEN"
	ClientTimeoutEnvVar  = "EQUINIX_API_TIMEOUT"
	MetalAuthTokenEnvVar = "METAL_AUTH_TOKEN"
)

// ProviderMeta allows passing additional metadata
// specific to our provider in provider config blocks
// See https://developer.hashicorp.com/terraform/internals/provider-meta
type ProviderMeta struct {
	ModuleName string `cty:"module_name"`
}

const (
	consumerToken = "aZ9GmqHTPtxevvFq9SK3Pi2yr9YCbRzduCSXF2SNem5sjB91mDq7Th3ZwTtRqMWZ"
	metalBasePath = "/metal/v1/"
	uaEnvVar      = "TF_APPEND_USER_AGENT"
)

var (
	// DefaultBaseURL is the default URL to use for API requests
	DefaultBaseURL = "https://api.equinix.com"
	// DefaultTimeout is the default request timeout to use for API requests
	DefaultTimeout = 30
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
	c.Metal = c.NewMetalClient()
	return nil
}

func (c *Config) newAuthClient() *http.Client {
	var authTransport http.RoundTripper
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

	authClient := http.Client{
		Timeout: c.requestTimeout(),
		//nolint:staticcheck // We should move to subsystem loggers, but that is a much bigger change
		Transport: logging.NewTransport("Equinix", authTransport),
	}
	return &authClient
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
	//nolint:staticcheck // We should move to subsystem loggers, but that is a much bigger change
	transport := logging.NewTransport("Equinix Fabric (fabricv4)", c.authClient.Transport)

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = transport
	retryClient.HTTPClient.Timeout = c.requestTimeout()
	retryClient.RetryMax = c.MaxRetries
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryWaitMax = c.MaxRetryWait
	standardClient := retryClient.StandardClient()

	baseURL, _ := url.Parse(c.BaseURL)

	configuration := fabricv4.NewConfiguration()
	configuration.Servers = fabricv4.ServerConfigurations{
		fabricv4.ServerConfiguration{
			URL: baseURL.String(),
		},
	}
	configuration.HTTPClient = standardClient
	configuration.AddDefaultHeader("X-SOURCE", "API")
	configuration.AddDefaultHeader("X-CORRELATION-ID", correlationId(25))
	client := fabricv4.NewAPIClient(configuration)

	return client
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

// NewMetalClientForSDK returns a new equinix-sdk-go client that enables
// an SDK resource to interact with the Metal API
func (c *Config) NewMetalClientForSDK(d *schema.ResourceData) *metalv1.APIClient {
	client := c.newMetalClient()

	baseUserAgent := c.tfSdkUserAgent(client.GetConfig().UserAgent)
	client.GetConfig().UserAgent = generateModuleUserAgentString(d, baseUserAgent)

	return client
}

// NewMetalClientForFramework returns a new equinix-sdk-go client that enables
// an framework resource to interact with the Metal API
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

// AddModuleToNEUserAgent injects the ModuleName into SDK resource metadata for analytics
func (c *Config) AddModuleToNEUserAgent(client *ne.Client, d *schema.ResourceData) {
	cli := *client
	rc := cli.(*ne.RestClient)
	rc.SetHeader("User-agent", generateModuleUserAgentString(d, c.neUserAgent))
	*client = rc
}

// AddFwModuleToMetalUserAgent injects the ModuleName into framework resource metadata for analytics
// TODO (ocobleseqx) - known issue, Metal services are initialized using the metal client pointer
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

// AddModuleToMetalUserAgent injects the ModuleName into SDK resource metadata for analytics
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
