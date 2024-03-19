package ssh_user

import (
	"testing"

	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNetworkSSHUser_resourceFromResourceData(t *testing.T) {
	// given
	rawData := map[string]interface{}{
		networkSSHUserSchemaNames["Username"]: "user",
		networkSSHUserSchemaNames["Password"]: "secret",
	}
	deviceUUIDs := []string{"52c00d7f-c310-458e-9426-1d7549e1f600", "5f1483f4-c479-424d-98c5-43a266aae25c"}
	d := schema.TestResourceDataRaw(t, createNetworkSSHUserResourceSchema(), rawData)
	d.Set(networkSSHUserSchemaNames["DeviceUUIDs"], deviceUUIDs)
	expected := ne.SSHUser{
		Username:    ne.String(rawData[networkSSHUserSchemaNames["Username"]].(string)),
		Password:    ne.String(rawData[networkSSHUserSchemaNames["Password"]].(string)),
		DeviceUUIDs: deviceUUIDs,
	}

	// when
	result := createNetworkSSHUser(d)

	// then
	assert.Equal(t, expected, result, "Result matches expected result")
}

func TestNetworkSSHUser_updateResourceData(t *testing.T) {
	// given
	d := schema.TestResourceDataRaw(t, createNetworkSSHUserResourceSchema(), make(map[string]interface{}))
	input := ne.SSHUser{
		Username:    ne.String("user"),
		Password:    ne.String("secret"),
		DeviceUUIDs: []string{"52c00d7f-c310-458e-9426-1d7549e1f600", "5f1483f4-c479-424d-98c5-43a266aae25c"},
	}
	// when
	err := updateNetworkSSHUserResource(&input, d)
	// then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ne.StringValue(input.Username), d.Get(networkSSHUserSchemaNames["Username"]), "Username matches")
	assert.Equal(t, ne.StringValue(input.Password), d.Get(networkSSHUserSchemaNames["Password"]), "Password matches")
	assert.Equal(t, input.DeviceUUIDs, converters.SetToStringList(d.Get(networkSSHUserSchemaNames["DeviceUUIDs"]).(*schema.Set)), "DeviceUUIDs matches")
}
