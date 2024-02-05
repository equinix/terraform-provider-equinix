package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRoutingProtocol() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoutingProtocolRead,
		Schema:      readFabricRoutingProtocolResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch routing protocol for a given UUID",
	}
}

func dataSourceRoutingProtocolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricRoutingProtocolRead(ctx, d, meta)
}
