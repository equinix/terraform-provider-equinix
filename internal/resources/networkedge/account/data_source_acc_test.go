package account

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/nprintf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNetworkAccount_basic(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	context := map[string]interface{}{
		"resourceName": "tf-account",
		"metro_code":   metro.(string),
		"status":       "active",
		"project_id":   "92cbcfd9-347b-4da5-901d-2cea82575941",
	}
	resourceName := fmt.Sprintf("data.equinix_network_account.%s", context["resourceName"].(string))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkAccountConfig_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "number"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "ucm_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func testAccDataSourceNetworkAccountConfig_basic(ctx map[string]interface{}) string {
	return nprintf.NPrintf(`
data "equinix_network_account" "%{resourceName}" {
  metro_code = "%{metro_code}"
  status     = "%{status}"
}
`, ctx)
}
