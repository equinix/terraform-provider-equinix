package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccDataSourceMetalSpotMarketRequest_Basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))
	var key packngo.SpotMarketRequest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSpotMarketRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalSpotMarketRequestConfig_Basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSpotMarketRequestExists("metal_spot_market_request.req", &key),
					resource.TestCheckResourceAttr(
						"data.metal_spot_market_request.dreq", "device_ids.#", "2"),
				),
			},
		},
	})
}

func testDataSourceMetalSpotMarketRequestConfig_Basic(projSuffix string) string {
	return fmt.Sprintf(`

resource "metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

resource "metal_spot_market_request" "req" {
  project_id    = "${metal_project.test.id}"
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

data "metal_spot_market_request" "dreq" {
  request_id = metal_spot_market_request.req.id
}
`, projSuffix)
}
