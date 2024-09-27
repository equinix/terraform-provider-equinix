package route_filter

import (
	"context"
	"errors"
	"fmt"
	"log"

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
	fabric := meta.NewFabricClientForTesting()

	name := fabricv4.ROUTEFILTERSSEARCHFILTERITEMPROPERTY_NAME
	equinixState := fabricv4.ROUTEFILTERSSEARCHFILTERITEMPROPERTY_STATE
	likeOperator := "like"
	equalOperator := "="
	limit := int32(100)
	routeFiltersSearchRequest := fabricv4.RouteFiltersSearchBase{
		Filter: &fabricv4.RouteFiltersSearchBaseFilter{
			And: []fabricv4.RouteFiltersSearchFilterItem{
				{
					Property: &name,
					Operator: &likeOperator,
					Values:   []string{"%_PFCR"},
				},
				{
					Property: &equinixState,
					Operator: &equalOperator,
					Values:   []string{string(fabricv4.ROUTEFILTERSTATE_PROVISIONED)},
				},
			},
		},
		Pagination: &fabricv4.Pagination{
			Limit: limit,
			Total: limit,
		},
	}

	fabricRouteFilters, _, err := fabric.RouteFiltersApi.SearchRouteFilters(ctx).RouteFiltersSearchBase(routeFiltersSearchRequest).Execute()
	if err != nil {
		return fmt.Errorf("error getting route filters list for sweeping fabric route filters: %s", err)
	}

	for _, routeFilter := range fabricRouteFilters.Data {
		if sweep.IsSweepableFabricTestResource(routeFilter.GetName()) {
			log.Printf("[DEBUG] Deleting Route Filter: %s", routeFilter.GetName())
			_, resp, err := fabric.RouteFiltersApi.DeleteRouteFilterByUuid(ctx, routeFilter.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric route filter: %s", err))
			}
		}
	}

	return errors.Join(errs...)
}
