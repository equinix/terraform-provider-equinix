package networkedge

import (
	"context"
	"net/http"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/equinix/terraform-provider-equinix/internal/networkedge"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sshUserResourceName = "equinix_network_ssh_user"

var sshUserDeviceIDs = []string{"52c00d7f-c310-458e-9426-1d7549e1f600", "5f1483f4-c479-424d-98c5-43a266aae25c"}

func sshUserResource(t *testing.T) *schema.Resource {
	t.Helper()
	r, ok := equinix.Provider().ResourcesMap[sshUserResourceName]
	require.True(t, ok, "provider has resource %s", sshUserResourceName)
	return r
}

func TestSSHUser_create(t *testing.T) {
	// given
	api := networkedge.NewMockAPI(t)
	d := schema.TestResourceDataRaw(t, sshUserResource(t).Schema, map[string]any{
		"username":   "user",
		"password":   "secret",
		"device_ids": []any{sshUserDeviceIDs[0], sshUserDeviceIDs[1]},
	})

	// when
	diags := sshUserResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)

	// The API takes one device in the create request. The other devices are added one by one.
	posts := api.RequestsFor(http.MethodPost, networkedge.SSHUsersPath)
	require.Len(t, posts, 1, "one create request was sent")
	assert.Equal(t, "user", posts[0].Body["username"], "username matches")
	assert.Equal(t, "secret", posts[0].Body["password"], "password matches")

	createdWith, ok := posts[0].Body["deviceUuid"].(string)
	require.True(t, ok, "create request names a device")
	remaining := sshUserDeviceIDs[0]
	if createdWith == sshUserDeviceIDs[0] {
		remaining = sshUserDeviceIDs[1]
	}
	assert.Contains(t, sshUserDeviceIDs, createdWith, "create request names one of the devices")
	assert.Len(t, api.RequestsFor(http.MethodPost, networkedge.SSHUsersPath+"/"+d.Id()+"/devices/"+remaining), 1,
		"the other device is added to the user")
}

func TestSSHUser_read(t *testing.T) {
	// given
	const id = "0cb9759d-58ab-44e6-9c10-6a3cfd18cefb"
	api := networkedge.NewMockAPI(t)
	api.SSHUsers.Seed(id, map[string]any{
		"username":    "user",
		"deviceUUIDs": []any{sshUserDeviceIDs[0], sshUserDeviceIDs[1]},
	})
	d := schema.TestResourceDataRaw(t, sshUserResource(t).Schema, map[string]any{})
	d.SetId(id)

	// when
	diags := sshUserResource(t).ReadContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "read does not return errors: %v", diags)
	assert.Equal(t, "user", d.Get("username"), "username matches")
	assert.ElementsMatch(t, sshUserDeviceIDs, converters.SetToStringList(d.Get("device_ids").(*schema.Set)), "device IDs match")
}
