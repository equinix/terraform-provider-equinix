package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func IDAttributeDefaultDescription() schema.StringAttribute {
	return IDAttribute("The unique identifier of the resource")
}

func IDAttribute(description string) schema.StringAttribute {
	att := schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Description: "The unique identifier of the resource",
	}
	if description != "" {
		att.Description = description
	}
	return att
}
