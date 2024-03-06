package organizationmember

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/packethost/packngo"
)

type ResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Invitee        types.String `tfsdk:"invitee"`
	InvitedBy      types.String `tfsdk:"invited_by"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ProjectsIDs    types.Set    `tfsdk:"projects_ids"`
	Nonce          types.String `tfsdk:"nonce"`
	Message        types.String `tfsdk:"message"`
	Created        types.String `tfsdk:"created"`
	Updated        types.String `tfsdk:"updated"`
	Roles          types.Set    `tfsdk:"roles"`
	State          types.String `tfsdk:"state"`
}

func (m *ResourceModel) parse(ctx context.Context, org *packngo.Invitation) diag.Diagnostics {
	var diags diag.Diagnostics
	m.Invitee = types.StringValue(org.Invitee)
	m.InvitedBy = types.StringValue(org.InvitedBy.Href)
	m.OrganizationID = types.StringValue(org.ID)

	ProjectList, _ := types.SetValueFrom(ctx, types.StringType, org.Projects)
	m.ProjectsIDs = ProjectList

	m.Nonce = types.StringValue(org.Nonce)
	m.Created = types.StringValue(org.CreatedAt.String())
	m.Updated = types.StringValue(org.UpdatedAt.String())

	rolesList, _ := types.SetValueFrom(ctx, types.StringType, org.Roles)
	m.Roles = rolesList
	m.State = types.StringValue("active")

	m.ID = types.StringValue(org.ID)
	return diags
}
