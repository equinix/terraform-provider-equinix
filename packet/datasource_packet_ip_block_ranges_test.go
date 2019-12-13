package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPacketIPBlockRanges_Basic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIPBlockRangesConfig_Basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.packet_ip_block_ranges.test", "ipv6.0"),
					resource.TestCheckResourceAttrPair(
						"packet_ip_attachment.test", "device_id",
						"packet_device.test", "id"),
				),
			},
			{
				ResourceName:      "packet_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testIPBlockRangesConfig_Basic(name string) string {
	return fmt.Sprintf(`

resource "packet_project" "test" {
    name = "tfacc-precreated_ip_block-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-device-ip-test"
  plan             = "t1.small.x86"
  facilities       = ["ewr1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = packet_project.test.id
}

data "packet_ip_block_ranges" "test" {
    facility         = "ewr1"
    project_id       = packet_device.test.project_id
}

resource "packet_ip_attachment" "test" {
    device_id = packet_device.test.id
    cidr_notation = cidrsubnet(data.packet_ip_block_ranges.test.ipv6.0, 8,2)
}
`, name)
}
