package packet

import (
	"fmt"
	"github.com/packethost/packngo"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func convertStringArr(ifaceArr []interface{}) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, v.(string))
	}
	return arr
}

func validateFacility(meta interface{}, facility string) (s []string, es []error) {
	for _, item := range packngo.Facilities {
		if item == facility {
			return
		}
	}
	client := meta.(*packngo.Client)
	facilities, _, err := client.Facilities.List(&packngo.ListOptions{})
	if err != nil {
		es = append(es, err)
	}
	for _, item := range facilities {
		if item.Code == facility {
			return
		}
	}
	es = append(es, fmt.Errorf("Could not find facility: %s", facility))
	return
}
