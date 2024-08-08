package device

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_device", &resource.Sweeper{
		Name: "equinix_metal_device",
		F:    testSweepDevices,
	})
}

func testSweepDevices(region string) error {
	var errs error
	log.Printf("[DEBUG] Sweeping devices")
	ctx := context.Background()
	config, err := sweep.GetConfigForMetal()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping devices: %s", err)
	}
	metal := config.NewMetalClientForTesting()
	ps, err := metal.ProjectsApi.FindProjects(ctx).ExecuteWithPagination()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweepeing devices: %s", err)
	}
	pids := []string{}
	for _, p := range ps.Projects {
		if sweep.IsSweepableTestResource(p.GetName()) {
			pids = append(pids, p.GetId())
		}
	}

	for _, pid := range pids {
		ds, _, err := metal.DevicesApi.FindProjectDevices(ctx, pid).Execute()
		if err != nil {
			log.Printf("Error listing devices to sweep: %s", err)
			continue
		}
		for _, d := range ds.Devices {
			err := sweepDevice(ctx, metal, d)
			if err != nil {
				errs = multierror.Append(errs, fmt.Errorf("Error deleting device %s", err))
			}
		}
	}

	return errs
}

func sweepDevice(ctx context.Context, metal *metalv1.APIClient, d metalv1.Device) error {
	if sweep.IsSweepableTestResource(d.GetHostname()) {
		nonSweepableDeviceStates := []metalv1.DeviceState{
			metalv1.DEVICESTATE_PROVISIONING,
			metalv1.DEVICESTATE_DEPROVISIONING,
		}

		if slices.Contains(nonSweepableDeviceStates, d.GetState()) {
			log.Printf("[WARNING] skipping sweep for device %s because it is still %s", d.GetId(), d.GetState())
		} else {
			if d.GetLocked() {
				// If a device is locked, we have to unlock it first before deleting
				unlockDevice := metalv1.DeviceUpdateInput{
					Locked: metalv1.PtrBool(false),
				}
				_, _, err := metal.DevicesApi.UpdateDevice(ctx, d.GetId()).DeviceUpdateInput(unlockDevice).Execute()
				if err != nil {
					return err
				}
			}

			_, err := metal.DevicesApi.DeleteDevice(ctx, d.GetId()).Execute()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
