package metal_project_ssh_key

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/metal_ssh_key"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	pkeySchema := metal_ssh_key.CommonFieldsResource()
	pkeySchema["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The ID of parent project",
		ForceNew:    true,
		Required:    true,
	}
	resource := metal_ssh_key.Resource()
	resource.Schema = pkeySchema
	return resource
}
