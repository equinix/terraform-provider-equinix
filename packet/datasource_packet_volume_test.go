package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccPacketVolumeDataSource_Basic(t *testing.T) {
	var volume packngo.Volume

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketVolumeDataSourceConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVolumeExists("packet_volume.test", &volume),
					resource.TestCheckResourceAttr(
						"packet_volume.test", "plan", "storage_1"),
					resource.TestCheckResourceAttr(
						"packet_volume.test", "billing_cycle", "hourly"),
					resource.TestCheckResourceAttr(
						"packet_volume.test", "size", "100"),
					resource.TestCheckResourceAttrPair(
						"data.packet_volume.test", "device_ids.0",
						"packet_device.test", "id"),
				),
			},
		},
	})
}

func testAccCheckPacketVolumeDataSourceConfig_basic(projectSuffix string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-volume-%s"
}

resource "packet_device" "test" {
    hostname         = "tfacc-device-volume-datasource"
    plan             = "t1.small.x86"
    facilities       = ["ewr1"]
    operating_system = "ubuntu_16_04"
    billing_cycle    = "hourly"
    project_id       = packet_project.test.id
}

resource "packet_volume" "test" {
    plan = "storage_1"
    billing_cycle = "hourly"
    size = 100
    project_id = "${packet_project.test.id}"
    facility = "ewr1"
    snapshot_policies { 
		snapshot_frequency = "1day"
		snapshot_count = 7
	}
}

resource "packet_volume_attachment" "test" {
	device_id = packet_device.test.id
	volume_id = packet_volume.test.id
}

data "packet_volume" "test" {
	volume_id = packet_volume_attachment.test.volume_id
}

data "packet_volume" "test2" {
	name = packet_volume.test.name
	project_id = packet_project.test.id
}`, projectSuffix)
}
