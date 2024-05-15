package equinix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricPortResourceSchemaUpdated() map[string]*schema.Schema {
	sch := FabricPortResourceSchema()
	sch["uuid"].Computed = true
	sch["uuid"].Optional = false
	sch["uuid"].Required = false
	return sch
}

func readGetPortsByNameQueryParamSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Query Parameter to Get Ports By Name",
		},
	}
}

func readFabricPortsResponseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Ports",
			Elem: &schema.Resource{
				Schema: readFabricPortResourceSchemaUpdated(),
			},
		},
		"filters": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "name",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: readGetPortsByNameQueryParamSch(),
			},
		},
	}
}

func dataSourceFabricPort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricPortRead,
		Schema:      FabricPortResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch port by uuid",
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
		Description: "Fabric V4 API compatible data resource that allow user to fetch port by name",
	}
}

func dataSourceFabricGetPortsByNameResponseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFabricPortGetByPortName(ctx, d, meta)
}
