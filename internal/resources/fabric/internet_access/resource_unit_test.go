package internetaccess

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const testServiceUUID = "11111111-2222-3333-4444-555555555555"

func init() {
	waitDelay = time.Millisecond
	waitMinInterval = time.Millisecond
}

func sampleService(t *testing.T) *fabricv4.InternetAccessService {
	t.Helper()
	svc := sampleServices(t).GetData()[0]
	return &svc
}

// resourceModelWithNulls returns a model whose nested values are typed nulls so that it
// can be written to a state of the resource schema.
func resourceModelWithNulls(ctx context.Context) ResourceModel {
	return ResourceModel{
		Timeouts:        timeouts.Value{Object: types.ObjectNull(resourceSchema(ctx).Attributes["timeouts"].GetType().(timeouts.Type).ObjectType.AttrTypes)},
		Change:          fwtypes.NewObjectValueOfNull[ChangeModel](ctx),
		Locations:       fwtypes.NewListNestedObjectValueOfNull[LocationModel](ctx),
		RoutingProtocol: fwtypes.NewObjectValueOfNull[ResourceRoutingModel](ctx),
		Billing:         fwtypes.NewObjectValueOfNull[ResourceBillingModel](ctx),
		Account:         fwtypes.NewObjectValueOfNull[ResourceAccountModel](ctx),
		Project:         fwtypes.NewObjectValueOfNull[ProjectModel](ctx),
		Order:           fwtypes.NewObjectValueOfNull[ResourceOrderModel](ctx),
		ChangeLog:       fwtypes.NewObjectValueOfNull[ChangeLogModel](ctx),
	}
}

func TestResourceSchemaMatchesModel(t *testing.T) {
	ctx := context.Background()
	s := resourceSchema(ctx)
	if diags := s.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema is not valid: %v", diags)
	}

	model := resourceModelWithNulls(ctx)
	if diags := model.parse(ctx, sampleService(t)); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	state := tfsdk.State{Schema: s}
	if diags := state.Set(ctx, &model); diags.HasError() {
		t.Fatalf("state does not match the schema: %v", diags)
	}
}

func TestResourceModelParse(t *testing.T) {
	ctx := context.Background()
	model := resourceModelWithNulls(ctx)
	if diags := model.parse(ctx, sampleService(t)); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}

	if model.ID.ValueString() != "svc-1" || model.UUID.ValueString() != "svc-1" {
		t.Errorf("id/uuid = %s/%s", model.ID, model.UUID)
	}
	if model.Name.ValueString() != "first" || model.Type.ValueString() != "SINGLE_IA" || model.State.ValueString() != "PROVISIONED" {
		t.Errorf("name/type/state = %s/%s/%s", model.Name, model.Type, model.State)
	}
	if model.Bandwidth.ValueInt32() != 15 || model.BandwidthCommit.ValueInt32() != 5 {
		t.Errorf("bandwidth = %d, commit = %d", model.Bandwidth.ValueInt32(), model.BandwidthCommit.ValueInt32())
	}

	var routing ResourceRoutingModel
	if diags := model.RoutingProtocol.As(ctx, &routing, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("routing: %v", diags)
	}
	if routing.Type.ValueString() != "DIRECT" {
		t.Errorf("routing type = %s", routing.Type)
	}
	conns, _ := routing.Connections.ToSlice(ctx)
	if len(conns) != 1 || conns[0].UUID.ValueString() != "conn-1" {
		t.Fatalf("connections = %+v", conns)
	}
	var v4 PeeringModel
	conns[0].PeeringIpv4.As(ctx, &v4, basetypes.ObjectAsOptions{})
	if v4.EquinixPeerIp.ValueString() != "137.174.168.34" {
		t.Errorf("equinix peer ip = %s", v4.EquinixPeerIp)
	}
	routes, _ := routing.CustomerRoutes.ToSlice(ctx)
	if len(routes) != 1 {
		t.Fatalf("customer routes = %d", len(routes))
	}

	var order ResourceOrderModel
	model.Order.As(ctx, &order, basetypes.ObjectAsOptions{})
	if order.PurchaseOrderNumber.ValueString() != "PO-1" || order.TermLength.ValueInt32() != 12 {
		t.Errorf("order = %+v", order)
	}
	var billing ResourceBillingModel
	model.Billing.As(ctx, &billing, basetypes.ObjectAsOptions{})
	if billing.Type.ValueString() != "FIXED" || !billing.Enabled.ValueBool() {
		t.Errorf("billing = %+v", billing)
	}
}

func TestResourceModelParse_KeepsConfiguredBgpOptions(t *testing.T) {
	ctx := context.Background()
	model := resourceModelWithNulls(ctx)
	model.RoutingProtocol = fwtypes.NewObjectValueOf[ResourceRoutingModel](ctx, &ResourceRoutingModel{
		Type:             types.StringValue("BGP"),
		CustomerRoutes:   fwtypes.NewListNestedObjectValueOfNull[CustomerRouteModel](ctx),
		Connections:      fwtypes.NewListNestedObjectValueOfNull[ResourceConnectionModel](ctx),
		ExportPolicy:     types.StringValue("FULL"),
		Asn:              fwtypes.NewObjectValueOf[ResourceAsnModel](ctx, &ResourceAsnModel{UUID: types.StringValue("asn-1")}),
		CustomerAsn:      types.Int64Value(65000),
		CustomerAsnRange: types.StringValue("BITS_16"),
		BgpAuthKey:       types.StringValue("secret"),
	})
	if diags := model.parse(ctx, sampleService(t)); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}

	var routing ResourceRoutingModel
	model.RoutingProtocol.As(ctx, &routing, basetypes.ObjectAsOptions{})
	if routing.ExportPolicy.ValueString() != "FULL" || routing.CustomerAsn.ValueInt64() != 65000 ||
		routing.CustomerAsnRange.ValueString() != "BITS_16" || routing.BgpAuthKey.ValueString() != "secret" {
		t.Errorf("BGP options were not kept: %+v", routing)
	}
	var asn ResourceAsnModel
	routing.Asn.As(ctx, &asn, basetypes.ObjectAsOptions{})
	if asn.UUID.ValueString() != "asn-1" {
		t.Errorf("asn uuid = %s", asn.UUID)
	}
}

func planFor(ctx context.Context, routingType string, extra func(*ResourceRoutingModel)) ResourceModel {
	conn := &ResourceConnectionModel{
		UUID: types.StringValue("conn-1"),
		Href: types.StringUnknown(),
		PeeringIpv4: fwtypes.NewObjectValueOf[PeeringModel](ctx, &PeeringModel{
			EquinixPeerIp: types.StringValue("10.0.0.1"),
			EquinixVrrpIp: types.StringValue("10.0.0.2"),
			Prefix:        types.StringUnknown(),
		}),
		PeeringIpv6: fwtypes.NewObjectValueOfUnknown[PeeringModel](ctx),
	}
	routing := ResourceRoutingModel{
		Type: types.StringValue(routingType),
		CustomerRoutes: fwtypes.NewListNestedObjectValueOfSlice[CustomerRouteModel](ctx, []*CustomerRouteModel{
			{IpBlock: fwtypes.NewObjectValueOf[IpBlockModel](ctx, &IpBlockModel{UUID: types.StringValue("blk-1"), Href: types.StringUnknown()})},
		}),
		Connections:      fwtypes.NewListNestedObjectValueOfSlice[ResourceConnectionModel](ctx, []*ResourceConnectionModel{conn}),
		ExportPolicy:     types.StringNull(),
		Asn:              fwtypes.NewObjectValueOfNull[ResourceAsnModel](ctx),
		CustomerAsn:      types.Int64Null(),
		CustomerAsnRange: types.StringNull(),
		BgpAuthKey:       types.StringNull(),
	}
	if extra != nil {
		extra(&routing)
	}

	model := resourceModelWithNulls(ctx)
	model.Type = types.StringValue("SINGLE_IA")
	model.Name = types.StringValue("eia_PFCR")
	model.Bandwidth = types.Int32Value(50)
	model.BandwidthCommit = types.Int32Unknown()
	model.RoutingProtocol = fwtypes.NewObjectValueOf[ResourceRoutingModel](ctx, &routing)
	model.Billing = fwtypes.NewObjectValueOf[ResourceBillingModel](ctx, &ResourceBillingModel{Type: types.StringValue("FIXED")})
	model.Account = fwtypes.NewObjectValueOf[ResourceAccountModel](ctx, &ResourceAccountModel{AccountNumber: types.StringValue("709976")})
	model.Project = fwtypes.NewObjectValueOf[ProjectModel](ctx, &ProjectModel{ProjectID: types.StringValue("proj-1")})
	model.Order = fwtypes.NewObjectValueOf[ResourceOrderModel](ctx, &ResourceOrderModel{PurchaseOrderNumber: types.StringValue("PO-9")})
	return model
}

func requestJSON(t *testing.T, plan ResourceModel) map[string]any {
	t.Helper()
	req, err := buildCreateRequest(context.Background(), plan)
	if err != nil {
		t.Fatalf("buildCreateRequest: %v", err)
	}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestBuildCreateRequest_Direct(t *testing.T) {
	got := requestJSON(t, planFor(context.Background(), "DIRECT", nil))

	if got["type"] != "SINGLE_IA" || got["name"] != "eia_PFCR" {
		t.Errorf("type/name = %v/%v", got["type"], got["name"])
	}
	if got["bandwidth"] != float64(50) {
		t.Errorf("bandwidth = %v", got["bandwidth"])
	}
	if _, ok := got["bandwidthCommit"]; ok {
		t.Error("unknown bandwidth_commit must not be sent")
	}
	if got["project"].(map[string]any)["projectId"] != "proj-1" {
		t.Errorf("project = %v", got["project"])
	}
	if got["account"].(map[string]any)["accountNumber"] != "709976" {
		t.Errorf("account = %v", got["account"])
	}
	if got["billing"].(map[string]any)["type"] != "FIXED" {
		t.Errorf("billing = %v", got["billing"])
	}
	if got["order"].(map[string]any)["purchaseOrderNumber"] != "PO-9" {
		t.Errorf("order = %v", got["order"])
	}

	rp := got["routingProtocol"].(map[string]any)
	if rp["type"] != "DIRECT" {
		t.Errorf("routing type = %v", rp["type"])
	}
	routes := rp["customerRoutes"].([]any)
	if len(routes) != 1 || routes[0].(map[string]any)["ipBlock"].(map[string]any)["uuid"] != "blk-1" {
		t.Errorf("customerRoutes = %v", routes)
	}
	conns := rp["connections"].([]any)
	if len(conns) != 1 {
		t.Fatalf("connections = %v", conns)
	}
	conn := conns[0].(map[string]any)
	if conn["uuid"] != "conn-1" {
		t.Errorf("connection uuid = %v", conn["uuid"])
	}
	v4 := conn["peeringIpv4"].(map[string]any)
	if v4["equinixPeerIp"] != "10.0.0.1" || v4["equinixVrrpIp"] != "10.0.0.2" {
		t.Errorf("peeringIpv4 = %v", v4)
	}
	if _, ok := conn["peeringIpv6"]; ok {
		t.Error("unknown peering_ipv6 must not be sent")
	}
	if _, ok := rp["exportPolicy"]; ok {
		t.Error("BGP options must not be sent for DIRECT")
	}
}

func TestBuildCreateRequest_Bgp(t *testing.T) {
	ctx := context.Background()
	got := requestJSON(t, planFor(ctx, "BGP", func(r *ResourceRoutingModel) {
		r.ExportPolicy = types.StringValue("FULL")
		r.Asn = fwtypes.NewObjectValueOf[ResourceAsnModel](ctx, &ResourceAsnModel{UUID: types.StringValue("asn-1")})
		r.CustomerAsn = types.Int64Value(65000)
		r.CustomerAsnRange = types.StringValue("BITS_16")
		r.BgpAuthKey = types.StringValue("secret")
	}))

	rp := got["routingProtocol"].(map[string]any)
	if rp["exportPolicy"] != "FULL" || rp["customerAsn"] != float64(65000) ||
		rp["customerAsnRange"] != "BITS_16" || rp["bgpAuthKey"] != "secret" {
		t.Errorf("BGP options = %v", rp)
	}
	if rp["asn"].(map[string]any)["uuid"] != "asn-1" {
		t.Errorf("asn = %v", rp["asn"])
	}
}

func TestBuildCreateRequest_OptionalOrderAndRoutesOmitted(t *testing.T) {
	ctx := context.Background()
	plan := planFor(ctx, "STATIC", func(r *ResourceRoutingModel) {
		r.CustomerRoutes = fwtypes.NewListNestedObjectValueOfUnknown[CustomerRouteModel](ctx)
	})
	plan.Order = fwtypes.NewObjectValueOfUnknown[ResourceOrderModel](ctx)

	got := requestJSON(t, plan)
	if _, ok := got["order"]; ok {
		t.Error("unknown order must not be sent")
	}
	rp := got["routingProtocol"].(map[string]any)
	if routes, _ := rp["customerRoutes"].([]any); len(routes) != 0 {
		t.Errorf("customerRoutes = %v, want empty", rp["customerRoutes"])
	}
}

func TestUpdateIsNotSupportedYet(t *testing.T) {
	var resp resource.UpdateResponse
	(&Resource{}).Update(context.Background(), resource.UpdateRequest{}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic")
	}
}

// ---- waiters ----

type fakeResponse struct {
	status int
	body   string
}

func newTestFabricClient(t *testing.T, responses ...fakeResponse) (*fabricv4.APIClient, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := int(calls.Add(1)) - 1
		if n >= len(responses) {
			n = len(responses) - 1
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responses[n].status)
		fmt.Fprint(w, responses[n].body)
	}))
	t.Cleanup(server.Close)

	cfg := fabricv4.NewConfiguration()
	cfg.Servers = fabricv4.ServerConfigurations{{URL: server.URL}}
	cfg.HTTPClient = server.Client()
	return fabricv4.NewAPIClient(cfg), &calls
}

func serviceWithState(state string) fakeResponse {
	return fakeResponse{http.StatusOK, fmt.Sprintf(`{
  "href": "h", "type": "SINGLE_IA", "uuid": %q, "name": "n", "state": %q,
  "change": {"href": "h"},
  "routingProtocol": {"type": "DIRECT", "customerRoutes": [], "connections": []},
  "billing": {"type": "FIXED"},
  "account": {"accountNumber": "1"},
  "project": {"projectId": "p"},
  "order": {},
  "changeLog": {"createdDateTime": "2026-04-23T15:19:07Z"},
  "useCase": "MAIN"
}`, testServiceUUID, state)}
}

func TestDeleteWaiter(t *testing.T) {
	t.Run("deprovisioned", func(t *testing.T) {
		client, _ := newTestFabricClient(t, serviceWithState("PROVISIONED"), serviceWithState("DEPROVISIONING"), serviceWithState("DEPROVISIONED"))
		if _, err := getDeleteWaiter(context.Background(), client, testServiceUUID, 5*time.Second).WaitForStateContext(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("failed is final", func(t *testing.T) {
		client, _ := newTestFabricClient(t, serviceWithState("FAILED"))
		if _, err := getDeleteWaiter(context.Background(), client, testServiceUUID, 5*time.Second).WaitForStateContext(context.Background()); err != nil {
			t.Fatalf("a failed service never deprovisions, the waiter must not hang: %v", err)
		}
	})
	t.Run("not found counts as deleted", func(t *testing.T) {
		client, _ := newTestFabricClient(t, fakeResponse{http.StatusNotFound, `[{"errorCode":"EQ-3040016","errorMessage":"not found"}]`})
		if _, err := getDeleteWaiter(context.Background(), client, testServiceUUID, 5*time.Second).WaitForStateContext(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// ---- getService ----

func TestGetService_ToleratesMissingOrder(t *testing.T) {
	// the API omits "order" while a service is being provisioned
	body := `{
  "href": "h", "type": "SINGLE_IA", "uuid": "` + testServiceUUID + `", "name": "n", "state": "PROVISIONING",
  "change": {"href": "h"},
  "routingProtocol": {"type": "STATIC", "customerRoutes": [], "connections": []},
  "billing": {"type": "FIXED"},
  "account": {"accountNumber": "1"},
  "project": {"projectId": "p"},
  "changeLog": {"createdDateTime": "2026-04-23T15:19:07Z"},
  "useCase": "MAIN"
}`
	client, _ := newTestFabricClient(t, fakeResponse{http.StatusOK, body})

	svc, _, err := getService(context.Background(), client, testServiceUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.GetState() != fabricv4.INTERNETACCESSSERVICESTATE_PROVISIONING || svc.GetUuid() != testServiceUUID {
		t.Errorf("service = %+v", svc)
	}
}

func TestGetService_KeepsOtherErrors(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		client, _ := newTestFabricClient(t, fakeResponse{http.StatusNotFound, `[{"errorCode":"EQ-3040016","errorMessage":"not found"}]`})
		_, resp, err := getService(context.Background(), client, testServiceUUID)
		if err == nil || resp == nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected a 404 error, got resp=%v err=%v", resp, err)
		}
	})
	t.Run("another missing property", func(t *testing.T) {
		client, _ := newTestFabricClient(t, fakeResponse{http.StatusOK, `{"uuid":"x","order":{}}`})
		if _, _, err := getService(context.Background(), client, testServiceUUID); err == nil {
			t.Fatal("a body that is broken for another reason must still be an error")
		}
	})
}
