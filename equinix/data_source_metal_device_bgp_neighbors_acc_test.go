package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalDeviceBgpNeighbors(t *testing.T) {
	projSuffix := fmt.Sprintf("ds-device-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalDeviceBgpNeighborsConfig(projSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_device_bgp_neighbors.test", "bgp_neighbors.#"),
				),
			},
		},
	})
}

func testAccDataSourceMetalDeviceBgpNeighborsConfig(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-project-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  termination_time = "%s"
}

data "equinix_metal_device" "test" {
  project_id       = equinix_metal_project.test.id
  hostname         = equinix_metal_device.test.hostname
}

data "equinix_metal_device_bgp_neighbors" "test" {
	device_id = equinix_metal_device.test.id
}

output "bgp_neighbors_listing" {
	value = data.equinix_metal_device_bgp_neighbors.test.bgp_neighbors
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, testDeviceTerminationTime())
}
