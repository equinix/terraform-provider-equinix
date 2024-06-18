package virtual_circuit

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/equinix/equinix-sdk-go/services/metalv1"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	metal := config.NewMetalClientForTesting()
	orgList, err := metal.OrganizationsApi.FindOrganizations(context.Background()).ExecuteWithPagination()
	if err != nil {
		return fmt.Errorf("error getting organization list for sweeping VirtualCircuits: %s", err)
	}
	for _, org := range orgList.Organizations {
		conns, _, err := metal.InterconnectionsApi.OrganizationListInterconnections(context.Background(), org.GetId()).Include([]string{"ports"}).Execute()
		if err != nil {
			errs = append(errs, fmt.Errorf("error getting connections list for sweeping VirtualCircuits: %s", err))
		}
		for _, conn := range conns.GetInterconnections() {
			if conn.GetType() == metalv1.INTERCONNECTIONTYPE_DEDICATED {
				for _, port := range conn.Ports {
					for _, vc := range port.VirtualCircuits {
						vcId := ""
						vcName := ""
						if vc.VlanVirtualCircuit != nil {
							vcId = vc.VlanVirtualCircuit.GetId()
							vcName = vc.VlanVirtualCircuit.GetName()
						} else {
							vcId = vc.VrfVirtualCircuit.GetId()
							vcName = vc.VlanVirtualCircuit.GetName()
						}
						if sweep.IsSweepableTestResource(vc.VlanVirtualCircuit.GetName()) {
							log.Printf("[INFO][SWEEPER_LOG] Deleting VirtualCircuit: %s", vcName)
							_, resp, err := metal.InterconnectionsApi.DeleteVirtualCircuit(context.Background(), vcId).Execute()
							if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
								errs = append(errs, fmt.Errorf("error deleting VirtualCircuit: %s", err))
							}
						}
					}
				}
			}
		}
	}

	return errors.Join(errs...)
}
