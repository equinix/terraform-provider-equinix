package device

import (
	"context"
	"fmt"
	"log"

	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_network_device", &resource.Sweeper{
		Name:         "equinix_network_device",
		Dependencies: []string{"equinix_network_device_link"},
		F:            testSweepNetworkDevice,
	})
}

func testSweepNetworkDevice(region string) error {
	config, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping Network devices: %s", err)
	}
	if err := config.Load(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading configuration: %s", err)
		return err
	}
	devices, err := config.Ne.GetDevices([]string{
		ne.DeviceStateInitializing,
		ne.DeviceStateProvisioned,
		ne.DeviceStateProvisioning,
		ne.DeviceStateWaitingSecondary,
		ne.DeviceStateWaitingClusterNodes,
		ne.DeviceStateClusterSetUpInProgress,
		ne.DeviceStateFailed,
	})
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error fetching NetworkDevice list: %s", err)
		return err
	}
	nonSweepableCount := 0
	for _, device := range devices {
		if !sweep.IsSweepableTestResource(ne.StringValue(device.Name)) {
			nonSweepableCount++
			continue
		}
		if ne.StringValue(device.RedundancyType) != "PRIMARY" {
			continue
		}
		if err := config.Ne.DeleteDevice(ne.StringValue(device.UUID)); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error deleting NetworkDevice resource %s (%s): %s", ne.StringValue(device.UUID), ne.StringValue(device.Name), err)
		} else {
			log.Printf("[INFO][SWEEPER_LOG] sent delete request for NetworkDevice resource %s (%s)", ne.StringValue(device.UUID), ne.StringValue(device.Name))
		}
	}
	if nonSweepableCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonSweepableCount)
	}
	return nil
}
