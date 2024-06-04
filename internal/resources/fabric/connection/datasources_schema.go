package connection

import (
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func readFabricConnectionResourceSchema() map[string]*schema.Schema {
	sch := fabricConnectionResourceSchema()
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

func readFabricConnectionSchemaUpdated() map[string]*schema.Schema {
	sch := readFabricConnectionResourceSchema()
	sch["uuid"].Computed = true
	sch["uuid"].Optional = false
	sch["uuid"].Required = false
	return sch
}

func readFabricConnectionSearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Cloud Routers",
			Elem: &schema.Resource{
				Schema: readFabricConnectionSchemaUpdated(),
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
						Type:        schema.TypeString,
						Required:    true,
						Description: fmt.Sprintf("Possible field names to use on filters. One of %v", fabricv4.AllowedSearchFieldNameEnumValues),
					},
					"operator": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Operators to use on your filtered field with the values given. One of [ =, !=, >, >=, <, <=, BETWEEN, NOT BETWEEN, LIKE, NOT LIKE, IN, NOT IN, IS NOT NULL, IS NULL]",
					},
					"values": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"group": {
						Type:         schema.TypeString,
						Optional:     true,
						Description:  "Optional custom id parameter to assign this filter to an inner AND or OR group. Group id must be prefixed with AND_ or OR_. Ensure intended grouped elements have the same given id. Ungrouped filters will be placed in the filter list group by themselves.",
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(AND_|OR_)`), "Given string does not start with AND_ or OR_"),
					},
				},
			},
		},
		"outer_operator": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Determines if the filter list will be grouped by AND or by OR. One of [AND, OR]",
			ValidateFunc: validation.StringInSlice([]string{"AND", "OR"}, false),
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
						Description: fmt.Sprintf("The property name to use in sorting. One of %v. Defaults to /changeLog/updatedDateTime", fabricv4.AllowedSortByEnumValues),
					},
				},
			},
		},
	}
}
