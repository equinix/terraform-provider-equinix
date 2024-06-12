package equinix

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricConnectionResourceSchema() map[string]*schema.Schema {
	sch := fabricConnectionResourceSchema()
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

func dataSourceFabricConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricConnectionRead,
		Schema:      readFabricConnectionResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch connection for a given UUID",
	}
}

func dataSourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricConnectionRead(ctx, d, meta)
}
