package networkedge

import (
	"context"
	"net/http"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/networkedge"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sshKeyResourceName = "equinix_network_ssh_key"

func sshKeyResource(t *testing.T) *schema.Resource {
	t.Helper()
	r, ok := equinix.Provider().ResourcesMap[sshKeyResourceName]
	require.True(t, ok, "provider has resource %s", sshKeyResourceName)
	return r
}

func sshKeyData(t *testing.T, id string, raw map[string]any) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, sshKeyResource(t).Schema, raw)
	d.SetId(id)
	return d
}

func TestSSHKey_create(t *testing.T) {
	// given
	api := networkedge.NewMockAPI(t)
	d := sshKeyData(t, "", map[string]any{
		"name":       "testKey",
		"public_key": "ssh-rsa AAAAtest",
		"type":       "RSA",
		"project_id": "68ccfd49-39b1-478e-957a-67c72f719d7a",
	})

	// when
	diags := sshKeyResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)

	posts := api.RequestsFor(http.MethodPost, networkedge.SSHKeysPath)
	require.Len(t, posts, 1, "one create request was sent")
	assert.Equal(t, map[string]any{
		"keyName":   "testKey",
		"keyValue":  "ssh-rsa AAAAtest",
		"keyType":   "RSA",
		"projectId": "68ccfd49-39b1-478e-957a-67c72f719d7a",
	}, posts[0].Body, "create request body matches the API contract")
}

func TestSSHKey_read(t *testing.T) {
	// given
	const id = "059c3020-aec5-44ca-816c-235435f16df9"
	api := networkedge.NewMockAPI(t)
	api.SSHKeys.Seed(id, map[string]any{
		"keyName":  "testKey",
		"keyValue": "ssh-rsa AAAAtest",
		"keyType":  "RSA",
	})
	d := sshKeyData(t, id, map[string]any{})

	// when
	diags := sshKeyResource(t).ReadContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "read does not return errors: %v", diags)
	assert.Equal(t, id, d.Get("uuid"), "UUID matches")
	assert.Equal(t, "testKey", d.Get("name"), "name matches")
	assert.Equal(t, "ssh-rsa AAAAtest", d.Get("public_key"), "public key matches")
	assert.Equal(t, "RSA", d.Get("type"), "type matches")
}
