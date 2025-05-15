// Package converters contains functions to convert between different data types and formats.
package converters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// StringArrToIfArr converts a slice of strings to a slice of any
func StringArrToIfArr(sli []string) []interface{} {
	var arr []interface{}
	for _, v := range sli {
		arr = append(arr, v)
	}
	return arr
}

// IfArrToStringArr converts a slice of any to a slice of strings
func IfArrToStringArr(ifaceArr []interface{}) []string {
	arr := make([]string, len(ifaceArr))
	for i, v := range ifaceArr {
		arr[i] = fmt.Sprint(v)
	}
	return arr
}

// IfArrToIntStringArr converts a slice of any to a slice of strings, where each element is an int
func IfArrToIntStringArr(ifaceArr []interface{}) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, strconv.Itoa(v.(int)))
	}
	return arr
}

// IfArrToIntArr converts a slice of any to a slice of ints
func IfArrToIntArr(ifaceArr []interface{}) []int {
	var arr []int
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, v.(int))
	}
	return arr
}

// ToLowerIf converts an any to a string and converts it to lowercase
func ToLowerIf(v interface{}) string {
	return strings.ToLower(v.(string))
}

// Difference returns the difference between two slices of strings.
// from https://stackoverflow.com/a/45428032
func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// ListToInt32List converts a slice of any to a slice of int32
func ListToInt32List(list []interface{}) []int32 {
	result := make([]int32, len(list))
	for i, v := range list {
		result[i] = int32(v.(int))
	}
	return result
}

// SetToStringList converts a schema.Set to a slice of strings
func SetToStringList(set *schema.Set) []string {
	list := set.List()
	return IfArrToStringArr(list)
}

// InterfaceMapToStringMap converts a string map of any values to a string map of strings
func InterfaceMapToStringMap(mapIn map[string]interface{}) map[string]string {
	mapOut := make(map[string]string)
	for k, v := range mapIn {
		mapOut[k] = fmt.Sprintf("%v", v)
	}
	return mapOut
}

// NetworkScopeArrayToStringArray converts a slice of fabricv4.NetworkScope to a slice of strings
func NetworkScopeArrayToStringArray(list []fabricv4.NetworkScope) []string {
	arr := make([]string, len(list))
	for _, v := range list {
		arr = append(arr, fmt.Sprint(v))
	}
	return arr
}
