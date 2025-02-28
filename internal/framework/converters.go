package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringSliceToAttrValue(x []string) []attr.Value {
	out := make([]attr.Value, len(x))
	for i, v := range x {
		out[i] = types.StringValue(v)
	}
	return out
}
