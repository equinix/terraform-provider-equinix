package equinix

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/equinix/ne-go"
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
		Name:                ne.String("device"),
		TypeCode:            ne.String("CSR1000V"),
		MetroCode:           ne.String("SV"),
		Throughput:          ne.Int(100),
		ThroughputUnit:      ne.String("Mbps"),
		HostName:            ne.String("test"),
		PackageCode:         ne.String("SEC"),
		Version:             ne.String("9.0.1"),
		IsBYOL:              ne.Bool(false),
		LicenseToken:        ne.String("sWf3df4gaAvbbexw45ga4f"),
		LicenseFile:         ne.String("/tmp/licenseFile"),
		ACLTemplateUUID:     ne.String("a624178c-6d59-4798-9a7f-2ddf2c7c5881"),
		AccountNumber:       ne.String("123456"),
		Notifications:       []string{"bla@bla.com"},
		PurchaseOrderNumber: ne.String("1234567890"),
		TermLength:          ne.Int(1),
		AdditionalBandwidth: ne.Int(50),
		OrderReference:      ne.String("12312121sddsf1231"),
		InterfaceCount:      ne.Int(10),
		WanInterfaceId:      ne.String("5"),
		CoreCount:           ne.Int(2),
		IsSelfManaged:       ne.Bool(false),
		VendorConfiguration: expectedPrimaryVendorConfig,
		UserPublicKey:       &expectedPrimaryUserKey,
	}
	rawData := map[string]interface{}{
		networkDeviceSchemaNames["Name"]:                ne.StringValue(expectedPrimary.Name),
		networkDeviceSchemaNames["TypeCode"]:            ne.StringValue(expectedPrimary.TypeCode),
		networkDeviceSchemaNames["MetroCode"]:           ne.StringValue(expectedPrimary.MetroCode),
		networkDeviceSchemaNames["Throughput"]:          ne.IntValue(expectedPrimary.Throughput),
		networkDeviceSchemaNames["Throughput"]:          ne.IntValue(expectedPrimary.Throughput),
		networkDeviceSchemaNames["ThroughputUnit"]:      ne.StringValue(expectedPrimary.ThroughputUnit),
		networkDeviceSchemaNames["HostName"]:            ne.StringValue(expectedPrimary.HostName),
		networkDeviceSchemaNames["PackageCode"]:         ne.StringValue(expectedPrimary.PackageCode),
		networkDeviceSchemaNames["Version"]:             ne.StringValue(expectedPrimary.Version),
		networkDeviceSchemaNames["IsBYOL"]:              ne.BoolValue(expectedPrimary.IsBYOL),
		networkDeviceSchemaNames["LicenseToken"]:        ne.StringValue(expectedPrimary.LicenseToken),
		networkDeviceSchemaNames["LicenseFile"]:         ne.StringValue(expectedPrimary.LicenseFile),
		networkDeviceSchemaNames["ACLTemplateUUID"]:     ne.StringValue(expectedPrimary.ACLTemplateUUID),
		networkDeviceSchemaNames["AccountNumber"]:       ne.StringValue(expectedPrimary.AccountNumber),
		networkDeviceSchemaNames["PurchaseOrderNumber"]: ne.StringValue(expectedPrimary.PurchaseOrderNumber),
		networkDeviceSchemaNames["TermLength"]:          ne.IntValue(expectedPrimary.TermLength),
		networkDeviceSchemaNames["AdditionalBandwidth"]: ne.IntValue(expectedPrimary.AdditionalBandwidth),
		networkDeviceSchemaNames["OrderReference"]:      ne.StringValue(expectedPrimary.OrderReference),
		networkDeviceSchemaNames["InterfaceCount"]:      ne.IntValue(expectedPrimary.InterfaceCount),
		networkDeviceSchemaNames["WanInterfaceId"]:      ne.StringValue(expectedPrimary.WanInterfaceId),
		networkDeviceSchemaNames["CoreCount"]:           ne.IntValue(expectedPrimary.CoreCount),
		networkDeviceSchemaNames["IsSelfManaged"]:       ne.BoolValue(expectedPrimary.IsSelfManaged),
	}
	d := schema.TestResourceDataRaw(t, createNetworkDeviceSchema(), rawData)
	d.Set(networkDeviceSchemaNames["Notifications"], expectedPrimary.Notifications)
	d.Set(networkDeviceSchemaNames["UserPublicKey"], flattenNetworkDeviceUserKeys([]*ne.DeviceUserPublicKey{&expectedPrimaryUserKey}))
	d.Set(networkDeviceSchemaNames["VendorConfiguration"], expectedPrimary.VendorConfiguration)

	//when
	primary, secondary := createNetworkDevices(d)

	//then
	assert.NotNil(t, primary, "Primary device is not nil")
	assert.Nil(t, secondary, "Secondary device is nil")
	assert.Equal(t, expectedPrimary, primary, "Primary device matches expected result")
}

func TestNetworkDevice_updateResourceData(t *testing.T) {
	//given
	inputPrimary := &ne.Device{
		Name:                ne.String("device"),
		TypeCode:            ne.String("CSR1000V"),
		MetroCode:           ne.String("SV"),
		Throughput:          ne.Int(100),
		ThroughputUnit:      ne.String("Mbps"),
		HostName:            ne.String("test"),
		PackageCode:         ne.String("SEC"),
		Version:             ne.String("9.0.1"),
		IsBYOL:              ne.Bool(true),
		LicenseToken:        ne.String("sWf3df4gaAvbbexw45ga4f"),
		ACLTemplateUUID:     ne.String("a624178c-6d59-4798-9a7f-2ddf2c7c5881"),
		AccountNumber:       ne.String("123456"),
		Notifications:       []string{"bla@bla.com"},
		PurchaseOrderNumber: ne.String("1234567890"),
		TermLength:          ne.Int(1),
		AdditionalBandwidth: ne.Int(50),
		OrderReference:      ne.String("12312121sddsf1231"),
		InterfaceCount:      ne.Int(10),
		CoreCount:           ne.Int(2),
		IsSelfManaged:       ne.Bool(true),
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
	d.Set(networkDeviceSchemaNames["Secondary"], flattenNetworkDeviceSecondary(&ne.Device{
		LicenseFile: ne.String(secondarySchemaLicenseFile),
	}))
	//when
	err := updateNetworkDeviceResource(inputPrimary, inputSecondary, d)

	//then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ne.StringValue(inputPrimary.Name), d.Get(networkDeviceSchemaNames["Name"]), "Name matches")
	assert.Equal(t, ne.StringValue(inputPrimary.TypeCode), d.Get(networkDeviceSchemaNames["TypeCode"]), "TypeCode matches")
	assert.Equal(t, ne.StringValue(inputPrimary.MetroCode), d.Get(networkDeviceSchemaNames["MetroCode"]), "MetroCode matches")
	assert.Equal(t, ne.IntValue(inputPrimary.Throughput), d.Get(networkDeviceSchemaNames["Throughput"]), "Throughput matches")
	assert.Equal(t, ne.StringValue(inputPrimary.ThroughputUnit), d.Get(networkDeviceSchemaNames["ThroughputUnit"]), "ThroughputUnit matches")
	assert.Equal(t, ne.StringValue(inputPrimary.HostName), d.Get(networkDeviceSchemaNames["HostName"]), "HostName matches")
	assert.Equal(t, ne.StringValue(inputPrimary.PackageCode), d.Get(networkDeviceSchemaNames["PackageCode"]), "PackageCode matches")
	assert.Equal(t, ne.StringValue(inputPrimary.Version), d.Get(networkDeviceSchemaNames["Version"]), "Version matches")
	assert.Equal(t, ne.BoolValue(inputPrimary.IsBYOL), d.Get(networkDeviceSchemaNames["IsBYOL"]), "IsBYOL matches")
	assert.Equal(t, ne.StringValue(inputPrimary.LicenseToken), d.Get(networkDeviceSchemaNames["LicenseToken"]), "LicenseToken matches")
	assert.Equal(t, ne.StringValue(inputPrimary.ACLTemplateUUID), d.Get(networkDeviceSchemaNames["ACLTemplateUUID"]), "ACLTemplateUUID matches")
	assert.Equal(t, ne.StringValue(inputPrimary.AccountNumber), d.Get(networkDeviceSchemaNames["AccountNumber"]), "AccountNumber matches")
	assert.Equal(t, inputPrimary.Notifications, expandSetToStringList(d.Get(networkDeviceSchemaNames["Notifications"]).(*schema.Set)), "Notifications matches")
	assert.Equal(t, ne.StringValue(inputPrimary.PurchaseOrderNumber), d.Get(networkDeviceSchemaNames["PurchaseOrderNumber"]), "PurchaseOrderNumber matches")
	assert.Equal(t, ne.IntValue(inputPrimary.TermLength), d.Get(networkDeviceSchemaNames["TermLength"]), "TermLength matches")
	assert.Equal(t, ne.IntValue(inputPrimary.AdditionalBandwidth), d.Get(networkDeviceSchemaNames["AdditionalBandwidth"]), "AdditionalBandwidth matches")
	assert.Equal(t, ne.StringValue(inputPrimary.OrderReference), d.Get(networkDeviceSchemaNames["OrderReference"]), "OrderReference matches")
	assert.Equal(t, ne.IntValue(inputPrimary.InterfaceCount), d.Get(networkDeviceSchemaNames["InterfaceCount"]), "InterfaceCount matches")
	assert.Equal(t, ne.IntValue(inputPrimary.CoreCount), d.Get(networkDeviceSchemaNames["CoreCount"]), "CoreCount matches")
	assert.Equal(t, ne.BoolValue(inputPrimary.IsSelfManaged), d.Get(networkDeviceSchemaNames["IsSelfManaged"]), "IsSelfManaged matches")
	assert.Equal(t, inputPrimary.VendorConfiguration, expandInterfaceMapToStringMap(d.Get(networkDeviceSchemaNames["VendorConfiguration"]).(map[string]interface{})), "VendorConfiguration matches")
	assert.Equal(t, inputPrimary.UserPublicKey, expandNetworkDeviceUserKeys(d.Get(networkDeviceSchemaNames["UserPublicKey"]).(*schema.Set))[0], "UserPublicKey matches")
	assert.Equal(t, ne.IntValue(inputPrimary.ASN), d.Get(networkDeviceSchemaNames["ASN"]), "ASN matches")
	assert.Equal(t, ne.StringValue(inputPrimary.ZoneCode), d.Get(networkDeviceSchemaNames["ZoneCode"]), "ZoneCode matches")
	assert.Equal(t, secondarySchemaLicenseFile, ne.StringValue(expandNetworkDeviceSecondary(d.Get(networkDeviceSchemaNames["Secondary"]).([]interface{})).LicenseFile), "Secondary LicenseFile matches")
}

func TestNetworkDevice_flattenSecondary(t *testing.T) {
	//given
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
		ACLTemplateUUID:     ne.String("a624178c-6d59-4798-9a7f-2ddf2c7c5881"),
		SSHIPAddress:        ne.String("1.1.1.1"),
		SSHIPFqdn:           ne.String("test-1.1.1.1-SV.test.equinix.com"),
		AccountNumber:       ne.String("123456"),
		Notifications:       []string{"bla@bla.com"},
		RedundancyType:      ne.String("PRIMARY"),
		RedundantUUID:       ne.String("c2a147a3-ff47-4a24-a6e5-d6d7ce6459f3"),
		AdditionalBandwidth: ne.Int(50),
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
			networkDeviceSchemaNames["UUID"]:                input.UUID,
			networkDeviceSchemaNames["Name"]:                input.Name,
			networkDeviceSchemaNames["Status"]:              input.Status,
			networkDeviceSchemaNames["LicenseStatus"]:       input.LicenseStatus,
			networkDeviceSchemaNames["MetroCode"]:           input.MetroCode,
			networkDeviceSchemaNames["IBX"]:                 input.IBX,
			networkDeviceSchemaNames["Region"]:              input.Region,
			networkDeviceSchemaNames["HostName"]:            input.HostName,
			networkDeviceSchemaNames["LicenseToken"]:        input.LicenseToken,
			networkDeviceSchemaNames["LicenseFileID"]:       input.LicenseFileID,
			networkDeviceSchemaNames["LicenseFile"]:         input.LicenseFile,
			networkDeviceSchemaNames["ACLTemplateUUID"]:     input.ACLTemplateUUID,
			networkDeviceSchemaNames["SSHIPAddress"]:        input.SSHIPAddress,
			networkDeviceSchemaNames["SSHIPFqdn"]:           input.SSHIPFqdn,
			networkDeviceSchemaNames["AccountNumber"]:       input.AccountNumber,
			networkDeviceSchemaNames["Notifications"]:       input.Notifications,
			networkDeviceSchemaNames["RedundancyType"]:      input.RedundancyType,
			networkDeviceSchemaNames["RedundantUUID"]:       input.RedundantUUID,
			networkDeviceSchemaNames["AdditionalBandwidth"]: input.AdditionalBandwidth,
			networkDeviceSchemaNames["Interfaces"]: []interface{}{
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
			networkDeviceSchemaNames["VendorConfiguration"]: map[string]string{
				"key": "value",
			},
			networkDeviceSchemaNames["UserPublicKey"]: []interface{}{
				map[string]interface{}{
					neDeviceUserKeySchemaNames["Username"]: input.UserPublicKey.Username,
					neDeviceUserKeySchemaNames["KeyName"]:  input.UserPublicKey.KeyName,
				},
			},
			networkDeviceSchemaNames["ASN"]:      input.ASN,
			networkDeviceSchemaNames["ZoneCode"]: input.ZoneCode,
		},
	}
	//when
	out := flattenNetworkDeviceSecondary(input)
	//then
	assert.NotNil(t, out, "Output is not nil")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestNetworkDevice_expandSecondary(t *testing.T) {
	//given
	f := func(i interface{}) int {
		str := fmt.Sprintf("%v", i)
		return schema.HashString(str)
	}
	input := []interface{}{
		map[string]interface{}{
			networkDeviceSchemaNames["UUID"]:                "0452fa68-8246-48b1-a1b2-817fb4baddcb",
			networkDeviceSchemaNames["Name"]:                "device",
			networkDeviceSchemaNames["MetroCode"]:           "SV",
			networkDeviceSchemaNames["HostName"]:            "SV5",
			networkDeviceSchemaNames["LicenseToken"]:        "sWf3df4gaAvbbexw45ga4f",
			networkDeviceSchemaNames["LicenseFile"]:         "/tmp/licenseFile",
			networkDeviceSchemaNames["ACLTemplateUUID"]:     "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
			networkDeviceSchemaNames["AccountNumber"]:       "123456",
			networkDeviceSchemaNames["Notifications"]:       schema.NewSet(schema.HashString, []interface{}{"bla@bla.com"}),
			networkDeviceSchemaNames["AdditionalBandwidth"]: 50,
			networkDeviceSchemaNames["VendorConfiguration"]: map[string]interface{}{
				"key": "value",
			},
			networkDeviceSchemaNames["UserPublicKey"]: schema.NewSet(f, []interface{}{
				map[string]interface{}{
					neDeviceUserKeySchemaNames["Username"]: "user",
					neDeviceUserKeySchemaNames["KeyName"]:  "testKey",
				},
			}),
		},
	}
	expected := &ne.Device{
		UUID:                ne.String(input[0].(map[string]interface{})[networkDeviceSchemaNames["UUID"]].(string)),
		Name:                ne.String(input[0].(map[string]interface{})[networkDeviceSchemaNames["Name"]].(string)),
		MetroCode:           ne.String(input[0].(map[string]interface{})[networkDeviceSchemaNames["MetroCode"]].(string)),
		HostName:            ne.String(input[0].(map[string]interface{})[networkDeviceSchemaNames["HostName"]].(string)),
		LicenseToken:        ne.String(input[0].(map[string]interface{})[networkDeviceSchemaNames["LicenseToken"]].(string)),
		LicenseFile:         ne.String(input[0].(map[string]interface{})[networkDeviceSchemaNames["LicenseFile"]].(string)),
		ACLTemplateUUID:     ne.String(input[0].(map[string]interface{})[networkDeviceSchemaNames["ACLTemplateUUID"]].(string)),
		AccountNumber:       ne.String(input[0].(map[string]interface{})[networkDeviceSchemaNames["AccountNumber"]].(string)),
		Notifications:       expandSetToStringList(input[0].(map[string]interface{})[networkDeviceSchemaNames["Notifications"]].(*schema.Set)),
		AdditionalBandwidth: ne.Int(input[0].(map[string]interface{})[networkDeviceSchemaNames["AdditionalBandwidth"]].(int)),
		VendorConfiguration: map[string]string{
			"key": "value",
		},
		UserPublicKey: expandNetworkDeviceUserKeys(input[0].(map[string]interface{})[networkDeviceSchemaNames["UserPublicKey"]].(*schema.Set))[0],
	}
	//when
	out := expandNetworkDeviceSecondary(input)
	//then
	assert.NotNil(t, out, "Output is not empty")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestNetworkDevice_uploadLicenseFile(t *testing.T) {
	//given
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
	//when
	err := uploadDeviceLicenseFile(openFunc, uploadFunc, ne.StringValue(device.TypeCode), device)
	//then
	assert.Nil(t, err, "License upload function does not return any error")
	assert.Equal(t, licenseFileID, ne.StringValue(device.LicenseFileID), "Device LicenseFileID matches")
	assert.Equal(t, ne.StringValue(device.MetroCode), rxMetroCode, "Received metroCode matches")
	assert.Equal(t, ne.StringValue(device.TypeCode), rxTypeCode, "Received typeCode matches")
	assert.Equal(t, fileName, rxFileName, "Received fileName matches")
	assert.Equal(t, ne.DeviceManagementTypeSelf, rxMgmtMode, "Received management mode matches")
	assert.Equal(t, ne.DeviceLicenseModeBYOL, rxLicMode, "Received management mode matches")
}

func TestNetworkDevice_statusProvisioningWaitConfiguration(t *testing.T) {
	//given
	deviceID := "test"
	var queriedDeviceID string
	fetchFunc := func(uuid string) (*ne.Device, error) {
		queriedDeviceID = uuid
		return &ne.Device{Status: ne.String(ne.DeviceStateProvisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceStatusProvisioningWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	//then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, queriedDeviceID, "Queried device ID matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_statusDeleteWaitConfiguration(t *testing.T) {
	//given
	deviceID := "test"
	var queriedDeviceID string
	fetchFunc := func(uuid string) (*ne.Device, error) {
		queriedDeviceID = uuid
		return &ne.Device{Status: ne.String(ne.DeviceStateDeprovisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceStatusDeleteWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	//then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, queriedDeviceID, "Queried device ID matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_licenseStatusWaitConfiguration(t *testing.T) {
	//given
	deviceID := "test"
	var queriedDeviceID string
	fetchFunc := func(uuid string) (*ne.Device, error) {
		queriedDeviceID = uuid
		return &ne.Device{LicenseStatus: ne.String(ne.DeviceLicenseStateApplied)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceLicenseStatusWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	//then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, queriedDeviceID, "Queried device ID matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_ACLStatusWaitConfiguration(t *testing.T) {
	//given
	aclID := "test"
	var receivedACLID string
	fetchFunc := func(uuid string) (*ne.ACLTemplate, error) {
		receivedACLID = uuid
		return &ne.ACLTemplate{DeviceACLStatus: ne.String(ne.ACLDeviceStatusProvisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceACLStatusWaitConfiguration(fetchFunc, aclID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	//then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, aclID, receivedACLID, "Queried ACL id matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}

func TestNetworkDevice_AdditionalBandwidthStatusWaitConfiguration(t *testing.T) {
	//given
	deviceID := "test"
	var receivedID string
	fetchFunc := func(uuid string) (*ne.DeviceAdditionalBandwidthDetails, error) {
		receivedID = uuid
		return &ne.DeviceAdditionalBandwidthDetails{Status: ne.String(ne.DeviceAdditionalBandwidthStatusProvisioned)}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceAdditionalBandwidthStatusWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForStateContext(context.Background())
	//then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, deviceID, receivedID, "Queried Additional Bandwidth device id matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Additional bandwidth status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Additional bandwidth wait configuration min timeout matches")
}
