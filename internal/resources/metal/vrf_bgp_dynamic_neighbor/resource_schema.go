package vrfbgpdynamicneighbor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description:         "This resource manages BGP dynamic neighbor ranges for an Equinix Metal VRF",
		MarkdownDescription: "This resource manages BGP dynamic neighbor ranges for an Equinix Metal VRF, but with markdown",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this the dynamic BGP neighbor",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gateway_id": schema.StringAttribute{
				Description: "The ID of the Equinix Metal VRF gateway for this dynamic BGP neighbor range",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"range": schema.StringAttribute{
				Description: "Network range of the dynamic BGP neighbor in CIDR format",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"asn": schema.Int64Attribute{
				Description: "The ASN of the dynamic BGP neighbor",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Description: "The state of the dynamic BGP neighbor",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				Description: "Tags attached to the dynamic BGP neighbor",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
