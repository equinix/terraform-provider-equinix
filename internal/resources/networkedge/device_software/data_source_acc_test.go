package device_software

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"
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
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
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
	return nprintf.Nprintf(`
data "equinix_network_device_software" "%{resourceName}" {
  device_type   = "%{device_type}"
  version_regex = "%{version_regex}"
  packages      = %{packages}
  most_recent   = %{most_recent}
}
`, ctx)
}
