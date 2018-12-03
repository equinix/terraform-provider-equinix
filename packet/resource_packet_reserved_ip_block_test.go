package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketReservedIPBlock_Basic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketReservedIPBlockConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "facility", "ewr1"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "quantity", "2"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "netmask", "255.255.255.254"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "public", "true"),
					resource.TestCheckResourceAttr(
						"packet_reserved_ip_block.test", "management", "false"),
				),
			},
		},
	})
}

func TestAccPacketReservedIPBlock_importBasic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketReservedIPBlockDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketReservedIPBlockConfig_basic(rs),
			},
			resource.TestStep{
				ResourceName:      "packet_reserved_ip_block.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPacketReservedIPBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_reserved_ip_block" {
			continue
		}
		if _, _, err := client.ProjectIPs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Reserved IP block still exists")
		}
	}

	return nil
}

func testAccCheckPacketReservedIPBlockConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "%s"
}

resource "packet_reserved_ip_block" "test" {
    project_id = "${packet_project.foobar.id}"
    facility = "ewr1"
	quantity = 2
}`, name)
}
