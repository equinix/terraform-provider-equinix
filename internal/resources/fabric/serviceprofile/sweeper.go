package serviceprofile

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_service_profile", &resource.Sweeper{
		Name:         "equinix_fabric_service_profile",
		Dependencies: []string{"equinix_fabric_connection"},
		F:            testSweepServiceProfiles,
	})
}

func testSweepServiceProfiles(_ string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Fabric Service Profiles")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Fabric Service Profiles: %s", err)
	}
	configLoadErr := meta.Load(ctx)
	if configLoadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping Fabric Service Profiles: %s", err)
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	limit := int32(100)
	offset := int32(0)
	equalOperator := string(fabricv4.EXPRESSIONOPERATOR_EQUAL)
	state := string(fabricv4.SERVICEPROFILESORTBY_STATE)

	serviceProfileSearchRequest := fabricv4.ServiceProfileSearchRequest{
		Filter: &fabricv4.ServiceProfileFilter{
			ServiceProfileAndFilter: &fabricv4.ServiceProfileAndFilter{
				And: []fabricv4.ServiceProfileSimpleExpression{
					{
						Property: &state,
						Operator: &equalOperator,
						Values:   []string{string(fabricv4.SERVICEPROFILESTATEENUM_ACTIVE)},
					},
				},
			},
		},
		Pagination: &fabricv4.PaginationRequest{
			Limit:  &limit,
			Offset: &offset,
		},
	}

	viewPoint := fabricv4.GETSERVICEPROFILESVIEWPOINTPARAMETER_Z_SIDE
	fabricServiceProfiles, _, err := fabric.ServiceProfilesApi.SearchServiceProfiles(ctx).ServiceProfileSearchRequest(serviceProfileSearchRequest).ViewPoint(viewPoint).Execute()

	if err != nil {
		return fmt.Errorf("error getting fabric service profiles list for sweeping fabric service profiles: %s", err)
	}

	for _, serviceProfile := range fabricServiceProfiles.Data {
		if sweep.IsSweepableFabricTestResource(serviceProfile.GetName()) {
			log.Printf("[DEBUG] Deleting Fabric Service Profiles: %s", serviceProfile.GetName())
			_, httpResponse, err := fabric.ServiceProfilesApi.DeleteServiceProfileByUuid(ctx, serviceProfile.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(httpResponse, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric Service Profiles: %s", err))
			}
		}
	}

	return errors.Join(errs...)
}
