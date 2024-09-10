package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricRoutingProtocolResourceSchema() map[string]*schema.Schema {
	sch := createFabricRoutingProtocolResourceSchema()
	for key := range sch {
		if key == "uuid" || key == "connection_uuid" {
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

func dataSourceRoutingProtocol() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoutingProtocolRead,
		Schema:      readFabricRoutingProtocolResourceSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch routing protocol for a given UUID

API documentation can be found here - https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#routing-protocols

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/connections/FCR-connect-azureQC.htm#ConfigureRoutingDetailsintheFabricPortal
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#routing-protocols`,
	}
}

func dataSourceRoutingProtocolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricRoutingProtocolRead(ctx, d, meta)
}
