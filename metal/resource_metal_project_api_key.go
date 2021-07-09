package metal

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func metalProjectAPIKey() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"read_only": {
			Type:     schema.TypeBool,
			Default:  true,
			ForceNew: true,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func resourceMetalProjectAPIKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalProjectAPIKeyCreate,
		Read:   resourceMetalProjectAPIKeyRead,
		Update: resourceMetalProjectAPIKeyUpdate,
		Delete: resourceMetalProjectAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: metalProjectAPIKey(),
	}
}

func resourceMetalProjectAPIKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	createRequest := &packngo.APIKeyCreateRequest{
		ProjectID:   d.Get("project_id").(string),
		ReadOnly:    d.Get("read_only").(bool),
		Description: d.Get("description").(string),
	}

	apiKey, _, err := client.APIKeys.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(apiKey.ID)

	return resourceMetalProjectAPIKeyRead(d, meta)
}

func resourceMetalProjectAPIKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	apiKey, err := client.APIKeys.ProjectGet(d.Id(), d.Get("project_id").(string), nil)
	if err != nil {
		err = friendlyError(err)

		// If the key is somehow already destroyed, mark as
		// succesfully gone
		if isNotFound(err) {
			log.Printf("[WARN] Project APIKey (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.SetId(apiKey.ID)
	d.Set("project_id", apiKey.Project.ID)
	d.Set("description", apiKey.Description)
	d.Set("read_only", apiKey.ReadOnly)
	d.Set("created", apiKey.Created)
	d.Set("updated", apiKey.Updated)

	return nil
}

func resourceMetalProjectAPIKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceMetalProjectAPIKeyRead(d, meta)
}

func resourceMetalProjectAPIKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	resp, err := client.APIKeys.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
