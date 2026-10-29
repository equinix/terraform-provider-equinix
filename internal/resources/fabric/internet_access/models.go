package internetaccess

import (
	"context"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const timeFormat = "2006-01-02T15:04:05.000Z"

// ---- datasource by UUID ----

type DataSourceByIDModel struct {
	InternetAccessServiceID types.String `tfsdk:"internet_access_service_id"`
	ID                      types.String `tfsdk:"id"`
	BaseInternetAccessModel
}

// ---- datasource search ----

type DataSourceAllInternetAccessModel struct {
	ID         types.String                                             `tfsdk:"id"`
	Data       fwtypes.ListNestedObjectValueOf[BaseInternetAccessModel] `tfsdk:"data"`
	Filter     types.Object                                             `tfsdk:"filter"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]                   `tfsdk:"pagination"`
	Sort       fwtypes.ObjectValueOf[SortModel]                         `tfsdk:"sort"`
}

type FilterModel struct {
	Property types.String   `tfsdk:"property"`
	Operator types.String   `tfsdk:"operator"`
	Values   []types.String `tfsdk:"values"`
}

type PaginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type SortModel struct {
	Direction types.String `tfsdk:"direction"`
	Property  types.String `tfsdk:"property"`
}

// ---- base model ----

type BaseInternetAccessModel struct {
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
	RoutingProtocol fwtypes.ObjectValueOf[RoutingProtocolModel]    `tfsdk:"routing_protocol"`
	Billing         fwtypes.ObjectValueOf[BillingModel]            `tfsdk:"billing"`
	Account         fwtypes.ObjectValueOf[AccountModel]            `tfsdk:"account"`
	Project         fwtypes.ObjectValueOf[ProjectModel]            `tfsdk:"project"`
	Order           fwtypes.ObjectValueOf[OrderModel]              `tfsdk:"order"`
	ChangeLog       fwtypes.ObjectValueOf[ChangeLogModel]          `tfsdk:"change_log"`
}

type ChangeModel struct {
	Href types.String `tfsdk:"href"`
}

type LocationModel struct {
	MetroHref types.String `tfsdk:"metro_href"`
	MetroCode types.String `tfsdk:"metro_code"`
	Region    types.String `tfsdk:"region"`
	Ibx       types.String `tfsdk:"ibx"`
}

type RoutingProtocolModel struct {
	Type           types.String                                        `tfsdk:"type"`
	CustomerRoutes fwtypes.ListNestedObjectValueOf[CustomerRouteModel] `tfsdk:"customer_routes"`
	Connections    fwtypes.ListNestedObjectValueOf[ConnectionModel]    `tfsdk:"connections"`
}

type CustomerRouteModel struct {
	IpBlock fwtypes.ObjectValueOf[IpBlockModel] `tfsdk:"ip_block"`
}

type IpBlockModel struct {
	Href types.String `tfsdk:"href"`
	UUID types.String `tfsdk:"uuid"`
}

type ConnectionModel struct {
	Href        types.String                        `tfsdk:"href"`
	UUID        types.String                        `tfsdk:"uuid"`
	PeeringIpv4 fwtypes.ObjectValueOf[PeeringModel] `tfsdk:"peering_ipv4"`
	PeeringIpv6 fwtypes.ObjectValueOf[PeeringModel] `tfsdk:"peering_ipv6"`
}

type PeeringModel struct {
	Prefix         types.String `tfsdk:"prefix"`
	PrefixLength   types.Int32  `tfsdk:"prefix_length"`
	EquinixPeerIp  types.String `tfsdk:"equinix_peer_ip"`
	CustomerPeerIp types.String `tfsdk:"customer_peer_ip"`
	EquinixVrrpIp  types.String `tfsdk:"equinix_vrrp_ip"`
	CustomerVrrpIp types.String `tfsdk:"customer_vrrp_ip"`
}

type BillingModel struct {
	Type      types.String `tfsdk:"type"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	StartDate types.String `tfsdk:"start_date"`
}

type AccountModel struct {
	AccountNumber types.String `tfsdk:"account_number"`
	Href          types.String `tfsdk:"href"`
}

type ProjectModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

type OrderModel struct {
	PurchaseOrderNumber     types.String `tfsdk:"purchase_order_number"`
	CustomerReferenceNumber types.String `tfsdk:"customer_reference_number"`
	BillingTier             types.String `tfsdk:"billing_tier"`
	OrderID                 types.String `tfsdk:"order_id"`
	OrderNumber             types.String `tfsdk:"order_number"`
	TermLength              types.Int32  `tfsdk:"term_length"`
	ContractedBandwidth     types.Int32  `tfsdk:"contracted_bandwidth"`
	Href                    types.String `tfsdk:"href"`
}

type ChangeLogModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	UpdatedDateTime types.String `tfsdk:"updated_date_time"`
	DeletedDateTime types.String `tfsdk:"deleted_date_time"`
}

// ---- parse functions ----

func (m *DataSourceByIDModel) parse(ctx context.Context, svc *fabricv4.InternetAccessService) diag.Diagnostics {
	m.InternetAccessServiceID = types.StringValue(svc.GetUuid())
	m.ID = types.StringValue(svc.GetUuid())
	return m.BaseInternetAccessModel.parse(ctx, svc)
}

func (m *DataSourceAllInternetAccessModel) parse(ctx context.Context, resp *fabricv4.InternetAccessServices) diag.Diagnostics {
	var diags diag.Diagnostics

	services := resp.GetData()
	if len(services) < 1 {
		diags.AddError("no data retrieved by internet access services data source",
			"either the account does not have any internet access services or the combination of filter, limit and offset needs to be updated")
		return diags
	}

	data := make([]*BaseInternetAccessModel, len(services))
	for i := range services {
		var model BaseInternetAccessModel
		diags = model.parse(ctx, &services[i])
		if diags.HasError() {
			return diags
		}
		data[i] = &model
	}

	responsePagination := resp.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}

	m.ID = types.StringValue(data[0].UUID.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[BaseInternetAccessModel](ctx, data)

	return diags
}

func (m *BaseInternetAccessModel) parse(ctx context.Context, svc *fabricv4.InternetAccessService) diag.Diagnostics {
	var diags diag.Diagnostics

	m.UUID = types.StringValue(svc.GetUuid())
	m.Href = types.StringValue(svc.GetHref())
	m.Type = types.StringValue(string(svc.GetType()))
	m.Name = types.StringValue(svc.GetName())
	m.State = types.StringValue(string(svc.GetState()))
	m.UseCase = types.StringValue(string(svc.GetUseCase()))

	m.Bandwidth = types.Int32Null()
	if v, ok := svc.GetBandwidthOk(); ok && v != nil {
		m.Bandwidth = types.Int32Value(*v)
	}
	m.BandwidthCommit = types.Int32Null()
	if v, ok := svc.GetBandwidthCommitOk(); ok && v != nil {
		m.BandwidthCommit = types.Int32Value(*v)
	}

	change := svc.GetChange()
	m.Change = fwtypes.NewObjectValueOf[ChangeModel](ctx, &ChangeModel{Href: types.StringValue(change.GetHref())})

	locations := svc.GetLocations()
	locationModels := make([]*LocationModel, len(locations))
	for i, loc := range locations {
		locationModels[i] = &LocationModel{
			MetroHref: types.StringValue(loc.GetMetroHref()),
			MetroCode: types.StringValue(loc.GetMetroCode()),
			Region:    types.StringValue(string(loc.GetRegion())),
			Ibx:       types.StringValue(loc.GetIbx()),
		}
	}
	m.Locations = fwtypes.NewListNestedObjectValueOfSlice[LocationModel](ctx, locationModels)

	m.RoutingProtocol = fwtypes.NewObjectValueOf[RoutingProtocolModel](ctx, routingProtocolToModel(ctx, svc.GetRoutingProtocol()))

	billing := svc.GetBilling()
	billingModel := BillingModel{
		Type:      types.StringValue(string(billing.GetType())),
		Enabled:   types.BoolNull(),
		StartDate: types.StringNull(),
	}
	if v, ok := billing.GetEnabledOk(); ok && v != nil {
		billingModel.Enabled = types.BoolValue(*v)
	}
	if v, ok := billing.GetStartDateOk(); ok && v != nil {
		billingModel.StartDate = types.StringValue(v.Format(timeFormat))
	}
	m.Billing = fwtypes.NewObjectValueOf[BillingModel](ctx, &billingModel)

	account := svc.GetAccount()
	m.Account = fwtypes.NewObjectValueOf[AccountModel](ctx, &AccountModel{
		AccountNumber: types.StringValue(account.GetAccountNumber()),
		Href:          types.StringValue(account.GetHref()),
	})

	project := svc.GetProject()
	m.Project = fwtypes.NewObjectValueOf[ProjectModel](ctx, &ProjectModel{ProjectID: types.StringValue(project.GetProjectId())})

	order := svc.GetOrder()
	orderModel := OrderModel{
		PurchaseOrderNumber:     types.StringValue(order.GetPurchaseOrderNumber()),
		CustomerReferenceNumber: types.StringValue(order.GetCustomerReferenceNumber()),
		BillingTier:             types.StringValue(order.GetBillingTier()),
		OrderID:                 types.StringValue(order.GetOrderId()),
		OrderNumber:             types.StringValue(order.GetOrderNumber()),
		TermLength:              types.Int32Null(),
		ContractedBandwidth:     types.Int32Null(),
		Href:                    types.StringValue(order.GetHref()),
	}
	if v, ok := order.GetTermLengthOk(); ok && v != nil {
		orderModel.TermLength = types.Int32Value(*v)
	}
	if v, ok := order.GetContractedBandwidthOk(); ok && v != nil {
		orderModel.ContractedBandwidth = types.Int32Value(*v)
	}
	m.Order = fwtypes.NewObjectValueOf[OrderModel](ctx, &orderModel)

	cl := svc.GetChangeLog()
	m.ChangeLog = fwtypes.NewObjectValueOf[ChangeLogModel](ctx, &ChangeLogModel{
		CreatedDateTime: optionalTime(cl.CreatedDateTime),
		UpdatedDateTime: optionalTime(cl.UpdatedDateTime),
		DeletedDateTime: optionalTime(cl.DeletedDateTime),
	})

	return diags
}

func routingProtocolToModel(ctx context.Context, rp fabricv4.InternetAccessRoutingProtocol) *RoutingProtocolModel {
	routes := rp.GetCustomerRoutes()
	routeModels := make([]*CustomerRouteModel, len(routes))
	for i, route := range routes {
		ipBlock := fwtypes.NewObjectValueOfNull[IpBlockModel](ctx)
		if b, ok := route.GetIpBlockOk(); ok && b != nil {
			ipBlock = fwtypes.NewObjectValueOf[IpBlockModel](ctx, &IpBlockModel{
				Href: types.StringValue(b.GetHref()),
				UUID: types.StringValue(b.GetUuid()),
			})
		}
		routeModels[i] = &CustomerRouteModel{IpBlock: ipBlock}
	}

	connections := rp.GetConnections()
	connectionModels := make([]*ConnectionModel, len(connections))
	for i, conn := range connections {
		connModel := &ConnectionModel{
			Href:        types.StringValue(conn.GetHref()),
			UUID:        types.StringValue(conn.GetUuid()),
			PeeringIpv4: fwtypes.NewObjectValueOfNull[PeeringModel](ctx),
			PeeringIpv6: fwtypes.NewObjectValueOfNull[PeeringModel](ctx),
		}
		if p, ok := conn.GetPeeringIpv4Ok(); ok && p != nil {
			connModel.PeeringIpv4 = fwtypes.NewObjectValueOf[PeeringModel](ctx, &PeeringModel{
				Prefix:         types.StringValue(p.GetPrefix()),
				PrefixLength:   optionalInt32(p.PrefixLength),
				EquinixPeerIp:  types.StringValue(p.GetEquinixPeerIp()),
				CustomerPeerIp: types.StringValue(p.GetCustomerPeerIp()),
				EquinixVrrpIp:  types.StringValue(p.GetEquinixVrrpIp()),
				CustomerVrrpIp: types.StringValue(p.GetCustomerVrrpIp()),
			})
		}
		if p, ok := conn.GetPeeringIpv6Ok(); ok && p != nil {
			connModel.PeeringIpv6 = fwtypes.NewObjectValueOf[PeeringModel](ctx, &PeeringModel{
				Prefix:         types.StringValue(p.GetPrefix()),
				PrefixLength:   optionalInt32(p.PrefixLength),
				EquinixPeerIp:  types.StringValue(p.GetEquinixPeerIp()),
				CustomerPeerIp: types.StringValue(p.GetCustomerPeerIp()),
				EquinixVrrpIp:  types.StringValue(p.GetEquinixVrrpIp()),
				CustomerVrrpIp: types.StringValue(p.GetCustomerVrrpIp()),
			})
		}
		connectionModels[i] = connModel
	}

	return &RoutingProtocolModel{
		Type:           types.StringValue(string(rp.GetType())),
		CustomerRoutes: fwtypes.NewListNestedObjectValueOfSlice[CustomerRouteModel](ctx, routeModels),
		Connections:    fwtypes.NewListNestedObjectValueOfSlice[ConnectionModel](ctx, connectionModels),
	}
}

func optionalTime(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(timeFormat))
}

func optionalInt32(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*v)
}
