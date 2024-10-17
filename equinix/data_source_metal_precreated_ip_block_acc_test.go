package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalPreCreatedIPBlock_basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalPreCreatedIPBlockConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_precreated_ip_block.test_fac_pubv6", "cidr_notation"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_precreated_ip_block.test_metro_priv4", "cidr_notation"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_ip_attachment.test", "device_id",
						"equinix_metal_device.test", "id"),
				),
			},
			{
				ResourceName:      "equinix_metal_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataSourceMetalPreCreatedIPBlockConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-precreated_ip_block-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-ip-blockt"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}

data "equinix_metal_precreated_ip_block" "test_fac_pubv6" {
    facility         = equinix_metal_device.test.deployed_facility
    project_id       = equinix_metal_device.test.project_id
    address_family   = 6
    public           = true
}

data "equinix_metal_precreated_ip_block" "test_metro_priv4" {
    metro            = equinix_metal_device.test.metro
    project_id       = equinix_metal_device.test.project_id
    address_family   = 4
    public           = false
}

resource "equinix_metal_ip_attachment" "test" {
    device_id = equinix_metal_device.test.id
    cidr_notation = cidrsubnet(data.equinix_metal_precreated_ip_block.test_fac_pubv6.cidr_notation,8,2)
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}
