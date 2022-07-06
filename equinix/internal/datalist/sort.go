package datalist

import (
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var sortAttributes = []string{"asc", "desc"}

type commonSort struct {
	attribute string
	direction string
}

func sortSchema(allowedAttributes []string) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"attribute": {
					Type:         schema.TypeString,
					Description:  "The attribute used to sort the results. Sort attributes are case-sensitive",
					Required:     true,
					ValidateFunc: validation.StringInSlice(allowedAttributes, false),
				},
				"direction": {
					Type:         schema.TypeString,
					Description:  "Sort results in ascending or descending order. Strings are sorted in alphabetical order. One of: asc, desc",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(sortAttributes, false),
				},
			},
		},
		Optional:    true,
		Description: "One or more attribute/direction pairs on which to sort results. If multiple sorts are provided, they will be applied in order",
	}
}

func expandSorts(rawSorts []interface{}) []commonSort {
	expandedSorts := make([]commonSort, len(rawSorts))
	for i, rawSort := range rawSorts {
		f := rawSort.(map[string]interface{})

		expandedSort := commonSort{
			attribute: f["attribute"].(string),
			direction: f["direction"].(string),
		}

		expandedSorts[i] = expandedSort
	}
	return expandedSorts
}

func applySorts(recordSchema map[string]*schema.Schema, records []map[string]interface{}, sorts []commonSort) []map[string]interface{} {
	sort.Slice(records, func(_i, _j int) bool {
		for _, s := range sorts {
			// Handle multiple sorts by applying them in order
			i := _i
			j := _j
			if strings.EqualFold(s.direction, "desc") {
				// If the direction is desc, reverse index to compare
				i = _j
				j = _i
			}

			value1 := records[i]
			value2 := records[j]
			cmp := compareValues(recordSchema[s.attribute], value1[s.attribute], value2[s.attribute])
			if cmp != 0 {
				return cmp < 0
			}
		}

		return true
	})

	return records
}
