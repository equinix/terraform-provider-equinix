package internetaccess

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ResourceModel struct {
	ID              types.String                                   `tfsdk:"id"`
	Timeouts        timeouts.Value                                 `tfsdk:"timeouts"`
	UUID            types.String                                   `tfsdk:"uuid"`
	Href            types.String                                   `tfsdk:"href"`
	Type            types.String                                   `tfsdk:"type"`
	Name            types.String                                   `tfsdk:"name"`
	Bandwidth       types.Int32                                    `tfsdk:"bandwidth"`
	BandwidthCommit types.Int32                                    `tfsdk:"bandwidth_commit"`
	State           types.String                                   `tfsdk:"state"`
	UseCase         types.String                                   `tfsdk:"use_case"`
	Change          fwtypes.ObjectValueOf[ChangeModel]             `tfsdk:"change"`
	Locations       fwtypes.ListNestedObjectValueOf[LocationModel] `tfsdk:"locations"`
	RoutingProtocol fwtypes.ObjectValueOf[ResourceRoutingModel]    `tfsdk:"routing_protocol"`
	Billing         fwtypes.ObjectValueOf[ResourceBillingModel]    `tfsdk:"billing"`
	Account         fwtypes.ObjectValueOf[ResourceAccountModel]    `tfsdk:"account"`
	Project         fwtypes.ObjectValueOf[ProjectModel]            `tfsdk:"project"`
	Order           fwtypes.ObjectValueOf[ResourceOrderModel]      `tfsdk:"order"`
	ChangeLog       fwtypes.ObjectValueOf[ChangeLogModel]          `tfsdk:"change_log"`
}

type ResourceBillingModel struct {
	Type      types.String `tfsdk:"type"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	StartDate types.String `tfsdk:"start_date"`
}

type ResourceAccountModel struct {
	AccountNumber types.String `tfsdk:"account_number"`
	Href          types.String `tfsdk:"href"`
}

type ResourceOrderModel struct {
	PurchaseOrderNumber     types.String `tfsdk:"purchase_order_number"`
	CustomerReferenceNumber types.String `tfsdk:"customer_reference_number"`
	BillingTier             types.String `tfsdk:"billing_tier"`
	OrderID                 types.String `tfsdk:"order_id"`
	OrderNumber             types.String `tfsdk:"order_number"`
	TermLength              types.Int32  `tfsdk:"term_length"`
	ContractedBandwidth     types.Int32  `tfsdk:"contracted_bandwidth"`
	Href                    types.String `tfsdk:"href"`
}

type ResourceRoutingModel struct {
	Type             types.String                                             `tfsdk:"type"`
	CustomerRoutes   fwtypes.ListNestedObjectValueOf[CustomerRouteModel]      `tfsdk:"customer_routes"`
	Connections      fwtypes.ListNestedObjectValueOf[ResourceConnectionModel] `tfsdk:"connections"`
	ExportPolicy     types.String                                             `tfsdk:"export_policy"`
	Asn              fwtypes.ObjectValueOf[ResourceAsnModel]                  `tfsdk:"asn"`
	CustomerAsn      types.Int64                                              `tfsdk:"customer_asn"`
	CustomerAsnRange types.String                                             `tfsdk:"customer_asn_range"`
	BgpAuthKey       types.String                                             `tfsdk:"bgp_auth_key"`
}

type ResourceAsnModel struct {
	UUID types.String `tfsdk:"uuid"`
}

type ResourceConnectionModel struct {
	UUID        types.String                        `tfsdk:"uuid"`
	Href        types.String                        `tfsdk:"href"`
	PeeringIpv4 fwtypes.ObjectValueOf[PeeringModel] `tfsdk:"peering_ipv4"`
	PeeringIpv6 fwtypes.ObjectValueOf[PeeringModel] `tfsdk:"peering_ipv6"`
}

func (m *ResourceModel) parse(ctx context.Context, svc *fabricv4.InternetAccessService) diag.Diagnostics {
	var diags diag.Diagnostics

	var base BaseInternetAccessModel
	diags.Append(base.parse(ctx, svc)...)
	if diags.HasError() {
		return diags
	}

	m.ID = base.UUID
	m.UUID = base.UUID
	m.Href = base.Href
	m.Type = base.Type
	m.Name = base.Name
	m.Bandwidth = base.Bandwidth
	m.BandwidthCommit = base.BandwidthCommit
	m.State = base.State
	m.UseCase = base.UseCase
	m.Change = base.Change
	m.Locations = base.Locations
	m.Project = base.Project
	m.ChangeLog = base.ChangeLog

	var billing BillingModel
	diags.Append(base.Billing.As(ctx, &billing, basetypes.ObjectAsOptions{})...)
	m.Billing = fwtypes.NewObjectValueOf[ResourceBillingModel](ctx, &ResourceBillingModel{
		Type: billing.Type, Enabled: billing.Enabled, StartDate: billing.StartDate,
	})

	var account AccountModel
	diags.Append(base.Account.As(ctx, &account, basetypes.ObjectAsOptions{})...)
	m.Account = fwtypes.NewObjectValueOf[ResourceAccountModel](ctx, &ResourceAccountModel{
		AccountNumber: account.AccountNumber, Href: account.Href,
	})

	var order OrderModel
	diags.Append(base.Order.As(ctx, &order, basetypes.ObjectAsOptions{})...)
	m.Order = fwtypes.NewObjectValueOf[ResourceOrderModel](ctx, &ResourceOrderModel{
		PurchaseOrderNumber:     order.PurchaseOrderNumber,
		CustomerReferenceNumber: order.CustomerReferenceNumber,
		BillingTier:             order.BillingTier,
		OrderID:                 order.OrderID,
		OrderNumber:             order.OrderNumber,
		TermLength:              order.TermLength,
		ContractedBandwidth:     order.ContractedBandwidth,
		Href:                    order.Href,
	})

	previous := ResourceRoutingModel{
		ExportPolicy:     types.StringNull(),
		CustomerAsn:      types.Int64Null(),
		CustomerAsnRange: types.StringNull(),
		BgpAuthKey:       types.StringNull(),
		Asn:              fwtypes.NewObjectValueOfNull[ResourceAsnModel](ctx),
	}
	if !m.RoutingProtocol.IsNull() && !m.RoutingProtocol.IsUnknown() {
		diags.Append(m.RoutingProtocol.As(ctx, &previous, basetypes.ObjectAsOptions{})...)
	}

	var apiRouting RoutingProtocolModel
	diags.Append(base.RoutingProtocol.As(ctx, &apiRouting, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	apiConnections, d := apiRouting.Connections.ToSlice(ctx)
	diags.Append(d...)
	connections := make([]*ResourceConnectionModel, len(apiConnections))
	for i, c := range apiConnections {
		connections[i] = &ResourceConnectionModel{
			UUID: c.UUID, Href: c.Href, PeeringIpv4: c.PeeringIpv4, PeeringIpv6: c.PeeringIpv6,
		}
	}

	m.RoutingProtocol = fwtypes.NewObjectValueOf[ResourceRoutingModel](ctx, &ResourceRoutingModel{
		Type:             apiRouting.Type,
		CustomerRoutes:   apiRouting.CustomerRoutes,
		Connections:      fwtypes.NewListNestedObjectValueOfSlice[ResourceConnectionModel](ctx, connections),
		ExportPolicy:     previous.ExportPolicy,
		Asn:              previous.Asn,
		CustomerAsn:      previous.CustomerAsn,
		CustomerAsnRange: previous.CustomerAsnRange,
		BgpAuthKey:       previous.BgpAuthKey,
	})

	return diags
}

func buildCreateRequest(ctx context.Context, plan ResourceModel) (fabricv4.InternetAccessPostRequest, error) {
	var routing ResourceRoutingModel
	if diags := plan.RoutingProtocol.As(ctx, &routing, basetypes.ObjectAsOptions{}); diags.HasError() {
		return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading routing_protocol from plan")
	}
	var billing ResourceBillingModel
	if diags := plan.Billing.As(ctx, &billing, basetypes.ObjectAsOptions{}); diags.HasError() {
		return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading billing from plan")
	}
	var account ResourceAccountModel
	if diags := plan.Account.As(ctx, &account, basetypes.ObjectAsOptions{}); diags.HasError() {
		return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading account from plan")
	}
	var project ProjectModel
	if diags := plan.Project.As(ctx, &project, basetypes.ObjectAsOptions{}); diags.HasError() {
		return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading project from plan")
	}

	routingReq := fabricv4.InternetAccessRoutingProtocolRequest{
		Type:                 fabricv4.InternetAccessRoutingProtocolType(routing.Type.ValueString()),
		CustomerRoutes:       []fabricv4.InternetAccessCustomerRouteRequest{},
		AdditionalProperties: map[string]interface{}{},
	}

	if !routing.CustomerRoutes.IsNull() && !routing.CustomerRoutes.IsUnknown() {
		routes, diags := routing.CustomerRoutes.ToSlice(ctx)
		if diags.HasError() {
			return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading customer_routes from plan")
		}
		for _, route := range routes {
			var block IpBlockModel
			if diags := route.IpBlock.As(ctx, &block, basetypes.ObjectAsOptions{}); diags.HasError() {
				return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading customer route ip_block from plan")
			}
			routingReq.CustomerRoutes = append(routingReq.CustomerRoutes, fabricv4.InternetAccessCustomerRouteRequest{
				IpBlock: fabricv4.InternetAccessIpBlockRequest{Uuid: block.UUID.ValueString()},
			})
		}
	}

	connections, diags := routing.Connections.ToSlice(ctx)
	if diags.HasError() {
		return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading connections from plan")
	}
	connReqs := make([]map[string]interface{}, 0, len(connections))
	for _, conn := range connections {
		item := map[string]interface{}{"uuid": conn.UUID.ValueString()}
		if v4 := peeringRequest(ctx, conn.PeeringIpv4); v4 != nil {
			item["peeringIpv4"] = v4
		}
		if v6 := peeringRequest(ctx, conn.PeeringIpv6); v6 != nil {
			item["peeringIpv6"] = v6
		}
		connReqs = append(connReqs, item)
	}
	routingReq.AdditionalProperties["connections"] = connReqs

	if routing.Type.ValueString() == string(fabricv4.INTERNETACCESSROUTINGPROTOCOLTYPE_BGP) {
		if !routing.ExportPolicy.IsNull() && !routing.ExportPolicy.IsUnknown() {
			routingReq.AdditionalProperties["exportPolicy"] = routing.ExportPolicy.ValueString()
		}
		if !routing.Asn.IsNull() && !routing.Asn.IsUnknown() {
			var asn ResourceAsnModel
			if diags := routing.Asn.As(ctx, &asn, basetypes.ObjectAsOptions{}); diags.HasError() {
				return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading asn from plan")
			}
			routingReq.AdditionalProperties["asn"] = map[string]interface{}{"uuid": asn.UUID.ValueString()}
		}
		if !routing.CustomerAsn.IsNull() && !routing.CustomerAsn.IsUnknown() {
			routingReq.AdditionalProperties["customerAsn"] = routing.CustomerAsn.ValueInt64()
		}
		if !routing.CustomerAsnRange.IsNull() && !routing.CustomerAsnRange.IsUnknown() {
			routingReq.AdditionalProperties["customerAsnRange"] = routing.CustomerAsnRange.ValueString()
		}
		if !routing.BgpAuthKey.IsNull() && !routing.BgpAuthKey.IsUnknown() {
			routingReq.AdditionalProperties["bgpAuthKey"] = routing.BgpAuthKey.ValueString()
		}
	}

	req := fabricv4.InternetAccessPostRequest{
		Type:            fabricv4.InternetAccessServiceType(plan.Type.ValueString()),
		Name:            plan.Name.ValueString(),
		RoutingProtocol: routingReq,
		Billing: fabricv4.InternetAccessPostRequestBilling{
			Type: fabricv4.InternetAccessBillingType(billing.Type.ValueString()),
		},
		Project: fabricv4.Project{ProjectId: project.ProjectID.ValueString()},
		Account: fabricv4.InternetAccessAccount{AccountNumber: account.AccountNumber.ValueString()},
	}
	if !plan.Bandwidth.IsNull() && !plan.Bandwidth.IsUnknown() {
		req.SetBandwidth(plan.Bandwidth.ValueInt32())
	}
	if !plan.BandwidthCommit.IsNull() && !plan.BandwidthCommit.IsUnknown() {
		req.SetBandwidthCommit(plan.BandwidthCommit.ValueInt32())
	}
	if !plan.Order.IsNull() && !plan.Order.IsUnknown() {
		var order ResourceOrderModel
		if diags := plan.Order.As(ctx, &order, basetypes.ObjectAsOptions{}); diags.HasError() {
			return fabricv4.InternetAccessPostRequest{}, fmt.Errorf("failed reading order from plan")
		}
		if !order.PurchaseOrderNumber.IsNull() && !order.PurchaseOrderNumber.IsUnknown() {
			orderReq := fabricv4.InternetAccessOrderRequest{}
			orderReq.SetPurchaseOrderNumber(order.PurchaseOrderNumber.ValueString())
			req.SetOrder(orderReq)
		}
	}

	return req, nil
}

func peeringRequest(ctx context.Context, p fwtypes.ObjectValueOf[PeeringModel]) map[string]interface{} {
	if p.IsNull() || p.IsUnknown() {
		return nil
	}
	var peering PeeringModel
	if diags := p.As(ctx, &peering, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil
	}
	out := map[string]interface{}{}
	if !peering.EquinixPeerIp.IsNull() && !peering.EquinixPeerIp.IsUnknown() {
		out["equinixPeerIp"] = peering.EquinixPeerIp.ValueString()
	}
	if !peering.EquinixVrrpIp.IsNull() && !peering.EquinixVrrpIp.IsUnknown() {
		out["equinixVrrpIp"] = peering.EquinixVrrpIp.ValueString()
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
