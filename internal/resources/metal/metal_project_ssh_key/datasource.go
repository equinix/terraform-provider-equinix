package metal_project_ssh_key

import (
	"github.com/equinix/terraform-provider-equinix/internal/config"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"fmt"
	"path"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/metal_ssh_key"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func DataSource() *schema.Resource {
	dsSchema := metal_ssh_key.CommonFieldsDataSource()
	dsSchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The ID of parent project",
		ForceNew:    true,
		Required:    true,
	}
	dataSource := &schema.Resource{
		Read:   dataSourceRead,
		Schema: dsSchema,
	}
	return dataSource
}

func dataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).Metal

	search := d.Get("search").(string)
	id := d.Get("id").(string)
	projectID := d.Get("project_id").(string)

	if id == "" && search == "" {
		return fmt.Errorf("You must supply either search or id")
	}

	var (
		key        packngo.SSHKey
		searchOpts *packngo.SearchOptions
	)

	if search != "" {
		searchOpts = &packngo.SearchOptions{Search: search}
	}
	keys, _, err := client.Projects.ListSSHKeys(projectID, searchOpts)
	if err != nil {
		err = fmt.Errorf("Error listing project ssh keys: %s", equinix_errors.FriendlyError(err))
		return err
	}

	for i := range keys {
		// use the first match for searches
		if search != "" {
			key = keys[i]
			break
		}

		// otherwise find the matching ID
		if keys[i].ID == id {
			key = keys[i]
			break
		}
	}

	if key.ID == "" {
		// Not Found
		return fmt.Errorf("Project %q SSH Key matching %q was not found", projectID, search)
	}

	ownerID := path.Base(key.Owner.Href)

	d.SetId(key.ID)
	d.Set("name", key.Label)
	d.Set("public_key", key.Key)
	d.Set("fingerprint", key.FingerPrint)
	d.Set("owner_id", ownerID)
	d.Set("created", key.Created)
	d.Set("updated", key.Updated)

	if strings.Contains(key.Owner.Href, "/projects/") {
		d.Set("project_id", ownerID)
	}

	return nil
}
