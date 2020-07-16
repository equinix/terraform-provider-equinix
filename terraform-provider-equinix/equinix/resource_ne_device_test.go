package equinix

import (
	"ne-go"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

var primaryDevFields = []string{"AccountNumber", "ACL", "AdditionalBandwidth", "DeviceTypeCode", "HostName", "LicenseFileID", "LicenseKey", "LicenseSecret", "LicenseToken", "LicenseType", "MetroCode", "Name", "Notifications", "PackageCode", "PurchaseOrderNumber", "TermLength", "Throughput", "ThroughputUnit", "Version", "CoreCount", "InterfaceCount"}
var secondaryDevFields = []string{"AccountNumber", "ACL", "AdditionalBandwidth", "HostName", "LicenseFileID", "LicenseKey", "LicenseSecret", "LicenseToken", "MetroCode", "Name", "Notifications"}
var vendorConfigFields = []string{"SiteID", "SystemIPAddress"}

func TestNeDevice_resourceDataFromDomain(t *testing.T) {
	//Given
	d := schema.TestResourceDataRaw(t, createNeDeviceResourceSchema(), make(map[string]interface{}))
	primary := ne.Device{
		AccountNumber:       "123456",
		ACL:                 []string{"1.1.1.1/32", "2.2.2.2/32"},
		AdditionalBandwidth: 100,
		DeviceTypeCode:      "CSR1000V",
		HostName:            "testDevice",
		LicenseFileID:       "someFileId",
		LicenseKey:          "someKey",
		LicenseSecret:       "someSecret",
		LicenseToken:        "someToken",
		LicenseType:         "BYOL",
		MetroCode:           "SV",
		Name:                "test",
		Notifications:       []string{"a@a.com", "b@b.com"},
		PackageCode:         "SEC",
		PurchaseOrderNumber: "aa-bb-cc",
		TermLength:          3,
		Throughput:          5,
		ThroughputUnit:      "Gbps",
		VendorConfig: &ne.DeviceVendorConfig{
			SiteID:          "someSystemID",
			SystemIPAddress: "someSystemIPAddress"},
		Version:        "16.12.1e",
		CoreCount:      4,
		InterfaceCount: 10,
	}
	secondary := ne.Device{
		AccountNumber:       "7890",
		ACL:                 []string{"3.3.3.3/32", "4.4.4.4/32"},
		AdditionalBandwidth: 200,
		HostName:            "testDevice-sec",
		LicenseFileID:       "someFileId",
		LicenseKey:          "someKey",
		LicenseSecret:       "someSecret",
		LicenseToken:        "someToken",
		MetroCode:           "DC",
		Name:                "test-sec",
		Notifications:       []string{"c@c.com", "d@d.com"},
		VendorConfig: &ne.DeviceVendorConfig{
			SiteID:          "someSystemID",
			SystemIPAddress: "someSystemIPAddress"},
	}

	//When
	err := updateNeDeviceResource(&primary, &secondary, d)

	//Then
	assert.Nil(t, err, "Schema update should not return an error")
	sourceMatchesTargetSchema(t, primary, primaryDevFields, d, neDeviceSchemaNames)

	vendorConf := d.Get(neDeviceSchemaNames["VendorConfig"]).(*schema.Set)
	assert.Equal(t, 1, vendorConf.Len(), "There is one vendor configuration defined")
	sourceMatchesTargetSchema(t, *primary.VendorConfig, vendorConfigFields, vendorConf.List()[0], neDeviceVendorSchemaNames)

	secConns := d.Get(neDeviceSchemaNames["Secondary"]).(*schema.Set)
	assert.Equal(t, 1, secConns.Len(), "There is one secondary connection defined")
	sourceMatchesTargetSchema(t, secondary, secondaryDevFields, secConns.List()[0], neDeviceSchemaNames)

	secVendorConf := secConns.List()[0].(map[string]interface{})[neDeviceSchemaNames["VendorConfig"]].(*schema.Set)
	assert.Equal(t, 1, secVendorConf.Len(), "There is one vendor configuration defined for secondary device")
	sourceMatchesTargetSchema(t, *secondary.VendorConfig, vendorConfigFields, secVendorConf.List()[0], neDeviceVendorSchemaNames)
}

func TestNeDevice_domainFromResourceData(t *testing.T) {
	//given
	d := schema.TestResourceDataRaw(t, createNeDeviceResourceSchema(), make(map[string]interface{}))
	d.Set(neDeviceSchemaNames["AccountNumber"], "123456")
	d.Set(neDeviceSchemaNames["ACL"], []string{"1.1.1.1/32", "2.2.2.2/32"})
	d.Set(neDeviceSchemaNames["AdditionalBandwidth"], 100)
	d.Set(neDeviceSchemaNames["DeviceTypeCode"], "CSR1000V")
	d.Set(neDeviceSchemaNames["HostName"], "testDevice")
	d.Set(neDeviceSchemaNames["LicenseFileID"], "someFileId")
	d.Set(neDeviceSchemaNames["LicenseKey"], "someKey")
	d.Set(neDeviceSchemaNames["LicenseSecret"], "someSecret")
	d.Set(neDeviceSchemaNames["LicenseToken"], "someToken")
	d.Set(neDeviceSchemaNames["LicenseType"], "BYOL")
	d.Set(neDeviceSchemaNames["MetroCode"], "SV")
	d.Set(neDeviceSchemaNames["Name"], "test")
	d.Set(neDeviceSchemaNames["Notifications"], []string{"a@a.com", "b@b.com"})
	d.Set(neDeviceSchemaNames["PackageCode"], "SEC")
	d.Set(neDeviceSchemaNames["PurchaseOrderNumber"], "aa-bb-cc")
	d.Set(neDeviceSchemaNames["TermLength"], 3)
	d.Set(neDeviceSchemaNames["Throughput"], 5)
	d.Set(neDeviceSchemaNames["ThroughputUnit"], "Gbps")
	d.Set(neDeviceSchemaNames["Version"], "16.12.1e")
	d.Set(neDeviceSchemaNames["CoreCount"], 2)
	d.Set(neDeviceSchemaNames["InterfaceCount"], 12)
	priVendorConf := flattenNeDeviceVendorConfig(&ne.DeviceVendorConfig{
		SiteID:          "someSiteID",
		SystemIPAddress: "someSystemIPAddress"})
	d.Set(neDeviceSchemaNames["VendorConfig"], priVendorConf)
	secDev := make(map[string]interface{})
	secDev[neDeviceSchemaNames["AccountNumber"]] = "7890"
	secDev[neDeviceSchemaNames["ACL"]] = []string{"3.3.3.3/32", "4.4.4.4/32"}
	secDev[neDeviceSchemaNames["AdditionalBandwidth"]] = 200
	secDev[neDeviceSchemaNames["HostName"]] = "testDevice-sec"
	secDev[neDeviceSchemaNames["LicenseFileID"]] = "someFileId"
	secDev[neDeviceSchemaNames["LicenseKey"]] = "someKey"
	secDev[neDeviceSchemaNames["LicenseSecret"]] = "someSecret"
	secDev[neDeviceSchemaNames["LicenseToken"]] = "someToken"
	secDev[neDeviceSchemaNames["MetroCode"]] = "DC"
	secDev[neDeviceSchemaNames["Name"]] = "test-sec"
	secDev[neDeviceSchemaNames["Notifications"]] = []string{"c@c.com", "d@d.com"}
	secVendorConf := flattenNeDeviceVendorConfig(&ne.DeviceVendorConfig{
		SiteID:          "someSiteIDSec",
		SystemIPAddress: "someSystemIPAddressSec"})
	secDev[neDeviceSchemaNames["VendorConfig"]] = secVendorConf
	d.Set(neDeviceSchemaNames["Secondary"], []map[string]interface{}{secDev})

	//when
	primary, secondary := createNeDevices(d)

	//then
	assert.NotNil(t, primary, "Primary device should be present")
	sourceMatchesTargetSchema(t, *primary, primaryDevFields, d, neDeviceSchemaNames)
	sourceMatchesTargetSchema(t, *primary.VendorConfig, vendorConfigFields, priVendorConf.([]map[string]interface{})[0], neDeviceVendorSchemaNames)

	assert.NotNil(t, secondary, "Secondary device should be present")
	sourceMatchesTargetSchema(t, *secondary, secondaryDevFields, secDev, neDeviceSchemaNames)
	sourceMatchesTargetSchema(t, *secondary.VendorConfig, vendorConfigFields, secVendorConf.([]map[string]interface{})[0], neDeviceVendorSchemaNames)
}
