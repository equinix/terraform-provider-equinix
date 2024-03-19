package device_link

import (
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNetworkDeviceLink_createFromResourceData(t *testing.T) {
	// given

	expected := ne.DeviceLinkGroup{
		Name:      ne.String("testGroup"),
		Subnet:    ne.String("10.10.1.0/24"),
		ProjectID: ne.String("68ccfd49-39b1-478e-957a-67c72f719d7a"),
		Devices: []ne.DeviceLinkGroupDevice{
			{
				DeviceID:    ne.String("3eee8518-b19d-4de5-afd8-afd9b67e6e8c"),
				ASN:         ne.Int(0),
				InterfaceID: ne.Int(5),
			}, {
				DeviceID:    ne.String("7c737fe8-3a9e-4ab9-afcc-c06d01db326d"),
				ASN:         ne.Int(0),
				InterfaceID: ne.Int(10),
			},
		},
		Links: []ne.DeviceLinkGroupLink{
			{
				AccountNumber:        ne.String(""),
				Throughput:           ne.String("1"),
				ThroughputUnit:       ne.String("Gbps"),
				SourceMetroCode:      ne.String("LD"),
				DestinationMetroCode: ne.String("AM"),
				SourceZoneCode:       ne.String(""),
				DestinationZoneCode:  ne.String(""),
			},
		},
	}

	rawData := map[string]interface{}{
		networkDeviceLinkSchemaNames["Name"]:      ne.StringValue(expected.Name),
		networkDeviceLinkSchemaNames["Subnet"]:    ne.StringValue(expected.Subnet),
		networkDeviceLinkSchemaNames["ProjectID"]: ne.StringValue(expected.ProjectID),
	}
	d := schema.TestResourceDataRaw(t, createNetworkDeviceLinkResourceSchema(), rawData)
	d.Set(networkDeviceLinkSchemaNames["Devices"], flattenNetworkDeviceLinkDevices(nil, expected.Devices))
	d.Set(networkDeviceLinkSchemaNames["Links"], flattenNetworkDeviceLinkConnections(nil, expected.Links))
	// when
	result := createNetworkDeviceLink(d)
	// then
	assert.Equal(t, expected, result, "Created device link matches expected result")
}

func TestNetworkDeviceLink_updateResourceData(t *testing.T) {
	// given
	input := ne.DeviceLinkGroup{
		UUID:   ne.String("aae04283-10f9-4edb-9395-33681176592b"),
		Name:   ne.String("testGroup"),
		Subnet: ne.String("10.10.1.0/24"),
		Status: ne.String(ne.DeviceLinkGroupStatusProvisioned),
		Devices: []ne.DeviceLinkGroupDevice{
			{
				DeviceID:    ne.String("3eee8518-b19d-4de5-afd8-afd9b67e6e8c"),
				ASN:         ne.Int(0),
				InterfaceID: ne.Int(5),
			},
			{
				DeviceID:    ne.String("7c737fe8-3a9e-4ab9-afcc-c06d01db326d"),
				ASN:         ne.Int(0),
				InterfaceID: ne.Int(10),
			},
		},
		Links: []ne.DeviceLinkGroupLink{
			{
				AccountNumber:        ne.String(""),
				Throughput:           ne.String("1"),
				ThroughputUnit:       ne.String("Gbps"),
				SourceMetroCode:      ne.String("LD"),
				DestinationMetroCode: ne.String("AM"),
				SourceZoneCode:       ne.String(""),
				DestinationZoneCode:  ne.String(""),
			},
		},
	}
	d := schema.TestResourceDataRaw(t, createNetworkDeviceLinkResourceSchema(), make(map[string]interface{}))
	// when
	err := updateNetworkDeviceLinkResource(&input, d)
	// then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ne.StringValue(input.UUID), d.Get(networkDeviceLinkSchemaNames["UUID"]), "UUID matches")
	assert.Equal(t, ne.StringValue(input.Name), d.Get(networkDeviceLinkSchemaNames["Name"]), "Name matches")
	assert.Equal(t, ne.StringValue(input.Subnet), d.Get(networkDeviceLinkSchemaNames["Subnet"]), "Subnet matches")
	assert.Equal(t, ne.StringValue(input.Status), d.Get(networkDeviceLinkSchemaNames["Status"]), "Status matches")
	assert.Equal(t, input.Devices, expandNetworkDeviceLinkDevices(d.Get(networkDeviceLinkSchemaNames["Devices"]).(*schema.Set)), "Device matches")
	assert.Equal(t, input.Links, expandNetworkDeviceLinkConnections(d.Get(networkDeviceLinkSchemaNames["Links"]).(*schema.Set)), "Links matches")
}
