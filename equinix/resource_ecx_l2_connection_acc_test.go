package equinix

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/ecx-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	priPortEnvVar              = "TF_ACC_FABRIC_PRI_PORT_NAME"
	secPortEnvVar              = "TF_ACC_FABRIC_SEC_PORT_NAME"
	awsSpEnvVar                = "TF_ACC_FABRIC_L2_AWS_SP_NAME"
	awsAuthKeyEnvVar           = "TF_ACC_FABRIC_L2_AWS_ACCOUNT_ID"
	azureSpEnvVar              = "TF_ACC_FABRIC_L2_AZURE_SP_NAME"
	azureXRServiceKeyEnvVar    = "TF_ACC_FABRIC_L2_AZURE_XROUTE_SERVICE_KEY"
	gcpOneSpEnvVar             = "TF_ACC_FABRIC_L2_GCP1_SP_NAME"
	gcpOneConnServiceKeyEnvVar = "TF_ACC_FABRIC_L2_GCP1_INTERCONN_SERVICE_KEY"
	gcpTwoSpEnvVar             = "TF_ACC_FABRIC_L2_GCP2_SP_NAME"
	gcpTwoConnServiceKeyEnvVar = "TF_ACC_FABRIC_L2_GCP2_INTERCONN_SERVICE_KEY"
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
		if !isSweepableTestResource(ecx.StringValue(conn.Name)) {
			nonSweepableCount++
			continue
		}
		if err := config.ecx.DeleteL2Connection(ecx.StringValue(conn.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting ECXL2Connection resource %s (%s): %s", ecx.StringValue(conn.UUID), ecx.StringValue(conn.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for ECXL2Connection resource %s (%s)", ecx.StringValue(conn.UUID), ecx.StringValue(conn.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}

func TestAccFabricL2Connection_Port_Single_AWS(t *testing.T) {
	portName, _ := schema.EnvDefaultFunc(priPortEnvVar, "sit-001-CX-SV1-NL-Dot1q-BO-10G-PRI-JUN-33")()
	spName, _ := schema.EnvDefaultFunc(awsSpEnvVar, "AWS Direct Connect")()
	authKey, _ := schema.EnvDefaultFunc(awsAuthKeyEnvVar, "123456789012")()
	context := map[string]interface{}{
		"port-resourceName":                "test",
		"port-name":                        portName.(string),
		"connection-resourceName":          "test",
		"connection-profile_name":          spName.(string),
		"connection-name":                  fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"connection-speed":                 50,
		"connection-speed_unit":            "MB",
		"connection-notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"connection-purchase_order_number": randString(10),
		"connection-vlan_stag":             randInt(2000),
		"connection-seller_region":         "us-west-1",
		"connection-seller_metro_code":     "SV",
		"connection-authorization_key":     authKey.(string),
	}
	resourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["connection-resourceName"].(string))
	var testConn ecx.L2Connection
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withPort().withConnection().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccFabricL2ConnectionExists(resourceName, &testConn),
					testAccFabricL2ConnectionAttributes(&testConn, context),
					resource.TestCheckResourceAttr(resourceName, "status", ecx.ConnectionStatusProvisioned),
					resource.TestCheckResourceAttrSet(resourceName, "provider_status"),
					resource.TestCheckResourceAttrSet(resourceName, "zside_port_uuid"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.required_data.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.required_data.0.key", "awsConnectionId"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFabricL2Connection_Port_HA_Azure(t *testing.T) {
	priPortName, _ := schema.EnvDefaultFunc(priPortEnvVar, "sit-001-CX-SV1-NL-Dot1q-BO-10G-PRI-JUN-33")()
	secPortName, _ := schema.EnvDefaultFunc(secPortEnvVar, "sit-001-CX-SV5-NL-Dot1q-BO-10G-SEC-JUN-36")()
	spName, _ := schema.EnvDefaultFunc(azureSpEnvVar, "Azure ExpressRoute")()
	serviceKey, _ := schema.EnvDefaultFunc(azureXRServiceKeyEnvVar, "ExpressRoute-ServiceKey")()
	context := map[string]interface{}{
		"port-resourceName":                "test",
		"port-name":                        priPortName.(string),
		"port-secondary_resourceName":      "test-sec",
		"port-secondary_name":              secPortName.(string),
		"connection-resourceName":          "test",
		"connection-profile_name":          spName.(string),
		"connection-name":                  fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"connection-speed":                 50,
		"connection-speed_unit":            "MB",
		"connection-notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"connection-purchase_order_number": randString(10),
		"connection-vlan_stag":             randInt(2000),
		"connection-seller_metro_code":     "SV",
		"connection-authorization_key":     serviceKey,
		"connection-named_tag":             "Public",
		"connection-secondary_name":        fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"connection-secondary_vlan_stag":   randInt(2000),
	}
	contextWithChanges := copyMap(context)
	contextWithChanges["connection-name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	contextWithChanges["connection-secondary_name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	resourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["connection-resourceName"].(string))
	var primary, secondary ecx.L2Connection
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withPort().withConnection().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccFabricL2ConnectionExists(resourceName, &primary),
					testAccFabricL2ConnectionAttributes(&primary, context),
					testAccFabricL2ConnectionSecondaryExists(&primary, &secondary),
					testAccFabricL2ConnectionSecondaryAttributes(&secondary, context),
					resource.TestCheckResourceAttr(resourceName, "status", ecx.ConnectionStatusPendingBGPPeering),
					resource.TestCheckResourceAttrSet(resourceName, "provider_status"),
					testAccFabricL2ConnectionRedundancyAttributes(&primary, &secondary),
				),
			},
			{
				Config: newTestAccConfig(contextWithChanges).withPort().withConnection().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccFabricL2ConnectionExists(resourceName, &primary),
					testAccFabricL2ConnectionAttributes(&primary, contextWithChanges),
					testAccFabricL2ConnectionSecondaryExists(&primary, &secondary),
					testAccFabricL2ConnectionSecondaryAttributes(&secondary, contextWithChanges),
					resource.TestCheckResourceAttr(resourceName, "status", ecx.ConnectionStatusPendingBGPPeering),
					resource.TestCheckResourceAttrSet(resourceName, "provider_status"),
					testAccFabricL2ConnectionRedundancyAttributes(&primary, &secondary),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFabricL2Connection_Device_HA_GCP(t *testing.T) {
	deviceMetro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	priSPName, _ := schema.EnvDefaultFunc(gcpOneSpEnvVar, "Google Cloud Partner Interconnect Zone 1")()
	secSPName, _ := schema.EnvDefaultFunc(gcpTwoSpEnvVar, "Google Cloud Partner Interconnect Zone 2")()
	priServiceKey, _ := schema.EnvDefaultFunc(gcpOneConnServiceKeyEnvVar, "Interconnect-ServiceKey")()
	secServiceKey, _ := schema.EnvDefaultFunc(gcpTwoConnServiceKeyEnvVar, "Interconnect-ServiceKey")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":                      "test",
		"device-account_name":                      accountName.(string),
		"device-self_managed":                      true,
		"device-byol":                              true,
		"device-name":                              fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-metro_code":                        deviceMetro.(string),
		"device-type_code":                         "PA-VM",
		"device-package_code":                      "VM100",
		"device-notifications":                     []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                          fmt.Sprintf("tf-%s", randString(6)),
		"device-term_length":                       1,
		"device-version":                           "9.0.4",
		"device-core_count":                        2,
		"device-purchase_order_number":             randString(10),
		"device-order_reference":                   randString(10),
		"device-secondary_name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-secondary_hostname":                fmt.Sprintf("tf-%s", randString(6)),
		"device-secondary_notifications":           []string{"secondary@equinix.com"},
		"sshkey-resourceName":                      "test",
		"sshkey-name":                              fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"sshkey-public_key":                        "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXdzXBHaVpKpdO0udnB+4JOgUq7APO2rPXfrevvlZrps98AtlwXXVWZ5duRH5NFNfU4G9HCSiAPsebgjY0fG85tcShpXfHfACLt0tBW8XhfLQP2T6S50FQ1brBdURMDCMsD7duOXqvc0dlbs2/KcswHvuUmqVzob3bz7n1bQ48wIHsPg4ARqYhy5LN3OkllJH/6GEfqi8lKZx01/P/gmJMORcJujuOyXRB+F2iXBVYdhjML3Qg4+tEekBcVZOxUbERRZ0pvQ52Y6wUhn2VsjljixyqeOdmD0m6DayDQgSWms6bKPpBqN7zhXXk4qe8bXT4tQQba65b2CQ2A91jw2KgM/YZNmjyUJ+Rf1cQosJf9twqbAZDZ6rAEmj9zzvQ5vD/CGuzxdVMkePLlUK4VGjPu7cVzhXrnq4318WqZ5/lNiCST8NQ0fssChN8ANUzr/p/wwv3faFMVNmjxXTZMsbMFT/fbb2MVVuqNFN65drntlg6/xEao8gZROuRYiakBx8= user@host",
		"connection-resourceName":                  "test",
		"connection-name":                          fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"connection-profile_name":                  priSPName.(string),
		"connection-speed":                         50,
		"connection-speed_unit":                    "MB",
		"connection-notifications":                 []string{"marry@equinix.com", "john@equinix.com"},
		"connection-purchase_order_number":         randString(10),
		"connection-seller_metro_code":             "SV",
		"connection-seller_region":                 "us-west2",
		"connection-authorization_key":             priServiceKey.(string),
		"connection-device_interface_id":           5,
		"connection-secondary_name":                fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"connection-secondary_profile_name":        secSPName.(string),
		"connection-secondary_speed":               100,
		"connection-secondary_speed_unit":          "MB",
		"connection-secondary_seller_metro_code":   "SV",
		"connection-secondary_seller_region":       "us-west2",
		"connection-secondary_authorization_key":   secServiceKey.(string),
		"connection-secondary_device_interface_id": 5,
	}
	connResourceName := fmt.Sprintf("equinix_ecx_l2_connection.%s", context["connection-resourceName"].(string))
	var primary, secondary ecx.L2Connection
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withSSHKey().
					withConnection().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccFabricL2ConnectionExists(connResourceName, &primary),
					testAccFabricL2ConnectionAttributes(&primary, context),
					testAccFabricL2ConnectionSecondaryExists(&primary, &secondary),
					testAccFabricL2ConnectionSecondaryAttributes(&secondary, context),
				),
			},
			{
				ResourceName:            connResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_interface_id", "secondary_connection.0.device_interface_id"},
			},
		},
	})
}

func testAccFabricL2ConnectionExists(resourceName string, conn *ecx.L2Connection) resource.TestCheckFunc {
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
		if ecx.StringValue(resp.UUID) != rs.Primary.ID {
			return fmt.Errorf("resource ID does not match %v - %v", rs.Primary.ID, resp.UUID)
		}
		*conn = *resp
		return nil
	}
}

func testAccFabricL2ConnectionSecondaryExists(primary *ecx.L2Connection, secondary *ecx.L2Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Config).ecx
		if ecx.StringValue(primary.RedundantUUID) == "" {
			return fmt.Errorf("primary connection has no RedundantUUID set")
		}
		resp, err := client.GetL2Connection(ecx.StringValue(primary.RedundantUUID))
		if err != nil {
			return fmt.Errorf("error when fetching L2 connection %v", err)
		}
		*secondary = *resp
		return nil
	}
}

func testAccFabricL2ConnectionAttributes(conn *ecx.L2Connection, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["connection-name"]; ok && ecx.StringValue(conn.Name) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ecx.StringValue(conn.Name), v)
		}
		if v, ok := ctx["connection-speed"]; ok && ecx.IntValue(conn.Speed) != v.(int) {
			return fmt.Errorf("speed does not match %v - %v", ecx.IntValue(conn.Speed), v)
		}
		if v, ok := ctx["connection-speed_unit"]; ok && ecx.StringValue(conn.SpeedUnit) != v.(string) {
			return fmt.Errorf("speedUnit does not match %v - %v", ecx.StringValue(conn.SpeedUnit), v)
		}
		if v, ok := ctx["connection-notifications"]; ok && !slicesMatch(conn.Notifications, v.([]string)) {
			return fmt.Errorf("notifications does not match %v - %v", conn.Notifications, v)
		}
		if v, ok := ctx["connection-purchase_order_number"]; ok && ecx.StringValue(conn.PurchaseOrderNumber) != v.(string) {
			return fmt.Errorf("purchaseOrderNumber does not match %v - %v", ecx.StringValue(conn.PurchaseOrderNumber), v)
		}
		if v, ok := ctx["connection-vlan_stag"]; ok && ecx.IntValue(conn.VlanSTag) != v.(int) {
			return fmt.Errorf("vlanSTag does not match %v - %v", ecx.IntValue(conn.VlanSTag), v)
		}
		if v, ok := ctx["connection-vlan_ctag"]; ok && ecx.IntValue(conn.VlanCTag) != v.(int) {
			return fmt.Errorf("vlanCTag does not match %v - %v", ecx.IntValue(conn.VlanCTag), v)
		}
		if v, ok := ctx["connection-zside_port_uuid"]; ok && ecx.StringValue(conn.ZSidePortUUID) != v.(string) {
			return fmt.Errorf("zSidePortUUID does not match %v - %v", ecx.StringValue(conn.ZSidePortUUID), v)
		}
		if v, ok := ctx["connection-zside_vlan_stag"]; ok && ecx.IntValue(conn.ZSideVlanSTag) != v.(int) {
			return fmt.Errorf("zSideVlanSTag does not match %v - %v", ecx.IntValue(conn.ZSideVlanSTag), v)
		}
		if v, ok := ctx["connection-zside_vlan_ctag"]; ok && ecx.IntValue(conn.ZSideVlanCTag) != v.(int) {
			return fmt.Errorf("zSideVlanCTag does not match %v - %v", ecx.IntValue(conn.ZSideVlanCTag), v)
		}
		if v, ok := ctx["connection-named_tag"]; ok && ecx.StringValue(conn.NamedTag) != v.(string) {
			return fmt.Errorf("named_tag does not match %v - %v", ecx.StringValue(conn.NamedTag), v)
		}
		if v, ok := ctx["connection-seller_region"]; ok && ecx.StringValue(conn.SellerRegion) != v.(string) {
			return fmt.Errorf("sellerRegion does not match %v - %v", ecx.StringValue(conn.SellerRegion), v)
		}
		if v, ok := ctx["connection-seller_metro_code"]; ok && ecx.StringValue(conn.SellerMetroCode) != v.(string) {
			return fmt.Errorf("sellerMetroCode does not match %v - %v", ecx.StringValue(conn.SellerMetroCode), v)
		}
		if v, ok := ctx["connection-authorization_key"]; ok && ecx.StringValue(conn.AuthorizationKey) != v.(string) {
			return fmt.Errorf("authorizationKey does not match %v - %v", ecx.StringValue(conn.AuthorizationKey), v)
		}
		return nil
	}
}

func testAccFabricL2ConnectionSecondaryAttributes(conn *ecx.L2Connection, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["connection-secondary_name"]; ok && ecx.StringValue(conn.Name) != v.(string) {
			return fmt.Errorf("connection secondary name does not match %v - %v", ecx.StringValue(conn.Name), v)
		}
		if v, ok := ctx["connection-secondary_speed"]; ok && ecx.IntValue(conn.Speed) != v.(int) {
			return fmt.Errorf("connection secondary speed does not match %v - %v", ecx.IntValue(conn.Speed), v)
		}
		if v, ok := ctx["connection-secondary_speed_unit"]; ok && ecx.StringValue(conn.SpeedUnit) != v.(string) {
			return fmt.Errorf("connection secondary speed unit does not match %v - %v", ecx.StringValue(conn.SpeedUnit), v)
		}
		if v, ok := ctx["connection-secondary_vlan_stag"]; ok && ecx.IntValue(conn.VlanSTag) != v.(int) {
			return fmt.Errorf("connection secondary vlanSTag does not match %v - %v", ecx.IntValue(conn.VlanSTag), v)
		}
		if v, ok := ctx["connection-secondary_vlan_ctag"]; ok && ecx.IntValue(conn.VlanCTag) != v.(int) {
			return fmt.Errorf("connection secondary vlanCTag does not match %v - %v", ecx.IntValue(conn.VlanCTag), v)
		}
		if v, ok := ctx["connection-secondary-seller_region"]; ok && ecx.StringValue(conn.SellerRegion) != v.(string) {
			return fmt.Errorf("connection secondary seller region does not match %v - %v", ecx.StringValue(conn.SellerRegion), v)
		}
		if v, ok := ctx["connection-secondary-seller_metro_code"]; ok && ecx.StringValue(conn.SellerMetroCode) != v.(string) {
			return fmt.Errorf("connection secondary seller metro code does not match %v - %v", ecx.StringValue(conn.SellerMetroCode), v)
		}
		if v, ok := ctx["connection-secondary_authorization_key"]; ok && ecx.StringValue(conn.AuthorizationKey) != v.(string) {
			return fmt.Errorf("connection secondary authorization_key code does not match %v - %v", ecx.StringValue(conn.AuthorizationKey), v)
		}
		return nil
	}
}

func testAccFabricL2ConnectionRedundancyAttributes(primary, secondary *ecx.L2Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ecx.StringValue(primary.RedundancyType) != "primary" {
			return fmt.Errorf("primary connection redundancy type does not match  %v - %v", ecx.StringValue(primary.RedundancyType), "primary")
		}
		if ecx.StringValue(primary.RedundantUUID) != ecx.StringValue(secondary.UUID) {
			return fmt.Errorf("primary connection redundant UUID does not match  %v - %v", ecx.StringValue(primary.RedundantUUID), ecx.StringValue(secondary.UUID))
		}
		if ecx.StringValue(secondary.RedundancyType) != "secondary" {
			return fmt.Errorf("secondary connection redundancy type does not match  %v - %v", ecx.StringValue(secondary.RedundancyType), "secondary")
		}
		if ecx.StringValue(secondary.RedundantUUID) != ecx.StringValue(primary.UUID) {
			return fmt.Errorf("secondary connection redundant UUID does not match  %v - %v", ecx.StringValue(secondary.RedundantUUID), ecx.StringValue(primary.UUID))
		}
		return nil
	}
}

func (t *testAccConfig) withConnection() *testAccConfig {
	t.config += testAccFabricL2Connection(t.ctx)
	return t
}

func (t *testAccConfig) withPort() *testAccConfig {
	t.config += testAccFabricPort(t.ctx)
	return t
}

func testAccFabricPort(ctx map[string]interface{}) string {
	var config string
	config += nprintf(`
data "equinix_ecx_port" "%{port-resourceName}" {
  name = "%{port-name}"
}`, ctx)

	if _, ok := ctx["port-secondary_resourceName"]; ok {
		config += nprintf(`
data "equinix_ecx_port" "%{port-secondary_resourceName}" {
  name = "%{port-secondary_name}"
}`, ctx)
	}
	return config
}

func testAccFabricL2Connection(ctx map[string]interface{}) string {
	var config string
	config += nprintf(`
data "equinix_ecx_l2_sellerprofile" "pri" {
  name = "%{connection-profile_name}"
}`, ctx)

	if _, ok := ctx["connection-secondary_profile_name"]; ok {
		config += nprintf(`
data "equinix_ecx_l2_sellerprofile" "sec" {
  name = "%{connection-secondary_profile_name}"
}`, ctx)
	}

	config += nprintf(`
resource "equinix_ecx_l2_connection" "%{connection-resourceName}" {
  name                  = "%{connection-name}"
  profile_uuid          = data.equinix_ecx_l2_sellerprofile.pri.id
  speed                 = %{connection-speed}
  speed_unit            = "%{connection-speed_unit}"
  notifications         = %{connection-notifications}
  seller_metro_code     = "%{connection-seller_metro_code}"
  authorization_key     = "%{connection-authorization_key}"`, ctx)
	if _, ok := ctx["connection-purchase_order_number"]; ok {
		config += nprintf(`
  purchase_order_number = "%{connection-purchase_order_number}"`, ctx)
	}
	if _, ok := ctx["connection-seller_region"]; ok {
		config += nprintf(`
  seller_region         = "%{connection-seller_region}"`, ctx)
	}
	if _, ok := ctx["port-resourceName"]; ok {
		config += nprintf(`
  port_uuid             = data.equinix_ecx_port.%{port-resourceName}.id`, ctx)
	}
	if _, ok := ctx["device-resourceName"]; ok {
		config += nprintf(`
  device_uuid           = equinix_network_device.%{device-resourceName}.id`, ctx)
	}
	if _, ok := ctx["connection-vlan_stag"]; ok {
		config += nprintf(`
  vlan_stag             = %{connection-vlan_stag}`, ctx)
	}
	if _, ok := ctx["connection-vlan_ctag"]; ok {
		config += nprintf(`
  vlan_ctag             = %{connection-vlan_ctag}`, ctx)
	}
	if _, ok := ctx["connection-named_tag"]; ok {
		config += nprintf(`
  named_tag             = "%{connection-named_tag}"`, ctx)
	}
	if _, ok := ctx["connection-device_interface_id"]; ok {
		config += nprintf(`
  device_interface_id   = %{connection-device_interface_id}`, ctx)
	}
	if _, ok := ctx["connection-secondary_name"]; ok {
		config += nprintf(`
  secondary_connection {
    name                = "%{connection-secondary_name}"`, ctx)
		if _, ok := ctx["connection-secondary_profile_name"]; ok {
			config += nprintf(`
    profile_uuid        = data.equinix_ecx_l2_sellerprofile.sec.id`, ctx)
		}
		if _, ok := ctx["port-secondary_resourceName"]; ok {
			config += nprintf(`
    port_uuid           = data.equinix_ecx_port.%{port-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["device-secondary_name"]; ok {
			config += nprintf(`
    device_uuid         = equinix_network_device.%{device-resourceName}.redundant_id`, ctx)
		}
		if _, ok := ctx["connection-secondary_vlan_stag"]; ok {
			config += nprintf(`
    vlan_stag           = %{connection-secondary_vlan_stag}`, ctx)
		}
		if _, ok := ctx["connection-secondary_vlan_ctag"]; ok {
			config += nprintf(`
    vlan_ctag           = %{connection-secondary_vlan_ctag}`, ctx)
		}
		if _, ok := ctx["connection-secondary_device_interface_id"]; ok {
			config += nprintf(`
    device_interface_id = %{connection-secondary_device_interface_id}`, ctx)
		}
		if _, ok := ctx["connection-secondary_speed"]; ok {
			config += nprintf(`
    speed               = %{connection-secondary_speed}`, ctx)
		}
		if _, ok := ctx["connection-secondary_speed_unit"]; ok {
			config += nprintf(`
    speed_unit          = "%{connection-secondary_speed_unit}"`, ctx)
		}
		if _, ok := ctx["connection-secondary_seller_metro_code"]; ok {
			config += nprintf(`
    seller_metro_code   = "%{connection-secondary_seller_metro_code}"`, ctx)
		}
		if _, ok := ctx["connection-secondary_seller_region"]; ok {
			config += nprintf(`
    seller_region       = "%{connection-secondary_seller_region}"`, ctx)
		}
		if _, ok := ctx["connection-secondary_authorization_key"]; ok {
			config += nprintf(`
    authorization_key   = "%{connection-secondary_authorization_key}"`, ctx)
		}
		config += `
  }`
	}
	config += `
}`
	return config
}
