package asn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	testASNUUID = "11111111-2222-3333-4444-555555555555"

	privateAsnJSON = `{
  "uuid": "11111111-2222-3333-4444-555555555555",
  "href": "https://api.equinix.com/fabric/v4/asns/11111111-2222-3333-4444-555555555555",
  "type": "PRIVATE_ASN",
  "state": "PROVISIONED",
  "number": 64529,
  "range": "16_BIT",
  "project": {"projectId": "proj-1"},
  "location": {"metroCode": "SV"},
  "assets": [],
  "changeLog": {"createdDateTime": "2026-11-02T18:27:57.976Z"}
}`

	publicAsnJSON = `{
  "uuid": "11111111-2222-3333-4444-555555555555",
  "type": "PUBLIC_ASN",
  "state": "PENDING",
  "number": 64497,
  "rirRegion": "RIPE_NCC",
  "project": {"projectId": "proj-1"},
  "ownershipChallengeDetails": {"type": "RPKI", "message": "eqx-rpki-challenge-64497"},
  "assets": [],
  "changeLog": {"createdDateTime": "2026-11-02T18:27:57.976Z"}
}`
)

func init() {
	waitDelay = time.Millisecond
	waitMinInterval = time.Millisecond
}

func unmarshalAsn(t *testing.T, raw string) *fabricv4.Asn {
	t.Helper()
	var asn fabricv4.Asn
	if err := json.Unmarshal([]byte(raw), &asn); err != nil {
		t.Fatalf("unmarshal ASN: %v", err)
	}
	return &asn
}

func TestResourceModelParse_Private(t *testing.T) {
	ctx := context.Background()
	var m ResourceModel
	if diags := m.parse(ctx, unmarshalAsn(t, privateAsnJSON)); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}

	if m.ID.ValueString() != testASNUUID || m.UUID.ValueString() != testASNUUID {
		t.Errorf("id/uuid = %v/%v", m.ID, m.UUID)
	}
	if m.Type.ValueString() != "PRIVATE_ASN" || m.State.ValueString() != "PROVISIONED" {
		t.Errorf("type/state = %v/%v", m.Type, m.State)
	}
	if m.Number.ValueInt64() != 64529 || m.Range.ValueString() != "16_BIT" {
		t.Errorf("number/range = %v/%v", m.Number, m.Range)
	}
	if !m.RirRegion.IsNull() {
		t.Errorf("rir_region should be null for a private ASN, got %v", m.RirRegion)
	}
	loc, _ := m.Location.ToPtr(ctx)
	if loc == nil || loc.MetroCode.ValueString() != "SV" {
		t.Errorf("location = %+v", loc)
	}
	if !m.OwnershipChallengeDetails.IsNull() {
		t.Errorf("ownership challenge should be null for a private ASN")
	}
}

func TestResourceModelParse_PublicKeepsSignature(t *testing.T) {
	ctx := context.Background()

	var m ResourceModel
	m.OwnershipChallengeDetails = fwtypes.NewObjectValueOf[ResourceChallengeModel](ctx, &ResourceChallengeModel{
		Type:      types.StringValue("RPKI"),
		Message:   types.StringUnknown(),
		Signature: types.StringValue("signed-by-user"),
	})
	if diags := m.parse(ctx, unmarshalAsn(t, publicAsnJSON)); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}

	if !m.Location.IsNull() || m.RirRegion.ValueString() != "RIPE_NCC" {
		t.Errorf("location/rir_region = %v/%v", m.Location, m.RirRegion)
	}
	challenge, _ := m.OwnershipChallengeDetails.ToPtr(ctx)
	if challenge == nil {
		t.Fatal("ownership challenge is null")
	}
	if challenge.Type.ValueString() != "RPKI" || challenge.Message.ValueString() != "eqx-rpki-challenge-64497" {
		t.Errorf("challenge = %+v", challenge)
	}
	if challenge.Signature.ValueString() != "signed-by-user" {
		t.Errorf("signature not preserved from the plan: %v", challenge.Signature)
	}
}

func TestResourceModelParse_PublicSignatureReturnedByAPI(t *testing.T) {
	ctx := context.Background()
	raw := strings.Replace(publicAsnJSON, `"message": "eqx-rpki-challenge-64497"`, `"message": "m", "signature": "from-api"`, 1)

	var m ResourceModel
	if diags := m.parse(ctx, unmarshalAsn(t, raw)); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	challenge, _ := m.OwnershipChallengeDetails.ToPtr(ctx)
	if challenge.Signature.ValueString() != "from-api" {
		t.Errorf("signature = %v, want from-api", challenge.Signature)
	}
}

func TestResourceModelMatchesSchema(t *testing.T) {
	ctx := context.Background()
	s := resourceSchema(ctx)
	if diags := s.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema is invalid: %v", diags)
	}

	for name, raw := range map[string]string{"private": privateAsnJSON, "public": publicAsnJSON} {
		var m ResourceModel
		m.Timeouts = timeouts.Value{Object: types.ObjectNull(s.Attributes["timeouts"].GetType().(timeouts.Type).ObjectType.AttrTypes)}
		if diags := m.parse(ctx, unmarshalAsn(t, raw)); diags.HasError() {
			t.Fatalf("%s parse: %v", name, diags)
		}
		state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
		if diags := state.Set(ctx, &m); diags.HasError() {
			t.Errorf("%s model does not match schema: %v", name, diags)
		}
	}
}

func TestBuildCreateRequest_Private(t *testing.T) {
	ctx := context.Background()
	plan := &ResourceModel{
		Type:                      types.StringValue("PRIVATE_ASN"),
		Project:                   fwtypes.NewObjectValueOf[ProjectModel](ctx, &ProjectModel{ProjectID: types.StringValue("proj-1")}),
		Number:                    types.Int64Unknown(),
		RirRegion:                 types.StringUnknown(),
		Range:                     types.StringValue("16_BIT"),
		Location:                  fwtypes.NewObjectValueOf[ResourceLocationModel](ctx, &ResourceLocationModel{MetroCode: types.StringValue("SV")}),
		OwnershipChallengeDetails: fwtypes.NewObjectValueOfUnknown[ResourceChallengeModel](ctx),
	}

	req, err := buildCreateRequest(ctx, plan)
	if err != nil {
		t.Fatalf("buildCreateRequest: %v", err)
	}
	body := marshalToMap(t, req)

	if body["type"] != "PRIVATE_ASN" || body["range"] != "16_BIT" {
		t.Errorf("body = %v", body)
	}
	if loc, _ := body["location"].(map[string]any); loc["metroCode"] != "SV" {
		t.Errorf("location = %v", body["location"])
	}
	if proj, _ := body["project"].(map[string]any); proj["projectId"] != "proj-1" {
		t.Errorf("project = %v", body["project"])
	}
	for _, key := range []string{"number", "rirRegion", "ownershipChallengeDetails"} {
		if _, ok := body[key]; ok {
			t.Errorf("%s must not be sent for a private ASN: %v", key, body)
		}
	}
}

func TestBuildCreateRequest_Public(t *testing.T) {
	ctx := context.Background()
	plan := &ResourceModel{
		Type:      types.StringValue("PUBLIC_ASN"),
		Project:   fwtypes.NewObjectValueOf[ProjectModel](ctx, &ProjectModel{ProjectID: types.StringValue("proj-1")}),
		Number:    types.Int64Value(64497),
		RirRegion: types.StringValue("RIPE_NCC"),
		Range:     types.StringUnknown(),
		Location:  fwtypes.NewObjectValueOfUnknown[ResourceLocationModel](ctx),
		OwnershipChallengeDetails: fwtypes.NewObjectValueOf[ResourceChallengeModel](ctx, &ResourceChallengeModel{
			Type:      types.StringValue("RPKI"),
			Message:   types.StringUnknown(),
			Signature: types.StringNull(),
		}),
	}

	req, err := buildCreateRequest(ctx, plan)
	if err != nil {
		t.Fatalf("buildCreateRequest: %v", err)
	}
	body := marshalToMap(t, req)

	if body["type"] != "PUBLIC_ASN" || body["rirRegion"] != "RIPE_NCC" || body["number"] != float64(64497) {
		t.Errorf("body = %v", body)
	}
	if challenge, _ := body["ownershipChallengeDetails"].(map[string]any); challenge["type"] != "RPKI" {
		t.Errorf("ownership challenge = %v", body["ownershipChallengeDetails"])
	}
	if _, ok := body["signature"]; ok {
		t.Errorf("signature must not be sent at creation: %v", body)
	}
}

func TestBuildCreateRequest_MissingProject(t *testing.T) {
	ctx := context.Background()
	if _, err := buildCreateRequest(ctx, &ResourceModel{Project: fwtypes.NewObjectValueOfNull[ProjectModel](ctx)}); err == nil {
		t.Error("expected an error without a project")
	}
}

func marshalToMap(t *testing.T, v any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestSignatureFromPlan(t *testing.T) {
	ctx := context.Background()
	with := func(sig types.String) *ResourceModel {
		return &ResourceModel{OwnershipChallengeDetails: fwtypes.NewObjectValueOf[ResourceChallengeModel](ctx, &ResourceChallengeModel{
			Type: types.StringValue("RPKI"), Message: types.StringNull(), Signature: sig,
		})}
	}
	tests := map[string]struct {
		model *ResourceModel
		want  string
	}{
		"set":     {with(types.StringValue("sig")), "sig"},
		"null":    {with(types.StringNull()), ""},
		"unknown": {with(types.StringUnknown()), ""},
		"no challenge": {&ResourceModel{
			OwnershipChallengeDetails: fwtypes.NewObjectValueOfNull[ResourceChallengeModel](ctx),
		}, ""},
	}
	for name, tc := range tests {
		if got := signatureFromPlan(ctx, tc.model); got != tc.want {
			t.Errorf("%s: signature = %q, want %q", name, got, tc.want)
		}
	}
}

func TestCreateTargetStates(t *testing.T) {
	if got := createTargetStates(false); len(got) != 2 {
		t.Errorf("without pending: %v", got)
	}
	got := createTargetStates(true)
	if len(got) != 3 || got[2] != string(fabricv4.ASNSTATE_PENDING) {
		t.Errorf("with pending: %v", got)
	}
}

// ---- API interaction against a fake server ----

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

func asnWithState(state string) fakeResponse {
	return fakeResponse{http.StatusOK, fmt.Sprintf(`{"uuid":%q,"state":%q}`, testASNUUID, state)}
}

func TestCreateWaiter_ReachesProvisioned(t *testing.T) {
	client, calls := newTestFabricClient(t, asnWithState("REQUESTED"), asnWithState("PROVISIONING"), asnWithState("PROVISIONED"))

	got, err := getWaiter(context.Background(), client, testASNUUID, 5*time.Second, createPendingStates(), createTargetStates(false)).WaitForStateContext(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state := got.(*fabricv4.Asn).GetState(); state != fabricv4.ASNSTATE_PROVISIONED {
		t.Errorf("state = %s", state)
	}
	if calls.Load() != 3 {
		t.Errorf("calls = %d, want 3", calls.Load())
	}
}

func TestCreateWaiter_PendingIsATargetForPublicAsn(t *testing.T) {
	client, _ := newTestFabricClient(t, asnWithState("REQUESTED"), asnWithState("PENDING"))

	if _, err := getWaiter(context.Background(), client, testASNUUID, 5*time.Second, createPendingStates(), createTargetStates(true)).WaitForStateContext(context.Background()); err != nil {
		t.Fatalf("PENDING should be accepted while waiting for a signature: %v", err)
	}
}

func TestCreateWaiter_FailedStateIsAnError(t *testing.T) {
	client, _ := newTestFabricClient(t, asnWithState("REQUESTED"), asnWithState("PROVISIONING_FAILED"))

	if _, err := getWaiter(context.Background(), client, testASNUUID, 5*time.Second, createPendingStates(), createTargetStates(false)).WaitForStateContext(context.Background()); err == nil {
		t.Fatal("expected an error for PROVISIONING_FAILED")
	}
}

func TestDeleteWaiter(t *testing.T) {
	t.Run("deprovisioned", func(t *testing.T) {
		client, _ := newTestFabricClient(t, asnWithState("PROVISIONED"), asnWithState("DEPROVISIONING"), asnWithState("DEPROVISIONED"))
		if _, err := getDeleteWaiter(context.Background(), client, testASNUUID, 5*time.Second).WaitForStateContext(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("not found counts as deleted", func(t *testing.T) {
		client, _ := newTestFabricClient(t, fakeResponse{http.StatusNotFound, `[{"errorCode":"EQ-3040016","errorMessage":"not found"}]`})
		if _, err := getDeleteWaiter(context.Background(), client, testASNUUID, 5*time.Second).WaitForStateContext(context.Background()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestPatchSignature(t *testing.T) {
	var method, path, body string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		method, path, body = r.Method, r.URL.Path, string(raw)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, `{"uuid":%q,"state":"PROVISIONING"}`, testASNUUID)
	}))
	t.Cleanup(server.Close)
	cfg := fabricv4.NewConfiguration()
	cfg.Servers = fabricv4.ServerConfigurations{{URL: server.URL}}
	cfg.HTTPClient = server.Client()

	if err := patchSignature(context.Background(), fabricv4.NewAPIClient(cfg), testASNUUID, "sig=="); err != nil {
		t.Fatalf("patchSignature: %v", err)
	}
	if method != http.MethodPatch || !strings.HasSuffix(path, "/asns/"+testASNUUID) {
		t.Errorf("request = %s %s", method, path)
	}
	var ops []map[string]any
	if err := json.Unmarshal([]byte(body), &ops); err != nil || len(ops) != 1 {
		t.Fatalf("body = %s (%v)", body, err)
	}
	if ops[0]["op"] != "replace" || ops[0]["path"] != "/ownershipChallengeDetails/signature" || ops[0]["value"] != "sig==" {
		t.Errorf("patch operation = %v", ops[0])
	}
}

func TestPatchSignature_APIError(t *testing.T) {
	client, _ := newTestFabricClient(t, fakeResponse{http.StatusBadRequest, `[{"errorCode":"EQ-3000001","errorMessage":"invalid signature"}]`})
	if err := patchSignature(context.Background(), client, testASNUUID, "bad"); err == nil {
		t.Error("expected an error")
	}
}
