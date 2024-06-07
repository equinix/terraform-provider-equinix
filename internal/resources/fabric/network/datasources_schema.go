package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricNetworkResourceSchema() map[string]*schema.Schema {
	sch := fabricNetworkResourceSchema()
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

func readFabricNetworkSearchSchema() map[string]*schema.Schema {
	sch := readFabricNetworkResourceSchema()
	sch["uuid"].Required = false
	return sch
}
