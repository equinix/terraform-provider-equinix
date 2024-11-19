package port_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalPort_byName(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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
}

data "equinix_metal_port" "test" {
    device_id = equinix_metal_device.test.id
    name      = "eth0"
}

`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), name, acceptance.TestDeviceTerminationTime())
}

func TestAccDataSourceMetalPort_byId(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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
}

data "equinix_metal_port" "test" {
  port_id        = equinix_metal_device.test.ports[0].id
}
`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), name, acceptance.TestDeviceTerminationTime())
}
