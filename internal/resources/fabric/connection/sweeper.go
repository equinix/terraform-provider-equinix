package connection

import (
	"context"
	"errors"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/sweep"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func AddTestSweeper() {
	resource.AddTestSweepers("equinix_fabric_connection", &resource.Sweeper{
		Name:         "equinix_fabric_connection",
		Dependencies: []string{},
		F:            testSweepConnections,
	})
}

func testSweepConnections(region string) error {
	var errs []error
	log.Printf("[DEBUG] Sweeping Fabric Connections")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Conections: %s", err)
	}
	meta.Load(ctx)
	fabric := meta.NewFabricClientForTesting()

	name := fabricv4.SEARCHFIELDNAME_NAME
	equinixStatus := fabricv4.SEARCHFIELDNAME_OPERATION_EQUINIX_STATUS
	likeOperator := fabricv4.EXPRESSIONOPERATOR_LIKE
	equalOperator := fabricv4.EXPRESSIONOPERATOR_EQUAL
	limit := int32(100)
	connectionsSearchRequest := fabricv4.SearchRequest{
		Filter: &fabricv4.Expression{
			And: []fabricv4.Expression{
				{
					Property: &name,
					Operator: &likeOperator,
					Values:   []string{"PFCR"},
				},
				{
					Property: &equinixStatus,
					Operator: &equalOperator,
					Values:   []string{"PROVISIONED"},
				},
			},
		},
		Pagination: &fabricv4.PaginationRequest{
			Limit: &limit,
		},
	}

	fabricConnections, _, err := fabric.ConnectionsApi.SearchConnections(ctx).SearchRequest(connectionsSearchRequest).Execute()
	if err != nil {
		return fmt.Errorf("error getting connections list for sweeping fabric connections: %s", err)
	}

	for _, connection := range fabricConnections.Data {
		if strings.HasSuffix(connection.GetName(), "_PFCR") {
			log.Printf("[DEBUG] Deleting Connection: %s", connection.GetName())
			_, resp, err := fabric.ConnectionsApi.DeleteConnectionByUuid(ctx, connection.GetUuid()).Execute()
			if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
				errs = append(errs, fmt.Errorf("error deleting fabric connection: %s", err))
			}
		}
	}

	return errors.Join(errs...)
}
