package equinix

import (
	"fmt"
	"testing"

	"ecx-go/v3"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	priPortEnvVar = "TF_ACC_ECX_PRI_DOT1Q_PORT_ID"
	secPortEnvVar = "TF_ACC_ECX_SEC_DOT1Q_PORT_ID"
	awsSpEnvVar   = "TF_ACC_ECX_L2_AWS_SP_ID"
	azureSpEnvVar = "TF_ACC_ECX_L2_AZURE_SP_ID"
)

func TestAccECXL2ConnectionAWSDot1Q(t *testing.T) {
	t.Parallel()
	portID, _ := schema.EnvDefaultFunc(priPortEnvVar, "0cbc3f5e-308e-4043-916f-2a2788818bf3")()
	spID, _ := schema.EnvDefaultFunc(awsSpEnvVar, "496e149f-cc04-4ef6-af8f-ac657da78be6")()
	context := map[string]interface{}{
		"resourceName":          "tf-aws-dot1q",
		"name":                  fmt.Sprintf("tf-tst-%s", randString(6)),
		"profile_uuid":          spID.(string),
		"speed":                 50,
		"speed_unit":            "MB",
		"notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"purchase_order_number": "1234567890",
		"port_uuid":             portID.(string),
		"vlan_stag":             randInt(2000),
		"seller_region":         "us-west-2",
		"seller_metro_code":     "SV",
		"authorization_key":     "123456789012",
	}
	resourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["resourceName"].(string))
	var testConn ecx.L2Connection
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2ConnectionAWSDot1Q(context),
				Check: resource.ComposeTestCheckFunc(
					testAccECXL2ConnectionExists(resourceName, &testConn),
					testAccECXL2ConnectionAttributes(&testConn, context),
				),
			},
		},
	})
}

func TestAccECXL2ConnectionAzureDot1QPub(t *testing.T) {
	t.Parallel()
	priPortID, _ := schema.EnvDefaultFunc(priPortEnvVar, "0cbc3f5e-308e-4043-916f-2a2788818bf3")()
	secPortID, _ := schema.EnvDefaultFunc(secPortEnvVar, "e3bad661-cdb0-484e-9699-ec5492bd623b")()
	spID, _ := schema.EnvDefaultFunc(azureSpEnvVar, "9915c38f-568f-4c06-b7a2-cef6b6e67847")()
	context := map[string]interface{}{
		"resourceName":          "tf-azure-dot1q-pub",
		"name":                  fmt.Sprintf("tf-tst-%s", randString(6)),
		"profile_uuid":          spID.(string),
		"speed":                 50,
		"speed_unit":            "MB",
		"notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"purchase_order_number": "1234567890",
		"port_uuid":             priPortID.(string),
		"vlan_stag":             randInt(2000),
		"seller_region":         "us-west-2",
		"seller_metro_code":     "SV",
		"authorization_key":     "123456789016",
		"named_tag":             "Public",
		"secondary_name":        fmt.Sprintf("tf-tst-%s", randString(6)),
		"secondary_port_uuid":   secPortID.(string),
		"secondary_vlan_stag":   randInt(2000),
	}
	resourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["resourceName"].(string))
	var primary, secondary ecx.L2Connection
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2ConnectionAzureDot1QPub(context),
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

func testAccECXL2ConnectionAWSDot1Q(ctx map[string]interface{}) string {
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
 seller_region = "%{seller_region}"
 seller_metro_code = "%{seller_metro_code}"
 authorization_key = "%{authorization_key}"
}
`, ctx)
}

func testAccECXL2ConnectionAzureDot1QPub(ctx map[string]interface{}) string {
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
  seller_region = "%{seller_region}"
  seller_metro_code = "%{seller_metro_code}"
  authorization_key = "%{authorization_key}"
  named_tag = "%{named_tag}"
  secondary_connection {
    name = "%{secondary_name}"
    port_uuid = "%{secondary_port_uuid}"
    vlan_stag = %{secondary_vlan_stag}
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
		if v, ok := ctx["named_tag"]; ok && conn.NamedTag != v.(string) {
			return fmt.Errorf("named_tag does not match %v - %v", conn.NamedTag, v)
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
