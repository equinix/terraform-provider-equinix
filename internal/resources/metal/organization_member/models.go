package organizationmember

import (
	"context"
	"path"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// func (m *ResourceModel) parse(ctx context.Context, invitee string, members []packngo.Member, invitations []packngo.Invitation) diag.Diagnostics {
func (m *ResourceModel) parse(ctx context.Context, member *member) diag.Diagnostics {
	var diags diag.Diagnostics

	if member.isMember() {
		projectsList, diag := types.SetValueFrom(ctx, types.StringType, member.Member.Projects)
		if diag.HasError() {
			return diag
		}
		m.ProjectsIDs = projectsList
		m.State = types.StringValue("active")

		rolesList, diag := types.SetValueFrom(ctx, types.StringType, member.Member.Roles)
		if diag.HasError() {
			return diag
		}
		m.Roles = rolesList
		m.OrganizationID = types.StringValue(member.Member.Organization.URL)

	} else if member.isInvitation() {

		projectsList, diag := types.SetValueFrom(ctx, types.StringType, member.Invitation.Projects)
		if diag.HasError() {
			return diag
		}
		m.ProjectsIDs = projectsList

		m.State = types.StringValue("invited")

		rolesList, diag := types.SetValueFrom(ctx, types.StringType, member.Invitation.Roles)
		if diag.HasError() {
			return diag
		}
		m.Roles = rolesList

		//m.OrganizationID = types.StringValue(member.Invitation.Organization.Href)
		m.OrganizationID = types.StringValue(path.Base(member.Invitation.Organization.Href))
		m.Created = types.StringValue(member.Invitation.CreatedAt.String())
		m.Updated = types.StringValue(member.Invitation.UpdatedAt.String())
		m.Nonce = types.StringValue(member.Invitation.Nonce)

		//m.InvitedBy = types.StringValue(member.Invitation.InvitedBy.Href)
		m.InvitedBy = types.StringValue(path.Base(member.Invitation.InvitedBy.Href))
		m.ID = types.StringValue(member.Invitation.ID)
	}
	return diags
}
