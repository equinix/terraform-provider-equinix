package organization

import (
	"fmt"
	"log"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_organization", &resource.Sweeper{
		Name:         "equinix_metal_organization",
		Dependencies: []string{"equinix_metal_project"},
		F:            testSweepOrganizations,
	})
}

func testSweepOrganizations(region string) error {
	log.Printf("[DEBUG] Sweeping organizations")
	config, err := sweep.GetConfigForMetal()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping organizations: %s", err)
	}
	metal := config.NewMetalClient()
	os, _, err := metal.Organizations.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting org list for sweeping organizations: %s", err)
	}
	oids := []string{}
	for _, o := range os {
		if sweep.IsSweepableTestResource(o.Name) {
			oids = append(oids, o.ID)
		}
	}
	for _, oid := range oids {
		log.Printf("Removing organization %s", oid)
		_, err := metal.Organizations.Delete(oid)
		if err != nil {
			return fmt.Errorf("Error deleting organization %s", err)
		}
	}
	return nil
}
