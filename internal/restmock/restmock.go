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
	"net/url"
	"sort"
	"strings"
	"sync"
	"testing"
)

// Request is one HTTP request that the server received.
type Request struct {
	Method string
	Path   string
	// Query holds the query parameters of the request.
	Query url.Values
	// Body is the decoded JSON request body. For a multipart form request it
	// holds the text fields of the form. It is nil if the request had no body.
	Body map[string]any
	// Files holds the files of a multipart form request, by form field name.
	Files map[string]File
}

// File is a file that a multipart form request uploaded.
type File struct {
	Name    string
	Content string
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
// collection path followed by an identifier returns one document, and PUT
// merges the request body into it.
type Collection struct {
	srv      *Server
	basePath string
	required []string
	items    map[string]map[string]any
	actions  []action

	// NewDocument converts a create request into the stored document. If it
	// is nil, the request body is stored as it is.
	NewDocument func(req Request) map[string]any

	// CreatedIDProperty is the property of the create response that holds the
	// identifier of the new document. The default is "uuid".
	CreatedIDProperty string
}

type action struct {
	method string
	prefix string
	fn     func(doc map[string]any, rest string)
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

// AddAction registers a handler for requests with a given method on
// basePath/{id}/{prefix}{rest}, where {id} is the identifier of a stored
// document. The handler receives the document, which it can change, and {rest}.
func (c *Collection) AddAction(method, prefix string, fn func(doc map[string]any, rest string)) {
	c.srv.mu.Lock()
	defer c.srv.mu.Unlock()
	c.actions = append(c.actions, action{method: method, prefix: prefix, fn: fn})
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
	rec := Request{Method: r.Method, Path: r.URL.Path, Query: r.URL.Query()}
	var err error
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err = readMultipart(r, &rec)
	} else {
		var raw []byte
		raw, err = io.ReadAll(r.Body)
		if err == nil && len(raw) > 0 {
			err = json.Unmarshal(raw, &rec.Body)
		}
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
		case r.URL.Path == c.basePath && r.Method == http.MethodGet:
			c.list(w)
			return
		case strings.HasPrefix(r.URL.Path, c.basePath+"/"):
			if c.serveItem(w, rec) {
				return
			}
		}
	}
	writeError(w, http.StatusNotFound, fmt.Sprintf("restmock has no route for %s %s", r.Method, r.URL.Path))
}

// list returns all documents in one page, in the pagination format of the
// Network Edge API. Tests that use it must keep the collection smaller than
// one page. The caller holds the server lock.
func (c *Collection) list(w http.ResponseWriter) {
	ids := make([]string, 0, len(c.items))
	for id := range c.items {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	data := make([]map[string]any, 0, len(ids))
	for _, id := range ids {
		data = append(data, c.items[id])
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"pagination": map[string]any{"offset": 0, "limit": len(data), "total": len(data)},
		"data":       data,
	})
}

// serveItem handles requests below a document path. It returns false if it has no handler
// for the request. The caller holds the server lock.
func (c *Collection) serveItem(w http.ResponseWriter, rec Request) bool {
	id, sub, hasSub := strings.Cut(strings.TrimPrefix(rec.Path, c.basePath+"/"), "/")
	if !hasSub && rec.Method == http.MethodGet {
		c.get(w, id)
		return true
	}
	if !hasSub && rec.Method == http.MethodPut {
		c.update(w, id, rec)
		return true
	}
	if !hasSub {
		return false
	}
	for _, a := range c.actions {
		rest, ok := strings.CutPrefix(sub, a.prefix)
		if !ok || a.method != rec.Method {
			continue
		}
		doc, found := c.items[id]
		if !found {
			writeError(w, http.StatusNotFound, "resource "+id+" does not exist")
			return true
		}
		a.fn(doc, rest)
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
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
	body := rec.Body
	if c.NewDocument != nil {
		body = c.NewDocument(rec)
	}
	doc := make(map[string]any, len(body)+1)
	for k, v := range body {
		doc[k] = v
	}
	doc["uuid"] = id
	c.items[id] = doc
	w.Header().Set("Location", c.basePath+"/"+id)
	idProperty := c.CreatedIDProperty
	if idProperty == "" {
		idProperty = "uuid"
	}
	writeJSON(w, http.StatusCreated, map[string]any{idProperty: id})
}

// update merges the request body into a document. The caller holds the server lock.
func (c *Collection) update(w http.ResponseWriter, id string, rec Request) {
	doc, ok := c.items[id]
	if !ok {
		writeError(w, http.StatusNotFound, "resource "+id+" does not exist")
		return
	}
	for k, v := range rec.Body {
		doc[k] = v
	}
	writeJSON(w, http.StatusOK, map[string]any{"uuid": id})
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

// readMultipart copies the fields and files of a multipart form request into rec.
func readMultipart(r *http.Request, rec *Request) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return err
	}
	rec.Body = make(map[string]any)
	for k, v := range r.MultipartForm.Value {
		rec.Body[k] = v[0]
	}
	rec.Files = make(map[string]File)
	for field, headers := range r.MultipartForm.File {
		f, err := headers[0].Open()
		if err != nil {
			return err
		}
		content, err := io.ReadAll(f)
		_ = f.Close()
		if err != nil {
			return err
		}
		rec.Files[field] = File{Name: headers[0].Filename, Content: string(content)}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"errorCode": "EQ-MOCK", "errorMessage": message})
}
