package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalDevice_basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalDeviceCheckDestroyed,
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
resource "equinix_metal_project" "test" {
    name = "tfacc-project-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
}

data "equinix_metal_device" "test" {
  project_id       = equinix_metal_project.test.id
  hostname         = equinix_metal_device.test.hostname
}`, projSuffix)
}

func TestAccDataSourceMetalDevice_byID(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-by-id-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalDeviceCheckDestroyed,
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
resource "equinix_metal_project" "test" {
    name = "tfacc-project-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
}

data "equinix_metal_device" "test" {
  device_id       = equinix_metal_device.test.id
}`, projSuffix)
}
