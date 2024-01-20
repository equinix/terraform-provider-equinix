package equinix

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//func setNestedAttributesToComputed(top *schema.Schema) *schema.Schema {
//	if reflect.TypeOf(top.Elem)
//}
//

// Recursive function to traverse nested schema definitions
//
//	func traverseSchema(attr string, sch *schema.Schema) {
//		switch t := sch.Elem.(type) {
//		case *schema.Resource:
//			// Recursive call for nested resources
//			for nestedAttr, nestedSchema := range t.Schema {
//				// Recursive call for each nested attribute
//
//				traverseSchema(nestedAttr, nestedSchema)
//			}
//		case *schema.Schema:
//			// Handle nested simple types (e.g., TypeMap)
//			if t.Elem != nil {
//				// Recursive call for nested elements within TypeList or TypeSet
//				traverseSchema(attr, t)
//			}
//		}
//	}
func readFabricConnectionResourceSchema() map[string]*schema.Schema {
	sch := FabricConnectionResourceSchema()
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

func dataSourceFabricConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricConnectionRead,
		Schema:      readFabricConnectionResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch connection for a given UUID\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
	}
}

func dataSourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricConnectionRead(ctx, d, meta)
}
