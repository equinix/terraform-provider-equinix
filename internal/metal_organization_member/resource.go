package metal_organization_member

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"github.com/equinix/terraform-provider-equinix/internal/helper"
)

type OrganizationMemberOrInvite struct {
	*packngo.Member
	*packngo.Invitation
}

func (m *OrganizationMemberOrInvite) isMember() bool {
	return m.Member != nil
}

func (m *OrganizationMemberOrInvite) isInvitation() bool {
	return m.Invitation != nil
}

type OrganizationMemberResourceModel struct {
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

func findMember(invitee string, members []packngo.Member, invitations []packngo.Invitation) (*OrganizationMemberOrInvite, error) {
	for _, mbr := range members {
		if mbr.User.Email == invitee {
			return &OrganizationMemberOrInvite{Member: &mbr}, nil
		}
	}

	for _, inv := range invitations {
		if inv.Invitee == invitee {
			return &OrganizationMemberOrInvite{Invitation: &inv}, nil
		}
	}
	return nil, fmt.Errorf("member not found")
}

func (rm *OrganizationMemberResourceModel) parse(ctx context.Context, m *OrganizationMemberOrInvite) diag.Diagnostics {
    var diags diag.Diagnostics

    if m.isMember() {
        // Parse member data
        rm.Invitee = types.StringValue(m.Member.User.Email)
        rm.OrganizationID = types.StringValue(path.Base(m.Member.Organization.URL))
        rm.State = types.StringValue("active")

        memberProjects := make([]string, len(m.Member.Projects))
        for i, project := range m.Member.Projects {
            memberProjects[i] = path.Base(project.URL)
        }

        projectIDs, diags := types.SetValueFrom(ctx, types.StringType, memberProjects)
        if diags.HasError() {
            return diags
        }
        rm.ProjectsIDs = projectIDs

        roles, diags := types.SetValueFrom(ctx, types.StringType, m.Member.Roles)
        if diags.HasError() {
            return diags
        }
        rm.Roles = roles
    } else if m.isInvitation() {
        // Parse invitation data
        rm.Invitee = types.StringValue(m.Invitation.Invitee)
        rm.OrganizationID = types.StringValue(path.Base(m.Invitation.Organization.Href))
        rm.State = types.StringValue("invited")
        rm.Created = types.StringValue(m.Invitation.CreatedAt.String())
        rm.Updated = types.StringValue(m.Invitation.UpdatedAt.String())
        rm.Nonce = types.StringValue(m.Invitation.Nonce)
        rm.InvitedBy = types.StringValue(path.Base(m.Invitation.InvitedBy.Href))

        projectIDs, diags := types.SetValueFrom(ctx, types.StringType, m.Invitation.Projects)
        if diags.HasError() {
            return diags
        }
        rm.ProjectsIDs = projectIDs

        roles, diags := types.SetValueFrom(ctx, types.StringType, m.Invitation.Roles)
        if diags.HasError() {
            return diags
        }
        rm.Roles = roles
    }

    // Construct the ID for the resource after rm.Invitee and rm.OrganizationID are updated
    id := fmt.Sprintf("%s:%s", rm.Invitee.ValueString(), rm.OrganizationID.ValueString())
    rm.ID = types.StringValue(id)

    return diags
}


func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_organization_member",
				Schema: &organizationMemberResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Create an instance of your resource model to hold the planned state
    var plan OrganizationMemberResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Prepare the data for the API request
    var roles []string
	resp.Diagnostics.Append(plan.Roles.ElementsAs(ctx, &roles, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
    
    var projectIDs []string
	resp.Diagnostics.Append(plan.ProjectsIDs.ElementsAs(ctx, &projectIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

    createRequest := &packngo.InvitationCreateRequest{
        Invitee:     plan.Invitee.ValueString(),
        Message:     strings.TrimSpace(plan.Message.ValueString()),
        Roles:       roles,
        ProjectsIDs: projectIDs,
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // API call to create the organization member or send an invitation
    invitation, _, err := client.Invitations.Create(plan.OrganizationID.ValueString(), createRequest, nil)
    if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error creating Organization Member",
            "Could not create organization member or send invitation: " + err.Error(),
        )
        return
    }

    // Invitation object wrapped in a member type required by the parse function
    m := &OrganizationMemberOrInvite{Invitation: invitation}

    // Parse API response into the Terraform state
    stateDiags := plan.parse(ctx, m)
    resp.Diagnostics.Append(stateDiags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set the state
    diags = resp.State.Set(ctx, &plan)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Retrieve the current state
    var state OrganizationMemberResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the invitee email and organization ID from the state
    parts := strings.Split(state.ID.ValueString(), ":")
    if len(parts) != 2 {
        resp.Diagnostics.AddError(
            "Invalid ID format",
            "Expected ID format is 'invitee:organizationID'. Got: " + state.ID.ValueString(),
        )
        return
    }
    invitee := parts[0]
    orgID := parts[1]

    // API calls to get the current state of the organization member
    invitations, _, err := client.Invitations.List(orgID, &packngo.ListOptions{Includes: []string{"user"}})
    if err != nil {
        err = helper.FriendlyError(err)
        // If the org was destroyed, mark as gone
		if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Organization Member",
				fmt.Sprintf("[WARN] Organization (%s) not found, removing Organization Member from state", orgID),
			)
			resp.State.RemoveResource(ctx)
			return
		}
        resp.Diagnostics.AddError(
            "Error reading Organization Invitations",
            "Could not read invitations for organization: " + err.Error(),
        )
        return
    }

    members, _, err := client.Members.List(orgID, &packngo.GetOptions{Includes: []string{"user"}})
    if err != nil {
        err = helper.FriendlyError(err)
        // If the org was destroyed, mark as gone
		if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Organization Member",
				fmt.Sprintf("[WARN] Organization (%s) not found, removing Organization Member from state", orgID),
			)
			resp.State.RemoveResource(ctx)
			return
		}
        resp.Diagnostics.AddError(
            "Error reading Organization Members",
            "Could not read members for organization: " + err.Error(),
        )
        return
    }

    member, err := findMember(invitee, members, invitations)
    // TODO (ocobles) we used to check here with legacy SDKv2
    // if !d.IsNewResource() && err != nil
    // to find out if the Read function was called during import
    // Now we can use Private state but not sure how to
    // https://github.com/hashicorp/terraform-plugin-sdk/issues/1005#issuecomment-1623695760
    if err != nil {
        resp.Diagnostics.AddError(
            "Error finding organization member",
            "Could not find member or invitation: " + err.Error(),
        )
        return
    }
    
    // Parse the API response into the Terraform state
    parseDiags := state.parse(ctx, member)
    resp.Diagnostics.Append(parseDiags...)
    if parseDiags.HasError() {
        return
    }

    // Update the Terraform state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Update(
    ctx context.Context,
    req resource.UpdateRequest,
    resp *resource.UpdateResponse,
) {
	// This resource does not support updates
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve the current state
    var state OrganizationMemberResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Fetch current members and invitations
    invitations, _, err := client.Invitations.List(state.OrganizationID.String(), &packngo.ListOptions{Includes: []string{"user"}})
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading organization invitations",
            "Could not read invitations for organization: " + err.Error(),
        )
        return
    }

    members, _, err := client.Members.List(state.OrganizationID.String(), &packngo.GetOptions{Includes: []string{"user"}})
    if err != nil {
        resp.Diagnostics.AddError(
            "Error reading organization members",
            "Could not read members for organization: " + err.Error(),
        )
        return
    }

    // Find the member or invitation to delete
    member, err := findMember(state.Invitee.String(), members, invitations)
    if err != nil {
        // If member or invitation is not found, it's already gone
        return
    }

    // Delete the member or invitation
    if member.isMember() {
        _, err = client.Members.Delete(state.OrganizationID.String(), member.Member.ID)
    } else if member.isInvitation() {
        _, err = client.Invitations.Delete(member.Invitation.ID)
    }

    if err != nil {
        err = helper.FriendlyError(err)
        // If the member/invitation was deleted, mark as gone.
        if helper.IsNotFound(err) {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error deleting organization member",
            "Could not delete member or invitation: " + err.Error(),
        )
        return
    }
}

var organizationMemberResourceSchema = schema.Schema{
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
            Description: "Organization roles (owner, collaborator, limited_collaborator, billing)",
            Required:    true,
            ElementType: types.StringType,
        },
        "state": schema.StringAttribute{
            Description: "The state of the membership ('invited' when an invitation is open, 'active' when the user is an organization member)",
            Computed:    true,
        },
    },
}
