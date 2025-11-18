package batch

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	ctx := context.Background()
	vb := NewVlanBatch("port_foo")
	vb.AddAssignment("vlan1")

	cases := []struct {
		name         string
		handler      http.HandlerFunc
		expectations func(*testing.T, *metalv1.PortVlanAssignmentBatch, *http.Response, error)
	}{
		{
			name: "successful-execution",
			handler: func() http.HandlerFunc {
				batchID := "batch-id"
				invocations := []http.HandlerFunc{
					func(resp http.ResponseWriter, req *http.Request) {
						b := metalv1.PortVlanAssignmentBatch{
							Id: &batchID,
						}

						if strings.Contains(req.URL.String(), "port_foo") {
							resp.Header().Add("content-type", "application/json")
							bytes, err := b.MarshalJSON()
							if err != nil {
								t.Fatalf("failed to marshal json: %v", err)
							}
							resp.Write(bytes)
						}
					},
					func(resp http.ResponseWriter, req *http.Request) {
						b := metalv1.PortVlanAssignmentBatch{
							Id: &batchID,
							State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED.Ptr(),
						}

						if strings.Contains(req.URL.String(), "port_foo") {
							resp.Header().Add("content-type", "application/json")
							bytes, err := b.MarshalJSON()
							if err != nil {
								t.Fatalf("failed to marshal json: %v", err)
							}
							resp.Write(bytes)
						}
					},

				}

				index := 0

				return func(resp http.ResponseWriter, req *http.Request) {
					invocations[index](resp, req)
					index++
				}
			}(),
			expectations: func(t *testing.T, batch *metalv1.PortVlanAssignmentBatch, _ *http.Response, err error) {
				assert.NoError(t, err)
				assert.Equal(t, metalv1.PORTVLANASSIGNMENTBATCHSTATE_COMPLETED, batch.GetState())
			},
		},
{
			name: "failure-after-excuting",
			handler: func() http.HandlerFunc {
				batchID := "batch-id"
				invocations := []http.HandlerFunc{
					func(resp http.ResponseWriter, req *http.Request) {
						b := metalv1.PortVlanAssignmentBatch{
							Id: &batchID,
						}

						if strings.Contains(req.URL.String(), "port_foo") {
							resp.Header().Add("content-type", "application/json")
							bytes, err := b.MarshalJSON()
							if err != nil {
								t.Fatalf("failed to marshal json: %v", err)
							}
							resp.Write(bytes)
						}
					},
					func(resp http.ResponseWriter, req *http.Request) {
						b := metalv1.PortVlanAssignmentBatch{
							Id: &batchID,
							State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_QUEUED.Ptr(),
						}

						if strings.Contains(req.URL.String(), "port_foo") {
							resp.Header().Add("content-type", "application/json")
							bytes, err := b.MarshalJSON()
							if err != nil {
								t.Fatalf("failed to marshal json: %v", err)
							}
							resp.Write(bytes)
						}
					},
					func(resp http.ResponseWriter, req *http.Request) {
						b := metalv1.PortVlanAssignmentBatch{
							Id: &batchID,
							State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_IN_PROGRESS.Ptr(),
						}

						if strings.Contains(req.URL.String(), "port_foo") {
							resp.Header().Add("content-type", "application/json")
							bytes, err := b.MarshalJSON()
							if err != nil {
								t.Fatalf("failed to marshal json: %v", err)
							}
							resp.Write(bytes)
						}
					},
					func(resp http.ResponseWriter, req *http.Request) {
						b := metalv1.PortVlanAssignmentBatch{
							Id: &batchID,
							State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED.Ptr(),
						}

						if strings.Contains(req.URL.String(), "port_foo") {
							resp.Header().Add("content-type", "application/json")
							bytes, err := b.MarshalJSON()
							if err != nil {
								t.Fatalf("failed to marshal json: %v", err)
							}
							resp.Write(bytes)
						}
					},
				}

				index := 0

				return func(resp http.ResponseWriter, req *http.Request) {
					invocations[index](resp, req)
					index++
				}
			}(),
			expectations: func(t *testing.T, batch *metalv1.PortVlanAssignmentBatch, _ *http.Response, err error) {
				assert.ErrorContains(t, err, fmt.Sprintf("vlan assignment batch %s provisioning failed", "batch-id"))
				assert.Equal(t, metalv1.PORTVLANASSIGNMENTBATCHSTATE_FAILED, batch.GetState())
				
			},
		},
		{
			name: "failed-to-create-batch",
			handler: func() http.HandlerFunc {
				i := 0
				return func(resp http.ResponseWriter, req *http.Request) {
					switch i {
					case 0:
						resp.WriteHeader(http.StatusInternalServerError)
					}
				}
			}(),

			expectations: func(t *testing.T, _ *metalv1.PortVlanAssignmentBatch, _ *http.Response, err error) {
				assert.ErrorContains(t, err, "failed to create batch")
			},
		},
		{
			name: "failed-to-poll-batch",
			handler: func() http.HandlerFunc {
				batchID := "batch-id"
				invocations := []http.HandlerFunc{
					func(resp http.ResponseWriter, req *http.Request) {
						b := metalv1.PortVlanAssignmentBatch{
							Id: &batchID,
						}

						if strings.Contains(req.URL.String(), "port_foo") {
							resp.Header().Add("content-type", "application/json")
							bytes, err := b.MarshalJSON()
							if err != nil {
								t.Fatalf("failed to marshal json: %v", err)
							}
							resp.Write(bytes)
						}
					},
					func(resp http.ResponseWriter, req *http.Request) {
						b := metalv1.PortVlanAssignmentBatch{
							Id: &batchID,
							State: metalv1.PORTVLANASSIGNMENTBATCHSTATE_QUEUED.Ptr(),
						}

						if strings.Contains(req.URL.String(), "port_foo") {
							resp.Header().Add("content-type", "application/json")
							bytes, err := b.MarshalJSON()
							if err != nil {
								t.Fatalf("failed to marshal json: %v", err)
							}
							resp.Write(bytes)
						}
					},
					func(resp http.ResponseWriter, req *http.Request) {
							resp.Header().Add("content-type", "application/json")
							resp.WriteHeader(http.StatusInternalServerError)
					},
						
				}


				index := 0
				return func(resp http.ResponseWriter, req *http.Request) {
					invocations[index](resp, req)
					index++
				}
			}(),

			expectations: func(t *testing.T, _ *metalv1.PortVlanAssignmentBatch, _ *http.Response, err error) {
				assert.ErrorContains(t, err, "could not be polled")
			},
		},

	}

	for _, tc := range cases {
		t.Run(tc.name, func(tt *testing.T) {
			mockAPI := httptest.NewServer(http.HandlerFunc(tc.handler))
			defer mockAPI.Close()
			meta := &config.Config{
				BaseURL: mockAPI.URL,
				Token:   "superNotRealToken",
			}

			err := meta.Load(ctx)

			client := meta.NewMetalClientForTesting()

			batchResp, httpResp, err := vb.Execute(ctx, client)

			tc.expectations(tt, batchResp, httpResp, err)
		})
	}

}
