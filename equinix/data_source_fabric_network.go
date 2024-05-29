package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricNetworkResourceSchema() map[string]*schema.Schema {
	sch := fabricNetworkResourceSchema()
	for key := range sch {
		if key == "uuid" {
			sch[key].Required = true
			sch[key].Optional = false
			sch[key].Computed = false
		} else {
			sch[key].Required = false
			sch[key].Optional = false
			sch[key].Computed = true
			sch[key].MaxItems = 0
			sch[key].ValidateFunc = nil
		}
	}
	return sch
}
func dataSourceFabricNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricNetworkRead,
		Schema:      readFabricNetworkResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch Fabric Network for a given UUID",
	}
}

func dataSourceFabricNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricNetworkRead(ctx, d, meta)
}
