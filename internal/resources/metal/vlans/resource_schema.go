package vlans

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this Metal Vlan",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "ID of parent project",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description string",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"facility": schema.StringAttribute{
				Description:        "Facility where to create the VLAN",
				DeprecationMessage: "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
				Optional:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("metro"),
					}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				// TODO: aayushrangwala to check if this is needed with the framework changes
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	suppress diff when unsetting facility
				//if len(old) > 0 && new == "" {
				//	return true
				//}
				//return old == new
				//},
			},
			"metro": schema.StringAttribute{
				Description: "Metro in which to create the VLAN",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("facility"),
					}...),
				},
				// TODO: aayushrangwala to check if this is needed with the framework changes
				//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//	_, facOk := d.GetOk("facility")
				// new - new val from template
				// old - old val from state
				//
				// suppress diff if metro is manually set for first time, and
				// facility is already set
				//if len(new) > 0 && old == "" && facOk {
				//	return facOk
				//}
				//return old == new
				//},
				// TODO: add statefunc in framework
				//StateFunc: converters.ToLowerIf,
			},
			"vxlan": schema.Int64Attribute{
				Description: "VLAN ID, must be unique in metro",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Optional: true,
				Computed: true,
			},
		},
	}
}
