package datalist

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	matchByStringComparison = []string{"in", "re", "substring"}
	matchByNumberComparison = []string{"less_than", "less_than_or_equal", "greater_than", "greater_than_or_equal"}
)

type commonFilter struct {
	attribute string
	values    []interface{}
	all       bool
	matchBy   string
}

func filterSchema(allowedAttributes []string) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"attribute": {
					Type:         schema.TypeString,
					Description:  "The attribute used to filter. Filter attributes are case-sensitive",
					Required:     true,
					ValidateFunc: validation.StringInSlice(allowedAttributes, false),
				},
				"values": {
					Type:        schema.TypeList,
					Description: "The filter values. Filter values are case-sensitive. If you specify multiple values for a filter, the values are joined with an OR by default, and the request returns all results that match any of the specified values",
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"all": {
					Type:        schema.TypeBool,
					Description: "If is set to true, the values are joined with an AND, and the requests returns only the results that match all specified values",
					Optional:    true,
					Default:     false,
				},
				"match_by": {
					Type:         schema.TypeString,
					Description:  "The type of comparison to apply. One of: in (default), re, substring, less_than, less_than_or_equal, greater_than, greater_than_or_equal",
					Optional:     true,
					Default:      "in",
					ValidateFunc: validation.StringInSlice(append(matchByNumberComparison, matchByStringComparison...), false),
				},
			},
		},
		Optional:    true,
		Description: "One or more attribute/values pairs on which to filter results",
	}
}

func expandFilters(recordSchema map[string]*schema.Schema, rawFilters []interface{}) ([]commonFilter, error) {
	expandedFilters := make([]commonFilter, len(rawFilters))

	for i, rawFilter := range rawFilters {
		f := rawFilter.(map[string]interface{})

		attr := f["attribute"].(string)
		s, ok := recordSchema[attr]
		if !ok {
			return nil, fmt.Errorf("field '%s' does not exist in record schema", attr)
		}

		matchBy := "in"
		if v, ok := f["match_by"].(string); ok {
			matchBy = v
		}

		expandedFilterValues, err := expandFilterValues(f["values"].([]interface{}), s, matchBy)
		if err != nil {
			return nil, err
		}

		for _, nc := range matchByNumberComparison {
			if matchBy == nc {
				if len(expandedFilterValues) != 1 {
					return nil, fmt.Errorf("field '%s' works with only one value", matchBy)
				}
			}
		}

		all := false
		if v, ok := f["all"]; ok {
			all = v.(bool)
		}

		expandedFilter := commonFilter{
			attribute: attr,
			values:    expandedFilterValues,
			all:       all,
			matchBy:   matchBy,
		}

		expandedFilters[i] = expandedFilter
	}

	return expandedFilters, nil
}

func isPrimitiveType(fieldType schema.ValueType) bool {
	switch fieldType {
	case schema.TypeString,
		schema.TypeBool,
		schema.TypeInt,
		schema.TypeFloat:
		return true
	}

	return false
}

// Expands a single filter value (which is a string) into the Go type that can actually be
// used for comparisons when filtering. This should not be called with container or
// composite types.
func expandPrimitiveFilterValue(
	filterValue string,
	fieldType schema.ValueType,
	matchBy string,
) (interface{}, error) {
	var expandedValue interface{}

	switch fieldType {
	case schema.TypeString:
		switch matchBy {
		case "in", "substring":
			expandedValue = filterValue
		case "re":
			re, err := regexp.Compile(filterValue)
			if err != nil {
				return nil, fmt.Errorf("unable to parse value as regular expression: %s: %s", filterValue, err)
			}
			expandedValue = re
		default:
			panic("unreachable")
		}

	case schema.TypeBool:
		boolValue, err := strconv.ParseBool(filterValue)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value as bool: %s: %s", filterValue, err)
		}
		expandedValue = boolValue

	case schema.TypeInt:
		intValue, err := strconv.Atoi(filterValue)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value as integer: %s: %s", filterValue, err)
		}
		expandedValue = intValue

	case schema.TypeFloat:
		floatValue, err := strconv.ParseFloat(filterValue, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse value as floating point: %s: %s", filterValue, err)
		}
		expandedValue = floatValue

	default:
		panic("unreachable")
	}

	return expandedValue, nil
}

// Takes the "raw" set of strings provided by the Terraform SDK and converts the values
// into the actual Go types used for comparisons when filtering.
func expandFilterValues(
	rawFilterValues []interface{},
	fieldSchema *schema.Schema,
	matchBy string,
) ([]interface{}, error) {
	expandedFilterValues := make([]interface{}, len(rawFilterValues))

	for i, rawFilterValue := range rawFilterValues {
		filterValue := rawFilterValue.(string)
		var expandedValue interface{}

		if isPrimitiveType(fieldSchema.Type) {
			ev, err := expandPrimitiveFilterValue(filterValue, fieldSchema.Type, matchBy)
			if err != nil {
				return nil, err
			}
			expandedValue = ev
		} else {
			if s, ok := fieldSchema.Elem.(*schema.Schema); ok {
				if isPrimitiveType(s.Type) {
					ev, err := expandPrimitiveFilterValue(filterValue, s.Type, matchBy)
					if err != nil {
						return nil, err
					}
					expandedValue = ev
				} else {
					return nil, fmt.Errorf("cannot filter on a non-primitive type")
				}
			} else {
				return nil, fmt.Errorf("cannot filter on aggregate type with non-Schema element type")
			}
		}

		expandedFilterValues[i] = expandedValue
	}

	return expandedFilterValues, nil
}

func applyFilters(recordSchema map[string]*schema.Schema, records []map[string]interface{}, filters []commonFilter) []map[string]interface{} {
	for _, f := range filters {
		// Handle multiple filters by applying them in order
		var filteredRecords []map[string]interface{}

		filterFunc := func(record map[string]interface{}) bool {
			result := f.all

			for _, filterValue := range f.values {
				thisValueMatches := valueMatches(recordSchema[f.attribute], record[f.attribute], filterValue, f.matchBy)
				if !f.all {
					result = result || thisValueMatches
				} else {
					result = result && thisValueMatches
				}
			}

			return result
		}

		for _, record := range records {
			if filterFunc(record) {
				filteredRecords = append(filteredRecords, record)
			}
		}

		records = filteredRecords
	}

	return records
}
