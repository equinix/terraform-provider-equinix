package port

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/tfacc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalPort_byName(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { tfacc.PreCheck(t) },
		Providers: tfacc.AccProviders,
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
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-port-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"

  lifecycle {
    ignore_changes = [
      plan,
      metro,
    ]
  }
}

data "equinix_metal_port" "test" {
    device_id = equinix_metal_device.test.id
    name      = "eth0"
}

`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func TestAccDataSourceMetalPort_byId(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { tfacc.PreCheck(t) },
		Providers: tfacc.AccProviders,
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
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-port-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-port"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"

  lifecycle {
    ignore_changes = [
      plan,
      metro,
    ]
  }
}

data "equinix_metal_port" "test" {
  port_id        = equinix_metal_device.test.ports[0].id
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}
