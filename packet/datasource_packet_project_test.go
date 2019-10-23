package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/packethost/packngo"
)

func testAccCheckPacketDataSourceProject_Basic(r string) string {
	return fmt.Sprintf(`
resource "packet_organization" "test" {
	name = "tfacc-organization-%s"
}

resource "packet_project" "foobar" {
	name = "tfacc-project-%s"
	organization_id = "${packet_organization.test.id}"
	bgp_config {
		deployment_type = "local"
		md5 = "2SFsdfsg43"
		asn = 65000
	}
}

data packet_project "test" {
	project_id = packet_project.foobar.id
}

data packet_project "test2" {
	name= packet_project.foobar.name
}

`, r, r)
}

func TestAccPacketDataSourceProject_Basic(t *testing.T) {
	var project packngo.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketDataSourceProject_Basic(rn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "name", fmt.Sprintf("tfacc-project-%s", rn)),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "bgp_config.0.md5",
						"2SFsdfsg43"),
					resource.TestCheckResourceAttrPair(
						"packet_project.foobar", "id",
						"data.packet_project.test", "id"),
					resource.TestCheckResourceAttrPair(
						"packet_project.foobar", "name",
						"data.packet_project.test2", "name"),
				),
			},
		},
	})
}
