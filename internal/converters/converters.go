package converters

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ConvertStringArr(ifaceArr []interface{}) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, v.(string))
	}
	return arr
}

func ConvertStringArrToIfArr(sli []string) []interface{} {
	var arr []interface{}
	for _, v := range sli {
		arr = append(arr, v)
	}
	return arr
}

func ConvertIntArr(ifaceArr []interface{}) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, strconv.Itoa(v.(int)))
	}
	return arr
}

func ConvertIntArr2(ifaceArr []interface{}) []int {
	var arr []int
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, v.(int))
	}
	return arr
}

func ToLower(v interface{}) string {
	return strings.ToLower(v.(string))
}

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

func ExpandListToStringList(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = fmt.Sprint(v)
	}
	return result
}

func ExpandSetToStringList(set *schema.Set) []string {
	list := set.List()
	return ExpandListToStringList(list)
}

func ExpandInterfaceMapToStringMap(mapIn map[string]interface{}) map[string]string {
	mapOut := make(map[string]string)
	for k, v := range mapIn {
		mapOut[k] = fmt.Sprintf("%v", v)
	}
	return mapOut
}
