package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/ecx-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFabricL2Connection_Port_Single_Accepter_AWS(t *testing.T) {
	portName, _ := schema.EnvDefaultFunc(priPortEnvVar, "sit-001-CX-SV1-NL-Dot1q-BO-10G-PRI-JUN-33")()
	spName, _ := schema.EnvDefaultFunc(awsSpEnvVar, "AWS Direct Connect")()
	context := map[string]interface{}{
		"port-resourceName":            "test",
		"port-name":                    portName.(string),
		"connection-resourceName":      "test",
		"connection-name":              fmt.Sprintf("tf-tst-%s", randString(6)),
		"connection-profile_name":      spName.(string),
		"connection-speed":             50,
		"connection-speed_unit":        "MB",
		"connection-notifications":     []string{"marry@equinix.com", "john@equinix.com"},
		"connection-vlan_stag":         randInt(2000),
		"connection-seller_region":     "us-west-2",
		"connection-seller_metro_code": "SV",
		"connection-authorization_key": "123456789012",
		"accepter-resourceName":        "test",
		"accepter-access_key":          "AKIAGGJKJU7BC3QJKYQ",
		"accepter-secret_key":          "CXGJW1lWbqENEqSkBAK",
	}
	connectionResourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["connection-resourceName"].(string))
	accepterResourceName := fmt.Sprintf("equinix_ecx_l2_connection_accepter.%s", context["accepter-resourceName"].(string))
	var testConn ecx.L2Connection
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withPort().withConnection().withAccepter().build(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(accepterResourceName, "connection_id"),
					resource.TestCheckResourceAttrPair(accepterResourceName, "connection_id", connectionResourceName, "id"),
					resource.TestCheckResourceAttrSet(accepterResourceName, "aws_connection_id"),
					testAccFabricL2ConnectionExists(connectionResourceName, &testConn),
					testAccFabricL2ConnectionAccepterStatus(&testConn),
				),
			},
			{
				ResourceName:      accepterResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func (t *testAccConfig) withAccepter() *testAccConfig {
	t.config += testAccFabricL2ConnectionAccepter(t.ctx)
	return t
}

func testAccFabricL2ConnectionAccepter(ctx map[string]interface{}) string {
	return nprintf(`
resource "equinix_ecx_l2_connection_accepter" "%{accepter-resourceName}" {
   connection_id = equinix_ecx_l2_connection.%{connection-resourceName}.id
   access_key    = "%{accepter-access_key}"
   secret_key    = "%{accepter-secret_key}"
}
`, ctx)
}

func testAccFabricL2ConnectionAccepterStatus(conn *ecx.L2Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ecx.StringValue(conn.ProviderStatus) != ecx.ConnectionStatusProvisioned {
			return fmt.Errorf("provider_status does not match %v - %v", ecx.ConnectionStatusProvisioned, conn.ProviderStatus)
		}
		return nil
	}
}
