package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccOrgDataSource_Basic(t *testing.T) {
	var org packngo.Organization
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalOrgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalOrgDataSourceConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalOrgExists("equinix_metal_organization.test", &org),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "name",
						fmt.Sprintf("tfacc-datasource-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "description", "quux"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_organization.test", "name",
						fmt.Sprintf("tfacc-datasource-org-%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckMetalOrgDataSourceConfigBasic(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_organization" "test" {
  name = "tfacc-datasource-org-%d"
  description = "quux"
}

data "equinix_metal_organization" "test" {
  organization_id = metal_organization.test.id
}


`, r)
}
