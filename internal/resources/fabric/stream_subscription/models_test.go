package streamsubscription

import (
	"context"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	testStreamID       = "11111111-1111-1111-1111-111111111111"
	testSubscriptionID = "22222222-2222-2222-2222-222222222222"
)

func TestDataSourceByIDsParseKeepsConfiguredIDs(t *testing.T) {
	ctx := context.Background()
	m := dataSourceByIDsModel{
		StreamID:       types.StringValue(testStreamID),
		SubscriptionID: types.StringValue(testSubscriptionID),
	}
	sub := fabricv4.StreamSubscription{}
	sub.SetUuid(testSubscriptionID)

	if diags := m.parse(ctx, &sub); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	if got := m.StreamID.ValueString(); got != testStreamID {
		t.Errorf("stream_id = %q, want %q", got, testStreamID)
	}
	if got := m.SubscriptionID.ValueString(); got != testSubscriptionID {
		t.Errorf("subscription_id = %q, want %q", got, testSubscriptionID)
	}
	if got := m.ID.ValueString(); got != testSubscriptionID {
		t.Errorf("id = %q, want %q", got, testSubscriptionID)
	}
}

func TestDataSourceAllParseKeepsConfiguredStreamID(t *testing.T) {
	ctx := context.Background()
	m := dataSourceAll{StreamID: types.StringValue(testStreamID)}
	sub := fabricv4.StreamSubscription{}
	sub.SetUuid(testSubscriptionID)
	resp := fabricv4.GetAllStreamSubscriptionResponse{Data: []fabricv4.StreamSubscription{sub}}

	if diags := m.parse(ctx, &resp); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	if got := m.StreamID.ValueString(); got != testStreamID {
		t.Errorf("stream_id = %q, want %q", got, testStreamID)
	}
}
