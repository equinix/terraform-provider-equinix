package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkDeviceTypeDataSource(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"resourceName": "router",
		"category":     "Router",
		"vendor":       "Cisco",
		"metro_codes":  []string{"DC"},
	}
	resourceName := fmt.Sprintf("data.equinix_network_device_type.%s", context["resourceName"].(string))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkDeviceTypeDataSource(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "code"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
		},
	})
}

func testAccNetworkDeviceTypeDataSource(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_device_type" "%{resourceName}" {
  category    = "%{category}"
  vendor      = "%{vendor}"
  metro_codes = %{metro_codes}
}
`, ctx)
}
