package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalDevice_Basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalDeviceConfig_Basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_device.test", "hostname", "tfacc-test-device"),
					resource.TestCheckResourceAttrPair(
						"metal_device.test", "id",
						"data.metal_device.test", "id"),
					resource.TestCheckResourceAttrPair(
						"metal_device.test", "operating_system",
						"data.metal_device.test", "operating_system"),
					resource.TestCheckResourceAttr(
						"data.metal_device.test", "always_pxe", "false"),
					resource.TestCheckResourceAttrSet(
						"data.metal_device.test", "access_public_ipv4"),
				),
			},
		},
	})
}

func testDataSourceMetalDeviceConfig_Basic(projSuffix string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-project-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
}

data "metal_device" "test" {
  project_id       = metal_project.test.id
  hostname         = metal_device.test.hostname
}`, projSuffix)
}

func TestAccDataSourceMetalDevice_ByID(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-by-id-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalDeviceConfig_ByID(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_device.test", "hostname", "tfacc-test-device"),
					resource.TestCheckResourceAttrPair(
						"metal_device.test", "id",
						"data.metal_device.test", "id"),
					resource.TestCheckResourceAttrPair(
						"metal_device.test", "operating_system",
						"data.metal_device.test", "operating_system"),
					resource.TestCheckResourceAttr(
						"data.metal_device.test", "always_pxe", "false"),
					resource.TestCheckResourceAttrSet(
						"data.metal_device.test", "access_public_ipv4"),
				),
			},
		},
	})
}

func testDataSourceMetalDeviceConfig_ByID(projSuffix string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-project-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
}

data "metal_device" "test" {
  device_id       = metal_device.test.id
}`, projSuffix)
}
