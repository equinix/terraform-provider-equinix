package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcesCloudRouterResourceSchema() map[string]*schema.Schema {
	sch := resourcesFabricCloudRouterResourceSchema()
	for key, _ := range sch {
		if key == "uuid" {
			sch[key].Required = true
			sch[key].Optional = false
			sch[key].Computed = false
		} else {
			sch[key].Required = false
			sch[key].Optional = false
			sch[key].Computed = true
		}
	}
	return sch
}

func dataSourceCloudRouter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudRouterRead,
		Schema:      dataSourcesCloudRouterResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch Fabric Cloud Router for a given UUID\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func dataSourceCloudRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceCloudRouterRead(ctx, d, meta)
}
