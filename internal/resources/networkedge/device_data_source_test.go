package networkedge

import (
	"context"
	"net/http"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/networkedge"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const deviceDataSourceName = "equinix_network_device"

func TestDeviceDataSource_validStatusList(t *testing.T) {
	ds, ok := equinix.Provider().DataSourcesMap[deviceDataSourceName]
	require.True(t, ok, "provider has data source %s", deviceDataSourceName)

	tests := []struct {
		name       string
		statusList string
		// wantStatus is the status filter that the provider sends when it lists devices.
		wantStatus string
		wantError  bool
	}{
		{name: "empty list", statusList: "", wantStatus: ""},
		{name: "one item", statusList: "provisioned", wantStatus: "PROVISIONED"},
		{name: "two items", statusList: "provisioning,provisioned", wantStatus: "PROVISIONING,PROVISIONED"},
		{name: "whitespace", statusList: "provisioning, \tprovisioned  ", wantStatus: "PROVISIONING,PROVISIONED"},
		{name: "capitalization", statusList: "provisioning, \tproVISioned  ", wantStatus: "PROVISIONING,PROVISIONED"},
		{name: "invalid item", statusList: "provisioning, provisioned,somethingInvalid ", wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			api := networkedge.NewMockAPI(t)
			api.Devices.Seed("6e0c2b09-37e9-4fd2-b9e4-e68b08d8b29d", map[string]any{
				"name":   "testDevice",
				"status": "PROVISIONED",
			})
			d := schema.TestResourceDataRaw(t, ds.Schema, map[string]any{
				"name":              "testDevice",
				"valid_status_list": tt.statusList,
			})

			// when
			diags := ds.ReadContext(context.Background(), d, api.Config)

			// then
			lists := api.RequestsFor(http.MethodGet, networkedge.DevicesPath)
			if tt.wantError {
				assert.True(t, diags.HasError(), "read returns an error")
				assert.Empty(t, lists, "no device list request was sent")
				return
			}
			require.False(t, diags.HasError(), "read does not return errors: %v", diags)
			require.Len(t, lists, 1, "one device list request was sent")
			assert.Equal(t, tt.wantStatus, lists[0].Query.Get("status"), "status filter matches")
		})
	}
}
