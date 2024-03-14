package organizationmember

import (
	"context"
	"fmt"
	"log"
	"strings"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/packethost/packngo"
)

type member struct {
	*packngo.Member
	*packngo.Invitation
}

func (m *member) isMember() bool {
	return m.Member != nil
}

func (m *member) isInvitation() bool {
	return m.Invitation != nil
}

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name:   "equinix_metal_organization_member",
				Schema: GetResourceSchema(),
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "importer Organization")

	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		return

	}
	invitee := parts[0]
	orgID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tfpath.Root("invitee"), invitee)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tfpath.Root("organization_id"), orgID)...)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invite := plan.Invitee.ValueString()

	roles := make([]string, 0)
	for _, elem := range plan.Roles.Elements() {
		if strValue, ok := elem.(types.String); ok {

			if !strValue.IsNull() {
				roles = append(roles, strValue.ValueString())
			}
		}
	}

	projects := make([]string, 0)
	for _, elem := range plan.ProjectsIDs.Elements() {
		if strValue, ok := elem.(types.String); ok {
			projects = append(projects, strValue.ValueString())
		}
	}

	createRequest := &packngo.InvitationCreateRequest{
		Invitee:     invite,
		Roles:       roles,
		ProjectsIDs: projects,
		Message:     strings.TrimSpace(plan.Message.ValueString()),
	}

	orgID := plan.OrganizationID.ValueString()
	_, _, err := client.Invitations.Create(orgID, createRequest, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create invitation",
			err.Error(),
		)
		return
	}

	listOpts := &packngo.ListOptions{Includes: []string{"user"}}
	invitations, _, err := client.Invitations.List(orgID, listOpts)
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		// If the org was destroyed, mark as gone.
		if equinix_errors.IsNotFound(err) {
			plan.OrganizationID = basetypes.NewStringNull()
			return
		}
		resp.Diagnostics.AddError(
			"Failed to list invitations",
			err.Error(),
		)
		return
	}

	members, _, err := client.Members.List(orgID, listOpts)
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		// If the org was destroyed, mark as gone.
		if equinix_errors.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Failed to List members",
			err.Error(),
		)
		return
	}

	member, err := findMember(invite, members, invitations)
	if err != nil {
		log.Printf("[WARN] Could not find member %s in organization, removing from state", plan.OrganizationID)
		plan.OrganizationID = basetypes.NewStringNull()
		resp.Diagnostics.AddError(
			"Failed to find members",
			err.Error(),
		)
		return
	}

	// Parse API response into the Terraform state
	if member != nil {
		diags := plan.parse(ctx, member)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read Organization")
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve values from plan
	var data ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invitee := data.Invitee.ValueString()
	orgID := data.OrganizationID.ValueString()

	listOpts := &packngo.ListOptions{Includes: []string{"user"}}
	invitations, _, err := client.Invitations.List(orgID, listOpts)
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		// If the org was destroyed, mark as gone.
		if equinix_errors.IsNotFound(err) {
			data.OrganizationID = basetypes.NewStringNull()
			return
		}
		resp.Diagnostics.AddError(
			"Failed to list invitations",
			err.Error(),
		)
		return
	}

	members, _, err := client.Members.List(orgID, &packngo.GetOptions{Includes: []string{"user"}})
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		// If the org was destroyed, mark as gone.
		if equinix_errors.IsNotFound(err) {
			data.OrganizationID = basetypes.NewStringNull()
			return
		}
		resp.Diagnostics.AddError(
			"Failed to list members",
			err.Error(),
		)
		return
	}
	member, err := findMember(invitee, members, invitations)
	if err != nil {
		log.Printf("[WARN] Could not find member %s in organization, removing from state", data.OrganizationID)
		data.OrganizationID = basetypes.NewStringNull()
		return
	}

	diags := data.parse(ctx, member)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete Organization")
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	var data ResourceModel

	listOpts := &packngo.ListOptions{Includes: []string{"user"}}
	invitations, _, err := client.Invitations.List(data.OrganizationID.ValueString(), listOpts)
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		// If the org was destroyed, mark as gone.
		if equinix_errors.IsNotFound(err) {
			data.OrganizationID = types.StringNull()
			return
		}
		resp.Diagnostics.AddError(
			"Failed to list invitations",
			err.Error(),
		)
		return
	}

	org, _, err := client.Organizations.Get(data.OrganizationID.ValueString(), &packngo.GetOptions{Includes: []string{"members", "members.user"}})
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		// If the org was destroyed, mark as gone.
		if equinix_errors.IsNotFound(err) {
			data.OrganizationID = types.StringNull()
			return
		}

		resp.Diagnostics.AddError(
			"Failed to get Organizations",
			err.Error(),
		)
		return
	}

	member, err := findMember(data.Invitee.ValueString(), org.Members, invitations)
	if err != nil {
		data.OrganizationID = types.StringNull()
		return
	}

	if member.isMember() {
		_, err = client.Members.Delete(data.OrganizationID.ValueString(), member.Member.ID)
		if err != nil {
			err = equinix_errors.FriendlyError(err)
			// If the member was deleted, mark as gone.
			if equinix_errors.IsNotFound(err) {
				data.OrganizationID = types.StringNull()
				return
			}
			resp.Diagnostics.AddError(
				"Failed to delete member",
				err.Error(),
			)
			return
		}
	} else if member.isInvitation() {
		_, err = client.Invitations.Delete(member.Invitation.ID)
		if err != nil {
			err = equinix_errors.FriendlyError(err)
			// If the invitation was deleted, mark as gone.
			if equinix_errors.IsNotFound(err) {
				data.OrganizationID = types.StringNull()
				return
			}

			resp.Diagnostics.AddError(
				"Failed to delete invitation",
				err.Error(),
			)
			return
		}
	}
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// This resource does not support updates
}

func findMember(invitee string, members []packngo.Member, invitations []packngo.Invitation) (*member, error) {
	for _, mbr := range members {
		if mbr.User.Email == invitee {
			return &member{Member: &mbr}, nil
		}
	}

	for _, inv := range invitations {
		if inv.Invitee == invitee {
			return &member{Invitation: &inv}, nil
		}
	}
	return nil, fmt.Errorf("member not found")
}
