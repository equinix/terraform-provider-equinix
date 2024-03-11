package user_api_key

import (
	"fmt"
	"log"
	"net/http"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_user_api_key", &resource.Sweeper{
		Name: "equinix_metal_user_api_key",
		F:    testSweepUserAPIKeys,
	})
}

// Will remove all User API keys with Description starting with "tfacc-"
func testSweepUserAPIKeys(region string) error {
	log.Printf("[DEBUG] Sweeping user_api keys")
	config, err := sweep.GetConfigForMetal()
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
		if sweep.IsSweepableTestResource(k.Description) {
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
