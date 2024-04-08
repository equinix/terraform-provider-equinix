package vlan

import (
	"fmt"
	"log"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_vlan", &resource.Sweeper{
		Name:         "equinix_metal_vlan",
		Dependencies: []string{"equinix_metal_connection", "equinix_metal_virtual_circuit", "equinix_metal_vrf", "equinix_metal_device"},
		F:            testSweepVlans,
	})
}

func testSweepVlans(region string) error {
	log.Printf("[DEBUG] Sweeping vlan")
	config, err := sweep.GetConfigForMetal()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping vlan: %s", err)
	}
	metal := config.NewMetalClient()
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping vlan: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if sweep.IsSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := metal.ProjectVirtualNetworks.List(pid, nil)
		if err != nil {
			log.Printf("Error listing vlan to sweep: %s", err)
			continue
		}
		for _, d := range ds.VirtualNetworks {
			if sweep.IsSweepableTestResource(d.Description) {
				dids = append(dids, d.ID)
			}
		}
	}

	for _, did := range dids {
		log.Printf("Removing vlan %s", did)
		_, err := metal.ProjectVirtualNetworks.Delete(did)
		if err != nil {
			return fmt.Errorf("Error deleting vlan %s", err)
		}
	}
	return nil
}
