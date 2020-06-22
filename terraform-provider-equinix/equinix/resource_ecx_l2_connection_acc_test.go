package equinix

import (
	"fmt"
	"testing"

	"ecx-go/v3"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccECXL2Connection_single(t *testing.T) {
	context := map[string]interface{}{
		"resourceName":          "aws_dot1q",
		"name":                  "tf-single-aws",
		"profile_uuid":          "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73",
		"speed":                 200,
		"speed_unit":            "MB",
		"notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"purchase_order_number": "1234567890",
		"port_uuid":             "febc9d80-11e0-4dc8-8eb8-c41b6b378df2",
		"vlan_stag":             777,
		"vlan_ctag":             1000,
		"seller_region":         "us-east-1",
		"seller_metro_code":     "SV",
		"authorization_key":     "1234456",
	}

	resourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["resourceName"].(string))
	var testConn ecx.L2Connection
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2Connection_single(context),
				Check: resource.ComposeTestCheckFunc(
					testAccECXL2ConnectionExists(resourceName, &testConn),
					testAccECXL2ConnectionAttributes(&testConn, context),
				),
			},
		},
	})
}

func TestAccECXL2Connection_redundant(t *testing.T) {
	context := map[string]interface{}{
		"resourceName":              "redundant_self",
		"name":                      "tf-redundant-self",
		"profile_uuid":              "2a4f7e27-dff8-4f15-aeda-a11ffe9ccf73",
		"speed":                     50,
		"speed_unit":                "MB",
		"notifications":             []string{"marry@equinix.com", "john@equinix.com"},
		"purchase_order_number":     "1234567890",
		"port_uuid":                 "febc9d80-11e0-4dc8-8eb8-c41b6b378df2",
		"vlan_stag":                 800,
		"zside_port_uuid":           "03a969b5-9cea-486d-ada0-2a4496ed72fb",
		"zside_vlan_stag":           1010,
		"seller_region":             "us-east-1",
		"seller_metro_code":         "SV",
		"secondary_name":            "tf-redundant-self-sec",
		"secondary_port_uuid":       "86872ae5-ca19-452b-8e69-bb1dd5f93bd1",
		"secondary_vlan_stag":       999,
		"secondary_vlan_ctag":       1000,
		"secondary_zside_port_uuid": "393b2f6e-9c66-4a39-adac-820120555420",
		"secondary_zside_vlan_stag": 1022,
	}

	resourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["resourceName"].(string))
	var primary, secondary ecx.L2Connection
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2Connection_redundant(context),
				Check: resource.ComposeTestCheckFunc(
					testAccECXL2ConnectionExists(resourceName, &primary),
					testAccECXL2ConnectionAttributes(&primary, context),
					testAccECXL2ConnectionSecondaryExists(&primary, &secondary),
					testAccECXL2ConnectionSecondaryAttributes(&secondary, context),
				),
			},
		},
	})
}

func testAccECXL2Connection_single(ctx map[string]interface{}) string {
	return nprintf(`
resource "equinix_ecx_l2_connection" "%{resourceName}" {
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

func testAccECXL2Connection_redundant(ctx map[string]interface{}) string {
	return nprintf(`
resource "equinix_ecx_l2_connection" "%{resourceName}" {
  name = "%{name}"
  profile_uuid = "%{profile_uuid}"
  speed = %{speed}
  speed_unit = "%{speed_unit}"
  notifications = %{notifications}
  purchase_order_number = "%{purchase_order_number}"
  port_uuid = "%{port_uuid}"
  vlan_stag = %{vlan_stag}
  zside_port_uuid = "%{zside_port_uuid}"
  zside_vlan_stag = %{zside_vlan_stag}
  seller_region = "%{seller_region}"
  seller_metro_code = "%{seller_metro_code}"
  secondary_connection {
    name = "%{secondary_name}"
    port_uuid = "%{secondary_port_uuid}"
    vlan_stag = %{secondary_vlan_stag}
    vlan_ctag = %{secondary_vlan_ctag}
    zside_port_uuid = "%{secondary_zside_port_uuid}"
    zside_vlan_stag = %{secondary_zside_vlan_stag}
  }
}
`, ctx)
}

func testAccECXL2ConnectionExists(resourceName string, conn *ecx.L2Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		client := testAccProvider.Meta().(*Config).ecx
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}

		resp, err := client.GetL2Connection(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching L2 connection %v", err)
		}
		if resp.UUID != rs.Primary.ID {
			return fmt.Errorf("resource ID does not match %v - %v", rs.Primary.ID, resp.UUID)
		}
		*conn = *resp
		return nil
	}
}

func testAccECXL2ConnectionSecondaryExists(primary *ecx.L2Connection, secondary *ecx.L2Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Config).ecx
		if primary.RedundantUUID == "" {
			return fmt.Errorf("primary connection has no RedundantUUID set")
		}
		resp, err := client.GetL2Connection(primary.RedundantUUID)
		if err != nil {
			return fmt.Errorf("error when fetching L2 connection %v", err)
		}
		*secondary = *resp
		return nil
	}
}

func testAccECXL2ConnectionAttributes(conn *ecx.L2Connection, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["name"]; ok && conn.Name != v.(string) {
			return fmt.Errorf("name does not match %v - %v", conn.Name, v)
		}
		if v, ok := ctx["profile_uuid"]; ok && conn.ProfileUUID != v.(string) {
			return fmt.Errorf("profileUUID does not match %v - %v", conn.ProfileUUID, v)
		}
		if v, ok := ctx["speed"]; ok && conn.Speed != v.(int) {
			return fmt.Errorf("speed does not match %v - %v", conn.Speed, v)
		}
		if v, ok := ctx["speed_unit"]; ok && conn.SpeedUnit != v.(string) {
			return fmt.Errorf("speedUnit does not match %v - %v", conn.SpeedUnit, v)
		}
		if v, ok := ctx["notifications"]; ok && !slicesMatch(conn.Notifications, v.([]string)) {
			return fmt.Errorf("notifications does not match %v - %v", conn.Notifications, v)
		}
		if v, ok := ctx["purchase_order_number"]; ok && conn.PurchaseOrderNumber != v.(string) {
			return fmt.Errorf("purchaseOrderNumber does not match %v - %v", conn.PurchaseOrderNumber, v)
		}
		if v, ok := ctx["port_uuid"]; ok && conn.PortUUID != v.(string) {
			return fmt.Errorf("portUUID does not match %v - %v", conn.PortUUID, v)
		}
		if v, ok := ctx["vlan_stag"]; ok && conn.VlanSTag != v.(int) {
			return fmt.Errorf("vlanSTag does not match %v - %v", conn.VlanSTag, v)
		}
		if v, ok := ctx["vlan_ctag"]; ok && conn.VlanCTag != v.(int) {
			return fmt.Errorf("vlanCTag does not match %v - %v", conn.VlanCTag, v)
		}
		if v, ok := ctx["zside_port_uuid"]; ok && conn.ZSidePortUUID != v.(string) {
			return fmt.Errorf("zSidePortUUID does not match %v - %v", conn.ZSidePortUUID, v)
		}
		if v, ok := ctx["zside_vlan_stag"]; ok && conn.ZSideVlanSTag != v.(int) {
			return fmt.Errorf("zSideVlanSTag does not match %v - %v", conn.ZSideVlanSTag, v)
		}
		if v, ok := ctx["zside_vlan_ctag"]; ok && conn.ZSideVlanCTag != v.(int) {
			return fmt.Errorf("zSideVlanCTag does not match %v - %v", conn.ZSideVlanCTag, v)
		}
		if v, ok := ctx["seller_region"]; ok && conn.SellerRegion != v.(string) {
			return fmt.Errorf("sellerRegion does not match %v - %v", conn.SellerRegion, v)
		}
		if v, ok := ctx["seller_metro_code"]; ok && conn.SellerMetroCode != v.(string) {
			return fmt.Errorf("sellerMetroCode does not match %v - %v", conn.SellerMetroCode, v)
		}
		if v, ok := ctx["authorization_key"]; ok && conn.AuthorizationKey != v.(string) {
			return fmt.Errorf("authorizationKey does not match %v - %v", conn.AuthorizationKey, v)
		}
		return nil
	}
}

func testAccECXL2ConnectionSecondaryAttributes(conn *ecx.L2Connection, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["secondary_name"]; ok && conn.Name != v.(string) {
			return fmt.Errorf("name does not match %v - %v", conn.Name, v)
		}
		if v, ok := ctx["secondary_port_uuid"]; ok && conn.PortUUID != v.(string) {
			return fmt.Errorf("portUUID does not match %v - %v", conn.PortUUID, v)
		}
		if v, ok := ctx["secondary_vlan_stag"]; ok && conn.VlanSTag != v.(int) {
			return fmt.Errorf("vlanSTag does not match %v - %v", conn.VlanSTag, v)
		}
		if v, ok := ctx["secondary_vlan_ctag"]; ok && conn.VlanCTag != v.(int) {
			return fmt.Errorf("vlanCTag does not match %v - %v", conn.VlanCTag, v)
		}
		if v, ok := ctx["secondary_zside_port_uuid"]; ok && conn.ZSidePortUUID != v.(string) {
			return fmt.Errorf("zSidePortUUID does not match %v - %v", conn.ZSidePortUUID, v)
		}
		if v, ok := ctx["secondary_zside_vlan_stag"]; ok && conn.ZSideVlanSTag != v.(int) {
			return fmt.Errorf("zSideVlanSTag does not match %v - %v", conn.ZSideVlanSTag, v)
		}
		if v, ok := ctx["secondary_zside_vlan_ctag"]; ok && conn.ZSideVlanCTag != v.(int) {
			return fmt.Errorf("zSideVlanCTag does not match %v - %v", conn.ZSideVlanCTag, v)
		}
		return nil
	}
}

func slicesMatch(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	visited := make([]bool, len(s1))
	for i := 0; i < len(s1); i++ {
		found := false
		for j := 0; j < len(s2); j++ {
			if visited[j] {
				continue
			}
			if s1[i] == s2[j] {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
