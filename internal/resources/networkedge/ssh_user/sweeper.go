package ssh_user

import (
	"context"
	"fmt"
	"log"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_network_ssh_user", &resource.Sweeper{
		Name: "equinix_network_ssh_user",
		F:    testSweepNetworkSSHUser,
	})
}

func testSweepNetworkSSHUser(region string) error {
	config, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping Network SSH users: %s", err)
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	users, err := config.Ne.GetSSHUsers()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching NetworkSSHUser list: %s", err)
		return err
	}
	for _, user := range users {
		if !sweep.IsSweepableTestResource(ne.StringValue(user.Username)) {
			continue
		}
		if err := config.Ne.DeleteSSHUser(ne.StringValue(user.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkSSHUser resource %s (%s): %s", ne.StringValue(user.UUID), ne.StringValue(user.Username), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkSSHUser resource %s (%s)", ne.StringValue(user.UUID), ne.StringValue(user.Username))
		}
	}
	return nil
}
