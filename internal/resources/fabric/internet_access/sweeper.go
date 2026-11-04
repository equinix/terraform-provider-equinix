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

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_internet_access_service", &resource.Sweeper{
		Name:         "equinix_fabric_internet_access_service",
		Dependencies: []string{},
		F:            testSweepInternetAccessServices,
	})
}

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
	states := []string{
		string(fabricv4.INTERNETACCESSSERVICESTATE_PROVISIONING),
		string(fabricv4.INTERNETACCESSSERVICESTATE_PROVISIONED),
	}
	filter := fabricv4.SearchExpression{
		And: []fabricv4.SearchExpression{{Property: &property, Operator: &operator, Values: states}},
	}

	var testServices []fabricv4.InternetAccessService
	for offset := int32(0); ; offset += sweepPageSize {
		limit, pageOffset := sweepPageSize, offset
		searchReq := fabricv4.InternetAccessSearchRequest{
			Filter:     &filter,
			Pagination: &fabricv4.PaginationRequest{Offset: &pageOffset, Limit: &limit},
		}
		result, _, searchErr := searchServices(ctx, fabric, searchReq)
		if searchErr != nil {
			return fmt.Errorf("error searching Internet Access services for sweep: %s", equinix_errors.FormatFabricError(searchErr))
		}

		services := result.GetData()
		for _, svc := range services {
			if sweep.IsSweepableFabricTestResource(svc.GetName()) {
				testServices = append(testServices, svc)
			}
		}
		if len(services) < int(sweepPageSize) {
			break
		}
	}

	for _, svc := range testServices {
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

	return errors.Join(errs...)
}
