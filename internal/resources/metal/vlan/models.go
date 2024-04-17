package vlan

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
)

type DataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ProjectID          types.String `tfsdk:"project_id"`
	VlanID             types.String `tfsdk:"vlan_id"`
	Vxlan              types.Int64  `tfsdk:"vxlan"`
	Facility           types.String `tfsdk:"facility"`
	Metro              types.String `tfsdk:"metro"`
	Description        types.String `tfsdk:"description"`
	AssignedDevicesIds types.List   `tfsdk:"assigned_devices_ids"`
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

	deviceIds := make([]types.String, 0, len(vlan.Instances))
	for _, device := range vlan.Instances {
		deviceIds = append(deviceIds, types.StringValue(device.ID))
	}

	return m.AssignedDevicesIds.ElementsAs(context.Background(), &deviceIds, false)
}

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
	return nil
}
