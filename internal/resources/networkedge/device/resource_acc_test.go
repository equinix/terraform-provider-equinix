package device

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	tstResourcePrefix = "tfacc"

	networkDeviceAccountNameEnvVar          = "TF_ACC_NETWORK_DEVICE_BILLING_ACCOUNT_NAME"
	networkDeviceSecondaryAccountNameEnvVar = "TF_ACC_NETWORK_DEVICE_SECONDARY_BILLING_ACCOUNT_NAME"
	networkDeviceMetroEnvVar                = "TF_ACC_NETWORK_DEVICE_METRO"
	networkDeviceSecondaryMetroEnvVar       = "TF_ACC_NETWORK_DEVICE_SECONDARY_METRO"
	networkDeviceCSRSDWANLicenseFileEnvVar  = "TF_ACC_NETWORK_DEVICE_CSRSDWAN_LICENSE_FILE"
	networkDeviceVSRXLicenseFileEnvVar      = "TF_ACC_NETWORK_DEVICE_VSRX_LICENSE_FILE"
	networkDeviceVersaController1EnvVar     = "TF_ACC_NETWORK_DEVICE_VERSA_CONTROLLER1"
	networkDeviceVersaController2EnvVar     = "TF_ACC_NETWORK_DEVICE_VERSA_CONTROLLER2"
	networkDeviceVersaLocalIDEnvVar         = "TF_ACC_NETWORK_DEVICE_VERSA_LOCALID"
	networkDeviceVersaRemoteIDEnvVar        = "TF_ACC_NETWORK_DEVICE_VERSA_REMOTEID"
	networkDeviceVersaSerialNumberEnvVar    = "TF_ACC_NETWORK_DEVICE_VERSA_SERIAL"
	networkDeviceCGENIXLicenseKeyEnvVar     = "TF_ACC_NETWORK_DEVICE_CGENIX_LICENSE_KEY"
	networkDeviceCGENIXLicenseSecretEnvVar  = "TF_ACC_NETWORK_DEVICE_CGENIX_LICENSE_SECRET"
	networkDevicePANWLicenseTokenEnvVar     = "TF_ACC_NETWORK_DEVICE_PANW_LICENSE_TOKEN"
)

func init() {
	resource.AddTestSweepers("equinix_network_device", &resource.Sweeper{
		Name:         "equinix_network_device",
		Dependencies: []string{"equinix_network_device_link"},
		F:            testSweepNetworkDevice,
	})
}

func testSweepNetworkDevice(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping Network devices: %s", err)
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	devices, err := config.Ne.GetDevices([]string{
		ne.DeviceStateInitializing,
		ne.DeviceStateProvisioned,
		ne.DeviceStateProvisioning,
		ne.DeviceStateWaitingSecondary,
		ne.DeviceStateWaitingClusterNodes,
		ne.DeviceStateClusterSetUpInProgress,
		ne.DeviceStateFailed,
	})
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
		if err := config.Ne.DeleteDevice(ne.StringValue(device.UUID)); err != nil {
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

func TestAccNetworkDevice_CSR1000V_HA_Managed_Sub(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-account_name":            accountName.(string),
		"device-self_managed":            false,
		"device-byol":                    false,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-throughput":              500,
		"device-throughput_unit":         "Mbps",
		"device-metro_code":              metro.(string),
		"device-type_code":               "CSR1000V",
		"device-package_code":            "SEC",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", acctest.RandString(41)),
		"device-term_length":             1,
		"device-version":                 "16.09.05",
		"device-core_count":              2,
		"device-purchase_order_number":   acctest.RandString(10),
		"device-order_reference":         acctest.RandString(10),
		"device-interface_count":         24,
		"device-additional_bandwidth":    0,
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"user-resourceName":              "tst-user",
		"user-username":                  fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"user-password":                  acctest.RandString(10),
	}

	contextWithACLs := copyMap(context)
	contextWithACLs["acl-resourceName"] = "acl-pri"
	contextWithACLs["acl-name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithACLs["acl-description"] = acctest.RandString(50)
	contextWithACLs["acl-secondary_resourceName"] = "acl-sec"
	contextWithACLs["acl-secondary_name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithACLs["acl-secondary_description"] = acctest.RandString(50)
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", context["user-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-secondary_resourceName"].(string))
	var primary, secondary ne.Device
	var user ne.SSHUser
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.ParallelTest(t, resource.TestCase{
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
					testAccNeDeviceHAAttributes(deviceResourceName),
					testAccNeSSHUserExists(userResourceName, &user),
					testAccNeSSHUserAttributes(&user, []*ne.Device{&primary, &secondary}, context),
					resource.TestCheckResourceAttrSet(userResourceName, "uuid"),
				),
			},
			{
				ResourceName:      deviceResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: newTestAccConfig(contextWithACLs).withDevice().
					withSSHUser().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "acl_template_id"),
					testAccNeDeviceACL(priACLResourceName, &primary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.acl_template_id"),
					testAccNeDeviceACL(secACLResourceName, &secondary),
				),
			},
		},
	})
}

func TestAccNetworkDevice_CSR1000V_HA_Self_BYOL(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-account_name":            accountName.(string),
		"device-self_managed":            true,
		"device-byol":                    true,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-throughput":              500,
		"device-throughput_unit":         "Mbps",
		"device-metro_code":              metro.(string),
		"device-type_code":               "CSR1000V",
		"device-package_code":            "SEC",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":             1,
		"device-version":                 "16.09.05",
		"device-core_count":              2,
		"device-purchase_order_number":   acctest.RandString(10),
		"device-order_reference":         acctest.RandString(10),
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"sshkey-resourceName":            "test",
		"sshkey-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"sshkey-public_key":              "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXdzXBHaVpKpdO0udnB+4JOgUq7APO2rPXfrevvlZrps98AtlwXXVWZ5duRH5NFNfU4G9HCSiAPsebgjY0fG85tcShpXfHfACLt0tBW8XhfLQP2T6S50FQ1brBdURMDCMsD7duOXqvc0dlbs2/KcswHvuUmqVzob3bz7n1bQ48wIHsPg4ARqYhy5LN3OkllJH/6GEfqi8lKZx01/P/gmJMORcJujuOyXRB+F2iXBVYdhjML3Qg4+tEekBcVZOxUbERRZ0pvQ52Y6wUhn2VsjljixyqeOdmD0m6DayDQgSWms6bKPpBqN7zhXXk4qe8bXT4tQQba65b2CQ2A91jw2KgM/YZNmjyUJ+Rf1cQosJf9twqbAZDZ6rAEmj9zzvQ5vD/CGuzxdVMkePLlUK4VGjPu7cVzhXrnq4318WqZ5/lNiCST8NQ0fssChN8ANUzr/p/wwv3faFMVNmjxXTZMsbMFT/fbb2MVVuqNFN65drntlg6/xEao8gZROuRYiakBx8= user@host",
	}

	contextWithACLs := copyMap(context)
	contextWithACLs["acl-resourceName"] = "acl-pri"
	contextWithACLs["acl-name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithACLs["acl-description"] = acctest.RandString(50)
	contextWithACLs["acl-secondary_resourceName"] = "acl-sec"
	contextWithACLs["acl-secondary_name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithACLs["acl-secondary_description"] = acctest.RandString(50)
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-secondary_resourceName"].(string))
	var primary, secondary ne.Device
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.ParallelTest(t, resource.TestCase{
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
					testAccNeDeviceHAAttributes(deviceResourceName),
				),
			},
			{
				Config: newTestAccConfig(contextWithACLs).withDevice().
					withSSHKey().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "acl_template_id"),
					testAccNeDeviceACL(priACLResourceName, &primary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.acl_template_id"),
					testAccNeDeviceACL(secACLResourceName, &secondary),
				),
			},
		},
	})
}

func TestAccNetworkDevice_vSRX_HA_Managed_Sub(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-account_name":            accountName.(string),
		"device-self_managed":            false,
		"device-byol":                    false,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-metro_code":              metro.(string),
		"device-type_code":               "VSRX",
		"device-package_code":            "STD",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":             1,
		"device-version":                 "19.2R2.7",
		"device-core_count":              2,
		"device-purchase_order_number":   acctest.RandString(10),
		"device-order_reference":         acctest.RandString(10),
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
	}

	contextWithChanges := copyMap(context)
	contextWithChanges["user-resourceName"] = "test"
	contextWithChanges["user-username"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithChanges["user-password"] = acctest.RandString(10)
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", contextWithChanges["user-resourceName"].(string))
	var primary, secondary ne.Device
	var user ne.SSHUser
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateRegistered),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateRegistered),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					testAccNeDeviceHAAttributes(deviceResourceName),
				),
			},
			{
				Config: newTestAccConfig(contextWithChanges).withDevice().withSSHUser().build(),
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

func TestAccNetworkDevice_vSRX_HA_Managed_BYOL(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	licenseFile, _ := schema.EnvDefaultFunc(networkDeviceVSRXLicenseFileEnvVar, "")()
	if licenseFile.(string) == "" {
		t.Skip("Skipping TestAccNetworkDevice_vSRX_HA_Managed_BYOL test since TF_ACC_NETWORK_DEVICE_VSRX_LICENSE_FILE env var is not defined with a valid license file")
	}
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-account_name":            accountName.(string),
		"device-self_managed":            false,
		"device-byol":                    true,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-license_file":            licenseFile.(string),
		"device-metro_code":              metro.(string),
		"device-type_code":               "VSRX",
		"device-package_code":            "STD",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":             1,
		"device-version":                 "19.2R2.7",
		"device-core_count":              2,
		"device-purchase_order_number":   acctest.RandString(10),
		"device-order_reference":         acctest.RandString(10),
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_license_file":  licenseFile.(string),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"acl-resourceName":               "acl-pri",
		"acl-name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-description":                acctest.RandString(50),
		"acl-secondary_resourceName":     "acl-sec",
		"acl-secondary_name":             fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-secondary_description":      acctest.RandString(50),
	}

	contextWithChanges := copyMap(context)
	contextWithChanges["device-name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithChanges["device-additional_bandwidth"] = 100
	contextWithChanges["device-notifications"] = []string{"jerry@equinix.com", "tom@equinix.com"}
	contextWithChanges["device-secondary_name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithChanges["device-secondary_additional_bandwidth"] = 100
	contextWithChanges["device-secondary_notifications"] = []string{"miki@equinix.com", "mini@equinix.com"}
	contextWithChanges["user-resourceName"] = "test"
	contextWithChanges["user-username"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithChanges["user-password"] = acctest.RandString(10)
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", contextWithChanges["user-resourceName"].(string))
	var primary, secondary ne.Device
	var user ne.SSHUser
	resource.ParallelTest(t, resource.TestCase{
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
					testAccNeDeviceHAAttributes(deviceResourceName),
					resource.TestCheckResourceAttrSet(deviceResourceName, "license_file_id"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.license_file_id"),
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

func TestAccNetworkDevice_vSRX_HA_Self_BYOL(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-account_name":            accountName.(string),
		"device-self_managed":            true,
		"device-byol":                    true,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-metro_code":              metro.(string),
		"device-type_code":               "VSRX",
		"device-package_code":            "STD",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":             1,
		"device-version":                 "19.2R2.7",
		"device-core_count":              2,
		"device-purchase_order_number":   acctest.RandString(10),
		"device-order_reference":         acctest.RandString(10),
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"acl-resourceName":               "acl-pri",
		"acl-name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-description":                acctest.RandString(50),
		"acl-secondary_resourceName":     "acl-sec",
		"acl-secondary_name":             fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-secondary_description":      acctest.RandString(50),
		"sshkey-resourceName":            "test",
		"sshkey-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"sshkey-public_key":              "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXdzXBHaVpKpdO0udnB+4JOgUq7APO2rPXfrevvlZrps98AtlwXXVWZ5duRH5NFNfU4G9HCSiAPsebgjY0fG85tcShpXfHfACLt0tBW8XhfLQP2T6S50FQ1brBdURMDCMsD7duOXqvc0dlbs2/KcswHvuUmqVzob3bz7n1bQ48wIHsPg4ARqYhy5LN3OkllJH/6GEfqi8lKZx01/P/gmJMORcJujuOyXRB+F2iXBVYdhjML3Qg4+tEekBcVZOxUbERRZ0pvQ52Y6wUhn2VsjljixyqeOdmD0m6DayDQgSWms6bKPpBqN7zhXXk4qe8bXT4tQQba65b2CQ2A91jw2KgM/YZNmjyUJ+Rf1cQosJf9twqbAZDZ6rAEmj9zzvQ5vD/CGuzxdVMkePLlUK4VGjPu7cVzhXrnq4318WqZ5/lNiCST8NQ0fssChN8ANUzr/p/wwv3faFMVNmjxXTZMsbMFT/fbb2MVVuqNFN65drntlg6/xEao8gZROuRYiakBx8= user@host",
	}

	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	var primary, secondary ne.Device
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withACL().withSSHKey().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					testAccNeDeviceHAAttributes(deviceResourceName),
				),
			},
		},
	})
}

func TestAccNetworkDevice_PaloAlto_HA_Managed_Sub(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-account_name":            accountName.(string),
		"device-self_managed":            false,
		"device-byol":                    false,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-metro_code":              metro.(string),
		"device-type_code":               "PA-VM",
		"device-package_code":            "VM100",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":             1,
		"device-version":                 "9.0.4",
		"device-core_count":              2,
		"device-purchase_order_number":   acctest.RandString(10),
		"device-order_reference":         acctest.RandString(10),
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"acl-resourceName":               "acl-pri",
		"acl-name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-description":                acctest.RandString(50),
		"acl-secondary_resourceName":     "acl-sec",
		"acl-secondary_name":             fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-secondary_description":      acctest.RandString(50),
	}

	contextWithChanges := copyMap(context)
	contextWithChanges["device-additional_bandwidth"] = 50
	contextWithChanges["device-secondary_additional_bandwidth"] = 50
	contextWithChanges["user-resourceName"] = "tst-user"
	contextWithChanges["user-username"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithChanges["user-password"] = acctest.RandString(10)
	var primary, secondary ne.Device
	var primaryACL, secondaryACL ne.ACLTemplate
	var user ne.SSHUser
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-secondary_resourceName"].(string))
	userResourceName := fmt.Sprintf("equinix_network_ssh_user.%s", contextWithChanges["user-resourceName"].(string))
	resource.ParallelTest(t, resource.TestCase{
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
					testAccNeDeviceHAAttributes(deviceResourceName),
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					resource.TestCheckResourceAttrSet(deviceResourceName, "acl_template_id"),
					testAccNeDeviceACL(priACLResourceName, &primary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.acl_template_id"),
					testAccNeDeviceACL(secACLResourceName, &secondary),
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
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":            "test",
		"device-account_name":            accountName.(string),
		"device-self_managed":            true,
		"connectivity":                   "PRIVATE",
		"device-byol":                    true,
		"device-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-metro_code":              metro.(string),
		"device-type_code":               "PA-VM",
		"device-package_code":            "VM100",
		"device-notifications":           []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":                fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":             1,
		"device-version":                 "9.0.4",
		"device-core_count":              2,
		"device-purchase_order_number":   acctest.RandString(10),
		"device-order_reference":         acctest.RandString(10),
		"device-secondary_name":          fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_hostname":      fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-secondary_notifications": []string{"secondary@equinix.com"},
		"sshkey-resourceName":            "test",
		"sshkey-name":                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"sshkey-public_key":              "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXdzXBHaVpKpdO0udnB+4JOgUq7APO2rPXfrevvlZrps98AtlwXXVWZ5duRH5NFNfU4G9HCSiAPsebgjY0fG85tcShpXfHfACLt0tBW8XhfLQP2T6S50FQ1brBdURMDCMsD7duOXqvc0dlbs2/KcswHvuUmqVzob3bz7n1bQ48wIHsPg4ARqYhy5LN3OkllJH/6GEfqi8lKZx01/P/gmJMORcJujuOyXRB+F2iXBVYdhjML3Qg4+tEekBcVZOxUbERRZ0pvQ52Y6wUhn2VsjljixyqeOdmD0m6DayDQgSWms6bKPpBqN7zhXXk4qe8bXT4tQQba65b2CQ2A91jw2KgM/YZNmjyUJ+Rf1cQosJf9twqbAZDZ6rAEmj9zzvQ5vD/CGuzxdVMkePLlUK4VGjPu7cVzhXrnq4318WqZ5/lNiCST8NQ0fssChN8ANUzr/p/wwv3faFMVNmjxXTZMsbMFT/fbb2MVVuqNFN65drntlg6/xEao8gZROuRYiakBx8= user@host",
	}

	contextWithACLs := copyMap(context)
	contextWithACLs["acl-resourceName"] = "acl-pri"
	contextWithACLs["acl-name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithACLs["acl-description"] = acctest.RandString(50)
	contextWithACLs["acl-secondary_resourceName"] = "acl-sec"
	contextWithACLs["acl-secondary_name"] = fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6))
	contextWithACLs["acl-secondary_description"] = acctest.RandString(50)
	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", contextWithACLs["acl-secondary_resourceName"].(string))
	var primary, secondary ne.Device
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.ParallelTest(t, resource.TestCase{
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
					testAccNeDeviceHAAttributes(deviceResourceName),
				),
			},
			{
				Config: newTestAccConfig(contextWithACLs).withDevice().withSSHKey().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "acl_template_id"),
					testAccNeDeviceACL(priACLResourceName, &primary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.acl_template_id"),
					testAccNeDeviceACL(secACLResourceName, &secondary),
				),
			},
		},
	})
}

func TestAccNetworkDevice_CSRSDWAN_HA_Self_BYOL(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	licFile, _ := schema.EnvDefaultFunc(networkDeviceCSRSDWANLicenseFileEnvVar, "test-fixtures/CSRSDWAN.cfg")()
	context := map[string]interface{}{
		"device-resourceName":                           "test",
		"device-account_name":                           accountName.(string),
		"device-self_managed":                           true,
		"device-byol":                                   true,
		"device-name":                                   fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-metro_code":                             metro.(string),
		"device-license_file":                           licFile.(string),
		"device-throughput":                             250,
		"device-throughput_unit":                        "Mbps",
		"device-type_code":                              "CSRSDWAN",
		"device-package_code":                           "ESSENTIALS",
		"device-notifications":                          []string{"marry@equinix.com", "john@equinix.com"},
		"device-term_length":                            1,
		"device-version":                                "16.12.3",
		"device-core_count":                             2,
		"device-purchase_order_number":                  acctest.RandString(10),
		"device-order_reference":                        acctest.RandString(10),
		"device-vendorConfig_enabled":                   true,
		"device-vendorConfig_siteId":                    "10",
		"device-vendorConfig_systemIpAddress":           "1.1.1.1",
		"device-secondary_name":                         fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_license_file":                 licFile.(string),
		"device-secondary_notifications":                []string{"secondary@equinix.com"},
		"device-secondary_vendorConfig_enabled":         true,
		"device-secondary_vendorConfig_siteId":          "20",
		"device-secondary_vendorConfig_systemIpAddress": "2.2.2.2",
		"acl-resourceName":                              "acl-pri",
		"acl-name":                                      fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-description":                               acctest.RandString(50),
		"acl-secondary_resourceName":                    "acl-sec",
		"acl-secondary_name":                            fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-secondary_description":                     acctest.RandString(50),
	}

	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-secondary_resourceName"].(string))
	var primary, secondary ne.Device
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					resource.TestCheckResourceAttrSet(deviceResourceName, "acl_template_id"),
					testAccNeDeviceACL(priACLResourceName, &primary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.acl_template_id"),
					testAccNeDeviceACL(secACLResourceName, &secondary),
					testAccNeDeviceHAAttributes(deviceResourceName),
					resource.TestCheckResourceAttrSet(deviceResourceName, "license_file_id"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.license_file_id"),
				),
			},
		},
	})
}

func TestAccNetworkDevice_Versa_HA_Self_BYOL(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	controller1, _ := schema.EnvDefaultFunc(networkDeviceVersaController1EnvVar, "1.1.1.1")()
	controller2, _ := schema.EnvDefaultFunc(networkDeviceVersaController2EnvVar, "2.2.2.2")()
	localID, _ := schema.EnvDefaultFunc(networkDeviceVersaLocalIDEnvVar, "test@versa.com")()
	remoteID, _ := schema.EnvDefaultFunc(networkDeviceVersaRemoteIDEnvVar, "test@versa.com")()
	serialNumber, _ := schema.EnvDefaultFunc(networkDeviceVersaSerialNumberEnvVar, "Test")()
	context := map[string]interface{}{
		"device-resourceName":                        "test",
		"device-account_name":                        accountName.(string),
		"device-self_managed":                        true,
		"device-byol":                                true,
		"device-name":                                fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-metro_code":                          metro.(string),
		"device-type_code":                           "VERSA_SDWAN",
		"device-package_code":                        "FLEX_VNF_2",
		"device-notifications":                       []string{"marry@equinix.com", "john@equinix.com"},
		"device-term_length":                         1,
		"device-version":                             "16.1R2S8",
		"device-core_count":                          2,
		"device-purchase_order_number":               acctest.RandString(10),
		"device-order_reference":                     acctest.RandString(10),
		"device-vendorConfig_enabled":                true,
		"device-vendorConfig_controller1":            controller1.(string),
		"device-vendorConfig_controller2":            controller2.(string),
		"device-vendorConfig_localId":                localID.(string),
		"device-vendorConfig_remoteId":               remoteID.(string),
		"device-vendorConfig_serialNumber":           serialNumber.(string),
		"device-secondary_name":                      fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_notifications":             []string{"secondary@equinix.com"},
		"device-secondary_vendorConfig_enabled":      true,
		"device-secondary_vendorConfig_controller1":  controller1.(string),
		"device-secondary_vendorConfig_controller2":  controller2.(string),
		"device-secondary_vendorConfig_localId":      localID.(string),
		"device-secondary_vendorConfig_remoteId":     remoteID.(string),
		"device-secondary_vendorConfig_serialNumber": serialNumber.(string),
		"acl-resourceName":                           "acl-pri",
		"acl-name":                                   fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-description":                            acctest.RandString(50),
		"acl-secondary_resourceName":                 "acl-sec",
		"acl-secondary_name":                         fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-secondary_description":                  acctest.RandString(50),
	}

	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-secondary_resourceName"].(string))
	var primary, secondary ne.Device
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					resource.TestCheckResourceAttrSet(deviceResourceName, "acl_template_id"),
					testAccNeDeviceACL(priACLResourceName, &primary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.acl_template_id"),
					testAccNeDeviceACL(secACLResourceName, &secondary),
					testAccNeDeviceHAAttributes(deviceResourceName),
				),
			},
		},
	})
}

func TestAccNetworkDevice_CGENIX_HA_Self_BYOL(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	licenseKey, _ := schema.EnvDefaultFunc(networkDeviceCGENIXLicenseKeyEnvVar, acctest.RandString(10))()
	licenseSecret, _ := schema.EnvDefaultFunc(networkDeviceCGENIXLicenseSecretEnvVar, acctest.RandString(10))()
	context := map[string]interface{}{
		"device-resourceName":                         "test",
		"device-account_name":                         accountName.(string),
		"device-self_managed":                         true,
		"device-byol":                                 true,
		"device-name":                                 fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-metro_code":                           metro.(string),
		"device-type_code":                            "CGENIXSDWAN",
		"device-package_code":                         "3102V",
		"device-notifications":                        []string{"marry@equinix.com", "john@equinix.com"},
		"device-term_length":                          1,
		"device-version":                              "5.2.1-b11",
		"device-core_count":                           2,
		"device-purchase_order_number":                acctest.RandString(10),
		"device-order_reference":                      acctest.RandString(10),
		"device-vendorConfig_enabled":                 true,
		"device-vendorConfig_licenseKey":              licenseKey.(string),
		"device-vendorConfig_licenseSecret":           licenseSecret.(string),
		"device-secondary_name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_notifications":              []string{"secondary@equinix.com"},
		"device-secondary_vendorConfig_enabled":       true,
		"device-secondary_vendorConfig_licenseKey":    licenseKey.(string),
		"device-secondary_vendorConfig_licenseSecret": licenseSecret.(string),
		"acl-resourceName":                            "acl-pri",
		"acl-name":                                    fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-description":                             acctest.RandString(50),
		"acl-secondary_resourceName":                  "acl-sec",
		"acl-secondary_name":                          fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-secondary_description":                   acctest.RandString(50),
	}

	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	priACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-resourceName"].(string))
	secACLResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-secondary_resourceName"].(string))
	var primary, secondary ne.Device
	var primaryACL, secondaryACL ne.ACLTemplate
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceSecondaryExists(&primary, &secondary),
					testAccNeDeviceSecondaryAttributes(&secondary, context),
					testAccNeDeviceStatusAttributes(&secondary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceRedundancyAttributes(&primary, &secondary),
					testAccNetworkACLTemplateExists(priACLResourceName, &primaryACL),
					testAccNetworkACLTemplateExists(secACLResourceName, &secondaryACL),
					resource.TestCheckResourceAttrSet(deviceResourceName, "acl_template_id"),
					testAccNeDeviceACL(priACLResourceName, &primary),
					resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.acl_template_id"),
					testAccNeDeviceACL(secACLResourceName, &secondary),
					testAccNeDeviceHAAttributes(deviceResourceName),
				),
			},
		},
	})
}

func TestAccNetworkDevice_PaloAlto_Cluster_Self_BYOL(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	licenseToken, _ := schema.EnvDefaultFunc(networkDevicePANWLicenseTokenEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":                "test",
		"device-account_name":                accountName.(string),
		"device-self_managed":                true,
		"device-byol":                        true,
		"device-name":                        fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-metro_code":                  metro.(string),
		"device-type_code":                   "PA-VM",
		"device-package_code":                "VM100",
		"device-notifications":               []string{"marry@equinix.com", "john@equinix.com"},
		"device-term_length":                 1,
		"device-version":                     "10.1.3",
		"device-interface_count":             10,
		"device-core_count":                  2,
		"device-cluster_name":                fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-node0_vendorConfig_enabled":  true,
		"device-node0_vendorConfig_hostname": fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-node0_license_token":         licenseToken.(string),
		"device-node1_vendorConfig_enabled":  true,
		"device-node1_vendorConfig_hostname": fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-node1_license_token":         licenseToken.(string),
		"sshkey-resourceName":                "test",
		"sshkey-name":                        fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"sshkey-public_key":                  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXdzXBHaVpKpdO0udnB+4JOgUq7APO2rPXfrevvlZrps98AtlwXXVWZ5duRH5NFNfU4G9HCSiAPsebgjY0fG85tcShpXfHfACLt0tBW8XhfLQP2T6S50FQ1brBdURMDCMsD7duOXqvc0dlbs2/KcswHvuUmqVzob3bz7n1bQ48wIHsPg4ARqYhy5LN3OkllJH/6GEfqi8lKZx01/P/gmJMORcJujuOyXRB+F2iXBVYdhjML3Qg4+tEekBcVZOxUbERRZ0pvQ52Y6wUhn2VsjljixyqeOdmD0m6DayDQgSWms6bKPpBqN7zhXXk4qe8bXT4tQQba65b2CQ2A91jw2KgM/YZNmjyUJ+Rf1cQosJf9twqbAZDZ6rAEmj9zzvQ5vD/CGuzxdVMkePLlUK4VGjPu7cVzhXrnq4318WqZ5/lNiCST8NQ0fssChN8ANUzr/p/wwv3faFMVNmjxXTZMsbMFT/fbb2MVVuqNFN65drntlg6/xEao8gZROuRYiakBx8= user@host",
		"acl-resourceName":                   "acl-cluster",
		"acl-name":                           fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"acl-description":                    acctest.RandString(50),
		"mgmtAcl-resourceName":               "mgmtAcl-cluster",
		"mgmtAcl-name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"mgmtAcl-description":                acctest.RandString(50),
	}

	deviceResourceName := fmt.Sprintf("equinix_network_device.%s", context["device-resourceName"].(string))
	clusterAclResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["acl-resourceName"].(string))
	clusterMgmtAclResourceName := fmt.Sprintf("equinix_network_acl_template.%s", context["mgmtAcl-resourceName"].(string))
	var primary ne.Device
	var wanAcl, mgmtAcl ne.ACLTemplate
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withSSHKey().withACL().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceExists(deviceResourceName, &primary),
					testAccNeDeviceAttributes(&primary, context),
					testAccNeDeviceStatusAttributes(&primary, ne.DeviceStateProvisioned, ne.DeviceLicenseStateApplied),
					testAccNeDeviceClusterAttributes(deviceResourceName),
					testAccNeDeviceClusterNodeAttributes(&primary, context),
					testAccNetworkACLTemplateExists(clusterAclResourceName, &wanAcl),
					testAccNetworkACLTemplateExists(clusterMgmtAclResourceName, &mgmtAcl),
					resource.TestCheckResourceAttrSet(deviceResourceName, "acl_template_id"),
					resource.TestCheckResourceAttrSet(deviceResourceName, "mgmt_acl_template_uuid"),
					testAccNeDeviceACL(clusterAclResourceName, &primary),
				),
			},
			{
				ResourceName:            deviceResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_details.0.node0.0.license_token", "cluster_details.0.node1.0.license_token", "mgmt_acl_template_uuid"},
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
		client := testAccProvider.Meta().(*config.Config).Ne
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
		client := testAccProvider.Meta().(*config.Config).Ne
		resp, err := client.GetDevice(ne.StringValue(primary.RedundantUUID))
		if err != nil {
			return fmt.Errorf("error when fetching network device '%s': %s", ne.StringValue(primary.RedundantUUID), err)
		}
		*secondary = *resp
		return nil
	}
}

func testAccNeDevicePairExists(resourceName string, primary, secondary *ne.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		client := testAccProvider.Meta().(*config.Config).Ne
		resp, err := client.GetDevice(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching primary network device '%s': %s", rs.Primary.ID, err)
		}
		*primary = *resp
		resp, err = client.GetDevice(ne.StringValue(resp.RedundantUUID))
		if err != nil {
			return fmt.Errorf("error when fetching secondary network device '%s': %s", rs.Primary.ID, err)
		}
		*secondary = *resp
		return nil
	}
}

func testAccNeDeviceAttributes(device *ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["device-name"]; ok && !(ne.StringValue(device.Name) == v.(string) || ne.StringValue(device.Name) == v.(string)+"-Node0") {
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
		if v, ok := ctx["device-notifications"]; ok && !comparisons.SlicesMatch(device.Notifications, v.([]string)) {
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
		if v, ok := ctx["device-additional_bandwidth"]; ok && ne.IntValue(device.AdditionalBandwidth) != v.(int) {
			return fmt.Errorf("device-additional_bandwidth does not match %v - %v", ne.IntValue(device.AdditionalBandwidth), v)
		}
		if v, ok := ctx["device-vendorConfig_siteId"]; ok && device.VendorConfiguration["siteId"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_siteId does not match %v - %v", device.VendorConfiguration["siteId"], v)
		}
		if v, ok := ctx["device-vendorConfig_systemIpAddress"]; ok && device.VendorConfiguration["systemIpAddress"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_systemIpAddress does not match %v - %v", device.VendorConfiguration["systemIpAddress"], v)
		}
		if v, ok := ctx["device-vendorConfig_licenseKey"]; ok && device.VendorConfiguration["licenseKey"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_licenseKey does not match %v - %v", device.VendorConfiguration["licenseKey"], v)
		}
		if v, ok := ctx["device-vendorConfig_licenseSecret"]; ok && device.VendorConfiguration["licenseSecret"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_licenseSecret does not match %v - %v", device.VendorConfiguration["licenseSecret"], v)
		}
		if v, ok := ctx["device-vendorConfig_controller1"]; ok && device.VendorConfiguration["controller1"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_controller1 does not match %v - %v", device.VendorConfiguration["controller1"], v)
		}
		if v, ok := ctx["device-vendorConfig_controller2"]; ok && device.VendorConfiguration["controller2"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_controller2 does not match %v - %v", device.VendorConfiguration["controller2"], v)
		}
		if v, ok := ctx["device-vendorConfig_localId"]; ok && device.VendorConfiguration["localId"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_localId does not match %v - %v", device.VendorConfiguration["localId"], v)
		}
		if v, ok := ctx["device-vendorConfig_remoteId"]; ok && device.VendorConfiguration["remoteId"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_remoteId does not match %v - %v", device.VendorConfiguration["remoteId"], v)
		}
		if v, ok := ctx["device-vendorConfig_serialNumber"]; ok && device.VendorConfiguration["serialNumber"] != v.(string) {
			return fmt.Errorf("device-vendorConfig_serialNumber does not match %v - %v", device.VendorConfiguration["serialNumber"], v)
		}
		if v, ok := ctx["connectivity"]; ok && ne.StringValue(device.Connectivity) != v.(string) {
			return fmt.Errorf("connectivity does not match %v - %v", ne.StringValue(device.Connectivity), v)
		}
		return nil
	}
}

func testAccNeDeviceSecondaryAttributes(device *ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	secCtx := make(map[string]interface{})
	for key, value := range ctx {
		secCtx[key] = value
	}
	if v, ok := ctx["device-secondary_name"]; ok {
		secCtx["device-name"] = v
	}
	if v, ok := ctx["device-secondary_hostname"]; ok {
		secCtx["device-hostname"] = v
	}
	if v, ok := ctx["device-secondary_notifications"]; ok {
		secCtx["device-notifications"] = v
	}
	if v, ok := ctx["device-secondary_additional_bandwidth"]; ok {
		secCtx["device-additional_bandwidth"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_siteId"]; ok {
		secCtx["device-vendorConfig_siteId"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_systemIpAddress"]; ok {
		secCtx["device-vendorConfig_systemIpAddress"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_licenseKey"]; ok {
		secCtx["device-vendorConfig_licenseKey"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_licenseSecret"]; ok {
		secCtx["device-vendorConfig_licenseSecret"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_controller1"]; ok {
		secCtx["device-vendorConfig_controller1"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_controller2"]; ok {
		secCtx["device-vendorConfig_controller2"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_localId"]; ok {
		secCtx["device-vendorConfig_localId"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_remoteId"]; ok {
		secCtx["device-vendorConfig_remoteId"] = v
	}
	if v, ok := ctx["device-secondary_vendorConfig_serialNumber"]; ok {
		secCtx["device-vendorConfig_serialNumber"] = v
	}
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

func testAccNeDeviceACL(resourceName string, device *ne.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		templateId := rs.Primary.ID
		deviceID := ne.StringValue(device.UUID)
		client := testAccProvider.Meta().(*config.Config).Ne
		if ne.StringValue(device.ACLTemplateUUID) != rs.Primary.ID {
			return fmt.Errorf("acl_template_id for device %s does not match %v - %v", deviceID, ne.StringValue(device.ACLTemplateUUID), templateId)
		}
		deviceACL, err := client.GetDeviceACLDetails(deviceID)
		if err != nil {
			return fmt.Errorf("error when fetching ACL details for device '%s': %s", deviceID, err)
		}
		if ne.StringValue(deviceACL.Status) != ne.ACLDeviceStatusProvisioned {
			return fmt.Errorf("device_acl_status for device %s does not match %v - %v", deviceID, ne.StringValue(deviceACL.Status), ne.ACLDeviceStatusProvisioned)
		}
		return nil
	}
}

func testAccNeDeviceHAAttributes(deviceResourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(deviceResourceName, "uuid"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "ibx"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "region"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "ssh_ip_address"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "ssh_ip_fqdn"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.uuid"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.ibx"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.region"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.ssh_ip_address"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "secondary_device.0.ssh_ip_fqdn"),
	)
}

func testAccNeDeviceClusterAttributes(deviceResourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(deviceResourceName, "uuid"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "cluster_details.0.cluster_id"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "cluster_details.0.num_of_nodes"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "cluster_details.0.node0.0.uuid"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "cluster_details.0.node0.0.name"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "cluster_details.0.node1.0.uuid"),
		resource.TestCheckResourceAttrSet(deviceResourceName, "cluster_details.0.node1.0.name"),
	)
}

func testAccNeDeviceClusterNodeAttributes(device *ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cluster := device.ClusterDetails
		if v, ok := ctx["device-cluster_name"]; ok && ne.StringValue(cluster.ClusterName) != v.(string) {
			return fmt.Errorf("cluster_name does not match %v - %v", ne.StringValue(cluster.ClusterName), v)
		}
		if v, ok := ctx["device-node0_vendorConfig_hostname"]; ok && cluster.Node0.VendorConfiguration["hostname"] != v.(string) {
			return fmt.Errorf("device-node0_vendorConfig_hostname does not match %v - %v", cluster.Node0.VendorConfiguration["hostname"], v)
		}
		if v, ok := ctx["device-node0_vendorConfig_adminPassword"]; ok && cluster.Node0.VendorConfiguration["adminPassword"] != v.(string) {
			return fmt.Errorf("device-node0_vendorConfig_adminPassword does not match %v - %v", cluster.Node0.VendorConfiguration["adminPassword"], v)
		}
		if v, ok := ctx["device-node0_vendorConfig_controller1"]; ok && cluster.Node0.VendorConfiguration["controller1"] != v.(string) {
			return fmt.Errorf("device-node0_vendorConfig_controller1 does not match %v - %v", cluster.Node0.VendorConfiguration["controller1"], v)
		}
		if v, ok := ctx["device-node0_vendorConfig_activationKey"]; ok && cluster.Node0.VendorConfiguration["activationKey"] != v.(string) {
			return fmt.Errorf("device-node0_vendorConfig_activationKey does not match %v - %v", cluster.Node0.VendorConfiguration["activationKey"], v)
		}
		if v, ok := ctx["device-node0_vendorConfig_controllerFqdn"]; ok && cluster.Node0.VendorConfiguration["controllerFqdn"] != v.(string) {
			return fmt.Errorf("device-node0_vendorConfig_controllerFqdn does not match %v - %v", cluster.Node0.VendorConfiguration["controllerFqdn"], v)
		}
		if v, ok := ctx["device-node0_vendorConfig_rootPassword"]; ok && cluster.Node0.VendorConfiguration["rootPassword"] != v.(string) {
			return fmt.Errorf("device-node0_vendorConfig_rootPassword does not match %v - %v", cluster.Node0.VendorConfiguration["rootPassword"], v)
		}
		if v, ok := ctx["device-node1_vendorConfig_hostname"]; ok && cluster.Node1.VendorConfiguration["hostname"] != v.(string) {
			return fmt.Errorf("device-node1_vendorConfig_hostname does not match %v - %v", cluster.Node1.VendorConfiguration["hostname"], v)
		}
		if v, ok := ctx["device-node1_vendorConfig_adminPassword"]; ok && cluster.Node1.VendorConfiguration["adminPassword"] != v.(string) {
			return fmt.Errorf("device-node1_vendorConfig_adminPassword does not match %v - %v", cluster.Node1.VendorConfiguration["adminPassword"], v)
		}
		if v, ok := ctx["device-node1_vendorConfig_controller1"]; ok && cluster.Node1.VendorConfiguration["controller1"] != v.(string) {
			return fmt.Errorf("device-node1_vendorConfig_controller1 does not match %v - %v", cluster.Node1.VendorConfiguration["controller1"], v)
		}
		if v, ok := ctx["device-node1_vendorConfig_activationKey"]; ok && cluster.Node1.VendorConfiguration["activationKey"] != v.(string) {
			return fmt.Errorf("device-node1_vendorConfig_activationKey does not match %v - %v", cluster.Node1.VendorConfiguration["activationKey"], v)
		}
		if v, ok := ctx["device-node1_vendorConfig_controllerFqdn"]; ok && cluster.Node1.VendorConfiguration["controllerFqdn"] != v.(string) {
			return fmt.Errorf("device-node1_vendorConfig_controllerFqdn does not match %v - %v", cluster.Node1.VendorConfiguration["controllerFqdn"], v)
		}
		if v, ok := ctx["device-node1_vendorConfig_rootPassword"]; ok && cluster.Node1.VendorConfiguration["rootPassword"] != v.(string) {
			return fmt.Errorf("device-node1_vendorConfig_rootPassword does not match %v - %v", cluster.Node1.VendorConfiguration["rootPassword"], v)
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
	config += nprintf.NPrintf(`
data "equinix_network_account" "test" {
  metro_code = "%{device-metro_code}"
  status     = "Active"`, ctx)
	if v, ok := ctx["device-account_name"]; ok && !comparisons.IsEmpty(v) {
		config += nprintf.NPrintf(`
  name = "%{device-account_name}"`, ctx)
	}
	config += nprintf.NPrintf(`
}`, ctx)
	if _, ok := ctx["device-secondary_metro_code"]; ok {
		config += nprintf.NPrintf(`
data "equinix_network_account" "test-secondary" {
  metro_code = "%{device-secondary_metro_code}"
  status     = "Active"`, ctx)
		if v, ok := ctx["device-secondary_account_name"]; ok && !comparisons.IsEmpty(v) {
			config += nprintf.NPrintf(`
  name = "%{device-secondary_account_name}"`, ctx)
		}
		config += nprintf.NPrintf(` 
}`, ctx)
	}
	config += nprintf.NPrintf(`
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
		config += nprintf.NPrintf(`
  purchase_order_number = "%{device-purchase_order_number}"`, ctx)
	}
	if _, ok := ctx["device-purchase_order_number"]; ok {
		config += nprintf.NPrintf(`
  order_reference       = "%{device-order_reference}"`, ctx)
	}
	if _, ok := ctx["device-additional_bandwidth"]; ok {
		config += nprintf.NPrintf(`
  additional_bandwidth  = %{device-additional_bandwidth}`, ctx)
	}
	if _, ok := ctx["device-throughput"]; ok {
		config += nprintf.NPrintf(`
  throughput            = %{device-throughput}
  throughput_unit       = "%{device-throughput_unit}"`, ctx)
	}
	if _, ok := ctx["device-hostname"]; ok {
		config += nprintf.NPrintf(`
  hostname              = "%{device-hostname}"`, ctx)
	}
	if _, ok := ctx["device-interface_count"]; ok {
		config += nprintf.NPrintf(`
  interface_count       = %{device-interface_count}`, ctx)
	}
	if _, ok := ctx["acl-resourceName"]; ok {
		config += nprintf.NPrintf(`
  acl_template_id       = equinix_network_acl_template.%{acl-resourceName}.id`, ctx)
	}
	if _, ok := ctx["mgmtAcl-resourceName"]; ok {
		config += nprintf.NPrintf(`
  mgmt_acl_template_uuid = equinix_network_acl_template.%{mgmtAcl-resourceName}.id`, ctx)
	}
	if _, ok := ctx["sshkey-resourceName"]; ok {
		config += nprintf.NPrintf(`
  ssh_key {
    username = "test"
    key_name = equinix_network_ssh_key.%{sshkey-resourceName}.name
  }`, ctx)
	}
	if _, ok := ctx["device-license_file"]; ok {
		config += nprintf.NPrintf(`
  license_file          = "%{device-license_file}"`, ctx)
	}
	if _, ok := ctx["device-vendorConfig_enabled"]; ok {
		config += nprintf.NPrintf(`
  vendor_configuration  = {`, ctx)
		if _, ok := ctx["device-vendorConfig_siteId"]; ok {
			config += nprintf.NPrintf(`
    siteId          = "%{device-vendorConfig_siteId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_systemIpAddress"]; ok {
			config += nprintf.NPrintf(`
    systemIpAddress = "%{device-vendorConfig_systemIpAddress}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_licenseKey"]; ok {
			config += nprintf.NPrintf(`
    licenseKey = "%{device-vendorConfig_licenseKey}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_licenseSecret"]; ok {
			config += nprintf.NPrintf(`
    licenseSecret = "%{device-vendorConfig_licenseSecret}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_controller1"]; ok {
			config += nprintf.NPrintf(`
    controller1 = "%{device-vendorConfig_controller1}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_controller2"]; ok {
			config += nprintf.NPrintf(`
    controller2 = "%{device-vendorConfig_controller2}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_localId"]; ok {
			config += nprintf.NPrintf(`
    localId = "%{device-vendorConfig_localId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_remoteId"]; ok {
			config += nprintf.NPrintf(`
    remoteId = "%{device-vendorConfig_remoteId}"`, ctx)
		}
		if _, ok := ctx["device-vendorConfig_serialNumber"]; ok {
			config += nprintf.NPrintf(`
    serialNumber = "%{device-vendorConfig_serialNumber}"`, ctx)
		}
		config += nprintf.NPrintf(`
  }`, ctx)
	}
	if _, ok := ctx["device-secondary_name"]; ok {
		config += nprintf.NPrintf(`
  secondary_device {
    name                 = "%{device-secondary_name}"`, ctx)
		if _, ok := ctx["device-secondary_metro_code"]; ok {
			config += nprintf.NPrintf(`
    metro_code           = "%{device-secondary_metro_code}"
    account_number       = data.equinix_network_account.test-secondary.number`, ctx)
		} else {
			config += nprintf.NPrintf(`
    metro_code           = "%{device-metro_code}"
    account_number       = data.equinix_network_account.test.number`, ctx)
		}
		config += nprintf.NPrintf(`
    notifications        = %{device-secondary_notifications}`, ctx)
		if _, ok := ctx["device-secondary_additional_bandwidth"]; ok {
			config += nprintf.NPrintf(`
    additional_bandwidth = %{device-secondary_additional_bandwidth}`, ctx)
		}
		if _, ok := ctx["device-secondary_hostname"]; ok {
			config += nprintf.NPrintf(`
    hostname             = "%{device-secondary_hostname}"`, ctx)
		}
		if _, ok := ctx["acl-secondary_resourceName"]; ok {
			config += nprintf.NPrintf(`
    acl_template_id      = equinix_network_acl_template.%{acl-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["mgmtAcl-secondary_resourceName"]; ok {
			config += nprintf.NPrintf(`
    mgmt_acl_template_uuid = equinix_network_acl_template.%{mgmtAcl-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["sshkey-resourceName"]; ok {
			config += nprintf.NPrintf(`
    ssh_key {
      username = "test"
      key_name = equinix_network_ssh_key.%{sshkey-resourceName}.name
    }`, ctx)
		}
		if _, ok := ctx["device-secondary_license_file"]; ok {
			config += nprintf.NPrintf(`
    license_file         = "%{device-secondary_license_file}"`, ctx)
		}
		if _, ok := ctx["device-secondary_vendorConfig_enabled"]; ok {
			config += nprintf.NPrintf(`
    vendor_configuration  = {`, ctx)
			if _, ok := ctx["device-secondary_vendorConfig_siteId"]; ok {
				config += nprintf.NPrintf(`
      siteId          = "%{device-secondary_vendorConfig_siteId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_systemIpAddress"]; ok {
				config += nprintf.NPrintf(`
      systemIpAddress = "%{device-secondary_vendorConfig_systemIpAddress}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_licenseKey"]; ok {
				config += nprintf.NPrintf(`
      licenseKey = "%{device-secondary_vendorConfig_licenseKey}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_licenseSecret"]; ok {
				config += nprintf.NPrintf(`
      licenseSecret = "%{device-secondary_vendorConfig_licenseSecret}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_controller1"]; ok {
				config += nprintf.NPrintf(`
      controller1 = "%{device-secondary_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_controller2"]; ok {
				config += nprintf.NPrintf(`
      controller2 = "%{device-secondary_vendorConfig_controller2}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_localId"]; ok {
				config += nprintf.NPrintf(`
      localId = "%{device-secondary_vendorConfig_localId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_remoteId"]; ok {
				config += nprintf.NPrintf(`
      remoteId = "%{device-secondary_vendorConfig_remoteId}"`, ctx)
			}
			if _, ok := ctx["device-secondary_vendorConfig_serialNumber"]; ok {
				config += nprintf.NPrintf(`
      serialNumber = "%{device-secondary_vendorConfig_serialNumber}"`, ctx)
			}
			config += nprintf.NPrintf(`
    }`, ctx)
		}
		config += `
  }`
	}
	if _, ok := ctx["device-cluster_name"]; ok {
		config += nprintf.NPrintf(`
  cluster_details {
    cluster_name        = "%{device-cluster_name}"`, ctx)
		config += `
    node0 {`
		if _, ok := ctx["device-node0_license_file_id"]; ok {
			config += nprintf.NPrintf(`
      license_file_id   = "%{device-node0_license_file_id}"`, ctx)
		}
		if _, ok := ctx["device-node0_license_token"]; ok {
			config += nprintf.NPrintf(`
      license_token     = "%{device-node0_license_token}"`, ctx)
		}
		if _, ok := ctx["device-node0_vendorConfig_enabled"]; ok {
			config += nprintf.NPrintf(`
      vendor_configuration {`, ctx)
			if _, ok := ctx["device-node0_vendorConfig_hostname"]; ok {
				config += nprintf.NPrintf(`
        hostname        = "%{device-node0_vendorConfig_hostname}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_adminPassword"]; ok {
				config += nprintf.NPrintf(`
        admin_password  = "%{device-node0_vendorConfig_adminPassword}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_controller1"]; ok {
				config += nprintf.NPrintf(`
        controller1     = "%{device-node0_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_activationKey"]; ok {
				config += nprintf.NPrintf(`
        activation_key  = "%{device-node0_vendorConfig_activationKey}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_controllerFqdn"]; ok {
				config += nprintf.NPrintf(`
        controller_fqdn = "%{device-node0_vendorConfig_controllerFqdn}"`, ctx)
			}
			if _, ok := ctx["device-node0_vendorConfig_rootPassword"]; ok {
				config += nprintf.NPrintf(`
        root_password   = "%{device-node0_vendorConfig_rootPassword}"`, ctx)
			}
			config += nprintf.NPrintf(`
      }`, ctx)
		}
		config += `
    }`
		config += `
    node1 {`
		if _, ok := ctx["device-node1_license_file_id"]; ok {
			config += nprintf.NPrintf(`
      license_file_id   = "%{device-node1_license_file_id}"`, ctx)
		}
		if _, ok := ctx["device-node1_license_token"]; ok {
			config += nprintf.NPrintf(`
      license_token     = "%{device-node1_license_token}"`, ctx)
		}
		if _, ok := ctx["device-node1_vendorConfig_enabled"]; ok {
			config += nprintf.NPrintf(`
      vendor_configuration {`, ctx)
			if _, ok := ctx["device-node1_vendorConfig_hostname"]; ok {
				config += nprintf.NPrintf(`
        hostname        = "%{device-node1_vendorConfig_hostname}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_adminPassword"]; ok {
				config += nprintf.NPrintf(`
        admin_password  = "%{device-node1_vendorConfig_adminPassword}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_controller1"]; ok {
				config += nprintf.NPrintf(`
        controller1     = "%{device-node1_vendorConfig_controller1}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_activationKey"]; ok {
				config += nprintf.NPrintf(`
        activation_key  = "%{device-node1_vendorConfig_activationKey}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_controllerFqdn"]; ok {
				config += nprintf.NPrintf(`
        controller_fqdn = "%{device-node1_vendorConfig_controllerFqdn}"`, ctx)
			}
			if _, ok := ctx["device-node1_vendorConfig_rootPassword"]; ok {
				config += nprintf.NPrintf(`
        root_password   = "%{device-node1_vendorConfig_rootPassword}"`, ctx)
			}
			config += nprintf.NPrintf(`
      }`, ctx)
		}
		config += `
    }`
		config += `
  }`
	}
	config += `
}`
	return config
}

func testAccNetworkDeviceACL(ctx map[string]interface{}) string {
	var config string
	if _, ok := ctx["acl-name"]; ok {
		config += nprintf.NPrintf(`
resource "equinix_network_acl_template" "%{acl-resourceName}" {
  name          = "%{acl-name}"
  description   = "%{acl-description}"
  inbound_rule {
    subnet   = "10.0.0.0/24"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	}
	if _, ok := ctx["mgmtAcl-name"]; ok {
		config += nprintf.NPrintf(`
resource "equinix_network_acl_template" "%{mgmtAcl-resourceName}" {
  name          = "%{mgmtAcl-name}"
  description   = "%{mgmtAcl-description}"
  inbound_rule {
    subnet   = "11.0.0.0/24"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	}
	if _, ok := ctx["acl-secondary_name"]; ok {
		config += nprintf.NPrintf(`
resource "equinix_network_acl_template" "%{acl-secondary_resourceName}" {
  name          = "%{acl-secondary_name}"
  description   = "%{acl-secondary_description}"
  inbound_rule {
    subnet   = "192.0.0.0/24"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	}
	if _, ok := ctx["mgmtAcl-secondary_name"]; ok {
		config += nprintf.NPrintf(`
resource "equinix_network_acl_template" "%{mgmtAcl-secondary_resourceName}" {
  name          = "%{mgmtAcl-secondary_name}"
  description   = "%{mgmtAcl-secondary_description}"
  inbound_rule {
    subnet   = "193.0.0.0/24"
    protocol = "IP"
    src_port = "any"
    dst_port = "any"
  }
}`, ctx)
	}
	return config
}

func testAccNetworkDeviceSSHKey(ctx map[string]interface{}) string {
	return nprintf.NPrintf(`
resource "equinix_network_ssh_key" "%{sshkey-resourceName}" {
  name       = "%{sshkey-name}"
  public_key = "%{sshkey-public_key}"
}
`, ctx)
}
