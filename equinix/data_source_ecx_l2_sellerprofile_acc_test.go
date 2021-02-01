package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const tstL2SellerProfileEnvVar = "TF_ACC_ECX_SELLER_PROFILE_NAME"

func TestAccECXL2SellerProfile(t *testing.T) {
	t.Parallel()
	profileName, _ := schema.EnvDefaultFunc(tstL2SellerProfileEnvVar, "AWS Service Profile")()
	context := map[string]interface{}{
		"resourceName": "tf-aws",
		"name":         profileName,
	}
	resourceName := fmt.Sprintf("data.equinix_ecx_l2_sellerprofile.%s", context["resourceName"].(string))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2SellerProfile(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "speed_from_api"),
					resource.TestCheckResourceAttrSet(resourceName, "speed_customization_allowed"),
					resource.TestCheckResourceAttrSet(resourceName, "redundancy_required"),
					resource.TestCheckResourceAttrSet(resourceName, "encapsulation"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_name"),
				),
			},
		},
	})
}

func testAccECXL2SellerProfile(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_ecx_l2_sellerprofile" "%{resourceName}" {
  name = "%{name}"
}
`, ctx)
}
