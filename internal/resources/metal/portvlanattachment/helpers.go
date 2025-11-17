package portvlanattachment

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type vlanAssignmentBatchEntry = metalv1.PortVlanAssignmentBatchCreateInputVlanAssignmentsInner

// vlanBatch helps to build the proper batch vlan assignment payload
// nothing is sent to the API until `createAndWaitForBatch` is called.
type vlanBatch struct {
	portID      string
	assignments []vlanAssignment
	result      *http.Response
}

func newVlanBatch(portID string) *vlanBatch {
	return &vlanBatch{
		portID:      portID,
		assignments: []vlanAssignment{},
	}
}

func (b *vlanBatch) assign(vlanID string) {
	state := metalv1.PORTVLANASSIGNMENTBATCHVLANASSIGNMENTSINNERSTATE_ASSIGNED.Ptr()
	assignments := append(b.assignments, vlanAssignment{vlanID: vlanID, state: state})
	b.assignments = assignments
}

func (b *vlanBatch) unassign(vlanID string) {
	state := metalv1.PORTVLANASSIGNMENTBATCHVLANASSIGNMENTSINNERSTATE_UNASSIGNED.Ptr()
	assignments := append(b.assignments, vlanAssignment{vlanID: vlanID, state: state})
	b.assignments = assignments
}

func (b vlanBatch) toBatchCreateInput() metalv1.PortVlanAssignmentBatchCreateInput {
	createInput := metalv1.NewPortVlanAssignmentBatchCreateInput()
	assignments := []metalv1.PortVlanAssignmentBatchCreateInputVlanAssignmentsInner{}
	for _, assignment := range b.assignments {
		assignments = append(assignments, assignment.toVlanAssignmentBatchEntry())
	}

	createInput.SetVlanAssignments(assignments)

	return *createInput
}

func (b vlanBatch) httpResponse() *http.Response {
	return b.result
}

func (b *vlanBatch) createAndWaitForBatch(ctx context.Context, start time.Time, client *metalv1.APIClient) error {
	batchReq := client.PortsApi.CreatePortVlanAssignmentBatch(ctx, b.portID)

	batch, resp, err := batchReq.PortVlanAssignmentBatchCreateInput(b.toBatchCreateInput()).Execute()
	if err != nil {
		b.result = resp
		return fmt.Errorf("failed to create batch for vlan assignment: %w", err)
	}

	deadline, _ := ctx.Deadline()

	ctxTimeout := deadline.Sub(start)

	poller := retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{string(metalv1.PORTVLANASSIGNMENTBATCHSTATE_QUEUED), string(metalv1.PORTVLANASSIGNMENTBATCHSTATE_IN_PROGRESS)},
		Target:     []string{string(metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED), string(metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED)},
		MinTimeout: 5 * time.Second,
		Timeout:    ctxTimeout - time.Since(start) - 30*time.Second,
		Refresh: func() (result any, state string, err error) {
			batchResp, resp, err := client.PortsApi.FindPortVlanAssignmentBatchByPortIdAndBatchId(ctx, b.portID, batch.GetId()).Execute()
			switch batchResp.GetState() {
			case metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED:
				return batchResp, string(metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED),
					fmt.Errorf("vlan assignment batch %s provisioning failed: %s", batchResp.GetId(), strings.Join(batchResp.ErrorMessages, "; "))
			case metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED:
				return batchResp, string(metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED), nil
			default:
				if err != nil {
					b.result = resp
					return batchResp, "", fmt.Errorf("vlan assignment batch %s could not be polled: %w", batch.GetId(), err)
				}
				return batchResp, string(batchResp.GetState()), err
			}

		},
	}

	if _, err = poller.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("vlan assignment batch %s is not complete after timeout: %w", batch.GetId(), err)
	}

	return nil
}

type vlanAssignment struct {
	vlanID string
	state  *metalv1.PortVlanAssignmentBatchVlanAssignmentsInnerState
}

func (vla vlanAssignment) toVlanAssignmentBatchEntry() vlanAssignmentBatchEntry {
	return vlanAssignmentBatchEntry{Vlan: &vla.vlanID, State: vla.state}
}
