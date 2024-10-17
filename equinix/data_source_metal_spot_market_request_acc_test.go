package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/packethost/packngo"
)

func TestAccDataSourceMetalSpotMarketRequest_basic(t *testing.T) {
	projectName := fmt.Sprintf("ds-device-%s", acctest.RandString(10))
	var (
		facKey packngo.SpotMarketRequest
		metKey packngo.SpotMarketRequest
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalSpotMarketRequestCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalSpotMarketRequestConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSpotMarketRequestExists("equinix_metal_spot_market_request.req", &facKey),
				),
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
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

resource "equinix_metal_spot_market_request" "req" {
  project_id    = "${equinix_metal_project.test.id}"
  max_bid_price = 0.01
  facilities    = local.facilities
  devices_min   = 1
  devices_max   = 1
  wait_for_devices = false

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = local.os
    plan             = local.plan
  }
}

data "equinix_metal_spot_market_request" "dreq" {
  request_id = equinix_metal_spot_market_request.req.id
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix)
}

func testAccDataSourceMetalSpotMarketRequestConfig_metro(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

resource "equinix_metal_spot_market_request" "req" {
  project_id    = equinix_metal_project.test.id
  max_bid_price = 0.01
  metro         = local.metro
  devices_min   = 1
  devices_max   = 1
  wait_for_devices = false

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = local.os
    plan             = local.plan
  }
}

data "equinix_metal_spot_market_request" "dreq" {
  request_id = equinix_metal_spot_market_request.req.id
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix)
}
