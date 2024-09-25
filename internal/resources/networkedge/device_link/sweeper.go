package device_link

import (
	"context"
	"fmt"
	"log"

	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_network_device_link", &resource.Sweeper{
		Name: "equinix_network_device_link",
		F:    testSweepNetworkDeviceLink,
	})
}

func testSweepNetworkDeviceLink(region string) error {
	config, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping Network devices link: %s", err)
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	links, err := config.Ne.GetDeviceLinkGroups()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching device links list: %s", err)
		return err
	}
	nonSweepableCount := 0
	for _, link := range links {
		if !sweep.IsSweepableTestResource(ne.StringValue(link.Name)) {
			nonSweepableCount++
			continue
		}
		if err := config.Ne.DeleteDeviceLinkGroup(ne.StringValue(link.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkDeviceLink resource %s (%s): %s", ne.StringValue(link.UUID), ne.StringValue(link.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkDeviceLink resource %s (%s)", ne.StringValue(link.UUID), ne.StringValue(link.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}
