package project

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
				Description: "Name of the connection resource",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("project_id"),
					}...),
				},
			},
			"project_id": schema.StringAttribute{
				Description: "ID of project to which the connection belongs",
				Optional:    true,
				Computed:    true,
			},
			"created": schema.StringAttribute{
				Description: "The timestamp for when the project was created",
				Computed:    true,
			},
			"updated": schema.StringAttribute{
				Description: "The timestamp for the last time the project was updated",
				Computed:    true,
			},
			"backend_transfer": schema.BoolAttribute{
				Description: "Whether Backend Transfer is enabled for this project",
				Computed:    true,
			},
			"payment_method_id": schema.StringAttribute{
				Description: "The UUID of payment method for this project",
				Computed:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The UUID of this project's parent organization",
				Computed:    true,
			},
			"user_ids": schema.ListAttribute{
				Description: "List of UUIDs of user accounts which belong to this project",
				ElementType: types.StringType,
				Computed:    true,
			},
			"bgp_config": schema.ListAttribute{
				Description: "Optional BGP settings. Refer to [Equinix Metal guide for BGP](https://metal.equinix.com/developers/docs/networking/local-global-bgp/)",
				CustomType:  fwtypes.NewListNestedObjectTypeOf[BGPConfigModel](ctx),
				ElementType: fwtypes.NewObjectTypeOf[BGPConfigModel](ctx),
				Computed:    true,
			},
		},
	}
}
