package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricCloudRouterResourceSchema() map[string]*schema.Schema {
	sch := fabricCloudRouterResourceSchema()
	for key := range sch {
		if key == "uuid" {
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

func dataSourceFabricCloudRouter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricCloudRouterRead,
		Schema:      readFabricCloudRouterResourceSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch Fabric Cloud Router for a given UUID

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-intro.htm#HowItWorks
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#fabric-cloud-routers`,
	}
}

func dataSourceFabricCloudRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricCloudRouterRead(ctx, d, meta)
}
