package connection

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/sweep"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// sweepableConnectionSearchStates are the states a connection can be in
// while it still counts against its dependents (cloud routers, networks,
// ports). Connections stuck mid-provisioning block those deletes just like
// active ones, so they are swept too. DEPROVISIONING connections are not
// deleted again but are waited on.
//
// The connections search only accepts a subset of fabricv4.ConnectionState
// (it rejects PROVISIONED and REPROVISIONING, which GET can still return).
var sweepableConnectionSearchStates = []string{
	string(fabricv4.CONNECTIONSTATE_ACTIVE),
	string(fabricv4.CONNECTIONSTATE_DEPROVISIONING),
	string(fabricv4.CONNECTIONSTATE_DRAFT),
	string(fabricv4.CONNECTIONSTATE_FAILED),
	string(fabricv4.CONNECTIONSTATE_PENDING),
	string(fabricv4.CONNECTIONSTATE_PROVISIONING),
}

// deprovisionPendingStates are the states GET may report while a swept
// connection is on its way to DEPROVISIONED.
var deprovisionPendingStates = append([]string{
	string(fabricv4.CONNECTIONSTATE_PROVISIONED),
	string(fabricv4.CONNECTIONSTATE_REPROVISIONING),
}, sweepableConnectionSearchStates...)

// connectionSweepTimeout bounds the whole connection sweep. Connections are
// swept concurrently, so this is the budget for the slowest one, not for
// each of them in turn.
const connectionSweepTimeout = 15 * time.Minute

// connectionSweepPollInterval is how often deprovisioning is polled.
var connectionSweepPollInterval = 10 * time.Second

// remainingSweepTime returns how long until ctx's deadline, less a margin so
// that retry/wait timeouts fire (and report the last state) before ctx does.
func remainingSweepTime(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return connectionSweepTimeout
	}
	return time.Until(deadline) - time.Second
}

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
					Values:   sweepableConnectionSearchStates,
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

	ctx, cancel := context.WithTimeout(ctx, connectionSweepTimeout)
	defer cancel()

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	for _, connection := range fabricConnections.Data {
		if !sweep.IsSweepableFabricTestResource(connection.GetName()) {
			continue
		}
		wg.Add(1)
		go func(connection fabricv4.Connection) {
			defer wg.Done()
			if err := sweepConnection(ctx, fabric, connection); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(connection)
	}
	wg.Wait()

	return errors.Join(errs...)
}

// sweepConnection deletes a connection and waits for it to finish
// deprovisioning, so that sweepers depending on this one (cloud routers,
// networks) don't fail with "active connections" errors. Deletes rejected
// because the connection is mid-operation are retried until ctx expires.
func sweepConnection(ctx context.Context, fabric *fabricv4.APIClient, connection fabricv4.Connection) error {
	name, uuid := connection.GetName(), connection.GetUuid()

	if connection.GetState() != fabricv4.CONNECTIONSTATE_DEPROVISIONING {
		log.Printf("[DEBUG] Deleting Connection: %s", name)
		alreadyDeleted := false
		err := retry.RetryContext(ctx, remainingSweepTime(ctx), func() *retry.RetryError {
			_, resp, err := fabric.ConnectionsApi.DeleteConnectionByUuid(ctx, uuid).Execute()
			switch {
			case isConnectionAlreadyDeleted(err):
				alreadyDeleted = true
				return nil
			case isConnectionInTransientState(err):
				log.Printf("[DEBUG] Connection %s (%s) is in a transient state, retrying delete", name, uuid)
				return retry.RetryableError(err)
			}
			if err := equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err); err != nil {
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error deleting fabric connection %s (%s): %s", name, uuid, err)
		}
		if alreadyDeleted {
			log.Printf("[DEBUG] Connection %s (%s) already deleted", name, uuid)
			return nil
		}
	}

	log.Printf("[DEBUG] Waiting for Connection %s (%s) to deprovision", name, uuid)
	_, err := waitForConnection(ctx, fabric, uuid,
		deprovisionPendingStates,
		[]string{string(fabricv4.CONNECTIONSTATE_DEPROVISIONED)},
		remainingSweepTime(ctx),
		withNotFoundAs(string(fabricv4.CONNECTIONSTATE_DEPROVISIONED)),
		withPollInterval(connectionSweepPollInterval),
	)
	if err != nil {
		return fmt.Errorf("error waiting for fabric connection %s (%s) to deprovision: %s", name, uuid, err)
	}
	return nil
}
