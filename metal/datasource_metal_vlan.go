package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			},
			"metro": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vlan_id", "facility"},
				Description:   "Metro where the VLAN is deployed",
				StateFunc:     toLower,
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

func dataSourceMetalVlanRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*packngo.Client)

	projectRaw, projectOk := d.GetOk("project_id")
	vxlanRaw, vxlanOk := d.GetOk("vxlan")
	vlanIdRaw, vlanIdOk := d.GetOk("vlan_id")
	metroRaw, metroOk := d.GetOk("metro")
	facilityRaw, facilityOk := d.GetOk("facility")

	if !(vlanIdOk || (vxlanOk || projectOk || metroOk || facilityOk)) {
		return friendlyError(fmt.Errorf("You must set either vlan_id or a combination of vxlan, project_id, and, metro or facility"))
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
		projectID := projectRaw.(string)
		vxlan := vxlanRaw.(int)
		metro := metroRaw.(string)
		facility := facilityRaw.(string)
		vlans, _, err := c.ProjectVirtualNetworks.List(
			projectRaw.(string),
			&packngo.GetOptions{Includes: []string{"assigned_to"}},
		)
		if err != nil {
			return friendlyError(err)
		}

		matches := 0
		for _, v := range vlans.VirtualNetworks {
			if vxlan != 0 && v.VXLAN != vxlan {
				continue
			}
			if facility != "" && v.FacilityCode != facility {
				continue
			}
			if metro != "" && v.MetroCode != metro {
				continue
			}
			matches++
			if matches > 1 {
				return friendlyError(fmt.Errorf("Project %s has more than one matching VLAN", projectID))
			}
			vlan = &v
		}

		if matches == 0 {
			return friendlyError(fmt.Errorf("Project %s does not have matching VLANs", projectID))
		}

		if vlan == nil {
			locName := "facility"
			loc := facility
			if metroOk {
				locName = "metro"
				loc = metro
			}
			return friendlyError(fmt.Errorf("Project %s doesn't contain VLAN with vxlan %d in %s %s", projectID, vxlan, locName, loc))
		}
	}
	assignedDevices := []string{}
	for _, d := range vlan.Instances {
		assignedDevices = append(assignedDevices, d.ID)
	}

	d.SetId(vlan.ID)

	return setMap(d, map[string]interface{}{
		"vlan_id":     vlan.ID,
		"project_id":  vlan.Project.ID,
		"vxlan":       vlan.VXLAN,
		"facility":    vlan.FacilityCode,
		"metro":       vlan.MetroCode,
		"description": vlan.Description,
	})
}
