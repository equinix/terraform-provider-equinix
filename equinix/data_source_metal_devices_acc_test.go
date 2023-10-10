package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalDevices(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalDevicesConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_devices.test_filter_tags", "devices.#", "1"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_devices.test_search", "devices.#", "1"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.dev_tags", "id",
						"data.equinix_metal_devices.test_filter_tags", "devices.0.device_id"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.dev_search", "id",
						"data.equinix_metal_devices.test_search", "devices.0.device_id"),
				),
			},
		},
	})
}

func testDataSourceMetalDevicesConfig_basic(projSuffix string) string {
	return fmt.Sprintf(`
%[1]s

resource "equinix_metal_project" "test" {
    name = "tfacc-project-%[2]s"
}

resource "equinix_metal_device" "dev_tags" {
  hostname         = "tfacc-test-device1"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  termination_time = "%[3]s"
  tags             = ["tag1", "tag2"]
}

resource "equinix_metal_device" "dev_search" {
  hostname         = "tfacc-test-device2-unlikelystring"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  termination_time = "%[3]s"
}

data "equinix_metal_devices" "test_filter_tags" {
  project_id = equinix_metal_project.test.id
  filter {
	attribute = "tags"
	values    = ["tag1"]
  }
  depends_on = [equinix_metal_device.dev_tags]
}

data "equinix_metal_devices" "test_search" {
  project_id = equinix_metal_project.test.id
  search     = "unlikelystring"
  depends_on = [equinix_metal_device.dev_search]
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, testDeviceTerminationTime())
}
