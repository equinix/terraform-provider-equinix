package route_aggregation

import (
	"context"
	"errors"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"log"
	"net/http"
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
	likeOperator := "like"
	equalOperator := "="
	limit := int32(100)
	routeAggregationsSearch := fabricv4.RouteAggregationsSearchBase{
		Filter: &fabricv4.RouteAggregationsSearchBaseFilter{
			And: []fabricv4.RouteAggregationsSearchFilterItem{
				{
					Property: &name,
					Operator: &likeOperator,
					Values:   []string{"%_PFCR", "%_PNFV"},
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

	routeAggregation, _, err := fabric.RouteAggregationsApi.SearchRouteAggregations(ctx).RouteAggregationsSearchBase(routeAggregationsSearch).Execute()
	if err != nil {
		return fmt.Errorf("error getting streams list for sweeping fabric route aggregations: %s", err)
	}

	for _, ra := range routeAggregation.GetData() {
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
