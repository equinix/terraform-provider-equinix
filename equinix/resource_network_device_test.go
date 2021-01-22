package equinix

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNetworkDevice_createFromResourceData(t *testing.T) {
	expectedPrimaryUserKey := ne.DeviceUserPublicKey{
		Username: "user",
		KeyName:  "key",
	}
	expectedPrimaryVendorConfig := map[string]string{
		"key": "value",
	}
	expectedPrimary := &ne.Device{
		Name:                "device",
		TypeCode:            "CSR1000V",
		MetroCode:           "SV",
		Throughput:          100,
		ThroughputUnit:      "Mbps",
		HostName:            "test",
		PackageCode:         "SEC",
		Version:             "9.0.1",
		IsBYOL:              true,
		LicenseToken:        "sWf3df4gaAvbbexw45ga4f",
		LicenseFile:         "/tmp/licenseFile",
		ACLTemplateUUID:     "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
		AccountNumber:       "123456",
		Notifications:       []string{"bla@bla.com"},
		PurchaseOrderNumber: "1234567890",
		TermLength:          1,
		AdditionalBandwidth: 50,
		OrderReference:      "12312121sddsf1231",
		InterfaceCount:      10,
		CoreCount:           2,
		IsSelfManaged:       true,
		VendorConfiguration: expectedPrimaryVendorConfig,
		UserPublicKey:       &expectedPrimaryUserKey,
	}
	rawData := map[string]interface{}{
		networkDeviceSchemaNames["Name"]:                expectedPrimary.Name,
		networkDeviceSchemaNames["TypeCode"]:            expectedPrimary.TypeCode,
		networkDeviceSchemaNames["MetroCode"]:           expectedPrimary.MetroCode,
		networkDeviceSchemaNames["Throughput"]:          expectedPrimary.Throughput,
		networkDeviceSchemaNames["ThroughputUnit"]:      expectedPrimary.ThroughputUnit,
		networkDeviceSchemaNames["HostName"]:            expectedPrimary.HostName,
		networkDeviceSchemaNames["PackageCode"]:         expectedPrimary.PackageCode,
		networkDeviceSchemaNames["Version"]:             expectedPrimary.Version,
		networkDeviceSchemaNames["IsBYOL"]:              expectedPrimary.IsBYOL,
		networkDeviceSchemaNames["LicenseToken"]:        expectedPrimary.LicenseToken,
		networkDeviceSchemaNames["LicenseFile"]:         expectedPrimary.LicenseFile,
		networkDeviceSchemaNames["ACLTemplateUUID"]:     expectedPrimary.ACLTemplateUUID,
		networkDeviceSchemaNames["AccountNumber"]:       expectedPrimary.AccountNumber,
		networkDeviceSchemaNames["PurchaseOrderNumber"]: expectedPrimary.PurchaseOrderNumber,
		networkDeviceSchemaNames["TermLength"]:          expectedPrimary.TermLength,
		networkDeviceSchemaNames["AdditionalBandwidth"]: expectedPrimary.AdditionalBandwidth,
		networkDeviceSchemaNames["OrderReference"]:      expectedPrimary.OrderReference,
		networkDeviceSchemaNames["InterfaceCount"]:      expectedPrimary.InterfaceCount,
		networkDeviceSchemaNames["CoreCount"]:           expectedPrimary.CoreCount,
		networkDeviceSchemaNames["IsSelfManaged"]:       expectedPrimary.IsSelfManaged,
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
		Name:                "device",
		TypeCode:            "CSR1000V",
		MetroCode:           "SV",
		Throughput:          100,
		ThroughputUnit:      "Mbps",
		HostName:            "test",
		PackageCode:         "SEC",
		Version:             "9.0.1",
		IsBYOL:              true,
		LicenseToken:        "sWf3df4gaAvbbexw45ga4f",
		ACLTemplateUUID:     "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
		AccountNumber:       "123456",
		Notifications:       []string{"bla@bla.com"},
		PurchaseOrderNumber: "1234567890",
		TermLength:          1,
		AdditionalBandwidth: 50,
		OrderReference:      "12312121sddsf1231",
		InterfaceCount:      10,
		CoreCount:           2,
		IsSelfManaged:       true,
		VendorConfiguration: map[string]string{
			"key": "value",
		},
		UserPublicKey: &ne.DeviceUserPublicKey{
			Username: "user",
			KeyName:  "key",
		},
	}
	inputSecondary := &ne.Device{}
	secondarySchemaLicenseFile := "/tmp/licenseFileSec"
	d := schema.TestResourceDataRaw(t, createNetworkDeviceSchema(), make(map[string]interface{}))
	d.Set(networkDeviceSchemaNames["Secondary"], flattenNetworkDeviceSecondary(&ne.Device{
		LicenseFile: secondarySchemaLicenseFile,
	}))
	//when
	err := updateNetworkDeviceResource(inputPrimary, inputSecondary, d)

	//then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, inputPrimary.Name, d.Get(networkDeviceSchemaNames["Name"]), "Name matches")
	assert.Equal(t, inputPrimary.TypeCode, d.Get(networkDeviceSchemaNames["TypeCode"]), "TypeCode matches")
	assert.Equal(t, inputPrimary.MetroCode, d.Get(networkDeviceSchemaNames["MetroCode"]), "MetroCode matches")
	assert.Equal(t, inputPrimary.Throughput, d.Get(networkDeviceSchemaNames["Throughput"]), "Throughput matches")
	assert.Equal(t, inputPrimary.ThroughputUnit, d.Get(networkDeviceSchemaNames["ThroughputUnit"]), "ThroughputUnit matches")
	assert.Equal(t, inputPrimary.HostName, d.Get(networkDeviceSchemaNames["HostName"]), "HostName matches")
	assert.Equal(t, inputPrimary.PackageCode, d.Get(networkDeviceSchemaNames["PackageCode"]), "PackageCode matches")
	assert.Equal(t, inputPrimary.Version, d.Get(networkDeviceSchemaNames["Version"]), "Version matches")
	assert.Equal(t, inputPrimary.IsBYOL, d.Get(networkDeviceSchemaNames["IsBYOL"]), "IsBYOL matches")
	assert.Equal(t, inputPrimary.LicenseToken, d.Get(networkDeviceSchemaNames["LicenseToken"]), "LicenseToken matches")
	assert.Equal(t, inputPrimary.ACLTemplateUUID, d.Get(networkDeviceSchemaNames["ACLTemplateUUID"]), "ACLTemplateUUID matches")
	assert.Equal(t, inputPrimary.AccountNumber, d.Get(networkDeviceSchemaNames["AccountNumber"]), "AccountNumber matches")
	assert.Equal(t, inputPrimary.Notifications, expandSetToStringList(d.Get(networkDeviceSchemaNames["Notifications"]).(*schema.Set)), "Notifications matches")
	assert.Equal(t, inputPrimary.PurchaseOrderNumber, d.Get(networkDeviceSchemaNames["PurchaseOrderNumber"]), "PurchaseOrderNumber matches")
	assert.Equal(t, inputPrimary.TermLength, d.Get(networkDeviceSchemaNames["TermLength"]), "TermLength matches")
	assert.Equal(t, inputPrimary.AdditionalBandwidth, d.Get(networkDeviceSchemaNames["AdditionalBandwidth"]), "AdditionalBandwidth matches")
	assert.Equal(t, inputPrimary.OrderReference, d.Get(networkDeviceSchemaNames["OrderReference"]), "OrderReference matches")
	assert.Equal(t, inputPrimary.InterfaceCount, d.Get(networkDeviceSchemaNames["InterfaceCount"]), "InterfaceCount matches")
	assert.Equal(t, inputPrimary.CoreCount, d.Get(networkDeviceSchemaNames["CoreCount"]), "CoreCount matches")
	assert.Equal(t, inputPrimary.IsSelfManaged, d.Get(networkDeviceSchemaNames["IsSelfManaged"]), "IsSelfManaged matches")
	assert.Equal(t, inputPrimary.VendorConfiguration, expandInterfaceMapToStringMap(d.Get(networkDeviceSchemaNames["VendorConfiguration"]).(map[string]interface{})), "VendorConfiguration matches")
	assert.Equal(t, inputPrimary.UserPublicKey, expandNetworkDeviceUserKeys(d.Get(networkDeviceSchemaNames["UserPublicKey"]).(*schema.Set))[0], "UserPublicKey matches")
	assert.Equal(t, secondarySchemaLicenseFile, expandNetworkDeviceSecondary(d.Get(networkDeviceSchemaNames["Secondary"]).([]interface{})).LicenseFile, "Secondary LicenseFile matches")
}

func TestNetworkDevice_flattenSecondary(t *testing.T) {
	//given
	input := &ne.Device{
		UUID:                "0452fa68-8246-48b1-a1b2-817fb4baddcb",
		Name:                "device",
		Status:              ne.DeviceStateProvisioned,
		LicenseStatus:       ne.DeviceLicenseStateApplied,
		MetroCode:           "SV",
		IBX:                 "SV5",
		Region:              "AMER",
		HostName:            "test",
		LicenseToken:        "sWf3df4gaAvbbexw45ga4f",
		LicenseFileID:       "d72dbe58-e596-4698-8b57-0a38e8077d25",
		LicenseFile:         "/tmp/myfile",
		ACLTemplateUUID:     "a624178c-6d59-4798-9a7f-2ddf2c7c5881",
		SSHIPAddress:        "1.1.1.1",
		SSHIPFqdn:           "test-1.1.1.1-SV.test.equinix.com",
		AccountNumber:       "123456",
		Notifications:       []string{"bla@bla.com"},
		RedundancyType:      "PRIMARY",
		RedundantUUID:       "c2a147a3-ff47-4a24-a6e5-d6d7ce6459f3",
		AdditionalBandwidth: 50,
		Interfaces: []ne.DeviceInterface{
			{
				ID:                1,
				Name:              "GigabitEthernet1",
				Status:            "AVAILABLE",
				OperationalStatus: "UP",
				MACAddress:        "58-0A-C9-7A-DA-E9",
				IPAddress:         "2.2.2.2",
				AssignedType:      "test-connection(AWS Direct Connect)",
				Type:              "DATA",
			},
		},
		VendorConfiguration: map[string]string{
			"key": "value",
		},
		UserPublicKey: &ne.DeviceUserPublicKey{
			Username: "user",
			KeyName:  "testKey",
		},
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
		UUID:                input[0].(map[string]interface{})[networkDeviceSchemaNames["UUID"]].(string),
		Name:                input[0].(map[string]interface{})[networkDeviceSchemaNames["Name"]].(string),
		MetroCode:           input[0].(map[string]interface{})[networkDeviceSchemaNames["MetroCode"]].(string),
		HostName:            input[0].(map[string]interface{})[networkDeviceSchemaNames["HostName"]].(string),
		LicenseToken:        input[0].(map[string]interface{})[networkDeviceSchemaNames["LicenseToken"]].(string),
		LicenseFile:         input[0].(map[string]interface{})[networkDeviceSchemaNames["LicenseFile"]].(string),
		ACLTemplateUUID:     input[0].(map[string]interface{})[networkDeviceSchemaNames["ACLTemplateUUID"]].(string),
		AccountNumber:       input[0].(map[string]interface{})[networkDeviceSchemaNames["AccountNumber"]].(string),
		Notifications:       expandSetToStringList(input[0].(map[string]interface{})[networkDeviceSchemaNames["Notifications"]].(*schema.Set)),
		AdditionalBandwidth: input[0].(map[string]interface{})[networkDeviceSchemaNames["AdditionalBandwidth"]].(int),
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
	device := &ne.Device{LicenseFile: "/path/to/" + fileName, MetroCode: "SV", TypeCode: "VSRX"}
	var rxMetroCode, rxFileName, rxTypeCode, rxMgmtMode, rxLicMode string
	uploadFunc := func(metroCode, deviceTypeCode, deviceManagementMode, licenseMode, fileName string, reader io.Reader) (string, error) {
		rxMetroCode = metroCode
		rxFileName = fileName
		rxTypeCode = deviceTypeCode
		rxMgmtMode = deviceManagementMode
		rxLicMode = licenseMode
		return licenseFileID, nil
	}
	openFunc := func(name string) (*os.File, error) {
		return &os.File{}, nil
	}
	//when
	err := uploadDeviceLicenseFile(openFunc, uploadFunc, device.TypeCode, device)
	//then
	assert.Nil(t, err, "License upload function does not return any error")
	assert.Equal(t, licenseFileID, device.LicenseFileID, "Device LicenseFileID matches")
	assert.Equal(t, device.MetroCode, rxMetroCode, "Received metroCode matches")
	assert.Equal(t, device.TypeCode, rxTypeCode, "Received typeCode matches")
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
		return &ne.Device{Status: ne.DeviceStateProvisioned}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceStatusProvisioningWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForState()
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
		return &ne.Device{Status: ne.DeviceStateDeprovisioned}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceStatusDeleteWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForState()
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
		return &ne.Device{LicenseStatus: ne.DeviceLicenseStateApplied}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceLicenseStatusWaitConfiguration(fetchFunc, deviceID, delay, timeout)
	_, err := waitConfig.WaitForState()
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
		return &ne.ACLTemplate{DeviceACLStatus: ne.ACLDeviceStatusProvisioned}, nil
	}
	delay := 100 * time.Millisecond
	timeout := 10 * time.Minute
	//when
	waitConfig := createNetworkDeviceACLStatusWaitConfiguration(fetchFunc, aclID, delay, timeout)
	_, err := waitConfig.WaitForState()
	//then
	assert.Nil(t, err, "WaitForState does not return an error")
	assert.Equal(t, aclID, receivedACLID, "Queried ACL id matches")
	assert.Equal(t, timeout, waitConfig.Timeout, "Device status wait configuration timeout matches")
	assert.Equal(t, delay, waitConfig.MinTimeout, "Device status wait configuration min timeout matches")
}
