// Package restmock provides an in-memory fake of a JSON REST API for unit
// tests. Tests point a real API client at the fake server, so the code under
// test sends real HTTP requests and parses real HTTP responses. Tests then
// assert on the requests that the server received and on the responses that
// the code under test parses.
//
// The package does not depend on any Equinix service. Service-specific
// packages register the collections (paths and required properties) of their
// own API.
package restmock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// Request is one HTTP request that the server received.
type Request struct {
	Method string
	Path   string
	// Body is the decoded JSON request body. It is nil if the request had no body.
	Body map[string]any
}

// Server is a fake REST API server.
type Server struct {
	srv *httptest.Server

	mu          sync.Mutex
	collections []*Collection
	nextID      int
	requests    []Request
}

// Collection stores the JSON documents of one REST collection, for example
// /v1/things. POST on the collection path creates a document. GET on the
// collection path followed by an identifier returns one document.
type Collection struct {
	srv      *Server
	basePath string
	required []string
	items    map[string]map[string]any
}

// NewServer starts a fake REST API server. The server stops when the test ends.
func NewServer(t *testing.T) *Server {
	t.Helper()
	s := &Server{}
	s.srv = httptest.NewServer(http.HandlerFunc(s.serve))
	t.Cleanup(s.srv.Close)
	return s
}

// URL returns the base URL of the server.
func (s *Server) URL() string {
	return s.srv.URL
}

// AddCollection registers a collection at basePath. POST requests must contain
// all properties in required.
func (s *Server) AddCollection(basePath string, required ...string) *Collection {
	c := &Collection{
		srv:      s,
		basePath: strings.TrimSuffix(basePath, "/"),
		required: required,
		items:    make(map[string]map[string]any),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.collections = append(s.collections, c)
	return c
}

// RequestsFor returns the requests with a given method and path that the
// server received, oldest first.
func (s *Server) RequestsFor(method, path string) []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Request
	for _, r := range s.requests {
		if r.Method == method && r.Path == path {
			out = append(out, r)
		}
	}
	return out
}

// Seed stores a document under a given identifier, as if something outside
// of the test had created it. The "uuid" property is set to id.
func (c *Collection) Seed(id string, doc map[string]any) {
	c.srv.mu.Lock()
	defer c.srv.mu.Unlock()
	cp := make(map[string]any, len(doc)+1)
	for k, v := range doc {
		cp[k] = v
	}
	cp["uuid"] = id
	c.items[id] = cp
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	rec := Request{Method: r.Method, Path: r.URL.Path}
	raw, err := io.ReadAll(r.Body)
	if err == nil && len(raw) > 0 {
		err = json.Unmarshal(raw, &rec.Body)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = append(s.requests, rec)
	if err != nil {
		writeError(w, http.StatusBadRequest, "request body is not valid JSON: "+err.Error())
		return
	}
	for _, c := range s.collections {
		switch {
		case r.URL.Path == c.basePath && r.Method == http.MethodPost:
			c.create(w, rec)
			return
		case strings.HasPrefix(r.URL.Path, c.basePath+"/") && r.Method == http.MethodGet:
			c.get(w, strings.TrimPrefix(r.URL.Path, c.basePath+"/"))
			return
		}
	}
	writeError(w, http.StatusNotFound, fmt.Sprintf("restmock has no route for %s %s", r.Method, r.URL.Path))
}

// create stores the request body as a new document. The caller holds the server lock.
func (c *Collection) create(w http.ResponseWriter, rec Request) {
	for _, prop := range c.required {
		if v, ok := rec.Body[prop]; !ok || v == nil || v == "" {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("required property %q is missing from the request body", prop))
			return
		}
	}
	c.srv.nextID++
	id := fmt.Sprintf("00000000-0000-4000-8000-%012d", c.srv.nextID)
	doc := make(map[string]any, len(rec.Body)+1)
	for k, v := range rec.Body {
		doc[k] = v
	}
	doc["uuid"] = id
	c.items[id] = doc
	w.Header().Set("Location", c.basePath+"/"+id)
	writeJSON(w, http.StatusCreated, map[string]any{"uuid": id})
}

// get returns one document. The caller holds the server lock.
func (c *Collection) get(w http.ResponseWriter, id string) {
	doc, ok := c.items[id]
	if !ok {
		writeError(w, http.StatusNotFound, "resource "+id+" does not exist")
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"errorCode": "EQ-MOCK", "errorMessage": message})
}
