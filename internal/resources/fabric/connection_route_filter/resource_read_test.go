package connection_route_filter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
)

const (
	testConnectionID    = "11111111-0000-0000-0000-000000000000"
	testRouteFilterUUID = "00000000-0000-0000-0000-000000000000"
)

func testFabricConfig(t *testing.T, status int) *config.Config {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`[{"errorCode":"EQ-3045601","errorMessage":"Connection route filter does not exist"}]`))
	}))
	t.Cleanup(srv.Close)
	return &config.Config{BaseURL: srv.URL, Token: "test-token"}
}

func TestResourceRead_notFoundRemovesFromState(t *testing.T) {
	d := Resource().TestResourceData()
	d.SetId(testRouteFilterUUID)
	if err := d.Set("connection_id", testConnectionID); err != nil {
		t.Fatal(err)
	}

	diags := resourceRead(context.Background(), d, testFabricConfig(t, http.StatusNotFound))

	if diags.HasError() {
		t.Fatalf("expected no error on 404, got %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("expected resource to be removed from state, got id %q", d.Id())
	}
}

func TestResourceRead_otherErrorsAreReturned(t *testing.T) {
	d := Resource().TestResourceData()
	d.SetId(testRouteFilterUUID)
	if err := d.Set("connection_id", testConnectionID); err != nil {
		t.Fatal(err)
	}

	diags := resourceRead(context.Background(), d, testFabricConfig(t, http.StatusForbidden))

	if !diags.HasError() {
		t.Fatal("expected an error on 403")
	}
}
