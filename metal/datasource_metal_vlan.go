package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalVlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalVlanRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vlan_id"},
				Description:   "ID of parent project of the VLAN. Use together with vxland",
			},
			"vxlan": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vlan_id"},
				Description:   "VXLAN numner of the VLAN. Unique in a project. Use with project_id",
			},
			"vlan_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"project_id", "vxlan"},
				Description:   "Metal UUID of the VLAN resource",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VLAN description text",
			},
			"facility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Facility where the VLAN is deployed",
			},
			"metro": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Metro where the VLAN is deployed",
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

func dataSourceMetalVlanRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*packngo.Client)

	projectRaw, projectOk := d.GetOk("project_id")
	vxlanRaw, vxlanOk := d.GetOk("vxlan")
	vlanIdRaw, vlanIdOk := d.GetOk("vlan_id")

	if !(vlanIdOk || vxlanOk || projectOk) {
		return friendlyError(fmt.Errorf("You must set either vlan_id or vxlan and project_id"))
	}

	var vlan *packngo.VirtualNetwork

	if vlanIdOk {
		var err error
		vlan, _, err = c.ProjectVirtualNetworks.Get(
			vlanIdRaw.(string),
			&packngo.GetOptions{Includes: []string{"assigned_to"}},
		)
		if err != nil {
			return friendlyError(err)
		}

	} else {
		if !(vxlanOk || projectOk) {
			return friendlyError(fmt.Errorf("If you set project_id, you also must set vxlan and vice versa"))
		}

		projectID := projectRaw.(string)
		vxlan := vxlanRaw.(int)
		vlans, _, err := c.ProjectVirtualNetworks.List(
			projectRaw.(string),
			&packngo.GetOptions{Includes: []string{"assigned_to"}},
		)
		if err != nil {
			return friendlyError(err)
		}
		for _, v := range vlans.VirtualNetworks {
			if v.VXLAN == vxlan {
				vlan = &v
			}
		}
		if vlan == nil {
			return friendlyError(fmt.Errorf("Project %s doesn't contain VLAN with vxlan %d", projectID, vxlan))
		}

	}
	assignedDevices := []string{}
	for _, d := range vlan.Instances {
		assignedDevices = append(assignedDevices, d.ID)
	}

	d.SetId(vlan.ID)

	return setMap(d, map[string]interface{}{
		"vlan_id":    vlan.ID,
		"project_id": vlan.Project.ID,
		"vxlan":      vlan.VXLAN,
		"facility":   vlan.FacilityCode,
		"metro":      vlan.MetroCode,
	})
}
