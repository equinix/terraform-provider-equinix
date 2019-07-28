package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccOrgDataSource_Basic(t *testing.T) {
	var org packngo.Organization

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketOrgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketOrgDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketOrgExists("packet_organization.test", &org),
					testAccCheckPacketOrgAttributes(&org),
					resource.TestCheckResourceAttr(
						"packet_organization.test", "name", "foobar"),
					resource.TestCheckResourceAttr(
						"packet_organization.test", "description", "quux"),
					resource.TestCheckResourceAttr(
						"data.packet_organization.test", "name", "foobar"),
					resource.TestCheckResourceAttrPair(
						"data.packet_organization.test2", "id", "packet_organization.test", "id"),
				),
			},
		},
	})
}

var testAccCheckPacketOrgDataSourceConfigBasic = fmt.Sprintf(`
resource "packet_organization" "test" {
		name = "foobar"
		description = "quux"
}

data "packet_organization" "test" {
    organization_id = packet_organization.test.id
}

data "packet_organization" "test2" {
    name = packet_organization.test.name
}

`)
