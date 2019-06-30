package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourcePacketDevice_Basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePacketDeviceConfig_Basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.packet_device.test", "hostname", "test-device"),
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "id",
						"data.packet_device.test", "id"),
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
  hostname         = "test-device"
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

/*
func TestAccDataSourcePacketDevice_Nonexistent(t *testing.T) {
	projectName := fmt.Sprintf("ds-no-device-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourcePacketDeviceConfig_Nonexistent(projectName),
				ExpectError: regexp.MustCompile("no device found with hostname.*"),
			},
		},
	})
}
*/
func testDataSourcePacketDeviceConfig_Nonexistent() string {
	return fmt.Sprintf(`
data "packet_device" "test" {
  project_id       = packet_project.test.id
  hostname         = "no-such-device"
}`)
}
