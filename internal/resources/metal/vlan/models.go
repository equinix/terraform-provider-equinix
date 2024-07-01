package vlan

import (
	"context"
	"fmt"
	"strings"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (m *DataSourceModel) parse(vlan *metalv1.VirtualNetwork) (d diag.Diagnostics) {
	m.ID = types.StringValue(vlan.GetId())
	m.VlanID = types.StringValue(vlan.GetId())
	m.Vxlan = types.Int64Value(int64(vlan.GetVxlan()))
	m.Facility = types.StringValue("")

	if vlan.GetDescription() != "" {
		m.Description = types.StringValue(vlan.GetDescription())
	}

	if vlan.AdditionalProperties["project"] != nil {
		project := vlan.AdditionalProperties["project"].(map[string]interface{})
		project_id := project["id"].(string)
		m.ProjectID = types.StringValue(project_id)
	}

	if vlan.Facility != nil {
		facility_code := vlan.Facility.AdditionalProperties["code"].(string)
		metro := vlan.Facility.AdditionalProperties["metro"].(map[string]interface{})
		metro_code := metro["code"].(string)
		m.Facility = types.StringValue(strings.ToLower(facility_code))
		m.Metro = types.StringValue(strings.ToLower(metro_code))
	}

	if vlan.Metro != nil {
		if m.Metro.IsNull() {
			m.Metro = types.StringValue(vlan.Metro.GetCode())
		} else if !strings.EqualFold(m.Metro.ValueString(), vlan.Metro.GetCode()) {
			d.AddWarning(
				"unexpected value for metro",
				fmt.Sprintf("expected vlan %v to have metro %v, but metro was %v",
					m.ID, m.Metro, vlan.Metro.Code))
		}
	}

	deviceIds := make([]types.String, 0, len(vlan.Instances))
	for _, device := range vlan.Instances {
		deviceIds = append(deviceIds, types.StringValue(device.GetId()))
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

func (m *ResourceModel) parse(vlan *metalv1.VirtualNetwork) (d diag.Diagnostics) {
	m.ID = types.StringValue(vlan.GetId())
	m.Vxlan = types.Int64Value(int64(vlan.GetVxlan()))
	m.Facility = types.StringValue("")

	if vlan.GetDescription() != "" {
		m.Description = types.StringValue(vlan.GetDescription())
	}

	if vlan.AdditionalProperties["project"] != nil {
		project := vlan.AdditionalProperties["project"].(map[string]interface{})
		project_id := project["id"].(string)
		m.ProjectID = types.StringValue(project_id)
	}

	if vlan.Facility != nil {
		facility_code := vlan.Facility.AdditionalProperties["code"].(string)
		metro := vlan.Facility.AdditionalProperties["metro"].(map[string]interface{})
		metro_code := metro["code"].(string)
		m.Facility = types.StringValue(strings.ToLower(facility_code))
		m.Metro = types.StringValue(strings.ToLower(metro_code))
	}

	if vlan.Metro != nil {
		if m.Metro.IsNull() {
			m.Metro = types.StringValue(vlan.Metro.GetCode())
		} else if !strings.EqualFold(m.Metro.ValueString(), vlan.Metro.GetCode()) {
			d.AddError(
				"unexpected value for metro",
				fmt.Sprintf("expected vlan %v to have metro %v, but metro was %v",
					m.ID, m.Metro, vlan.Metro.Code))
		}
	}
	return nil
}
