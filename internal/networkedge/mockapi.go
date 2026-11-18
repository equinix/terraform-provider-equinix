package networkedge

import (
	"context"
	"net/http"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/restmock"
)

// Paths of the Network Edge API collections that MockAPI serves.
const (
	SSHKeysPath  = "/ne/v1/publicKeys"
	SSHUsersPath = "/ne/v1/sshUsers"
)

// MockAPI is a fake Network Edge API for unit tests. Config holds a provider
// configuration whose Network Edge client sends its requests to the fake server.
type MockAPI struct {
	*restmock.Server

	// Config is the provider configuration to pass to resource and data source functions.
	Config *config.Config

	// SSHKeys holds the SSH public keys.
	SSHKeys *restmock.Collection

	// SSHUsers holds the SSH users. A user created through the API has the
	// properties username and deviceUUIDs. The API does not return passwords.
	SSHUsers *restmock.Collection
}

// NewMockAPI starts a fake Network Edge API that stops when the test ends.
func NewMockAPI(t *testing.T) *MockAPI {
	t.Helper()
	srv := restmock.NewServer(t)

	cfg := &config.Config{BaseURL: srv.URL(), Token: "test-token"}
	if err := cfg.Load(context.Background()); err != nil {
		t.Fatalf("loading provider config for the mock Network Edge API: %v", err)
	}

	sshUsers := srv.AddCollection(SSHUsersPath, "username", "password", "deviceUuid")
	// The create request names one device and a password. The stored user lists its devices.
	sshUsers.NewDocument = func(body map[string]any) map[string]any {
		return map[string]any{
			"username":    body["username"],
			"deviceUUIDs": []any{body["deviceUuid"]},
		}
	}
	// POST /ne/v1/sshUsers/{id}/devices/{deviceId} gives the user access to a device.
	sshUsers.AddAction(http.MethodPost, "devices/", func(doc map[string]any, deviceID string) {
		devices, _ := doc["deviceUUIDs"].([]any)
		doc["deviceUUIDs"] = append(devices, deviceID)
	})

	return &MockAPI{
		Server:   srv,
		Config:   cfg,
		SSHKeys:  srv.AddCollection(SSHKeysPath, "keyName", "keyValue"),
		SSHUsers: sshUsers,
	}
}
