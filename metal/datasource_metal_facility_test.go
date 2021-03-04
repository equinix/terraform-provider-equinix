package metal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
