package precisiontime

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceByIDModel struct {
	EptServiceID types.String `tfsdk:"ept_service_id"`
	ID           types.String `tfsdk:"id"`
	basePrecisionTimeModel
}

type dataSourceAllEptServicesModel struct {
	ID         types.String                                            `tfsdk:"id"`
	Data       fwtypes.ListNestedObjectValueOf[basePrecisionTimeModel] `tfsdk:"data"`
	Filter     fwtypes.ListNestedObjectValueOf[filterModel]            `tfsdk:"filters"`
	Pagination fwtypes.ObjectValueOf[paginationModel]                  `tfsdk:"pagination"`
	Sort       fwtypes.ListNestedObjectValueOf[sortModel]              `tfsdk:"sort"`
}

type filterModel struct {
	Property types.String                      `tfsdk:"property"`
	Operator types.String                      `tfsdk:"operator"`
	Values   fwtypes.ListValueOf[types.String] `tfsdk:"values"`
	Or       types.Bool                        `tfsdk:"or"`
}

type paginationModel struct {
	Offset types.Int32 `tfsdk:"offset"`
	Limit  types.Int32 `tfsdk:"limit"`
}

type sortModel struct {
	Direction types.String `tfsdk:"direction"`
	Property  types.String `tfsdk:"property"`
}

type resourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	basePrecisionTimeModel
}
type basePrecisionTimeModel struct {
	Type                    types.String                                                  `tfsdk:"type"`
	Name                    types.String                                                  `tfsdk:"name"`
	Package                 fwtypes.ObjectValueOf[packageModel]                           `tfsdk:"package"`
	Connections             fwtypes.ListNestedObjectValueOf[connectionModel]              `tfsdk:"connections"`
	Ipv4                    fwtypes.ObjectValueOf[ipv4Model]                              `tfsdk:"ipv4"`
	NtpAdvanceConfiguration fwtypes.ListNestedObjectValueOf[ntpAdvanceConfigurationModel] `tfsdk:"ntp_advanced_configuration"`
	PtpAdvanceConfiguration fwtypes.ObjectValueOf[ptpAdvanceConfigurationModel]           `tfsdk:"ptp_advanced_configuration"`
	UUID                    types.String                                                  `tfsdk:"uuid"`
	Href                    types.String                                                  `tfsdk:"href"`
	State                   types.String                                                  `tfsdk:"state"`
	Project                 fwtypes.ObjectValueOf[projectModel]                           `tfsdk:"project"`
	Account                 fwtypes.ObjectValueOf[accountModel]                           `tfsdk:"account"`
	Order                   fwtypes.ObjectValueOf[orderModel]                             `tfsdk:"order"`
	PrecisionTimePrice      fwtypes.ObjectValueOf[precisionTimePriceModel]                `tfsdk:"precision_time_price"`
	ChangeLog               fwtypes.ObjectValueOf[changeLogModel]                         `tfsdk:"change_log"`
}

type packageModel struct {
	Code types.String `tfsdk:"code"`
	Href types.String `tfsdk:"href"`
}

type connectionModel struct {
	UUID types.String `tfsdk:"uuid"`
	Href types.String `tfsdk:"href"`
	Type types.String `tfsdk:"type"`
}

type ipv4Model struct {
	Primary        types.String `tfsdk:"primary"`
	Secondary      types.String `tfsdk:"secondary"`
	NetworkMask    types.String `tfsdk:"network_mask"`
	DefaultGateway types.String `tfsdk:"default_gateway"`
}

type ntpAdvanceConfigurationModel struct {
	Type      types.String `tfsdk:"type"`
	KeyNumber types.Int32  `tfsdk:"key_number"`
	Key       types.String `tfsdk:"key"`
}

type ptpAdvanceConfigurationModel struct {
	TimeScale           types.String `tfsdk:"time_scale"`
	Domain              types.Int32  `tfsdk:"domain"`
	Priority1           types.Int32  `tfsdk:"priority1"`
	Priority2           types.Int32  `tfsdk:"priority2"`
	LogAnnounceInterval types.Int32  `tfsdk:"log_announce_interval"`
	LogSyncInterval     types.Int32  `tfsdk:"log_sync_interval"`
	LogDelayReqInterval types.Int32  `tfsdk:"log_delay_req_interval"`
	TransportMode       types.String `tfsdk:"transport_mode"`
	GrantTime           types.Int32  `tfsdk:"grant_time"`
}

type projectModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

type accountModel struct {
	AccountNumber          types.Int64  `tfsdk:"account_number"`
	AccountName            types.String `tfsdk:"account_name"`
	OrgID                  types.Int64  `tfsdk:"org_id"`
	OrganizationName       types.String `tfsdk:"organization_name"`
	GlobalOrgID            types.String `tfsdk:"global_org_id"`
	GlobalOrganizationName types.String `tfsdk:"global_organization_name"`
	UcmID                  types.String `tfsdk:"ucm_id"`
	GlobalCustID           types.String `tfsdk:"global_cust_id"`
	ResellerAccountNumber  types.Int64  `tfsdk:"reseller_account_number"`
	ResellerAccountName    types.String `tfsdk:"reseller_account_name"`
	ResellerUcmID          types.String `tfsdk:"reseller_ucm_id"`
	ResellerOrgID          types.Int64  `tfsdk:"reseller_org_id"`
}

type orderModel struct {
	PurchaseOrderNumber     types.String `tfsdk:"purchase_order_number"`
	CustomerReferenceNumber types.String `tfsdk:"customer_reference_number"`
	OrderNumber             types.String `tfsdk:"order_number"`
}

type precisionTimePriceModel struct {
	Currency types.String                                  `tfsdk:"currency"`
	Charges  fwtypes.ListNestedObjectValueOf[chargesModel] `tfsdk:"charges"`
}

type chargesModel struct {
	Type  types.String  `tfsdk:"type"`
	Price types.Float32 `tfsdk:"price"`
}

type changeLogModel struct {
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

func (m *dataSourceByIDModel) parse(ctx context.Context, ept *fabricv4.PrecisionTimeServiceResponse) diag.Diagnostics {
	m.EptServiceID = types.StringValue(ept.GetUuid())
	m.ID = types.StringValue(ept.GetUuid())
	diags := m.basePrecisionTimeModel.parse(ctx, ept)
	return diags
}

func (m *dataSourceAllEptServicesModel) parse(ctx context.Context, eptResponse *fabricv4.ServiceSearchResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(eptResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by precision time services data source",
			"either the account does not have any precision time services data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]basePrecisionTimeModel, len(eptResponse.GetData()))
	streams := eptResponse.GetData()
	for index, stream := range streams {
		var streamModel basePrecisionTimeModel
		diags = streamModel.parse(ctx, &stream)
		if diags.HasError() {
			return diags
		}
		data[index] = streamModel
	}
	responsePagination := eptResponse.GetPagination()
	pagination := paginationModel{
		Offset: types.Int32Value(responsePagination.GetOffset()),
		Limit:  types.Int32Value(responsePagination.GetLimit()),
	}

	m.ID = types.StringValue(data[0].UUID.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[paginationModel](ctx, &pagination)
	m.Data = fwtypes.NewListNestedObjectValueOfValueSlice[basePrecisionTimeModel](ctx, data)

	return diags
}

func (m *resourceModel) parse(ctx context.Context, routeAggregation *fabricv4.PrecisionTimeServiceResponse) diag.Diagnostics {
	m.ID = types.StringValue(routeAggregation.GetUuid())
	diags := m.basePrecisionTimeModel.parse(ctx, routeAggregation)
	return diags
}

func (m *basePrecisionTimeModel) parse(ctx context.Context, ept *fabricv4.PrecisionTimeServiceResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Type = types.StringValue(string(ept.GetType()))
	m.Name = types.StringValue(ept.GetName())

	m.State = types.StringValue(string(ept.GetState()))

	m.Package, diags = parsePackage(ctx, ept.GetPackage())
	if diags.HasError() {
		return diags
	}

	m.Connections, diags = parseConnections(ctx, ept.GetConnections())
	if diags.HasError() {
		return diags
	}

	m.Ipv4, diags = parseIpv4(ctx, ept.GetIpv4())
	if diags.HasError() {
		return diags
	}

	m.NtpAdvanceConfiguration, diags = parseNtpAdvanceConfiguration(ctx, ept.GetNtpAdvancedConfiguration())
	if diags.HasError() {
		return diags
	}

	m.PtpAdvanceConfiguration, diags = parsePtpAdvancedConfiguration(ctx, ept.GetPtpAdvancedConfiguration())
	if diags.HasError() {
		return diags
	}

	m.UUID = types.StringValue(ept.GetUuid())
	m.Href = types.StringValue(ept.GetHref())

	eptProject := ept.GetProject()
	result := projectModel{
		ProjectID: types.StringValue(eptProject.GetProjectId()),
	}
	m.Project = fwtypes.NewObjectValueOf[projectModel](ctx, &result)

	m.Account, diags = parseAccount(ctx, ept.GetAccount())
	if diags.HasError() {
		return diags
	}

	m.Order, diags = parseOrder(ctx, ept.GetOrder())
	if diags.HasError() {
		return diags
	}

	m.PrecisionTimePrice, diags = parsePrecisionTimePrice(ctx, ept.GetPricing())
	if diags.HasError() {
		return diags
	}

	m.ChangeLog, diags = parseChangeLog(ctx, ept.GetChangeLog())
	if diags.HasError() {
		return diags
	}

	return diags
}

func parsePackage(ctx context.Context, packageEpt fabricv4.PrecisionTimePackagePostResponse) (fwtypes.ObjectValueOf[packageModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := packageModel{}

	result.Code = types.StringValue(string(packageEpt.GetCode()))

	result.Href = types.StringValue(packageEpt.GetHref())
	return fwtypes.NewObjectValueOf[packageModel](ctx, &result), diags
}

func parseChangeLog(ctx context.Context, changeLog fabricv4.Changelog) (fwtypes.ObjectValueOf[changeLogModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}

	result := changeLogModel{}
	const TIMEFORMAT = "2008-02-02T14:02:02.000Z"
	if changeLog.GetCreatedBy() != "" {
		result.CreatedBy = types.StringValue(changeLog.GetCreatedBy())
	}
	if changeLog.GetCreatedByFullName() != "" {
		result.CreatedByFullName = types.StringValue(changeLog.GetCreatedByFullName())
	}
	if changeLog.GetCreatedByEmail() != "" {
		result.CreatedByEmail = types.StringValue(changeLog.GetCreatedByEmail())
	}
	result.CreatedDateTime = types.StringValue(changeLog.GetCreatedDateTime().Format(TIMEFORMAT))
	if changeLog.GetUpdatedBy() != "" {
		result.UpdatedBy = types.StringValue(changeLog.GetUpdatedBy())
	}
	if changeLog.GetUpdatedByFullName() != "" {
		result.UpdatedByFullName = types.StringValue(changeLog.GetUpdatedByFullName())
	}
	if changeLog.GetUpdatedByEmail() != "" {
		result.UpdatedByEmail = types.StringValue(changeLog.GetUpdatedByEmail())
	}
	result.UpdatedDateTime = types.StringValue(changeLog.GetUpdatedDateTime().Format(TIMEFORMAT))
	if changeLog.GetDeletedBy() != "" {
		result.DeletedBy = types.StringValue(changeLog.GetDeletedBy())
	}
	if changeLog.GetDeletedByFullName() != "" {
		result.DeletedByFullName = types.StringValue(changeLog.GetDeletedByFullName())
	}
	if changeLog.GetDeletedByEmail() != "" {
		result.DeletedByEmail = types.StringValue(changeLog.GetDeletedByEmail())
	}
	result.DeletedDateTime = types.StringValue(changeLog.GetDeletedDateTime().Format(TIMEFORMAT))

	return fwtypes.NewObjectValueOf[changeLogModel](ctx, &result), diags
}

func parseConnections(ctx context.Context, connections []fabricv4.VirtualConnectionTimeServiceResponse) (fwtypes.ListNestedObjectValueOf[connectionModel], diag.Diagnostics) {
	connectionModels := make([]connectionModel, len(connections))

	for index, connection := range connections {
		connectionModel := connectionModel{
			UUID: types.StringValue(connection.GetUuid()),
		}
		if connection.GetHref() != "" {
			connectionModel.Href = types.StringValue(connection.GetHref())
		}
		if connection.GetType() != "" {
			connectionModel.Type = types.StringValue(connection.GetType())
		}
		connectionModels[index] = connectionModel
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, connectionModels), nil
}

func parseIpv4(ctx context.Context, ipv4 fabricv4.Ipv4) (fwtypes.ObjectValueOf[ipv4Model], diag.Diagnostics) {
	result := ipv4Model{}
	if ipv4.GetPrimary() != "" {
		result.Primary = types.StringValue(ipv4.GetPrimary())
	}
	if ipv4.GetSecondary() != "" {
		result.Secondary = types.StringValue(ipv4.GetSecondary())
	}
	if ipv4.GetDefaultGateway() != "" {
		result.DefaultGateway = types.StringValue(ipv4.GetDefaultGateway())
	}
	if ipv4.GetNetworkMask() != "" {
		result.NetworkMask = types.StringValue(ipv4.GetNetworkMask())
	}
	return fwtypes.NewObjectValueOf[ipv4Model](ctx, &result), nil
}

func parseNtpAdvanceConfiguration(ctx context.Context, ntp []fabricv4.Md5) (fwtypes.ListNestedObjectValueOf[ntpAdvanceConfigurationModel], diag.Diagnostics) {
	ntpModels := make([]ntpAdvanceConfigurationModel, 0, len(ntp))
	diags := diag.Diagnostics{}

	if len(ntp) == 0 {
		return fwtypes.NewListNestedObjectValueOfNull[ntpAdvanceConfigurationModel](ctx), diags
	}

	for _, md5 := range ntp {
		var configModel ntpAdvanceConfigurationModel

		configModel.Type = types.StringNull()
		if md5.Type != nil {
			configModel.Type = types.StringValue(string(md5.GetType()))
		}

		configModel.KeyNumber = types.Int32Null()
		if md5.KeyNumber != nil {
			configModel.KeyNumber = types.Int32Value(md5.GetKeyNumber())
		}

		configModel.Key = types.StringNull()
		if md5.Key != nil {
			configModel.Key = types.StringValue(md5.GetKey())
		}

		if !configModel.Type.IsNull() || !configModel.Key.IsNull() || !configModel.KeyNumber.IsNull() {
			ntpModels = append(ntpModels, configModel)
		}
	}
	if len(ntpModels) > 0 {
		return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, ntpModels), diags
	}
	emptySlice := make([]ntpAdvanceConfigurationModel, 0)
	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, emptySlice), diags
}

func parsePtpAdvancedConfiguration(ctx context.Context, ptp fabricv4.PtpAdvanceConfiguration) (fwtypes.ObjectValueOf[ptpAdvanceConfigurationModel], diag.Diagnostics) {
	result := ptpAdvanceConfigurationModel{}
	hasValue := false
	if ptp.GetTimeScale() != "" {
		result.TimeScale = types.StringValue(string(ptp.GetTimeScale()))
	}
	if ptp.GetDomain() > 0 {
		result.Domain = types.Int32Value(ptp.GetDomain())
	}
	if ptp.GetPriority1() > 0 {
		result.Priority1 = types.Int32Value(ptp.GetPriority1())
	}
	if ptp.GetPriority2() > 0 {
		result.Priority2 = types.Int32Value(ptp.GetPriority2())
	}
	if ptp.GetLogAnnounceInterval() > 0 {
		result.LogAnnounceInterval = types.Int32Value(int32(ptp.GetLogAnnounceInterval()))
	}
	if ptp.GetLogSyncInterval() > 0 {
		result.LogSyncInterval = types.Int32Value(int32(ptp.GetLogSyncInterval()))
	}
	if ptp.GetLogDelayReqInterval() > 0 {
		result.LogDelayReqInterval = types.Int32Value(int32(ptp.GetLogDelayReqInterval()))
	}
	if ptp.GetTransportMode() != "" {
		result.TransportMode = types.StringValue(string(ptp.GetTransportMode()))
	}
	if ptp.GetGrantTime() > 0 {
		result.GrantTime = types.Int32Value(ptp.GetGrantTime())
	}

	if !hasValue {
		return fwtypes.NewObjectValueOfNull[ptpAdvanceConfigurationModel](ctx), nil
	}
	return fwtypes.NewObjectValueOf[ptpAdvanceConfigurationModel](ctx, &result), nil

}

func parseAccount(ctx context.Context, account fabricv4.SimplifiedAccount) (fwtypes.ObjectValueOf[accountModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := accountModel{}

	if account.GetAccountNumber() != 0 {
		result.AccountNumber = types.Int64Value(account.GetAccountNumber())
	}
	if account.GetAccountName() != "" {
		result.AccountName = types.StringValue(account.GetAccountName())
	}
	if account.OrgId != nil {
		result.OrgID = types.Int64Value(account.GetOrgId())
	}
	if account.OrganizationName != nil {
		result.OrganizationName = types.StringValue(account.GetAccountName())
	}
	if account.GlobalOrgId != nil {
		result.GlobalOrgID = types.StringValue(account.GetGlobalOrgId())
	}
	if account.GlobalOrganizationName != nil {
		result.GlobalOrganizationName = types.StringValue(account.GetGlobalOrganizationName())
	}
	if account.UcmId != nil {
		result.UcmID = types.StringValue(account.GetUcmId())
	}
	if account.GlobalCustId != nil {
		result.GlobalCustID = types.StringValue(account.GetGlobalCustId())
	}
	if account.ResellerAccountNumber != nil {
		result.ResellerAccountNumber = types.Int64Value(account.GetResellerAccountNumber())
	}
	if account.ResellerAccountName != nil {
		result.ResellerAccountName = types.StringValue(account.GetResellerAccountName())
	}
	if account.ResellerUcmId != nil {
		result.ResellerUcmID = types.StringValue(account.GetResellerUcmId())
	}
	if account.ResellerOrgId != nil {
		result.ResellerOrgID = types.Int64Value(account.GetResellerOrgId())
	}

	return fwtypes.NewObjectValueOf[accountModel](ctx, &result), diags
}

func parseOrder(ctx context.Context, order fabricv4.PrecisionTimeOrder) (fwtypes.ObjectValueOf[orderModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}
	result := orderModel{}

	if order.PurchaseOrderNumber != nil {
		result.PurchaseOrderNumber = types.StringValue(order.GetPurchaseOrderNumber())
	}
	if order.CustomerReferenceNumber != nil {
		result.PurchaseOrderNumber = types.StringValue(order.GetCustomerReferenceNumber())
	}
	if order.OrderNumber != nil {
		result.OrderNumber = types.StringValue(order.GetOrderNumber())
	}
	return fwtypes.NewObjectValueOf[orderModel](ctx, &result), diags
}

func parsePrecisionTimePrice(ctx context.Context, price fabricv4.PrecisionTimePrice) (fwtypes.ObjectValueOf[precisionTimePriceModel], diag.Diagnostics) {
	result := precisionTimePriceModel{}

	if price.Currency != nil {
		result.Currency = types.StringValue(price.GetCurrency())
	}
	charges := price.GetCharges()
	parsedCharges, diags := parseCharges(ctx, charges)
	if diags.HasError() {
		return fwtypes.NewObjectValueOf[precisionTimePriceModel](ctx, &result), diags
	}
	result.Charges = parsedCharges

	return fwtypes.NewObjectValueOf[precisionTimePriceModel](ctx, &result), diags
}

func parseCharges(ctx context.Context, charges []fabricv4.PriceCharge) (fwtypes.ListNestedObjectValueOf[chargesModel], diag.Diagnostics) {
	chargesModels := make([]chargesModel, len(charges))
	for index, charge := range charges {
		chargesModel := chargesModel{
			Type: types.StringValue(string(charge.GetType())),
		}
		if charge.GetPrice() > 0 {
			chargesModel.Price = types.Float32Value(float32(charge.GetPrice()))
		}
		chargesModels[index] = chargesModel
	}
	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, chargesModels), nil
}
