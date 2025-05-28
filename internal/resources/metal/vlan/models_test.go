package vlan

import (
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceModel_Parse(t *testing.T) {
	vlan := &packngo.VirtualNetwork{
		ID:          "vlan-123",
		VXLAN:       1001,
		Description: "Test VLAN",
		Project:     &packngo.Project{ID: "project-abc"},
		Facility: &packngo.Facility{
			Code: "SV",
			Metro: &packngo.Metro{
				Code: "sv",
			},
		},
		Instances: []*packngo.Device{
			{ID: "device-1"},
			{ID: "device-2"},
		},
	}

	var model DataSourceModel
	diags := model.parse(vlan)

	assert.False(t, diags.HasError())
	assert.Equal(t, "vlan-123", model.ID.ValueString())
	assert.Equal(t, int64(1001), model.Vxlan.ValueInt64())
	assert.Equal(t, "project-abc", model.ProjectID.ValueString())
	assert.Equal(t, "sv", model.Metro.ValueString())
	assert.Equal(t, "sv", model.Facility.ValueString())
	assert.Equal(t, "Test VLAN", model.Description.ValueString())
	assert.Equal(t, 2, len(model.AssignedDevicesIDs.Elements()))
}

func TestDataSourceModel_Parse_MetroMismatch(t *testing.T) {
	vlan := &packngo.VirtualNetwork{
		ID:      "vlan-456",
		VXLAN:   2002,
		Project: &packngo.Project{ID: "project-def"},
		Facility: &packngo.Facility{
			Code: "NY",
			Metro: &packngo.Metro{
				Code: "ny",
			},
		},
		Metro:     &packngo.Metro{Code: "la"}, // Mismatch
		Instances: []*packngo.Device{},
	}

	var model DataSourceModel
	diags := model.parse(vlan)

	fmt.Printf("diags is %v\n", diags)

	assert.Equal(t, diags.WarningsCount(), 1)
	assert.Contains(t, diags.Warnings()[0].Summary(), "unexpected value for metro")
}

func TestDataSourceModel_Parse_MissingFacilityAndMetro(t *testing.T) {
	vlan := &packngo.VirtualNetwork{
		ID:          "vlan-789",
		VXLAN:       3003,
		Description: "No location info",
		Project:     &packngo.Project{ID: "project-xyz"},
		Instances:   []*packngo.Device{},
	}

	var model DataSourceModel
	diags := model.parse(vlan)
	fmt.Printf("model is %v\n", model)
	assert.False(t, diags.HasError())
	assert.Equal(t, "", model.Facility.ValueString())
	assert.Equal(t, "", model.Metro.ValueString())
}

func TestResourceModel_Parse(t *testing.T) {
	vlan := &packngo.VirtualNetwork{
		ID:          "vlan-456",
		VXLAN:       2002,
		Description: "Another VLAN",
		Project:     &packngo.Project{ID: "project-def"},
		Facility: &packngo.Facility{
			Code: "NY",
			Metro: &packngo.Metro{
				Code: "ny",
			},
		},
		Metro: &packngo.Metro{Code: "ny"},
	}

	var model ResourceModel
	diags := model.parse(vlan)

	assert.False(t, diags.HasError())
	assert.Equal(t, "vlan-456", model.ID.ValueString())
	assert.Equal(t, int64(2002), model.Vxlan.ValueInt64())
	assert.Equal(t, "project-def", model.ProjectID.ValueString())
	assert.Equal(t, "ny", model.Metro.ValueString())
	assert.Equal(t, "ny", model.Facility.ValueString())
	assert.Equal(t, "Another VLAN", model.Description.ValueString())
}

func TestResourceModel_Parse_MetroMismatch(t *testing.T) {
	vlan := &packngo.VirtualNetwork{
		ID:          "vlan-456",
		VXLAN:       2002,
		Description: "Another VLAN",
		Project:     &packngo.Project{ID: "project-def"},
		Metro:       &packngo.Metro{Code: "ny"},
	}

	model := ResourceModel{
		Metro: types.StringValue("la"), // Intentionally different
	}

	diags := model.parse(vlan)

	assert.True(t, diags.HasError())
	assert.Contains(t, diags.Errors()[0].Summary(), "unexpected value for metro")
	assert.Equal(t, "vlan-456", model.ID.ValueString())
	assert.Equal(t, int64(2002), model.Vxlan.ValueInt64())
	assert.Equal(t, "project-def", model.ProjectID.ValueString())
	assert.Equal(t, "la", model.Metro.ValueString()) // Remains unchanged
	assert.Equal(t, "Another VLAN", model.Description.ValueString())
}

func TestResourceModel_ParseMetalV1_Positive(t *testing.T) {
	vlanID := "vlan-789"
	var vxlan int32 = 3003
	description := "Metal VLAN"
	projectID := "project-uuid"
	metroCode := "la"

	vlan := &metalv1.VirtualNetwork{
		Id:          &vlanID,
		Vxlan:       &vxlan,
		Description: &description,
		AssignedTo:  &metalv1.Project{Id: &projectID},
		Facility: &metalv1.Href{
			Href:                 "/metal/v1/facility/facility-uuid",
			AdditionalProperties: map[string]interface{}{"facility_code": "LA1"},
		},
		MetroCode: &metroCode,
	}

	var model ResourceModel
	diags := model.parseMetalV1(vlan)

	assert.False(t, diags.HasError())
	assert.Equal(t, "vlan-789", model.ID.ValueString())
	assert.Equal(t, int64(3003), model.Vxlan.ValueInt64())
	assert.Equal(t, "project-uuid", model.ProjectID.ValueString())
	assert.Equal(t, "la1", model.Facility.ValueString())
	assert.Equal(t, "la", model.Metro.ValueString())
	assert.Equal(t, "Metal VLAN", model.Description.ValueString())
}

func TestResourceModel_ParseMetalV1_MetroMismatch(t *testing.T) {
	// Metro mismatch to trigger a diagnostic error
	vlanID := "vlan-999"
	var vxlan int32 = 4004
	description := "Mismatched Metro VLAN"
	projectID := "project-xyz"
	metroCode := "la"

	vlan := &metalv1.VirtualNetwork{
		Id:          &vlanID,
		Vxlan:       &vxlan,
		Description: &description,
		AssignedTo:  &metalv1.Project{Id: &projectID},
		Facility: &metalv1.Href{
			AdditionalProperties: map[string]interface{}{"facility_code": "LA1"},
		},
		MetroCode: &metroCode, // Mismatch expected
	}

	model := ResourceModel{
		Metro: types.StringValue("ny"), // Intentionally different
	}

	diags := model.parseMetalV1(vlan)

	assert.True(t, diags.HasError())
	assert.Contains(t, diags.Errors()[0].Summary(), "unexpected value for metro")
	assert.Equal(t, "vlan-999", model.ID.ValueString())
	assert.Equal(t, "ny", model.Metro.ValueString()) // Remains unchanged
}
