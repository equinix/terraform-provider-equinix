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

const deviceLinkResourceName = "equinix_network_device_link"

const (
	deviceLinkDeviceA = "3eee8518-b19d-4de5-afd8-afd9b67e6e8c"
	deviceLinkDeviceB = "7c737fe8-3a9e-4ab9-afcc-c06d01db326d"
)

func deviceLinkResource(t *testing.T) *schema.Resource {
	t.Helper()
	r, ok := equinix.Provider().ResourcesMap[deviceLinkResourceName]
	require.True(t, ok, "provider has resource %s", deviceLinkResourceName)
	return r
}

func TestDeviceLink_create(t *testing.T) {
	// given
	api := networkedge.NewMockAPI(t)
	// Reading a link looks up the ASN of each linked device.
	api.Devices.Seed(deviceLinkDeviceA, map[string]any{"asn": 0})
	api.Devices.Seed(deviceLinkDeviceB, map[string]any{"asn": 0})
	d := schema.TestResourceDataRaw(t, deviceLinkResource(t).Schema, map[string]any{
		"name":            "testGroup",
		"subnet":          "10.10.1.0/24",
		"project_id":      "68ccfd49-39b1-478e-957a-67c72f719d7a",
		"redundancy_type": "PRIMARY",
		"device": []any{
			map[string]any{"id": deviceLinkDeviceA, "asn": 0, "interface_id": 5},
			map[string]any{"id": deviceLinkDeviceB, "asn": 0, "interface_id": 10},
		},
		"metro_link": []any{
			map[string]any{"account_number": "", "metro_code": "MX", "throughput": "10", "throughput_unit": "Mbps"},
			map[string]any{"account_number": "", "metro_code": "SV", "throughput": "100", "throughput_unit": "Mbps"},
		},
	})

	// when
	diags := deviceLinkResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)

	posts := api.RequestsFor(http.MethodPost, networkedge.LinksPath)
	require.Len(t, posts, 1, "one create request was sent")
	body := posts[0].Body
	assert.Equal(t, "testGroup", body["groupName"], "group name matches")
	assert.Equal(t, "10.10.1.0/24", body["subnet"], "subnet matches")
	assert.Equal(t, "68ccfd49-39b1-478e-957a-67c72f719d7a", body["projectId"], "project ID matches")
	assert.Equal(t, "PRIMARY", body["redundancyType"], "redundancy type matches")
	assert.ElementsMatch(t, []any{
		map[string]any{"deviceUuid": deviceLinkDeviceA, "asn": float64(0), "interfaceId": float64(5)},
		map[string]any{"deviceUuid": deviceLinkDeviceB, "asn": float64(0), "interfaceId": float64(10)},
	}, body["linkDevices"], "devices match")
	assert.ElementsMatch(t, []any{
		map[string]any{"accountNumber": "", "metroCode": "MX", "throughput": "10", "throughputUnit": "Mbps"},
		map[string]any{"accountNumber": "", "metroCode": "SV", "throughput": "100", "throughputUnit": "Mbps"},
	}, body["metroLinks"], "metro links match")
}

func TestDeviceLink_read(t *testing.T) {
	// given
	const id = "aae04283-10f9-4edb-9395-33681176592b"
	api := networkedge.NewMockAPI(t)
	api.Devices.Seed(deviceLinkDeviceA, map[string]any{"asn": 0})
	api.Devices.Seed(deviceLinkDeviceB, map[string]any{"asn": 0})
	api.Links.Seed(id, map[string]any{
		"groupName":      "testGroup",
		"subnet":         "10.10.1.0/24",
		"redundancyType": "PRIMARY",
		"status":         "PROVISIONED",
		"linkDevices": []any{
			map[string]any{"deviceUuid": deviceLinkDeviceA, "asn": 0, "interfaceId": 5},
			map[string]any{"deviceUuid": deviceLinkDeviceB, "asn": 0, "interfaceId": 10},
		},
		"metroLinks": []any{
			map[string]any{"accountNumber": "", "metroCode": "MX", "throughput": "10", "throughputUnit": "Mbps"},
			map[string]any{"accountNumber": "", "metroCode": "LD", "throughput": "10", "throughputUnit": "Mbps"},
		},
	})
	d := schema.TestResourceDataRaw(t, deviceLinkResource(t).Schema, map[string]any{})
	d.SetId(id)

	// when
	diags := deviceLinkResource(t).ReadContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "read does not return errors: %v", diags)
	assert.Equal(t, id, d.Get("uuid"), "UUID matches")
	assert.Equal(t, "testGroup", d.Get("name"), "name matches")
	assert.Equal(t, "10.10.1.0/24", d.Get("subnet"), "subnet matches")
	assert.Equal(t, "PROVISIONED", d.Get("status"), "status matches")

	devices := d.Get("device").(*schema.Set).List()
	require.Len(t, devices, 2, "two devices")
	deviceIDs := map[string]int{}
	for _, v := range devices {
		m := v.(map[string]any)
		deviceIDs[m["id"].(string)] = m["interface_id"].(int)
		assert.Equal(t, 0, m["asn"], "device ASN matches")
	}
	assert.Equal(t, map[string]int{deviceLinkDeviceA: 5, deviceLinkDeviceB: 10}, deviceIDs, "devices match")

	var metroCodes []string
	for _, v := range d.Get("metro_link").(*schema.Set).List() {
		m := v.(map[string]any)
		metroCodes = append(metroCodes, m["metro_code"].(string))
		assert.Equal(t, "10", m["throughput"], "throughput matches")
		assert.Equal(t, "Mbps", m["throughput_unit"], "throughput unit matches")
	}
	assert.ElementsMatch(t, []string{"MX", "LD"}, metroCodes, "metro links match")
}
