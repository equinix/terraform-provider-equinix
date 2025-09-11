// Package sts provides a context-aware token source for Equinix STS (Secure Token Service).
package sts

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/equinix/equinix-sdk-go/services/stsv1alpha"
	"golang.org/x/oauth2"
)

// ContextAwareTokenSource Implements the refresh token source
type ContextAwareTokenSource struct {
	conf   *Config
	client *stsv1alpha.APIClient
	mu     sync.Mutex
	token  *oauth2.Token
}

// OidcTokenExchange performs an OIDC token exchange using the configured STS client and settings.
// It ensures thread safety, validates required configuration, and caches the token until expiry.
// Returns a valid OAuth2 token or an error if the exchange fails.
func (s *ContextAwareTokenSource) OidcTokenExchange(ctx context.Context) (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.token.Valid() {
		if s.conf.StsAuthScope == "" {
			return nil, fmt.Errorf("authorization scope cannot be empty for OIDC token exchange")
		}

		if s.conf.StsSourceToken == "" {
			return nil, fmt.Errorf("sts source token cannot be empty for OIDC token exchange")
		}

		// Execute token exchange
		response, httpResp, err := s.client.UseApi.UseTokenPost(ctx).
			GrantType(stsv1alpha.USETOKENPOSTREQUESTGRANTTYPE_URN_IETF_PARAMS_OAUTH_GRANT_TYPE_TOKEN_EXCHANGE).
			Scope(s.conf.StsAuthScope).
			SubjectToken(s.conf.StsSourceToken).
			SubjectTokenType(stsv1alpha.USETOKENPOSTREQUESTSUBJECTTOKENTYPE_URN_IETF_PARAMS_OAUTH_TOKEN_TYPE_ID_TOKEN).
			Execute()

		if err != nil {
			var httpRespBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					httpRespBody = string(bodyBytes)
				} else {
					httpRespBody = fmt.Sprintf("failed to read http response body: %v", readErr)
				}
				// Optionally reset the body if it needs to be read again elsewhere
				err := httpResp.Body.Close()
				if err != nil {
					return nil, err
				}
			}
			return nil, fmt.Errorf("sts token exchange failed with response body: %s and error: %s", httpRespBody, err)
		}

		token := oauth2.Token{
			AccessToken: response.AccessToken,
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
		}

		if token.AccessToken == "" {
			return nil, fmt.Errorf("sts server response missing access token")
		}
		s.token = &token
	}
	return s.token, nil
}
