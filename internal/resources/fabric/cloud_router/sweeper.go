package cloud_router

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_cloud_router", &resource.Sweeper{
		Name:         "equinix_fabric_cloud_router",
		Dependencies: []string{},
		F:            testSweepCloudRouters,
	})
}

func testSweepCloudRouters(region string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Fabric Cloud Routers")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Fabric Cloud Routers: %s", err)
	}
	err = meta.Load(ctx)
	if err != nil {
		log.Printf("Error loading meta: %v", err)
		return err
	}
	fabric := meta.NewFabricClientForTesting()

	equinixStatus := "/state"
	equalOperator := string(fabricv4.EXPRESSIONOPERATOR_EQUAL)
	limit := int32(100)
	cloudRouterSearchRequest := fabricv4.CloudRouterSearchRequest{
		Filter: &fabricv4.CloudRouterFilters{
			And: []fabricv4.CloudRouterFilter{
				{
					CloudRouterSimpleExpression: &fabricv4.CloudRouterSimpleExpression{
						Property: &equinixStatus,
						Operator: &equalOperator,
						Values:   []string{string(fabricv4.EQUINIXSTATUS_PROVISIONED)},
					},
				},
			},
		},
		Pagination: &fabricv4.PaginationRequest{
			Limit: &limit,
		},
	}

	fabricCloudRouters, _, err := fabric.CloudRoutersApi.SearchCloudRouters(ctx).CloudRouterSearchRequest(cloudRouterSearchRequest).Execute()
	if err != nil {
		return fmt.Errorf("error getting cloud router list for sweeping fabric cloud routers: %s", err)
	}

	for _, cloudRouter := range fabricCloudRouters.Data {
		if sweep.IsSweepableFabricTestResource(cloudRouter.GetName()) {
			log.Printf("[DEBUG] Deleting Cloud Routers: %s", cloudRouter.GetName())
			resp, err := fabric.CloudRoutersApi.DeleteCloudRouterByUuid(ctx, cloudRouter.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric Cloud Router: %s", err))
			}
		}
	}

	return errors.Join(errs...)
}
