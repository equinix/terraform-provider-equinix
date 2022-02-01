package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataSourceNetworkDeviceType_basic(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"resourceName": "router",
		"category":     "Router",
		"vendor":       "Cisco",
		"metro_codes":  []string{metro.(string)},
	}
	resourceName := fmt.Sprintf("data.equinix_network_device_type.%s", context["resourceName"].(string))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkDeviceTypeConfig_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "code"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
		},
	})
}

func testAccDataSourceNetworkDeviceTypeConfig_basic(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_device_type" "%{resourceName}" {
  category    = "%{category}"
  vendor      = "%{vendor}"
  metro_codes = %{metro_codes}
}
`, ctx)
}
