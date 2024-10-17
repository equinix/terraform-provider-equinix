package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalIPBlockRanges_basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-precreated_ip_block-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-ip-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}

data "equinix_metal_ip_block_ranges" "test" {
    facility   = equinix_metal_device.test.deployed_facility
    project_id = equinix_metal_device.test.project_id
}

resource "equinix_metal_ip_attachment" "test" {
    device_id = equinix_metal_device.test.id
    cidr_notation = cidrsubnet(data.equinix_metal_ip_block_ranges.test.ipv6.0, 8, 2)
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}
