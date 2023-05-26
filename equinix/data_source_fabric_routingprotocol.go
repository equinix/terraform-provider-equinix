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
		Description: "Fabric V4 API compatible data resource that allow user to fetch routing protocol for a given UUID\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func dataSourceRoutingProtocolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	//connection_uuid, _ := d.Get("connection_uuid").(string) //fixme: is this how you set a new variable from input???
	d.SetId(uuid)
	return resourceFabricRoutingProtocolRead(ctx, d, meta)
}
