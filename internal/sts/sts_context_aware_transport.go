package sts

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type StsContextAwareTransport struct {
	// Source supplies the token to add to outgoing requests'
	// Authorization headers.
	Source *StsContextAwareTokenSource

	// Base is the base RoundTripper used to make HTTP requests.
	// If nil, http.DefaultTransport is used.
	Base http.RoundTripper
}

func (c *Config) New() *StsContextAwareTransport {
	return &StsContextAwareTransport{
		Source: c.StsTokenSource(),
	}
}

// RoundTrip authorizes and authenticates the request with an
// access token from ContextAwareTransport's Source.
func (t *StsContextAwareTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				//nolint:errcheck // Inherited from upstream; disabling lint to avoid a larger refactor
				req.Body.Close()
			}
		}()
	}

	// passing in the existing request context
	token, err := t.Source.OidcTokenExchange(req.Context())
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	req2 := cloneRequest(req) // per RoundTripper contract
	token.SetAuthHeader(req2)

	// req.Body is assumed to be closed by the base RoundTripper.
	reqBodyClosed = true
	return t.base().RoundTrip(req2)
}

var cancelOnce sync.Once

// CancelRequest does nothing. It used to be a legacy cancellation mechanism
// but now only it only logs on first use to warn that it's deprecated.
//
// Deprecated: use contexts for cancellation instead.
func (t *StsContextAwareTransport) CancelRequest(_ *http.Request) {
	cancelOnce.Do(func() {
		log.Printf("deprecated: golang.org/x/oauth2: StsContextAwareTransport.CancelRequest no longer does anything; use contexts")
	})
}

func (t *StsContextAwareTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
