package streamattachment

import (
	"context"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestMetricsEnabledForAsset(t *testing.T) {
	tests := []struct {
		asset     fabricv4.Asset
		requested bool
		want      bool
	}{
		{fabricv4.ASSET_PROJECTS, true, false},
		{fabricv4.ASSET_PROJECTS, false, false},
		{fabricv4.ASSET_ORGANIZATIONS, true, false},
		{fabricv4.ASSET_ORGANIZATIONS, false, false},
		{fabricv4.ASSET_CONNECTIONS, true, true},
		{fabricv4.ASSET_CONNECTIONS, false, false},
		{fabricv4.ASSET_PORTS, true, true},
		{fabricv4.ASSET_PORTS, false, false},
		{fabricv4.ASSET_ROUTERS, true, true},
		{fabricv4.ASSET_ROUTERS, false, false},
	}

	for _, tt := range tests {
		got := metricsEnabledForAsset(tt.asset, tt.requested)
		if got != tt.want {
			t.Errorf("metricsEnabledForAsset(%q, %v) = %v, want %v", tt.asset, tt.requested, got, tt.want)
		}
	}
}

// minimalSchemaForModifierTest builds the smallest tfsdk.Config that satisfies
// metricsEnabledModifier: it only needs "asset" to be readable.
func minimalConfig(t *testing.T, assetValue string, metricsEnabled bool) tfsdk.Config {
	t.Helper()

	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"asset":           schema.StringAttribute{Required: true},
			"metrics_enabled": schema.BoolAttribute{Optional: true, Computed: true},
		},
	}

	raw := tftypes.NewValue(
		tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"asset":           tftypes.String,
			"metrics_enabled": tftypes.Bool,
		}},
		map[string]tftypes.Value{
			"asset":           tftypes.NewValue(tftypes.String, assetValue),
			"metrics_enabled": tftypes.NewValue(tftypes.Bool, metricsEnabled),
		},
	)

	return tfsdk.Config{Raw: raw, Schema: s}
}

func TestMetricsEnabledModifier(t *testing.T) {
	tests := []struct {
		name        string
		asset       string
		planValue   bool
		wantValue   bool
		wantWarning bool
	}{
		{
			name:        "projects with true gets forced to false with warning",
			asset:       string(fabricv4.ASSET_PROJECTS),
			planValue:   true,
			wantValue:   false,
			wantWarning: true,
		},
		{
			name:        "projects with false stays false no warning",
			asset:       string(fabricv4.ASSET_PROJECTS),
			planValue:   false,
			wantValue:   false,
			wantWarning: false,
		},
		{
			name:        "organizations with true gets forced to false with warning",
			asset:       string(fabricv4.ASSET_ORGANIZATIONS),
			planValue:   true,
			wantValue:   false,
			wantWarning: true,
		},
		{
			name:        "connections with true is unchanged",
			asset:       string(fabricv4.ASSET_CONNECTIONS),
			planValue:   true,
			wantValue:   true,
			wantWarning: false,
		},
		{
			name:        "connections with false is unchanged",
			asset:       string(fabricv4.ASSET_CONNECTIONS),
			planValue:   false,
			wantValue:   false,
			wantWarning: false,
		},
		{
			name:        "ports with true is unchanged",
			asset:       string(fabricv4.ASSET_PORTS),
			planValue:   true,
			wantValue:   true,
			wantWarning: false,
		},
	}

	modifier := metricsEnabledModifier{}
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := planmodifier.BoolRequest{
				Path:      path.Root("metrics_enabled"),
				PlanValue: types.BoolValue(tt.planValue),
				Config:    minimalConfig(t, tt.asset, tt.planValue),
			}
			var resp planmodifier.BoolResponse
			resp.PlanValue = req.PlanValue

			modifier.PlanModifyBool(ctx, req, &resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics.Errors())
			}
			if tt.wantWarning && resp.Diagnostics.WarningsCount() == 0 {
				t.Error("expected a warning but got none")
			}
			if !tt.wantWarning && resp.Diagnostics.WarningsCount() > 0 {
				t.Errorf("unexpected warning: %v", resp.Diagnostics.Warnings())
			}
			if resp.PlanValue.ValueBool() != tt.wantValue {
				t.Errorf("PlanValue = %v, want %v", resp.PlanValue.ValueBool(), tt.wantValue)
			}
		})
	}
}
