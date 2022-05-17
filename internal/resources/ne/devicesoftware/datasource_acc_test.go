package devicesoftware

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/tfacc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkDeviceSoftware_versionRegex(t *testing.T) {
	context := map[string]interface{}{
		"resourceName":  "csrLatest",
		"device_type":   "CSR1000V",
		"version_regex": "^16.09.+",
		"packages":      []string{"IPBASE"},
		"most_recent":   true,
	}
	resourceName := fmt.Sprintf("data.equinix_network_device_software.%s", context["resourceName"].(string))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { tfacc.PreCheck(t) },
		Providers: tfacc.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkDeviceSoftwareConfig_versionRegex(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "date"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "stable"),
					resource.TestCheckResourceAttrSet(resourceName, "release_notes_link"),
				),
			},
		},
	})
}

func testAccDataSourceNetworkDeviceSoftwareConfig_versionRegex(ctx map[string]interface{}) string {
	return tfacc.NPrintf(`
data "equinix_network_device_software" "%{resourceName}" {
  device_type   = "%{device_type}"
  version_regex = "%{version_regex}"
  packages      = %{packages}
  most_recent   = %{most_recent}
}
`, ctx)
}
