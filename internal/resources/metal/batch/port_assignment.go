package batch

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

type VlanBatch struct {
	portID      string
	assignments []assignment
}

func NewVlanBatch(portID string) *VlanBatch {
	return &VlanBatch{
		portID:      portID,
		assignments: []assignment{},
	}
}

func (vb *VlanBatch) AddAssignment(vlanID string) {
	state := metalv1.PORTVLANASSIGNMENTBATCHVLANASSIGNMENTSINNERSTATE_ASSIGNED.Ptr()
	assignments := append(vb.assignments, vlanAssignment{vlanID: vlanID, state: state})
	vb.assignments = assignments

}

func (vb *VlanBatch) AddNativeAssignment(vlanID string) {
	state := metalv1.PORTVLANASSIGNMENTBATCHVLANASSIGNMENTSINNERSTATE_ASSIGNED.Ptr()
	native := true
	assignments := append(vb.assignments, vlanAssignment{vlanID: vlanID, state: state, native: &native})
	vb.assignments = assignments
}

func (vb *VlanBatch) RemoveAssignment(vlanID string) {
	state := metalv1.PORTVLANASSIGNMENTBATCHVLANASSIGNMENTSINNERSTATE_UNASSIGNED.Ptr()
	assignments := append(vb.assignments, vlanAssignment{vlanID: vlanID, state: state})
	vb.assignments = assignments

}

func (vb *VlanBatch) toBatchCreateInput() metalv1.PortVlanAssignmentBatchCreateInput {
	createInput := metalv1.NewPortVlanAssignmentBatchCreateInput()
	assignments := []metalv1.PortVlanAssignmentBatchCreateInputVlanAssignmentsInner{}
	for _, assignment := range vb.assignments {
		assignments = append(assignments, assignment.toVlanAssignmentBatchEntry())
	}

	createInput.SetVlanAssignments(assignments)

	return *createInput
}

func (vb *VlanBatch) Execute(ctx context.Context, client *metalv1.APIClient) (*metalv1.PortVlanAssignmentBatch, *http.Response, error) {
	start := time.Now()
	batchReq := client.PortsApi.CreatePortVlanAssignmentBatch(ctx, vb.portID)

	batch, resp, err := batchReq.PortVlanAssignmentBatchCreateInput(vb.toBatchCreateInput()).Execute()
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create batch for vlan assignment: %w", err)
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
			batchResp, resp, err := client.PortsApi.FindPortVlanAssignmentBatchByPortIdAndBatchId(ctx, vb.portID, batch.GetId()).Execute()
			switch batchResp.GetState() {
			case metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED:
				return batchResp, string(metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED),
					fmt.Errorf("vlan assignment batch %s provisioning failed: %s", batchResp.GetId(), strings.Join(batchResp.ErrorMessages, "; "))
			case metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED:
				return batchResp, string(metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED), nil
			default:
				if err != nil {
					return resp, "", fmt.Errorf("vlan assignment batch %s could not be polled: %w", batch.GetId(), err)
				}
				return batchResp, string(batchResp.GetState()), err
			}

		},
	}

	res, err := poller.WaitForStateContext(ctx)
	if err != nil {
		switch value := res.(type) {
		case *http.Response:
			return nil, value, err
		case *metalv1.PortVlanAssignmentBatch:
			return value, nil, err
		default:
			return nil, nil, err
		}
	}

	return res.(*metalv1.PortVlanAssignmentBatch), nil, err
}

type assignment interface {
	toVlanAssignmentBatchEntry() vlanAssignmentBatchEntry
}

type vlanAssignment struct {
	vlanID string
	state  *metalv1.PortVlanAssignmentBatchVlanAssignmentsInnerState
	native *bool
}

func (vla vlanAssignment) toVlanAssignmentBatchEntry() vlanAssignmentBatchEntry {
	res := vlanAssignmentBatchEntry{Vlan: &vla.vlanID, State: vla.state}

	if vla.native != nil {
		res.Native = vla.native
	}

	return res
}
