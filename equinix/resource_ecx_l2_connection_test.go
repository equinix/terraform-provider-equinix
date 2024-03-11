package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/ecx-go/v2"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/stretchr/testify/assert"
)

func TestFabricL2Connection_createFromResourceData(t *testing.T) {
	rawData := map[string]interface{}{
		ecxL2ConnectionSchemaNames["Name"]:                "kekewrmMwe",
		ecxL2ConnectionSchemaNames["ProfileUUID"]:         "5d113752-996b-4b59-8e21-8927e7b98058",
		ecxL2ConnectionSchemaNames["Speed"]:               50,
		ecxL2ConnectionSchemaNames["SpeedUnit"]:           "MB",
		ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]: "234242353",
		ecxL2ConnectionSchemaNames["PortUUID"]:            "52c00d7f-c310-458e-9426-1d7549e1f600",
		ecxL2ConnectionSchemaNames["DeviceUUID"]:          "5f1483f4-c479-424d-98c5-43a266aae25c",
		ecxL2ConnectionSchemaNames["DeviceInterfaceID"]:   5,
		ecxL2ConnectionSchemaNames["VlanSTag"]:            1043,
		ecxL2ConnectionSchemaNames["VlanCTag"]:            2045,
		ecxL2ConnectionSchemaNames["NamedTag"]:            "Public",
		ecxL2ConnectionSchemaNames["ZSidePortUUID"]:       "52c00d7f-c310-458e-9426-1d7549e1f600",
		ecxL2ConnectionSchemaNames["ZSideVlanSTag"]:       420,
		ecxL2ConnectionSchemaNames["ZSideVlanCTag"]:       1056,
		ecxL2ConnectionSchemaNames["SellerRegion"]:        "werwerewr",
		ecxL2ConnectionSchemaNames["SellerMetroCode"]:     "SV",
		ecxL2ConnectionSchemaNames["AuthorizationKey"]:    "123456789012",
		ecxL2ConnectionSchemaNames["ServiceToken"]:        "1c356a7b-d632-18a5-c357-a33146cab65d",
		ecxL2ConnectionSchemaNames["ZSideServiceToken"]:   "febc9d80-11e0-4dc8-8eb8-c41b6b378df2",
	}
	d := schema.TestResourceDataRaw(t, createECXL2ConnectionResourceSchema(), rawData)
	d.Set(ecxL2ConnectionSchemaNames["Notifications"], []string{"test@test.com"})
	expectedPrimary := &ecx.L2Connection{
		Name:                ecx.String(rawData[ecxL2ConnectionSchemaNames["Name"]].(string)),
		ProfileUUID:         ecx.String(rawData[ecxL2ConnectionSchemaNames["ProfileUUID"]].(string)),
		Speed:               ecx.Int(rawData[ecxL2ConnectionSchemaNames["Speed"]].(int)),
		SpeedUnit:           ecx.String(rawData[ecxL2ConnectionSchemaNames["SpeedUnit"]].(string)),
		Notifications:       []string{"test@test.com"},
		PurchaseOrderNumber: ecx.String(rawData[ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]].(string)),
		PortUUID:            ecx.String(rawData[ecxL2ConnectionSchemaNames["PortUUID"]].(string)),
		DeviceUUID:          ecx.String(rawData[ecxL2ConnectionSchemaNames["DeviceUUID"]].(string)),
		DeviceInterfaceID:   ecx.Int(rawData[ecxL2ConnectionSchemaNames["DeviceInterfaceID"]].(int)),
		VlanSTag:            ecx.Int(rawData[ecxL2ConnectionSchemaNames["VlanSTag"]].(int)),
		VlanCTag:            ecx.Int(rawData[ecxL2ConnectionSchemaNames["VlanCTag"]].(int)),
		NamedTag:            ecx.String(rawData[ecxL2ConnectionSchemaNames["NamedTag"]].(string)),
		ZSidePortUUID:       ecx.String(rawData[ecxL2ConnectionSchemaNames["ZSidePortUUID"]].(string)),
		ZSideVlanSTag:       ecx.Int(rawData[ecxL2ConnectionSchemaNames["ZSideVlanSTag"]].(int)),
		ZSideVlanCTag:       ecx.Int(rawData[ecxL2ConnectionSchemaNames["ZSideVlanCTag"]].(int)),
		SellerRegion:        ecx.String(rawData[ecxL2ConnectionSchemaNames["SellerRegion"]].(string)),
		SellerMetroCode:     ecx.String(rawData[ecxL2ConnectionSchemaNames["SellerMetroCode"]].(string)),
		AuthorizationKey:    ecx.String(rawData[ecxL2ConnectionSchemaNames["AuthorizationKey"]].(string)),
		ServiceToken:        ecx.String(rawData[ecxL2ConnectionSchemaNames["ServiceToken"]].(string)),
		ZSideServiceToken:   ecx.String(rawData[ecxL2ConnectionSchemaNames["ZSideServiceToken"]].(string)),
	}

	// when
	primary, secondary := createECXL2Connections(d)

	// then
	assert.NotNil(t, primary, "Primary connection is not nil")
	assert.Nil(t, secondary, "Secondary connection is nil")
	assert.Equal(t, expectedPrimary, primary, "Primary connection matches expected result")
}

func TestFabricL2Connection_updateResourceData(t *testing.T) {
	// given
	d := schema.TestResourceDataRaw(t, createECXL2ConnectionResourceSchema(), make(map[string]interface{}))
	input := &ecx.L2Connection{
		UUID:                ecx.String(acctest.RandString(36)),
		Name:                ecx.String(acctest.RandString(36)),
		ProfileUUID:         ecx.String(acctest.RandString(36)),
		Speed:               ecx.Int(50),
		SpeedUnit:           ecx.String("MB"),
		Status:              ecx.String(ecx.ConnectionStatusProvisioned),
		ProviderStatus:      ecx.String(ecx.ConnectionStatusProvisioned),
		Notifications:       []string{"bla@bla.com"},
		PurchaseOrderNumber: ecx.String(acctest.RandString(10)),
		PortUUID:            ecx.String(acctest.RandString(36)),
		DeviceUUID:          ecx.String(acctest.RandString(36)),
		VendorToken:         ecx.String(acctest.RandString(36)),
		VlanSTag:            ecx.Int(acctest.RandIntRange(0, 2000)),
		VlanCTag:            ecx.Int(acctest.RandIntRange(0, 2000)),
		NamedTag:            ecx.String(acctest.RandString(100)),
		AdditionalInfo:      []ecx.L2ConnectionAdditionalInfo{{Name: ecx.String(acctest.RandString(10)), Value: ecx.String(acctest.RandString(10))}},
		ZSidePortUUID:       ecx.String(acctest.RandString(36)),
		ZSideVlanCTag:       ecx.Int(acctest.RandIntRange(0, 2000)),
		ZSideVlanSTag:       ecx.Int(acctest.RandIntRange(0, 2000)),
		SellerRegion:        ecx.String(acctest.RandString(10)),
		SellerMetroCode:     ecx.String(acctest.RandString(2)),
		AuthorizationKey:    ecx.String(acctest.RandString(10)),
		RedundancyGroup:     ecx.String(acctest.RandString(36)),
		RedundancyType:      ecx.String(acctest.RandString(10)),
	}
	prevServiceToken := ecx.String(acctest.RandString(20))
	d.Set(ecxL2ConnectionSchemaNames["ServiceToken"], prevServiceToken)
	prevZsideServiceToken := ecx.String(acctest.RandString(20))
	d.Set(ecxL2ConnectionSchemaNames["ZSideServiceToken"], prevZsideServiceToken)

	// when
	err := updateECXL2ConnectionResource(input, nil, d)

	// then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ecx.StringValue(input.UUID), d.Get(ecxL2ConnectionSchemaNames["UUID"]), "UUID matches")
	assert.Equal(t, ecx.StringValue(input.Name), d.Get(ecxL2ConnectionSchemaNames["Name"]), "Name matches")
	assert.Equal(t, ecx.StringValue(input.ProfileUUID), d.Get(ecxL2ConnectionSchemaNames["ProfileUUID"]), "ProfileUUID matches")
	assert.Equal(t, ecx.IntValue(input.Speed), d.Get(ecxL2ConnectionSchemaNames["Speed"]), "Speed matches")
	assert.Equal(t, ecx.StringValue(input.SpeedUnit), d.Get(ecxL2ConnectionSchemaNames["SpeedUnit"]), "SpeedUnit matches")
	assert.Equal(t, ecx.StringValue(input.Status), d.Get(ecxL2ConnectionSchemaNames["Status"]), "Status matches")
	assert.Equal(t, ecx.StringValue(input.ProviderStatus), d.Get(ecxL2ConnectionSchemaNames["ProviderStatus"]), "ProviderStatus matches")
	assert.Equal(t, input.Notifications, converters.SetToStringList(d.Get(ecxL2ConnectionSchemaNames["Notifications"]).(*schema.Set)), "Notifications matches")
	assert.Equal(t, ecx.StringValue(input.PurchaseOrderNumber), d.Get(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]), "PurchaseOrderNumber matches")
	assert.Equal(t, ecx.StringValue(input.PortUUID), d.Get(ecxL2ConnectionSchemaNames["PortUUID"]), "PortUUID matches")
	assert.Equal(t, ecx.StringValue(input.DeviceUUID), d.Get(ecxL2ConnectionSchemaNames["DeviceUUID"]), "DeviceUUID matches")
	assert.Equal(t, ecx.IntValue(input.DeviceInterfaceID), d.Get(ecxL2ConnectionSchemaNames["DeviceInterfaceID"]), "DeviceInterfaceID matches")
	assert.Equal(t, ecx.StringValue(input.VendorToken), d.Get(ecxL2ConnectionSchemaNames["VendorToken"]), "VendorToken matches")
	assert.Equal(t, ecx.StringValue(input.VendorToken), d.Get(ecxL2ConnectionSchemaNames["ServiceToken"]), "ServiceToken and VendorToken match")
	assert.NotEqual(t, ecx.StringValue(prevServiceToken), d.Get(ecxL2ConnectionSchemaNames["ServiceToken"]), "ServiceToken updated")
	assert.Equal(t, ecx.StringValue(input.VendorToken), d.Get(ecxL2ConnectionSchemaNames["ZSideServiceToken"]), "ZSideServiceToken and VendorToken match")
	assert.NotEqual(t, ecx.StringValue(prevZsideServiceToken), d.Get(ecxL2ConnectionSchemaNames["ZSideServiceToken"]), "ZSideServiceToken updated")
	assert.Equal(t, ecx.IntValue(input.VlanSTag), d.Get(ecxL2ConnectionSchemaNames["VlanSTag"]), "VlanSTag matches")
	assert.Equal(t, ecx.IntValue(input.VlanCTag), d.Get(ecxL2ConnectionSchemaNames["VlanCTag"]), "VlanCTag matches")
	assert.Equal(t, ecx.StringValue(input.NamedTag), d.Get(ecxL2ConnectionSchemaNames["NamedTag"]), "NamedTag matches")
	assert.Equal(t, input.AdditionalInfo, expandECXL2ConnectionAdditionalInfo(d.Get(ecxL2ConnectionSchemaNames["AdditionalInfo"]).(*schema.Set)), "AdditionalInfo matches")
	assert.Equal(t, ecx.StringValue(input.ZSidePortUUID), d.Get(ecxL2ConnectionSchemaNames["ZSidePortUUID"]), "ZSidePortUUID matches")
	assert.Equal(t, ecx.IntValue(input.ZSideVlanCTag), d.Get(ecxL2ConnectionSchemaNames["ZSideVlanCTag"]), "ZSideVlanCTag matches")
	assert.Equal(t, ecx.IntValue(input.ZSideVlanSTag), d.Get(ecxL2ConnectionSchemaNames["ZSideVlanSTag"]), "ZSideVlanSTag matches")
	assert.Equal(t, ecx.StringValue(input.SellerRegion), d.Get(ecxL2ConnectionSchemaNames["SellerRegion"]), "SellerRegion matches")
	assert.Equal(t, ecx.StringValue(input.SellerMetroCode), d.Get(ecxL2ConnectionSchemaNames["SellerMetroCode"]), "SellerMetroCode matches")
	assert.Equal(t, ecx.StringValue(input.AuthorizationKey), d.Get(ecxL2ConnectionSchemaNames["AuthorizationKey"]), "AuthorizationKey matches")
	assert.Equal(t, ecx.StringValue(input.RedundancyGroup), d.Get(ecxL2ConnectionSchemaNames["RedundancyGroup"]), "RedundancyGroup matches")
	assert.Equal(t, ecx.StringValue(input.RedundancyType), d.Get(ecxL2ConnectionSchemaNames["RedundancyType"]), "RedundancyType matches")
}

func TestFabricL2Connection_flattenSecondary(t *testing.T) {
	// given
	input := &ecx.L2Connection{
		UUID:             ecx.String(acctest.RandString(36)),
		Name:             ecx.String(acctest.RandString(36)),
		ProfileUUID:      ecx.String(acctest.RandString(36)),
		Speed:            ecx.Int(50),
		SpeedUnit:        ecx.String("MB"),
		Status:           ecx.String(ecx.ConnectionStatusProvisioned),
		ProviderStatus:   ecx.String(ecx.ConnectionStatusProvisioned),
		PortUUID:         ecx.String(acctest.RandString(36)),
		DeviceUUID:       ecx.String(acctest.RandString(36)),
		VlanSTag:         ecx.Int(acctest.RandIntRange(0, 2000)),
		VlanCTag:         ecx.Int(acctest.RandIntRange(0, 2000)),
		ZSidePortUUID:    ecx.String(acctest.RandString(36)),
		ZSideVlanCTag:    ecx.Int(acctest.RandIntRange(0, 2000)),
		ZSideVlanSTag:    ecx.Int(acctest.RandIntRange(0, 2000)),
		SellerRegion:     ecx.String(acctest.RandString(10)),
		SellerMetroCode:  ecx.String(acctest.RandString(2)),
		AuthorizationKey: ecx.String(acctest.RandString(10)),
		RedundancyGroup:  ecx.String(acctest.RandString(10)),
		RedundancyType:   ecx.String(acctest.RandString(10)),
		VendorToken:      ecx.String(acctest.RandString(36)),
	}
	previousInput := &ecx.L2Connection{
		DeviceInterfaceID: ecx.Int(acctest.RandIntRange(0, 10)),
		ServiceToken:      ecx.String(acctest.RandString(36)),
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
			ecxL2ConnectionSchemaNames["DeviceInterfaceID"]: previousInput.DeviceInterfaceID,
			ecxL2ConnectionSchemaNames["VlanSTag"]:          input.VlanSTag,
			ecxL2ConnectionSchemaNames["VlanCTag"]:          input.VlanCTag,
			ecxL2ConnectionSchemaNames["ZSidePortUUID"]:     input.ZSidePortUUID,
			ecxL2ConnectionSchemaNames["ZSideVlanCTag"]:     input.ZSideVlanCTag,
			ecxL2ConnectionSchemaNames["ZSideVlanSTag"]:     input.ZSideVlanSTag,
			ecxL2ConnectionSchemaNames["SellerRegion"]:      input.SellerRegion,
			ecxL2ConnectionSchemaNames["SellerMetroCode"]:   input.SellerMetroCode,
			ecxL2ConnectionSchemaNames["AuthorizationKey"]:  input.AuthorizationKey,
			ecxL2ConnectionSchemaNames["RedundancyGroup"]:   input.RedundancyGroup,
			ecxL2ConnectionSchemaNames["RedundancyType"]:    input.RedundancyType,
			ecxL2ConnectionSchemaNames["Actions"]:           []interface{}{},
			ecxL2ConnectionSchemaNames["ServiceToken"]:      input.VendorToken,
			ecxL2ConnectionSchemaNames["VendorToken"]:       input.VendorToken,
		},
	}

	// when
	out := flattenECXL2ConnectionSecondary(previousInput, input)

	// then
	assert.NotNil(t, out, "Output is not nil")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestFabricL2Connection_expandSecondary(t *testing.T) {
	// given
	input := []interface{}{
		map[string]interface{}{
			ecxL2ConnectionSchemaNames["Name"]:              "testName",
			ecxL2ConnectionSchemaNames["ProfileUUID"]:       "529574df-1dfb-4fad-b904-8edd3920e8b7",
			ecxL2ConnectionSchemaNames["Speed"]:             50,
			ecxL2ConnectionSchemaNames["SpeedUnit"]:         "MB",
			ecxL2ConnectionSchemaNames["PortUUID"]:          "8640622d-e4fd-4118-8e0e-566fc5af8f6a",
			ecxL2ConnectionSchemaNames["DeviceUUID"]:        "af93a177-5f3d-4102-b231-c15fc49ca099",
			ecxL2ConnectionSchemaNames["DeviceInterfaceID"]: 6,
			ecxL2ConnectionSchemaNames["VlanSTag"]:          434,
			ecxL2ConnectionSchemaNames["VlanCTag"]:          0,
			ecxL2ConnectionSchemaNames["SellerRegion"]:      "",
			ecxL2ConnectionSchemaNames["SellerMetroCode"]:   "SV",
			ecxL2ConnectionSchemaNames["AuthorizationKey"]:  "123456789012",
			ecxL2ConnectionSchemaNames["ServiceToken"]:      "5cc17b61-ac44-4ed4-b035-71a8d46e449f",
		},
	}
	expected := &ecx.L2Connection{
		Name:              ecx.String(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["Name"]].(string)),
		ProfileUUID:       ecx.String(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["ProfileUUID"]].(string)),
		Speed:             ecx.Int(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["Speed"]].(int)),
		SpeedUnit:         ecx.String(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["SpeedUnit"]].(string)),
		PortUUID:          ecx.String(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["PortUUID"]].(string)),
		DeviceUUID:        ecx.String(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["DeviceUUID"]].(string)),
		DeviceInterfaceID: ecx.Int(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["DeviceInterfaceID"]].(int)),
		VlanSTag:          ecx.Int(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["VlanSTag"]].(int)),
		VlanCTag:          nil,
		SellerRegion:      nil,
		SellerMetroCode:   ecx.String(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["SellerMetroCode"]].(string)),
		AuthorizationKey:  ecx.String(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["AuthorizationKey"]].(string)),
		ServiceToken:      ecx.String(input[0].(map[string]interface{})[ecxL2ConnectionSchemaNames["ServiceToken"]].(string)),
	}

	// when
	out := expandECXL2ConnectionSecondary(input)

	// then
	assert.NotNil(t, out, "Output is not empty")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestFabricL2Connection_flattenAdditionalInfo(t *testing.T) {
	// given
	input := []ecx.L2ConnectionAdditionalInfo{
		{
			Name:  ecx.String(acctest.RandString(32)),
			Value: ecx.String(acctest.RandString(32)),
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			ecxL2ConnectionAdditionalInfoSchemaNames["Name"]:  input[0].Name,
			ecxL2ConnectionAdditionalInfoSchemaNames["Value"]: input[0].Value,
		},
	}
	// when
	out := flattenECXL2ConnectionAdditionalInfo(input)
	// then
	assert.NotNil(t, out, "Output is not empty")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestFabricL2Connection_expandAdditionalInfo(t *testing.T) {
	f := func(i interface{}) int {
		str := fmt.Sprintf("%v", i)
		return schema.HashString(str)
	}
	// given
	input := schema.NewSet(f, []interface{}{
		map[string]interface{}{
			ecxL2ConnectionAdditionalInfoSchemaNames["Name"]:  acctest.RandString(36),
			ecxL2ConnectionAdditionalInfoSchemaNames["Value"]: acctest.RandString(36),
		},
	})
	inputList := input.List()
	expected := []ecx.L2ConnectionAdditionalInfo{
		{
			Name:  ecx.String(inputList[0].(map[string]interface{})[ecxL2ConnectionAdditionalInfoSchemaNames["Name"]].(string)),
			Value: ecx.String(inputList[0].(map[string]interface{})[ecxL2ConnectionAdditionalInfoSchemaNames["Value"]].(string)),
		},
	}
	// when
	out := expandECXL2ConnectionAdditionalInfo(input)
	// then
	assert.NotNil(t, out, "Output is not empty")
	assert.Equal(t, expected, out, "Output matches expected result")
}

func TestFabricL2Connection_flattenActions(t *testing.T) {
	// given
	input := []ecx.L2ConnectionAction{
		{
			Type:        ecx.String(acctest.RandString(32)),
			OperationID: ecx.String(acctest.RandString(32)),
			Message:     ecx.String(acctest.RandString(32)),
			RequiredData: []ecx.L2ConnectionActionData{
				{
					Key:               ecx.String(acctest.RandString(10)),
					Label:             ecx.String(acctest.RandString(10)),
					Value:             ecx.String(acctest.RandString(10)),
					IsEditable:        ecx.Bool(true),
					ValidationPattern: ecx.String(acctest.RandString(10)),
				},
			},
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			ecxL2ConnectionActionsSchemaNames["Type"]:        input[0].Type,
			ecxL2ConnectionActionsSchemaNames["OperationID"]: input[0].OperationID,
			ecxL2ConnectionActionsSchemaNames["Message"]:     input[0].Message,
			ecxL2ConnectionActionsSchemaNames["RequiredData"]: []interface{}{
				map[string]interface{}{
					ecxL2ConnectionActionDataSchemaNames["Key"]:               input[0].RequiredData[0].Key,
					ecxL2ConnectionActionDataSchemaNames["Label"]:             input[0].RequiredData[0].Label,
					ecxL2ConnectionActionDataSchemaNames["Value"]:             input[0].RequiredData[0].Value,
					ecxL2ConnectionActionDataSchemaNames["IsEditable"]:        input[0].RequiredData[0].IsEditable,
					ecxL2ConnectionActionDataSchemaNames["ValidationPattern"]: input[0].RequiredData[0].ValidationPattern,
				},
			},
		},
	}
	// when
	out := flattenECXL2ConnectionActions(input)
	// then
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
	// given
	updateReq := mockedL2ConnectionUpdateRequest{}
	changes := map[string]interface{}{
		ecxL2ConnectionSchemaNames["Name"]:      acctest.RandString(32),
		ecxL2ConnectionSchemaNames["Speed"]:     50,
		ecxL2ConnectionSchemaNames["SpeedUnit"]: "MB",
	}
	// when
	fillFabricL2ConnectionUpdateRequest(&updateReq, changes)
	// then
	assert.Equal(t, changes[ecxL2ConnectionSchemaNames["Name"]], updateReq.name, "Update request name matches")
	assert.Equal(t, changes[ecxL2ConnectionSchemaNames["Speed"]], updateReq.speed, "Update request speed matches")
	assert.Equal(t, changes[ecxL2ConnectionSchemaNames["SpeedUnit"]], updateReq.speedUnit, "Update speed unit matches")
}
