package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMetalPort_ByName(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPortConfig_ByName(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_port.test", "bond_port", "false"),
				),
			},
		},
	})
}

func testPortConfig_ByName(name string) string {
	return fmt.Sprintf(`

resource "metal_project" "test" {
    name = "tfacc-port-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = "c3.medium.x86"
  facilities       = ["sv15"]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = metal_project.test.id
}

data "metal_port" "test" {
    device_id        = metal_device.test.id
    port_name        = "eth0"
}

`, name)
}

func TestAccMetalPort_ById(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPortConfig_ById(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.metal_port.test", "port_name"),
				),
			},
		},
	})
}

func testPortConfig_ById(name string) string {
	return fmt.Sprintf(`

resource "metal_project" "test" {
    name = "tfacc-port-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = "c3.medium.x86"
  facilities       = ["sv15"]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = metal_project.test.id
}

data "metal_port" "test" {
  id        = metal_device.test.ports[0].id
}

`, name)
}
