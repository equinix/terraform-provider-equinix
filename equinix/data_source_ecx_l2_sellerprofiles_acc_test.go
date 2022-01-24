package equinix

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccECXL2SellerProfiles(t *testing.T) {
	context := map[string]interface{}{
		"resourceName":             "test",
		"name_regex":               ".+Direct Connect.*",
		"metro_codes":              []string{"SV", "DC"},
		"speed_bands":              []string{"1GB", "100MB"},
		"organization_global_name": "AWS",
	}
	resourceName := fmt.Sprintf("data.equinix_ecx_l2_sellerprofiles.%s", context["resourceName"].(string))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccECXL2SellerProfiles(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckECXL2SellerProfiles(resourceName),
				),
			},
		},
	})
}

func testAccECXL2SellerProfiles(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_ecx_l2_sellerprofiles" "%{resourceName}" {
  name_regex               = "%{name_regex}"
  metro_codes              = %{metro_codes}
  speed_bands              = %{speed_bands}
  organization_global_name = "%{organization_global_name}"
}
`, ctx)
}

func testAccCheckECXL2SellerProfiles(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		profilesNumber, ok := rs.Primary.Attributes["profiles.#"]
		if !ok {
			return fmt.Errorf("profiles are not set")
		}
		if profilesNumberInt, _ := strconv.Atoi(profilesNumber); profilesNumberInt < 1 {
			return fmt.Errorf("number of profiles should be at least 1 but is %v", profilesNumberInt)
		}
		return nil
	}
}
