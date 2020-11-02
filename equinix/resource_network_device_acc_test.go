package equinix

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	networkDeviceMetroEnvVar = "TF_ACC_NETWORK_DEVICE_METRO"
)

func init() {
	resource.AddTestSweepers("NetworkDevice", &resource.Sweeper{
		Name: "NetworkDevice",
		F:    testSweepNetworkDevice,
	})
}

func testSweepNetworkDevice(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	devices, err := config.ne.GetDevices([]string{
		ne.DeviceStateInitializing,
		ne.DeviceStateProvisioned,
		ne.DeviceStateProvisioning,
		ne.DeviceStateWaitingSecondary,
		ne.DeviceStateFailed})
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching NetworkDevice list: %s", err)
		return err
	}
	for _, device := range devices {
		if !isSweepableTestResource(device.Name) {
			continue
		}
		if device.RedundancyType != "PRIMARY" {
			continue
		}
		if err := config.ne.DeleteDevice(device.UUID); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkDevice resource %s (%s): %s", device.UUID, device.Name, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkDevice resource %s (%s)", device.UUID, device.Name)
		}
	}
	return nil
}

func TestAccNetworkDeviceAndUser(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"resourceName":            "tst-csr1000v",
		"name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"throughput":              500,
		"throughput_unit":         "Mbps",
		"metro_code":              metro.(string),
		"type_code":               "CSR1000V",
		"package_code":            "SEC",
		"notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"hostname":                fmt.Sprintf("tf-%s", randString(6)),
		"term_length":             1,
		"version":                 "16.09.05",
		"core_count":              2,
		"purchase_order_number":   randString(10),
		"order_reference":         randString(10),
		"interface_count":         24,
		"secondary-name":          fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"secondary-hostname":      fmt.Sprintf("tf-%s", randString(6)),
		"secondary-notifications": []string{"secondary@equinix.com"},
		"userResourceName":        "tst-user",
		"username":                fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"password":                randString(10),
		"acl-resourceName":        "acl-pri",
		"acl-name":                fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"acl-description":         randString(50),
		"acl2-resourceName":       "acl-sec",
		"acl2-name":               fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"acl2-description":        randString(50),
	}
	resourceName := fmt.Sprintf("equinix_network_device.%s", context["resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", context["userResourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl2-resourceName"].(string))
	var primary, secondary ne.Device
	var user ne.SSHUser
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkDeviceAndUser(context),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(resourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "status", ne.DeviceStateProvisioned),
					resource.TestCheckResourceAttr(resourceName, "license_status", ne.DeviceLicenseStateRegistered),
					resource.TestCheckResourceAttrSet(resourceName, "ibx"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_ip_fqdn"),
					testAccNeSSHUserExists(userResourceName, &user),
					testAccNeSSHUserAttributes(&user, []*ne.Device{&primary, &secondary}, context),
					resource.TestCheckResourceAttrSet(userResourceName, "uuid"),
				),
			},
			{
				Config: testAccNetworkDeviceAndUserAddACLs(context),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					testAccNeDeviceACLs(&primary, &secondary, &primaryACL, &secondaryACL),
				),
			},
		},
	})
}

func testAccNetworkDeviceAndUser(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_account" "test" {
  metro_code = "%{metro_code}"
  status     = "Active"
}

resource "equinix_network_device" "%{resourceName}" {
	name                  = "%{name}"
	throughput            = %{throughput}
	throughput_unit       = "%{throughput_unit}"
	metro_code            = data.equinix_network_account.test.metro_code
	type_code             = "%{type_code}"
	package_code          = "%{package_code}"
	notifications         = %{notifications}
	hostname              = "%{hostname}"
	term_length           = %{term_length}
	account_number        = data.equinix_network_account.test.number
	version               = "%{version}"
	core_count            = %{core_count}
	purchase_order_number = "%{purchase_order_number}"
	order_reference       = "%{order_reference}"
	interface_count       = %{interface_count}
	secondary_device {
		name           = "%{secondary-name}"
		metro_code     = data.equinix_network_account.test.metro_code
		hostname       = "%{secondary-hostname}"
		notifications  = %{secondary-notifications}
		account_number = data.equinix_network_account.test.number
	  }
}

resource "equinix_network_ssh_user" "%{userResourceName}" {
	username = "%{username}"
	password = "%{password}"
	device_ids = [
	  equinix_network_device.%{resourceName}.id,
	  equinix_network_device.%{resourceName}.redundant_id
	]
}
`, ctx)
}

func testAccNetworkDeviceAndUserAddACLs(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_account" "test" {
  metro_code = "%{metro_code}"
  status     = "Active"
}

resource "equinix_network_acl_template" "%{acl-resourceName}" {
	name          = "%{acl-name}"
	description   = "%{acl-description}"
	metro_code    = data.equinix_network_account.test.metro_code
	inbound_rule {
		subnets  = ["10.0.0.0/24"]
		protocol = "IP"
		src_port = "any"
		dst_port = "any"
	}
}

resource "equinix_network_acl_template" "%{acl2-resourceName}" {
	name          = "%{acl2-name}"
	description   = "%{acl2-description}"
	metro_code    = data.equinix_network_account.test.metro_code
	inbound_rule {
		subnets  = ["192.0.0.0/24"]
		protocol = "IP"
		src_port = "any"
		dst_port = "any"
	}
}

resource "equinix_network_device" "%{resourceName}" {
	name                  = "%{name}"
	throughput            = %{throughput}
	throughput_unit       = "%{throughput_unit}"
	metro_code            = data.equinix_network_account.test.metro_code
	type_code             = "%{type_code}"
	package_code          = "%{package_code}"
	notifications         = %{notifications}
	hostname              = "%{hostname}"
	term_length           = %{term_length}
	account_number        = data.equinix_network_account.test.number
	version               = "%{version}"
	core_count            = %{core_count}
	purchase_order_number = "%{purchase_order_number}"
	order_reference       = "%{order_reference}"
	interface_count       = %{interface_count}
	acl_template_id       = equinix_network_acl_template.%{acl-resourceName}.id
	secondary_device {
		name            = "%{secondary-name}"
		metro_code      = data.equinix_network_account.test.metro_code
		hostname        = "%{secondary-hostname}"
		notifications   = %{secondary-notifications}
		account_number  = data.equinix_network_account.test.number
		acl_template_id = equinix_network_acl_template.%{acl2-resourceName}.id
	  }
}

resource "equinix_network_ssh_user" "%{userResourceName}" {
	username = "%{username}"
	password = "%{password}"
	device_ids = [
	  equinix_network_device.%{resourceName}.id,
	  equinix_network_device.%{resourceName}.redundant_id
	]
}
`, ctx)
}

func testAccNeDeviceExists(resourceName string, device *ne.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		client := testAccProvider.Meta().(*Config).ne
		resp, err := client.GetDevice(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching network device '%s': %s", rs.Primary.ID, err)
		}
		*device = *resp
		return nil
	}
}

func testAccNeDeviceSecondaryExists(primary, secondary *ne.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if primary.RedundantUUID == "" {
			return fmt.Errorf("secondary device UUID is not set")
		}
		client := testAccProvider.Meta().(*Config).ne
		resp, err := client.GetDevice(primary.RedundantUUID)
		if err != nil {
			return fmt.Errorf("error when fetching network device '%s': %s", primary.RedundantUUID, err)
		}
		*secondary = *resp
		return nil
	}
}

func testAccNeDeviceAttributes(device *ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["name"]; ok && device.Name != v.(string) {
			return fmt.Errorf("name does not match %v - %v", device.Name, v)
		}
		if v, ok := ctx["throughput"]; ok && device.Throughput != v.(int) {
			return fmt.Errorf("throughput does not match %v - %v", device.Throughput, v)
		}
		if v, ok := ctx["throughput_unit"]; ok && device.ThroughputUnit != v.(string) {
			return fmt.Errorf("throughput_unit does not match %v - %v", device.ThroughputUnit, v)
		}
		if v, ok := ctx["metro_code"]; ok && device.MetroCode != v.(string) {
			return fmt.Errorf("metro_code does not match %v - %v", device.MetroCode, v)
		}
		if v, ok := ctx["type_code"]; ok && device.TypeCode != v.(string) {
			return fmt.Errorf("metro_code does not match %v - %v", device.TypeCode, v)
		}
		if v, ok := ctx["package_code"]; ok && device.PackageCode != v.(string) {
			return fmt.Errorf("package_code does not match %v - %v", device.PackageCode, v)
		}
		if v, ok := ctx["notifications"]; ok && !slicesMatch(device.Notifications, v.([]string)) {
			return fmt.Errorf("notifications does not match %v - %v", device.Notifications, v)
		}
		if v, ok := ctx["hostname"]; ok && device.HostName != v.(string) {
			return fmt.Errorf("hostname does not match %v - %v", device.HostName, v)
		}
		if v, ok := ctx["term_length"]; ok && device.TermLength != v.(int) {
			return fmt.Errorf("term_length does not match %v - %v", device.TermLength, v)
		}
		if v, ok := ctx["version"]; ok && device.Version != v.(string) {
			return fmt.Errorf("version does not match %v - %v", device.Version, v)
		}
		if v, ok := ctx["core_count"]; ok && device.CoreCount != v.(int) {
			return fmt.Errorf("version does not match %v - %v", device.CoreCount, v)
		}
		if v, ok := ctx["purchase_order_number"]; ok && device.PurchaseOrderNumber != v.(string) {
			return fmt.Errorf("purchase_order_number does not match %v - %v", device.PurchaseOrderNumber, v)
		}
		if v, ok := ctx["order_reference"]; ok && device.OrderReference != v.(string) {
			return fmt.Errorf("order_reference does not match %v - %v", device.OrderReference, v)
		}
		if v, ok := ctx["interface_count"]; ok && device.InterfaceCount != v.(int) {
			return fmt.Errorf("interface_count does not match %v - %v", device.InterfaceCount, v)
		}
		if device.IsBYOL != false {
			return fmt.Errorf("byol does not match  %v - %v", device.IsBYOL, false)
		}
		if device.IsSelfManaged != false {
			return fmt.Errorf("self_managed does not match  %v - %v", device.IsSelfManaged, false)
		}
		return nil
	}
}

func testAccNeDeviceSecondaryAttributes(device *ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	secCtx := make(map[string]interface{})
	for key, value := range ctx {
		secCtx[key] = value
	}
	secCtx["name"] = ctx["secondary-name"]
	secCtx["hostname"] = ctx["secondary-hostname"]
	secCtx["notifications"] = ctx["secondary-notifications"]
	return testAccNeDeviceAttributes(device, secCtx)
}

func testAccNeDeviceRedundancyAttributes(primary, secondary *ne.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if primary.RedundancyType != "PRIMARY" {
			return fmt.Errorf("redundancy_type does not match %v - %v", primary.RedundancyType, "PRIMARY")
		}
		if primary.RedundantUUID != secondary.UUID {
			return fmt.Errorf("redundant_id does not match %v - %v", primary.RedundantUUID, secondary.UUID)
		}
		if secondary.RedundancyType != "SECONDARY" {
			return fmt.Errorf("redundancy_type does not match %v - %v", secondary.RedundancyType, "SECONDARY")
		}
		if secondary.RedundantUUID != primary.UUID {
			return fmt.Errorf("redundant_id does not match %v - %v", secondary.RedundantUUID, primary.UUID)
		}
		return nil
	}
}

func testAccNeDeviceACLs(primary, secondary *ne.Device, primaryACL, secondaryACL *ne.ACLTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if primary.ACLTemplateUUID != primaryACL.UUID {
			return fmt.Errorf("Primary device %s template UUID does not match %v - %v", primary.UUID, primary.ACLTemplateUUID, primaryACL.UUID)
		}
		if secondary.ACLTemplateUUID != secondary.UUID {
			return fmt.Errorf("Secondary device %s template UUID does not match %v - %v", secondary.UUID, secondary.ACLTemplateUUID, secondaryACL.UUID)
		}
		if primaryACL.DeviceUUID != primary.UUID {
			return fmt.Errorf("Primary ACL %s device UUID does not match %v - %v", primaryACL.UUID, primaryACL.DeviceUUID, primary.UUID)
		}
		if secondaryACL.DeviceUUID != secondary.UUID {
			return fmt.Errorf("Secondary ACL %s device UUID does not match %v - %v", secondaryACL.UUID, secondaryACL.DeviceUUID, secondary.UUID)
		}
		return nil
	}
}
