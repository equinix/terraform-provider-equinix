package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalMetro() *schema.Resource {
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
		},
	}
}

func dataSourceMetalMetroRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	code := d.Get("code").(string)

	if code == "" {
		return fmt.Errorf("Error Metro code is required")
	}

	metros, _, err := client.Metros.List(nil)
	if err != nil {
		return fmt.Errorf("Error listing Metros: %s", err)
	}

	for _, m := range metros {
		if m.Code == code {
			d.SetId(m.ID)
			return setMap(d, map[string]interface{}{
				"id":      m.ID,
				"code":    m.Code,
				"name":    m.Name,
				"country": m.Country,
			})
		}
	}

	return fmt.Errorf("Metro %s was not found", code)
}
