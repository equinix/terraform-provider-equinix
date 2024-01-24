package metal_project_ssh_key

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"project_id": schema.StringAttribute{
			Description: "The ID of parent project",
			Required:    true,
		},
		"search": schema.StringAttribute{
			Description:  "The name, fingerprint, id, or public_key of the SSH Key to search for in the Equinix Metal project",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.AtLeastOneOf(path.Expressions{
					path.MatchRoot("id"),
				}...),
			},
		},
		"id": schema.StringAttribute{
			Description: "The id of the SSH Key",
			Optional:    true,
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The label of the Equinix Metal SSH Key",
			Computed:    true,
		},
		"public_key": schema.StringAttribute{
			Description: "The public key",
			Computed:    true,
		},
		"fingerprint": schema.StringAttribute{
			Description: "The fingerprint of the SSH key",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "The timestamp for when the SSH key was created",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "The timestamp for the last time the SSH key was updated",
			Computed:    true,
		},
		"owner_id": schema.StringAttribute{
			Description: "The UUID of the Equinix Metal API User who owns this key",
			Computed:    true,
		},
	},
}
