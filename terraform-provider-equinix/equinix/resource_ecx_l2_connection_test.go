package equinix

import (
	ecx "ecx-go-client/v3"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

var primaryConnFields = []string{"UUID", "Name", "ProfileUUID", "Speed", "SpeedUnit", "Status", "Notifications", "PurchaseOrderNumber",
	"PortUUID", "VlanSTag", "VlanCTag", "ZSidePortUUID", "ZSideVlanSTag", "ZSideVlanCTag", "SellerRegion", "SellerMetroCode",
	"AuthorizationKey", "RedundantUUID"}

var secondaryConnFields = []string{"UUID", "Name", "Status", "PortUUID", "VlanSTag", "VlanCTag",
	"ZSidePortUUID", "ZSideVlanSTag", "ZSideVlanCTag"}

func TestECXL2Connection_resourceDataFromConnections(t *testing.T) {
	//Given
	d := schema.TestResourceDataRaw(t, createECXL2ConnectionResourceSchema(), make(map[string]interface{}))
	primary := ecx.L2Connection{
		UUID:                "uu-id",
		Name:                "name",
		ProfileUUID:         "profileUUID",
		Speed:               666,
		SpeedUnit:           "MB",
		Status:              "DELETED",
		Notifications:       []string{"janek@equinix.com", "marek@equinix.com"},
		PurchaseOrderNumber: "1234",
		PortUUID:            "primaryPortUUID",
		VlanSTag:            100,
		VlanCTag:            101,
		ZSidePortUUID:       "primaryZSidePortUUID",
		ZSideVlanSTag:       200,
		ZSideVlanCTag:       201,
		SellerRegion:        "EMEA",
		SellerMetroCode:     "AM",
		AuthorizationKey:    "authorizationKey",
		RedundantUUID:       "sec-uu-id"}
	secondary := ecx.L2Connection{
		UUID:          "sec-uu-id",
		Name:          "secName",
		Status:        "DELETED",
		PortUUID:      "secondaryPortUUID",
		VlanSTag:      690,
		VlanCTag:      691,
		ZSidePortUUID: "secondaryZSidePortUUID",
		ZSideVlanSTag: 717,
		ZSideVlanCTag: 718}

	//When
	err := updateECXL2ConnectionResource(&primary, &secondary, d)

	//Then
	assert.Nil(t, err, "Schema update should not return an error")
	sourceMatchesTargetSchema(t, primary, primaryConnFields, d, ecxL2ConnectionSchemaNames)

	secConns := d.Get(ecxL2ConnectionSchemaNames["SecondaryConnection"])
	assert.IsType(t, &schema.Set{}, secConns, "Secondary connection schema type is set")
	secConnsList := secConns.(*schema.Set).List()
	assert.Equal(t, 1, len(secConnsList), "There is only one secondary connection")
	sourceMatchesTargetSchema(t, secondary, secondaryConnFields, secConnsList[0], ecxL2ConnectionSchemaNames)
}

func TestECXL2Connection_connectionsFromResourceData(t *testing.T) {
	//Given
	d := schema.TestResourceDataRaw(t, createECXL2ConnectionResourceSchema(), make(map[string]interface{}))
	d.Set(ecxL2ConnectionSchemaNames["UUID"], "uu-id")
	d.Set(ecxL2ConnectionSchemaNames["Name"], "name")
	d.Set(ecxL2ConnectionSchemaNames["ProfileUUID"], "ProfileUUID")
	d.Set(ecxL2ConnectionSchemaNames["Speed"], 100)
	d.Set(ecxL2ConnectionSchemaNames["SpeedUnit"], "SpeedUnit")
	d.Set(ecxL2ConnectionSchemaNames["Status"], "PROVISIONED")
	d.Set(ecxL2ConnectionSchemaNames["Notifications"], []string{"janek@equinix.com", "marek@equinix.com"})
	d.Set(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"], "1234")
	d.Set(ecxL2ConnectionSchemaNames["PortUUID"], "portUUID")
	d.Set(ecxL2ConnectionSchemaNames["VlanSTag"], 100)
	d.Set(ecxL2ConnectionSchemaNames["VlanCTag"], 200)
	d.Set(ecxL2ConnectionSchemaNames["ZSidePortUUID"], "zSidePortUUID")
	d.Set(ecxL2ConnectionSchemaNames["ZSideVlanSTag"], 500)
	d.Set(ecxL2ConnectionSchemaNames["ZSideVlanCTag"], 600)
	d.Set(ecxL2ConnectionSchemaNames["SellerRegion"], "EMEA")
	d.Set(ecxL2ConnectionSchemaNames["SellerMetroCode"], "AM")
	d.Set(ecxL2ConnectionSchemaNames["AuthorizationKey"], "authorizationKey")
	d.Set(ecxL2ConnectionSchemaNames["RedundantUUID"], "sec-uu-id")
	secConn := make(map[string]interface{})
	secConn[ecxL2ConnectionSchemaNames["UUID"]] = "sec-uu-id"
	secConn[ecxL2ConnectionSchemaNames["Name"]] = "sec-name"
	secConn[ecxL2ConnectionSchemaNames["Status"]] = "PROVISIONED"
	secConn[ecxL2ConnectionSchemaNames["PortUUID"]] = "sec-portUUID"
	secConn[ecxL2ConnectionSchemaNames["VlanSTag"]] = 1000
	secConn[ecxL2ConnectionSchemaNames["VlanCTag"]] = 2000
	secConn[ecxL2ConnectionSchemaNames["ZSidePortUUID"]] = "zSidePortUUID"
	secConn[ecxL2ConnectionSchemaNames["ZSideVlanSTag"]] = 5000
	secConn[ecxL2ConnectionSchemaNames["ZSideVlanCTag"]] = 6000
	d.Set(ecxL2ConnectionSchemaNames["SecondaryConnection"], []map[string]interface{}{secConn})

	//when
	primary, secondary := createECXL2Connections(d)

	//then
	assert.NotNil(t, primary, "Primary connection should be present")
	sourceMatchesTargetSchema(t, *primary, primaryConnFields, d, ecxL2ConnectionSchemaNames)
	assert.NotNil(t, secondary, "Secondary connection should be present")
	sourceMatchesTargetSchema(t, *secondary, secondaryConnFields, secConn, ecxL2ConnectionSchemaNames)
}
