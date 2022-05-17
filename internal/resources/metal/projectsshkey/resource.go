package projectsshkey

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/sshkey"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	pkeySchema := sshkey.MetalSSHKeyCommonFields()
	pkeySchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The ID of parent project",
		ForceNew:    true,
		Required:    true,
	}
	return &schema.Resource{
		Create: sshkey.ResourceMetalSSHKeyCreate,
		Read:   sshkey.ResourceMetalSSHKeyRead,
		Update: sshkey.ResourceMetalSSHKeyUpdate,
		Delete: sshkey.ResourceMetalSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: pkeySchema,
	}
}
