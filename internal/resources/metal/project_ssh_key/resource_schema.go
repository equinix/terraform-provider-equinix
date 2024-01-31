package project_ssh_key

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	metal_ssh_key "github.com/equinix/terraform-provider-equinix/internal/resources/metal/ssh_key"
)

func GetResourceSchema() *schema.Schema {
    sch := metal_ssh_key.GetResourceSchema()
	sch.Attributes["project_id"] = schema.StringAttribute{
		Description:   "The ID of parent project",
		Required:      true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	return sch
}
