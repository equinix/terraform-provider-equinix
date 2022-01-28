package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetalIPBlockRanges_Basic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIPBlockRangesConfig_Basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.metal_ip_block_ranges.test", "ipv6.0"),
					resource.TestCheckResourceAttrPair(
						"metal_ip_attachment.test", "device_id",
						"metal_device.test", "id"),
				),
			},
			{
				ResourceName:      "metal_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testIPBlockRangesConfig_Basic(name string) string {
	return fmt.Sprintf(`

resource "metal_project" "test" {
    name = "tfacc-precreated_ip_block-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-ip-test"
  plan             = "t1.small.x86"
  facilities       = ["ewr1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = metal_project.test.id
}

data "metal_ip_block_ranges" "test" {
    facility         = "ewr1"
    project_id       = metal_device.test.project_id
}

resource "metal_ip_attachment" "test" {
    device_id = metal_device.test.id
    cidr_notation = cidrsubnet(data.metal_ip_block_ranges.test.ipv6.0, 8,2)
}
`, name)
}
