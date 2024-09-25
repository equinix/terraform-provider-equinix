package ssh_key

import (
	"context"
	"fmt"
	"log"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_network_ssh_key", &resource.Sweeper{
		Name: "equinix_network_ssh_key",
		F:    testSweepNetworkSSHKey,
	})
}

func testSweepNetworkSSHKey(region string) error {
	config, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping Network SSH keys: %s", err)
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	keys, err := config.Ne.GetSSHPublicKeys()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching NetworkSSHKey list: %s", err)
		return err
	}
	nonSweepableCount := 0
	for _, key := range keys {
		if !sweep.IsSweepableTestResource(ne.StringValue(key.Name)) {
			nonSweepableCount++
			continue
		}
		if err := config.Ne.DeleteSSHPublicKey(ne.StringValue(key.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkSSHKey resource %s (%s): %s", ne.StringValue(key.UUID), ne.StringValue(key.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkSSHKey resource %s (%s)", ne.StringValue(key.UUID), ne.StringValue(key.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}
