package equinix

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalVRF() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Read:   resourceMetalVRFRead,
		Create: resourceMetalVRFCreate,
		Update: resourceMetalVRFUpdate,
		Delete: resourceMetalVRFDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User-supplied name of the VRF, unique to the project",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the VRF",
			},
			"metro": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Metro Code",
			},
			"local_asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "",
			},
			"ip_ranges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "All IPv4 and IPv6 Ranges that will be available to BGP Peers. IPv4 addresses must be /8 or smaller with a minimum size of /29. IPv6 must be /56 or smaller with a minimum size of /64. Ranges must not overlap other ranges within the VRF.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID",
			},
			// TODO: created_by, created_at, updated_at, href
		},
	}
}

func resourceMetalVRFCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal

	createRequest := &packngo.VRFCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Metro:       d.Get("metro").(string),
		LocalASN:    d.Get("local_asn").(int),
		IPRanges:    expandSetToStringList(d.Get("ip_ranges").(*schema.Set)),
	}

	projectId := d.Get("project_id").(string)
	vrf, _, err := client.VRFs.Create(projectId, createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(vrf.ID)

	return resourceMetalVRFRead(d, meta)
}

func resourceMetalVRFUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal

	sPtr := func(s string) *string { return &s }
	iPtr := func(i int) *int { return &i }

	updateRequest := &packngo.VRFUpdateRequest{}
	if d.HasChange("name") {
		updateRequest.Name = sPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateRequest.Description = sPtr(d.Get("description").(string))
	}
	if d.HasChange("local_asn") {
		updateRequest.LocalASN = iPtr(d.Get("local_asn").(int))
	}
	if d.HasChange("ip_ranges") {
		ipRanges := expandSetToStringList(d.Get("ip_ranges").(*schema.Set))
		updateRequest.IPRanges = &ipRanges
	}

	_, _, err := client.VRFs.Update(d.Id(), updateRequest)
	if err != nil {
		return friendlyError(err)
	}

	return resourceMetalVRFRead(d, meta)
}

func resourceMetalVRFRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal

	getOpts := &packngo.GetOptions{Includes: []string{"project", "metro"}}

	vrf, _, err := client.VRFs.Get(d.Id(), getOpts)
	if err != nil {
		if isNotFound(err) || isForbidden(err) {
			log.Printf("[WARN] VRF (%s) not accessible, removing from state", d.Id())
			d.SetId("")

			return nil
		}
		return err
	}
	m := map[string]interface{}{
		"name":        vrf.Name,
		"description": vrf.Description,
		"metro":       vrf.Metro.Code,
		"local_asn":   vrf.LocalASN,
		"ip_ranges":   vrf.IPRanges,
		"project_id":  vrf.Project.ID,
	}

	return setMap(d, m)
}

func resourceMetalVRFDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).metal

	resp, err := client.VRFs.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) == nil {
		d.SetId("")
	}

	return friendlyError(err)
}
