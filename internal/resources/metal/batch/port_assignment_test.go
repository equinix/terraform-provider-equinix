package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/config"
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
			true,
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

func TestVlanBatchWithMockAPI(t *testing.T) {
	batchID := "batch-id"
	cases := []struct {
		name         string
		responses    []json.Marshaler
		statusCodes  []int
		expectations func(*testing.T, *metalv1.PortVlanAssignmentBatch, *http.Response, error)
	}{
		{
			name: "successful-execution",
			responses: []json.Marshaler{
				metalv1.PortVlanAssignmentBatch{Id: &batchID},
				metalv1.PortVlanAssignmentBatch{Id: &batchID, State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED.Ptr()},
			},
			statusCodes: []int{http.StatusCreated, http.StatusOK},
			expectations: func(t *testing.T, batch *metalv1.PortVlanAssignmentBatch, _ *http.Response, err error) {
				assert.NoError(t, err)
				assert.Equal(t, metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED, batch.GetState())
			},
		},
		{
			name: "failure-after-excuting",
			responses: []json.Marshaler{
				metalv1.PortVlanAssignmentBatch{Id: &batchID},
				metalv1.PortVlanAssignmentBatch{Id: &batchID, State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_QUEUED.Ptr()},
				metalv1.PortVlanAssignmentBatch{Id: &batchID, State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_IN_PROGRESS.Ptr()},
				metalv1.PortVlanAssignmentBatch{Id: &batchID, State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED.Ptr()},
			},
			statusCodes: []int{http.StatusCreated, http.StatusOK, http.StatusOK, http.StatusOK},
			expectations: func(t *testing.T, batch *metalv1.PortVlanAssignmentBatch, _ *http.Response, err error) {
				assert.ErrorContains(t, err, fmt.Sprintf("vlan assignment batch %s provisioning failed", "batch-id"))
				assert.Equal(t, metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED, batch.GetState())

			},
		},
		{
			name:        "failed-to-create-batch",
			responses:   []json.Marshaler{nil},
			statusCodes: []int{http.StatusInternalServerError},
			expectations: func(t *testing.T, _ *metalv1.PortVlanAssignmentBatch, _ *http.Response, err error) {
				assert.ErrorContains(t, err, "failed to create batch")
			},
		},
		{
			name: "failed-to-poll-batch",
			responses: []json.Marshaler{
				metalv1.PortVlanAssignmentBatch{Id: &batchID},
				nil,
			},
			statusCodes: []int{http.StatusCreated, http.StatusInternalServerError},
			expectations: func(t *testing.T, _ *metalv1.PortVlanAssignmentBatch, _ *http.Response, err error) {
				assert.ErrorContains(t, err, "could not be polled")
			},
		},
	}

	ctx := context.Background()
	vb := NewVlanBatch("port_foo")
	vb.AddAssignment("vlan1")
	vb.SetRetryTimeouts(time.Millisecond, time.Millisecond)

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			callIndex := 0
			handler := func(w http.ResponseWriter, _ *http.Request) {
				resp := tc.responses[callIndex]
				status := tc.statusCodes[callIndex]
				switch status {
				case http.StatusOK, http.StatusCreated:
					body, err := resp.MarshalJSON()
					assert.NoError(tt, err, "failed to marshal json call %d", callIndex)
					w.Header().Add("content-type", "application/json")
					w.WriteHeader(status)
					_, err = w.Write(body)
					assert.NoError(tt, err, "failed to write body for call %d", callIndex)
				default:
					w.Header().Add("content-type", "application/json")
					w.WriteHeader(status)
				}

				callIndex++
			}
			mockAPI := httptest.NewServer(http.HandlerFunc(handler))
			defer mockAPI.Close()
			meta := &config.Config{
				BaseURL: mockAPI.URL,
				Token:   "superNotRealToken",
			}

			err := meta.Load(ctx)
			assert.NoError(tt, err)

			client := meta.NewMetalClientForTesting()

			batchResp, httpResp, err := vb.Execute(ctx, client)

			tc.expectations(tt, batchResp, httpResp, err)
		})
	}

}

func TestVlanBatchNoAssignments(t *testing.T) {
	vb := NewVlanBatch("port_foo")

	// empty client because we shouldn't hit it at all
	client := &metalv1.APIClient{}
	_, _, err := vb.Execute(context.Background(), client)

	// the previous batch assignment logic returned nil so that this
	// could be a no op for port updates that weren't updating
	// vlan assignments
	assert.NoError(t, err)
}
