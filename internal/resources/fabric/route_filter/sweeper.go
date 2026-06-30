package route_filter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_route_filter", &resource.Sweeper{
		Name:         "equinix_fabric_route_filter",
		Dependencies: []string{"equinix_fabric_connection"},
		F:            testSweepRouteFilters,
	})
}

func testSweepRouteFilters(region string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Fabric Route Filters")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Route Filters: %s", err)
	}
	configLoadErr := meta.Load(ctx)
	if configLoadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping Route Filters: %s", err)
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	routeFiltersSearchRequest := fabricv4.RouteFiltersSearchRequest{
		Filter: &fabricv4.SearchFilter{
			SearchAndExpression: &fabricv4.SearchAndExpression{And: []fabricv4.SearchFilterExpression{
				{
					SearchSimpleExpression: &fabricv4.SearchSimpleExpression{
						Property: "/name",
						Operator: "LIKE",
						Values:   []string{"%_PFCR"},
					},
				},
				{
					SearchSimpleExpression: &fabricv4.SearchSimpleExpression{
						Property: "/state",
						Operator: "=",
						Values:   []string{string(fabricv4.ROUTEFILTERSTATE_PROVISIONED)},
					},
				},
			},
			},
		},
	}

	fabricRouteFilters, _, err := fabric.RouteFiltersApi.SearchRouteFilters(ctx).RouteFiltersSearchRequest(routeFiltersSearchRequest).Execute()
	if err != nil {
		return fmt.Errorf("error getting route filters list for sweeping fabric route filters: %s", err)
	}

	for _, routeFilter := range fabricRouteFilters.Data {
		if sweep.IsSweepableFabricTestResource(routeFilter.GetName()) {
			log.Printf("[DEBUG] Deleting Route Filter: %s", routeFilter.GetName())
			_, resp, err := fabric.RouteFiltersApi.DeleteRouteFilterByUuid(ctx, routeFilter.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric route filter: %s", err))
			}
		}
	}

	return errors.Join(errs...)
}
