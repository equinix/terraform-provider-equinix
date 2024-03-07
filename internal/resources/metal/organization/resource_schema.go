package organization

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func GetResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"name": schema.StringAttribute{
				Description: "The name of the Organization",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description string",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"website": schema.StringAttribute{
				Description: "Website link",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"twitter": schema.StringAttribute{
				Description: "Twitter handle",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"logo": schema.StringAttribute{
				Description: "Logo URL",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"address": schema.ListNestedBlock{
				Description: "Address information block",
				Validators: []validator.List{
					listvalidator.IsRequired(),
					listvalidator.SizeAtLeast(1),
					listvalidator.SizeAtMost(1),
				},
				CustomType: fwtypes.NewListNestedObjectTypeOf[AddressResourceModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Description: "Postal address",
							Required:    true,
						},
						"city": schema.StringAttribute{
							Description: "City name",
							Required:    true,
						},
						"zip_code": schema.StringAttribute{
							Description: "Zip Code",
							Required:    true,
						},
						"country": schema.StringAttribute{
							Description: "Two letter country code (ISO 3166-1 alpha-2), e.g. US",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 18),
								equinix_validation.StringIsCountryCode,
							},
						},
						"state": schema.StringAttribute{
							Description: "State name",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
		},
	}
}
