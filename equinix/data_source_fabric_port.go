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
* Getting Started: https://docs.equinix.com/fabric/ports/managing-fabric-ports/
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Ports`,
	}
}

func dataSourceFabricPortRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricPortRead(ctx, d, meta)
}

func dataSourceFabricGetPortsByName() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricGetPortsByNameResponseRead,
		Schema:      readFabricPortsResponseSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch ports by name or uuid

Additional documentation:
* Getting Started: https://docs.equinix.com/fabric/ports/managing-fabric-ports/
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Ports`,
	}
}

func dataSourceFabricGetPortsByNameResponseRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return resourceFabricPortGetByPortName(ctx, d, meta)
}

func dataSourceFabricGetPortVlans() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricPortVlansRead,
		Schema:      readFabricPortVlansResponseSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch the vlans of a given port

Additional documentation:
* Getting Started: https://docs.equinix.com/fabric/ports/managing-fabric-ports/
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/Ports`,
	}
}

func dataSourceFabricPortVlansRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFabricPortVlansRead(ctx, d, meta)
}
