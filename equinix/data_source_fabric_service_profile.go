package equinix

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFabricServiceProfileReadByUuid() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricServiceProfileRead,
		Schema:      readFabricServiceProfileSchema(),
	}
}

func dataSourceFabricServiceProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricServiceProfileRead(ctx, d, meta)
}

func dataSourceFabricSearchServiceProfilesByName() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricSearchServiceProfilesRead,
		Schema:      readFabricServiceProfilesSearchSchema(),
	}
}

func dataSourceFabricSearchServiceProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceServiceProfilesSearchRequest(ctx, d, meta)
}
