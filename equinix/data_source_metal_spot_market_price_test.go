package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetalSpotPrice_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSpotMarketRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetalSpotMarketPrice(),
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

func testDataSourceMetalSpotMarketPrice() string {
	return fmt.Sprintf(`
data "equinix_metal_spot_market_price" "metro" {
	metro = "sv"
	plan     = "c3.medium.x86"
}

data "equinix_metal_spot_market_price" "facility" {
	facility = "sjc1"
	plan     = "c3.medium.x86"
}
`)
}
