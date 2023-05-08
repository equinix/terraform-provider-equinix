package equinix

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/equinix/ecx-go/v2"
	"github.com/equinix/ne-go"
	"github.com/equinix/oauth2-go"
	"github.com/equinix/terraform-provider-equinix/version"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/packethost/packngo"
	xoauth2 "golang.org/x/oauth2"
)

type DumpTransport struct {
	r http.RoundTripper
}

func (d *DumpTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	dump, _ := httputil.DumpRequestOut(h, true)
	fmt.Printf("****REQUEST****\n%q\n", dump)
	resp, err := d.r.RoundTrip(h)
	dump, _ = httputil.DumpResponse(resp, true)
	fmt.Printf("****RESPONSE****\n%q\n****************\n\n", dump)
	return resp, err
}

const (
	consumerToken         = "aZ9GmqHTPtxevvFq9SK3Pi2yr9YCbRzduCSXF2SNem5sjB91mDq7Th3ZwTtRqMWZ"
	metalBasePath         = "/metal/v1/"
	uaEnvVar              = "TF_APPEND_USER_AGENT"
	emptyCredentialsError = `the provider needs to be configured with the proper credentials before it
can be used.

One of pair "client_id" - "client_secret" or "token" must be set in the provider
configuration to interact with Equinix Fabric and Network Edge services, and
"auth_token" to interact with Equinix Metal. These can also be configured using
environment variables.

Please note that while the authentication arguments are individually optional to allow
interaction with the different services independently, trying to provision the resources
of a service without the required credentials will return an API error referring to
'Invalid authentication token' or 'error when acquiring token'.

More information on the provider configuration can be found here:
https://registry.terraform.io/providers/equinix/equinix/latest/docs`
)

var (
	DefaultBaseURL   = "https://uatapi.equinix.com"
	DefaultTimeout   = 30
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)
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

	ecx   ecx.Client
	ne    ne.Client
	metal *packngo.Client

	ecxUserAgent   string
	neUserAgent    string
	metalUserAgent string

	terraformVersion string
	fabricClient     *v4.APIClient
	FabricAuthToken  string
}

// Load function validates configuration structure fields and configures
// all required API clients.
func (c *Config) Load(ctx context.Context) error {
	if c.BaseURL == "" {
		return fmt.Errorf("'baseURL' cannot be empty")
	}

	if c.Token == "" && (c.ClientID == "" || c.ClientSecret == "") && c.AuthToken == "" {
		return fmt.Errorf(emptyCredentialsError)
	}

	var authClient *http.Client
	if c.Token != "" {
		tokenSource := xoauth2.StaticTokenSource(&xoauth2.Token{AccessToken: c.Token})
		oauthTransport := &xoauth2.Transport{
			Source: tokenSource,
		}
		authClient = &http.Client{
			Transport: oauthTransport,
		}
	} else {
		authConfig := oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			BaseURL:      c.BaseURL,
		}
		authClient = authConfig.New(ctx)

		if c.ClientID != "" && c.ClientSecret != "" {
			tke, err := authConfig.TokenSource(ctx, authClient).Token()
			if err != nil {
				if err != nil {
					return err
				}
			}
			if tke != nil {
				c.FabricAuthToken = tke.AccessToken
			}
		}
	}

	if c.FabricAuthToken == "" {
		c.FabricAuthToken = c.Token
	}
	authClient.Timeout = c.requestTimeout()
	authClient.Transport = logging.NewTransport("Equinix", authClient.Transport)
	ecxClient := ecx.NewClient(ctx, c.BaseURL, authClient)
	neClient := ne.NewClient(ctx, c.BaseURL, authClient)

	if c.PageSize > 0 {
		ecxClient.SetPageSize(c.PageSize)
		neClient.SetPageSize(c.PageSize)
	}
	c.ecxUserAgent = c.fullUserAgent("equinix/ecx-go")
	ecxClient.SetHeaders(map[string]string{
		"User-agent": c.ecxUserAgent,
	})
	c.neUserAgent = c.fullUserAgent("equinix/ecx-go")
	neClient.SetHeaders(map[string]string{
		"User-agent": c.neUserAgent,
	})

	c.ecx = ecxClient
	c.ne = neClient
	c.metal = c.NewMetalClient()
	c.fabricClient = c.NewFabricClient()
	return nil
}

// NewFabricClient returns a new client for accessing Equinix Fabric's v4 API.
func (c *Config) NewFabricClient() *v4.APIClient {
	transport := logging.NewTransport("Equinix Fabric", http.DefaultTransport)
	authClient := &http.Client{
		Transport: transport,
	}
	authClient.Timeout = c.requestTimeout()
	fabricHeaderMap := map[string]string{
		"X-SOURCE":         "API",
		"X-CORRELATION-ID": CorrelationId(25),
	}
	v4Configuration := v4.Configuration{
		BasePath:      c.BaseURL,
		DefaultHeader: fabricHeaderMap,
		UserAgent:     "equinix/fabric-go",
		HTTPClient:    authClient,
	}
	client := v4.NewAPIClient(&v4Configuration)
	return client
}

// NewMetalClient returns a new client for accessing Equinix Metal's API.
func (c *Config) NewMetalClient() *packngo.Client {
	transport := http.DefaultTransport
	// transport = &DumpTransport{http.DefaultTransport} // Debug only
	transport = logging.NewTransport("Equinix Metal", transport)
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = transport
	retryClient.RetryMax = c.MaxRetries
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryWaitMax = c.MaxRetryWait
	retryClient.CheckRetry = MetalRetryPolicy
	standardClient := retryClient.StandardClient()
	baseURL, _ := url.Parse(c.BaseURL)
	baseURL.Path = path.Join(baseURL.Path, metalBasePath) + "/"
	client, _ := packngo.NewClientWithBaseURL(consumerToken, c.AuthToken, standardClient, baseURL.String())
	client.UserAgent = c.fullUserAgent(client.UserAgent)
	c.metalUserAgent = client.UserAgent
	return client
}

func (c *Config) requestTimeout() time.Duration {
	if c.RequestTimeout == 0 {
		return 5 * time.Second
	}
	return c.RequestTimeout
}

func MetalRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, nil
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, nil
			}
		}
		// The error is likely recoverable so retry.
		return true, nil
	}
	return false, nil
}

func terraformUserAgent(version string) string {
	ua := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s",
		version, meta.SDKVersionString())

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

func (c *Config) addModuleToECXUserAgent(client *ecx.Client, d *schema.ResourceData) {
	cli := *client
	rc := cli.(*ecx.RestClient)
	rc.SetHeader("User-agent", generateModuleUserAgentString(d, c.ecxUserAgent))
	*client = rc
}

func (c *Config) addModuleToNEUserAgent(client *ne.Client, d *schema.ResourceData) {
	cli := *client
	rc := cli.(*ne.RestClient)
	rc.SetHeader("User-agent", generateModuleUserAgentString(d, c.neUserAgent))
	*client = rc
}

// TODO (ocobleseqx) - known issue, Metal services are initialized using the metal client pointer
// if two or more modules in same project interact with metal resources they will override
// the UserAgent resulting in swapped UserAgent.
// This can be fixed by letting the headers be overwritten on the initialized Packngo ServiceOp
// clients on a query-by-query basis.
func (c *Config) addModuleToMetalUserAgent(d *schema.ResourceData) {
	c.metal.UserAgent = generateModuleUserAgentString(d, c.metalUserAgent)
}

func generateModuleUserAgentString(d *schema.ResourceData, baseUserAgent string) string {
	var m providerMeta
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

func (c *Config) fullUserAgent(suffix string) string {
	tfUserAgent := terraformUserAgent(c.terraformVersion)
	userAgent := fmt.Sprintf("%s terraform-provider-equinix/%s %s", tfUserAgent, version.ProviderVersion, suffix)
	return strings.TrimSpace(userAgent)
}
