package service_token

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceBaseSchema() map[string]*schema.Schema {
	sch := resourceSchema()
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

func dataSourceBaseSchemaUpdated() map[string]*schema.Schema {
	sch := dataSourceBaseSchema()
	sch["uuid"].Computed = true
	sch["uuid"].Optional = false
	sch["uuid"].Required = false
	return sch
}

func paginationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The page offset for the pagination request. Index of the first element. Default is 0.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of elements to be requested per page. Number must be between 1 and 100. Default is 20",
			},
			"total": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Total number of elements returned.",
			},
			"next": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL relative to the last item in the response.",
			},
			"previous": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL relative to the first item in the response.",
			},
		},
	}
}

func dataSourceSearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Service Tokens",
			Elem: &schema.Resource{
				Schema: dataSourceBaseSchemaUpdated(),
			},
		},
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Filters for the Data Source Search Request",
			MaxItems:    10,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The API response property which you want to filter your request on. Can be one of the following: \"/type\", \"/name\", \"/project/projectId\", \"/uuid\", \"/state\"",
						ValidateFunc: validation.StringInSlice([]string{"/uuid", "/state", "/name", "/project/projectId"}, true),
					},
					"operator": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Possible operators to use on the filter property. Can be one of the following: [ \"=\", \"!=\", \"[NOT] LIKE\", \"[NOT] IN\", \"ILIKE\" ]",
						ValidateFunc: validation.StringInSlice([]string{"=", "!=", "[NOT] LIKE", "[NOT] IN", "ILIKE"}, true),
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
			Elem:        paginationSchema(),
		},
	}
}
