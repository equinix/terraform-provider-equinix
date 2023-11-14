package equinix

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func resourceMetalOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalOrganizationMemberCreate,
		Read:   resourceMetalOrganizationMemberRead,
		Delete: resourceMetalOrganizationMemberDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), ":")
				invitee := parts[0]
				orgID := parts[1]
				d.SetId(d.Id())
				d.Set("invitee", invitee)
				d.Set("organization_id", orgID)
				if err := resourceMetalOrganizationMemberRead(d, meta); err != nil {
					return nil, err
				}
				if d.Id() == "" {
					return nil, fmt.Errorf("Member %s does not exist in organization %s.", invitee, orgID)
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"invitee": {
				Type:         schema.TypeString,
				Description:  "The email address of the user to invite",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"invited_by": {
				Type:        schema.TypeString,
				Description: "The user id of the user that sent the invitation (only known in the invitation stage)",
				Computed:    true,
			},
			"organization_id": {
				Type:         schema.TypeString,
				Description:  "The organization to invite the user to",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"projects_ids": {
				Type:        schema.TypeSet,
				Description: "Project IDs the member has access to within the organization. If the member is an 'owner', the projects list should be empty.",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				// TODO: Update should be supported. packngo.InvitationService does not offer an Update	method.
				ForceNew: true,
			},
			"nonce": {
				Type:        schema.TypeString,
				Description: "The nonce for the invitation (only known in the invitation stage)",
				Computed:    true,
			},
			"message": {
				Type:        schema.TypeString,
				Description: "A message to the invitee (only used during the invitation stage)",
				Optional:    true,
				ForceNew:    true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "When the invitation was created (only known in the invitation stage)",
				Computed:    true,
			},
			"updated": {
				Type:        schema.TypeString,
				Description: "When the invitation was updated (only known in the invitation stage)",
				Computed:    true,
			},
			"roles": {
				Type:        schema.TypeSet,
				Description: "Organization roles (owner, collaborator, limited_collaborator, billing)",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				// TODO: Update should be supported. packngo.InvitationService does not offer an Update	method.
				ForceNew: true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "The state of the membership ('invited' when an invitation is open, 'active' when the user is an organization member)",
				Computed:    true,
			},
		},
	}
}

func resourceMetalOrganizationMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).Metal

	email := d.Get("invitee").(string)
	createRequest := &packngo.InvitationCreateRequest{
		Invitee:     email,
		Roles:       convertStringArr(d.Get("roles").(*schema.Set).List()),
		ProjectsIDs: convertStringArr(d.Get("projects_ids").(*schema.Set).List()),
		Message:     strings.TrimSpace(d.Get("message").(string)),
	}

	orgID := d.Get("organization_id").(string)
	_, _, err := client.Invitations.Create(orgID, createRequest, nil)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(fmt.Sprintf("%s:%s", email, orgID))

	return resourceMetalOrganizationMemberRead(d, meta)
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

func resourceMetalOrganizationMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).Metal
	parts := strings.Split(d.Id(), ":")
	invitee := parts[0]
	orgID := parts[1]

	listOpts := &packngo.ListOptions{Includes: []string{"user"}}
	invitations, _, err := client.Invitations.List(orgID, listOpts)
	if err != nil {
		err = friendlyError(err)
		// If the org was destroyed, mark as gone.
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	members, _, err := client.Members.List(orgID, &packngo.GetOptions{Includes: []string{"user"}})
	if err != nil {
		err = friendlyError(err)
		// If the org was destroyed, mark as gone.
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	member, err := findMember(invitee, members, invitations)
	if !d.IsNewResource() && err != nil {
		log.Printf("[WARN] Could not find member %s in organization, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if member.isMember() {
		projectIDs := []string{}
		for _, project := range member.Member.Projects {
			projectIDs = append(projectIDs, path.Base(project.URL))
		}
		return setMap(d, map[string]interface{}{
			"state":           "active",
			"roles":           stringArrToIfArr(member.Member.Roles),
			"projects_ids":    stringArrToIfArr(projectIDs),
			"organization_id": path.Base(member.Member.Organization.URL),
		})
	} else if member.isInvitation() {
		projectIDs := []string{}
		for _, project := range member.Invitation.Projects {
			projectIDs = append(projectIDs, path.Base(project.Href))
		}
		return setMap(d, map[string]interface{}{
			"state":           "invited",
			"organization_id": path.Base(member.Invitation.Organization.Href),
			"roles":           member.Invitation.Roles,
			"projects_ids":    projectIDs,
			"created":         member.Invitation.CreatedAt.String(),
			"updated":         member.Invitation.UpdatedAt.String(),
			"nonce":           member.Invitation.Nonce,
			"invited_by":      path.Base(member.Invitation.InvitedBy.Href),
		})
	}
	return fmt.Errorf("got an invalid member object")
}

func resourceMetalOrganizationMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).Metal

	listOpts := &packngo.ListOptions{Includes: []string{"user"}}
	invitations, _, err := client.Invitations.List(d.Get("organization_id").(string), listOpts)
	if err != nil {
		err = friendlyError(err)
		// If the org was destroyed, mark as gone.
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	orgID := d.Get("organization_id").(string)
	org, _, err := client.Organizations.Get(orgID, &packngo.GetOptions{Includes: []string{"members", "members.user"}})
	if err != nil {
		err = friendlyError(err)
		// If the org was destroyed, mark as gone.
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	member, err := findMember(d.Get("invitee").(string), org.Members, invitations)
	if err != nil {
		d.SetId("")
		return nil
	}

	if member.isMember() {
		_, err = client.Members.Delete(orgID, member.Member.ID)
		if err != nil {
			err = friendlyError(err)
			// If the member was deleted, mark as gone.
			if isNotFound(err) {
				d.SetId("")
				return nil
			}
			return err
		}
	} else if member.isInvitation() {
		_, err = client.Invitations.Delete(member.Invitation.ID)
		if err != nil {
			err = friendlyError(err)
			// If the invitation was deleted, mark as gone.
			if isNotFound(err) {
				d.SetId("")
				return nil
			}
			return err
		}
	}
	return nil
}
