package metal

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFacilityDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalFacilityDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_facility.test", "code", "ewr1"),
				),
			},
		},
	})
}

func testAccCheckMetalFacilityDataSourceConfigBasic() string {
	return `
data "metal_facility" "test" {
    code = "ewr1"
}
`
}

func testAccDataSourceFacilityConfigCapacity() string {
	return `
data "metal_facility" "test" {
    code = "ewr1"
    capacity {
        plan = "t1.small.x86"
        quantity = 1000
    }
}
`
}

var matchErrNoCapacity = regexp.MustCompile(`Not enough capacity.*`)

func TestAccDataSourceFacility_Capacity(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceFacilityConfigCapacity(),
				ExpectError: matchErrNoCapacity,
			},
		},
	})
}
