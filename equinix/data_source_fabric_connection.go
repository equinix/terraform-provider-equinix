package equinix

import (
	"context"
	"fmt"
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
	fmt.Println(" inside fabric connection read data ")
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricConnectionRead(ctx, d, meta)
}
