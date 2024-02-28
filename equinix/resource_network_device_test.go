package equinix

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/equinix/ne-go"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNetworkDevice_createFromResourceData(t *testing.T) {
	expectedPrimaryUserKey := ne.DeviceUserPublicKey{
		Username: ne.String("user"),
		KeyName:  ne.String("key"),
	}
	expectedPrimaryVendorConfig := map[string]string{
		"key": "value",
	}
	expectedPrimary := &ne.Device{
		Name:                  ne.String("device"),
		TypeCode:              ne.String("CSR1000V"),
		MetroCode:             ne.String("SV"),
		Throughput:            ne.Int(100),
		ThroughputUnit:        ne.String("Mbps"),
		HostName:              ne.String("test"),
		PackageCode:           ne.String("SEC"),
		Version:               ne.String("9.0.1"),
		IsBYOL:                ne.Bool(false),
		LicenseToken:          ne.String("sWf3df4gaAvbbexw45ga4f"),
		LicenseFile:           ne.String("/tmp/licenseFile"),
		ACLTemplateUUID:       ne.String("a624178c-6d59-4798-9a7f-2ddf2c7c5881"),
		AccountNumber:         ne.String("123456"),
		Notifications:         []string{"bla@bla.com"},
		PurchaseOrderNumber:   ne.String("1234567890"),
		TermLength:            ne.Int(1),
		AdditionalBandwidth:   ne.Int(50),
		OrderReference:        ne.String("12312121sddsf1231"),
		InterfaceCount:        ne.Int(10),
		WanInterfaceId:        ne.String("5"),
		CoreCount:             ne.Int(2),
		IsSelfManaged:         ne.Bool(false),
		VendorConfiguration:   expectedPrimaryVendorConfig,
		UserPublicKey:         &expectedPrimaryUserKey,
		Connectivity:          ne.String("INTERNET-ACCESS"),
		ProjectID:             ne.String("68ccfd49-39b1-478e-957a-67c72f719d7a"),
		DiverseFromDeviceUUID: ne.String("ed7891bd-15b4-4f72-ac56-d96cfdacddcc"),
	}
	rawData := map[string]interface{}{
		neDeviceSchemaNames["Name"]:                  ne.StringValue(expectedPrimary.Name),
		neDeviceSchemaNames["TypeCode"]:              ne.StringValue(expectedPrimary.TypeCode),
		neDeviceSchemaNames["MetroCode"]:             ne.StringValue(expectedPrimary.MetroCode),
		neDeviceSchemaNames["Throughput"]:            ne.IntValue(expectedPrimary.Throughput),
		neDeviceSchemaNames["Throughput"]:            ne.IntValue(expectedPrimary.Throughput),
		neDeviceSchemaNames["ThroughputUnit"]:        ne.StringValue(expectedPrimary.ThroughputUnit),
		neDeviceSchemaNames["HostName"]:              ne.StringValue(expectedPrimary.HostName),
		neDeviceSchemaNames["PackageCode"]:           ne.StringValue(expectedPrimary.PackageCode),
		neDeviceSchemaNames["Version"]:               ne.StringValue(expectedPrimary.Version),
		neDeviceSchemaNames["IsBYOL"]:                ne.BoolValue(expectedPrimary.IsBYOL),
		neDeviceSchemaNames["LicenseToken"]:          ne.StringValue(expectedPrimary.LicenseToken),
		neDeviceSchemaNames["LicenseFile"]:           ne.StringValue(expectedPrimary.LicenseFile),
		neDeviceSchemaNames["ACLTemplateUUID"]:       ne.StringValue(expectedPrimary.ACLTemplateUUID),
		neDeviceSchemaNames["AccountNumber"]:         ne.StringValue(expectedPrimary.AccountNumber),
		neDeviceSchemaNames["PurchaseOrderNumber"]:   ne.StringValue(expectedPrimary.PurchaseOrderNumber),
		neDeviceSchemaNames["TermLength"]:            ne.IntValue(expectedPrimary.TermLength),
		neDeviceSchemaNames["AdditionalBandwidth"]:   ne.IntValue(expectedPrimary.AdditionalBandwidth),
		neDeviceSchemaNames["OrderReference"]:        ne.StringValue(expectedPrimary.OrderReference),
		neDeviceSchemaNames["InterfaceCount"]:        ne.IntValue(expectedPrimary.InterfaceCount),
		neDeviceSchemaNames["WanInterfaceId"]:        ne.StringValue(expectedPrimary.WanInterfaceId),
		neDeviceSchemaNames["CoreCount"]:             ne.IntValue(expectedPrimary.CoreCount),
		neDeviceSchemaNames["IsSelfManaged"]:         ne.BoolValue(expectedPrimary.IsSelfManaged),
		neDeviceSchemaNames["ProjectID"]:             ne.StringValue(expectedPrimary.ProjectID),
		neDeviceSchemaNames["DiverseFromDeviceUUID"]: ne.StringValue(expectedPrimary.DiverseFromDeviceUUID),
	}
	d := schema.TestResourceDataRaw(t, createNetworkDeviceSchema(), rawData)
	d.Set(neDeviceSchemaNames["Notifications"], expectedPrimary.Notifications)
	d.Set(neDeviceSchemaNames["UserPublicKey"], flattenNetworkDeviceUserKeys([]*ne.DeviceUserPublicKey{&expectedPrimaryUserKey}))
	d.Set(neDeviceSchemaNames["VendorConfiguration"], expectedPrimary.VendorConfiguration)

	// when
	primary, secondary := createNetworkDevices(d)

	// then
	assert.NotNil(t, primary, "Primary device is not nil")
	assert.Nil(t, secondary, "Secondary device is nil")
	assert.Equal(t, expectedPrimary, primary, "Primary device matches expected result")
}

func TestNetworkDevice_updateResourceData(t *testing.T) {
	// given
	inputPrimary := &ne.Device{
		Name:                  ne.String("device"),
		TypeCode:              ne.String("CSR1000V"),
		ProjectID:             ne.String("68ccfd49-39b1-478e-957a-67c72f719d7a"),
		MetroCode:             ne.String("SV"),
		Throughput:            ne.Int(100),
		ThroughputUnit:        ne.String("Mbps"),
		HostName:              ne.String("test"),
		PackageCode:           ne.String("SEC"),
		Version:               ne.String("9.0.1"),
		IsBYOL:                ne.Bool(true),
		LicenseToken:          ne.String("sWf3df4gaAvbbexw45ga4f"),
		ACLTemplateUUID:       ne.String("a624178c-6d59-4798-9a7f-2ddf2c7c5881"),
		AccountNumber:         ne.String("123456"),
		Notifications:         []string{"bla@bla.com"},
		PurchaseOrderNumber:   ne.String("1234567890"),
		TermLength:            ne.Int(1),
		AdditionalBandwidth:   ne.Int(50),
		OrderReference:        ne.String("12312121sddsf1231"),
		InterfaceCount:        ne.Int(10),
		WanInterfaceId:        ne.String("6"),
		CoreCount:             ne.Int(2),
		IsSelfManaged:         ne.Bool(true),
		DiverseFromDeviceUUID: ne.String("68ccfd49-39b1-478e-957a-67c72f719d7a"),
		DiverseFromDeviceName: ne.String("diverseFromDeviceName"),
		VendorConfiguration: map[string]string{
			"key": "value",
		},
		UserPublicKey: &ne.DeviceUserPublicKey{
			Username: ne.String("user"),
			KeyName:  ne.String("key"),
		},
		ASN:      ne.Int(11222),
		ZoneCode: ne.String("Zone2"),
	}
	inputSecondary := &ne.Device{}
	secondarySchemaLicenseFile := "/tmp/licenseFileSec"
	d := schema.TestResourceDataRaw(t, createNetworkDeviceSchema(), make(map[string]interface{}))
	d.Set(neDeviceSchemaNames["Secondary"], flattenNetworkDeviceSecondary(&ne.Device{
		LicenseFile: ne.String(secondarySchemaLicenseFile),
	}))
	// when
	err := updateNetworkDeviceResource(inputPrimary, inputSecondary, d)

	// then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ne.StringValue(inputPrimary.Name), d.Get(neDeviceSchemaNames["Name"]), "Name matches")
	assert.Equal(t, ne.StringValue(inputPrimary.TypeCode), d.Get(neDeviceSchemaNames["TypeCode"]), "TypeCode matches")
	assert.Equal(t, ne.StringValue(inputPrimary.MetroCode), d.Get(neDeviceSchemaNames["MetroCode"]), "MetroCode matches")
	assert.Equal(t, ne.IntValue(inputPrimary.Throughput), d.Get(neDeviceSchemaNames["Throughput"]), "Throughput matches")
	assert.Equal(t, ne.StringValue(inputPrimary.ThroughputUnit), d.Get(neDeviceSchemaNames["ThroughputUnit"]), "ThroughputUnit matches")
	assert.Equal(t, ne.StringValue(inputPrimary.HostName), d.Get(neDeviceSchemaNames["HostName"]), "HostName matches")
	assert.Equal(t, ne.StringValue(inputPrimary.PackageCode), d.Get(neDeviceSchemaNames["PackageCode"]), "PackageCode matches")
	assert.Equal(t, ne.StringValue(inputPrimary.Version), d.Get(neDeviceSchemaNames["Version"]), "Version matches")
	assert.Equal(t, ne.BoolValue(inputPrimary.IsBYOL), d.Get(neDeviceSchemaNames["IsBYOL"]), "IsBYOL matches")
	assert.Empty(t, d.Get(neDeviceSchemaNames["LicenseToken"]), "LicenseToken is empty")
	assert.Equal(t, ne.StringValue(inputPrimary.ACLTemplateUUID), d.Get(neDeviceSchemaNames["ACLTemplateUUID"]), "ACLTemplateUUID matches")
	assert.Equal(t, ne.StringValue(inputPrimary.AccountNumber), d.Get(neDeviceSchemaNames["AccountNumber"]), "AccountNumber matches")
	assert.Equal(t, inputPrimary.Notifications, converters.SetToStringList(d.Get(neDeviceSchemaNames["Notifications"]).(*schema.Set)), "Notifications matches")
	assert.Equal(t, ne.StringValue(inputPrimary.PurchaseOrderNumber), d.Get(neDeviceSchemaNames["PurchaseOrderNumber"]), "PurchaseOrderNumber matches")
	assert.Equal(t, ne.IntValue(inputPrimary.TermLength), d.Get(neDeviceSchemaNames["TermLength"]), "TermLength matches")
	assert.Equal(t, ne.IntValue(inputPrimary.AdditionalBandwidth), d.Get(neDeviceSchemaNames["AdditionalBandwidth"]), "AdditionalBandwidth matches")
	assert.Equal(t, ne.StringValue(inputPrimary.OrderReference), d.Get(neDeviceSchemaNames["OrderReference"]), "OrderReference matches")
	assert.Equal(t, ne.IntValue(inputPrimary.InterfaceCount), d.Get(neDeviceSchemaNames["InterfaceCount"]), "InterfaceCount matches")
	assert.Empty(t, d.Get(neDeviceSchemaNames["WanInterfaceId"]), "Wan Interface Id is empty")
	assert.Equal(t, ne.IntValue(inputPrimary.CoreCount), d.Get(neDeviceSchemaNames["CoreCount"]), "CoreCount matches")
	assert.Equal(t, ne.BoolValue(inputPrimary.IsSelfManaged), d.Get(neDeviceSchemaNames["IsSelfManaged"]), "IsSelfManaged matches")
	assert.Equal(t, inputPrimary.VendorConfiguration, converters.InterfaceMapToStringMap(d.Get(neDeviceSchemaNames["VendorConfiguration"]).(map[string]interface{})), "VendorConfiguration matches")
	assert.Equal(t, inputPrimary.UserPublicKey, expandNetworkDeviceUserKeys(d.Get(neDeviceSchemaNames["UserPublicKey"]).(*schema.Set))[0], "UserPublicKey matches")
	assert.Equal(t, ne.IntValue(inputPrimary.ASN), d.Get(neDeviceSchemaNames["ASN"]), "ASN matches")
	assert.Equal(t, ne.StringValue(inputPrimary.ZoneCode), d.Get(neDeviceSchemaNames["ZoneCode"]), "ZoneCode matches")
	assert.Equal(t, ne.StringValue(inputPrimary.ProjectID), d.Get(neDeviceSchemaNames["ProjectID"]), "ProjectID matches")
	assert.Equal(t, ne.StringValue(inputPrimary.DiverseFromDeviceUUID), d.Get(neDeviceSchemaNames["DiverseFromDeviceUUID"]), "DiverseFromDeviceUUID matches")
	assert.Equal(t, ne.StringValue(inputPrimary.DiverseFromDeviceName), d.Get(neDeviceSchemaNames["DiverseFromDeviceName"]), "DiverseFromDeviceName matches")
	assert.Equal(t, secondarySchemaLicenseFile, ne.StringValue(expandNetworkDeviceSecondary(d.Get(neDeviceSchemaNames["Secondary"]).([]interface{})).LicenseFile), "Secondary LicenseFile matches")
}

func TestNetworkDevice_flattenSecondary(t *testing.T) {
	// given
	input := &ne.Device{
		UUID:                ne.String("0452fa68-8246-48b1-a1b2-817fb4baddcb"),
		Name:                ne.String("device"),
		Status:              ne.String(ne.DeviceStateProvisioned),
		LicenseStatus:       ne.String(ne.DeviceLicenseStateApplied),
		MetroCode:           ne.String("SV"),
		IBX:                 ne.String("SV5"),
		Region:              ne.String("AMER"),
		HostName:            ne.String("test"),
		LicenseToken:        ne.String("sWf3df4gaAvbbexw45ga4f"),
		LicenseFileID:       ne.String("d72dbe58-e596-4698-8b57-0a38e8077d25"),
		LicenseFile:         ne.String("/tmp/myfile"),
		CloudInitFileID:     ne.String("59a29dac-73fb-4e09-8aba-215b30fed63e"),
		ACLTemplateUUID:     ne.String("a624178c-6d59-4798-9a7f-2ddf2c7c5881"),
		SSHIPAddress:        ne.String("1.1.1.1"),
		SSHIPFqdn:           ne.String("test-1.1.1.1-SV.test.equinix.com"),
		AccountNumber:       ne.String("123456"),
		Notifications:       []string{"bla@bla.com"},
		RedundancyType:      ne.String("PRIMARY"),
		RedundantUUID:       ne.String("c2a147a3-ff47-4a24-a6e5-d6d7ce6459f3"),
		AdditionalBandwidth: ne.Int(50),
		ProjectID:           ne.String("68ccfd49-39b1-478e-957a-67c72f719d7a"),

		Interfaces: []ne.DeviceInterface{
			{
				ID:                ne.Int(1),
				Name:              ne.String("GigabitEthernet1"),
				Status:            ne.String("AVAILABLE"),
				OperationalStatus: ne.String("UP"),
				MACAddress:        ne.String("58-0A-C9-7A-DA-E9"),
				IPAddress:         ne.String("2.2.2.2"),
				AssignedType:      ne.String("test-connection(AWS Direct Connect)"),
				Type:              ne.String("DATA"),
			},
		},
		VendorConfiguration: map[string]string{
			"key": "value",
		},
		UserPublicKey: &ne.DeviceUserPublicKey{
			Username: ne.String("user"),
			KeyName:  ne.String("testKey"),
		},
		ASN:      ne.Int(11222),
		ZoneCode: ne.String("Zone2"),
	}
	expected := []interface{}{
		map[string]interface{}{
			neDeviceSchemaNames["UUID"]:                input.UUID,
			neDeviceSchemaNames["Name"]:                input.Name,
			neDeviceSchemaNames["Status"]:              input.Status,
			neDeviceSchemaNames["LicenseStatus"]:       input.LicenseStatus,
			neDeviceSchemaNames["MetroCode"]:           input.MetroCode,
			neDeviceSchemaNames["IBX"]:                 input.IBX,
			neDeviceSchemaNames["Region"]:              input.Region,
			neDeviceSchemaNames["HostName"]:            input.HostName,
			neDeviceSchemaNames["LicenseToken"]:        input.LicenseToken,
			neDeviceSchemaNames["LicenseFileID"]:       input.LicenseFileID,
			neDeviceSchemaNames["LicenseFile"]:         input.LicenseFile,
			neDeviceSchemaNames["CloudInitFileID"]:     input.CloudInitFileID,
			neDeviceSchemaNames["ACLTemplateUUID"]:     input.ACLTemplateUUID,
			neDeviceSchemaNames["SSHIPAddress"]:        input.SSHIPAddress,
			neDeviceSchemaNames["SSHIPFqdn"]:           input.SSHIPFqdn,
			neDeviceSchemaNames["AccountNumber"]:       input.AccountNumber,
			neDeviceSchemaNames["Notifications"]:       input.Notifications,
			neDeviceSchemaNames["RedundancyType"]:      input.RedundancyType,
			neDeviceSchemaNames["RedundantUUID"]:       input.RedundantUUID,
			neDeviceSchemaNames["AdditionalBandwidth"]: input.AdditionalBandwidth,
			neDeviceSchemaNames["ProjectID"]:           input.ProjectID,
			neDeviceSchemaNames["Interfaces"]: []interface{}{
				map[string]interface{}{
					neDeviceInterfaceSchemaNames["ID"]:                input.Interfaces[0].ID,
					neDeviceInterfaceSchemaNames["Name"]:              input.Interfaces[0].Name,
					neDeviceInterfaceSchemaNames["Status"]:            input.Interfaces[0].Status,
					neDeviceInterfaceSchemaNames["OperationalStatus"]: input.Interfaces[0].OperationalStatus,
					neDeviceInterfaceSchemaNames["MACAddress"]:        input.Interfaces[0].MACAddress,
					neDeviceInterfaceSchemaNames["IPAddress"]:         input.Interfaces[0].IPAddress,
					neDeviceInterfaceSchemaNames["AssignedType"]:      input.Interfaces[0].AssignedType,
					neDeviceInterfaceSchemaNames["Type"]:              input.Interfaces[0].Type,
				},
			},
			neDeviceSchemaNames["VendorConfiguration"]: map[string]string{
				"key": "value",
			},
			neDeviceSchemaNames["UserPublicKey"]: []interface{}{
				map[string]interface{}{
					neDeviceUserKeySchemaNames["Username"]: input.UserPublicKey.Username,
					neDeviceUserKeySchemaNames["KeyName"]:  input.UserPublicKey.KeyName,
				},
			},
			neDeviceSchemaNames["ASN"]:      input.ASN,
			neDeviceSchemaNames["ZoneCode"]: input.ZoneCode,
		},
	}
	// when
	out := flattenNetworkDeviceSecondary(input)
	// then
	assert.NotNil(t, out, "Output is not nil")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestNetworkDevice_expandSecondary(t *testing.T) {
	// given
	f := func(i interface{}) int {
		str := fmt.Sprintf("%v", i)
		return schema.HashString(str)
	}
	input := []interface{}{
		map[string]interface{}{
			neDeviceSchemaNames["UUID"]:                "0452fa68-8246-48b1-a1b2-817fb4baddcb",
			neDeviceSchemaNames["Name"]:                "device",
			neDeviceSchemaNames["MetroCode"]:           "SV",
			neDeviceSchemaNames["HostName"]:            "SV5",
			neDeviceSchemaNames["LicenseToken"]:        "sWf3df4gaAvbbexw45ga4f",
			neDeviceSchemaNames["LicenseFile"]:         "/tmp/licenseFile",
			neDeviceSchemaNames["ACLTemplateUUID"]:     "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
			neDeviceSchemaNames["AccountNumber"]:       "123456",
			neDeviceSchemaNames["Notifications"]:       schema.NewSet(schema.HashString, []interface{}{"bla@bla.com"}),
			neDeviceSchemaNames["AdditionalBandwidth"]: 50,
			neDeviceSchemaNames["VendorConfiguration"]: map[string]interface{}{
				"key": "value",
			},
			neDeviceSchemaNames["UserPublicKey"]: schema.NewSet(f, []interface{}{
				map[string]interface{}{
					neDeviceUserKeySchemaNames["Username"]: "user",
					neDeviceUserKeySchemaNames["KeyName"]:  "testKey",
				},
			}),
		},
	}
	expected := &ne.Device{
		UUID:                ne.String(input[0].(map[string]interface{})[neDeviceSchemaNames["UUID"]].(string)),
		Name:                ne.String(input[0].(map[string]interface{})[neDeviceSchemaNames["Name"]].(string)),
		MetroCode:           ne.String(input[0].(map[string]interface{})[neDeviceSchemaNames["MetroCode"]].(string)),
		HostName:            ne.String(input[0].(map[string]interface{})[neDeviceSchemaNames["HostName"]].(string)),
		LicenseToken:        ne.String(input[0].(map[string]interface{})[neDeviceSchemaNames["LicenseToken"]].(string)),
		LicenseFile:         ne.String(input[0].(map[string]interface{})[neDeviceSchemaNames["LicenseFile"]].(string)),
		ACLTemplateUUID:     ne.String(input[0].(map[string]interface{})[neDeviceSchemaNames["ACLTemplateUUID"]].(string)),
		AccountNumber:       ne.String(input[0].(map[string]interface{})[neDeviceSchemaNames["AccountNumber"]].(string)),
		Notifications:       converters.SetToStringList(input[0].(map[string]interface{})[neDeviceSchemaNames["Notifications"]].(*schema.Set)),
		AdditionalBandwidth: ne.Int(input[0].(map[string]interface{})[neDeviceSchemaNames["AdditionalBandwidth"]].(int)),
		VendorConfiguration: map[string]string{
			"key": "value",
		},
		UserPublicKey: expandNetworkDeviceUserKeys(input[0].(map[string]interface{})[neDeviceSchemaNames["UserPublicKey"]].(*schema.Set))[0],
	}
	// when
	out := expandNetworkDeviceSecondary(input)
	// then
	assert.NotNil(t, out, "Output is not empty")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestNetworkDevice_uploadLicenseFile(t *testing.T) {
	// given
	fileName := "test.lic"
	licenseFileID := "someTestID"
	device := &ne.Device{LicenseFile: ne.String("/path/to/" + fileName), MetroCode: ne.String("SV"), TypeCode: ne.String("VSRX")}
	var rxMetroCode, rxFileName, rxTypeCode, rxMgmtMode, rxLicMode string
	uploadFunc := func(metroCode, deviceTypeCode, deviceManagementMode, licenseMode, fileName string, reader io.Reader) (*string, error) {
		rxMetroCode = metroCode
		rxFileName = fileName
		rxTypeCode = deviceTypeCode
		rxMgmtMode = deviceManagementMode
		rxLicMode = licenseMode
		return &licenseFileID, nil
	}
	openFunc := func(name string) (*os.File, error) {
		return &os.File{}, nil
	}
	// when
	err := uploadDeviceLicenseFile(openFunc, uploadFunc, ne.StringValue(device.TypeCode), device)
	// then
	assert.Nil(t, err, "License upload function does not return any error")
	assert.Equal(t, licenseFileID, ne.StringValue(device.LicenseFileID), "Device LicenseFileID matches")
	assert.Equal(t, ne.StringValue(device.MetroCode), rxMetroCode, "Received metroCode matches")
	assert.Equal(t, ne.StringValue(device.TypeCode), rxTypeCode, "Received typeCode matches")
	assert.Equal(t, fileName, rxFileName, "Received fileName matches")
	assert.Equal(t, ne.DeviceManagementTypeSelf, rxMgmtMode, "Received management mode matches")
	assert.Equal(t, ne.DeviceLicenseModeBYOL, rxLicMode, "Received management mode matches")
}

func TestNetworkDevice_statusProvisioningWaitConfiguration(t *testing.T) {
	// given
	deviceID := "test"
	var queriedDeviceID string
	fetchFunc := func(uuid string) (*ne.Device, error) {
		queriedDeviceID = uuid
		return &ne.Device{Status: ne.String(ne.DeviceStateProvisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	// when
	waitConfig := createNetworkDeviceStatusProvisioningWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	// then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, queriedDeviceID, "Queried device ID matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_statusDeleteWaitConfiguration(t *testing.T) {
	// given
	deviceID := "test"
	var queriedDeviceID string
	fetchFunc := func(uuid string) (*ne.Device, error) {
		queriedDeviceID = uuid
		return &ne.Device{Status: ne.String(ne.DeviceStateDeprovisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	// when
	waitConfig := createNetworkDeviceStatusDeleteWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	// then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, queriedDeviceID, "Queried device ID matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_statusResourceUpgradeWaitConfiguration(t *testing.T) {
	// given
	deviceID := "test"
	var queriedDeviceID string
	fetchFunc := func(uuid string) (*ne.Device, error) {
		queriedDeviceID = uuid
		return &ne.Device{Status: ne.String(ne.DeviceStateProvisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	// when
	waitConfig := createNetworkDeviceStatusResourceUpgradeWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	// then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, queriedDeviceID, "Queried device ID matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_licenseStatusWaitConfiguration(t *testing.T) {
	// given
	deviceID := "test"
	var queriedDeviceID string
	fetchFunc := func(uuid string) (*ne.Device, error) {
		queriedDeviceID = uuid
		return &ne.Device{LicenseStatus: ne.String(ne.DeviceLicenseStateApplied)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	// when
	waitConfig := createNetworkDeviceLicenseStatusWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	// then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, queriedDeviceID, "Queried device ID matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_ACLStatusWaitConfiguration(t *testing.T) {
	// given
	deviceUUID := "test"
	var receivedDeviceUUID string
	fetchFunc := func(uuid string) (*ne.DeviceACLDetails, error) {
		receivedDeviceUUID = uuid
		return &ne.DeviceACLDetails{Status: ne.String(ne.ACLDeviceStatusProvisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	// when
	waitConfig := createNetworkDeviceACLStatusWaitConfiguration(fetchFunc, deviceUUID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	// then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceUUID, receivedDeviceUUID, "Queried Device id matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_AdditionalBandwidthStatusWaitConfiguration(t *testing.T) {
	// given
	deviceID := "test"
	var receivedID string
	fetchFunc := func(uuid string) (*ne.DeviceAdditionalBandwidthDetails, error) {
		receivedID = uuid
		return &ne.DeviceAdditionalBandwidthDetails{Status: ne.String(ne.DeviceAdditionalBandwidthStatusProvisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	// when
	waitConfig := createNetworkDeviceAdditionalBandwidthStatusWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	// then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, receivedID, "Queried Additional Bandwidth device id matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Additional bandwidth status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Additional bandwidth wait configuration min timeout matches")
}
