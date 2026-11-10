package asn

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Autonomous System Numbers (ASNs).

A PRIVATE_ASN is provisioned by Equinix in the given metro and range. A PUBLIC_ASN is registered by number and RIR region and stays in the PENDING state until the signed ownership challenge is provided in ` + "`ownership_challenge_details.signature`" + `.

Additional Documentation:
* API: https://docs.equinix.com/api-catalog/fabricv4/#tag/ASNs`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"type": schema.StringAttribute{
				Description: "ASN type. One of: PUBLIC_ASN, PRIVATE_ASN",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("PUBLIC_ASN", "PRIVATE_ASN"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project": schema.SingleNestedAttribute{
				Description: "Project this ASN belongs to",
				Required:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ProjectModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Description: "Equinix-assigned project ID",
						Required:    true,
					},
				},
			},
			"number": schema.Int64Attribute{
				Description: "Autonomous System Number. Required for a PUBLIC_ASN; assigned by Equinix for a PRIVATE_ASN",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"rir_region": schema.StringAttribute{
				Description: "Regional Internet Registry region of a PUBLIC_ASN. One of: AFRINIC, APNIC, ARIN, LACNIC, RIPE_NCC",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("AFRINIC", "APNIC", "ARIN", "LACNIC", "RIPE_NCC"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"range": schema.StringAttribute{
				Description: "Range of a PRIVATE_ASN. One of: 16_BIT, 32_BIT",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("16_BIT", "32_BIT"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.SingleNestedAttribute{
				Description: "Location of a PRIVATE_ASN",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ResourceLocationModel](ctx),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"metro_code": schema.StringAttribute{
						Description: "Two-letter metro code (e.g. SV, DC, AM)",
						Required:    true,
					},
				},
			},
			"ownership_challenge_details": schema.SingleNestedAttribute{
				Description: "Ownership challenge of a PUBLIC_ASN. Sign the returned message and set the signature to complete the verification",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ResourceChallengeModel](ctx),
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Ownership challenge type. One of: RPKI, LOA",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("RPKI", "LOA"),
						},
					},
					"message": schema.StringAttribute{
						Description: "Challenge message to sign",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"signature": schema.StringAttribute{
						Description: "Signed challenge value. Can be set at creation or updated in place",
						Optional:    true,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			"uuid": schema.StringAttribute{
				Description: "Equinix-assigned UUID of the ASN",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"href": schema.StringAttribute{
				Description: "URL of the ASN resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "State of the ASN. One of: REQUESTED, PROVISIONING, PENDING, PROVISIONED, ATTACHED, FAILED, PROVISIONING_FAILED, DEPROVISIONING, DEPROVISIONING_FAILED, DEPROVISIONED",
				Computed:    true,
			},
			"change": schema.SingleNestedAttribute{
				Description: "Current state of the latest change on the ASN",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ChangeModel](ctx),
				Attributes: map[string]schema.Attribute{
					"href": schema.StringAttribute{Description: "URI to the change resource", Computed: true},
				},
			},
			"assets": schema.ListNestedAttribute{
				Description: "Assets (services) the ASN is attached to",
				Computed:    true,
				CustomType:  fwtypes.NewListNestedObjectTypeOf[AssetModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: assetAttributes(ctx),
				},
			},
			"change_log": schema.SingleNestedAttribute{
				Description: "Details of the last change on the ASN",
				Computed:    true,
				CustomType:  fwtypes.NewObjectTypeOf[ChangeLogModel](ctx),
				Attributes: map[string]schema.Attribute{
					"created_date_time": schema.StringAttribute{Description: "Creation date and time", Computed: true},
					"updated_date_time": schema.StringAttribute{Description: "Last update date and time", Computed: true},
					"deleted_date_time": schema.StringAttribute{Description: "Deletion date and time", Computed: true},
				},
			},
		},
	}
}

func assetAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"uuid":             schema.StringAttribute{Description: "Asset identifier", Computed: true},
		"announcing_metro": schema.StringAttribute{Description: "Metro announcing the ASN", Computed: true},
		"service": schema.SingleNestedAttribute{
			Description: "Service the ASN is attached to",
			Computed:    true,
			CustomType:  fwtypes.NewObjectTypeOf[AssetServiceModel](ctx),
			Attributes: map[string]schema.Attribute{
				"href": schema.StringAttribute{Description: "URI to the service", Computed: true},
				"type": schema.StringAttribute{Description: "Service type. One of: IA", Computed: true},
				"uuid": schema.StringAttribute{Description: "Service identifier", Computed: true},
			},
		},
		"ip_blocks": schema.ListNestedAttribute{
			Description: "IP blocks announced with the ASN",
			Computed:    true,
			CustomType:  fwtypes.NewListNestedObjectTypeOf[AssetIpBlockModel](ctx),
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"href": schema.StringAttribute{Description: "URI to the IP block", Computed: true},
					"type": schema.StringAttribute{Description: "IP block reference type. One of: IPV4_IP_BLOCK, IPV6_IP_BLOCK", Computed: true},
					"uuid": schema.StringAttribute{Description: "IP block identifier", Computed: true},
				},
			},
		},
	}
}
