package equinix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMetalUserAPIKey() *schema.Resource {
	userKeySchema := schemaMetalAPIKey()
	userKeySchema["user_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "UUID of user owning this key",
	}
	return &schema.Resource{
		Create: resourceMetalAPIKeyCreate,
		Read:   resourceMetalAPIKeyRead,
		Delete: resourceMetalAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: userKeySchema,
	}
}
