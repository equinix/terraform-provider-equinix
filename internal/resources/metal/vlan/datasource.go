package vlan

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/converters"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceMetalVlanRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
	client := meta.(*config.Config).Metalgo

	projectRaw, projectOk := d.GetOk("project_id")
	vxlanRaw, vxlanOk := d.GetOk("vxlan")
	vlanIdRaw, vlanIdOk := d.GetOk("vlan_id")
	metroRaw, metroOk := d.GetOk("metro")
	facilityRaw, facilityOk := d.GetOk("facility")

	if !(vlanIdOk || (vxlanOk || projectOk || metroOk || facilityOk)) {
		return diag.Errorf("You must set either vlan_id or a combination of vxlan, project_id, and, metro or facility")
	}

	var vlan *metalv1.VirtualNetwork

	if vlanIdOk {
		var err error
		vlan, _, err = client.VLANsApi.
			GetVirtualNetwork(ctx, vlanIdRaw.(string)).
			Include([]string{"assigned_to"}).Execute()

		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}

	} else {
		projectID := projectRaw.(string)
		vxlan := vxlanRaw.(int)
		metro := metroRaw.(string)
		facility := facilityRaw.(string)
		vlans, _, err := client.VLANsApi.
			FindVirtualNetworks(ctx, projectRaw.(string)).
			Include([]string{"assigned_to"}).Execute()

		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}

		vlan, err = MatchingVlan(vlans.VirtualNetworks, vxlan, projectID, facility, metro)
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}
	}

	assignedDevices := []string{}
	for _, d := range vlan.Instances {
		assignedDevices = append(assignedDevices, d.GetId()) // instances is a list of href, should be list of device?
	}

	d.SetId(vlan.GetId())

	return diag.FromErr(equinix_schema.SetMap(d, map[string]interface{}{
		"vlan_id":     vlan.GetId(),
		"project_id":  vlan.AssignedTo.GetId(), // vlan assigned_to is an href; should be project?
		"vxlan":       vlan.GetVxlan(),
		"facility":    vlan.FacilityCode, // facility is deprecated, vlan is metro-scoped; remove this attr?
		"metro":       vlan.MetroCode,
		"description": vlan.Description,
	}))
}
