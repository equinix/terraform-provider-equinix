package metal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetroDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalMetroDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.metal_metro.test", "code", "sv"),
				),
			},
		},
	})
}

func testAccCheckMetalMetroDataSourceConfigBasic() string {
	return `
data "metal_metro" "test" {
    code = "sv"
}
`
}
