package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketVolumeAttachment_Basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketVolumeAttachmentConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"packet_volume_attachment.test", "volume_id",
						"packet_volume.test", "id"),
					resource.TestCheckResourceAttrPair(
						"packet_volume_attachment.test", "device_id",
						"packet_device.test", "id"),
				),
			},
			resource.TestStep{
				ResourceName:      "packet_volume_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPacketVolumeAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_volume_attachment" {
			continue
		}
		if _, _, err := client.VolumeAttachments.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("VolumeAttachment still exists")
		}
	}

	return nil
}

func testAccCheckPacketVolumeAttachmentConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "%s"
}

resource "packet_device" "test" {
    hostname         = "terraform-test-device-va"
    plan             = "baremetal_0"
    facility         = "ewr1"
    operating_system = "ubuntu_16_04"
    billing_cycle    = "hourly"
    project_id       = "${packet_project.test.id}"
}

resource "packet_volume" "test" {
    plan = "storage_1"
    billing_cycle = "hourly"
    size = 100
    project_id = "${packet_project.test.id}"
    facility = "ewr1"
    snapshot_policies = { snapshot_frequency = "1day", snapshot_count = 7 }
}

resource "packet_volume_attachment" "test" {
	device_id = "${packet_device.test.id}"
	volume_id = "${packet_volume.test.id}"
}`, name)
}
