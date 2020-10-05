package equinix

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/ecx-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	priPortEnvVar = "TF_ACC_ECX_PRI_DOT1Q_PORT_NAME"
	secPortEnvVar = "TF_ACC_ECX_SEC_DOT1Q_PORT_NAME"
	awsSpEnvVar   = "TF_ACC_ECX_L2_AWS_SP_NAME"
	azureSpEnvVar = "TF_ACC_ECX_L2_AZURE_SP_NAME"
)

func init() {
	resource.AddTestSweepers("ECXL2Connection", &resource.Sweeper{
		Name: "ECXL2Connection",
		F:    testSweepECXL2Connections,
	})
}

func testSweepECXL2Connections(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	conns, err := config.ecx.GetL2OutgoingConnections([]string{
		ecx.ConnectionStatusNotAvailable,
		ecx.ConnectionStatusPendingAutoApproval,
		ecx.ConnectionStatusPendingBGPPeering,
		ecx.ConnectionStatusProvisioned,
		ecx.ConnectionStatusProvisioning,
		ecx.ConnectionStatusRejected,
	})
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching ECXL2Connection list: %s", err)
		return err
	}
	nonSweepableCount := 0
	for _, conn := range conns {
		if !isSweepableTestResource(conn.Name) {
			nonSweepableCount++
			continue
		}
		if err := config.ecx.DeleteL2Connection(conn.UUID); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting ECXL2Connection resource %s (%s): %s", conn.UUID, conn.Name, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for ECXL2Connection resource %s (%s)", conn.UUID, conn.Name)
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}

func TestAccECXL2ConnectionAWSDot1Q(t *testing.T) {
	t.Parallel()
	portName, _ := schema.EnvDefaultFunc(priPortEnvVar, "sit-001-CX-SV1-NL-Dot1q-BO-10G-PRI-JUN-33")()
	spName, _ := schema.EnvDefaultFunc(awsSpEnvVar, "AWS Direct Connect")()
	context := map[string]interface{}{
		"port_name":             portName.(string),
		"profile_name":          spName.(string),
		"resourceName":          "tst-aws-dot1q",
		"name":                  fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"speed":                 50,
		"speed_unit":            "MB",
		"notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"purchase_order_number": randString(10),
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
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_status"),
				),
			},
		},
	})
}

func TestAccECXL2ConnectionAzureDot1QPub(t *testing.T) {
	t.Parallel()
	priPortName, _ := schema.EnvDefaultFunc(priPortEnvVar, "sit-001-CX-SV1-NL-Dot1q-BO-10G-PRI-JUN-33")()
	secPortName, _ := schema.EnvDefaultFunc(secPortEnvVar, "sit-001-CX-SV5-NL-Dot1q-BO-10G-SEC-JUN-36")()
	spName, _ := schema.EnvDefaultFunc(azureSpEnvVar, "Azure Express Route")()
	context := map[string]interface{}{
		"pri_port_name":         priPortName.(string),
		"sec_port_name":         secPortName.(string),
		"profile_name":          spName.(string),
		"resourceName":          "tf-azure-dot1q-pub",
		"name":                  fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"speed":                 50,
		"speed_unit":            "MB",
		"notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"purchase_order_number": randString(10),
		"vlan_stag":             randInt(2000),
		"seller_metro_code":     "SV",
		"authorization_key":     "123456789069",
		"named_tag":             "Public",
		"secondary_name":        fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
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
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_status"),
					resource.TestCheckResourceAttrSet(resourceName, "redundancy_type"),
					resource.TestCheckResourceAttrSet(resourceName, "redundant_uuid"),
				),
			},
		},
	})
}

func testAccECXL2ConnectionAWSDot1Q(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_ecx_l2_sellerprofile" "profile" {
  name = "%{profile_name}"
}

data "equinix_ecx_port" "port" {
  name = "%{port_name}"
}

resource "equinix_ecx_l2_connection" "%{resourceName}" {
  name = "%{name}"
  profile_uuid = data.equinix_ecx_l2_sellerprofile.profile.uuid
  speed = %{speed}
  speed_unit = "%{speed_unit}"
  notifications = %{notifications}
  purchase_order_number = "%{purchase_order_number}"
  port_uuid = data.equinix_ecx_port.port.uuid
  vlan_stag = %{vlan_stag}
  seller_region = "%{seller_region}"
  seller_metro_code = "%{seller_metro_code}"
  authorization_key = "%{authorization_key}"
}
`, ctx)
}

func testAccECXL2ConnectionAzureDot1QPub(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_ecx_l2_sellerprofile" "profile" {
  name = "%{profile_name}"
}
	  
data "equinix_ecx_port" "port-pri" {
  name = "%{pri_port_name}"
}

data "equinix_ecx_port" "port-sec" {
  name = "%{sec_port_name}"
}

resource "equinix_ecx_l2_connection" "%{resourceName}" {
  name = "%{name}"
  profile_uuid = data.equinix_ecx_l2_sellerprofile.profile.uuid
  speed = %{speed}
  speed_unit = "%{speed_unit}"
  notifications = %{notifications}
  purchase_order_number = "%{purchase_order_number}"
  port_uuid = data.equinix_ecx_port.port-pri.uuid
  vlan_stag = %{vlan_stag}
  seller_metro_code = "%{seller_metro_code}"
  authorization_key = "%{authorization_key}"
  named_tag = "%{named_tag}"
  secondary_connection {
    name = "%{secondary_name}"
    port_uuid = data.equinix_ecx_port.port-sec.uuid
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
