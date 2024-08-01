package vrfbgpdynamicneighbor

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Model struct {
	ID        types.String `tfsdk:"id"`
	GatewayID types.String `tfsdk:"gateway_id"`
	Range     types.String `tfsdk:"range"`
	ASN       types.Int64  `tfsdk:"asn"`
	State     types.String `tfsdk:"state"`
	Tags      types.List   `tfsdk:"tags"` // List of strings
}

func (m *Model) parse(ctx context.Context, neighbor *metalv1.BgpDynamicNeighbor) (d diag.Diagnostics) {
	m.ID = types.StringValue(neighbor.GetId())

	m.GatewayID = types.StringValue(neighbor.MetalGateway.GetId())
	m.Range = types.StringValue(neighbor.GetBgpNeighborRange())
	m.ASN = types.Int64Value(neighbor.GetBgpNeighborAsn())
	m.State = types.StringValue(string(neighbor.GetState()))

	tags, diags := types.ListValueFrom(ctx, types.StringType, neighbor.GetTags())
	if diags.HasError() {
		return diags
	}

	m.Tags = tags

	return nil
}
