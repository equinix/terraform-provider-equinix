package connection

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
)

const testConnectionID = "00000000-0000-0000-0000-000000000000"

func testFabricConfig(t *testing.T, status int) *config.Config {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`[{"errorCode":"EQ-3142102","errorMessage":"Connection does not exist"}]`))
	}))
	t.Cleanup(srv.Close)
	return &config.Config{BaseURL: srv.URL, Token: "test-token"}
}

func TestResourceFabricConnectionRead_notFoundRemovesFromState(t *testing.T) {
	d := Resource().TestResourceData()
	d.SetId(testConnectionID)

	diags := resourceFabricConnectionRead(context.Background(), d, testFabricConfig(t, http.StatusNotFound))

	if diags.HasError() {
		t.Fatalf("expected no error on 404, got %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("expected connection to be removed from state, got id %q", d.Id())
	}
}

func TestResourceFabricConnectionRead_otherErrorsAreReturned(t *testing.T) {
	d := Resource().TestResourceData()
	d.SetId(testConnectionID)

	diags := resourceFabricConnectionRead(context.Background(), d, testFabricConfig(t, http.StatusForbidden))

	if !diags.HasError() {
		t.Fatal("expected an error on 403")
	}
}

func TestDataSourceFabricConnectionRead_notFoundIsAnError(t *testing.T) {
	d := DataSource().TestResourceData()
	if err := d.Set("uuid", testConnectionID); err != nil {
		t.Fatal(err)
	}

	diags := dataSourceFabricConnectionRead(context.Background(), d, testFabricConfig(t, http.StatusNotFound))

	if !diags.HasError() {
		t.Fatal("expected the data source to fail when the connection does not exist")
	}
}
