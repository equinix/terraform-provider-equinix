package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalMetro_basic(t *testing.T) {
	testMetro := "da"
	resource.ParallelTest(t, resource.TestCase{ // Step 3/4, expected an error with pattern, no match on: Error running pre-apply refresh: exit status 1
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalMetroConfig_basic(testMetro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_metro.test", "code", testMetro),
				),
			},
			{
				Config: testAccDataSourceMetalMetroConfig_capacityReasonable(testMetro),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_metro.test", "code", testMetro),
				),
			},
			{
				Config:      testAccDataSourceMetalMetroConfig_capacityUnreasonable(testMetro),
				ExpectError: matchErrNoCapacity,
			},
			{
				Config:      testAccDataSourceMetalMetroConfig_capacityUnreasonableMultiple(testMetro),
				ExpectError: matchErrNoCapacity,
			},
		},
	})
}

func testAccDataSourceMetalMetroConfig_basic(facCode string) string {
	return fmt.Sprintf(`
data "equinix_metal_metro" "test" {
    code = "%s"
}
`, facCode)
}

func testAccDataSourceMetalMetroConfig_capacityUnreasonable(facCode string) string {
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

func testAccDataSourceMetalMetroConfig_capacityReasonable(facCode string) string {
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

func testAccDataSourceMetalMetroConfig_capacityUnreasonableMultiple(facCode string) string {
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
