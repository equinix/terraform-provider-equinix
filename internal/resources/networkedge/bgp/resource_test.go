package bgp

import (
	"context"
	"testing"
	"time"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNetworkBGP_createFromResourceData(t *testing.T) {
	// given
	expected := ne.BGPConfiguration{
		ConnectionUUID:    ne.String("6ca8d0df-c71a-4475-a835-53c2df1e6667"),
		LocalIPAddress:    ne.String("1.1.1.1/32"),
		LocalASN:          ne.Int(15344),
		RemoteIPAddress:   ne.String("2.2.2.2"),
		RemoteASN:         ne.Int(60421),
		AuthenticationKey: ne.String("secret"),
	}
	rawData := map[string]interface{}{
		networkBGPSchemaNames["ConnectionUUID"]:    ne.StringValue(expected.ConnectionUUID),
		networkBGPSchemaNames["LocalIPAddress"]:    ne.StringValue(expected.LocalIPAddress),
		networkBGPSchemaNames["LocalASN"]:          ne.IntValue(expected.LocalASN),
		networkBGPSchemaNames["RemoteIPAddress"]:   ne.StringValue(expected.RemoteIPAddress),
		networkBGPSchemaNames["RemoteASN"]:         ne.IntValue(expected.RemoteASN),
		networkBGPSchemaNames["AuthenticationKey"]: ne.StringValue(expected.AuthenticationKey),
	}
	d := schema.TestResourceDataRaw(t, createNetworkBGPResourceSchema(), rawData)
	// when
	result := createNetworkBGPConfiguration(d)
	// then
	assert.Equal(t, expected, result, "Created BGP configuration matches expected result")
}

func TestNetworkBGP_updateResourceData(t *testing.T) {
	// when
	input := ne.BGPConfiguration{
		UUID:               ne.String("0cb9759d-58ab-44e6-9c10-6a3cfd18cefb"),
		DeviceUUID:         ne.String("8895983f-00f9-42f1-a387-85248f2aab49"),
		ConnectionUUID:     ne.String("6ca8d0df-c71a-4475-a835-53c2df1e6667"),
		LocalIPAddress:     ne.String("1.1.1.1/32"),
		LocalASN:           ne.Int(15344),
		RemoteIPAddress:    ne.String("2.2.2.2"),
		RemoteASN:          ne.Int(60421),
		AuthenticationKey:  ne.String("secret"),
		State:              ne.String("established"),
		ProvisioningStatus: ne.String(ne.BGPProvisioningStatusProvisioned),
	}
	d := schema.TestResourceDataRaw(t, createNetworkBGPResourceSchema(), make(map[string]interface{}))
	// when
	err := updateNetworkBGPResource(&input, d)
	// then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ne.StringValue(input.UUID), d.Get(networkBGPSchemaNames["UUID"]), "UUID matches")
	assert.Equal(t, ne.StringValue(input.DeviceUUID), d.Get(networkBGPSchemaNames["DeviceUUID"]), "DeviceUUID matches")
	assert.Equal(t, ne.StringValue(input.ConnectionUUID), d.Get(networkBGPSchemaNames["ConnectionUUID"]), "ConnectionUUID matches")
	assert.Equal(t, ne.StringValue(input.LocalIPAddress), d.Get(networkBGPSchemaNames["LocalIPAddress"]), "LocalIPAddress matches")
	assert.Equal(t, ne.IntValue(input.LocalASN), d.Get(networkBGPSchemaNames["LocalASN"]), "LocalASN matches")
	assert.Equal(t, ne.StringValue(input.RemoteIPAddress), d.Get(networkBGPSchemaNames["RemoteIPAddress"]), "RemoteIPAddress matches")
	assert.Equal(t, ne.IntValue(input.RemoteASN), d.Get(networkBGPSchemaNames["RemoteASN"]), "RemoteASN matches")
	assert.Equal(t, ne.StringValue(input.AuthenticationKey), d.Get(networkBGPSchemaNames["AuthenticationKey"]), "AuthenticationKey matches")
	assert.Equal(t, ne.StringValue(input.State), d.Get(networkBGPSchemaNames["State"]), "State matches")
	assert.Equal(t, ne.StringValue(input.ProvisioningStatus), d.Get(networkBGPSchemaNames["ProvisioningStatus"]), "ProvisioningStatus matches")
}

type mockedBGPUpdateRequest struct {
	uuid string
	data map[string]interface{}
}

func (r *mockedBGPUpdateRequest) WithLocalIPAddress(v string) ne.BGPUpdateRequest {
	r.data["localIPAddress"] = v
	return r
}

func (r *mockedBGPUpdateRequest) WithLocalASN(v int) ne.BGPUpdateRequest {
	r.data["localASN"] = v
	return r
}

func (r *mockedBGPUpdateRequest) WithRemoteASN(v int) ne.BGPUpdateRequest {
	r.data["remoteASN"] = v
	return r
}

func (r *mockedBGPUpdateRequest) WithRemoteIPAddress(v string) ne.BGPUpdateRequest {
	r.data["remoteIPAddress"] = v
	return r
}

func (r *mockedBGPUpdateRequest) WithAuthenticationKey(v string) ne.BGPUpdateRequest {
	r.data["authenticationKey"] = v
	return r
}

func (r *mockedBGPUpdateRequest) Execute() error {
	return nil
}

func TestNetworkBGP_createUpdateRequest(t *testing.T) {
	// given
	req := &mockedBGPUpdateRequest{data: make(map[string]interface{})}
	f := func(uuid string) ne.BGPUpdateRequest {
		req.uuid = uuid
		return req
	}
	bgp := ne.BGPConfiguration{
		LocalIPAddress:    ne.String("1.1.1.1/32"),
		LocalASN:          ne.Int(15344),
		RemoteIPAddress:   ne.String("2.2.2.2"),
		RemoteASN:         ne.Int(60421),
		AuthenticationKey: ne.String("secret"),
	}
	// when
	createNetworkBGPUpdateRequest(f, &bgp)
	// then
	assert.Equal(t, ne.StringValue(bgp.RemoteIPAddress), req.data["remoteIPAddress"], "RemoteIPAddress matches")
	assert.Equal(t, ne.IntValue(bgp.RemoteASN), req.data["remoteASN"], "RemoteASN matches")
	assert.Equal(t, ne.StringValue(bgp.LocalIPAddress), req.data["localIPAddress"], "LocalIPAddress matches")
	assert.Equal(t, ne.IntValue(bgp.LocalASN), req.data["localASN"], "LocalASN matches")
	assert.Equal(t, ne.StringValue(bgp.AuthenticationKey), req.data["authenticationKey"], "AuthenticationKey matches")
}

func TestNetworkBGP_statusProvisioningWaitConfiguration(t *testing.T) {
	// given
	bgpID := "test"
	var queriedID string
	fetchFunc := func(uuid string) (*ne.BGPConfiguration, error) {
		queriedID = uuid
		return &ne.BGPConfiguration{ProvisioningStatus: ne.String(ne.BGPProvisioningStatusProvisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	// when
	waitConfig := createBGPConfigStatusProvisioningWaitConfiguration(fetchFunc, bgpID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	// then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, bgpID, queriedID, "Queried device ID matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}
