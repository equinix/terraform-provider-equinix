package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkDevicePlatform_basic(t *testing.T) {
	context := map[string]interface{}{
		"resourceName": "csrLarge",
		"device_type":  "CSR1000V",
		"flavor":       "large",
		"packages":     []string{"IPBASE"},
	}
	resourceName := fmt.Sprintf("data.equinix_network_device_platform.%s", context["resourceName"].(string))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkDevicePlatformConfig_basic(context),
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

func testAccDataSourceNetworkDevicePlatformConfig_basic(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_device_platform" "%{resourceName}" {
  device_type = "%{device_type}"
  flavor      = "%{flavor}"
  packages    = %{packages}
}
`, ctx)
}
