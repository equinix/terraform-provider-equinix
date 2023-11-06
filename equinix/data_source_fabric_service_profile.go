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
		Description: "Fabric V4 API compatible data resource that allow user to fetch Service Profile by UUID filter criteria\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability",
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
		Description: "Fabric V4 API compatible data resource that allow user to fetch Service Profile by name filter criteria\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func dataSourceFabricSearchServiceProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceServiceProfilesSearchRequest(ctx, d, meta)
}
