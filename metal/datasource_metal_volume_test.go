package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccMetalVolumeDataSource_Basic(t *testing.T) {
	var volume packngo.Volume

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVolumeDataSourceConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalVolumeExists("metal_volume.test", &volume),
					resource.TestCheckResourceAttr(
						"metal_volume.test", "plan", "storage_1"),
					resource.TestCheckResourceAttr(
						"metal_volume.test", "billing_cycle", "hourly"),
					resource.TestCheckResourceAttr(
						"metal_volume.test", "size", "100"),
					resource.TestCheckResourceAttrPair(
						"data.metal_volume.test", "device_ids.0",
						"metal_device.test", "id"),
				),
			},
		},
	})
}

func testAccCheckMetalVolumeDataSourceConfig_basic(projectSuffix string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-volume-%s"
}

resource "metal_device" "test" {
    hostname         = "tfacc-device-volume-datasource"
    plan             = "t1.small.x86"
    facilities       = ["ewr1"]
    operating_system = "ubuntu_16_04"
    billing_cycle    = "hourly"
    project_id       = metal_project.test.id
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
	device_id = metal_device.test.id
	volume_id = metal_volume.test.id
}

data "metal_volume" "test" {
	volume_id = metal_volume_attachment.test.volume_id
}

data "metal_volume" "test2" {
	name = metal_volume.test.name
	project_id = metal_project.test.id
}`, projectSuffix)
}
