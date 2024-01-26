package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalDevice_basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalDeviceConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_device.test", "hostname", "tfacc-test-device"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.test", "id",
						"data.equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.test", "operating_system",
						"data.equinix_metal_device.test", "operating_system"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_device.test", "always_pxe", "false"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_device.test", "access_public_ipv4"),
				),
			},
		},
	})
}

func testDataSourceMetalDeviceConfig_basic(projSuffix string) string {
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
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, testDeviceTerminationTime())
}

func TestAccDataSourceMetalDevice_byID(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-by-id-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalDeviceConfig_byID(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_device.test", "hostname", "tfacc-test-device"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.test", "id",
						"data.equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.test", "operating_system",
						"data.equinix_metal_device.test", "operating_system"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_device.test", "always_pxe", "false"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_device.test", "access_public_ipv4"),
				),
			},
		},
	})
}

func testDataSourceMetalDeviceConfig_byID(projSuffix string) string {
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
  device_id       = equinix_metal_device.test.id
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, testDeviceTerminationTime())
}
