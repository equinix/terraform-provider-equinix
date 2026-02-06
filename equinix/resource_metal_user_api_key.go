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
		Create:             resourceMetalAPIKeyCreate,
		Read:               resourceMetalAPIKeyRead,
		Delete:             resourceMetalAPIKeyDelete,
		DeprecationMessage: "Metal platform reaches end-of-life June 30, 2026. Removal scheduled for provider version 5.0.0. Continue using 4.x releases through the sunset period. Reference: https://docs.equinix.com/metal/",
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: userKeySchema,
	}
}
