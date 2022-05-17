package userapikey

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	apikey "github.com/equinix/terraform-provider-equinix/internal/resources/metal/projectapikey"
)

func Resource() *schema.Resource {
	userKeySchema := apikey.SchemaMetalAPIKey()
	userKeySchema["user_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "UUID of user owning this key",
	}
	return &schema.Resource{
		Create: apikey.ResourceMetalAPIKeyCreate,
		Read:   apikey.ResourceMetalAPIKeyRead,
		Delete: apikey.ResourceMetalAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: userKeySchema,
	}
}
