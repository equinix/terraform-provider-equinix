package metal

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-metal/version"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/packethost/packngo"
)

const (
	consumerToken = "aZ9GmqHTPtxevvFq9SK3Pi2yr9YCbRzduCSXF2SNem5sjB91mDq7Th3ZwTtRqMWZ"
	uaEnvVar      = "TF_APPEND_USER_AGENT"
)

type Config struct {
	terraformVersion string
	AuthToken        string
	MaxRetries       int
	MaxRetryWait     time.Duration
}

var redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

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

// Client returns a new client for accessing Equinix Metal's API.
func (c *Config) Client() *packngo.Client {
	transport := logging.NewTransport("Equinix Metal", http.DefaultTransport)
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = transport
	retryClient.RetryMax = c.MaxRetries
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryWaitMax = c.MaxRetryWait
	retryClient.CheckRetry = MetalRetryPolicy
	standardClient := retryClient.StandardClient()

	client := packngo.NewClientWithAuth(consumerToken, c.AuthToken, standardClient)
	tfUserAgent := terraformUserAgent(c.terraformVersion)
	userAgent := fmt.Sprintf("%s terraform-provider-metal/%s %s",
		tfUserAgent, version.ProviderVersion, client.UserAgent)

	client.UserAgent = strings.TrimSpace(userAgent)

	return client
}
