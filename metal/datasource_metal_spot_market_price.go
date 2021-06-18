package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceSpotMarketPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalSpotMarketPriceRead,
		Schema: map[string]*schema.Schema{
			"facility": {
				Type:          schema.TypeString,
				Description:   "Name of the facility",
				ConflictsWith: []string{"metro"},
				Optional:      true,
			},
			"metro": {
				Type:          schema.TypeString,
				Description:   "Name of the metro",
				ConflictsWith: []string{"facility"},
				Optional:      true,
				StateFunc:     toLower,
			},
			"plan": {
				Type:        schema.TypeString,
				Description: "Name of the plan",
				Required:    true,
			},
			"price": {
				Type:        schema.TypeFloat,
				Description: "Current spot market price for given plan in given facility",
				Computed:    true,
			},
		},
	}
}

func dataSourceMetalSpotMarketPriceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	sms := client.SpotMarket.(*packngo.SpotMarketServiceOp)
	facility := d.Get("facility").(string)
	metro := d.Get("metro").(string)
	plan := d.Get("plan").(string)

	if facility != "" && metro != "" {
		return fmt.Errorf("Parameters facility and metro cannot be used together")
	}

	filter := facility
	fn := sms.PricesByFacility
	filterType := "facility"

	if metro != "" {
		filter = metro
		fn = sms.PricesByMetro
		filterType = "metro"
	}

	prices, _, err := fn()
	if err != nil {
		return err
	}

	match, ok := prices[filter]
	if !ok {
		return fmt.Errorf("Cannot find %s %s", filterType, filter)
	}

	price, ok := match[plan]
	if !ok {
		return fmt.Errorf("Cannot find price for plan %s in %s %s", plan, filterType, filter)
	}

	d.Set("price", price)
	d.SetId(fmt.Sprintf("%s-%s-%s", filterType, filter, plan))
	return nil
}
