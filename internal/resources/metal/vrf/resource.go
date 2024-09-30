package vrf

import (
	"context"
	"log"
	"net/http"

	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
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
				Description: "Metro ID or Code where the VRF will be deployed",
				ForceNew:    true,
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
				Description: "Project ID where the VRF will be deployed",
				ForceNew:    true,
			},
			// TODO: created_by, created_at, updated_at, href
		},
	}
}

func resourceMetalVRFCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	createRequest := metalv1.VrfCreateInput{
		Name:        d.Get("name").(string),
		Description: metalv1.PtrString(d.Get("description").(string)),
		Metro:       d.Get("metro").(string),
		IpRanges:    converters.SetToStringList(d.Get("ip_ranges").(*schema.Set)),
	}

	if value, ok := d.GetOk("local_asn"); ok {
		createRequest.LocalAsn = metalv1.PtrInt64(int64(value.(int)))
	}

	projectId := d.Get("project_id").(string)
	vrf, _, err := client.VRFsApi.
		CreateVrf(ctx, projectId).
		VrfCreateInput(createRequest).
		Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vrf.GetId())

	return resourceMetalVRFRead(ctx, d, meta)
}

func resourceMetalVRFUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	updateRequest := metalv1.VrfUpdateInput{}
	if d.HasChange("name") {
		updateRequest.SetName(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateRequest.SetDescription(d.Get("description").(string))
	}
	if d.HasChange("local_asn") {
		updateRequest.SetLocalAsn(int32(d.Get("local_asn").(int)))
	}
	if d.HasChange("ip_ranges") {
		ipRanges := converters.SetToStringList(d.Get("ip_ranges").(*schema.Set))
		updateRequest.SetIpRanges(ipRanges)
	}

	_, _, err := client.VRFsApi.
		UpdateVrf(ctx, d.Id()).
		VrfUpdateInput(updateRequest).
		Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMetalVRFRead(ctx, d, meta)
}

func resourceMetalVRFRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	vrf, _, err := client.VRFsApi.
		FindVrfById(ctx, d.Id()).
		Include([]string{"project", "metro"}).
		Execute()
	if err != nil {
		if equinix_errors.IsNotFound(err) || equinix_errors.IsForbidden(err) {
			log.Printf("[WARN] VRF (%s) not accessible, removing from state", d.Id())
			d.SetId("")

			return nil
		}
		return diag.FromErr(err)
	}
	m := map[string]interface{}{
		"name":        vrf.GetName(),
		"description": vrf.GetDescription(),
		"metro":       vrf.Metro.GetCode(),
		"local_asn":   vrf.GetLocalAsn(),
		"ip_ranges":   vrf.GetIpRanges(),
		"project_id":  vrf.Project.GetId(),
	}

	return diag.FromErr(equinix_schema.SetMap(d, m))
}

func resourceMetalVRFDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	resp, err := client.VRFsApi.DeleteVrf(ctx, d.Id()).Execute()
	if equinix_errors.IgnoreHttpResponseErrors(http.StatusForbidden, http.StatusNotFound)(resp, err) == nil {
		d.SetId("")
	}

	return diag.FromErr(err)
}
