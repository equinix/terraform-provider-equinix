package connection

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/packethost/packngo"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_metal_connection", &resource.Sweeper{
		Name:         "equinix_metal_connection",
		Dependencies: []string{},
		F:            testSweepMetalConnections,
	})
}

func testSweepMetalConnections(region string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Connections")
	ctx := context.Background()
	config, err := sweep.GetConfigForMetal()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Conections: %s", err)
	}
	metal := config.NewMetalClientForTesting()
	orgList, err := metal.OrganizationsApi.FindOrganizations(ctx).Exclude([]string{"address", "billing_address"}).ExecuteWithPagination()
	if err != nil {
		return fmt.Errorf("error getting organization list for sweeping Connections: %s", err)
	}

	for _, org := range orgList.Organizations {
		conns, _, err := metal.InterconnectionsApi.OrganizationListInterconnections(ctx, org.GetId()).Execute()
		if err != nil {
			errs = append(errs, fmt.Errorf("error getting connections list for sweeping Connections: %s", err))
			continue
		}
		for _, conn := range conns.GetInterconnections() {
			if sweep.IsSweepableTestResource(conn.GetName()) {
				if packngo.ConnectionType(conn.GetType()) != packngo.ConnectionType(metalv1.INTERCONNECTIONTYPE_DEDICATED) {
					log.Printf("[INFO][SWEEPER_LOG] Deleting Connection: %s", conn.GetId())
					_, _, err := metal.InterconnectionsApi.DeleteInterconnection(ctx, conn.GetId()).Execute()
					if err != nil {
						errs = append(errs, fmt.Errorf("error deleting VirtualCircuit: %s", err))
					}
				}
			}
		}
	}

	return errors.Join(errs...)
}
