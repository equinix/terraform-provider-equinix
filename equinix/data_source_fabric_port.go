package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFabricPort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricPortRead,
		Schema:      readFabricPortResourceSchema(),
	}
}

func dataSourceFabricPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricPortRead(ctx, d, meta)
}

func dataSourceFabricGetPortsByName() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricGetPortsByNameResponseRead,
		Schema:      readFabricPortsResponseSchema(),
	}
}

func dataSourceFabricGetPortsByNameResponseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFabricPortGetByPortName(ctx, d, meta)
}
