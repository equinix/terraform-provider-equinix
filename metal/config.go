package metal

import (
	"bytes"
	"context"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/packethost/packngo"
)

const (
	consumerToken = "aZ9GmqHTPtxevvFq9SK3Pi2yr9YCbRzduCSXF2SNem5sjB91mDq7Th3ZwTtRqMWZ"
)

type Config struct {
	AuthToken string
}

func newRetryableTransport(defaultTransport http.RoundTripper) http.RoundTripper {
	c := retryablehttp.NewClient()
	c.HTTPClient = &http.Client{Transport: defaultTransport}

	c.RetryMax = 10
	c.RetryWaitMax = 30 * time.Second
	c.RetryWaitMin = time.Second
	c.CheckRetry = MetalRetryPolicy
	return &retryableTransport{c}
}

type retryableTransport struct {
	*retryablehttp.Client
}

// RoundTrip wraps calling an HTTP method with retries.
func (c *retryableTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	var body io.ReadSeeker
	if r.Body != nil {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bs)
	}
	req, err := retryablehttp.NewRequest(r.Method, r.URL.String(), body)
	if err != nil {
		return nil, err
	}
	for key, val := range r.Header {
		req.Header.Set(key, val[0])
	}
	return c.Client.Do(req)
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

// Client returns a new client for accessing Equinix Metal's API.
func (c *Config) Client() *packngo.Client {
	transport := logging.NewTransport("Equinix Metal", http.DefaultTransport)
	httpClient := &http.Client{Transport: newRetryableTransport(transport)}
	return packngo.NewClientWithAuth(consumerToken, c.AuthToken, httpClient)
}
