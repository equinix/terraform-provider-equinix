package equinix

import (
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

var sshUserFields = []string{"Username", "Password", "DeviceUUIDs"}

func TestNeSSHUser_resourceDataFromDomain(t *testing.T) {
	//given
	d := schema.TestResourceDataRaw(t, createNeSSHUserResourceSchema(), make(map[string]interface{}))
	user := ne.SSHUser{
		Username:    "test",
		Password:    "qwerty",
		DeviceUUIDs: []string{"one", "two", "three"}}

	//when
	err := updateNeSSHUserResource(&user, d)

	//then
	assert.Nil(t, err, "Schema update should not return an error")
	sourceMatchesTargetSchema(t, user, sshUserFields, d, neSSHUserSchemaNames)
}

func TestNeSSHUser_domainFromResourceData(t *testing.T) {
	//given
	d := schema.TestResourceDataRaw(t, createNeSSHUserResourceSchema(), make(map[string]interface{}))
	d.Set(neSSHUserSchemaNames["Username"], "test")
	d.Set(neSSHUserSchemaNames["Password"], "qwerty")
	d.Set(neSSHUserSchemaNames["DeviceUUIDs"], []string{"one", "two", "three"})

	//when
	user := createNeSSHUser(d)

	//then
	sourceMatchesTargetSchema(t, user, sshUserFields, d, neSSHUserSchemaNames)
}
