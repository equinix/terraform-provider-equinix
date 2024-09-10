package equinix

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricServiceProfileResourceSchema() map[string]*schema.Schema {
	sch := fabricServiceProfileSchema()
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

func readFabricServiceProfileSearchResourceSchema() map[string]*schema.Schema {
	sch := fabricServiceProfileSchema()
	for key := range sch {
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
		"and_filters": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Optional boolean flag to indicate if the filters will be AND'd together. Defaults to false",
			Default:     false,
		},
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Filters for the Data Source Search Request (If and_filters is not set to true you cannot provide more than one filter block)",
			MaxItems:    10,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:        schema.TypeString,
						Required:    true,
						Description: fmt.Sprintf("Property to apply operator and values to. One of %v", []string{"/name", "/uuid", "/state", "/metros/code", "/visibility", "/type", "/project/projectId"}),
					},
					"operator": {
						Type:        schema.TypeString,
						Required:    true,
						Description: fmt.Sprintf("Operators to use on your filtered field with the values given. One of %v", []string{"="}),
					},
					"values": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"pagination": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Pagination details for the Data Source Search Request",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"offset": {
						Type:        schema.TypeInt,
						Optional:    true,
						Default:     0,
						Description: "The page offset for the pagination request. Index of the first element. Default is 0.",
					},
					"limit": {
						Type:        schema.TypeInt,
						Optional:    true,
						Default:     20,
						Description: "Number of elements to be requested per page. Number must be between 1 and 100. Default is 20",
					},
				},
			},
		},
		"sort": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Filters for the Data Source Search Request",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"direction": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "DESC",
						Description:  "The sorting direction. Can be one of: [DESC, ASC], Defaults to DESC",
						ValidateFunc: validation.StringInSlice([]string{"DESC", "ASC"}, true),
					},
					"property": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "/changeLog/updatedDateTime",
						Description: fmt.Sprintf("The property name to use in sorting. One of %v. Defaults to /changeLog/updatedDateTime", []string{"/name", "/uuid", "/state", "/location/metroCode", "/location/metroName", "/package/code", "/changeLog/createdDateTime", "/changeLog/updatedDateTime"}),
					},
				},
			},
		},
	}
}

func dataSourceFabricServiceProfileReadByUuid() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricServiceProfileRead,
		Schema:      readFabricServiceProfileResourceSchema(),
		Description: `Fabric V4 API compatible data resource that allow user to fetch Service Profile by UUID filter criteria

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/IMPLEMENTATION/fabric-Sprofiles-implement.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#service-profiles`,
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
		Description: `Fabric V4 API compatible data resource that allow user to fetch Service Profile by name filter criteria

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/Fabric/IMPLEMENTATION/fabric-Sprofiles-implement.htm
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#service-profiles`,
	}
}

func dataSourceFabricSearchServiceProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceServiceProfilesSearchRequest(ctx, d, meta)
}
