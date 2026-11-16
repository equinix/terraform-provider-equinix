package streamsubscription

import (
	"context"
	"testing"

	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestStreamSubscriptionWriteOnlyCredentialSchema(t *testing.T) {
	ctx := context.Background()
	s := resourceSchema(ctx)

	if diags := s.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema implementation: %v", diags)
	}

	writeOnly := s.Attributes["sink_credential_wo"].(schema.SingleNestedAttribute)
	if !writeOnly.IsWriteOnly() {
		t.Error("sink_credential_wo is not write-only")
	}
	for _, name := range []string{"access_token", "integration_key", "api_key", "password"} {
		attr := writeOnly.Attributes[name]
		if !attr.IsWriteOnly() || !attr.IsSensitive() {
			t.Errorf("sink_credential_wo.%s: write-only %t, sensitive %t", name, attr.IsWriteOnly(), attr.IsSensitive())
		}
	}
	if s.Attributes["sink_credential_wo_version"].IsWriteOnly() {
		t.Error("sink_credential_wo_version must be stored so that changing it plans an update")
	}
}

func testPlan(ctx context.Context, password types.String) resourceModel {
	credential := sinkCredentialModel{
		Type:     types.StringValue("USERNAME_PASSWORD"),
		Username: types.StringValue("stream"),
		Password: password,
	}
	sink := sinkModel{
		Type:       types.StringValue("WEBHOOK"),
		Credential: fwtypes.NewObjectValueOf(ctx, &credential),
		Settings:   fwtypes.NewObjectValueOfNull[sinkSettingsModel](ctx),
	}
	plan := resourceModel{}
	plan.Type = types.StringValue("STREAM_SUBSCRIPTION")
	plan.Name = types.StringValue("name")
	plan.Description = types.StringValue("description")
	plan.Enabled = types.BoolValue(true)
	plan.MetricSelector = fwtypes.NewObjectValueOfNull[selectorModel](ctx)
	plan.EventSelector = fwtypes.NewObjectValueOfNull[selectorModel](ctx)
	plan.Sink = fwtypes.NewObjectValueOf(ctx, &sink)
	return plan
}

func TestBuildRequestsUseWriteOnlyPassword(t *testing.T) {
	ctx := context.Background()
	writeOnly := fwtypes.NewObjectValueOf(ctx, &sinkCredentialWriteOnlyModel{
		AccessToken:    types.StringNull(),
		IntegrationKey: types.StringNull(),
		APIKey:         types.StringNull(),
		Password:       types.StringValue("write-only-secret"),
	})

	post, diags := buildCreateRequest(ctx, testPlan(ctx, types.StringValue("")), writeOnly)
	if diags.HasError() {
		t.Fatalf("buildCreateRequest: %v", diags)
	}
	credential := post.GetSink().Credential
	if got := credential.GetPassword(); got != "write-only-secret" {
		t.Errorf("create password = %q, want the write-only value", got)
	}
	if got := credential.GetUsername(); got != "stream" {
		t.Errorf("create username = %q, want %q", got, "stream")
	}

	put, diags := buildUpdateRequest(ctx, testPlan(ctx, types.StringNull()), writeOnly)
	if diags.HasError() {
		t.Fatalf("buildUpdateRequest: %v", diags)
	}
	if got := put.GetSink().Credential.GetPassword(); got != "write-only-secret" {
		t.Errorf("update password = %q, want the write-only value", got)
	}
}

func TestBuildCreateRequestKeepsStatefulPassword(t *testing.T) {
	ctx := context.Background()
	post, diags := buildCreateRequest(ctx, testPlan(ctx, types.StringValue("stateful-secret")), fwtypes.NewObjectValueOfNull[sinkCredentialWriteOnlyModel](ctx))
	if diags.HasError() {
		t.Fatalf("buildCreateRequest: %v", diags)
	}
	if got := post.GetSink().Credential.GetPassword(); got != "stateful-secret" {
		t.Errorf("password = %q, want the stateful value", got)
	}
}

func TestClearWhenWriteOnlySet(t *testing.T) {
	ctx := context.Background()
	writeOnlyType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"password": tftypes.String}}
	configType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"sink_credential_wo": writeOnlyType}}
	configSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"sink_credential_wo": schema.SingleNestedAttribute{
				Optional:  true,
				WriteOnly: true,
				Attributes: map[string]schema.Attribute{
					"password": schema.StringAttribute{Optional: true, WriteOnly: true},
				},
			},
		},
	}
	modifier := clearWhenWriteOnlySet(path.Root("sink_credential_wo").AtName("password"))

	cases := map[string]struct {
		writeOnly tftypes.Value
		want      types.String
	}{
		"write-only password set": {
			writeOnly: tftypes.NewValue(writeOnlyType, map[string]tftypes.Value{"password": tftypes.NewValue(tftypes.String, "new")}),
			want:      types.StringValue(""),
		},
		"write-only password null": {
			writeOnly: tftypes.NewValue(writeOnlyType, map[string]tftypes.Value{"password": tftypes.NewValue(tftypes.String, nil)}),
			want:      types.StringValue("old"),
		},
		"write-only object null": {
			writeOnly: tftypes.NewValue(writeOnlyType, nil),
			want:      types.StringValue("old"),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			config := tfsdk.Config{
				Schema: configSchema,
				Raw:    tftypes.NewValue(configType, map[string]tftypes.Value{"sink_credential_wo": tc.writeOnly}),
			}
			req := planmodifier.StringRequest{Config: config, PlanValue: types.StringValue("old")}
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
			modifier.PlanModifyString(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("diagnostics: %v", resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(tc.want) {
				t.Errorf("plan = %s, want %s", resp.PlanValue, tc.want)
			}
		})
	}
}
