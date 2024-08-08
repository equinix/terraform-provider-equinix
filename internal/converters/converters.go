package converters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func StringArrToIfArr(sli []string) []interface{} {
	var arr []interface{}
	for _, v := range sli {
		arr = append(arr, v)
	}
	return arr
}

func IfArrToStringArr(ifaceArr []interface{}) []string {
	arr := make([]string, len(ifaceArr))
	for i, v := range ifaceArr {
		arr[i] = fmt.Sprint(v)
	}
	return arr
}

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

func ToLowerIf(v interface{}) string {
	return strings.ToLower(v.(string))
}

// Difference returns the elements in `a` that aren't in `b`.
func Difference[T comparable](a, b []T) []T {
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []T
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func ListToInt32List(list []interface{}) []int32 {
	result := make([]int32, len(list))
	for i, v := range list {
		result[i] = int32(v.(int))
	}
	return result
}

func SetToStringList(set *schema.Set) []string {
	list := set.List()
	return IfArrToStringArr(list)
}

func InterfaceMapToStringMap(mapIn map[string]interface{}) map[string]string {
	mapOut := make(map[string]string)
	for k, v := range mapIn {
		mapOut[k] = fmt.Sprintf("%v", v)
	}
	return mapOut
}
