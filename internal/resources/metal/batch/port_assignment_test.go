package batch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVlanBatch(t *testing.T) {
	vb := NewVlanBatch("port_foo")
	vb.AddAssignment("vlan1")
	vb.AddNativeAssignment("vlan2")
	vb.RemoveAssignment("vlan3")

	createInput := vb.toBatchCreateInput()
	assignments := createInput.GetVlanAssignments()
	assert.Equal(t, 3, len(assignments))
	expectedAssignments := []struct {
		vlan     string
		assigned bool
		native   bool
	}{
		{
			"vlan1",
			false,
			false,
		},
		{
			"vlan2",
			true,
			true,
		},
		{
			"vlan3",
			false,
			false,
		},
	}

	for i, assignment := range createInput.GetVlanAssignments() {
		expectation := expectedAssignments[i]
		assert.Equal(t, expectation.vlan, assignment.GetVlan(), "vlan does not match expectation case %d", i)
		assert.Equal(t, expectation.assigned, assignment.GetState() == "assigned", "state does not match expectation case %d", i)
		assert.Equal(t, expectation.native, assignment.GetNative(), "native status does not match expectation case %d", i)
	}
}
