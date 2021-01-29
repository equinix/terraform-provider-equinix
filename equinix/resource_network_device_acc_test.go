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
	networkDeviceMetroEnvVar       = "TF_ACC_NETWORK_DEVICE_METRO"
	networkDeviceLicenseFileEnvVar = "TF_ACC_NETWORK_DEVICE_LICENSE_FILE"
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
	nonSweepableCount := 0
	for _, device := range devices {
		if !isSweepableTestResource(ne.StringValue(device.Name)) {
			nonSweepableCount++
			continue
		}
		if ne.StringValue(device.RedundancyType) != "PRIMARY" {
			continue
		}
		if err := config.ne.DeleteDevice(ne.StringValue(device.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkDevice resource %s (%s): %s", ne.StringValue(device.UUID), ne.StringValue(device.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkDevice resource %s (%s)", ne.StringValue(device.UUID), ne.StringValue(device.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}

func TestAccNetworkDevice_CSR100V_HA_Managed_Sub(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-self_managed":            false,
		"device-byol":                    false,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-throughput":              500,
		"device-throughput_unit":         "Mbps",
		"device-metro_code":              metro.(string),
		"device-type_code":               "CSR1000V",
		"device-package_code":            "SEC",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", randString(6)),
		"device-term_length":             1,
		"device-version":                 "16.09.05",
		"device-core_count":              2,
		"device-purchase_order_number":   randString(10),
		"device-order_reference":         randString(10),
		"device-interface_count":         24,
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", randString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"user-resourceName":              "tst-user",
		"user-username":                  fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"user-password":                  randString(10),
	}
	contextWithACLs := copyMap(context)
	contextWithACLs["acl-resourceName"] = "acl-pri"
	contextWithACLs["acl-name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	contextWithACLs["acl-description"] = randString(50)
	contextWithACLs["acl-metroCode"] = metro.(string)
	contextWithACLs["acl-secondary_resourceName"] = "acl-sec"
	contextWithACLs["acl-secondary_name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	contextWithACLs["acl-secondary_description"] = randString(50)
	contextWithACLs["acl-secondary_metroCode"] = metro.(string)
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", context["user-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-secondary_resourceName"].(string))
	var primary, secondary ne.Device
	var user ne.SSHUser
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withSSHUser().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateRegistered),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateRegistered),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "uuid"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ibx"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "region"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ssh_ip_address"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ssh_ip_fqdn"),
					testAccNeSSHUserExists(userResourceName, &user),
					testAccNeSSHUserAttributes(&user, []*ne.Device{&primary, &secondary}, context),
					resource.TestCheckResourceAttrSet(userResourceName, "uuid"),
				),
			},
			{
				Config: newTestAccConfig(contextWithACLs).withDevice().
					withSSHUser().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceACLs(&primary, &secondary, &primaryACL, &secondaryACL),
				),
			},
		},
	})
}

func TestAccNetworkDevice_vSRX_HA_Managed_BYOL(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	licFile, _ := schema.EnvDefaultFunc(networkDeviceLicenseFileEnvVar, "jnpr.lic")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-self_managed":            false,
		"device-byol":                    true,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-license_file":            licFile.(string),
		"device-metro_code":              metro.(string),
		"device-type_code":               "VSRX",
		"device-package_code":            "STD",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", randString(6)),
		"device-term_length":             1,
		"device-version":                 "19.2R2.7",
		"device-core_count":              2,
		"device-purchase_order_number":   randString(10),
		"device-order_reference":         randString(10),
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-secondary_license_file":  licFile.(string),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", randString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"acl-resourceName":               "acl-pri",
		"acl-name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"acl-description":                randString(50),
		"acl-metroCode":                  metro.(string),
		"acl-secondary_resourceName":     "acl-sec",
		"acl-secondary_name":             fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"acl-secondary_description":      randString(50),
		"acl-secondary_metroCode":        metro.(string),
	}
	contextWithChanges := copyMap(context)
	contextWithChanges["device-name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	contextWithChanges["device-additional_bandwidth"] = 100
	contextWithChanges["device-notifications"] = []string{"jerry@equinix.com", "tom@equinix.com"}
	contextWithChanges["device-secondary_name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	contextWithChanges["device-secondary_additional_bandwidth"] = 100
	contextWithChanges["device-secondary_notifications"] = []string{"miki@equinix.com", "mini@equinix.com"}
	contextWithChanges["user-resourceName"] = "test"
	contextWithChanges["user-username"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	contextWithChanges["user-password"] = randString(10)
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", contextWithChanges["user-resourceName"].(string))
	var primary, secondary ne.Device
	var user ne.SSHUser
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateRegistered),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateRegistered),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "uuid"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ibx"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "region"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ssh_ip_address"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ssh_ip_fqdn"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "license_file_id"),
				),
			},
			{
				Config: newTestAccConfig(contextWithChanges).withDevice().withACL().withSSHUser().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, contextWithChanges),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateRegistered),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, contextWithChanges),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateRegistered),
					testAccNeSSHUserExists(userResourceName, &user),
					testAccNeSSHUserAttributes(&user, []*ne.Device{&primary, &secondary}, contextWithChanges),
				),
			},
		},
	})
}

func TestAccNetworkDevice_PaloAlto_HA_Self_BYOL(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-self_managed":            true,
		"device-byol":                    true,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-metro_code":              metro.(string),
		"device-type_code":               "PA-VM",
		"device-package_code":            "VM100",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", randString(6)),
		"device-term_length":             1,
		"device-version":                 "9.0.4",
		"device-core_count":              2,
		"device-purchase_order_number":   randString(10),
		"device-order_reference":         randString(10),
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", randString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"sshkey-resourceName":            "test",
		"sshkey-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"sshkey-public_key":              "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXdzXBHaVpKpdO0udnB+4JOgUq7APO2rPXfrevvlZrps98AtlwXXVWZ5duRH5NFNfU4G9HCSiAPsebgjY0fG85tcShpXfHfACLt0tBW8XhfLQP2T6S50FQ1brBdURMDCMsD7duOXqvc0dlbs2/KcswHvuUmqVzob3bz7n1bQ48wIHsPg4ARqYhy5LN3OkllJH/6GEfqi8lKZx01/P/gmJMORcJujuOyXRB+F2iXBVYdhjML3Qg4+tEekBcVZOxUbERRZ0pvQ52Y6wUhn2VsjljixyqeOdmD0m6DayDQgSWms6bKPpBqN7zhXXk4qe8bXT4tQQba65b2CQ2A91jw2KgM/YZNmjyUJ+Rf1cQosJf9twqbAZDZ6rAEmj9zzvQ5vD/CGuzxdVMkePLlUK4VGjPu7cVzhXrnq4318WqZ5/lNiCST8NQ0fssChN8ANUzr/p/wwv3faFMVNmjxXTZMsbMFT/fbb2MVVuqNFN65drntlg6/xEao8gZROuRYiakBx8= user@host",
	}
	contextWithACLs := copyMap(context)
	contextWithACLs["acl-resourceName"] = "acl-pri"
	contextWithACLs["acl-name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	contextWithACLs["acl-description"] = randString(50)
	contextWithACLs["acl-metroCode"] = metro.(string)
	contextWithACLs["acl-secondary_resourceName"] = "acl-sec"
	contextWithACLs["acl-secondary_name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6))
	contextWithACLs["acl-secondary_description"] = randString(50)
	contextWithACLs["acl-secondary_metroCode"] = metro.(string)
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-secondary_resourceName"].(string))
	var primary, secondary ne.Device
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withSSHKey().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "uuid"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ibx"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "region"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ssh_ip_address"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "ssh_ip_fqdn"),
				),
			},
			{
				Config: newTestAccConfig(contextWithACLs).withDevice().withSSHKey().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceACLs(&primary, &secondary, &primaryACL, &secondaryACL),
				),
			},
		},
	})
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
		if ne.StringValue(primary.RedundantUUID) == "" {
			return fmt.Errorf("secondary device UUID is not set")
		}
		client := testAccProvider.Meta().(*Config).ne
		resp, err := client.GetDevice(ne.StringValue(primary.RedundantUUID))
		if err != nil {
			return fmt.Errorf("error when fetching network device '%s': %s", ne.StringValue(primary.RedundantUUID), err)
		}
		*secondary = *resp
		return nil
	}
}

func testAccNeDeviceAttributes(device *ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["device-name"]; ok && ne.StringValue(device.Name) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ne.StringValue(device.Name), v)
		}
		if v, ok := ctx["device-self_managed"]; ok && ne.BoolValue(device.IsSelfManaged) != v.(bool) {
			return fmt.Errorf("self_managed does not match %v - %v", ne.BoolValue(device.IsSelfManaged), v)
		}
		if v, ok := ctx["device-byol"]; ok && ne.BoolValue(device.IsBYOL) != v.(bool) {
			return fmt.Errorf("byol does not match %v - %v", ne.BoolValue(device.IsBYOL), v)
		}
		if v, ok := ctx["device-throughput"]; ok && ne.IntValue(device.Throughput) != v.(int) {
			return fmt.Errorf("throughput does not match %v - %v", ne.IntValue(device.Throughput), v)
		}
		if v, ok := ctx["device-throughput_unit"]; ok && ne.StringValue(device.ThroughputUnit) != v.(string) {
			return fmt.Errorf("throughput_unit does not match %v - %v", ne.StringValue(device.ThroughputUnit), v)
		}
		if v, ok := ctx["device-metro_code"]; ok && ne.StringValue(device.MetroCode) != v.(string) {
			return fmt.Errorf("metro_code does not match %v - %v", ne.StringValue(device.MetroCode), v)
		}
		if v, ok := ctx["device-type_code"]; ok && ne.StringValue(device.TypeCode) != v.(string) {
			return fmt.Errorf("type_code does not match %v - %v", ne.StringValue(device.TypeCode), v)
		}
		if v, ok := ctx["device-package_code"]; ok && ne.StringValue(device.PackageCode) != v.(string) {
			return fmt.Errorf("device-package_code does not match %v - %v", ne.StringValue(device.PackageCode), v)
		}
		if v, ok := ctx["device-notifications"]; ok && !slicesMatch(device.Notifications, v.([]string)) {
			return fmt.Errorf("device-notifications does not match %v - %v", device.Notifications, v)
		}
		if v, ok := ctx["device-hostname"]; ok && ne.StringValue(device.HostName) != v.(string) {
			return fmt.Errorf("device-hostname does not match %v - %v", ne.StringValue(device.HostName), v)
		}
		if v, ok := ctx["device-term_length"]; ok && ne.IntValue(device.TermLength) != v.(int) {
			return fmt.Errorf("device-term_length does not match %v - %v", ne.IntValue(device.TermLength), v)
		}
		if v, ok := ctx["device-version"]; ok && ne.StringValue(device.Version) != v.(string) {
			return fmt.Errorf("device-version does not match %v - %v", ne.StringValue(device.Version), v)
		}
		if v, ok := ctx["device-core_count"]; ok && ne.IntValue(device.CoreCount) != v.(int) {
			return fmt.Errorf("device-core_count does not match %v - %v", ne.IntValue(device.CoreCount), v)
		}
		if v, ok := ctx["device-purchase_order_number"]; ok && ne.StringValue(device.PurchaseOrderNumber) != v.(string) {
			return fmt.Errorf("device-purchase_order_number does not match %v - %v", ne.StringValue(device.PurchaseOrderNumber), v)
		}
		if v, ok := ctx["device-order_reference"]; ok && ne.StringValue(device.OrderReference) != v.(string) {
			return fmt.Errorf("device-order_reference does not match %v - %v", ne.StringValue(device.OrderReference), v)
		}
		if v, ok := ctx["device-interface_count"]; ok && ne.IntValue(device.InterfaceCount) != v.(int) {
			return fmt.Errorf("device-interface_count does not match %v - %v", ne.IntValue(device.InterfaceCount), v)
		}
		return nil
	}
}

func testAccNeDeviceSecondaryAttributes(device *ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	secCtx := make(map[string]interface{})
	for key, value := range ctx {
		secCtx[key] = value
	}
	secCtx["device-name"] = ctx["device-secondary_name"]
	secCtx["device-hostname"] = ctx["device-secondary_hostname"]
	secCtx["device-notifications"] = ctx["device-secondary_notifications"]
	return testAccNeDeviceAttributes(device, secCtx)
}

func testAccNeDeviceRedundancyAttributes(primary, secondary *ne.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ne.StringValue(primary.RedundancyType) != "PRIMARY" {
			return fmt.Errorf("redundancy_type does not match %v - %v", ne.StringValue(primary.RedundancyType), "PRIMARY")
		}
		if ne.StringValue(primary.RedundantUUID) != ne.StringValue(secondary.UUID) {
			return fmt.Errorf("redundant_id does not match %v - %v", ne.StringValue(primary.RedundantUUID), secondary.UUID)
		}
		if ne.StringValue(secondary.RedundancyType) != "SECONDARY" {
			return fmt.Errorf("redundancy_type does not match %v - %v", ne.StringValue(secondary.RedundancyType), "SECONDARY")
		}
		if ne.StringValue(secondary.RedundantUUID) != ne.StringValue(primary.UUID) {
			return fmt.Errorf("redundant_id does not match %v - %v", ne.StringValue(secondary.RedundantUUID), primary.UUID)
		}
		return nil
	}
}

func testAccNeDeviceStatusAttributes(device *ne.Device, provStatus, licStatus string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ne.StringValue(device.Status) != provStatus {
			return fmt.Errorf("status for device %q does not match  %v - %v", ne.StringValue(device.UUID), ne.StringValue(device.Status), provStatus)
		}
		if ne.StringValue(device.LicenseStatus) != licStatus {
			return fmt.Errorf("license_status for device %q does not match  %v - %v", ne.StringValue(device.UUID), ne.StringValue(device.LicenseStatus), licStatus)
		}
		return nil
	}
}

func testAccNeDeviceACLs(primary, secondary *ne.Device, primaryACL, secondaryACL *ne.ACLTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ne.StringValue(primaryACL.DeviceUUID) != ne.StringValue(primary.UUID) {
			return fmt.Errorf("Primary ACL %s device UUID does not match %v - %v", ne.StringValue(primaryACL.UUID), ne.StringValue(primaryACL.DeviceUUID), ne.StringValue(primary.UUID))
		}
		if ne.StringValue(secondaryACL.DeviceUUID) != ne.StringValue(secondary.UUID) {
			return fmt.Errorf("Secondary ACL %s device UUID does not match %v - %v", ne.StringValue(secondaryACL.UUID), ne.StringValue(secondaryACL.DeviceUUID), ne.StringValue(secondary.UUID))
		}
		if ne.StringValue(primaryACL.DeviceACLStatus) != ne.ACLDeviceStatusProvisioned {
			return fmt.Errorf("Primary ACL %s device_acl_status does not match %v - %v", ne.StringValue(primaryACL.UUID), ne.StringValue(primaryACL.DeviceACLStatus), ne.ACLDeviceStatusProvisioned)
		}
		if ne.StringValue(secondaryACL.DeviceACLStatus) != ne.ACLDeviceStatusProvisioned {
			return fmt.Errorf("Secondary ACL %s device_acl_status does not match %v - %v", ne.StringValue(secondaryACL.UUID), ne.StringValue(secondaryACL.DeviceACLStatus), ne.ACLDeviceStatusProvisioned)
		}
		return nil
	}
}

func (t *testAccConfig) withDevice() *testAccConfig {
	t.config += testAccNetworkDevice(t.ctx)
	return t
}

func (t *testAccConfig) withACL() *testAccConfig {
	t.config += testAccNetworkDeviceACL(t.ctx)
	return t
}

func (t *testAccConfig) withSSHKey() *testAccConfig {
	t.config += testAccNetworkDeviceSSHKey(t.ctx)
	return t
}

func testAccNetworkDevice(ctx map[string]interface{}) string {
	var config string
	config += nprintf(`
data "equinix_network_account" "test" {
  metro_code = "%{device-metro_code}"
  status     = "Active"
}`, ctx)

	config += nprintf(`
resource "equinix_network_device" "%{device-resourceName}" {
  self_managed          = %{device-self_managed}
  byol                  = %{device-byol}
  name                  = "%{device-name}"
  metro_code            = "%{device-metro_code}"
  type_code             = "%{device-type_code}"
  package_code          = "%{device-package_code}"
  notifications         = %{device-notifications}
  term_length           = %{device-term_length}
  account_number        = data.equinix_network_account.test.number
  version               = "%{device-version}"
  core_count            = %{device-core_count}`, ctx)
	if _, ok := ctx["device-purchase_order_number"]; ok {
		config += nprintf(`
  purchase_order_number = "%{device-purchase_order_number}"`, ctx)
	}
	if _, ok := ctx["device-purchase_order_number"]; ok {
		config += nprintf(`
  order_reference       = "%{device-order_reference}"`, ctx)
	}
	if _, ok := ctx["device-additional_bandwidth"]; ok {
		config += nprintf(`
  additional_bandwidth  = "%{device-additional_bandwidth}"`, ctx)
	}
	if _, ok := ctx["device-throughput"]; ok {
		config += nprintf(`
  throughput            = %{device-throughput}
  throughput_unit       = "%{device-throughput_unit}"`, ctx)
	}
	if _, ok := ctx["device-hostname"]; ok {
		config += nprintf(`
  hostname              = "%{device-hostname}"`, ctx)
	}
	if _, ok := ctx["device-interface_count"]; ok {
		config += nprintf(`
  interface_count       = %{device-interface_count}`, ctx)
	}
	if _, ok := ctx["acl-resourceName"]; ok {
		config += nprintf(`
  acl_template_id       = equinix_network_acl_template.%{acl-resourceName}.id`, ctx)
	}
	if _, ok := ctx["sshkey-resourceName"]; ok {
		config += nprintf(`
  ssh_key {
    username = "test"
    key_name = equinix_network_ssh_key.%{sshkey-resourceName}.name
  }`, ctx)
	}
	if _, ok := ctx["device-license_file"]; ok {
		config += nprintf(`
  license_file          = "%{device-license_file}"`, ctx)
	}
	if _, ok := ctx["device-secondary_name"]; ok {
		config += nprintf(`
  secondary_device {
    name                 = "%{device-secondary_name}"
    metro_code           = "%{device-metro_code}"
    notifications        = %{device-secondary_notifications}
	account_number       = data.equinix_network_account.test.number`, ctx)
		if _, ok := ctx["device-secondary_additional_bandwidth"]; ok {
			config += nprintf(`
    additional_bandwidth = "%{device-secondary_additional_bandwidth}"`, ctx)
		}
		if _, ok := ctx["device-secondary_hostname"]; ok {
			config += nprintf(`
    hostname             = "%{device-secondary_hostname}"`, ctx)
		}
		if _, ok := ctx["acl-secondary_resourceName"]; ok {
			config += nprintf(`
    acl_template_id      = equinix_network_acl_template.%{acl-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["sshkey-resourceName"]; ok {
			config += nprintf(`
    ssh_key {
      username = "test"
      key_name = equinix_network_ssh_key.%{sshkey-resourceName}.name
    }`, ctx)
		}
		if _, ok := ctx["device-secondary_license_file"]; ok {
			config += nprintf(`
    license_file         = "%{device-secondary_license_file}"`, ctx)
		}
		config += `
  }`
	}
	config += `
}`
	return config
}

func testAccNetworkDeviceACL(ctx map[string]interface{}) string {
	config := nprintf(`
resource "equinix_network_acl_template" "%{acl-resourceName}" {
  name          = "%{acl-name}"
  description   = "%{acl-description}"
  metro_code    = "%{acl-metroCode}"
  inbound_rule {
    subnets  = ["10.0.0.0/24"]
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	if _, ok := ctx["acl-secondary_name"]; ok {
		config += nprintf(`
resource "equinix_network_acl_template" "%{acl-secondary_resourceName}" {
  name          = "%{acl-secondary_name}"
  description   = "%{acl-secondary_description}"
  metro_code    = "%{acl-secondary_metroCode}"
  inbound_rule {
     subnets  = ["192.0.0.0/24"]
     protocol = "IP"
     src_port = "any"
     dst_port = "any"
  }
}`, ctx)
	}
	return config
}

func testAccNetworkDeviceSSHKey(ctx map[string]interface{}) string {
	return nprintf(`
resource "equinix_network_ssh_key" "%{sshkey-resourceName}" {
  name       = "%{sshkey-name}"
  public_key = "%{sshkey-public_key}"
}
`, ctx)
}
