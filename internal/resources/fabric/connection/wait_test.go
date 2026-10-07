package connection

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

const testConnectionUUID = "11111111-2222-3333-4444-555555555555"

// newTestFabricClient returns a fabricv4 client pointed at a server that
// answers GET /fabric/v4/connections/{uuid} with responses[n] on the nth
// call, repeating the last response once they run out.
func newTestFabricClient(t *testing.T, responses ...func(w http.ResponseWriter)) (*fabricv4.APIClient, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := int(calls.Add(1)) - 1
		if n >= len(responses) {
			n = len(responses) - 1
		}
		w.Header().Set("Content-Type", "application/json")
		responses[n](w)
	}))
	t.Cleanup(server.Close)

	cfg := fabricv4.NewConfiguration()
	cfg.Servers = fabricv4.ServerConfigurations{{URL: server.URL}}
	cfg.HTTPClient = server.Client()
	return fabricv4.NewAPIClient(cfg), &calls
}

// connectionJSON renders a minimal connection that satisfies the SDK's
// required properties, plus the given extra fields.
func connectionJSON(fields string) func(w http.ResponseWriter) {
	return func(w http.ResponseWriter) {
		fmt.Fprintf(w, `{"uuid":%q,"type":"EVPL_VC","name":"tfacc_wait","bandwidth":1,"aSide":{},"zSide":{},%s}`, testConnectionUUID, fields)
	}
}

func withState(state fabricv4.ConnectionState) func(w http.ResponseWriter) {
	return connectionJSON(fmt.Sprintf(`"state":%q`, state))
}

func statusResponse(code int, body string) func(w http.ResponseWriter) {
	return func(w http.ResponseWriter) {
		w.WriteHeader(code)
		fmt.Fprint(w, body)
	}
}

func fabricErrorBody(code, message string) string {
	return fmt.Sprintf(`[{"errorCode":%q,"errorMessage":%q}]`, code, message)
}

// fastPoll keeps the StateChangeConf delays out of the test runtime.
var fastPoll = withPollInterval(time.Millisecond)

func TestWaitForConnection_ReachesTargetState(t *testing.T) {
	client, calls := newTestFabricClient(t,
		withState(fabricv4.CONNECTIONSTATE_PROVISIONING),
		withState(fabricv4.CONNECTIONSTATE_PROVISIONING),
		withState(fabricv4.CONNECTIONSTATE_PROVISIONED),
	)

	conn, err := waitForConnection(context.Background(), client, testConnectionUUID,
		[]string{string(fabricv4.CONNECTIONSTATE_PROVISIONING)},
		[]string{string(fabricv4.CONNECTIONSTATE_PROVISIONED)},
		time.Second, fastPoll,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := conn.GetState(); got != fabricv4.CONNECTIONSTATE_PROVISIONED {
		t.Errorf("state = %q, want %q", got, fabricv4.CONNECTIONSTATE_PROVISIONED)
	}
	if got := calls.Load(); got != 3 {
		t.Errorf("GET calls = %d, want 3", got)
	}
}

func TestWaitForConnection_UnexpectedStateFails(t *testing.T) {
	client, _ := newTestFabricClient(t, withState(fabricv4.CONNECTIONSTATE_FAILED))

	_, err := waitForConnection(context.Background(), client, testConnectionUUID,
		[]string{string(fabricv4.CONNECTIONSTATE_PROVISIONING)},
		[]string{string(fabricv4.CONNECTIONSTATE_PROVISIONED)},
		time.Second, fastPoll,
	)
	if err == nil || !strings.Contains(err.Error(), "unexpected state 'FAILED'") {
		t.Fatalf("err = %v, want unexpected state 'FAILED'", err)
	}
}

func TestWaitForConnection_EmptyPendingWaitsThroughAnyState(t *testing.T) {
	client, _ := newTestFabricClient(t,
		withState(fabricv4.CONNECTIONSTATE_DRAFT),
		withState(fabricv4.CONNECTIONSTATE_PENDING),
	)

	_, err := waitForConnection(context.Background(), client, testConnectionUUID,
		nil,
		[]string{string(fabricv4.CONNECTIONSTATE_PENDING)},
		time.Second, fastPoll,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWaitForConnection_Timeout(t *testing.T) {
	client, _ := newTestFabricClient(t, withState(fabricv4.CONNECTIONSTATE_DEPROVISIONING))

	_, err := waitForConnection(context.Background(), client, testConnectionUUID,
		[]string{string(fabricv4.CONNECTIONSTATE_DEPROVISIONING)},
		[]string{string(fabricv4.CONNECTIONSTATE_DEPROVISIONED)},
		50*time.Millisecond, fastPoll,
	)
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("err = %v, want timeout", err)
	}
}

func TestWaitForConnection_NotFound(t *testing.T) {
	notFound := statusResponse(http.StatusNotFound, fabricErrorBody("EQ-3142102", "Connection not found"))
	pending := []string{string(fabricv4.CONNECTIONSTATE_DEPROVISIONING)}
	target := []string{string(fabricv4.CONNECTIONSTATE_DEPROVISIONED)}

	t.Run("fails by default", func(t *testing.T) {
		client, _ := newTestFabricClient(t, notFound)
		_, err := waitForConnection(context.Background(), client, testConnectionUUID,
			pending, target, time.Second, fastPoll,
		)
		if err == nil {
			t.Fatal("expected an error for 404 without withNotFoundAs")
		}
	})

	t.Run("withNotFoundAs treats 404 as target", func(t *testing.T) {
		client, calls := newTestFabricClient(t,
			withState(fabricv4.CONNECTIONSTATE_DEPROVISIONING),
			notFound,
		)
		conn, err := waitForConnection(context.Background(), client, testConnectionUUID,
			pending, target, time.Second, fastPoll,
			withNotFoundAs(string(fabricv4.CONNECTIONSTATE_DEPROVISIONED)),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if conn == nil {
			t.Error("expected a non-nil connection")
		}
		if got := calls.Load(); got != 2 {
			t.Errorf("GET calls = %d, want 2", got)
		}
	})
}

func TestWaitForConnection_ServerErrorFails(t *testing.T) {
	client, _ := newTestFabricClient(t,
		statusResponse(http.StatusInternalServerError, fabricErrorBody("EQ-3000000", "boom")),
	)

	_, err := waitForConnection(context.Background(), client, testConnectionUUID,
		[]string{string(fabricv4.CONNECTIONSTATE_PROVISIONING)},
		[]string{string(fabricv4.CONNECTIONSTATE_PROVISIONED)},
		time.Second, fastPoll,
		withNotFoundAs(string(fabricv4.CONNECTIONSTATE_PROVISIONED)),
	)
	if err == nil {
		t.Fatal("expected an error for a 500 response")
	}
}

func TestWaitForConnection_WithConnectionStatus(t *testing.T) {
	client, _ := newTestFabricClient(t,
		connectionJSON(`"state":"PROVISIONED","operation":{"providerStatus":"PROVISIONING"}`),
		connectionJSON(`"state":"PROVISIONED","operation":{"providerStatus":"PROVISIONED"}`),
	)

	conn, err := waitForConnection(context.Background(), client, testConnectionUUID,
		[]string{string(fabricv4.PROVIDERSTATUS_PROVISIONING)},
		[]string{string(fabricv4.PROVIDERSTATUS_PROVISIONED)},
		time.Second, fastPoll,
		withConnectionStatus(func(c *fabricv4.Connection) string {
			operation := c.GetOperation()
			return string(operation.GetProviderStatus())
		}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	operation := conn.GetOperation()
	if got := operation.GetProviderStatus(); got != fabricv4.PROVIDERSTATUS_PROVISIONED {
		t.Errorf("providerStatus = %q, want %q", got, fabricv4.PROVIDERSTATUS_PROVISIONED)
	}
}

func TestIsConnectionAlreadyDeleted(t *testing.T) {
	tests := []struct {
		name     string
		response func(w http.ResponseWriter)
		want     bool
	}{
		{
			name:     "already deleted",
			response: statusResponse(http.StatusBadRequest, fabricErrorBody("EQ-3142509", "Connection already deleted")),
			want:     true,
		},
		{
			name:     "other bad request",
			response: statusResponse(http.StatusBadRequest, fabricErrorBody("EQ-3142231", "Seller port doesn't have bandwidth available")),
			want:     false,
		},
		{
			name:     "not found",
			response: statusResponse(http.StatusNotFound, fabricErrorBody("EQ-3142102", "Connection not found")),
			want:     false,
		},
		{
			name:     "success",
			response: withState(fabricv4.CONNECTIONSTATE_DEPROVISIONING),
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newTestFabricClient(t, tt.response)
			_, _, err := client.ConnectionsApi.DeleteConnectionByUuid(context.Background(), testConnectionUUID).Execute()
			if got := isConnectionAlreadyDeleted(err); got != tt.want {
				t.Errorf("isConnectionAlreadyDeleted(%v) = %v, want %v", err, got, tt.want)
			}
		})
	}

	if isConnectionAlreadyDeleted(fmt.Errorf("plain error")) {
		t.Error("isConnectionAlreadyDeleted(plain error) = true, want false")
	}
}
