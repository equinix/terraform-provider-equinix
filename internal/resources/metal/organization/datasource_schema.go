package organization

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{

			"id": framework.IDAttributeDefaultDescription(),
			"name": schema.StringAttribute{
				Description: "The name of the Organization",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("organization_id"), path.MatchRoot("name")),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "The UUID of the organization resource",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("organization_id"), path.MatchRoot("name")),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description string",
				Optional:    true,
			},
			"website": schema.StringAttribute{
				Description: "Website link",
				Computed:    true,
			},
			"twitter": schema.StringAttribute{
				Description: "Twitter handle",
				Computed:    true,
			},
			"logo": schema.StringAttribute{
				DeprecationMessage: "Empty string is the only valid value. Support for this field has been removed from the API",
				Description: "**Deprecated** Logo URL",
				Computed:    true,
			},

			"project_ids": schema.ListAttribute{
				Description: "UUIDs of project resources which belong to this organization",
				Computed:    true,
				ElementType: types.StringType,
			},
			"address": schema.ListAttribute{
				Description: "Business' address",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[AddressResourceModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[AddressResourceModel](ctx),
				Computed:    true,
			},
		},
	}
}
