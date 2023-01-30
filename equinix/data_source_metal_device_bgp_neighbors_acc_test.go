package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalDeviceBgpNeighbors(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalDeviceBgpNeighborsConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_device_bgp_neighbors.test", "bgp_neighbors.#"),
				),
			},
		},
	})
}

func testAccDataSourceMetalDeviceBgpNeighborsConfig(projectName string) string {
	return fmt.Sprintf(`
%s

data "equinix_metal_device_bgp_neighbors" "test" {
	device_id = equinix_metal_device.test.id
}

output "bgp_neighbors_listing" {
	value = data.equinix_metal_device_bgp_neighbors.test.bgp_neighbors
}
`, testDataSourceMetalDeviceConfig_basic(projectName))
}
