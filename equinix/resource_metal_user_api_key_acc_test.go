package equinix_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping user_api keys: %s", err)
	}
	metal := config.NewMetalClient()
	userApiKeys, _, err := metal.APIKeys.UserList(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting list for sweeping user_api keys: %s", err)
	}
	ids := []string{}
	for _, k := range userApiKeys {
		if isSweepableTestResource(k.Description) {
			ids = append(ids, k.ID)
		}
	}
	for _, id := range ids {
		log.Printf("Removing user api key %s", id)
		resp, err := metal.APIKeys.Delete(id)
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Error deleting user_api key %s", err)
		}
	}
	return nil
}

// Commented out because it lists existing user API keys in debug log

/*

func TestAccMetalUserAPIKey_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:    acceptance.TestExternalProviders,
		Providers:    acceptance.TestAccProviders,
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

func testAccMetalUserAPIKeyCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal
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
*/
