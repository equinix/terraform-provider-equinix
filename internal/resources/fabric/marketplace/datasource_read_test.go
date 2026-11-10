package marketplace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
)

const testSubscriptionID = "00000000-0000-0000-0000-000000000000"

func testFabricConfig(t *testing.T, status int) *config.Config {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`[{"errorCode":"EQ-3041501","errorMessage":"Marketplace subscription does not exist"}]`))
	}))
	t.Cleanup(srv.Close)
	return &config.Config{BaseURL: srv.URL, Token: "test-token"}
}

func TestDataSourceRead_notFoundIsAnError(t *testing.T) {
	d := DataSourceFabricMarketplaceSubscription().TestResourceData()
	if err := d.Set("uuid", testSubscriptionID); err != nil {
		t.Fatal(err)
	}

	diags := dataSourceFabricMarketplaceSubscriptionRead(context.Background(), d, testFabricConfig(t, http.StatusNotFound))

	if !diags.HasError() {
		t.Fatal("expected the data source to fail when the subscription does not exist")
	}
}

func TestDataSourceRead_otherErrorsAreReturned(t *testing.T) {
	d := DataSourceFabricMarketplaceSubscription().TestResourceData()
	if err := d.Set("uuid", testSubscriptionID); err != nil {
		t.Fatal(err)
	}

	diags := dataSourceFabricMarketplaceSubscriptionRead(context.Background(), d, testFabricConfig(t, http.StatusForbidden))

	if !diags.HasError() {
		t.Fatal("expected an error on 403")
	}
}
