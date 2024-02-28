package ssh_key

import (
	"fmt"
	"log"
	"net/http"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_ssh_key", &resource.Sweeper{
		Name: "equinix_metal_ssh_key",
		F:    testSweepSSHKeys,
	})
}

func testSweepSSHKeys(region string) error {
	log.Printf("[DEBUG] Sweeping ssh keys")
	config, err := sweep.GetConfigForMetal()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping ssh keys: %s", err)
	}
	metal := config.NewMetalClient()
	sshkeys, _, err := metal.SSHKeys.List()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting list for sweeping ssh keys: %s", err)
	}
	ids := []string{}
	for _, k := range sshkeys {
		if sweep.IsSweepableTestResource(k.Label) {
			ids = append(ids, k.ID)
		}
	}
	for _, id := range ids {
		log.Printf("Removing ssh key %s", id)
		resp, err := metal.SSHKeys.Delete(id)
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Error deleting ssh key %s", err)
		}
	}
	return nil
}
