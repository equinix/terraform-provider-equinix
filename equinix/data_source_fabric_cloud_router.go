package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceCloudRouter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudRouterRead,
		Schema:      readCloudRouterResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch Fabric Cloud Router for a given UUID",
	}
}

func dataSourceCloudRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceCloudRouterRead(ctx, d, meta)
}
