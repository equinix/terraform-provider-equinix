package organization

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			},
			"website": schema.StringAttribute{
				Description: "Website link",
				Optional:    true,
			},
			"twitter": schema.StringAttribute{
				Description: "Twitter handle",
				Optional:    true,
			},
			"logo": schema.StringAttribute{
				Description: "Logo URL",
				Optional:    true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"address": schema.SingleNestedBlock{
				Description: "Address information block",
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
							stringvalidator.RegexMatches(StringToRegex("(?i)^[a-z]{2}$"), "Address country must be a two letter code (ISO 3166-1 alpha-2)"),
						},
					},
					"state": schema.StringAttribute{
						Description: "State name",
						Optional:    true,
					},
				},
			},
		},
	}
}
