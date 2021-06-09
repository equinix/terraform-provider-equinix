package metal

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
					testAccCheckMetalOrgExists("metal_organization.test", &org),
					resource.TestCheckResourceAttr(
						"metal_organization.test", "name",
						fmt.Sprintf("tfacc-datasource-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"metal_organization.test", "description", "quux"),
					resource.TestCheckResourceAttr(
						"data.metal_organization.test", "name",
						fmt.Sprintf("tfacc-datasource-org-%d", rInt)),
					resource.TestCheckResourceAttrPair(
						"data.metal_organization.test2", "id", "metal_organization.test", "id"),
				),
			},
		},
	})
}

func testAccCheckMetalOrgDataSourceConfigBasic(r int) string {
	return fmt.Sprintf(`
resource "metal_organization" "test" {
		name = "tfacc-datasource-org-%d"
		description = "quux"
}

data "metal_organization" "test" {
    organization_id = metal_organization.test.id
}

data "metal_organization" "test2" {
    name = "${metal_organization.test.name}"
}

`, r)
}
