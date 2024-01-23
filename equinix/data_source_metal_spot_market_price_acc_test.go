package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalSpotMarketPrice_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalSpotMarketRequestCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalSpotMarketPriceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_spot_market_price.metro", "price"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_metal_spot_market_price.facility", "price"),
				),
			},
		},
	})
}

func testAccDataSourceMetalSpotMarketPriceConfig_basic() string {
	return fmt.Sprintf(`
data "equinix_metal_spot_market_price" "metro" {
	metro    = "sv"
	plan     = "c3.medium.x86"
}

data "equinix_metal_spot_market_price" "facility" {
	facility = "sv15"
	plan     = "c3.medium.x86"
}
`)
}
