package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalPort_byName(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalPortConfig_byName(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_port.test", "bond_name", "bond0"),
				),
			},
		},
	})
}

func testAccDataSourceMetalPortConfig_byName(name string) string {
	return fmt.Sprintf(`

resource "equinix_metal_project" "test" {
    name = "tfacc-port-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = "c3.medium.x86"
  metro            = "sv"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
}

data "equinix_metal_port" "test" {
    device_id = equinix_metal_device.test.id
    name      = "eth0"
}

`, name)
}

func TestAccDataSourceMetalPort_byId(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalPortConfig_byId(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_port.test", "name"),
				),
			},
		},
	})
}

func testAccDataSourceMetalPortConfig_byId(name string) string {
	return fmt.Sprintf(`

resource "equinix_metal_project" "test" {
    name = "tfacc-port-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = "c3.medium.x86"
  metro            = "sv"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
}

data "equinix_metal_port" "test" {
  port_id        = equinix_metal_device.test.ports[0].id
}
`, name)
}
