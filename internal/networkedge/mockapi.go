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
	FilesPath    = "/ne/v1/files"
	BGPPath      = "/ne/v1/bgp"
	DevicesPath  = "/ne/v1/devices"
	LinksPath    = "/ne/v1/links"

	ACLTemplatesPath = "/ne/v1/aclTemplates"
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

	// Files holds the uploaded files. A file is uploaded with a multipart
	// form request and has the properties fileName, metroCode, deviceTypeCode,
	// processType and status.
	Files *restmock.Collection

	// BGP holds the BGP configurations. They are provisioned as soon as they are created.
	BGP *restmock.Collection

	// Devices holds the virtual devices. Tests add the devices that other
	// resources refer to with Seed.
	Devices *restmock.Collection

	// Links holds the device link groups. They are provisioned as soon as they are created.
	Links *restmock.Collection

	// ACLTemplates holds the device ACL templates.
	ACLTemplates *restmock.Collection
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
	sshUsers.NewDocument = func(req restmock.Request) map[string]any {
		return map[string]any{
			"username":    req.Body["username"],
			"deviceUUIDs": []any{req.Body["deviceUuid"]},
		}
	}
	// POST /ne/v1/sshUsers/{id}/devices/{deviceId} gives the user access to a device.
	sshUsers.AddAction(http.MethodPost, "devices/", func(doc map[string]any, deviceID string) {
		devices, _ := doc["deviceUUIDs"].([]any)
		doc["deviceUUIDs"] = append(devices, deviceID)
	})

	files := srv.AddCollection(FilesPath, "metroCode", "deviceTypeCode", "processType")
	files.CreatedIDProperty = "fileUuid"
	files.NewDocument = func(req restmock.Request) map[string]any {
		return map[string]any{
			"fileName":       req.Files["file"].Name,
			"metroCode":      req.Body["metroCode"],
			"deviceTypeCode": req.Body["deviceTypeCode"],
			"processType":    req.Body["processType"],
			"status":         "SUCCESS",
		}
	}

	bgp := srv.AddCollection(BGPPath, "connectionUuid", "localIpAddress", "localAsn", "remoteIpAddress", "remoteAsn")
	bgp.NewDocument = func(req restmock.Request) map[string]any {
		doc := map[string]any{"provisioningStatus": "PROVISIONED"}
		for k, v := range req.Body {
			doc[k] = v
		}
		return doc
	}

	links := srv.AddCollection(LinksPath, "groupName")
	links.NewDocument = func(req restmock.Request) map[string]any {
		doc := map[string]any{"status": "PROVISIONED"}
		for k, v := range req.Body {
			doc[k] = v
		}
		return doc
	}

	return &MockAPI{
		Server:   srv,
		Config:   cfg,
		SSHKeys:  srv.AddCollection(SSHKeysPath, "keyName", "keyValue"),
		SSHUsers: sshUsers,
		Files:    files,
		BGP:      bgp,
		Devices:  srv.AddCollection(DevicesPath),
		Links:    links,

		ACLTemplates: srv.AddCollection(ACLTemplatesPath, "name", "description"),
	}
}
