package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func TestAccMetalIPAttachment_Basic(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalIPAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalIPAttachmentConfig_Basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_ip_attachment.test", "public", "true"),
					resource.TestCheckResourceAttrPair(
						"metal_ip_attachment.test", "device_id",
						"metal_device.test", "id"),
				),
			},
			{
				ResourceName:      "metal_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMetalIPAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_ip_attachment" {
			continue
		}
		if _, _, err := client.ProjectIPs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("IP attachment still exists")
		}
	}

	return nil
}

func testAccCheckMetalIPAttachmentConfig_Basic(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-ip_attachment-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-ip-attachment-test"
  plan             = "t1.small.x86"
  facilities       = ["ewr1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
}

resource "metal_reserved_ip_block" "test" {
    project_id = "${metal_project.test.id}"
    facility = "ewr1"
	quantity = 2
}


resource "metal_ip_attachment" "test" {
	device_id = "${metal_device.test.id}"
	cidr_notation = "${cidrhost(metal_reserved_ip_block.test.cidr_notation,0)}/32"
}`, name)
}
