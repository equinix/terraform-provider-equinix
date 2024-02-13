package vlans

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceMetalVlanRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vlan_id"},
				Description:   "ID of parent project of the VLAN. Use together with vxlan and metro or facility",
			},
			"vxlan": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vlan_id"},
				Description:   "VXLAN numner of the VLAN. Unique in a project and facility or metro. Use with project_id",
			},
			"facility": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vlan_id", "metro"},
				Description:   "Facility where the VLAN is deployed",
				Deprecated:    "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
			},
			"metro": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vlan_id", "facility"},
				Description:   "Metro where the VLAN is deployed",
				StateFunc:     converters.ToLowerIf,
			},
			"vlan_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"project_id", "vxlan", "metro", "facility"},
				Description:   "Metal UUID of the VLAN resource",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VLAN description text",
			},
			"assigned_devices_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of device IDs to which this VLAN is assigned",
			},
		},
	}
}

func dataSourceMetalVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).Metal

	projectRaw, projectOk := d.GetOk("project_id")
	vxlanRaw, vxlanOk := d.GetOk("vxlan")
	vlanIdRaw, vlanIdOk := d.GetOk("vlan_id")
	metroRaw, metroOk := d.GetOk("metro")
	facilityRaw, facilityOk := d.GetOk("facility")

	if !(vlanIdOk || (vxlanOk || projectOk || metroOk || facilityOk)) {
		return diag.FromErr(equinix_errors.FriendlyError(fmt.Errorf("You must set either vlan_id or a combination of vxlan, project_id, and, metro or facility")))
	}

	var vlan *packngo.VirtualNetwork

	if vlanIdOk {
		var err error
		vlan, _, err = client.ProjectVirtualNetworks.Get(
			vlanIdRaw.(string),
			&packngo.GetOptions{Includes: []string{"assigned_to"}},
		)
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}

	} else {
		projectID := projectRaw.(string)
		vxlan := vxlanRaw.(int)
		metro := metroRaw.(string)
		facility := facilityRaw.(string)
		vlans, _, err := client.ProjectVirtualNetworks.List(
			projectRaw.(string),
			&packngo.GetOptions{Includes: []string{"assigned_to"}},
		)
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}

		vlan, err = matchingVlan(vlans.VirtualNetworks, vxlan, projectID, facility, metro)
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}
	}

	assignedDevices := []string{}
	for _, d := range vlan.Instances {
		assignedDevices = append(assignedDevices, d.ID)
	}

	d.SetId(vlan.ID)

	return diag.FromErr(equinix_schema.SetMap(d, map[string]interface{}{
		"vlan_id":     vlan.ID,
		"project_id":  vlan.Project.ID,
		"vxlan":       vlan.VXLAN,
		"facility":    vlan.FacilityCode,
		"metro":       vlan.MetroCode,
		"description": vlan.Description,
	}))
}

func matchingVlan(vlans []packngo.VirtualNetwork, vxlan int, projectID, facility, metro string) (*packngo.VirtualNetwork, error) {
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
