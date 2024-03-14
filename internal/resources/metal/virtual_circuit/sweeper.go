package virtual_circuit

import (
	"errors"
	"fmt"
	"log"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/packethost/packngo"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_virtual_circuit", &resource.Sweeper{
		Name:         "equinix_metal_virtual_circuit",
		Dependencies: []string{},
		F:            testSweepVirtualCircuits,
	})
}

func testSweepVirtualCircuits(region string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping VirtualCircuits")
	config, err := sweep.GetConfigForMetal()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping VirtualCircuits: %s", err)
	}
	metal := config.NewMetalClient()
	orgList, _, err := metal.Organizations.List(nil)
	if err != nil {
		return fmt.Errorf("error getting organization list for sweeping VirtualCircuits: %s", err)
	}
	vcs := map[string]*packngo.VirtualCircuit{}
	for _, org := range orgList {
		conns, _, err := metal.Connections.OrganizationList(org.ID, &packngo.GetOptions{Includes: []string{"ports"}})
		if err != nil {
			errs = append(errs, fmt.Errorf("error getting connections list for sweeping VirtualCircuits: %s", err))
		}
		for _, conn := range conns {
			if conn.Type == packngo.ConnectionDedicated {
				for _, port := range conn.Ports {
					for _, vc := range port.VirtualCircuits {
						if sweep.IsSweepableTestResource(vc.Name) {
							vcs[vc.ID] = &vc
						}
					}
				}
			}
		}
	}
	for _, vc := range vcs {
		log.Printf("[INFO][SWEEPER_LOG] Deleting VirtualCircuit: %s", vc.Name)
		_, err := metal.VirtualCircuits.Delete(vc.ID)
		if err != nil {
			errs = append(errs, fmt.Errorf("error deleting VirtualCircuit: %s", err))
		}
	}

	return errors.Join(errs...)
}
