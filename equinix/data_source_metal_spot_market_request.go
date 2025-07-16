package equinix

import (
	"context"
	"sort"
	"strings"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalSpotMarketRequest() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "The Spot Market Requests API has been sunset.",
		ReadContext:        dataSourceMetalSpotMarketRequestRead,

		Schema: map[string]*schema.Schema{
			"request_id": {
				Type:        schema.TypeString,
				Description: "The id of the Spot Market Request",
				Required:    true,
			},
			"device_ids": {
				Type:        schema.TypeList,
				Description: "List of IDs of devices spawned by the referenced Spot Market Request",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"devices_min": {
				Type:        schema.TypeInt,
				Description: "Miniumum number devices to be created",
				Computed:    true,
			},
			"devices_max": {
				Type:        schema.TypeInt,
				Description: "Maximum number devices to be created",
				Computed:    true,
			},
			"max_bid_price": {
				Type:        schema.TypeFloat,
				Description: "Maximum price user is willing to pay per hour per device",
				Computed:    true,
			},
			"facilities": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Facility IDs where devices should be created",
				Deprecated:  "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
				Computed:    true,
			},
			"metro": {
				Type:        schema.TypeString,
				Description: "Metro where devices should be created.",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "Project ID",
				Computed:    true,
			},
			"plan": {
				Type:        schema.TypeString,
				Description: "The device plan slug.",
				Computed:    true,
			},
			"end_at": {
				Type:        schema.TypeString,
				Description: "Date and time When the spot market request will be ended.",
				Computed:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(60 * time.Minute),
			Update:  schema.DefaultTimeout(60 * time.Minute),
			Delete:  schema.DefaultTimeout(60 * time.Minute),
			Default: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func dataSourceMetalSpotMarketRequestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).Metal
	id := d.Get("request_id").(string)

	smr, _, err := client.SpotMarketRequests.Get(id, &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		if equinix_errors.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	deviceIDs := make([]string, len(smr.Devices))
	for i, d := range smr.Devices {
		deviceIDs[i] = d.ID
	}

	facs := smr.Facilities
	facCodes := []string{}

	for _, f := range facs {
		facCodes = append(facCodes, f.Code)
	}
	sort.Strings(facCodes) // avoid changes if we get the same facilities in a different order

	d.SetId(id)

	err = equinix_schema.SetMap(d, map[string]interface{}{
		"device_ids": deviceIDs,
		"end_at": func(d *schema.ResourceData, k string) error {
			if smr.EndAt != nil {
				return d.Set(k, smr.EndAt.Format(time.RFC3339))
			}
			return nil
		},
		"devices_max": smr.DevicesMax,
		"devices_min": smr.DevicesMin,
		"facilities":  facCodes,
		"metro": func(d *schema.ResourceData, k string) error {
			if smr.Metro != nil {
				return d.Set(k, strings.ToLower(smr.Metro.Code))
			}
			return nil
		},
		"max_bid_price": smr.MaxBidPrice,
		"plan":          smr.Plan.Slug,
		"project_id":    smr.Project.ID,
		// TODO: created_at is not in packngo
	})

	return diag.FromErr(err)
}
