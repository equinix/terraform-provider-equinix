package equinix

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricServiceProfileResourceSchema() map[string]*schema.Schema {
	sch := fabricServiceProfileSchema()
	for key, _ := range sch {
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

func readFabricServiceProfileSearchResourceSchema() map[string]*schema.Schema {
	sch := fabricServiceProfileSchema()
	for key, _ := range sch {
		sch[key].Required = false
		sch[key].Optional = false
		sch[key].Computed = true
		sch[key].MaxItems = 0
		sch[key].ValidateFunc = nil
	}
	return sch
}

func readFabricServiceProfilesSearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Service Profiles",
			Elem: &schema.Resource{
				Schema: readFabricServiceProfileSearchResourceSchema(),
			},
		},
		"view_point": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "flips view between buyer and seller representation. Available values : aSide, zSide. Default value : aSide",
		},
		"filter": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Service Profile Search Filter",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createServiceProfilesSearchFilterSch(),
			},
		},
		"sort": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Service Profile Sort criteria for Search Request response payload",
			Elem: &schema.Resource{
				Schema: createServiceProfilesSearchSortCriteriaSch(),
			},
		},
	}
}

func createServiceProfilesSearchFilterSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"property": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Search Criteria for Service Profile - /name, /uuid, /state, /metros/code, /visibility, /type",
		},
		"operator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Possible operator to use on filters = - equal",
		},
		"values": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Values",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func createServiceProfilesSearchSortCriteriaSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"direction": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"DESC", "ASC"}, true),
			Description:  "Priority type- DESC, ASC",
		},
		"property": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"/name", "/state", "/changeLog/createdDateTime", "/changeLog/updatedDateTime"}, true),
			Description:  "Search operation sort criteria /name /state /changeLog/createdDateTime /changeLog/updatedDateTime",
		},
	}
}

func dataSourceFabricServiceProfileReadByUuid() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricServiceProfileRead,
		Schema:      readFabricServiceProfileResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch Service Profile by UUID filter criteria",
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
		Description: "Fabric V4 API compatible data resource that allow user to fetch Service Profile by name filter criteria",
	}
}

func dataSourceFabricSearchServiceProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceServiceProfilesSearchRequest(ctx, d, meta)
}
