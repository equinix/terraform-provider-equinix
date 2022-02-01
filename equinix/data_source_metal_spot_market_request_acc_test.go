package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func TestAccDataSourceMetalSpotMarketRequest_basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))
	var (
		facKey packngo.SpotMarketRequest
		metKey packngo.SpotMarketRequest
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalSpotMarketRequestCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalSpotMarketRequestConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSpotMarketRequestExists("equinix_metal_spot_market_request.req", &facKey),
				),
			},
			{
				Config:             testAccDataSourceMetalSpotMarketRequestConfig_metro(projectName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccDataSourceMetalSpotMarketRequestConfig_metro(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSpotMarketRequestExists("equinix_metal_spot_market_request.req", &metKey),
					func(_ *terraform.State) error {
						if metKey.ID == facKey.ID {
							return fmt.Errorf("Expected a new spot_market_request")
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccDataSourceMetalSpotMarketRequestConfig_basic(projSuffix string) string {
	return fmt.Sprintf(`

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

resource "equinix_metal_spot_market_request" "req" {
  project_id    = "${equinix_metal_project.test.id}"
  max_bid_price = 0.01
  facilities    = ["sjc1"]
  devices_min   = 1
  devices_max   = 1
  wait_for_devices = false

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_20_04"
    plan             = "c3.small.x86"
  }
}

data "equinix_metal_spot_market_request" "dreq" {
  request_id = metal_spot_market_request.req.id
}
`, projSuffix)
}

func testAccDataSourceMetalSpotMarketRequestConfig_metro(projSuffix string) string {
	return fmt.Sprintf(`

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

resource "equinix_metal_spot_market_request" "req" {
  project_id    = "${equinix_metal_project.test.id}"
  max_bid_price = 0.01
  metro = "sv"
  devices_min   = 1
  devices_max   = 1
  wait_for_devices = false

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_20_04"
    plan             = "c3.small.x86"
  }
}

data "equinix_metal_spot_market_request" "dreq" {
  request_id = metal_spot_market_request.req.id
}
`, projSuffix)
}
