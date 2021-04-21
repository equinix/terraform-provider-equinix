package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMetalPreCreatedIPBlock() *schema.Resource {
	s := metalIPComputedFields()
	s["project_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the project where the searched block should be.",
	}
	s["global"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether to look for global block. Default is false for backward compatibility.",
	}
	s["public"] = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Whether to look for public or private block.",
	}

	s["facility"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Facility of the searched block. (for non-global blocks).",
	}

	s["metro"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Metro of the searched block (for non-global blocks).",
	}

	s["address_family"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "4 or 6, depending on which block you are looking for.",
	}
	s["cidr_notation"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "CIDR notation of the looked up block.",
	}
	s["quantity"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s["type"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Read:   dataSourceMetalReservedIPBlockRead,
		Schema: s,
	}
}
