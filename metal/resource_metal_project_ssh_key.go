package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMetalProjectSSHKey() *schema.Resource {
	pkeySchema := metalSSHKeyCommonFields()
	pkeySchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The ID of parent project",
		ForceNew:    true,
		Required:    true,
	}
	return &schema.Resource{
		Create: resourceMetalSSHKeyCreate,
		Read:   resourceMetalSSHKeyRead,
		Update: resourceMetalSSHKeyUpdate,
		Delete: resourceMetalSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: pkeySchema,
	}
}
