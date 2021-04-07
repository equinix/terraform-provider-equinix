package metal

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"facility": {
				Type:          schema.TypeString,
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
			},
			"vxlan": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceMetalVlanCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*packngo.Client)

	facRaw, facOk := d.GetOk("facility")
	metroRaw, metroOk := d.GetOk("metro")

	if !facOk && !metroOk {
		return friendlyError(errors.New("one of facility or metro must be configured"))
	}

	createRequest := &packngo.VirtualNetworkCreateRequest{
		ProjectID:   d.Get("project_id").(string),
		Description: d.Get("description").(string),
	}
	if metroOk {
		createRequest.Metro = metroRaw.(string)
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

	resp, err := client.ProjectVirtualNetworks.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}
	return nil
}
