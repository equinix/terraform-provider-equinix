package vlan

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"strings"
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

func (m *DataSourceModel) parse(vlan *packngo.VirtualNetwork) diag.Diagnostics {
	m.ID = types.StringValue(vlan.ID)
	m.VlanID = types.StringValue(vlan.ID)
	m.Description = types.StringValue(vlan.Description)
	m.Vxlan = types.Int64Value(int64(vlan.VXLAN))
	m.Facility = types.StringValue("")

	if vlan.Project.ID != "" {
		m.ProjectID = types.StringValue(vlan.Project.ID)
	}

	if vlan.Facility != nil {
		m.Facility = types.StringValue(strings.ToLower(vlan.Facility.Code))
		m.Metro = types.StringValue(strings.ToLower(vlan.Facility.Metro.Code))
	}

	if vlan.Metro != nil {
		m.Metro = types.StringValue(strings.ToLower(vlan.Metro.Code))
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

func (m *ResourceModel) parse(vlan *packngo.VirtualNetwork) diag.Diagnostics {
	m.ID = types.StringValue(vlan.ID)
	m.Description = types.StringValue(vlan.Description)
	m.Vxlan = types.Int64Value(int64(vlan.VXLAN))
	m.Facility = types.StringValue("")

	if vlan.Project.ID != "" {
		m.ProjectID = types.StringValue(vlan.Project.ID)
	}

	if vlan.Facility != nil {
		m.Facility = types.StringValue(strings.ToLower(vlan.Facility.Code))
		m.Metro = types.StringValue(strings.ToLower(vlan.Facility.Metro.Code))
	}

	if vlan.Metro != nil {
		m.Metro = types.StringValue(strings.ToLower(vlan.Metro.Code))
	}

	return nil
}
