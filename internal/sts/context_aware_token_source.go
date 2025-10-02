// Package sts provides a context-aware token source for Equinix STS (Secure Token Service).
package sts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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

	if s.token != nil && s.token.Valid() {
		return s.token, nil
	}

	if err := s.validateConfig(); err != nil {
		return nil, err
	}

	token, err := s.executeTokenExchangeWithRetry(ctx)
	if err != nil {
		return nil, err
	}

	s.token = token
	return s.token, nil
}

func (s *ContextAwareTokenSource) executeTokenExchangeWithRetry(ctx context.Context) (*oauth2.Token, error) {
	maxRetries := 5
	baseDelay := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		response, httpResp, err := s.client.UseApi.UseTokenPost(ctx).
			GrantType(stsv1alpha.USETOKENPOSTREQUESTGRANTTYPE_URN_IETF_PARAMS_OAUTH_GRANT_TYPE_TOKEN_EXCHANGE).
			Scope(s.conf.StsAuthScope).
			SubjectToken(s.conf.StsSourceToken).
			SubjectTokenType(stsv1alpha.USETOKENPOSTREQUESTSUBJECTTOKENTYPE_URN_IETF_PARAMS_OAUTH_TOKEN_TYPE_ID_TOKEN).
			Execute()

		if err == nil {
			return s.createTokenFromResponse(response)
		}

		if !s.shouldRetry(httpResp) {
			return nil, s.formatError(httpResp, err)
		}

		// Apply backoff delay
		delay := time.Duration(1<<attempt) * baseDelay
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("max retries exceeded")
}

func (s *ContextAwareTokenSource) validateConfig() error {
	if s.conf.StsAuthScope == "" {
		return fmt.Errorf("authorization scope cannot be empty for OIDC token exchange")
	}
	if s.conf.StsSourceToken == "" {
		return fmt.Errorf("sts source token cannot be empty for OIDC token exchange")
	}
	return nil
}

func (s *ContextAwareTokenSource) shouldRetry(httpResp *http.Response) bool {
	return httpResp != nil && httpResp.StatusCode == 409
}

func (s *ContextAwareTokenSource) createTokenFromResponse(response *stsv1alpha.TokenExchangeResponse) (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: response.AccessToken,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("sts server response missing access token")
	}
	return token, nil
}

func (s *ContextAwareTokenSource) formatError(httpResp *http.Response, err error) error {
	if httpResp == nil {
		return fmt.Errorf("STS token exchange failed: %w", err)
	}

	// Read response body for additional error details
	var bodyBytes []byte
	if httpResp.Body != nil {
		bodyBytes, _ = io.ReadAll(httpResp.Body)
		_ = httpResp.Body.Close()
	}

	errorMsg := fmt.Sprintf("STS token exchange failed with status %d", httpResp.StatusCode)

	if len(bodyBytes) > 0 {
		errorMsg += fmt.Sprintf(": %s", string(bodyBytes))
	}

	if err != nil {
		errorMsg += fmt.Sprintf(" (underlying error: %v)", err)
	}

	return errors.New(errorMsg)
}
