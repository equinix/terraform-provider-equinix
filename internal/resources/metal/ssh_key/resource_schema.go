package ssh_key

import (
	"github.com/equinix/terraform-provider-equinix/internal/deprecations"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func GetResourceSchema() *schema.Schema {
	sch := GetCommonFieldsSchema()
	sch.Attributes["name"] = schema.StringAttribute{
		Description: "The name of the SSH key for identification",
		Required:    true,
	}
	sch.Attributes["public_key"] = schema.StringAttribute{
		Description: "The public key",
		Required:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	return sch
}

func GetCommonFieldsSchema() *schema.Schema {
	return &schema.Schema{
		DeprecationMessage: deprecations.MetalDeprecationMessage,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"fingerprint": schema.StringAttribute{
				Description: "The fingerprint of the SSH key",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Description: "The timestamp for when the SSH key was created",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated": schema.StringAttribute{
				Description: "The timestamp for the last time the SSH key was updated",
				Computed:    true,
			},
			"owner_id": schema.StringAttribute{
				Description: "The UUID of the Equinix Metal API User who owns this key",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
