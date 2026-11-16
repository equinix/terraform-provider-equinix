package streamattachment

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// metricsEnabledForAsset returns false for asset types that the API does not allow
// metricsEnabled=true for (projects, organizations), and returns the requested value otherwise.
func metricsEnabledForAsset(asset fabricv4.Asset, requested bool) bool {
	if asset == fabricv4.ASSET_PROJECTS || asset == fabricv4.ASSET_ORGANIZATIONS {
		return false
	}
	return requested
}

// metricsEnabledModifier forces metrics_enabled to false for asset types that do not
// support metrics (projects, organizations). The API returns EQ-3045892 (HTTP 400) when
// metricsEnabled=true is sent for those asset types.
type metricsEnabledModifier struct{}

func (m metricsEnabledModifier) Description(_ context.Context) string {
	return "Forces metrics_enabled to false for projects and organizations, which do not support metrics."
}

func (m metricsEnabledModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m metricsEnabledModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	var asset types.String
	diags := req.Config.GetAttribute(ctx, path.Root("asset"), &asset)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if asset.ValueString() == string(fabricv4.ASSET_PROJECTS) || asset.ValueString() == string(fabricv4.ASSET_ORGANIZATIONS) {
		if req.PlanValue.ValueBool() {
			resp.Diagnostics.AddWarning(
				"metrics_enabled is not supported for projects or organizations",
				"The Equinix Fabric API does not allow metricsEnabled=true for project or organization stream attachments. "+
					"The value will be set to false. To get metrics, attach individual connections or ports with metrics_enabled=true.",
			)
		}
		resp.PlanValue = types.BoolValue(false)
	}
}
