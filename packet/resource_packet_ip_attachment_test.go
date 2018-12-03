package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketIPAttachment_Basic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketIPAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketIPAttachmentConfig_Basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_ip_attachment.test", "public", "true"),
					resource.TestCheckResourceAttrPair(
						"packet_ip_attachment.test", "device_id",
						"packet_device.test", "id"),
				),
			},
			resource.TestStep{
				ResourceName:      "packet_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPacketIPAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_ip_attachment" {
			continue
		}
		if _, _, err := client.ProjectIPs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("IP attachment still exists")
		}
	}

	return nil
}

func testAccCheckPacketIPAttachmentConfig_Basic(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "%s"
}

resource "packet_device" "test" {
  hostname         = "test"
  plan             = "baremetal_0"
  facility         = "ewr1"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}

resource "packet_reserved_ip_block" "test" {
    project_id = "${packet_project.test.id}"
    facility = "ewr1"
	quantity = 2
}


resource "packet_ip_attachment" "test" {
	device_id = "${packet_device.test.id}"
	cidr_notation = "${cidrhost(packet_reserved_ip_block.test.cidr_notation,0)}/32"
}`, name)
}
