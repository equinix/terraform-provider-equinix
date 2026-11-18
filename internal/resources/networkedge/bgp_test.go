package networkedge

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/networkedge"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const bgpResourceName = "equinix_network_bgp"

func bgpResource(t *testing.T) *schema.Resource {
	t.Helper()
	r, ok := equinix.Provider().ResourcesMap[bgpResourceName]
	require.True(t, ok, "provider has resource %s", bgpResourceName)
	return r
}

func TestBGP_create(t *testing.T) {
	// given
	api := networkedge.NewMockAPI(t)
	d := schema.TestResourceDataRaw(t, bgpResource(t).Schema, map[string]any{
		"connection_id":      "6ca8d0df-c71a-4475-a835-53c2df1e6667",
		"local_ip_address":   "1.1.1.1/32",
		"local_asn":          15344,
		"remote_ip_address":  "2.2.2.2",
		"remote_asn":         60421,
		"authentication_key": "secret",
	})

	// when
	diags := bgpResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)

	posts := api.RequestsFor(http.MethodPost, networkedge.BGPPath)
	require.Len(t, posts, 1, "one create request was sent")
	assert.Equal(t, map[string]any{
		"connectionUuid":    "6ca8d0df-c71a-4475-a835-53c2df1e6667",
		"localIpAddress":    "1.1.1.1/32",
		"localAsn":          float64(15344),
		"remoteIpAddress":   "2.2.2.2",
		"remoteAsn":         float64(60421),
		"authenticationKey": "secret",
	}, posts[0].Body, "create request body matches the API contract")
	assert.NotEmpty(t, d.Id(), "resource ID is set")
	assert.NotEmpty(t, api.RequestsFor(http.MethodGet, networkedge.BGPPath+"/"+d.Id()),
		"provisioning status of the new configuration was checked")
	assert.Equal(t, "PROVISIONED", d.Get("provisioning_status"), "provisioning status matches")
}

func TestBGP_read(t *testing.T) {
	// given
	const id = "0cb9759d-58ab-44e6-9c10-6a3cfd18cefb"
	api := networkedge.NewMockAPI(t)
	api.BGP.Seed(id, map[string]any{
		"virtualDeviceUuid":  "8895983f-00f9-42f1-a387-85248f2aab49",
		"connectionUuid":     "6ca8d0df-c71a-4475-a835-53c2df1e6667",
		"localIpAddress":     "1.1.1.1/32",
		"localAsn":           15344,
		"remoteIpAddress":    "2.2.2.2",
		"remoteAsn":          60421,
		"authenticationKey":  "secret",
		"state":              "established",
		"provisioningStatus": "PROVISIONED",
	})
	d := schema.TestResourceDataRaw(t, bgpResource(t).Schema, map[string]any{})
	d.SetId(id)

	// when
	diags := bgpResource(t).ReadContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "read does not return errors: %v", diags)
	assert.Equal(t, id, d.Get("uuid"), "UUID matches")
	assert.Equal(t, "8895983f-00f9-42f1-a387-85248f2aab49", d.Get("device_id"), "device ID matches")
	assert.Equal(t, "6ca8d0df-c71a-4475-a835-53c2df1e6667", d.Get("connection_id"), "connection ID matches")
	assert.Equal(t, "1.1.1.1/32", d.Get("local_ip_address"), "local IP address matches")
	assert.Equal(t, 15344, d.Get("local_asn"), "local ASN matches")
	assert.Equal(t, "2.2.2.2", d.Get("remote_ip_address"), "remote IP address matches")
	assert.Equal(t, 60421, d.Get("remote_asn"), "remote ASN matches")
	assert.Equal(t, "secret", d.Get("authentication_key"), "authentication key matches")
	assert.Equal(t, "established", d.Get("state"), "state matches")
	assert.Equal(t, "PROVISIONED", d.Get("provisioning_status"), "provisioning status matches")
}

func TestBGP_update(t *testing.T) {
	// given
	const id = "0cb9759d-58ab-44e6-9c10-6a3cfd18cefb"
	api := networkedge.NewMockAPI(t)
	api.BGP.Seed(id, map[string]any{"provisioningStatus": "PROVISIONED"})
	d := schema.TestResourceDataRaw(t, bgpResource(t).Schema, map[string]any{
		"local_ip_address":   "1.1.1.1/32",
		"local_asn":          15344,
		"remote_ip_address":  "2.2.2.2",
		"remote_asn":         60421,
		"authentication_key": "secret",
	})
	d.SetId(id)
	require.NoError(t, d.Set("uuid", id))

	// when
	diags := bgpResource(t).UpdateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "update does not return errors: %v", diags)
	puts := api.RequestsFor(http.MethodPut, networkedge.BGPPath+"/"+id)
	require.Len(t, puts, 1, "one update request was sent")
	assert.Equal(t, map[string]any{
		"localIpAddress":    "1.1.1.1/32",
		"localAsn":          float64(15344),
		"remoteIpAddress":   "2.2.2.2",
		"remoteAsn":         float64(60421),
		"authenticationKey": "secret",
	}, puts[0].Body, "update request body matches the API contract")
}

func newBGPCreateData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, bgpResource(t).Schema, map[string]any{
		"connection_id":     "6ca8d0df-c71a-4475-a835-53c2df1e6667",
		"local_ip_address":  "1.1.1.1/32",
		"local_asn":         15344,
		"remote_ip_address": "2.2.2.2",
		"remote_asn":        60421,
	})
}

func TestBGP_createWaitsUntilProvisioned(t *testing.T) {
	t.Parallel()
	// given
	api := networkedge.NewMockAPI(t)
	api.BGP.Sequence("", "provisioningStatus", "PROVISIONING", "PROVISIONED")
	d := newBGPCreateData(t)

	// when
	diags := bgpResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.False(t, diags.HasError(), "create does not return errors: %v", diags)
	assert.Equal(t, "PROVISIONED", d.Get("provisioning_status"), "provisioning status matches")
	// Two status checks while waiting, and one read.
	assert.GreaterOrEqual(t, len(api.RequestsFor(http.MethodGet, networkedge.BGPPath+"/"+d.Id())), 3,
		"create kept checking the status until it was provisioned")
}

func TestBGP_createFailsOnUnexpectedStatus(t *testing.T) {
	t.Parallel()
	// given
	api := networkedge.NewMockAPI(t)
	api.BGP.Sequence("", "provisioningStatus", "FAILED")
	d := newBGPCreateData(t)

	// when
	diags := bgpResource(t).CreateContext(context.Background(), d, api.Config)

	// then
	require.True(t, diags.HasError(), "create returns an error")
	assert.Contains(t, diags[0].Summary, "error waiting for BGP configuration", "error comes from the wait")
}

func TestBGP_createTimesOutWhileProvisioning(t *testing.T) {
	t.Parallel()
	// given
	api := networkedge.NewMockAPI(t)
	api.BGP.Sequence("", "provisioningStatus", "PROVISIONING") // the status never changes
	r := bgpResource(t)
	timeout := time.Second
	r.Timeouts.Create = &timeout
	// Resource data from TestResourceDataRaw ignores the timeouts of the resource.
	d := r.Data(nil)
	for k, v := range map[string]any{
		"connection_id":     "6ca8d0df-c71a-4475-a835-53c2df1e6667",
		"local_ip_address":  "1.1.1.1/32",
		"local_asn":         15344,
		"remote_ip_address": "2.2.2.2",
		"remote_asn":        60421,
	} {
		require.NoError(t, d.Set(k, v))
	}

	// when
	diags := r.CreateContext(context.Background(), d, api.Config)

	// then
	require.True(t, diags.HasError(), "create returns an error")
	assert.Contains(t, diags[0].Summary, "error waiting for BGP configuration", "error comes from the wait")
	assert.Contains(t, diags[0].Summary, "timeout while waiting for state", "wait gave up because of the timeout")
}
