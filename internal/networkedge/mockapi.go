package networkedge

import (
	"context"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/restmock"
)

// Paths of the Network Edge API collections that MockAPI serves.
const (
	SSHKeysPath = "/ne/v1/publicKeys"
)

// MockAPI is a fake Network Edge API for unit tests. Config holds a provider
// configuration whose Network Edge client sends its requests to the fake server.
type MockAPI struct {
	*restmock.Server

	// Config is the provider configuration to pass to resource and data source functions.
	Config *config.Config

	// SSHKeys holds the SSH public keys.
	SSHKeys *restmock.Collection
}

// NewMockAPI starts a fake Network Edge API that stops when the test ends.
func NewMockAPI(t *testing.T) *MockAPI {
	t.Helper()
	srv := restmock.NewServer(t)

	cfg := &config.Config{BaseURL: srv.URL(), Token: "test-token"}
	if err := cfg.Load(context.Background()); err != nil {
		t.Fatalf("loading provider config for the mock Network Edge API: %v", err)
	}

	return &MockAPI{
		Server:  srv,
		Config:  cfg,
		SSHKeys: srv.AddCollection(SSHKeysPath, "keyName", "keyValue"),
	}
}
