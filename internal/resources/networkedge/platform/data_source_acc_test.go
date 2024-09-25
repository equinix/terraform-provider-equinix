package platform

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"
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
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
	return nprintf.Nprintf(`
data "equinix_network_device_platform" "%{resourceName}" {
  device_type = "%{device_type}"
  flavor      = "%{flavor}"
  packages    = %{packages}
}
`, ctx)
}
