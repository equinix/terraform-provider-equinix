package equinix

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricCloudRouterResourceSchema() map[string]*schema.Schema {
	sch := fabricCloudRouterResourceSchema()
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

func readFabricCloudRouterResourceSchemaUpdated() map[string]*schema.Schema {
	sch := readFabricCloudRouterResourceSchema()
	sch["uuid"].Computed = true
	sch["uuid"].Optional = false
	sch["uuid"].Required = false
	return sch
}

func readFabricCloudRouterSearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Cloud Routers",
			Elem: &schema.Resource{
				Schema: readFabricCloudRouterResourceSchemaUpdated(),
			},
		},
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Filters for the Data Source Search Request. Maximum of 8 total filters.",
			MaxItems:    10,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The API response property which you want to filter your request on. Can be one of the following: \"/project/projectId\", \"/name\", \"/uuid\", \"/state\", \"/location/metroCode\", \"/location/metroName\", \"/package/code\", \"/*\"",
						ValidateFunc: validation.StringInSlice([]string{"/project/projectId", "/name", "/uuid", "/state", "/location/metroCode", "/location/metroName", "/package/code", "/*"}, true),
					},
					"operator": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Possible operators to use on the filter property. Can be one of the following: = - equal\n!= - not equal\n> - greater than\n>= - greater than or equal to\n< - less than\n<= - less than or equal to\n[NOT] BETWEEN - (not) between\n[NOT] LIKE - (not) like\n[NOT] IN - (not) in",
						ValidateFunc: validation.StringInSlice([]string{"=", "!=", ">", ">=", "<", "<=", "[NOT] BETWEEN", "[NOT] LIKE", "[NOT] IN"}, true),
					},
					"values": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"or": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Boolean flag indicating whether this filter is included in the OR group. There can only be one OR group and it can have a maximum of 3 filters. The OR group only counts as 1 of the 8 possible filters",
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
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "/changeLog/updatedDateTime",
						Description:  "The property name to use in sorting. Can be one of the following: [/name, /uuid, /state, /location/metroCode, /location/metroName, /package/code, /changeLog/createdDateTime, /changeLog/updatedDateTime], Defaults to /changeLog/updatedDateTime",
						ValidateFunc: validation.StringInSlice([]string{"/name", "/uuid", "/state", "/location/metroCode", "/location/metroName", "/package/code", "/changeLog/createdDateTime", "/changeLog/updatedDateTime"}, true),
					},
				},
			},
		},
	}
}

func dataSourceFabricCloudRouter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricCloudRouterRead,
		Schema:      readFabricCloudRouterResourceSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch Fabric Cloud Router for a given UUID",
	}
}

func dataSourceFabricCloudRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	uuid, _ := d.Get("uuid").(string)
	d.SetId(uuid)
	return resourceFabricCloudRouterRead(ctx, d, meta)
}

func dataSourceFabricGetCloudRouters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFabricGetCloudRoutersRead,
		Schema:      readFabricCloudRouterSearchSchema(),
		Description: "Fabric V4 API compatible data resource that allow user to fetch port by name",
	}
}

func dataSourceFabricGetCloudRoutersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFabricCloudRoutersSearch(ctx, d, meta)
}
