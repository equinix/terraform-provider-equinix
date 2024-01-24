package metal_ssh_key

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func commonFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

func CommonFieldsResource() map[string]*schema.Schema {
	resourceSchema := commonFields()
	resourceSchema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the SSH key for identification",
		Required:    true,
	}
	resourceSchema["public_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The public key that will be authorized for SSH access on Equinix Metal devices provisioned with this key",
		Required:    true,
		ForceNew:    true,
	}
	return resourceSchema
}

func CommonFieldsDataSource() map[string]*schema.Schema {
	dsSchema := commonFields()
	dsSchema["search"] = &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The name, fingerprint, id, or public_key of the SSH Key to search for in the Equinix Metal project",
		Optional:     true,
		ValidateFunc: validation.NoZeroValues,
	}
	dsSchema["id"] = &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The id of the SSH Key",
		Optional:     true,
		ValidateFunc: validation.NoZeroValues,
		Computed:     true,
	}
	dsSchema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The label of the Equinix Metal SSH Key",
		Computed:    true,
	}
	dsSchema["public_key"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The public SSH key that is authorized for SSH access on Equinix Metal devices provisioned with this key",
		Computed:    true,
	}
	return dsSchema
}
