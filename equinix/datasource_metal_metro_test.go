package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMetro_Basic(t *testing.T) {
	testMetro := "da"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetroConfigBasic(testMetro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_metro.test", "code", testMetro),
				),
			},
			{
				Config: testAccDataSourceMetroConfigCapacityReasonable(testMetro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_metro.test", "code", testMetro),
				),
			},
			{
				Config:      testAccDataSourceMetroConfigCapacityUnreasonable(testMetro),
				ExpectError: matchErrNoCapacity,
			},
			{
				Config:      testAccDataSourceMetroConfigCapacityUnreasonableMultiple(testMetro),
				ExpectError: matchErrNoCapacity,
			},
		},
	})
}

func testAccDataSourceMetroConfigBasic(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
}
`, facCode)
}

func testAccDataSourceMetroConfigCapacityUnreasonable(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1000
    }
}
`, facCode)
}

func testAccDataSourceMetroConfigCapacityReasonable(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1
    }
    capacity {
        plan = "c3.medium.x86"
        quantity = 1
    }
}
`, facCode)
}

func testAccDataSourceMetroConfigCapacityUnreasonableMultiple(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
    capacity {
        plan = "c3.small.x86"
        quantity = 1
    }
    capacity {
        plan = "c3.medium.x86"
        quantity = 1000
    }
}
`, facCode)
}
