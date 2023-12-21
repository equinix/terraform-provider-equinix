package equinix

import (
	"context"
	"log"

	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func ResourceMetalVRF() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout:   resourceMetalVRFRead,
		CreateWithoutTimeout: resourceMetalVRFCreate,
		UpdateWithoutTimeout: resourceMetalVRFUpdate,
		DeleteWithoutTimeout: resourceMetalVRFDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "The 4-byte ASN set on the VRF.",
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

func resourceMetalVRFCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	createRequest := &packngo.VRFCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Metro:       d.Get("metro").(string),
		LocalASN:    d.Get("local_asn").(int),
		IPRanges:    converters.SetToStringList(d.Get("ip_ranges").(*schema.Set)),
	}

	projectId := d.Get("project_id").(string)
	vrf, _, err := client.VRFs.Create(projectId, createRequest)
	if err != nil {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	}

	d.SetId(vrf.ID)

	return resourceMetalVRFRead(ctx, d, meta)
}

func resourceMetalVRFUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

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
		ipRanges := converters.SetToStringList(d.Get("ip_ranges").(*schema.Set))
		updateRequest.IPRanges = &ipRanges
	}

	_, _, err := client.VRFs.Update(d.Id(), updateRequest)
	if err != nil {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	}

	return resourceMetalVRFRead(ctx, d, meta)
}

func resourceMetalVRFRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	getOpts := &packngo.GetOptions{Includes: []string{"project", "metro"}}

	vrf, _, err := client.VRFs.Get(d.Id(), getOpts)
	if err != nil {
		if equinix_errors.IsNotFound(err) || equinix_errors.IsForbidden(err) {
			log.Printf("[WARN] VRF (%s) not accessible, removing from state", d.Id())
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}
	m := map[string]interface{}{
		"name":        vrf.Name,
		"description": vrf.Description,
		"metro":       vrf.Metro.Code,
		"local_asn":   vrf.LocalASN,
		"ip_ranges":   vrf.IPRanges,
		"project_id":  vrf.Project.ID,
	}

	return diag.FromErr(equinix_schema.SetMap(d, m))
}

func resourceMetalVRFDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	resp, err := client.VRFs.Delete(d.Id())
	if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) == nil {
		d.SetId("")
	}

	return diag.FromErr(equinix_errors.FriendlyError(err))
}
