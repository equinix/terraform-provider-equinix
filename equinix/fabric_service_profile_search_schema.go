package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func createServiceProfilesSearchExpressionSch() map[string]*schema.Schema {
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
			Description:  "Priority type- DESC, ASC"},
		"property": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"/name", "/state", "/changeLog/createdDateTime", "/changeLog/updatedDateTime"}, true),
			Description:  "Search operation sort criteria /name /state /changeLog/createdDateTime /changeLog/updatedDateTime"},
	}
}
