package internetaccess

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const searchResponseJSON = `{
  "pagination": {"offset": 0, "limit": 20, "total": 2, "next": "/next"},
  "sort": [{"direction": "ASC", "property": "/name"}],
  "data": [
    {
      "href": "https://api.equinix.com/fabric/v4/internetAccessServices/svc-1",
      "type": "SINGLE_IA",
      "uuid": "svc-1",
      "name": "first",
      "bandwidth": 15,
      "bandwidthCommit": 5,
      "state": "PROVISIONED",
      "change": {"href": "https://api.equinix.com/fabric/v4/internetAccessServices/svc-1/cloudEvents"},
      "locations": [{"metroHref": "https://api.equinix.com/fabric/v4/metros/SV", "metroCode": "SV", "region": "AMER", "ibx": "SV3"}],
      "routingProtocol": {
        "type": "DIRECT",
        "customerRoutes": [{"ipBlock": {"href": "https://api.equinix.com/fabric/v4/ipblocks/blk-1", "uuid": "blk-1"}}],
        "connections": [{
          "href": "https://api.equinix.com/fabric/v4/connections/conn-1",
          "uuid": "conn-1",
          "peeringIpv4": {"prefix": "137.174.168.34/31", "prefixLength": 31, "equinixPeerIp": "137.174.168.34", "customerPeerIp": "137.174.168.35"},
          "peeringIpv6": {"prefix": "2001:db8::/127", "prefixLength": 127, "equinixPeerIp": "2001:db8::1"}
        }]
      },
      "billing": {"type": "FIXED", "enabled": true, "startDate": "2026-04-23T15:19:07Z"},
      "account": {"accountNumber": "709976", "href": "https://api.equinix.com/billing/v2/billingAccounts/accountNumber/709976"},
      "project": {"projectId": "proj-1"},
      "order": {"purchaseOrderNumber": "PO-1", "termLength": 12, "contractedBandwidth": 100, "orderId": "ord-1"},
      "changeLog": {"createdDateTime": "2026-04-23T15:19:07.202Z", "updatedDateTime": "2026-04-24T15:19:07.202Z"},
      "useCase": "MAIN"
    },
    {
      "href": "https://api.equinix.com/fabric/v4/internetAccessServices/svc-2",
      "type": "DUAL_IA",
      "uuid": "svc-2",
      "name": "second",
      "state": "FAILED",
      "change": {"href": "https://api.equinix.com/fabric/v4/internetAccessServices/svc-2/cloudEvents"},
      "routingProtocol": {"type": "STATIC", "customerRoutes": [], "connections": [{"href": "h", "uuid": "conn-2"}]},
      "billing": {"type": "USAGE_BASED"},
      "account": {"accountNumber": "709976"},
      "project": {"projectId": "proj-1"},
      "order": {},
      "changeLog": {"createdDateTime": "2026-05-04T18:16:07Z", "deletedDateTime": "2026-05-15T23:45:10Z"},
      "useCase": "MANAGEMENT_ACCESS"
    }
  ]
}`

func sampleServices(t *testing.T) *fabricv4.InternetAccessServices {
	t.Helper()
	var resp fabricv4.InternetAccessServices
	if err := json.Unmarshal([]byte(searchResponseJSON), &resp); err != nil {
		t.Fatalf("unmarshal sample response: %v", err)
	}
	return &resp
}

func TestBaseInternetAccessModelParse(t *testing.T) {
	ctx := context.Background()
	resp := sampleServices(t)
	data := resp.GetData()

	var full BaseInternetAccessModel
	if diags := full.parse(ctx, &data[0]); diags.HasError() {
		t.Fatalf("parse returned errors: %v", diags)
	}
	if got := full.UUID.ValueString(); got != "svc-1" {
		t.Errorf("uuid = %q, want svc-1", got)
	}
	if full.Bandwidth.ValueInt32() != 15 || full.BandwidthCommit.ValueInt32() != 5 {
		t.Errorf("bandwidth = %v/%v, want 15/5", full.Bandwidth, full.BandwidthCommit)
	}
	if full.State.ValueString() != "PROVISIONED" || full.UseCase.ValueString() != "MAIN" {
		t.Errorf("state/use_case = %v/%v", full.State, full.UseCase)
	}

	rp, _ := full.RoutingProtocol.ToPtr(ctx)
	if rp.Type.ValueString() != "DIRECT" {
		t.Errorf("routing protocol type = %v, want DIRECT", rp.Type)
	}
	conns, _ := rp.Connections.ToSlice(ctx)
	if len(conns) != 1 || conns[0].UUID.ValueString() != "conn-1" {
		t.Fatalf("connections = %+v", conns)
	}
	v4, _ := conns[0].PeeringIpv4.ToPtr(ctx)
	if v4 == nil || v4.PrefixLength.ValueInt32() != 31 || v4.CustomerPeerIp.ValueString() != "137.174.168.35" {
		t.Errorf("peering ipv4 = %+v", v4)
	}
	routes, _ := rp.CustomerRoutes.ToSlice(ctx)
	if len(routes) != 1 {
		t.Fatalf("customer routes = %d, want 1", len(routes))
	}
	block, _ := routes[0].IpBlock.ToPtr(ctx)
	if block == nil || block.UUID.ValueString() != "blk-1" {
		t.Errorf("ip block = %+v", block)
	}

	billing, _ := full.Billing.ToPtr(ctx)
	if !billing.Enabled.ValueBool() || billing.StartDate.IsNull() {
		t.Errorf("billing = %+v", billing)
	}
	order, _ := full.Order.ToPtr(ctx)
	if order.TermLength.ValueInt32() != 12 || order.ContractedBandwidth.ValueInt32() != 100 {
		t.Errorf("order = %+v", order)
	}
	cl, _ := full.ChangeLog.ToPtr(ctx)
	if cl.CreatedDateTime.IsNull() || cl.UpdatedDateTime.IsNull() || !cl.DeletedDateTime.IsNull() {
		t.Errorf("change log = %+v", cl)
	}

	// second service has most optional fields absent
	var sparse BaseInternetAccessModel
	if diags := sparse.parse(ctx, &data[1]); diags.HasError() {
		t.Fatalf("parse returned errors: %v", diags)
	}
	if !sparse.Bandwidth.IsNull() || !sparse.BandwidthCommit.IsNull() {
		t.Errorf("absent bandwidth should be null, got %v/%v", sparse.Bandwidth, sparse.BandwidthCommit)
	}
	sparseRP, _ := sparse.RoutingProtocol.ToPtr(ctx)
	sparseConns, _ := sparseRP.Connections.ToSlice(ctx)
	if len(sparseConns) != 1 || !sparseConns[0].PeeringIpv4.IsNull() || !sparseConns[0].PeeringIpv6.IsNull() {
		t.Errorf("sparse connection should have null peering, got %+v", sparseConns)
	}
	sparseBilling, _ := sparse.Billing.ToPtr(ctx)
	if !sparseBilling.Enabled.IsNull() || !sparseBilling.StartDate.IsNull() {
		t.Errorf("sparse billing = %+v", sparseBilling)
	}
	sparseCL, _ := sparse.ChangeLog.ToPtr(ctx)
	if sparseCL.DeletedDateTime.IsNull() || !sparseCL.UpdatedDateTime.IsNull() {
		t.Errorf("sparse change log = %+v", sparseCL)
	}
}

func TestDataSourceAllInternetAccessModelParse(t *testing.T) {
	ctx := context.Background()

	var model DataSourceAllInternetAccessModel
	if diags := model.parse(ctx, sampleServices(t)); diags.HasError() {
		t.Fatalf("parse returned errors: %v", diags)
	}
	if model.ID.ValueString() != "svc-1" {
		t.Errorf("id = %q, want svc-1", model.ID.ValueString())
	}
	items, _ := model.Data.ToSlice(ctx)
	if len(items) != 2 {
		t.Fatalf("data length = %d, want 2", len(items))
	}
	pagination, _ := model.Pagination.ToPtr(ctx)
	if pagination.Total.ValueInt32() != 2 || pagination.Limit.ValueInt32() != 20 || pagination.Next.ValueString() != "/next" {
		t.Errorf("pagination = %+v", pagination)
	}

	var empty DataSourceAllInternetAccessModel
	if diags := empty.parse(ctx, &fabricv4.InternetAccessServices{}); !diags.HasError() {
		t.Error("parse of an empty response should return an error")
	}
}

func TestDataSourceByIDModelParse(t *testing.T) {
	ctx := context.Background()
	data := sampleServices(t).GetData()

	var model DataSourceByIDModel
	if diags := model.parse(ctx, &data[0]); diags.HasError() {
		t.Fatalf("parse returned errors: %v", diags)
	}
	if model.ID.ValueString() != "svc-1" || model.InternetAccessServiceID.ValueString() != "svc-1" {
		t.Errorf("id/internet_access_service_id = %v/%v", model.ID, model.InternetAccessServiceID)
	}
}

// The parsed models must be assignable to state built from the data source schemas;
// this catches tfsdk tag / schema mismatches without needing the API.
func TestModelsMatchSchemas(t *testing.T) {
	ctx := context.Background()
	resp := sampleServices(t)
	data := resp.GetData()

	single := dataSourceSingleInternetAccessSchema(ctx)
	all := dataSourceAllInternetAccessSchema(ctx)
	for name, s := range map[string]schema.Schema{"single": single, "all": all} {
		if diags := s.ValidateImplementation(ctx); diags.HasError() {
			t.Fatalf("%s schema is invalid: %v", name, diags)
		}
	}

	var byID DataSourceByIDModel
	if diags := byID.parse(ctx, &data[0]); diags.HasError() {
		t.Fatalf("parse by id: %v", diags)
	}
	byIDState := tfsdk.State{Schema: single, Raw: tftypes.NewValue(single.Type().TerraformType(ctx), nil)}
	if diags := byIDState.Set(ctx, &byID); diags.HasError() {
		t.Errorf("by id model does not match schema: %v", diags)
	}

	var search DataSourceAllInternetAccessModel
	if diags := search.parse(ctx, resp); diags.HasError() {
		t.Fatalf("parse search: %v", diags)
	}
	search.Filter = types.ObjectNull(map[string]attr.Type{
		"property": types.StringType,
		"operator": types.StringType,
		"values":   types.ListType{ElemType: types.StringType},
	})
	search.Sort = fwtypes.NewObjectValueOfNull[SortModel](ctx)
	searchState := tfsdk.State{Schema: all, Raw: tftypes.NewValue(all.Type().TerraformType(ctx), nil)}
	if diags := searchState.Set(ctx, &search); diags.HasError() {
		t.Errorf("search model does not match schema: %v", diags)
	}
}
