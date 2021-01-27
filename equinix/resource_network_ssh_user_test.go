package equinix

import (
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

var sshUserFields = []string{"Username", "Password", "DeviceUUIDs"}

func TestNetworkSSHUser_resourceDataFromDomain(t *testing.T) {
	//given
	d := schema.TestResourceDataRaw(t, createNetworkSSHUserResourceSchema(), make(map[string]interface{}))
	user := ne.SSHUser{
		Username:    ne.String("test"),
		Password:    ne.String("qwerty"),
		DeviceUUIDs: []string{"one", "two", "three"}}

	//when
	err := updateNetworkSSHUserResource(&user, d)

	//then
	assert.Nil(t, err, "Schema update should not return an error")
	//sourceMatchesTargetSchema(t, user, sshUserFields, d, networkSSHUserSchemaNames)
}

func TestNetworkSSHUser_domainFromResourceData(t *testing.T) {
	//given
	d := schema.TestResourceDataRaw(t, createNetworkSSHUserResourceSchema(), make(map[string]interface{}))
	d.Set(networkSSHUserSchemaNames["Username"], "test")
	d.Set(networkSSHUserSchemaNames["Password"], "qwerty")
	d.Set(networkSSHUserSchemaNames["DeviceUUIDs"], []string{"one", "two", "three"})

	//when
	user := createNetworkSSHUser(d)

	//then
	assert.NotNil(t, user, "User is not nil")
	//sourceMatchesTargetSchema(t, user, sshUserFields, d, networkSSHUserSchemaNames)
}
