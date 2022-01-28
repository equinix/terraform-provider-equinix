//Package oauth2 provides support for making oAuth2 authorized and authenticated HTTP requests
//for interactions with Equinix APIs, in particular Equinix specific client credencials grant type
package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

const (
	tokenPath       = "/oauth2/v1/token"
	defTokenTimeout = 3600
)

//Config describes oauth2 client credentials flow
type Config struct {
	// ClientID is the application's ID.
	ClientID string
	// ClientSecret is the application's secret.
	ClientSecret string
	// BaseURL is the base endpoint of a server that  token endpoint
	BaseURL string
}

//Error describes oauth2 err
type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("oauth2: error when acquiring token: code: %v, message %v", e.Code, e.Message)
}

//New creates *http.Client with Equinix oAuth2 tokensource.
//The returned client is not valid beyond the lifetime of the context.
func (c *Config) New(ctx context.Context) *http.Client {
	return c.NewWithClient(ctx, nil)
}

//NewWithClient creates *http.Client with Equinix oAuth2 tokensource and custom *http.Client.
//The returned client is not valid beyond the lifetime of the context.
func (c *Config) NewWithClient(ctx context.Context, hc *http.Client) *http.Client {
	return oauth2.NewClient(ctx, c.TokenSource(ctx, hc))
}

//TokenSource returns a TokenSource that returns t until t expires,
//automatically refreshing it as necessary using the provided context and the
//client ID and client secret.
func (c *Config) TokenSource(ctx context.Context, hc *http.Client) oauth2.TokenSource {
	var restClient *resty.Client
	if hc == nil {
		restClient = resty.New()
	} else {
		restClient = resty.NewWithClient(hc)
	}
	source := &tokenSource{
		ctx,
		c,
		restClient}
	return oauth2.ReuseTokenSource(nil, source)
}

type tokenSource struct {
	ctx  context.Context
	conf *Config
	*resty.Client
}

type tokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	TokenTimeout string `json:"token_timeout"`
	RefreshToken string `json:"refresh_token"`
}

type tokenError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (c *tokenSource) Token() (*oauth2.Token, error) {
	req := tokenRequest{"client_credentials", c.conf.ClientID, c.conf.ClientSecret}
	result := &tokenResponse{}
	resp, err := c.R().
		SetContext(c.ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-agent", "equinix/oauth2-go").
		SetBody(&req).
		SetResult(result).
		SetError(&tokenError{}).
		Post(c.conf.BaseURL + tokenPath)

	if err != nil {
		return nil, fmt.Errorf("oauth2: failed to fetch token: %s", err)
	}
	if resp.IsError() {
		respError := resp.Error().(*tokenError)
		return nil, Error{Code: respError.ErrorCode, Message: respError.ErrorMessage}
	}
	token := oauth2.Token{
		AccessToken:  result.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: result.RefreshToken}

	timeout, err := strconv.Atoi(result.TokenTimeout)
	if err != nil {
		timeout = defTokenTimeout
	}
	if timeout != 0 {
		token.Expiry = time.Now().Add(time.Duration(timeout) * time.Second)
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("oauth2: server response missing access_token")
	}
	return &token, nil
}
