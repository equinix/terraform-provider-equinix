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

type DataSourceByIDModel struct {
	InternetAccessServiceID types.String `tfsdk:"internet_access_service_id"`
	ID                      types.String `tfsdk:"id"`
	BaseInternetAccessServiceModel
}

func (m *DataSourceByIDModel) parse(ctx context.Context, svc *fabricv4.InternetAccessService) diag.Diagnostics {
	m.InternetAccessServiceID = types.StringValue(svc.GetUuid())
	m.ID = types.StringValue(svc.GetUuid())
	return m.BaseInternetAccessServiceModel.parse(ctx, svc)
}

type DataSourceAllInternetAccessServicesModel struct {
	ID         types.String                                                    `tfsdk:"id"`
	Data       fwtypes.ListNestedObjectValueOf[BaseInternetAccessServiceModel] `tfsdk:"data"`
	Filter     types.List                                                      `tfsdk:"filter"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]                          `tfsdk:"pagination"`
	Sort       fwtypes.ObjectValueOf[SortModel]                                `tfsdk:"sort"`
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

type BaseInternetAccessServiceModel struct {
	Href            types.String                                   `tfsdk:"href"`
	Type            types.String                                   `tfsdk:"type"`
	UUID            types.String                                   `tfsdk:"uuid"`
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
	Connections    fwtypes.ListNestedObjectValueOf[ConnectionRefModel] `tfsdk:"connections"`
}

type CustomerRouteModel struct {
	IpBlock fwtypes.ObjectValueOf[IpBlockRefModel] `tfsdk:"ip_block"`
}

type IpBlockRefModel struct {
	Href types.String `tfsdk:"href"`
	UUID types.String `tfsdk:"uuid"`
}

type ConnectionRefModel struct {
	Href types.String `tfsdk:"href"`
	UUID types.String `tfsdk:"uuid"`
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
	CreatedBy         types.String `tfsdk:"created_by"`
	CreatedByFullName types.String `tfsdk:"created_by_full_name"`
	CreatedByEmail    types.String `tfsdk:"created_by_email"`
	CreatedDateTime   types.String `tfsdk:"created_date_time"`
	UpdatedBy         types.String `tfsdk:"updated_by"`
	UpdatedByFullName types.String `tfsdk:"updated_by_full_name"`
	UpdatedByEmail    types.String `tfsdk:"updated_by_email"`
	UpdatedDateTime   types.String `tfsdk:"updated_date_time"`
	DeletedBy         types.String `tfsdk:"deleted_by"`
	DeletedByFullName types.String `tfsdk:"deleted_by_full_name"`
	DeletedByEmail    types.String `tfsdk:"deleted_by_email"`
	DeletedDateTime   types.String `tfsdk:"deleted_date_time"`
}

func (m *DataSourceAllInternetAccessServicesModel) parse(ctx context.Context, resp *fabricv4.InternetAccessServices) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(resp.GetData()) < 1 {
		diags.AddError("no data retrieved by internet access services data source",
			"either the account does not have any internet access services or the combination of filter, limit and offset needs to be updated")
		return diags
	}

	services := resp.GetData()
	data := make([]*BaseInternetAccessServiceModel, len(services))
	for i, svc := range services {
		var model BaseInternetAccessServiceModel
		diags = model.parse(ctx, &svc)
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
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[BaseInternetAccessServiceModel](ctx, data)

	return diags
}

func (m *BaseInternetAccessServiceModel) parse(ctx context.Context, svc *fabricv4.InternetAccessService) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Href = types.StringValue(svc.GetHref())
	m.Type = types.StringValue(string(svc.GetType()))
	m.UUID = types.StringValue(svc.GetUuid())
	m.Name = types.StringValue(svc.GetName())
	m.Bandwidth = types.Int32Value(svc.GetBandwidth())
	m.BandwidthCommit = types.Int32Value(svc.GetBandwidthCommit())
	m.State = types.StringValue(string(svc.GetState()))
	m.UseCase = types.StringValue(string(svc.GetUseCase()))

	svcChange := svc.GetChange()
	m.Change = fwtypes.NewObjectValueOf[ChangeModel](ctx, &ChangeModel{
		Href: types.StringValue(svcChange.GetHref()),
	})

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

	rp := svc.GetRoutingProtocol()
	customerRoutes := rp.GetCustomerRoutes()
	customerRouteModels := make([]*CustomerRouteModel, len(customerRoutes))
	for i, cr := range customerRoutes {
		ipBlock := cr.GetIpBlock()
		customerRouteModels[i] = &CustomerRouteModel{
			IpBlock: fwtypes.NewObjectValueOf[IpBlockRefModel](ctx, &IpBlockRefModel{
				Href: types.StringValue(ipBlock.GetHref()),
				UUID: types.StringValue(ipBlock.GetUuid()),
			}),
		}
	}
	connections := rp.GetConnections()
	connectionModels := make([]*ConnectionRefModel, len(connections))
	for i, conn := range connections {
		connectionModels[i] = &ConnectionRefModel{
			Href: types.StringValue(conn.GetHref()),
			UUID: types.StringValue(conn.GetUuid()),
		}
	}
	m.RoutingProtocol = fwtypes.NewObjectValueOf[RoutingProtocolModel](ctx, &RoutingProtocolModel{
		Type:           types.StringValue(string(rp.GetType())),
		CustomerRoutes: fwtypes.NewListNestedObjectValueOfSlice[CustomerRouteModel](ctx, customerRouteModels),
		Connections:    fwtypes.NewListNestedObjectValueOfSlice[ConnectionRefModel](ctx, connectionModels),
	})

	billing := svc.GetBilling()
	billingModel := BillingModel{
		Type:    types.StringValue(string(billing.GetType())),
		Enabled: types.BoolValue(billing.GetEnabled()),
	}
	if startDate, ok := billing.GetStartDateOk(); ok && startDate != nil {
		billingModel.StartDate = types.StringValue(startDate.Format(time.RFC3339))
	} else {
		billingModel.StartDate = types.StringNull()
	}
	m.Billing = fwtypes.NewObjectValueOf[BillingModel](ctx, &billingModel)

	account := svc.GetAccount()
	m.Account = fwtypes.NewObjectValueOf[AccountModel](ctx, &AccountModel{
		AccountNumber: types.StringValue(account.GetAccountNumber()),
		Href:          types.StringValue(account.GetHref()),
	})

	project := svc.GetProject()
	m.Project = fwtypes.NewObjectValueOf[ProjectModel](ctx, &ProjectModel{
		ProjectID: types.StringValue(project.GetProjectId()),
	})

	order := svc.GetOrder()
	m.Order = fwtypes.NewObjectValueOf[OrderModel](ctx, &OrderModel{
		PurchaseOrderNumber:     types.StringValue(order.GetPurchaseOrderNumber()),
		CustomerReferenceNumber: types.StringValue(order.GetCustomerReferenceNumber()),
		BillingTier:             types.StringValue(order.GetBillingTier()),
		OrderID:                 types.StringValue(order.GetOrderId()),
		OrderNumber:             types.StringValue(order.GetOrderNumber()),
		TermLength:              types.Int32Value(order.GetTermLength()),
		ContractedBandwidth:     types.Int32Value(order.GetContractedBandwidth()),
		Href:                    types.StringValue(order.GetHref()),
	})

	changeLog := svc.GetChangeLog()
	m.ChangeLog = fwtypes.NewObjectValueOf[ChangeLogModel](ctx, &ChangeLogModel{
		CreatedBy:         types.StringValue(changeLog.GetCreatedBy()),
		CreatedByFullName: types.StringValue(changeLog.GetCreatedByFullName()),
		CreatedByEmail:    types.StringValue(changeLog.GetCreatedByEmail()),
		CreatedDateTime:   types.StringValue(changeLog.GetCreatedDateTime().Format(timeFormat)),
		UpdatedBy:         types.StringValue(changeLog.GetUpdatedBy()),
		UpdatedByFullName: types.StringValue(changeLog.GetUpdatedByFullName()),
		UpdatedByEmail:    types.StringValue(changeLog.GetUpdatedByEmail()),
		UpdatedDateTime:   types.StringValue(changeLog.GetUpdatedDateTime().Format(timeFormat)),
		DeletedBy:         types.StringValue(changeLog.GetDeletedBy()),
		DeletedByFullName: types.StringValue(changeLog.GetDeletedByFullName()),
		DeletedByEmail:    types.StringValue(changeLog.GetDeletedByEmail()),
		DeletedDateTime:   types.StringValue(changeLog.GetDeletedDateTime().Format(timeFormat)),
	})

	return diags
}
