package portvlanattachment

import (
	"testing"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/stretchr/testify/assert"
)

func TestPortVlanAttachmentVlanBatch(t *testing.T) {
	batch := newVlanBatch("foo")
	batch.assign("vlan1")
	batch.unassign("vlan2")

	createInput := batch.toBatchCreateInput()

	assert.Equal(t, 2, len(createInput.VlanAssignments))
	matchAssignment(t, createInput, "vlan1", "assigned")
	matchAssignment(t, createInput, "vlan2", "unassigned")
}

func matchAssignment(t *testing.T, createInput metalv1.PortVlanAssignmentBatchCreateInput, vlanID string, state string) {
	var assignment metalv1.PortVlanAssignmentBatchCreateInputVlanAssignmentsInner

	for _, asgn := range createInput.VlanAssignments {
		if asgn.GetVlan() == vlanID {
			assignment = asgn
			break
		}
	}

	assert.NotEmpty(t, assignment, "assignment not found")
	assert.Equal(t, vlanID, assignment.GetVlan())
	assert.Equal(t, state, string(assignment.GetState()))
}
