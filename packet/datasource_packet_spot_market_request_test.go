package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccDataSourcePacketSpotMarketRequest_Basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))
	var key packngo.SpotMarketRequest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketSpotMarketRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePacketSpotMarketRequestConfig_Basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketSpotMarketRequestExists("packet_spot_market_request.req", &key),
					resource.TestCheckResourceAttr(
						"data.packet_spot_market_request.dreq", "device_ids.#", "2"),
				),
			},
		},
	})
}

func testDataSourcePacketSpotMarketRequestConfig_Basic(projSuffix string) string {
	return fmt.Sprintf(`

resource "packet_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

resource "packet_spot_market_request" "req" {
  project_id    = "${packet_project.test.id}"
  max_bid_price = 0.2
  facilities    = ["sjc1"]
  devices_min   = 2
  devices_max   = 2
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_16_04"
    plan             = "t1.small.x86"
  }
}

data "packet_spot_market_request" "dreq" {
  request_id = packet_spot_market_request.req.id
}
`, projSuffix)
}
