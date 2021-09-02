package metal

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalFacility() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalFacilityRead,
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the Facility to match",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of this Facility.",
				Computed:    true,
			},
			"metro": {
				Type:        schema.TypeString,
				Description: "This facility's metro code.",
				Computed:    true,
			},
			"features": {
				Type:        schema.TypeList,
				Description: "The features of this Facility.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"capacity": {
				Type:        schema.TypeList,
				Description: "Optional capacity specification",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan": {
							Type:        schema.TypeString,
							Description: "Plan which has to be available in selected facility",
							Required:    true,
						},
						"quantity": {
							Type:     schema.TypeInt,
							Default:  1,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func getCapacityInput(capacitySpecs []interface{}, baseServerInfo packngo.ServerInfo) *packngo.CapacityInput {
	ci := packngo.CapacityInput{Servers: []packngo.ServerInfo{}}
	for _, v := range capacitySpecs {
		item := v.(map[string]interface{})
		spec := baseServerInfo
		spec.Plan = item["plan"].(string)
		spec.Quantity = item["quantity"].(int)
		ci.Servers = append(ci.Servers, spec)
	}
	return &ci
}

func dataSourceMetalFacilityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	code := d.Get("code").(string)

	_, capacityOk := d.GetOk("capacity")
	if capacityOk {
		ci := getCapacityInput(
			d.Get("capacity").([]interface{}),
			packngo.ServerInfo{Facility: code},
		)
		res, _, err := client.CapacityService.Check(ci)
		if err != nil {
			return err
		}
		for _, s := range res.Servers {
			if !s.Available {
				return fmt.Errorf("Not enough capacity in facility %s for %d device(s) of plan %s", s.Facility, s.Quantity, s.Plan)
			}
		}
	}

	if code == "" {
		return fmt.Errorf("Error Facility code is required")
	}

	facilities, _, err := client.Facilities.List(nil)
	if err != nil {
		return fmt.Errorf("Error listing Facilities: %s", err)
	}

	for _, f := range facilities {
		if f.Code == code {
			d.SetId(f.ID)
			return setMap(d, map[string]interface{}{
				"code":     f.Code,
				"name":     f.Name,
				"features": f.Features,
				"metro": func(d *schema.ResourceData, k string) error {
					if f.Metro != nil {
						return d.Set(k, strings.ToLower(f.Metro.Code))
					}
					return nil
				},
			})
		}
	}

	return fmt.Errorf("Facility %s was not found", code)
}
