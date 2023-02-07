package equinix

import (
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNetworkFile_createFromResourceData(t *testing.T) {
	// given
	expected := map[string]string{
		"FileName":             "testFileName",
		"Content":              "testContent",
		"MetroCode":            "testMetroCode",
		"DeviceTypeCode":       "testDeviceTypeCode",
		"ProcessType":          "testProcessType",
		"DeviceManagementType": ne.DeviceManagementTypeSelf,
		"LicenseMode":          ne.DeviceLicenseModeBYOL,
	}
	rawData := map[string]interface{}{
		networkFileSchemaNames["FileName"]:       expected["FileName"],
		networkFileSchemaNames["Content"]:        expected["Content"],
		networkFileSchemaNames["MetroCode"]:      expected["MetroCode"],
		networkFileSchemaNames["DeviceTypeCode"]: expected["DeviceTypeCode"],
		networkFileSchemaNames["ProcessType"]:    expected["ProcessType"],
		networkFileSchemaNames["IsSelfManaged"]:  true,
		networkFileSchemaNames["IsBYOL"]:         true,
	}
	d := schema.TestResourceDataRaw(t, createNetworkFileSchema(), rawData)
	// when
	fileRequest := createFileRequest(d)
	// then
	assert.Equal(t, expected, fileRequest, "Created file request matches expected result")
}

func TestNetworkFile_updateResourceData(t *testing.T) {
	// given
	input := &ne.File{
		UUID:           ne.String("183b9c9a-251b-4720-85aa-aa80269a0ffe"),
		FileName:       ne.String("testFileName"),
		MetroCode:      ne.String("testMetroCode"),
		DeviceTypeCode: ne.String("testDeviceTypeCode"),
		ProcessType:    ne.String("testProcessType"),
		Status:         ne.String("testStatus"),
	}
	d := schema.TestResourceDataRaw(t, createNetworkFileSchema(), make(map[string]interface{}))
	// when
	err := updateFileResource(input, d)
	// then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ne.StringValue(input.UUID), d.Get(networkFileSchemaNames["UUID"]), "UUID matches")
	assert.Equal(t, ne.StringValue(input.FileName), d.Get(networkFileSchemaNames["FileName"]), "FileName matches")
	assert.Equal(t, ne.StringValue(input.MetroCode), d.Get(networkFileSchemaNames["MetroCode"]), "MetroCode matches")
	assert.Equal(t, ne.StringValue(input.DeviceTypeCode), d.Get(networkFileSchemaNames["DeviceTypeCode"]), "DeviceTypeCode matches")
	assert.Equal(t, ne.StringValue(input.ProcessType), d.Get(networkFileSchemaNames["ProcessType"]), "ProcessType matches")
	assert.Equal(t, ne.StringValue(input.Status), d.Get(networkFileSchemaNames["Status"]), "Status matches")
}
