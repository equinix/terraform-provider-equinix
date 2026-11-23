package stream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func testRead(t *testing.T, status int, body string) resource.ReadResponse {
	t.Helper()
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	r := NewResource().(*Resource)
	r.Meta = &config.Config{BaseURL: srv.URL, Token: "test-token"}

	sch := resourceSchema(ctx)
	state := tfsdk.State{Schema: sch, Raw: tftypes.NewValue(sch.Type().TerraformType(ctx), nil)}
	for attr, value := range map[string]string{"id": "00000000-0000-0000-0000-000000000001"} {
		if diags := state.SetAttribute(ctx, path.Root(attr), value); diags.HasError() {
			t.Fatalf("setting %s: %v", attr, diags)
		}
	}

	metaSchema := metaschema.Schema{Attributes: map[string]metaschema.Attribute{
		"module_name": metaschema.StringAttribute{Optional: true},
	}}
	providerMeta := tfsdk.Config{Schema: metaSchema, Raw: tftypes.NewValue(metaSchema.Type().TerraformType(ctx), nil)}

	resp := resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state, ProviderMeta: providerMeta}, &resp)
	return resp
}

func TestResourceRead_goneRemovesFromState(t *testing.T) {
	cases := map[string]struct {
		status int
		body   string
	}{
		"404":           {status: http.StatusNotFound, body: `[{"errorCode":"EQ-3000039","errorMessage":"Not found"}]`},
		"403":           {status: http.StatusForbidden, body: `[{"errorCode":"EQ-3000039","errorMessage":"Forbidden"}]`},
		"deprovisioned": {status: http.StatusOK, body: `{"uuid":"00000000-0000-0000-0000-000000000001","state":"DEPROVISIONED"}`},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resp := testRead(t, tc.status, tc.body)
			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no error, got %v", resp.Diagnostics)
			}
			if !resp.State.Raw.IsNull() {
				t.Error("expected the resource to be removed from state")
			}
		})
	}
}

func TestResourceRead_otherErrorsAreReturned(t *testing.T) {
	resp := testRead(t, http.StatusBadRequest, `[{"errorCode":"EQ-3000039","errorMessage":"Bad request"}]`)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error on 400")
	}
	if resp.State.Raw.IsNull() {
		t.Error("expected the resource to stay in state")
	}
}
