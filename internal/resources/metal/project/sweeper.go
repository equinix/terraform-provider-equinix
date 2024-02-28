package project

import (
	"fmt"
	"log"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_project", &resource.Sweeper{
		Name:         "equinix_metal_project",
		Dependencies: []string{"equinix_metal_device", "equinix_metal_vlan"},
		F:            testSweepProjects,
	})
}

func testSweepProjects(region string) error {
	var errs error
	log.Printf("[DEBUG] Sweeping projects")
	config, err := sweep.GetConfigForMetal()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping projects: %s", err)
	}
	metal := config.NewMetalClient()
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping projects: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if sweep.IsSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	for _, pid := range pids {
		log.Printf("Removing project %s", pid)
		_, err := metal.Projects.Delete(pid)
		if err != nil {
			errs = multierror.Append(errs, fmt.Errorf("Error deleting project %s", err))
		}
	}
	return errs
}
