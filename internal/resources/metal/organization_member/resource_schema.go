package organizationmember

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetResourceSchema() *schema.Schema {
	return &schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the organization member.",
				Computed:    true,
			},
			"invitee": schema.StringAttribute{
				Description: "The email address of the user to invite",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"invited_by": schema.StringAttribute{
				Description: "The user id of the user that sent the invitation (only known in the invitation stage)",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The organization to invite the user to",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"projects_ids": schema.SetAttribute{
				Description: "Project IDs the member has access to within the organization. If the member is an 'owner', the projects list should be empty.",
				Required:    true,
				ElementType: types.StringType,
			},
			"nonce": schema.StringAttribute{
				Description: "The nonce for the invitation (only known in the invitation stage)",
				Computed:    true,
			},
			"message": schema.StringAttribute{
				Description: "A message to the invitee (only used during the invitation stage)",
				Optional:    true,
			},
			"created": schema.StringAttribute{
				Description: "When the invitation was created (only known in the invitation stage)",
				Computed:    true,
			},
			"updated": schema.StringAttribute{
				Description: "When the invitation was updated (only known in the invitation stage)",
				Computed:    true,
			},
			"roles": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "Organization roles (owner, collaborator, limited_collaborator, billing)",
				Required:    true,
			},
			"state": schema.StringAttribute{
				Description: "The state of the membership ('invited' when an invitation is open, 'active' when the user is an organization member)",
				Computed:    true,
			},
		},
	}
}
