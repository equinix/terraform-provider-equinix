package vlan

import (
	"context"
	"fmt"
	"strings"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
)

// DataSourceModel represents the schema for reading VLAN data from the Equinix Metal API
// in a Terraform data source context. It includes fields that describe the VLAN's identity,
// location, and associated devices, and is used to populate Terraform state during read operations.
type DataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ProjectID          types.String `tfsdk:"project_id"`
	VlanID             types.String `tfsdk:"vlan_id"`
	Vxlan              types.Int64  `tfsdk:"vxlan"`
	Facility           types.String `tfsdk:"facility"`
	Metro              types.String `tfsdk:"metro"`
	Description        types.String `tfsdk:"description"`
	AssignedDevicesIDs types.List   `tfsdk:"assigned_devices_ids"`
}

func (m *DataSourceModel) parse(vlan *packngo.VirtualNetwork) (d diag.Diagnostics) {
	m.ID = types.StringValue(vlan.ID)
	m.VlanID = types.StringValue(vlan.ID)
	m.Vxlan = types.Int64Value(int64(vlan.VXLAN))
	m.Facility = types.StringValue("")

	if vlan.Description != "" {
		m.Description = types.StringValue(vlan.Description)
	}

	if vlan.Project.ID != "" {
		m.ProjectID = types.StringValue(vlan.Project.ID)
	}

	if vlan.Facility != nil {
		m.Facility = types.StringValue(strings.ToLower(vlan.Facility.Code))
		m.Metro = types.StringValue(strings.ToLower(vlan.Facility.Metro.Code))
	}

	deviceIDs := make([]types.String, 0, len(vlan.Instances))
	for _, device := range vlan.Instances {
		deviceIDs = append(deviceIDs, types.StringValue(device.ID))
	}

	m.AssignedDevicesIDs, d = types.ListValueFrom(context.Background(), types.StringType, deviceIDs)

	if vlan.Metro != nil {
		if m.Metro.IsNull() {
			m.Metro = types.StringValue(vlan.Metro.Code)
		} else if !strings.EqualFold(m.Metro.ValueString(), vlan.Metro.Code) {
			d.AddWarning(
				"unexpected value for metro",
				fmt.Sprintf("expected vlan %v to have metro %v, but metro was %v",
					m.ID, m.Metro, vlan.Metro.Code))
		}
	}

	return d
}

// ResourceModel defines the schema for managing VLAN resources in Terraform.
// It is used during create, read, update, and delete operations to represent
// the desired and actual state of a VLAN in Equinix Metal.
type ResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Vxlan       types.Int64  `tfsdk:"vxlan"`
	Facility    types.String `tfsdk:"facility"`
	Metro       types.String `tfsdk:"metro"`
	Description types.String `tfsdk:"description"`
}

func (m *ResourceModel) parse(vlan *packngo.VirtualNetwork) (d diag.Diagnostics) {
	m.ID = types.StringValue(vlan.ID)
	m.Vxlan = types.Int64Value(int64(vlan.VXLAN))
	m.Facility = types.StringValue("")

	if vlan.Description != "" {
		m.Description = types.StringValue(vlan.Description)
	}

	if vlan.Project.ID != "" {
		m.ProjectID = types.StringValue(vlan.Project.ID)
	}

	if vlan.Facility != nil {
		m.Facility = types.StringValue(strings.ToLower(vlan.Facility.Code))
		m.Metro = types.StringValue(strings.ToLower(vlan.Facility.Metro.Code))
	}

	if vlan.Metro != nil {
		if m.Metro.IsNull() {
			m.Metro = types.StringValue(vlan.Metro.Code)
		} else if !strings.EqualFold(m.Metro.ValueString(), vlan.Metro.Code) {
			d.AddError(
				"unexpected value for metro",
				fmt.Sprintf("expected vlan %v to have metro %v, but metro was %v",
					m.ID, m.Metro, vlan.Metro.Code))
		}
	}
	return d
}

func (m *ResourceModel) parseMetalV1(vlan *metalv1.VirtualNetwork) (d diag.Diagnostics) {
	m.ID = types.StringValue(vlan.GetId())
	m.Vxlan = types.Int64Value(int64(vlan.GetVxlan()))
	m.Facility = types.StringValue("")

	if vlan.GetDescription() != "" {
		m.Description = types.StringValue(vlan.GetDescription())
	}

	if vlan.AssignedTo != nil {
		m.ProjectID = types.StringValue(vlan.AssignedTo.GetId())
	}

	if vlan.Facility != nil {
		facilityCode := vlan.Facility.AdditionalProperties["facility_code"].(string)
		m.Facility = types.StringValue(strings.ToLower(facilityCode))
	}

	if vlan.MetroCode != nil {
		if m.Metro.IsNull() {
			m.Metro = types.StringValue(vlan.GetMetroCode())
		} else if !strings.EqualFold(m.Metro.ValueString(), vlan.GetMetroCode()) {
			d.AddError(
				"unexpected value for metro",
				fmt.Sprintf("expected vlan %v to have metro %v, but metro was %v",
					m.ID, m.Metro, vlan.GetMetroCode()))
		}
	}

	return d
}
