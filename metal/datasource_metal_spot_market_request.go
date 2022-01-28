package metal

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalSpotMarketRequest() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalSpotMarketRequestRead,

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
		Timeouts: resourceDefaultTimeouts,
	}
}
func dataSourceMetalSpotMarketRequestRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	id := d.Get("request_id").(string)

	smr, _, err := client.SpotMarketRequests.Get(id, &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
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

	d.SetId(id)

	return setMap(d, map[string]interface{}{
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
}
