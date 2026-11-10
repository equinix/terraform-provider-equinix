package precisiontime

import (
	"context"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestParseAccountOrganizationName(t *testing.T) {
	ctx := context.Background()
	account := fabricv4.SimplifiedAccount{}
	account.SetAccountName("my-account")
	account.SetOrganizationName("my-organization")

	obj, diags := parseAccount(ctx, account)
	if diags.HasError() {
		t.Fatalf("parseAccount: %v", diags)
	}
	var got accountModel
	if diags := obj.As(ctx, &got, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("As: %v", diags)
	}
	if v := got.AccountName.ValueString(); v != "my-account" {
		t.Errorf("account_name = %q, want %q", v, "my-account")
	}
	if v := got.OrganizationName.ValueString(); v != "my-organization" {
		t.Errorf("organization_name = %q, want %q", v, "my-organization")
	}
}

func TestParseOrderReferenceNumbers(t *testing.T) {
	ctx := context.Background()
	order := fabricv4.PrecisionTimeOrder{}
	order.SetPurchaseOrderNumber("PO-1")
	order.SetCustomerReferenceNumber("REF-1")
	order.SetOrderNumber("ORD-1")

	obj, diags := parseOrder(ctx, order)
	if diags.HasError() {
		t.Fatalf("parseOrder: %v", diags)
	}
	var got orderModel
	if diags := obj.As(ctx, &got, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("As: %v", diags)
	}
	if v := got.PurchaseOrderNumber.ValueString(); v != "PO-1" {
		t.Errorf("purchase_order_number = %q, want %q", v, "PO-1")
	}
	if v := got.CustomerReferenceNumber.ValueString(); v != "REF-1" {
		t.Errorf("customer_reference_number = %q, want %q", v, "REF-1")
	}
	if v := got.OrderNumber.ValueString(); v != "ORD-1" {
		t.Errorf("order_number = %q, want %q", v, "ORD-1")
	}
}
