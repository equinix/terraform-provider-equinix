package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalIPBlockRanges_basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalIPBlockRangesConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_ip_block_ranges.test", "ipv6.0"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_ip_attachment.test", "device_id",
						"equinix_metal_device.test", "id"),
				),
			},
			{
				ResourceName:      "equinix_metal_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataSourceMetalIPBlockRangesConfig_basic(name string) string {
	return fmt.Sprintf(`

resource "equinix_metal_project" "test" {
    name = "tfacc-precreated_ip_block-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-ip-test"
  plan             = "t1.small.x86"
  facilities       = ["ny5"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
}

data "equinix_metal_ip_block_ranges" "test" {
    facility         = "ny5"
    project_id       = equinix_metal_device.test.project_id
}

resource "equinix_metal_ip_attachment" "test" {
    device_id = equinix_metal_device.test.id
    cidr_notation = cidrsubnet(data.equinix_metal_ip_block_ranges.test.ipv6.0, 8,2)
}
`, name)
}
