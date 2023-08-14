package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMetalVRF() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMetalVRFRead,

		Schema: map[string]*schema.Schema{
			"vrf_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the VRF to lookup",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User-supplied name of the VRF, unique to the project",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the VRF",
			},
			"metro": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Metro Code",
			},
			"local_asn": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The 4-byte ASN set on the VRF.",
			},
			"ip_ranges": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "All IPv4 and IPv6 Ranges that will be available to BGP Peers. IPv4 addresses must be /8 or smaller with a minimum size of /29. IPv6 must be /56 or smaller with a minimum size of /64. Ranges must not overlap other ranges within the VRF.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID",
			},
		},
	}
}

func dataSourceMetalVRFRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vrfId, _ := d.Get("vrf_id").(string)

	d.SetId(vrfId)
	return resourceMetalVRFRead(ctx, d, meta)
}
