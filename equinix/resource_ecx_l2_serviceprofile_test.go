package equinix

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal"

	"github.com/equinix/ecx-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestFabricL2ServiceProfile_createFromResourceData(t *testing.T) {
	rawData := map[string]interface{}{
		ecxL2ServiceProfileSchemaNames["UUID"]:                         "5d113752-996b-4b59-8e21-8927e7b98058",
		ecxL2ServiceProfileSchemaNames["State"]:                        "active",
		ecxL2ServiceProfileSchemaNames["AlertPercentage"]:              50.30,
		ecxL2ServiceProfileSchemaNames["AllowCustomSpeed"]:             true,
		ecxL2ServiceProfileSchemaNames["AllowOverSubscription"]:        true,
		ecxL2ServiceProfileSchemaNames["APIAvailable"]:                 true,
		ecxL2ServiceProfileSchemaNames["AuthKeyLabel"]:                 "testAuthKeyLabel",
		ecxL2ServiceProfileSchemaNames["ConnectionNameLabel"]:          "testConnectionLabel",
		ecxL2ServiceProfileSchemaNames["CTagLabel"]:                    "testCTagLabel",
		ecxL2ServiceProfileSchemaNames["Description"]:                  "testDescription",
		ecxL2ServiceProfileSchemaNames["EnableAutoGenerateServiceKey"]: true,
		ecxL2ServiceProfileSchemaNames["EquinixManagedPortAndVlan"]:    true,
		ecxL2ServiceProfileSchemaNames["IntegrationID"]:                "testIntegrationID",
		ecxL2ServiceProfileSchemaNames["Name"]:                         "testName",
		ecxL2ServiceProfileSchemaNames["OverSubscription"]:             "2x",
		ecxL2ServiceProfileSchemaNames["Private"]:                      true,
		ecxL2ServiceProfileSchemaNames["RequiredRedundancy"]:           true,
		ecxL2ServiceProfileSchemaNames["SpeedFromAPI"]:                 true,
		ecxL2ServiceProfileSchemaNames["TagType"]:                      "CTAGED",
		ecxL2ServiceProfileSchemaNames["VlanSameAsPrimary"]:            true,
	}
	d := schema.TestResourceDataRaw(t, createECXL2ServiceProfileResourceSchema(), rawData)
	testEmails := []string{"marry@equinix.com"}
	d.Set(ecxL2ServiceProfileSchemaNames["OnBandwidthThresholdNotification"], testEmails)
	d.Set(ecxL2ServiceProfileSchemaNames["OnProfileApprovalRejectNotification"], testEmails)
	d.Set(ecxL2ServiceProfileSchemaNames["OnVcApprovalRejectionNotification"], testEmails)
	d.Set(ecxL2ServiceProfileSchemaNames["PrivateUserEmails"], testEmails)
	d.Set(ecxL2ServiceProfileSchemaNames["Features"], flattenECXL2ServiceProfileFeatures(
		&ecx.L2ServiceProfileFeatures{CloudReach: ecx.Bool(true), TestProfile: ecx.Bool(true)},
		ecx.L2ServiceProfileFeatures{CloudReach: ecx.Bool(true), TestProfile: ecx.Bool(false)},
	))
	d.Set(ecxL2ServiceProfileSchemaNames["Port"], flattenECXL2ServiceProfilePorts([]ecx.L2ServiceProfilePort{
		{
			ID:        ecx.String("testPortIDOne"),
			MetroCode: ecx.String("SV"),
		},
		{
			ID:        ecx.String("testPortIDTwo"),
			MetroCode: ecx.String("DC"),
		},
	}))
	d.Set(ecxL2ServiceProfileSchemaNames["SpeedBand"], flattenECXL2ServiceProfileSpeedBands([]ecx.L2ServiceProfileSpeedBand{
		{
			Speed:     ecx.Int(100),
			SpeedUnit: ecx.String("MB"),
		},
		{
			Speed:     ecx.Int(1),
			SpeedUnit: ecx.String("GB"),
		},
	}))
	expected := ecx.L2ServiceProfile{
		UUID:                                ecx.String(rawData[ecxL2ServiceProfileSchemaNames["UUID"]].(string)),
		State:                               ecx.String(rawData[ecxL2ServiceProfileSchemaNames["State"]].(string)),
		AlertPercentage:                     ecx.Float64(rawData[ecxL2ServiceProfileSchemaNames["AlertPercentage"]].(float64)),
		AllowCustomSpeed:                    ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["AllowCustomSpeed"]].(bool)),
		AllowOverSubscription:               ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["AllowOverSubscription"]].(bool)),
		APIAvailable:                        ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["APIAvailable"]].(bool)),
		AuthKeyLabel:                        ecx.String(rawData[ecxL2ServiceProfileSchemaNames["AuthKeyLabel"]].(string)),
		ConnectionNameLabel:                 ecx.String(rawData[ecxL2ServiceProfileSchemaNames["ConnectionNameLabel"]].(string)),
		CTagLabel:                           ecx.String(rawData[ecxL2ServiceProfileSchemaNames["CTagLabel"]].(string)),
		Description:                         ecx.String(rawData[ecxL2ServiceProfileSchemaNames["Description"]].(string)),
		EnableAutoGenerateServiceKey:        ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["EnableAutoGenerateServiceKey"]].(bool)),
		EquinixManagedPortAndVlan:           ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["EquinixManagedPortAndVlan"]].(bool)),
		IntegrationID:                       ecx.String(rawData[ecxL2ServiceProfileSchemaNames["IntegrationID"]].(string)),
		Name:                                ecx.String(rawData[ecxL2ServiceProfileSchemaNames["Name"]].(string)),
		OverSubscription:                    ecx.String(rawData[ecxL2ServiceProfileSchemaNames["OverSubscription"]].(string)),
		Private:                             ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["Private"]].(bool)),
		RequiredRedundancy:                  ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["RequiredRedundancy"]].(bool)),
		SpeedFromAPI:                        ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["SpeedFromAPI"]].(bool)),
		TagType:                             ecx.String(rawData[ecxL2ServiceProfileSchemaNames["TagType"]].(string)),
		VlanSameAsPrimary:                   ecx.Bool(rawData[ecxL2ServiceProfileSchemaNames["VlanSameAsPrimary"]].(bool)),
		OnBandwidthThresholdNotification:    testEmails,
		OnProfileApprovalRejectNotification: testEmails,
		OnVcApprovalRejectionNotification:   testEmails,
		PrivateUserEmails:                   testEmails,
		Features:                            expandECXL2ServiceProfileFeatures(d.Get(ecxL2ServiceProfileSchemaNames["Features"]).(*schema.Set).List()),
		Ports:                               expandECXL2ServiceProfilePorts(d.Get(ecxL2ServiceProfileSchemaNames["Port"]).(*schema.Set)),
		SpeedBands:                          expandECXL2ServiceProfileSpeedBands(d.Get(ecxL2ServiceProfileSchemaNames["SpeedBand"]).(*schema.Set)),
	}
	// when
	result := createECXL2ServiceProfile(d)
	// then
	assert.NotNil(t, result, "Result is not nil")
	assert.Equal(t, expected, *result, "Result matches expected value")
}

func TestFabricL2ServiceProfile_updateResourceData(t *testing.T) {
	// given
	d := schema.TestResourceDataRaw(t, createECXL2ServiceProfileResourceSchema(), make(map[string]interface{}))
	input := ecx.L2ServiceProfile{
		UUID:                                ecx.String("5d113752-996b-4b59-8e21-8927e7b98058"),
		State:                               ecx.String("active"),
		AlertPercentage:                     ecx.Float64(43.6),
		AllowCustomSpeed:                    ecx.Bool(true),
		AllowOverSubscription:               ecx.Bool(true),
		APIAvailable:                        ecx.Bool(true),
		AuthKeyLabel:                        ecx.String("testAuthKeyLabel"),
		ConnectionNameLabel:                 ecx.String("testConnectionLabel"),
		CTagLabel:                           ecx.String("CTAGED"),
		EnableAutoGenerateServiceKey:        ecx.Bool(true),
		EquinixManagedPortAndVlan:           ecx.Bool(true),
		IntegrationID:                       ecx.String("testIntegrationID"),
		Name:                                ecx.String("testName"),
		OnBandwidthThresholdNotification:    []string{"marry@equinix.com"},
		OnProfileApprovalRejectNotification: []string{"marry@equinix.com"},
		OnVcApprovalRejectionNotification:   []string{"marry@equinix.com"},
		OverSubscription:                    ecx.String("2x"),
		Private:                             ecx.Bool(true),
		PrivateUserEmails:                   []string{"marry@equinix.com"},
		RequiredRedundancy:                  ecx.Bool(true),
		SpeedFromAPI:                        ecx.Bool(true),
		TagType:                             ecx.String("bla"),
		VlanSameAsPrimary:                   ecx.Bool(true),
		Description:                         ecx.String("testDescription"),
		Features: ecx.L2ServiceProfileFeatures{
			CloudReach:  ecx.Bool(true),
			TestProfile: ecx.Bool(true),
		},
		Ports: []ecx.L2ServiceProfilePort{
			{
				ID:        ecx.String("one"),
				MetroCode: ecx.String("DC"),
			},
			{
				ID:        ecx.String("one"),
				MetroCode: ecx.String("SV"),
			},
		},
		SpeedBands: []ecx.L2ServiceProfileSpeedBand{
			{
				Speed:     ecx.Int(100),
				SpeedUnit: ecx.String("MB"),
			},
			{
				Speed:     ecx.Int(1),
				SpeedUnit: ecx.String("GB"),
			},
		},
	}
	// when
	err := updateECXL2ServiceProfileResource(&input, d)
	// then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ecx.StringValue(input.UUID), d.Get(ecxL2ServiceProfileSchemaNames["UUID"]), "UUID matches")
	assert.Equal(t, ecx.StringValue(input.State), d.Get(ecxL2ServiceProfileSchemaNames["State"]), "State matches")
	assert.Equal(t, ecx.Float64Value(input.AlertPercentage), d.Get(ecxL2ServiceProfileSchemaNames["AlertPercentage"]), "AlertPercentage matches")
	assert.Equal(t, ecx.BoolValue(input.AllowCustomSpeed), d.Get(ecxL2ServiceProfileSchemaNames["AllowCustomSpeed"]), "AllowCustomSpeed matches")
	assert.Equal(t, ecx.BoolValue(input.AllowOverSubscription), d.Get(ecxL2ServiceProfileSchemaNames["AllowOverSubscription"]), "AllowOverSubscription matches")
	assert.Equal(t, ecx.BoolValue(input.APIAvailable), d.Get(ecxL2ServiceProfileSchemaNames["APIAvailable"]), "APIAvailable matches")
	assert.Equal(t, ecx.StringValue(input.AuthKeyLabel), d.Get(ecxL2ServiceProfileSchemaNames["AuthKeyLabel"]), "AuthKeyLabel matches")
	assert.Equal(t, ecx.StringValue(input.ConnectionNameLabel), d.Get(ecxL2ServiceProfileSchemaNames["ConnectionNameLabel"]), "ConnectionNameLabel matches")
	assert.Equal(t, ecx.StringValue(input.CTagLabel), d.Get(ecxL2ServiceProfileSchemaNames["CTagLabel"]), "CTagLabel matches")
	assert.Equal(t, ecx.BoolValue(input.EnableAutoGenerateServiceKey), d.Get(ecxL2ServiceProfileSchemaNames["EnableAutoGenerateServiceKey"]), "EnableAutoGenerateServiceKey matches")
	assert.Equal(t, ecx.BoolValue(input.EquinixManagedPortAndVlan), d.Get(ecxL2ServiceProfileSchemaNames["EquinixManagedPortAndVlan"]), "EquinixManagedPortAndVlan matches")
	assert.Equal(t, ecx.StringValue(input.IntegrationID), d.Get(ecxL2ServiceProfileSchemaNames["IntegrationID"]), "IntegrationID matches")
	assert.Equal(t, ecx.StringValue(input.Name), d.Get(ecxL2ServiceProfileSchemaNames["Name"]), "Name matches")
	assert.Equal(t, input.OnBandwidthThresholdNotification, internal.ExpandSetToStringList(d.Get(ecxL2ServiceProfileSchemaNames["OnBandwidthThresholdNotification"]).(*schema.Set)), "OnBandwidthThresholdNotification matches")
	assert.Equal(t, input.OnProfileApprovalRejectNotification, internal.ExpandSetToStringList(d.Get(ecxL2ServiceProfileSchemaNames["OnProfileApprovalRejectNotification"]).(*schema.Set)), "OnProfileApprovalRejectNotification matches")
	assert.Equal(t, input.OnVcApprovalRejectionNotification, internal.ExpandSetToStringList(d.Get(ecxL2ServiceProfileSchemaNames["OnVcApprovalRejectionNotification"]).(*schema.Set)), "OnVcApprovalRejectionNotification matches")
	assert.Equal(t, ecx.StringValue(input.OverSubscription), d.Get(ecxL2ServiceProfileSchemaNames["OverSubscription"]), "OverSubscription matches")
	assert.Equal(t, ecx.BoolValue(input.Private), d.Get(ecxL2ServiceProfileSchemaNames["Private"]), "Private matches")
	assert.Equal(t, input.PrivateUserEmails, internal.ExpandSetToStringList(d.Get(ecxL2ServiceProfileSchemaNames["PrivateUserEmails"]).(*schema.Set)), "PrivateUserEmails matches")
	assert.Equal(t, ecx.BoolValue(input.RequiredRedundancy), d.Get(ecxL2ServiceProfileSchemaNames["RequiredRedundancy"]), "RequiredRedundancy matches")
	assert.Equal(t, ecx.BoolValue(input.SpeedFromAPI), d.Get(ecxL2ServiceProfileSchemaNames["SpeedFromAPI"]), "SpeedFromAPI matches")
	assert.Equal(t, ecx.StringValue(input.TagType), d.Get(ecxL2ServiceProfileSchemaNames["TagType"]), "TagType matches")
	assert.Equal(t, ecx.BoolValue(input.VlanSameAsPrimary), d.Get(ecxL2ServiceProfileSchemaNames["VlanSameAsPrimary"]), "VlanSameAsPrimary matches")
	assert.Equal(t, ecx.StringValue(input.Description), d.Get(ecxL2ServiceProfileSchemaNames["Description"]), "Description matches")
	assert.Equal(t, input.Features, expandECXL2ServiceProfileFeatures(d.Get(ecxL2ServiceProfileSchemaNames["Features"]).(*schema.Set).List()), "Features matches")
	assert.Equal(t, input.Ports, expandECXL2ServiceProfilePorts(d.Get(ecxL2ServiceProfileSchemaNames["Port"]).(*schema.Set)), "Ports matches")
	assert.Equal(t, input.SpeedBands, expandECXL2ServiceProfileSpeedBands(d.Get(ecxL2ServiceProfileSchemaNames["SpeedBand"]).(*schema.Set)), "SpeedBand matches")
}
