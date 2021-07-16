package metal

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func schemaMetalAPIKey() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"read_only": {
			Type:     schema.TypeBool,
			ForceNew: true,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"token": {
			Type:      schema.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}

func resourceMetalProjectAPIKey() *schema.Resource {
	projectKeySchema := schemaMetalAPIKey()
	projectKeySchema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}
	return &schema.Resource{
		Create: resourceMetalAPIKeyCreate,
		Read:   resourceMetalAPIKeyRead,
		Delete: resourceMetalAPIKeyDelete,
		Schema: projectKeySchema,
	}
}

func resourceMetalAPIKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	projectId := ""

	projectIdRaw, projectIdOk := d.GetOk("project_id")
	if projectIdOk {
		projectId = projectIdRaw.(string)
	}

	createRequest := &packngo.APIKeyCreateRequest{
		ProjectID:   projectId,
		ReadOnly:    d.Get("read_only").(bool),
		Description: d.Get("description").(string),
	}

	apiKey, _, err := client.APIKeys.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(apiKey.ID)

	return resourceMetalAPIKeyRead(d, meta)
}

func projectIdFromResourceData(d *schema.ResourceData) string {
	projectIdRaw, projectIdOk := d.GetOk("project_id")
	if projectIdOk {
		return projectIdRaw.(string)
	}
	return ""
}

func resourceMetalAPIKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	projectId := projectIdFromResourceData(d)

	var apiKey *packngo.APIKey
	var err error

	// if project has been set in the resource, look up project API key
	// (this is the reason project API key can't be imported)
	if projectId != "" {
		apiKey, err = client.APIKeys.ProjectGet(projectId, d.Id(),
			&packngo.GetOptions{Includes: []string{"project"}})
	} else {
		apiKey, err = client.APIKeys.UserGet(d.Id(),
			&packngo.GetOptions{Includes: []string{"user"}})
	}

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
	attrMap := map[string]interface{}{
		"description": apiKey.Description,
		"read_only":   apiKey.ReadOnly,
		"token":       apiKey.Token,
	}

	// this is kind of unnecessary as the project ID most likely already set,
	// because project API key can't be imported. But let's refresh the
	// project ID for future-proofing
	if apiKey.Project != nil && apiKey.Project.ID != "" {
		attrMap["project_id"] = apiKey.Project.ID
	}
	if apiKey.User != nil && apiKey.User.ID != "" {
		attrMap["user_id"] = apiKey.User.ID
	}

	return setMap(d, attrMap)
}

func resourceMetalAPIKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	resp, err := client.APIKeys.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
