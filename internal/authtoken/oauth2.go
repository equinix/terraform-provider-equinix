package authtoken

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/equinix/equinix-sdk-go/services/accesstokenv1"
	"golang.org/x/oauth2"
)

const (
	defTokenTimeout = 3600
)

// Config describes oauth2 client credentials flow
type Config struct {
	// ClientID is the application's ID.
	ClientID string
	// ClientSecret is the application's secret.
	ClientSecret string
	// BaseURL is the base endpoint of a server that  token endpoint
	BaseURL string
}

// Error describes oauth2 err
type Error struct {
	Code    string
	Message string
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided context and the
// client ID and client secret.
func (c *Config) TokenSource() *TokenSource {
	config := accesstokenv1.NewConfiguration()
	config.Servers = accesstokenv1.ServerConfigurations{
		accesstokenv1.ServerConfiguration{
			URL: c.BaseURL,
		},
	}
	restClient := accesstokenv1.NewAPIClient(config)
	source := TokenSource{
		c,
		restClient,
		sync.Mutex{},
		nil,
	}
	return &source
}

type TokenSource struct {
	conf   *Config
	client *accesstokenv1.APIClient
	mu     sync.Mutex
	token  *oauth2.Token
}

func (s *TokenSource) TokenWithContext(ctx context.Context) (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.token.Valid() {
		req := accesstokenv1.Oauth2TokenRequest{
			GrantType:    accesstokenv1.PtrString("client_credentials"),
			ClientId:     s.conf.ClientID,
			ClientSecret: s.conf.ClientSecret,
		}
		result, _, err := s.client.OAuth2TokenApi.GetOAuth2AccessToken(ctx).Payload(req).Execute()

		if err != nil {
			return nil, fmt.Errorf("oauth2: failed to fetch token: %s", err)
		}

		token := oauth2.Token{
			AccessToken:  result.AccessToken,
			TokenType:    "Bearer",
			RefreshToken: result.GetRefreshToken(),
		}

		timeout, err := strconv.Atoi(result.TokenTimeout)
		if err != nil {
			timeout = defTokenTimeout
		}
		if timeout != 0 {
			token.Expiry = time.Now().Add(time.Duration(timeout) * time.Second)
		}
		if token.AccessToken == "" {
			return nil, fmt.Errorf("oauth2: server response missing access token")
		}
		s.token = &token
	}

	return s.token, nil
}
