package servicetoken

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

// AddTestSweeper /* AddTestSweeper registers a test sweeper for cleaning up resources
func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_service_token", &resource.Sweeper{
		Name:         "equinix_fabric_service_token",
		Dependencies: []string{},
		F:            testSweepServiceToken,
	})
}

func testSweepServiceToken(_ string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Fabric Service Tokens")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Service Tokens: %s", err)
	}
	err = meta.Load(ctx)
	if err != nil {
		log.Printf("Error loading meta: %v", err)
		return err
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	state := fabricv4.SERVICETOKENSEARCHFIELDNAME_STATE
	limit := int32(100)
	offset := int32(0)
	equalOperator := fabricv4.SERVICETOKENSEARCHEXPRESSIONOPERATOR_EQUAL

	serviceTokensRequest := fabricv4.ServiceTokenSearchRequest{
		Filter: &fabricv4.ServiceTokenSearchExpression{
			And: []fabricv4.ServiceTokenSearchExpression{
				{
					Property: &state,
					Operator: &equalOperator,
					Values:   []string{string(fabricv4.SERVICETOKENSTATE_ACTIVE)},
				},
			},
		},
		Pagination: &fabricv4.PaginationRequest{
			Limit:  &limit,
			Offset: &offset,
		},
	}

	fabriceServiceTokens, _, err := fabric.ServiceTokensApi.SearchServiceTokens(ctx).ServiceTokenSearchRequest(serviceTokensRequest).Execute()
	if err != nil {
		return fmt.Errorf("error getting service tokens list for sweeping fabric service tokens: %s", err)
	}

	for _, serviceToken := range fabriceServiceTokens.Data {
		if sweep.IsSweepableFabricTestResource(serviceToken.GetName()) {
			log.Printf("[DEBUG] Deleting serviceToken: %s", serviceToken.GetName())
			_, resp, err := fabric.ServiceTokensApi.DeleteServiceTokenByUuid(ctx, serviceToken.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric serviceToken: %s", err))
			}
		}
	}

	return errors.Join(errs...)
}
