package streamattachment

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func resourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Stream Attachments

Additional Documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/KnowledgeCenter/Fabric/GettingStarted/Integrating-with-Fabric-V4-APIs/IntegrateWithSink.htm
* API: https://developer.equinix.com/catalog/fabricv4#tag/Streams`,
		Attributes: map[string]schema.Attribute{
			"id": framework.IDAttributeDefaultDescription(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"stream_id": schema.StringAttribute{
				Description: "UUID of the stream that is the target of this asset attachment",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"asset": schema.StringAttribute{
				Description: "Equinix defined asset category. Matches the product name the asset is a part of",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"asset_id": schema.StringAttribute{
				Description: "Equinix defined UUID of the asset being attached to the stream",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metrics_enabled": schema.BoolAttribute{
				Description: "Boolean value indicating enablement of metrics for this asset stream attachment",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Equinix defined type for the asset stream attachment",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"href": schema.StringAttribute{
				Description: "Equinix auto generated URI to the stream attachment in Equinix Portal",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "Equinix-assigned unique id for the stream attachment",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attachment_status": schema.StringAttribute{
				Description: "Value representing status for the stream attachment",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
