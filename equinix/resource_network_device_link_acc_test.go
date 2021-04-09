package equinix

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("NetworkDeviceLink", &resource.Sweeper{
		Name: "NetworkDeviceLink",
		F:    testSweepNetworkDeviceLink,
	})
}

func testSweepNetworkDeviceLink(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	links, err := config.ne.GetDeviceLinkGroups()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching device links list: %s", err)
		return err
	}
	nonSweepableCount := 0
	for _, link := range links {
		if !isSweepableTestResource(ne.StringValue(link.Name)) {
			nonSweepableCount++
			continue
		}
		if err := config.ne.DeleteDeviceLinkGroup(ne.StringValue(link.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkDeviceLink resource %s (%s): %s", ne.StringValue(link.UUID), ne.StringValue(link.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkDeviceLink resource %s (%s)", ne.StringValue(link.UUID), ne.StringValue(link.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}

func TestAccNetworkDeviceLink(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	metroSecondary, _ := schema.EnvDefaultFunc(networkDeviceSecondaryMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"device-resourceName":               "test",
		"device-self_managed":               false,
		"device-byol":                       false,
		"device-name":                       fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-throughput":                 500,
		"device-throughput_unit":            "Mbps",
		"device-metro_code":                 metro.(string),
		"device-type_code":                  "CSR1000V",
		"device-package_code":               "SEC",
		"device-notifications":              []string{"test@equinix.com"},
		"device-hostname":                   fmt.Sprintf("tf-%s", randString(6)),
		"device-term_length":                1,
		"device-version":                    "16.09.05",
		"device-core_count":                 2,
		"device-secondary_name":             fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-secondary_metro_code":       metroSecondary.(string),
		"device-secondary_hostname":         fmt.Sprintf("tf-%s", randString(6)),
		"device-secondary_notifications":    []string{"test@equinix.com"},
		"link-resourceName":                 "test",
		"link-name":                         fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"link-subnet":                       "10.69.1.0/24",
		"link-device_1_asn":                 23404,
		"link-device_1_interface_id":        6,
		"link-device_2_asn":                 24040,
		"link-device_2_interface_id":        6,
		"link-connection_1_throughput":      "50",
		"link-connection_1_throughput_unit": "Mbps",
	}
	deviceResourceName := "equinix_network_device." + context["device-resourceName"].(string)
	linkResourceName := "equinix_network_device_link." + context["link-resourceName"].(string)
	var deviceLink ne.DeviceLinkGroup
	var primaryDevice, secondaryDevice ne.Device
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withDeviceLink().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeDevicePairExists(deviceResourceName, &primaryDevice, &secondaryDevice),
					testAccNeDeviceLinkExists(linkResourceName, &deviceLink),
					testAccNeDeviceLinkAttributes(&deviceLink, context),
					resource.TestCheckResourceAttrSet(linkResourceName, "uuid"),
					resource.TestCheckResourceAttr(linkResourceName, "status", ne.DeviceLinkGroupStatusProvisioned),
					testAccNeDeviceLinkDeviceConnections(&deviceLink, &primaryDevice, &secondaryDevice, context),
				),
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
	config += nprintf(`
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
  link {
    account_number  = equinix_network_device.%{device-resourceName}.account_number
    throughput      = "%{link-connection_1_throughput}"
    throughput_unit = "%{link-connection_1_throughput_unit}"
    src_metro_code  = equinix_network_device.%{device-resourceName}.metro_code
    dst_metro_code  = equinix_network_device.%{device-resourceName}.secondary_device[0].metro_code
    src_zone_code   = equinix_network_device.%{device-resourceName}.zone_code
    dst_zone_code   = equinix_network_device.%{device-resourceName}.secondary_device[0].zone_code
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
		client := testAccProvider.Meta().(*Config).ne
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
		if ne.StringValue(deviceLink.Links[0].SourceZoneCode) != ne.StringValue(primaryDevice.ZoneCode) {
			return fmt.Errorf("link #1 src_zone_code does not match %v - %v", ne.StringValue(deviceLink.Links[0].SourceZoneCode), ne.StringValue(primaryDevice.ZoneCode))
		}
		if ne.StringValue(deviceLink.Links[0].DestinationMetroCode) != ne.StringValue(secondaryDevice.MetroCode) {
			return fmt.Errorf("link #1 dst_metro_code does not match %v - %v", ne.StringValue(deviceLink.Links[0].DestinationMetroCode), ne.StringValue(secondaryDevice.MetroCode))
		}
		if ne.StringValue(deviceLink.Links[0].DestinationZoneCode) != ne.StringValue(secondaryDevice.ZoneCode) {
			return fmt.Errorf("link #1 dst_zone_code does not match %v - %v", ne.StringValue(deviceLink.Links[0].DestinationZoneCode), ne.StringValue(secondaryDevice.ZoneCode))
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
