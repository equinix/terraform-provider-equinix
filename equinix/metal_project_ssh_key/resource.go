package metal_project_ssh_key

import (
	"github.com/equinix/terraform-provider-equinix/equinix/metal_ssh_key"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	pkeySchema := metal_ssh_key.MetalSSHKeyCommonFields()
	pkeySchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The ID of parent project",
		ForceNew:    true,
		Required:    true,
	}
	return &schema.Resource{
		Create: metal_ssh_key.ResourceMetalSSHKeyCreate,
		Read:   metal_ssh_key.ResourceMetalSSHKeyRead,
		Update: metal_ssh_key.ResourceMetalSSHKeyUpdate,
		Delete: metal_ssh_key.ResourceMetalSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: pkeySchema,
	}
}
