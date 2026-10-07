package connection

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// sweepableConnectionStates are the states a connection can be in while it
// still counts against its dependents (cloud routers, networks, ports).
// Connections stuck mid-provisioning block those deletes just like
// provisioned ones, so they are swept too. DEPROVISIONING connections are
// not deleted again but are waited on.
var sweepableConnectionStates = []string{
	string(fabricv4.CONNECTIONSTATE_ACTIVE),
	string(fabricv4.CONNECTIONSTATE_DEPROVISIONING),
	string(fabricv4.CONNECTIONSTATE_DRAFT),
	string(fabricv4.CONNECTIONSTATE_FAILED),
	string(fabricv4.CONNECTIONSTATE_PENDING),
	string(fabricv4.CONNECTIONSTATE_PROVISIONED),
	string(fabricv4.CONNECTIONSTATE_PROVISIONING),
	string(fabricv4.CONNECTIONSTATE_REPROVISIONING),
}

const connectionDeprovisionTimeout = 15 * time.Minute

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
	meta, err := testinghelpers.GetConfigForFabric()
	if err != nil {
		return fmt.Errorf("error getting configuration for sweeping Connections: %s", err)
	}
	err = meta.Load(ctx)
	if err != nil {
		log.Printf("Error loading meta: %v", err)
		return err
	}
	fabric := meta.NewFabricClientForTesting(ctx)

	name := fabricv4.SEARCHFIELDNAME_NAME
	state := fabricv4.SEARCHFIELDNAME_STATE
	likeOperator := fabricv4.EXPRESSIONOPERATOR_LIKE
	inOperator := fabricv4.EXPRESSIONOPERATOR_IN
	limit := int32(100)
	connectionsSearchRequest := fabricv4.SearchRequest{
		Filter: &fabricv4.Expression{
			And: []fabricv4.Expression{
				{
					Property: &name,
					Operator: &likeOperator,
					Values:   sweep.FabricTestResourceSuffixes,
				},
				{
					Property: &state,
					Operator: &inOperator,
					Values:   sweepableConnectionStates,
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

	var deleted []string
	for _, connection := range fabricConnections.Data {
		if !sweep.IsSweepableFabricTestResource(connection.GetName()) {
			continue
		}
		if connection.GetState() == fabricv4.CONNECTIONSTATE_DEPROVISIONING {
			deleted = append(deleted, connection.GetUuid())
			continue
		}
		log.Printf("[DEBUG] Deleting Connection: %s", connection.GetName())
		_, resp, err := fabric.ConnectionsApi.DeleteConnectionByUuid(ctx, connection.GetUuid()).Execute()
		if isConnectionAlreadyDeleted(err) {
			log.Printf("[DEBUG] Connection %s (%s) already deleted", connection.GetName(), connection.GetUuid())
			continue
		}
		if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) != nil {
			errs = append(errs, fmt.Errorf("error deleting fabric connection: %s", err))
			continue
		}
		deleted = append(deleted, connection.GetUuid())
	}

	// Connection deletes are asynchronous. Sweepers that depend on this one
	// (cloud routers, networks) fail with "active connections" errors unless
	// the connections have finished deprovisioning first.
	for _, uuid := range deleted {
		if err := waitForConnectionDeprovisioned(ctx, fabric, uuid); err != nil {
			errs = append(errs, fmt.Errorf("error waiting for fabric connection %s to deprovision: %s", uuid, err))
		}
	}

	return errors.Join(errs...)
}

// isConnectionAlreadyDeleted reports whether err is the API's
// EQ-3142509 "Connection already deleted" validation error.
func isConnectionAlreadyDeleted(err error) bool {
	var genericError *fabricv4.GenericOpenAPIError
	if !errors.As(err, &genericError) {
		return false
	}
	fabricErrs, ok := genericError.Model().([]fabricv4.Error)
	return ok && equinix_errors.HasErrorCode(fabricErrs, "EQ-3142509")
}

func waitForConnectionDeprovisioned(ctx context.Context, fabric *fabricv4.APIClient, uuid string) error {
	stateConf := &retry.StateChangeConf{
		Pending: sweepableConnectionStates,
		Target: []string{
			string(fabricv4.CONNECTIONSTATE_DEPROVISIONED),
		},
		Refresh: func() (any, string, error) {
			conn, resp, err := fabric.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return uuid, string(fabricv4.CONNECTIONSTATE_DEPROVISIONED), nil
				}
				return nil, "", equinix_errors.FormatFabricError(err)
			}
			return conn, string(conn.GetState()), nil
		},
		Timeout:    connectionDeprovisionTimeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
