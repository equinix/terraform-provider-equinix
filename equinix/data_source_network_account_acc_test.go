package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNetworkAccount_basic(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"resourceName": "tf-account",
		"metro_code":   metro.(string),
		"status":       "active",
	}
	resourceName := fmt.Sprintf("data.equinix_network_account.%s", context["resourceName"].(string))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkAccountConfig_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "number"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "ucm_id"),
				),
			},
		},
	})
}

func testAccDataSourceNetworkAccountConfig_basic(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_account" "%{resourceName}" {
  metro_code = "%{metro_code}"
  status     = "%{status}"
}
`, ctx)
}
