package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccOrgDataSource_Basic(t *testing.T) {
	var org packngo.Organization
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketOrgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketOrgDataSourceConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketOrgExists("packet_organization.test", &org),
					resource.TestCheckResourceAttr(
						"packet_organization.test", "name",
						fmt.Sprintf("tfacc-datasource-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"packet_organization.test", "description", "quux"),
					resource.TestCheckResourceAttr(
						"data.packet_organization.test", "name",
						fmt.Sprintf("tfacc-datasource-org-%d", rInt)),
					resource.TestCheckResourceAttrPair(
						"data.packet_organization.test2", "id", "packet_organization.test", "id"),
				),
			},
		},
	})
}

func testAccCheckPacketOrgDataSourceConfigBasic(r int) string {
	return fmt.Sprintf(`
resource "packet_organization" "test" {
		name = "tfacc-datasource-org-%d"
		description = "quux"
}

data "packet_organization" "test" {
    organization_id = packet_organization.test.id
}

data "packet_organization" "test2" {
    name = "${packet_organization.test.name}"
}

`, r)
}
