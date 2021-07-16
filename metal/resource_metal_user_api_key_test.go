package metal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("metal_user_api_key", &resource.Sweeper{
		Name: "metal_user_api_key",
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
	client := meta.(*packngo.Client)

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

func testAccMetalUserAPIKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_user_api_key" {
			continue
		}
		if _, err := client.APIKeys.UserGet(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("UserAPI key still exists")
		}
	}
	return nil
}

func testAccMetalUserAPIKeyConfig_Basic() string {
	return fmt.Sprintf(`
resource "metal_user_api_key" "test" {
    description = "tfacc-user-key"
    read_only   = true
}`)
}

// Commented out because it lists existing user API keys in debug log

/*

func TestAccMetalUserAPIKey_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalUserAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalUserAPIKeyConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"metal_user_api_key.test", "token"),
					resource.TestCheckResourceAttrSet(
						"metal_user_api_key.test", "user_id"),
				),
			},
		},
	})
}

*/
