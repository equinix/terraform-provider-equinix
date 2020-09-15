package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNeAccount(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"resourceName": "tf-account",
		"metro_code":   "SV",
		"status":       "active",
	}
	resourceName := fmt.Sprintf("data.equinix_ne_account.%s", context["resourceName"].(string))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNeAccount(context),
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

func testAccNeAccount(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_ne_account" "%{resourceName}" {
  metro_code = "%{metro_code}"
  status     = "%{status}"
}
`, ctx)
}
