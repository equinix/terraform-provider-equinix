package metro

import (
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/maps"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/facility"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalMetroRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of this Metro.",
				Computed:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the Metro to match",
				Required:    true,
			},
			"country": {
				Type:        schema.TypeString,
				Description: "The country of this Metro.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of this Metro.",
				Computed:    true,
			},
			"capacity": facility.CapacitySchema(),
		},
	}
}

func dataSourceMetalMetroRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).MetalClient
	code := d.Get("code").(string)

	_, capacityOk := d.GetOk("capacity")
	if capacityOk {
		ci := facility.GetCapacityInput(
			d.Get("capacity").([]interface{}),
			packngo.ServerInfo{Metro: code},
		)
		res, _, err := client.CapacityService.CheckMetros(ci)
		if err != nil {
			return err
		}
		for _, s := range res.Servers {
			if !s.Available {
				return fmt.Errorf("Not enough capacity in metro %s for %d device(s) of plan %s", s.Facility, s.Quantity, s.Plan)
			}
		}
		if err != nil {
			return err
		}
	}

	metros, _, err := client.Metros.List(nil)
	if err != nil {
		return fmt.Errorf("Error listing Metros: %s", err)
	}

	for _, m := range metros {
		if m.Code == code {
			d.SetId(m.ID)
			return maps.SetMap(d, map[string]interface{}{
				"id":      m.ID,
				"code":    m.Code,
				"name":    m.Name,
				"country": m.Country,
			})
		}
	}

	return fmt.Errorf("Metro %s was not found", code)
}
