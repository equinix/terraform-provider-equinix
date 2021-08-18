package metal

import (
	"errors"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalVlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalVlanCreate,
		Read:   resourceMetalVlanRead,
		Delete: resourceMetalVlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "ID of parent project",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description string",
				Optional:    true,
				ForceNew:    true,
			},
			"facility": {
				Type:          schema.TypeString,
				Description:   "Facility where to create the VLAN",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"metro"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// suppress diff when unsetting facility
					if len(old) > 0 && new == "" {
						return true
					}
					return old == new
				},
			},
			"metro": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"facility"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, facOk := d.GetOk("facility")
					// new - new val from template
					// old - old val from state
					//
					// suppress diff if metro is manually set for first time, and
					// facility is already set
					if len(new) > 0 && old == "" && facOk {
						return facOk
					}
					return old == new
				},
				StateFunc: toLower,
			},
			"vxlan": {
				Type:        schema.TypeInt,
				Description: "VLAN ID, must be unique in metro",
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceMetalVlanCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*packngo.Client)

	facRaw, facOk := d.GetOk("facility")
	metroRaw, metroOk := d.GetOk("metro")
	vxlanRaw, vxlanOk := d.GetOk("vxlan")

	if !facOk && !metroOk {
		return friendlyError(errors.New("one of facility or metro must be configured"))
	}
	if facOk && vxlanOk {
		return friendlyError(errors.New("you can set vxlan only for metro vlans"))
	}

	createRequest := &packngo.VirtualNetworkCreateRequest{
		ProjectID:   d.Get("project_id").(string),
		Description: d.Get("description").(string),
	}
	if metroOk {
		createRequest.Metro = metroRaw.(string)
		createRequest.VXLAN = vxlanRaw.(int)
	}
	if facOk {
		createRequest.Facility = facRaw.(string)
	}
	vlan, _, err := c.ProjectVirtualNetworks.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}
	d.SetId(vlan.ID)
	return resourceMetalVlanRead(d, meta)
}

func resourceMetalVlanRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*packngo.Client)

	vlan, _, err := c.ProjectVirtualNetworks.Get(d.Id(),
		&packngo.GetOptions{Includes: []string{"assigned_to"}})
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err

	}
	d.Set("description", vlan.Description)
	d.Set("project_id", vlan.Project.ID)
	d.Set("vxlan", vlan.VXLAN)
	d.Set("facility", vlan.FacilityCode)
	d.Set("metro", vlan.MetroCode)
	return nil
}

func resourceMetalVlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	id := d.Id()
	vlan, resp, err := client.ProjectVirtualNetworks.Get(id, &packngo.GetOptions{Includes: []string{"instances", "instances.network_ports.virtual_networks", "internet_gateway"}})
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	} else if err != nil {
		// missing vlans are deleted
		return nil
	}

	// all device ports must be unassigned before delete
	for _, i := range vlan.Instances {
		for _, p := range i.NetworkPorts {
			for _, a := range p.AttachedVirtualNetworks {
				// a.ID is not set despite including instaces.network_ports.virtual_networks
				// TODO(displague) packngo should offer GetID() that uses ID or Href
				aID := path.Base(a.Href)

				if aID == id {
					_, resp, err := client.Ports.Unassign(p.ID, id)

					if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
						return friendlyError(err)
					}
				}
			}
		}
	}

	// TODO(displague) do we need to unassign gateway connections before delete?

	return friendlyError(ignoreResponseErrors(httpForbidden, httpNotFound)(client.ProjectVirtualNetworks.Delete(id)))
}
