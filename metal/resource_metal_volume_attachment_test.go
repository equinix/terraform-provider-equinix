package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func TestAccMetalVolumeAttachment_Basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVolumeAttachmentConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_volume_attachment.test", "volume_id",
						"metal_volume.test", "id"),
					resource.TestCheckResourceAttrPair(
						"metal_volume_attachment.test", "device_id",
						"metal_device.test", "id"),
				),
			},
			{
				ResourceName:      "metal_volume_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMetalVolumeAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_volume_attachment" {
			continue
		}
		if _, _, err := client.VolumeAttachments.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("VolumeAttachment still exists")
		}
	}

	return nil
}

func testAccCheckMetalVolumeAttachmentConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-volume_attachment-%s"
}

resource "metal_device" "test" {
    hostname         = "tfacc-test-device-va"
    plan             = "t1.small.x86"
    facilities       = ["ewr1"]
    operating_system = "ubuntu_16_04"
    billing_cycle    = "hourly"
    project_id       = "${metal_project.test.id}"
}

resource "metal_volume" "test" {
    plan = "storage_1"
    billing_cycle = "hourly"
    size = 100
    project_id = "${metal_project.test.id}"
    facility = "ewr1"
    snapshot_policies {
	    snapshot_frequency = "1day"
		snapshot_count = 7
	}
}

resource "metal_volume_attachment" "test" {
	device_id = "${metal_device.test.id}"
	volume_id = "${metal_volume.test.id}"
}`, name)
}
