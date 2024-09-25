package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFabricPort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricPortRead,
		Schema:      FabricPortResourceSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch port by uuid

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/IMPLEMENTATION/fabric-ports-implement.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#ports`,
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
		Description: `Fabric V4 API compatible data resource that allow user to fetch port by name

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/IMPLEMENTATION/fabric-ports-implement.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#ports`,
	}
}

func dataSourceFabricGetPortsByNameResponseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFabricPortGetByPortName(ctx, d, meta)
}
