package equinix

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestAccECXL2ConnectionAccepter(t *testing.T) {
	t.Parallel()
	portID, _ := schema.EnvDefaultFunc(priPortEnvVar, "0cbc3f5e-308e-4043-916f-2a2788818bf3")()
	spID, _ := schema.EnvDefaultFunc(awsSpEnvVar, "496e149f-cc04-4ef6-af8f-ac657da78be6")()
	context := map[string]interface{}{
		"connectionResourceName": "tf-aws-dot1q",
		"name":                   fmt.Sprintf("tf-tst-%s", randString(6)),
		"profile_uuid":           spID.(string),
		"speed":                  50,
		"speed_unit":             "MB",
		"notifications":          []string{"marry@equinix.com", "john@equinix.com"},
		"purchase_order_number":  "1234567890",
		"port_uuid":              portID.(string),
		"vlan_stag":              randInt(2000),
		"seller_region":          "us-west-2",
		"seller_metro_code":      "SV",
		"authorization_key":      "123456789012",
		"accepterResourceName":   "tf-accepter",
		"access_key":             "AKIAGGJKJU7BC3QJKYQ",
		"secret_key":             "CXGJW1lWbqENEqSkBAK",
	}
	log.Printf("Config is %s", testAccECXL2ConnectionAccepter(context))
	connectionResourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["connectionResourceName"].(string))
	accepterResourceName := fmt.Sprintf("equinix_ecx_l2_connection_accepter.%s", context["accepterResourceName"].(string))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2ConnectionAccepter(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(accepterResourceName, "connection_id"),
					resource.TestCheckResourceAttrPair(accepterResourceName, "connection_id", connectionResourceName, "id"),
				),
			},
		},
	})
}

func testAccECXL2ConnectionAccepter(ctx map[string]interface{}) string {
	return nprintf(`
# Connection
resource "equinix_ecx_l2_connection" "%{connectionResourceName}" {
  name = "%{name}"
  profile_uuid = "%{profile_uuid}"
  speed = %{speed}
  speed_unit = "%{speed_unit}"
  notifications = %{notifications}
  purchase_order_number = "%{purchase_order_number}"
  port_uuid = "%{port_uuid}"
  vlan_stag = %{vlan_stag}
  seller_region = "%{seller_region}"
  seller_metro_code = "%{seller_metro_code}"
  authorization_key = "%{authorization_key}"
}
# Accepter
resource "equinix_ecx_l2_connection_accepter" "%{accepterResourceName}" {
   connection_id = equinix_ecx_l2_connection.%{connectionResourceName}.id
   access_key = "%{access_key}"
   secret_key = "%{secret_key}"
}
`, ctx)
}
