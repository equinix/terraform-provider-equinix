package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccECXL2ConnectionAccepter(t *testing.T) {
	context := map[string]interface{}{
		"connectionResourceName": "test_connection",
		"name":                   "tf-single-aws",
		"profile_uuid":           "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73",
		"speed":                  200,
		"speed_unit":             "MB",
		"Notifications":          []string{"marry@equinix.com", "john@equinix.com"},
		"purchase_order_number":  "1234567890",
		"port_uuid":              "febc9d80-11e0-4dc8-8eb8-c41b6b378df2",
		"vlan_stag":              777,
		"vlan_ctag":              1000,
		"seller_region":          "us-east-1",
		"seller_metro_code":      "SV",
		"authorization_key":      "1234456",
		"accepterResourceName":   "test_accepter",
		"access_key":             "AKIAGGJKJU7BC3QJKYQ",
		"secret_key":             "CXGJW1lWbqENEqSkBAK",
	}

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


func testAccECXl2Connection_base(ctx map[string]interface{}) string {
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
	vlan_ctag = %{vlan_ctag}
	seller_region = "%{seller_region}"
	seller_metro_code = "%{seller_metro_code}"
	authorization_key = "%{authorization_key}"
}
`, ctx)
}

func testAccECXL2ConnectionAccepter(ctx map[string]interface{}) string {
	return testAccECXl2Connection_base(ctx) + nprintf(`
# Accepter
resource "equinix_ecx_l2_connection_accepter" "%{accepterResourceName}" {
   connection_id = "${equinix_ecx_l2_connection.test.id}"
   access_key = "%{access_key}"
   secret_key = %{secret_key}
}
`, ctx)
}