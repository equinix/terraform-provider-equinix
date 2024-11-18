package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourcePlans_Basic(t *testing.T) {
	testSlug := "m2.xlarge.x86"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePlansConfigBasic(testSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_metal_plans.test", "plans.0.slug", testSlug),
				),
			},
		},
	})
}

func testAccDataSourcePlansConfigBasic(slug string) string {
	return fmt.Sprintf(`
data "equinix_metal_plans" "test" {
    filter {
        attribute = "slug"
        values    = ["%s"]
    }
}

output "test" {
    value = data.equinix_metal_plans.test
}
`, slug)
}
