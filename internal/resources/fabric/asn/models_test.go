package asn

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
  "sort": [{"direction": "DESC", "property": "/changeLog/createdDateTime"}],
  "data": [
    {
      "href": "https://api.equinix.com/fabric/v4/asns/asn-1",
      "uuid": "asn-1",
      "state": "ATTACHED",
      "number": 4200000001,
      "rirRegion": "ARIN",
      "range": "32_BIT",
      "project": {"projectId": "proj-1"},
      "location": {"metroHref": "https://api.equinix.com/fabric/v4/metros/SV", "metroCode": "SV", "metroName": "Silicon Valley", "region": "AMER", "ibx": "SV1"},
      "ownershipChallengeDetails": {"type": "RPKI"},
      "change": {"href": "https://api.equinix.com/fabric/v4/asns/asn-1/changes/c-1"},
      "assets": [{
        "uuid": "asset-1",
        "announcingMetro": "SV",
        "service": {"href": "https://api.equinix.com/fabric/v4/internetAccessServices/svc-1", "type": "IA", "uuid": "svc-1"},
        "ipBlocks": [{"href": "https://api.equinix.com/fabric/v4/ipBlocks/blk-1", "type": "IPV4_IP_BLOCK", "uuid": "blk-1"}]
      }],
      "changeLog": {"createdDateTime": "2026-10-29T12:30:41.008Z", "updatedDateTime": "2026-10-30T12:30:41.008Z"}
    },
    {
      "href": "https://api.equinix.com/fabric/v4/asns/asn-2",
      "uuid": "asn-2",
      "state": "PENDING",
      "number": 432,
      "rirRegion": "ARIN",
      "project": {"projectId": "proj-1"},
      "assets": [],
      "changeLog": {"createdDateTime": "2026-10-29T12:30:41Z", "deletedDateTime": "2026-11-01T00:00:00Z"}
    }
  ]
}`

func sampleResponse(t *testing.T) *fabricv4.AsnSearchResponseBody {
	t.Helper()
	var resp fabricv4.AsnSearchResponseBody
	if err := json.Unmarshal([]byte(searchResponseJSON), &resp); err != nil {
		t.Fatalf("unmarshal sample response: %v", err)
	}
	return &resp
}

func TestBaseAsnModelParse_Full(t *testing.T) {
	ctx := context.Background()
	data := sampleResponse(t).GetData()

	var model BaseAsnModel
	if diags := model.parse(ctx, &data[0]); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}

	if model.UUID.ValueString() != "asn-1" || model.State.ValueString() != "ATTACHED" {
		t.Errorf("uuid/state = %v/%v", model.UUID, model.State)
	}
	if model.Number.ValueInt64() != 4200000001 {
		t.Errorf("number = %v, want 4200000001 (32-bit values must not overflow)", model.Number)
	}
	if model.RirRegion.ValueString() != "ARIN" || model.Range.ValueString() != "32_BIT" {
		t.Errorf("rir_region/range = %v/%v", model.RirRegion, model.Range)
	}

	project, _ := model.Project.ToPtr(ctx)
	if project.ProjectID.ValueString() != "proj-1" {
		t.Errorf("project_id = %v", project.ProjectID)
	}

	loc, _ := model.Location.ToPtr(ctx)
	if loc.MetroCode.ValueString() != "SV" || loc.MetroName.ValueString() != "Silicon Valley" || loc.Region.ValueString() != "AMER" {
		t.Errorf("location = %+v", loc)
	}

	challenge, _ := model.OwnershipChallengeDetails.ToPtr(ctx)
	if challenge.Type.ValueString() != "RPKI" {
		t.Errorf("ownership challenge type = %v", challenge.Type)
	}

	assets, _ := model.Assets.ToSlice(ctx)
	if len(assets) != 1 {
		t.Fatalf("assets = %d, want 1", len(assets))
	}
	service, _ := assets[0].Service.ToPtr(ctx)
	if service == nil || service.UUID.ValueString() != "svc-1" || service.Type.ValueString() != "IA" {
		t.Errorf("asset service = %+v", service)
	}
	if blocks, _ := assets[0].IpBlocks.ToSlice(ctx); len(blocks) != 1 || blocks[0].UUID.ValueString() != "blk-1" || blocks[0].Type.ValueString() != "IPV4_IP_BLOCK" {
		t.Errorf("asset ip blocks = %+v", blocks)
	}

	cl, _ := model.ChangeLog.ToPtr(ctx)
	if cl.CreatedDateTime.ValueString() != "2026-10-29T12:30:41.008Z" || !cl.DeletedDateTime.IsNull() {
		t.Errorf("change log = %+v", cl)
	}
}

func TestBaseAsnModelParse_Sparse(t *testing.T) {
	ctx := context.Background()
	data := sampleResponse(t).GetData()

	var model BaseAsnModel
	if diags := model.parse(ctx, &data[1]); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	if model.State.ValueString() != "PENDING" || model.Number.ValueInt64() != 432 {
		t.Errorf("state/number = %v/%v", model.State, model.Number)
	}
	if len(model.Assets.Elements()) != 0 {
		t.Errorf("assets should be empty, got %d", len(model.Assets.Elements()))
	}
	cl, _ := model.ChangeLog.ToPtr(ctx)
	if cl.DeletedDateTime.ValueString() != "2026-11-01T00:00:00.000Z" || !cl.UpdatedDateTime.IsNull() {
		t.Errorf("change log = %+v", cl)
	}
}

func TestBaseAsnModelParse_NoNumber(t *testing.T) {
	var model BaseAsnModel
	if diags := model.parse(context.Background(), &fabricv4.Asn{}); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	if !model.Number.IsNull() {
		t.Errorf("number should be null when absent, got %v", model.Number)
	}
}

func TestDataSourceAllAsnsModelParse(t *testing.T) {
	ctx := context.Background()
	var model DataSourceAllAsnsModel
	if diags := model.parse(ctx, sampleResponse(t)); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	if model.ID.ValueString() != "asn-1" {
		t.Errorf("id = %v, want first uuid", model.ID)
	}
	if got := len(model.Data.Elements()); got != 2 {
		t.Errorf("data = %d, want 2", got)
	}
	p, _ := model.Pagination.ToPtr(ctx)
	if p.Total.ValueInt32() != 2 || p.Limit.ValueInt32() != 20 || p.Next.ValueString() != "/next" {
		t.Errorf("pagination = %+v", p)
	}
}

func TestDataSourceAllAsnsModelParse_Empty(t *testing.T) {
	var model DataSourceAllAsnsModel
	if diags := model.parse(context.Background(), &fabricv4.AsnSearchResponseBody{}); !diags.HasError() {
		t.Error("expected an error for an empty search response")
	}
}

func TestDataSourceByIDModelParse(t *testing.T) {
	data := sampleResponse(t).GetData()
	var model DataSourceByIDModel
	if diags := model.parse(context.Background(), &data[0]); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	if model.ID.ValueString() != "asn-1" || model.AsnID.ValueString() != "asn-1" {
		t.Errorf("id/asn_id = %v/%v", model.ID, model.AsnID)
	}
}

// The parsed models must be assignable to state built from the data source schemas;
// this catches tfsdk tag / schema mismatches without needing the API.
func TestModelsMatchSchemas(t *testing.T) {
	ctx := context.Background()
	resp := sampleResponse(t)
	data := resp.GetData()

	single := dataSourceSingleAsnSchema(ctx)
	all := dataSourceAllAsnsSchema(ctx)
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

	var search DataSourceAllAsnsModel
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
