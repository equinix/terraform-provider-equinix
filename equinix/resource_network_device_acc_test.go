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
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "DC")()
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
		"secondary-hostname":      randString(6),
		"secondary-notifications": []string{"secondary@equinix.com"},
		"userResourceName":        "tst-user",
		"username":                fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"password":                randString(10),
	}
	resourceName := fmt.Sprintf("equinix_network_device.%s", context["resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", context["userResourceName"].(string))
	var primary, secondary ne.Device
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkDeviceAndUser(context),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(resourceName, &primary),
					testAccNeDeviceAttributes(primary, context),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "status", ne.DeviceStateProvisioned),
					resource.TestCheckResourceAttrSet(resourceName, "license_status"),
					resource.TestCheckResourceAttrSet(resourceName, "ibx"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_ip_fqdn"),
					resource.TestCheckResourceAttrSet(resourceName, "redundant_id"),
					resource.TestCheckResourceAttr(resourceName, "redundancy_type", "PRIMARY"),
					testAccNeDeviceSecondaryExists(primary.RedundantUUID, &secondary),
					testAccNeDeviceSecondaryAttributes(secondary, context),
					resource.TestCheckResourceAttrSet(userResourceName, "uuid"),
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
	  equinix_network_device.%{resourceName}.uuid,
	  equinix_network_device.%{resourceName}.redundant_uuid
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

func testAccNeDeviceSecondaryExists(uuid string, device *ne.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if uuid == "" {
			return fmt.Errorf("secondary device UUID is not set")
		}
		client := testAccProvider.Meta().(*Config).ne
		resp, err := client.GetDevice(uuid)
		if err != nil {
			return fmt.Errorf("error when fetching network device '%s': %s", uuid, err)
		}
		*device = *resp
		return nil
	}
}

func testAccNeDeviceAttributes(device ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
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

func testAccNeDeviceSecondaryAttributes(device ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	secCtx := make(map[string]interface{})
	for key, value := range ctx {
		secCtx[key] = value
	}
	secCtx["name"] = ctx["secondary-name"]
	secCtx["hostname"] = ctx["secondary-hostname"]
	secCtx["notifications"] = ctx["secondary-notifications"]
	return testAccNeDeviceAttributes(device, secCtx)
}
