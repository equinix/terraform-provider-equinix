package equinix

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("equinix_metal_user_api_key", &resource.Sweeper{
		Name: "equinix_metal_user_api_key",
		F:    testSweepUserAPIKeys,
	})
}

// Will remove all User API keys with Description starting with "tfacc-"
func testSweepUserAPIKeys(region string) error {
	log.Printf("[DEBUG] Sweeping user_api keys")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping user_api keys: %s", err)
	}
	client := meta.Client()

	userApiKeys, _, err := client.APIKeys.UserList(nil)
	if err != nil {
		return fmt.Errorf("Error getting list for sweeping user_api keys: %s", err)
	}
	ids := []string{}
	for _, k := range userApiKeys {
		if strings.HasPrefix(k.Description, "tfacc-") {
			ids = append(ids, k.ID)
		}
	}
	for _, id := range ids {
		log.Printf("Removing user api key %s", id)
		resp, err := client.APIKeys.Delete(id)
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Error deleting user_api key %s", err)
		}
	}
	return nil
}

func testAccMetalUserAPIKeyCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).Client()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_user_api_key" {
			continue
		}
		if _, err := client.APIKeys.UserGet(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal UserAPI key still exists")
		}
	}
	return nil
}

func testAccMetalUserAPIKeyConfig_basic() string {
	return fmt.Sprintf(`
resource "equinix_metal_user_api_key" "test" {
    description = "tfacc-user-key"
    read_only   = true
}`)
}

// Commented out because it lists existing user API keys in debug log

/*

func TestAccMetalUserAPIKey_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalUserAPIKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalUserAPIKeyConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"equinix_metal_user_api_key.test", "token"),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_user_api_key.test", "user_id"),
				),
			},
		},
	})
}

*/
