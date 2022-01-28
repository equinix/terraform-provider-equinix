package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetalDatasourcePreCreatedIPBlock_Basic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDatasourcePreCreatedIPBlockConfig_Basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.metal_precreated_ip_block.test", "cidr_notation"),
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

func testDatasourcePreCreatedIPBlockConfig_Basic(name string) string {
	return fmt.Sprintf(`

resource "equinix_metal_project" "test" {
    name = "tfacc-precreated_ip_block-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-ip-blockt"
  plan             = "c3.small.x86"
  metro            = "ny"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = metal_project.test.id
}

data "equinix_metal_precreated_ip_block" "test" {
    facility         = metal_device.test.deployed_facility
    project_id       = metal_device.test.project_id
    address_family   = 6
    public           = true
}

resource "equinix_metal_ip_attachment" "test" {
    device_id = metal_device.test.id
    cidr_notation = cidrsubnet(data.metal_precreated_ip_block.test.cidr_notation,8,2)
}
`, name)
}
