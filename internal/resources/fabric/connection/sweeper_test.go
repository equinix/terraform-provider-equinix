package connection

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

// expectMethod fails the test if the request isn't the given method, then
// writes response.
func expectMethod(t *testing.T, method string, response testResponse) testResponse {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("got %s %s, want %s", r.Method, r.URL.Path, method)
		}
		response(w, r)
	}
}

func sweepTestConnection(state fabricv4.ConnectionState) fabricv4.Connection {
	conn := fabricv4.Connection{}
	conn.SetUuid(testConnectionUUID)
	conn.SetName("tfacc_sweep_PFCR")
	conn.SetState(state)
	return conn
}

func useFastSweepPolling(t *testing.T) {
	t.Helper()
	orig := connectionSweepPollInterval
	connectionSweepPollInterval = time.Millisecond
	t.Cleanup(func() { connectionSweepPollInterval = orig })
}

func TestSweepConnection_DeletesAndWaits(t *testing.T) {
	useFastSweepPolling(t)
	client, calls := newTestFabricClient(t,
		expectMethod(t, http.MethodDelete, withState(fabricv4.CONNECTIONSTATE_DEPROVISIONING)),
		expectMethod(t, http.MethodGet, withState(fabricv4.CONNECTIONSTATE_DEPROVISIONING)),
		expectMethod(t, http.MethodGet, withState(fabricv4.CONNECTIONSTATE_DEPROVISIONED)),
	)

	if err := sweepConnection(context.Background(), client, sweepTestConnection(fabricv4.CONNECTIONSTATE_ACTIVE)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := calls.Load(); got != 3 {
		t.Errorf("requests = %d, want 3", got)
	}
}

func TestSweepConnection_RetriesTransientState(t *testing.T) {
	useFastSweepPolling(t)
	transient := statusResponse(http.StatusBadRequest, fabricErrorBody("EQ-3142510", "Connection is in transient state"))
	client, calls := newTestFabricClient(t,
		expectMethod(t, http.MethodDelete, transient),
		expectMethod(t, http.MethodDelete, withState(fabricv4.CONNECTIONSTATE_DEPROVISIONING)),
		expectMethod(t, http.MethodGet, withState(fabricv4.CONNECTIONSTATE_DEPROVISIONED)),
	)

	if err := sweepConnection(context.Background(), client, sweepTestConnection(fabricv4.CONNECTIONSTATE_PROVISIONING)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := calls.Load(); got != 3 {
		t.Errorf("requests = %d, want 3", got)
	}
}

func TestSweepConnection_TransientStateUntilDeadline(t *testing.T) {
	useFastSweepPolling(t)
	client, _ := newTestFabricClient(t,
		expectMethod(t, http.MethodDelete,
			statusResponse(http.StatusBadRequest, fabricErrorBody("EQ-3142510", "Connection is in transient state"))),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	err := sweepConnection(ctx, client, sweepTestConnection(fabricv4.CONNECTIONSTATE_PROVISIONING))
	if err == nil {
		t.Fatal("expected an error when the connection stays in a transient state")
	}
	for _, want := range []string{"tfacc_sweep_PFCR", testConnectionUUID, "transient state"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
}

func TestSweepConnection_AlreadyDeletedSkipsWait(t *testing.T) {
	useFastSweepPolling(t)
	client, calls := newTestFabricClient(t,
		expectMethod(t, http.MethodDelete,
			statusResponse(http.StatusBadRequest, fabricErrorBody("EQ-3142509", "Connection already deleted"))),
	)

	if err := sweepConnection(context.Background(), client, sweepTestConnection(fabricv4.CONNECTIONSTATE_ACTIVE)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("requests = %d, want 1 (no GET after already-deleted)", got)
	}
}

func TestSweepConnection_DeprovisioningOnlyWaits(t *testing.T) {
	useFastSweepPolling(t)
	client, calls := newTestFabricClient(t,
		expectMethod(t, http.MethodGet,
			statusResponse(http.StatusNotFound, fabricErrorBody("EQ-3142102", "Connection not found"))),
	)

	if err := sweepConnection(context.Background(), client, sweepTestConnection(fabricv4.CONNECTIONSTATE_DEPROVISIONING)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("requests = %d, want 1 (no DELETE for DEPROVISIONING)", got)
	}
}

func TestSweepConnection_DeleteErrorIsNotRetried(t *testing.T) {
	useFastSweepPolling(t)
	client, calls := newTestFabricClient(t,
		expectMethod(t, http.MethodDelete,
			statusResponse(http.StatusBadRequest, fabricErrorBody("EQ-3142231", "some other validation error"))),
	)

	err := sweepConnection(context.Background(), client, sweepTestConnection(fabricv4.CONNECTIONSTATE_ACTIVE))
	if err == nil || !strings.Contains(err.Error(), "error deleting fabric connection tfacc_sweep_PFCR") {
		t.Fatalf("err = %v, want delete error naming the connection", err)
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("requests = %d, want 1", got)
	}
}

func TestSweepConnection_ForbiddenIsIgnored(t *testing.T) {
	useFastSweepPolling(t)
	client, _ := newTestFabricClient(t,
		expectMethod(t, http.MethodDelete, statusResponse(http.StatusForbidden, fabricErrorBody("EQ-3155114", "Authorization error"))),
		expectMethod(t, http.MethodGet, withState(fabricv4.CONNECTIONSTATE_DEPROVISIONED)),
	)

	if err := sweepConnection(context.Background(), client, sweepTestConnection(fabricv4.CONNECTIONSTATE_ACTIVE)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemainingSweepTime(t *testing.T) {
	if got := remainingSweepTime(context.Background()); got != connectionSweepTimeout {
		t.Errorf("no deadline: got %v, want %v", got, connectionSweepTimeout)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if got := remainingSweepTime(ctx); got <= 8*time.Second || got > 9*time.Second {
		t.Errorf("10s deadline: got %v, want just under 9s", got)
	}
}
