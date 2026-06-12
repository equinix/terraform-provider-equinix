package routeaggregation

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
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Route Aggregations: %s", err)
	}
	configLoadErr := meta.Load(ctx)
	if configLoadErr != nil {
		return fmt.Errorf("error loading configuration for sweeping Route Aggregations: %s", err)
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	name := fabricv4.ROUTEFILTERSSEARCHFILTERITEMPROPERTY_NAME
	equinixState := fabricv4.ROUTEFILTERSSEARCHFILTERITEMPROPERTY_STATE
	likeOperator := string(fabricv4.EXPRESSIONOPERATOR_LIKE)
	equalOperator := "="
	pfcrRouteAggregationsSearch := fabricv4.RouteAggregationsSearchBase{
		Filter: &fabricv4.RouteAggregationsSearchBaseFilter{
			And: []fabricv4.RouteAggregationsSearchFilterItem{
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
	}

	pfcrRouteAggr, _, err := fabric.RouteAggregationsApi.SearchRouteAggregations(ctx).RouteAggregationsSearchBase(pfcrRouteAggregationsSearch).Execute()
	if err != nil {
		return fmt.Errorf("error getting route aggregations list for sweeping fabric route aggregations: %s", err)
	}

	pnfvRouteAggregationsSearch := fabricv4.RouteAggregationsSearchBase{
		Filter: &fabricv4.RouteAggregationsSearchBaseFilter{
			And: []fabricv4.RouteAggregationsSearchFilterItem{
				{
					Property: &name,
					Operator: &likeOperator,
					Values:   []string{"%_PNFV"},
				},
				{
					Property: &equinixState,
					Operator: &equalOperator,
					Values:   []string{string(fabricv4.ROUTEFILTERSTATE_PROVISIONED)},
				},
			},
		},
	}

	pnfvRouteAggr, _, err := fabric.RouteAggregationsApi.SearchRouteAggregations(ctx).RouteAggregationsSearchBase(pnfvRouteAggregationsSearch).Execute()
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
