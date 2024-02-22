package vrf

import (
	"fmt"
	"log"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_vrf", &resource.Sweeper{
		Name: "equinix_metal_vrf",
		Dependencies: []string{
			"equinix_metal_device",
			"equinix_metal_virtual_circuit",
			// TODO: add sweeper when offered
			// "equinix_metal_reserved_ip_block",
		},
		F: testSweepVRFs,
	})
}

func testSweepVRFs(region string) error {
	log.Printf("[DEBUG] Sweeping VRFs")
	config, err := sweep.GetConfigForMetal()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping VRFs: %s", err)
	}
	metal := config.NewMetalClient()
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping VRFs: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if sweep.IsSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := metal.VRFs.List(pid, nil)
		if err != nil {
			log.Printf("Error listing VRFs to sweep: %s", err)
			continue
		}
		for _, d := range ds {
			if sweep.IsSweepableTestResource(d.Name) {
				dids = append(dids, d.ID)
			}
		}
	}

	for _, did := range dids {
		log.Printf("Removing VRFs %s", did)
		_, err := metal.VRFs.Delete(did)
		if err != nil {
			return fmt.Errorf("Error deleting VRFs %s", err)
		}
	}
	return nil
}
