package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFabricConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricConnectionRead,
		Schema:      readFabricConnectionResourceSchema(),
	}
}

func dataSourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricConnectionRead(ctx, d, meta)
}
