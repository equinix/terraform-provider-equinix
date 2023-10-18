package metal_ssh_key

import (
	"log"
	"path"
	"strings"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func MetalSSHKeyCommonFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the SSH key for identification",
			Required:    true,
		},

		"public_key": {
			Type:        schema.TypeString,
			Description: "The public key. If this is a file, it",
			Required:    true,
			ForceNew:    true,
		},
		"fingerprint": {
			Type:        schema.TypeString,
			Description: "The fingerprint of the SSH key",
			Computed:    true,
		},

		"created": {
			Type:        schema.TypeString,
			Description: "The timestamp for when the SSH key was created",
			Computed:    true,
		},

		"updated": {
			Type:        schema.TypeString,
			Description: "The timestamp for the last time the SSH key was updated",
			Computed:    true,
		},
		"owner_id": {
			Type:        schema.TypeString,
			Description: "The UUID of the Equinix Metal API User who owns this key",
			Computed:    true,
		},
	}
}

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: ResourceMetalSSHKeyCreate,
		Read:   ResourceMetalSSHKeyRead,
		Update: ResourceMetalSSHKeyUpdate,
		Delete: ResourceMetalSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: MetalSSHKeyCommonFields(),
	}
}

func ResourceMetalSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	createRequest := &packngo.SSHKeyCreateRequest{
		Label: d.Get("name").(string),
		Key:   d.Get("public_key").(string),
	}

	projectID, isProjectKey := d.GetOk("project_id")
	if isProjectKey {
		createRequest.ProjectID = projectID.(string)
	}

	key, _, err := client.SSHKeys.Create(createRequest)
	if err != nil {
		return equinix_errors.FriendlyError(err)
	}

	d.SetId(key.ID)

	return ResourceMetalSSHKeyRead(d, meta)
}

func ResourceMetalSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	key, _, err := client.SSHKeys.Get(d.Id(), nil)
	if err != nil {
		err = equinix_errors.FriendlyError(err)

		// If the key is somehow already destroyed, mark as
		// succesfully gone
		if equinix_errors.IsNotFound(err) {
			log.Printf("[WARN] SSHKey (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
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

func ResourceMetalSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	updateRequest := &packngo.SSHKeyUpdateRequest{}

	if d.HasChange("name") {
		kName := d.Get("name").(string)
		updateRequest.Label = &kName
	}

	if d.HasChange("public_key") {
		kKey := d.Get("public_key").(string)
		updateRequest.Key = &kKey
	}

	_, _, err := client.SSHKeys.Update(d.Id(), updateRequest)
	if err != nil {
		return equinix_errors.FriendlyError(err)
	}

	return ResourceMetalSSHKeyRead(d, meta)
}

func ResourceMetalSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	resp, err := client.SSHKeys.Delete(d.Id())
	if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
		return equinix_errors.FriendlyError(err)
	}

	d.SetId("")
	return nil
}
