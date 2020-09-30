package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNeDevicePlatform(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"resourceName": "csrLarge",
		"device_type":  "CSR1000V",
		"flavor":       "large",
		"packages":     []string{"IPBASE"},
	}
	resourceName := fmt.Sprintf("data.equinix_ne_device_platform.%s", context["resourceName"].(string))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNeDevicePlatform(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "flavor"),
					resource.TestCheckResourceAttrSet(resourceName, "core_count"),
					resource.TestCheckResourceAttrSet(resourceName, "memory"),
					resource.TestCheckResourceAttrSet(resourceName, "memory_unit"),
				),
			},
		},
	})
}

func testAccNeDevicePlatform(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_ne_device_platform" "%{resourceName}" {
  device_type = "%{device_type}"
  flavor      = "%{flavor}"
  packages    = %{packages}
}
`, ctx)
}
