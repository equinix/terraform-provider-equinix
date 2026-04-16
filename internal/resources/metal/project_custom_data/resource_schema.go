package project_custom_data

import (
	"github.com/equinix/terraform-provider-equinix/internal/deprecations"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func getResourceSchema() *schema.Schema {
	return &schema.Schema{
		DeprecationMessage: deprecations.MetalDeprecationMessage,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"project_id": schema.StringAttribute{
				Description: "The Equinix Metal project UUID to update",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"custom_data": schema.StringAttribute{
				Description: "Project custom data as a JSON object string",
				Required:    true,
			},
		},
	}
}
