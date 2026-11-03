package internetaccess

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/sweep"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	sweepPageSize    = int32(100)
	sweepDeleteLimit = 10 * time.Minute
)

// AddTestSweeper registers the sweeper for Internet Access services created by acceptance tests.
func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_internet_access_service", &resource.Sweeper{
		Name:         "equinix_fabric_internet_access_service",
		Dependencies: []string{},
		F:            testSweepInternetAccessServices,
	})
}

// isTestService tells services created by acceptance tests from everybody else's by the name suffix.
func isTestService(svc *fabricv4.InternetAccessService) bool {
	return sweep.IsSweepableFabricTestResource(svc.GetName())
}

// sweepableStates are the states in which a service still has to be deleted. FAILED is left out:
// the API keeps failed services in that state after a DELETE, so they could never be swept away.
func sweepableStates() []string {
	return []string{
		string(fabricv4.INTERNETACCESSSERVICESTATE_PROVISIONING),
		string(fabricv4.INTERNETACCESSSERVICESTATE_PROVISIONED),
	}
}

// testSweepInternetAccessServices deprovisions Internet Access services whose name ends with one
// of sweep.FabricTestResourceSuffixes. It waits for the deprovisioning because the connections
// the services use cannot be removed before that (the connection sweeper runs afterwards).
func testSweepInternetAccessServices(_ string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Internet Access services")

	ctx := context.Background()
	meta, err := testinghelpers.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Internet Access services: %s", err)
	}
	if loadErr := meta.Load(ctx); loadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping Internet Access services: %s", loadErr)
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	property, operator := "/state", fabricv4.SearchExpressionOperator("IN")
	filter := fabricv4.SearchExpression{
		And: []fabricv4.SearchExpression{{Property: &property, Operator: &operator, Values: sweepableStates()}},
	}

	for offset := int32(0); ; {
		limit, pageOffset := sweepPageSize, offset
		searchReq := fabricv4.InternetAccessSearchRequest{
			Filter:     &filter,
			Pagination: &fabricv4.PaginationRequest{Offset: &pageOffset, Limit: &limit},
		}
		result, _, searchErr := fabric.InternetAccessServicesApi.SearchEiaServices(ctx).InternetAccessSearchRequest(searchReq).Execute()
		if searchErr != nil {
			return fmt.Errorf("error searching Internet Access services for sweep: %s", equinix_errors.FormatFabricError(searchErr))
		}

		services := result.GetData()
		for i := range services {
			svc := &services[i]
			if !isTestService(svc) {
				continue
			}
			log.Printf("[DEBUG] Deleting Internet Access service: %s (%s)", svc.GetName(), svc.GetUuid())
			_, resp, delErr := fabric.InternetAccessServicesApi.DeleteEiaService(ctx, svc.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, delErr) != nil {
				errs = append(errs, fmt.Errorf("error deleting Internet Access service %s: %s", svc.GetUuid(), delErr))
				continue
			}
			if _, waitErr := getDeleteWaiter(ctx, fabric, svc.GetUuid(), sweepDeleteLimit).WaitForStateContext(ctx); waitErr != nil {
				errs = append(errs, fmt.Errorf("error waiting for Internet Access service %s to be deleted: %s", svc.GetUuid(), waitErr))
			}
		}

		if len(services) < int(sweepPageSize) {
			break
		}
		// deleted services leave the filtered result set, so only skip the ones that were kept
		offset += int32(countKept(services))
	}

	return errors.Join(errs...)
}

// countKept counts the services of a page that the sweeper leaves alone and that therefore
// stay in the result set of the next search.
func countKept(services []fabricv4.InternetAccessService) int {
	kept := 0
	for i := range services {
		if !isTestService(&services[i]) {
			kept++
		}
	}
	return kept
}
