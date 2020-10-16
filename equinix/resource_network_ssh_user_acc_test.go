package equinix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// SSH User acc tests are in Device acc test tests
// reason: SSH User requires device to be provisioned and that is time consuming operation

func init() {
	resource.AddTestSweepers("NetworkSSHUser", &resource.Sweeper{
		Name: "NetworkSSHUser",
		F:    testSweepNetworkSSHUser,
	})
}

func testSweepNetworkSSHUser(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	users, err := config.ne.GetSSHUsers()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching NetworkSSHUser list: %s", err)
		return err
	}
	for _, user := range users {
		if !isSweepableTestResource(user.Username) {
			continue
		}
		if err := config.ne.DeleteSSHUser(user.UUID); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkSSHUser resource %s (%s): %s", user.UUID, user.Username, err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkSSHUser resource %s (%s)", user.UUID, user.Username)
		}
	}
	return nil
}
