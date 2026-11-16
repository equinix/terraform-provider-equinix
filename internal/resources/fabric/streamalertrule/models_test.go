package streamalertrule

import (
	"context"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	testStreamID    = "11111111-1111-1111-1111-111111111111"
	testAlertRuleID = "33333333-3333-3333-3333-333333333333"
)

func TestDataSourceByIDsParseKeepsConfiguredIDs(t *testing.T) {
	ctx := context.Background()
	m := dataSourceStreamAlertRuleByIDsModel{
		StreamID:    types.StringValue(testStreamID),
		AlertRuleID: types.StringValue(testAlertRuleID),
	}
	rule := fabricv4.StreamAlertRule{}
	rule.SetUuid(testAlertRuleID)

	if diags := m.parse(ctx, &rule); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	if got := m.StreamID.ValueString(); got != testStreamID {
		t.Errorf("stream_id = %q, want %q", got, testStreamID)
	}
	if got := m.AlertRuleID.ValueString(); got != testAlertRuleID {
		t.Errorf("alert_rule_id = %q, want %q", got, testAlertRuleID)
	}
	if got := m.ID.ValueString(); got != testAlertRuleID {
		t.Errorf("id = %q, want %q", got, testAlertRuleID)
	}
}

func TestDataSourceAllParseKeepsConfiguredStreamID(t *testing.T) {
	ctx := context.Background()
	m := dataSourceAllStreamAlertRulesModel{StreamID: types.StringValue(testStreamID)}
	rule := fabricv4.StreamAlertRule{}
	rule.SetUuid(testAlertRuleID)
	resp := fabricv4.GetAllStreamAlertRuleResponse{Data: []fabricv4.StreamAlertRule{rule}}

	if diags := m.parse(ctx, &resp); diags.HasError() {
		t.Fatalf("parse: %v", diags)
	}
	if got := m.StreamID.ValueString(); got != testStreamID {
		t.Errorf("stream_id = %q, want %q", got, testStreamID)
	}
}
