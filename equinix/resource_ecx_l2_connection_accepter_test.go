package equinix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveAWSCredentials_Basic(t *testing.T) {
	// Given
	key := "testKey"
	secret := "testSecret"
	profileName := "testProfile"
	d := schema.TestResourceDataRaw(t, createECXL2ConnectionAccepterResourceSchema(),
		map[string]interface{}{
			ecxL2ConnectionAccepterSchemaNames["AccessKey"]: key,
			ecxL2ConnectionAccepterSchemaNames["SecretKey"]: secret,
			ecxL2ConnectionAccepterSchemaNames["Profile"]:   profileName,
		})

	// when
	creds, err := retrieveAWSCredentials(d)

	// then
	assert.Nil(t, err, "Error is not returned")
	assert.NotNil(t, creds, "Credentials value is returned")
	assert.Equal(t, key, creds.AccessKeyID, "AccessKeyID matches")
	assert.Equal(t, secret, creds.SecretAccessKey, "SecretAccessKey matches")
}
