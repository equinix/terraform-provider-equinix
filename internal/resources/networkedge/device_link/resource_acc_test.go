package device_link

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccNeDevicePairExists(resourceName string, primary, secondary *ne.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		client := acceptance.TestAccProvider.Meta().(*config.Config).Ne
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

func TestAccNetworkDeviceLink(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	metroSecondary, _ := schema.EnvDefaultFunc(networkDeviceSecondaryMetroEnvVar, metro)()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	accountNameSecondary, _ := schema.EnvDefaultFunc(networkDeviceSecondaryAccountNameEnvVar, accountName)()
	context := map[string]interface{}{
		"device-resourceName":               "test",
		"device-account_name":               accountName.(string),
		"device-self_managed":               false,
		"device-byol":                       false,
		"device-name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-throughput":                 500,
		"device-throughput_unit":            "Mbps",
		"device-metro_code":                 metro.(string),
		"device-type_code":                  "CSR1000V",
		"device-package_code":               "SEC",
		"device-notifications":              []string{"test@equinix.com"},
		"device-hostname":                   fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":                1,
		"device-version":                    "16.09.05",
		"device-core_count":                 2,
		"device-secondary_name":             fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-secondary_account_name":     accountNameSecondary.(string),
		"device-secondary_metro_code":       metroSecondary.(string),
		"device-secondary_hostname":         fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-secondary_notifications":    []string{"test@equinix.com"},
		"link-resourceName":                 "test",
		"link-name":                         fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"link-subnet":                       "10.69.1.0/24",
		"link-device_1_asn":                 23404,
		"link-device_1_interface_id":        6,
		"link-device_2_asn":                 24040,
		"link-device_2_interface_id":        6,
		"link-connection_1_throughput":      "50",
		"link-connection_1_throughput_unit": "Mbps",
		"metro-link_1_account_number":       "1234",
		"metro-link_1_metro_code":           metro.(string),
		"metro-link_1_throughput":           10,
		"metro-link_1_throughput_unit":      "Mbps",
		"metro-link_2_account_number":       "1432",
		"metro-link_2_metro_code":           metroSecondary.(string),
		"metro-link_2_throughput":           10,
		"metro-link_2_throughput_unit":      "Mbps",
		"link_redundancy-type":              "PRIMARY",
	}
	deviceResourceName := "equinix_network_device." + context["device-resourceName"].(string)
	linkResourceName := "equinix_network_device_link." + context["link-resourceName"].(string)
	var deviceLink ne.DeviceLinkGroup
	var primaryDevice, secondaryDevice ne.Device
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withDeviceLink().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDeviceLinkExists(linkResourceName, &deviceLink),
					testAccNeDeviceLinkAttributes(&deviceLink, context),
					resource.TestCheckResourceAttrSet(linkResourceName, "uuid"),
					resource.TestCheckResourceAttr(linkResourceName, "status", ne.DeviceLinkGroupStatusProvisioned),
					testAccNeDevicePairExists(deviceResourceName, &primaryDevice, &secondaryDevice),
					testAccNeDeviceLinkDeviceConnections(&deviceLink, &primaryDevice, &secondaryDevice, context),
				),
			},
			{
				ResourceName:      linkResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func (t *testAccConfig) withDeviceLink() *testAccConfig {
	t.config += testAccNetworkDeviceLink(t.ctx)
	return t
}

func testAccNetworkDeviceLink(ctx map[string]interface{}) string {
	var config string
	config += nprintf.Nprintf(`
resource "equinix_network_device_link" "%{link-resourceName}" {
  name   = "%{link-name}"
  subnet = "%{link-subnet}"
  device {
    id           = equinix_network_device.%{device-resourceName}.id
    asn          = %{link-device_1_asn}
    interface_id = %{link-device_1_interface_id}
  }
  device {
    id           = equinix_network_device.%{device-resourceName}.secondary_device[0].uuid
    asn          = %{link-device_2_asn}
    interface_id = %{link-device_2_interface_id}
  }
  # link block not required if metro_code is the same for both devices
  dynamic "link" {
    for_each = equinix_network_device.%{device-resourceName}.metro_code == equinix_network_device.%{device-resourceName}.secondary_device[0].metro_code ? [] : [1]
    content {
	  account_number  = equinix_network_device.%{device-resourceName}.account_number
	  throughput      = "%{link-connection_1_throughput}"
	  throughput_unit = "%{link-connection_1_throughput_unit}"
	  src_metro_code  = equinix_network_device.%{device-resourceName}.metro_code
	  dst_metro_code  = equinix_network_device.%{device-resourceName}.secondary_device[0].metro_code
    }
  }
  "metro_link" {
	  account_number  = "%{metro-link_1_account_number}"
	  throughput      = "%{metro-link_1_throughput}"
	  throughput_unit = "%{metro-link_1_throughput_unit}"
	  metro_code      = "%{metro-link_1_metro_code}"
  }
  "metro_link" {
	account_number  = "%{metro-link_2_account_number}"
	throughput      = "%{metro-link_2_throughput}"
	throughput_unit = "%{metro-link_2_throughput_unit}"
	metro_code      = "%{metro-link_2_metro_code}"
}
}`, ctx)
	return config
}

func testAccNeDeviceLinkExists(resourceName string, deviceLink *ne.DeviceLinkGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		client := acceptance.TestAccProvider.Meta().(*config.Config).Ne
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		resp, err := client.GetDeviceLinkGroup(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching device link '%s': %s", rs.Primary.ID, err)
		}
		*deviceLink = *resp
		return nil
	}
}

func testAccNeDeviceLinkAttributes(deviceLink *ne.DeviceLinkGroup, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["link-name"]; ok && ne.StringValue(deviceLink.Name) != v.(string) {
			return fmt.Errorf("name does not match %v - %v", ne.StringValue(deviceLink.Name), v)
		}
		if v, ok := ctx["link-subnet"]; ok && ne.StringValue(deviceLink.Subnet) != v.(string) {
			return fmt.Errorf("subnet does not match %v - %v", ne.StringValue(deviceLink.Subnet), v)
		}
		return nil
	}
}

func testAccNeDeviceLinkDeviceConnections(deviceLink *ne.DeviceLinkGroup, primaryDevice, secondaryDevice *ne.Device, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(deviceLink.Devices) != 2 {
			return fmt.Errorf("number of devices does not match %v - %v", len(deviceLink.Devices), 2)
		}
		deviceLinkDeviceMap := make(map[string]*ne.DeviceLinkGroupDevice)
		for i := range deviceLink.Devices {
			deviceLinkDeviceMap[ne.StringValue(deviceLink.Devices[i].DeviceID)] = &deviceLink.Devices[i]
		}
		if _, ok := deviceLinkDeviceMap[ne.StringValue(primaryDevice.UUID)]; ok {
			if err := testAccNeDeviceLinkDeviceAttributes(deviceLink, primaryDevice, "link-device_1", ctx); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("link does not contain primary device %v", ne.StringValue(primaryDevice.UUID))
		}
		if _, ok := deviceLinkDeviceMap[ne.StringValue(secondaryDevice.UUID)]; ok {
			if err := testAccNeDeviceLinkDeviceAttributes(deviceLink, secondaryDevice, "link-device_2", ctx); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("link does not contain secondary device %v", ne.StringValue(secondaryDevice.UUID))
		}
		if ne.StringValue(primaryDevice.MetroCode) == ne.StringValue(secondaryDevice.MetroCode) {
			if len(deviceLink.Links) != 0 {
				return fmt.Errorf("number of links for devices in same metro does not match %v - %v", len(deviceLink.Links), 0)
			}
		} else {
			if len(deviceLink.Links) != 1 {
				return fmt.Errorf("number of links does not match %v - %v", len(deviceLink.Links), 1)
			}
			if v, ok := ctx["link-connection_1_throughput"]; ok && ne.StringValue(deviceLink.Links[0].Throughput) != v.(string) {
				return fmt.Errorf("link #1 throughput does not match %v - %v", ne.StringValue(deviceLink.Links[0].Throughput), v)
			}
			if v, ok := ctx["link-connection_1_throughput_unit"]; ok && ne.StringValue(deviceLink.Links[0].ThroughputUnit) != v.(string) {
				return fmt.Errorf("link #1 throughput_unit does not match %v - %v", ne.StringValue(deviceLink.Links[0].ThroughputUnit), v)
			}
			if ne.StringValue(deviceLink.Links[0].SourceMetroCode) != ne.StringValue(primaryDevice.MetroCode) {
				return fmt.Errorf("link #1 src_metro_code does not match %v - %v", ne.StringValue(deviceLink.Links[0].SourceMetroCode), ne.StringValue(primaryDevice.MetroCode))
			}
			if ne.StringValue(deviceLink.Links[0].DestinationMetroCode) != ne.StringValue(secondaryDevice.MetroCode) {
				return fmt.Errorf("link #1 dst_metro_code does not match %v - %v", ne.StringValue(deviceLink.Links[0].DestinationMetroCode), ne.StringValue(secondaryDevice.MetroCode))
			}
		}
		return nil
	}
}

func testAccNeDeviceLinkDeviceAttributes(deviceLink *ne.DeviceLinkGroup, device *ne.Device, ctxPrefix string, ctx map[string]interface{}) error {
	if v, ok := ctx[ctxPrefix+"_asn"]; ok && ne.IntValue(device.ASN) != v.(int) {
		return fmt.Errorf("device %v ASN does not match %v - %v", ne.StringValue(device.UUID), ne.IntValue(device.ASN), v)
	}
	if v, ok := ctx[ctxPrefix+"_interface_id"]; ok && v.(int) > 0 && v.(int) <= len(device.Interfaces) {
		deviceInterfaceIdx := v.(int) - 1
		deviceInterface := device.Interfaces[deviceInterfaceIdx]
		if ne.StringValue(deviceInterface.AssignedType) != ne.StringValue(deviceLink.Name) {
			return fmt.Errorf("device %v interface #%d assignedType does not match link name %v - %v", ne.StringValue(device.UUID), v, ne.StringValue(deviceInterface.AssignedType), ne.StringValue(deviceLink.Name))
		}
	}
	return nil
}
