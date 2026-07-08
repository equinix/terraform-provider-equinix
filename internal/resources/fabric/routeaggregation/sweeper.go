package routeaggregation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/sweep"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_route_aggregation", &resource.Sweeper{
		Name:         "equinix_fabric_route_aggregation",
		Dependencies: []string{},
		F:            testSweepRouteAggregations,
	})
}

func testSweepRouteAggregations(_ string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Route Aggregations")
	ctx := context.Background()
	meta, err := testinghelpers.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Route Aggregations: %s", err)
	}
	configLoadErr := meta.Load(ctx)
	if configLoadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping Route Aggregations: %s", err)
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	pfcrRouteAggregationsSearch := fabricv4.RouteAggregationsSearchRequest{
		Filter: &fabricv4.SearchFilter{
			SearchAndExpression: &fabricv4.SearchAndExpression{
				And: []fabricv4.SearchFilterExpression{
					{
						SearchSimpleExpression: &fabricv4.SearchSimpleExpression{
							Property: "/name",
							Operator: "LIKE",
							Values:   []string{"%_PFCR"},
						}},
					{
						SearchSimpleExpression: &fabricv4.SearchSimpleExpression{
							Property: "/state",
							Operator: "=",
							Values:   []string{string(fabricv4.ROUTEFILTERSTATE_PROVISIONED)},
						}},
				},
			},
		},
	}

	pfcrRouteAggr, _, err := fabric.RouteAggregationsApi.SearchRouteAggregations(ctx).RouteAggregationsSearchRequest(pfcrRouteAggregationsSearch).Execute()
	if err != nil {
		return fmt.Errorf("error getting route aggregations list for sweeping fabric route aggregations: %s", err)
	}

	pnfvRouteAggregationsSearch := fabricv4.RouteAggregationsSearchRequest{
		Filter: &fabricv4.SearchFilter{
			SearchAndExpression: &fabricv4.SearchAndExpression{
				And: []fabricv4.SearchFilterExpression{
					{
						SearchSimpleExpression: &fabricv4.SearchSimpleExpression{
							Property: "/name",
							Operator: "LIKE",
							Values:   []string{"%_PNFV"},
						}},
					{
						SearchSimpleExpression: &fabricv4.SearchSimpleExpression{
							Property: "/state",
							Operator: "=",
							Values:   []string{string(fabricv4.ROUTEFILTERSTATE_PROVISIONED)},
						}},
				},
			},
		},
	}

	pnfvRouteAggr, _, err := fabric.RouteAggregationsApi.SearchRouteAggregations(ctx).RouteAggregationsSearchRequest(pnfvRouteAggregationsSearch).Execute()
	if err != nil {
		return fmt.Errorf("error getting route aggregations list for sweeping fabric route aggregations: %s", err)
	}

	for _, ra := range append(pfcrRouteAggr.GetData(), pnfvRouteAggr.GetData()...) {
		if sweep.IsSweepableFabricTestResource(ra.GetName()) {
			log.Printf("[DEBUG] Deleting route aggregation: %s", ra.GetName())
			_, resp, err := fabric.RouteAggregationsApi.DeleteRouteAggregationByUuid(ctx, ra.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric route aggregation: %s", err))
			}
		}
	}

	return errors.Join(errs...)
}
