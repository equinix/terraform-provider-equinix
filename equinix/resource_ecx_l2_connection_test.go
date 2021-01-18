package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/ecx-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestFabricL2Connection_createFromResourceData(t *testing.T) {
	rawData := map[string]interface{}{
		ecxL2ConnectionSchemaNames["Name"]:                randString(36),
		ecxL2ConnectionSchemaNames["ProfileUUID"]:         randString(36),
		ecxL2ConnectionSchemaNames["Speed"]:               50,
		ecxL2ConnectionSchemaNames["SpeedUnit"]:           "MB",
		ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]: randString(36),
		ecxL2ConnectionSchemaNames["PortUUID"]:            randString(36),
		ecxL2ConnectionSchemaNames["DeviceUUID"]:          randString(36),
		ecxL2ConnectionSchemaNames["DeviceInterfaceID"]:   randInt(10),
		ecxL2ConnectionSchemaNames["VlanSTag"]:            randInt(2000),
		ecxL2ConnectionSchemaNames["VlanCTag"]:            randInt(2000),
		ecxL2ConnectionSchemaNames["NamedTag"]:            randString(36),
		ecxL2ConnectionSchemaNames["ZSidePortUUID"]:       randString(36),
		ecxL2ConnectionSchemaNames["ZSideVlanSTag"]:       randInt(2000),
		ecxL2ConnectionSchemaNames["ZSideVlanCTag"]:       randInt(2000),
		ecxL2ConnectionSchemaNames["SellerRegion"]:        randString(10),
		ecxL2ConnectionSchemaNames["SellerMetroCode"]:     randString(2),
		ecxL2ConnectionSchemaNames["AuthorizationKey"]:    randString(10),
	}
	d := schema.TestResourceDataRaw(t, createECXL2ConnectionResourceSchema(), rawData)
	d.Set(ecxL2ConnectionSchemaNames["Notifications"], []string{"test@test.com"})
	expectedPrimary := &ecx.L2Connection{
		Name:                rawData[ecxL2ConnectionSchemaNames["Name"]].(string),
		ProfileUUID:         rawData[ecxL2ConnectionSchemaNames["ProfileUUID"]].(string),
		Speed:               rawData[ecxL2ConnectionSchemaNames["Speed"]].(int),
		SpeedUnit:           rawData[ecxL2ConnectionSchemaNames["SpeedUnit"]].(string),
		Notifications:       []string{"test@test.com"},
		PurchaseOrderNumber: rawData[ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]].(string),
		PortUUID:            rawData[ecxL2ConnectionSchemaNames["PortUUID"]].(string),
		DeviceUUID:          rawData[ecxL2ConnectionSchemaNames["DeviceUUID"]].(string),
		DeviceInterfaceID:   rawData[ecxL2ConnectionSchemaNames["DeviceInterfaceID"]].(int),
		VlanSTag:            rawData[ecxL2ConnectionSchemaNames["VlanSTag"]].(int),
		VlanCTag:            rawData[ecxL2ConnectionSchemaNames["VlanCTag"]].(int),
		NamedTag:            rawData[ecxL2ConnectionSchemaNames["NamedTag"]].(string),
		ZSidePortUUID:       rawData[ecxL2ConnectionSchemaNames["ZSidePortUUID"]].(string),
		ZSideVlanSTag:       rawData[ecxL2ConnectionSchemaNames["ZSideVlanSTag"]].(int),
		ZSideVlanCTag:       rawData[ecxL2ConnectionSchemaNames["ZSideVlanCTag"]].(int),
		SellerRegion:        rawData[ecxL2ConnectionSchemaNames["SellerRegion"]].(string),
		SellerMetroCode:     rawData[ecxL2ConnectionSchemaNames["SellerMetroCode"]].(string),
		AuthorizationKey:    rawData[ecxL2ConnectionSchemaNames["AuthorizationKey"]].(string),
	}

	//when
	primary, secondary := createECXL2Connections(d)

	//then
	assert.NotNil(t, primary, "Primary connection is not nil")
	assert.Nil(t, secondary, "Secondary connection is nil")
	assert.Equal(t, expectedPrimary, primary, "Primary connection matches expected result")
}

func TestFabricL2Connection_updateResourceData(t *testing.T) {
	//given
	d := schema.TestResourceDataRaw(t, createECXL2ConnectionResourceSchema(), make(map[string]interface{}))
	input := &ecx.L2Connection{
		UUID:                randString(36),
		Name:                randString(36),
		ProfileUUID:         randString(36),
		Speed:               50,
		SpeedUnit:           "MB",
		Status:              ecx.ConnectionStatusProvisioned,
		ProviderStatus:      ecx.ConnectionStatusProvisioned,
		Notifications:       []string{"bla@bla.com"},
		PurchaseOrderNumber: randString(10),
		PortUUID:            randString(36),
		DeviceUUID:          randString(36),
		DeviceInterfaceID:   randInt(10),
		VlanSTag:            randInt(2000),
		VlanCTag:            randInt(2000),
		NamedTag:            randString(100),
		AdditionalInfo:      []ecx.L2ConnectionAdditionalInfo{{Name: randString(10), Value: randString(10)}},
		ZSidePortUUID:       randString(36),
		ZSideVlanCTag:       randInt(2000),
		ZSideVlanSTag:       randInt(2000),
		SellerRegion:        randString(10),
		SellerMetroCode:     randString(2),
		AuthorizationKey:    randString(10),
		RedundantUUID:       randString(36),
		RedundancyType:      randString(10),
	}
	//when
	err := updateECXL2ConnectionResource(input, nil, d)

	//then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, input.UUID, d.Get(ecxL2ConnectionSchemaNames["UUID"]), "UUID matches")
	assert.Equal(t, input.Name, d.Get(ecxL2ConnectionSchemaNames["Name"]), "Name matches")
	assert.Equal(t, input.ProfileUUID, d.Get(ecxL2ConnectionSchemaNames["ProfileUUID"]), "ProfileUUID matches")
	assert.Equal(t, input.Speed, d.Get(ecxL2ConnectionSchemaNames["Speed"]), "Speed matches")
	assert.Equal(t, input.SpeedUnit, d.Get(ecxL2ConnectionSchemaNames["SpeedUnit"]), "SpeedUnit matches")
	assert.Equal(t, input.Status, d.Get(ecxL2ConnectionSchemaNames["Status"]), "Status matches")
	assert.Equal(t, input.ProviderStatus, d.Get(ecxL2ConnectionSchemaNames["ProviderStatus"]), "ProviderStatus matches")
	assert.Equal(t, input.Notifications, expandSetToStringList(d.Get(ecxL2ConnectionSchemaNames["Notifications"]).(*schema.Set)), "Notifications matches")
	assert.Equal(t, input.PurchaseOrderNumber, d.Get(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]), "PurchaseOrderNumber matches")
	assert.Equal(t, input.PortUUID, d.Get(ecxL2ConnectionSchemaNames["PortUUID"]), "PortUUID matches")
	assert.Equal(t, input.DeviceUUID, d.Get(ecxL2ConnectionSchemaNames["DeviceUUID"]), "DeviceUUID matches")
	assert.Equal(t, input.DeviceInterfaceID, d.Get(ecxL2ConnectionSchemaNames["DeviceInterfaceID"]), "DeviceInterfaceID matches")
	assert.Equal(t, input.VlanSTag, d.Get(ecxL2ConnectionSchemaNames["VlanSTag"]), "VlanSTag matches")
	assert.Equal(t, input.VlanCTag, d.Get(ecxL2ConnectionSchemaNames["VlanCTag"]), "VlanCTag matches")
	assert.Equal(t, input.NamedTag, d.Get(ecxL2ConnectionSchemaNames["NamedTag"]), "NamedTag matches")
	assert.Equal(t, input.AdditionalInfo, expandECXL2ConnectionAdditionalInfo(d.Get(ecxL2ConnectionSchemaNames["AdditionalInfo"]).(*schema.Set)), "AdditionalInfo matches")
	assert.Equal(t, input.ZSidePortUUID, d.Get(ecxL2ConnectionSchemaNames["ZSidePortUUID"]), "ZSidePortUUID matches")
	assert.Equal(t, input.ZSideVlanCTag, d.Get(ecxL2ConnectionSchemaNames["ZSideVlanCTag"]), "ZSideVlanCTag matches")
	assert.Equal(t, input.ZSideVlanSTag, d.Get(ecxL2ConnectionSchemaNames["ZSideVlanSTag"]), "ZSideVlanSTag matches")
	assert.Equal(t, input.SellerRegion, d.Get(ecxL2ConnectionSchemaNames["SellerRegion"]), "SellerRegion matches")
	assert.Equal(t, input.SellerMetroCode, d.Get(ecxL2ConnectionSchemaNames["SellerMetroCode"]), "SellerMetroCode matches")
	assert.Equal(t, input.AuthorizationKey, d.Get(ecxL2ConnectionSchemaNames["AuthorizationKey"]), "AuthorizationKey matches")
	assert.Equal(t, input.RedundantUUID, d.Get(ecxL2ConnectionSchemaNames["RedundantUUID"]), "RedundantUUID matches")
	assert.Equal(t, input.RedundancyType, d.Get(ecxL2ConnectionSchemaNames["RedundancyType"]), "RedundancyType matches")
}

func TestFabricL2Connection_flattenSecondary(t *testing.T) {
	//given
	input := &ecx.L2Connection{
		UUID:              randString(36),
		Name:              randString(36),
		ProfileUUID:       randString(36),
		Speed:             50,
		SpeedUnit:         "MB",
		Status:            ecx.ConnectionStatusProvisioned,
		ProviderStatus:    ecx.ConnectionStatusProvisioned,
		PortUUID:          randString(36),
		DeviceUUID:        randString(36),
		DeviceInterfaceID: randInt(10),
		VlanSTag:          randInt(2000),
		VlanCTag:          randInt(2000),
		ZSidePortUUID:     randString(36),
		ZSideVlanCTag:     randInt(2000),
		ZSideVlanSTag:     randInt(2000),
		SellerRegion:      randString(10),
		SellerMetroCode:   randString(2),
		AuthorizationKey:  randString(10),
		RedundantUUID:     randString(36),
		RedundancyType:    randString(10),
	}
	expected := []interface{}{
		map[string]interface{}{
			ecxL2ConnectionSchemaNames["UUID"]:              input.UUID,
			ecxL2ConnectionSchemaNames["Name"]:              input.Name,
			ecxL2ConnectionSchemaNames["ProfileUUID"]:       input.ProfileUUID,
			ecxL2ConnectionSchemaNames["Speed"]:             input.Speed,
			ecxL2ConnectionSchemaNames["SpeedUnit"]:         input.SpeedUnit,
			ecxL2ConnectionSchemaNames["Status"]:            input.Status,
			ecxL2ConnectionSchemaNames["ProviderStatus"]:    input.ProviderStatus,
			ecxL2ConnectionSchemaNames["PortUUID"]:          input.PortUUID,
			ecxL2ConnectionSchemaNames["DeviceUUID"]:        input.DeviceUUID,
			ecxL2ConnectionSchemaNames["DeviceInterfaceID"]: input.DeviceInterfaceID,
			ecxL2ConnectionSchemaNames["VlanSTag"]:          input.VlanSTag,
			ecxL2ConnectionSchemaNames["VlanCTag"]:          input.VlanCTag,
			ecxL2ConnectionSchemaNames["ZSidePortUUID"]:     input.ZSidePortUUID,
			ecxL2ConnectionSchemaNames["ZSideVlanCTag"]:     input.ZSideVlanCTag,
			ecxL2ConnectionSchemaNames["ZSideVlanSTag"]:     input.ZSideVlanSTag,
			ecxL2ConnectionSchemaNames["SellerRegion"]:      input.SellerRegion,
			ecxL2ConnectionSchemaNames["SellerMetroCode"]:   input.SellerMetroCode,
			ecxL2ConnectionSchemaNames["AuthorizationKey"]:  input.AuthorizationKey,
			ecxL2ConnectionSchemaNames["RedundantUUID"]:     input.RedundantUUID,
			ecxL2ConnectionSchemaNames["RedundancyType"]:    input.RedundancyType,
		},
	}

	//when
	out := flattenECXL2ConnectionSecondary(input)

	//then
	assert.NotNil(t, out, "Output is not nil")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestFabricL2Connection_expandSecondary(t *testing.T) {
	//given
	input := []interface{}{
		map[string]interface{}{
			ecxL2ConnectionSchemaNames["Name"]:              randString(36),
			ecxL2ConnectionSchemaNames["ProfileUUID"]:       randString(36),
			ecxL2ConnectionSchemaNames["Speed"]:             50,
			ecxL2ConnectionSchemaNames["SpeedUnit"]:         "MB",
			ecxL2ConnectionSchemaNames["PortUUID"]:          randString(36),
			ecxL2ConnectionSchemaNames["DeviceUUID"]:        randString(36),
			ecxL2ConnectionSchemaNames["DeviceInterfaceID"]: randInt(10),
			ecxL2ConnectionSchemaNames["VlanSTag"]:          randInt(2000),
			ecxL2ConnectionSchemaNames["VlanCTag"]:          randInt(2000),
			ecxL2ConnectionSchemaNames["SellerRegion"]:      randString(10),
			ecxL2ConnectionSchemaNames["SellerMetroCode"]:   randString(2),
			ecxL2ConnectionSchemaNames["AuthorizationKey"]:  randString(10),
		},
	}
	expected := &ecx.L2Connection{
		Name:              input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["Name"]].(string),
		ProfileUUID:       input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["ProfileUUID"]].(string),
		Speed:             input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["Speed"]].(int),
		SpeedUnit:         input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["SpeedUnit"]].(string),
		PortUUID:          input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["PortUUID"]].(string),
		DeviceUUID:        input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["DeviceUUID"]].(string),
		DeviceInterfaceID: input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["DeviceInterfaceID"]].(int),
		VlanSTag:          input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["VlanSTag"]].(int),
		VlanCTag:          input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["VlanCTag"]].(int),
		SellerRegion:      input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["SellerRegion"]].(string),
		SellerMetroCode:   input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["SellerMetroCode"]].(string),
		AuthorizationKey:  input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["AuthorizationKey"]].(string),
	}

	//when
	out := expandECXL2ConnectionSecondary(input)

	//then
	assert.NotNil(t, out, "Output is not empty")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestFabricL2Connection_flattenAdditionalInfo(t *testing.T) {
	//given
	input := []ecx.L2ConnectionAdditionalInfo{
		{
			Name:  randString(32),
			Value: randString(32),
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			ecxL2ConnectionAdditionalInfoSchemaNames["Name"]:  input[0].Name,
			ecxL2ConnectionAdditionalInfoSchemaNames["Value"]: input[0].Value,
		},
	}
	//when
	out := flattenECXL2ConnectionAdditionalInfo(input)
	//then
	assert.NotNil(t, out, "Output is not empty")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestFabricL2Connection_expandAdditionalInfo(t *testing.T) {
	f := func(i interface{}) int {
		str := fmt.Sprintf("%v", i)
		return schema.HashString(str)
	}
	//given
	input := schema.NewSet(f, []interface{}{
		map[string]interface{}{
			ecxL2ConnectionAdditionalInfoSchemaNames["Name"]:  randString(36),
			ecxL2ConnectionAdditionalInfoSchemaNames["Value"]: randString(36),
		},
	})
	inputList := input.List()
	expected := []ecx.L2ConnectionAdditionalInfo{
		{
			Name:  inputList[0].(map[string]interface{})[ecxL2ConnectionAdditionalInfoSchemaNames["Name"]].(string),
			Value: inputList[0].(map[string]interface{})[ecxL2ConnectionAdditionalInfoSchemaNames["Value"]].(string),
		},
	}
	//when
	out := expandECXL2ConnectionAdditionalInfo(input)
	//then
	assert.NotNil(t, out, "Output is not empty")
	assert.Equal(t, expected, out, "Output matches expected result")
}

type mockedL2ConnectionUpdateRequest struct {
	name      string
	speed     int
	speedUnit string
}

func (m *mockedL2ConnectionUpdateRequest) WithName(name string) ecx.L2ConnectionUpdateRequest {
	m.name = name
	return m
}

func (m *mockedL2ConnectionUpdateRequest) WithBandwidth(speed int, speedUnit string) ecx.L2ConnectionUpdateRequest {
	m.speed = speed
	m.speedUnit = speedUnit
	return m
}

func (m *mockedL2ConnectionUpdateRequest) WithSpeed(speed int) ecx.L2ConnectionUpdateRequest {
	m.speed = speed
	return m
}

func (m *mockedL2ConnectionUpdateRequest) WithSpeedUnit(speedUnit string) ecx.L2ConnectionUpdateRequest {
	m.speedUnit = speedUnit
	return m
}

func (m *mockedL2ConnectionUpdateRequest) Execute() error {
	return nil
}

func TestFabricL2Connection_fillUpdateRequest(t *testing.T) {
	//given
	updateReq := mockedL2ConnectionUpdateRequest{}
	changes := map[string]interface{}{
		ecxL2ConnectionSchemaNames["Name"]:      randString(32),
		ecxL2ConnectionSchemaNames["Speed"]:     50,
		ecxL2ConnectionSchemaNames["SpeedUnit"]: "MB",
	}
	//when
	fillFabricL2ConnectionUpdateRequest(&updateReq, changes)
	//then
	assert.Equal(t, changes[ecxL2ConnectionSchemaNames["Name"]], updateReq.name, "Update request name matches")
	assert.Equal(t, changes[ecxL2ConnectionSchemaNames["Speed"]], updateReq.speed, "Update request speed matches")
	assert.Equal(t, changes[ecxL2ConnectionSchemaNames["SpeedUnit"]], updateReq.speedUnit, "Update speed unit matches")
}
