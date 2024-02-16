package vlans

import (
	"context"
	"fmt"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/packethost/packngo"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_metal_vlan",
			},
		),
	}
}

func (r *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	s := dataSourceSchema()
	if s.Blocks == nil {
		s.Blocks = make(map[string]schema.Block)
	}
	resp.Schema = s
}

type DataSource struct {
	framework.BaseDataSource
	framework.WithTimeouts
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.VlanID.IsNull() &&
		(data.Vxlan.IsNull() && data.ProjectID.IsNull() && data.Metro.IsNull() && data.Facility.IsNull()) {
		resp.Diagnostics.AddError("Error fetching Vlan datasource",
			equinix_errors.
				FriendlyError(fmt.Errorf("You must set either vlan_id or a combination of vxlan, project_id, and, metro or facility")).
				Error())
		return
	}

	var vlan *packngo.VirtualNetwork

	if !data.VlanID.IsNull() {
		var err error
		vlan, _, err = client.ProjectVirtualNetworks.Get(
			data.VlanID.ValueString(),
			&packngo.GetOptions{Includes: []string{"assigned_to"}},
		)
		if err != nil {
			resp.Diagnostics.AddError("Error fetching Vlan using vlanId", equinix_errors.FriendlyError(err).Error())
			return
		}

	} else {
		vlans, _, err := client.ProjectVirtualNetworks.List(
			data.ProjectID.ValueString(),
			&packngo.GetOptions{Includes: []string{"assigned_to"}},
		)
		if err != nil {
			resp.Diagnostics.AddError("Error fetching vlans list for projectId",
				equinix_errors.FriendlyError(err).Error())
			return
		}

		vlan, err = MatchingVlan(vlans.VirtualNetworks, int(data.Vxlan.ValueInt64()), data.ProjectID.ValueString(),
			data.Facility.ValueString(), data.Metro.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error expected vlan not found", equinix_errors.FriendlyError(err).Error())
			return
		}
	}

	assignedDevices := []string{}
	for _, d := range vlan.Instances {
		assignedDevices = append(assignedDevices, d.ID)
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(data.parse(vlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func MatchingVlan(vlans []packngo.VirtualNetwork, vxlan int, projectID, facility, metro string) (*packngo.VirtualNetwork, error) {
	matches := []packngo.VirtualNetwork{}
	for _, v := range vlans {
		if vxlan != 0 && v.VXLAN != vxlan {
			continue
		}
		if facility != "" && v.FacilityCode != facility {
			continue
		}
		if metro != "" && v.MetroCode != metro {
			continue
		}
		matches = append(matches, v)
	}
	if len(matches) > 1 {
		return nil, equinix_errors.FriendlyError(fmt.Errorf("Project %s has more than one matching VLAN", projectID))
	}

	if len(matches) == 0 {
		return nil, equinix_errors.FriendlyError(fmt.Errorf("Project %s does not have matching VLANs", projectID))
	}
	return &matches[0], nil
}
