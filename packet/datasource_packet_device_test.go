package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourcePacketDevice_Basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePacketDeviceConfig_Basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.packet_device.test", "hostname", "tfacc-test-device"),
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "id",
						"data.packet_device.test", "id"),
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "operating_system",
						"data.packet_device.test", "operating_system"),
					resource.TestCheckResourceAttr(
						"data.packet_device.test", "always_pxe", "false"),
					resource.TestCheckResourceAttrSet(
						"data.packet_device.test", "access_public_ipv4"),
				),
			},
		},
	})
}

func testDataSourcePacketDeviceConfig_Basic(projSuffix string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-project-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}

data "packet_device" "test" {
  project_id       = packet_project.test.id
  hostname         = packet_device.test.hostname
}`, projSuffix)
}

func TestAccDataSourcePacketDevice_ByID(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-by-id-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePacketDeviceConfig_ByID(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.packet_device.test", "hostname", "tfacc-test-device"),
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "id",
						"data.packet_device.test", "id"),
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "operating_system",
						"data.packet_device.test", "operating_system"),
					resource.TestCheckResourceAttr(
						"data.packet_device.test", "always_pxe", "false"),
					resource.TestCheckResourceAttrSet(
						"data.packet_device.test", "access_public_ipv4"),
				),
			},
		},
	})
}

func testDataSourcePacketDeviceConfig_ByID(projSuffix string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-project-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}

data "packet_device" "test" {
  device_id       = packet_device.test.id
}`, projSuffix)
}
