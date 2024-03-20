package project

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	equinix_planmodifiers "github.com/equinix/terraform-provider-equinix/internal/planmodifiers"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"name": schema.StringAttribute{
				Description: "The name of the project. The maximum length is 80 characters",
				Required:    true,
			},
			"created": schema.StringAttribute{
				Description: "The timestamp for when the project was created",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated": schema.StringAttribute{
				Description: "The timestamp for the last time the project was updated",
				Computed:    true,
			},
			"backend_transfer": schema.BoolAttribute{
				MarkdownDescription: "Enable or disable [Backend Transfer](https://metal.equinix.com/developers/docs/networking/backend-transfer/), default is false",
				Description:         "Enable or disable Backend Transfer, default is false",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"payment_method_id": schema.StringAttribute{
				Description: "The UUID of payment method for this project. The payment method and the project need to belong to the same organization (passed with organization_id, or default)",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					equinix_validation.UUID(),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "The UUID of organization under which the project is created",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					equinix_validation.UUID(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"bgp_config": schema.ListNestedBlock{
				Description: "Address information block",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(1),
				},
				PlanModifiers: []planmodifier.List{
					equinix_planmodifiers.ImmutableList(),
				},
				CustomType: fwtypes.NewListNestedObjectTypeOf[BGPConfigModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: bgpConfigSchema,
				},
			},
		},
	}
}

var bgpConfigSchema = map[string]schema.Attribute{
	"deployment_type": schema.StringAttribute{
		MarkdownDescription: "The BGP deployment type, either 'local' or 'global'. The local is likely to be usable immediately, the global will need to be review by Equinix Metal engineers",
		Description:         "The BGP deployment type, either 'local' or 'global'",
		Required:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("local", "global"),
		},
	},
	"asn": schema.Int64Attribute{
		Description: "Autonomous System Number for local BGP deployment",
		Required:    true,
		PlanModifiers: []planmodifier.Int64{
			equinix_planmodifiers.ImmutableInt64(),
		},
	},
	"md5": schema.StringAttribute{
		Description: "Password for BGP session in plaintext (not a checksum)",
		Sensitive:   true,
		Optional:    true,
	},
	"status": schema.StringAttribute{
		Description: "Status of BGP configuration in the project",
		Computed:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	},
	"max_prefix": schema.Int64Attribute{
		Description: "The maximum number of route filters allowed per server",
		Computed:    true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	},
}
