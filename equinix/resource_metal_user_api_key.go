package equinix

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func resourceMetalUserAPIKey() *schema.Resource {
	userKeySchema := schemaMetalAPIKey()
	userKeySchema["user_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "UUID of user owning this key",
	}
	return &schema.Resource{
		Create: resourceMetalAPIKeyCreate,
		Read:   resourceMetalAPIKeyRead,
		Delete: resourceMetalAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: userKeySchema,
	}
}

func addMetalUserAPIKeySweeper() {
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
