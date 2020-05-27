package equinix

import (
	"ecx-go-client/v3"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

var spFieldsToTest = []string{"UUID", "State", "AlertPercentage", "AllowCustomSpeed", "AllowOverSubscription", "APIAvailable", "AuthKeyLabel",
	"ConnectionNameLabel", "CTagLabel", "EnableAutoGenerateServiceKey", "EquinixManagedPortAndVlan", "IntegrationID", "Name",
	"OnBandwidthThresholdNotification", "OnProfileApprovalRejectNotification", "OnVcApprovalRejectionNotification", "OverSubscription",
	"Private", "PrivateUserEmails", "RequiredRedundancy", "SpeedFromAPI", "TagType", "VlanSameAsPrimary"}
var featuresFieldsToTest = []string{"CloudReach", "TestProfile"}
var portFieldsToTest = []string{"ID", "MetroCode"}
var bandFieldsToTest = []string{"Speed", "SpeedUnit"}

func TestECXL2ServiceProfile_resourceDataFromDomain(t *testing.T) {
	//Given
	d := schema.TestResourceDataRaw(t, createECXL2ServiceProfileResourceSchema(), make(map[string]interface{}))
	testProfile := ecx.L2ServiceProfile{
		AlertPercentage:              30.3,
		AllowCustomSpeed:             true,
		AllowOverSubscription:        false,
		APIAvailable:                 true,
		AuthKeyLabel:                 "authKeyLabel",
		ConnectionNameLabel:          "connectionNameLabel",
		CTagLabel:                    "cTagLabel",
		EnableAutoGenerateServiceKey: false,
		EquinixManagedPortAndVlan:    false,
		Features: ecx.L2ServiceProfileFeatures{
			CloudReach:  true,
			TestProfile: true,
		},
		IntegrationID:                       "integrationID",
		Name:                                "name",
		OnBandwidthThresholdNotification:    []string{"miro@equinix.com", "jane@equinix.com"},
		OnProfileApprovalRejectNotification: []string{"miro@equinix.com", "jane@equinix.com"},
		OnVcApprovalRejectionNotification:   []string{"miro@equinix.com", "jane@equinix.com"},
		OverSubscription:                    "2x",
		Ports: []ecx.L2ServiceProfilePort{
			{
				ID:        "port-id1",
				MetroCode: "FR",
			}, {
				ID:        "port-id2",
				MetroCode: "AM",
			},
		},
		Private:            true,
		PrivateUserEmails:  []string{"miro@equinix.com", "jane@equinix.com"},
		RequiredRedundancy: false,
		SpeedBands: []ecx.L2ServiceProfileSpeedBand{
			{
				Speed:     100,
				SpeedUnit: "MB",
			}, {
				Speed:     1000,
				SpeedUnit: "MB",
			},
		},
		SpeedFromAPI:      false,
		TagType:           "tagType",
		VlanSameAsPrimary: false,
	}

	//When
	err := updateECXL2ServiceProfileResource(&testProfile, d)

	//Then
	assert.Nil(t, err, "Schema update should not return an error")
	sourceMatchesTargetSchema(t, testProfile, spFieldsToTest, d, ecxL2ServiceProfileSchemaNames)

	features := d.Get(ecxL2ServiceProfileSchemaNames["Features"])
	featuresList := features.(*schema.Set).List()
	assert.Equal(t, 1, len(featuresList), "There is one features element")
	sourceMatchesTargetSchema(t, testProfile.Features, featuresFieldsToTest, featuresList[0], ecxL2ServiceProfileFeaturesSchemaNames)

	ports := d.Get(ecxL2ServiceProfileSchemaNames["Port"]).(*schema.Set)
	assert.Equal(t, len(testProfile.Ports), ports.Len(), "Number of ports matches")
	for i := range testProfile.Ports {
		assert.True(t, ports.Contains(structToSchemaMap(testProfile.Ports[i], ecxL2ServiceProfilePortSchemaNames)), "Profile port is defined in schema")
	}

	bands := d.Get(ecxL2ServiceProfileSchemaNames["SpeedBand"]).(*schema.Set)
	assert.Equal(t, len(testProfile.SpeedBands), bands.Len(), "Number of speed bands matches")
	for i := range testProfile.SpeedBands {
		assert.True(t, bands.Contains(structToSchemaMap(testProfile.SpeedBands[i], ecxL2ServiceProfileSpeedBandSchemaNames)), "Profile speedband is defined in schema")
	}
}

func TestECXL2ServiceProfile_domainFromResourceData(t *testing.T) {
	//Given
	d := schema.TestResourceDataRaw(t, createECXL2ServiceProfileResourceSchema(), make(map[string]interface{}))
	d.Set(ecxL2ServiceProfileSchemaNames["UUID"], "uuid")
	d.Set(ecxL2ServiceProfileSchemaNames["State"], "state")
	d.Set(ecxL2ServiceProfileSchemaNames["AlertPercentage"], 30.3)
	d.Set(ecxL2ServiceProfileSchemaNames["AllowCustomSpeed"], false)
	d.Set(ecxL2ServiceProfileSchemaNames["AllowOverSubscription"], false)
	d.Set(ecxL2ServiceProfileSchemaNames["APIAvailable"], true)
	d.Set(ecxL2ServiceProfileSchemaNames["AuthKeyLabel"], "authKeyLabel")
	d.Set(ecxL2ServiceProfileSchemaNames["ConnectionNameLabel"], "connectionNameLabel")
	d.Set(ecxL2ServiceProfileSchemaNames["CTagLabel"], "cTagLabel")
	d.Set(ecxL2ServiceProfileSchemaNames["EnableAutoGenerateServiceKey"], false)
	d.Set(ecxL2ServiceProfileSchemaNames["EquinixManagedPortAndVlan"], false)
	d.Set(ecxL2ServiceProfileSchemaNames["IntegrationID"], "integrationID")
	d.Set(ecxL2ServiceProfileSchemaNames["Name"], "name")
	d.Set(ecxL2ServiceProfileSchemaNames["OnBandwidthThresholdNotification"], []string{"janek@equinix.com", "marek@equinix.com"})
	d.Set(ecxL2ServiceProfileSchemaNames["OnProfileApprovalRejectNotification"], []string{"janek@equinix.com", "marek@equinix.com"})
	d.Set(ecxL2ServiceProfileSchemaNames["OnVcApprovalRejectionNotification"], []string{"janek@equinix.com", "marek@equinix.com"})
	d.Set(ecxL2ServiceProfileSchemaNames["OverSubscription"], "2x")
	d.Set(ecxL2ServiceProfileSchemaNames["Private"], true)
	d.Set(ecxL2ServiceProfileSchemaNames["PrivateUserEmails"], []string{"janek@equinix.com", "marek@equinix.com"})
	d.Set(ecxL2ServiceProfileSchemaNames["RequiredRedundancy"], false)
	d.Set(ecxL2ServiceProfileSchemaNames["SpeedFromAPI"], true)
	d.Set(ecxL2ServiceProfileSchemaNames["TagType"], "tagType")
	d.Set(ecxL2ServiceProfileSchemaNames["VlanSameAsPrimary"], true)
	d.Set(ecxL2ServiceProfileSchemaNames["Features"], flattenECXL2ServiceProfileFeatures(
		ecx.L2ServiceProfileFeatures{CloudReach: true, TestProfile: false}))
	d.Set(ecxL2ServiceProfileSchemaNames["Port"], flattenECXL2ServiceProfilePorts([]ecx.L2ServiceProfilePort{
		{ID: "p1", MetroCode: "FR"},
		{ID: "p2", MetroCode: "AM"},
	}))
	d.Set(ecxL2ServiceProfileSchemaNames["SpeedBand"], flattenECXL2ServiceProfileSpeedBands([]ecx.L2ServiceProfileSpeedBand{
		{Speed: 666, SpeedUnit: "MB"},
		{Speed: 1000, SpeedUnit: "MB"},
	}))

	//when
	profile := createECXL2ServiceProfile(d)

	//then
	assert.NotNil(t, profile, "Profile should not be nil")
	sourceMatchesTargetSchema(t, *profile, spFieldsToTest, d, ecxL2ServiceProfileSchemaNames)

	assert.NotNil(t, profile.Features, "Features should not be nil")
	featuresList := d.Get(ecxL2ServiceProfileSchemaNames["Features"]).(*schema.Set).List()
	sourceMatchesTargetSchema(t, profile.Features, featuresFieldsToTest, featuresList[0], ecxL2ServiceProfileFeaturesSchemaNames)

	ports := d.Get(ecxL2ServiceProfileSchemaNames["Port"]).(*schema.Set)
	assert.Equal(t, ports.Len(), len(profile.Ports), "Number of ports matches")
	for i := range profile.Ports {
		assert.True(t, ports.Contains(structToSchemaMap(profile.Ports[i], ecxL2ServiceProfilePortSchemaNames)), "Profile port is defined in schema")
	}

	bands := d.Get(ecxL2ServiceProfileSchemaNames["SpeedBand"]).(*schema.Set)
	assert.Equal(t, bands.Len(), len(profile.SpeedBands), "Number of speedbands matches")
	for i := range profile.SpeedBands {
		assert.True(t, bands.Contains(structToSchemaMap(profile.SpeedBands[i], ecxL2ServiceProfileSpeedBandSchemaNames)), "Speedband is defined in schema")
	}
}
