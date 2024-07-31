package precision_time

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ResourceModel struct {
	ID                   types.String                                     `tfsdk:"id"`
	Type                 types.String                                     `tfsdk:"type"`
	Href                 types.String                                     `tfsdk:"href"`
	Uuid                 types.String                                     `tfsdk:"uuid"`
	Name                 types.String                                     `tfsdk:"name"`
	Description          types.String                                     `tfsdk:"description"`
	State                types.String                                     `tfsdk:"state"`
	Package              fwtypes.ObjectValueOf[PackageModel]              `tfsdk:"package"`
	Connections          fwtypes.ListNestedObjectValueOf[ConnectionModel] `tfsdk:"connections"`
	Ipv4                 fwtypes.ObjectValueOf[Ipv4Model]                 `tfsdk:"ipv4"`
	Account              fwtypes.ObjectValueOf[AccountModel]              `tfsdk:"account"`
	AdvanceConfiguration fwtypes.ObjectValueOf[AdvanceConfigurationModel] `tfsdk:"advance_configuration"`
	Project              fwtypes.ObjectValueOf[ProjectModel]              `tfsdk:"project"`
}

type DataSourceModel struct {
	ID                   types.String                                     `tfsdk:"id"`
	Type                 types.String                                     `tfsdk:"type"`
	Href                 types.String                                     `tfsdk:"href"`
	Uuid                 types.String                                     `tfsdk:"uuid"`
	Name                 types.String                                     `tfsdk:"name"`
	Description          types.String                                     `tfsdk:"description"`
	State                types.String                                     `tfsdk:"state"`
	Package              fwtypes.ObjectValueOf[PackageModel]              `tfsdk:"package"`
	Connections          fwtypes.ListNestedObjectValueOf[ConnectionModel] `tfsdk:"connections"`
	Ipv4                 fwtypes.ObjectValueOf[Ipv4Model]                 `tfsdk:"ipv4"`
	Account              fwtypes.ObjectValueOf[AccountModel]              `tfsdk:"account"`
	AdvanceConfiguration fwtypes.ObjectValueOf[AdvanceConfigurationModel] `tfsdk:"advance_configuration"`
	Project              fwtypes.ObjectValueOf[ProjectModel]              `tfsdk:"project"`
}

type PackageModel struct {
	Code                 types.String `tfsdk:"code"`
	Href                 types.String `tfsdk:"href"`
	Type                 types.String `tfsdk:"type"`
	Bandwidth            types.Int64  `tfsdk:"bandwidth"`
	ClientsPerSecondMax  types.Int64  `tfsdk:"clients_per_second_max"`
	RedundancySupported  types.Bool   `tfsdk:"redundancy_supported"`
	MultiSubnetSupported types.Bool   `tfsdk:"multi_subnet_supported"`
	AccuracyUnit         types.String `tfsdk:"accuracy_unit"`
	AccuracySla          types.Int64  `tfsdk:"accuracy_sla"`
	AccuracyAvgMin       types.Int64  `tfsdk:"accuracy_avg_min"`
	AccuracyAvgMax       types.Int64  `tfsdk:"accuracy_avg_max"`
}

type ConnectionModel struct {
	Uuid types.String `tfsdk:"uuid"`
	Href types.String `tfsdk:"href"`
	Type types.String `tfsdk:"type"`
}

type Ipv4Model struct {
	Primary        types.String `tfsdk:"primary"`
	Secondary      types.String `tfsdk:"secondary"`
	NetworkMask    types.String `tfsdk:"network_mask"`
	DefaultGateway types.String `tfsdk:"default_gateway"`
}

type AccountModel struct {
	AccountNumber     types.Int64  `tfsdk:"account_number"`
	IsResellerAccount types.Bool   `tfsdk:"is_reseller_account"`
	OrgId             types.String `tfsdk:"org_id"`
	GlobalOrgId       types.String `tfsdk:"global_org_id"`
}

type AdvanceConfigurationModel struct {
	Ntp fwtypes.ListNestedObjectValueOf[MD5Model] `tfsdk:"ntp"`
	Ptp fwtypes.ObjectValueOf[PTPModel]           `tfsdk:"ptp"`
}

type MD5Model struct {
	Type     types.String `tfsdk:"type"`
	Id       types.String `tfsdk:"id"`
	Password types.String `tfsdk:"password"`
}

type PTPModel struct {
	TimeScale           types.String `tfsdk:"time_scale"`
	Domain              types.Int64  `tfsdk:"domain"`
	Priority1           types.Int64  `tfsdk:"priority_1"`
	Priority2           types.Int64  `tfsdk:"priority_2"`
	LogAnnounceInterval types.Int64  `tfsdk:"log_announce_interval"`
	LogSyncInterval     types.Int64  `tfsdk:"log_sync_interval"`
	LogDelayReqInterval types.Int64  `tfsdk:"log_delay_req_interval"`
	TransportMode       types.String `tfsdk:"transport_mode"`
	GrantTime           types.Int64  `tfsdk:"grant_time"`
}

type ProjectModel struct {
	ProjectId types.String `tfsdk:"project_id"`
}

func (m *ResourceModel) parse(ctx context.Context, ept *fabricv4.PrecisionTimeServiceCreateResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = parsePrecisionTime(ctx, ept,
		&m.ID, &m.Type, &m.Href, &m.Uuid, &m.Name, &m.Description,
		&m.State, &m.Package, &m.Ipv4, &m.Project,
		&m.Account,
		&m.AdvanceConfiguration,
		&m.Connections,
	)

	return diags
}

func (m *DataSourceModel) parse(ctx context.Context, ept *fabricv4.PrecisionTimeServiceCreateResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = parsePrecisionTime(ctx, ept,
		&m.ID, &m.Type, &m.Href, &m.Uuid, &m.Name, &m.Description,
		&m.State, &m.Package, &m.Ipv4, &m.Project,
		&m.Account,
		&m.AdvanceConfiguration,
		&m.Connections,
	)

	return diags
}

func parsePrecisionTime(
	ctx context.Context,
	ept *fabricv4.PrecisionTimeServiceCreateResponse,
	id, type_, href, uuid, name, description, state *basetypes.StringValue,
	package_ *fwtypes.ObjectValueOf[PackageModel],
	ipv4 *fwtypes.ObjectValueOf[Ipv4Model],
	project *fwtypes.ObjectValueOf[ProjectModel],
	account *fwtypes.ObjectValueOf[AccountModel],
	advanceConfiguration *fwtypes.ObjectValueOf[AdvanceConfigurationModel],
	connections *fwtypes.ListNestedObjectValueOf[ConnectionModel],
) diag.Diagnostics {
	var diags diag.Diagnostics

	*id = types.StringValue(ept.GetUuid())
	*type_ = types.StringValue(string(ept.GetType()))
	*href = types.StringValue(ept.GetHref())
	*uuid = types.StringValue(ept.GetUuid())
	*name = types.StringValue(ept.GetName())
	*description = types.StringValue(ept.GetDescription())
	*state = types.StringValue(string(ept.GetState()))

	eptPackage := ept.GetPackage()
	parsedEptPackage, diags := parsePackage(ctx, &eptPackage)
	if diags.HasError() {
		return diags
	}
	*package_ = parsedEptPackage

	parsedEptConnections, diags := parseConnections(ctx, ept.GetConnections())
	if diags.HasError() {
		return diags
	}
	*connections = parsedEptConnections

	eptIpv4 := ept.GetIpv4()
	parsedEptIpv4, diags := parseIpv4(ctx, &eptIpv4)
	if diags.HasError() {
		return diags
	}
	*ipv4 = parsedEptIpv4

	eptAccount := ept.GetAccount()
	parsedEptAccount, diags := parseAccount(ctx, &eptAccount)
	if diags.HasError() {
		return diags
	}
	*account = parsedEptAccount

	eptAdvanceConfiguration := ept.GetAdvanceConfiguration()
	parsedEptAdvanceConfiguration, diags := parseAdvanceConfiguration(ctx, &eptAdvanceConfiguration)
	if diags.HasError() {
		return diags
	}
	*advanceConfiguration = parsedEptAdvanceConfiguration

	eptProject := ept.GetProject()
	parsedEptProject, diags := parseProject(ctx, &eptProject)
	if diags.HasError() {
		return diags
	}
	*project = parsedEptProject

	return diags
}

func parsePackage(ctx context.Context, package_ *fabricv4.PrecisionTimePackageResponse) (fwtypes.ObjectValueOf[PackageModel], diag.Diagnostics) {
	packageModel := &PackageModel{}

	packageModel.Code = types.StringValue(string(package_.GetCode()))
	packageModel.Href = types.StringValue(package_.GetHref())

	return fwtypes.NewObjectValueOf[PackageModel](ctx, packageModel), nil
}

func parseConnections(ctx context.Context, connections []fabricv4.FabricConnectionUuid) (fwtypes.ListNestedObjectValueOf[ConnectionModel], diag.Diagnostics) {
	connectionModels := make([]ConnectionModel, len(connections))

	for index, connection := range connections {
		connectionModels[index] = ConnectionModel{
			Uuid: types.StringValue(connection.GetUuid()),
			Href: types.StringValue(connection.GetHref()),
			Type: types.StringValue(connection.GetType()),
		}
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, connectionModels), nil
}

func parseIpv4(ctx context.Context, ipv4 *fabricv4.Ipv4) (fwtypes.ObjectValueOf[Ipv4Model], diag.Diagnostics) {
	ipv4Model := &Ipv4Model{}

	ipv4Model.Primary = types.StringValue(ipv4.GetPrimary())
	ipv4Model.Secondary = types.StringValue(ipv4.GetSecondary())
	ipv4Model.DefaultGateway = types.StringValue(ipv4.GetDefaultGateway())
	ipv4Model.NetworkMask = types.StringValue(ipv4.GetNetworkMask())

	return fwtypes.NewObjectValueOf[Ipv4Model](ctx, ipv4Model), nil
}

func parseAccount(ctx context.Context, account *fabricv4.Account) (fwtypes.ObjectValueOf[AccountModel], diag.Diagnostics) {
	accountModel := &AccountModel{}

	if account.GetAccountNumber() != 0 {
		accountModel.AccountNumber = types.Int64Value(int64(account.GetAccountNumber()))
	}
	if account.IsResellerAccount != nil {
		accountModel.IsResellerAccount = types.BoolValue(account.GetIsResellerAccount())
	}
	if account.OrgId != nil {
		accountModel.OrgId = types.StringValue(account.GetOrgId())
	}
	if account.GlobalOrgId != nil {
		accountModel.GlobalOrgId = types.StringValue(account.GetGlobalOrgId())
	}

	return fwtypes.NewObjectValueOf[AccountModel](ctx, accountModel), nil
}

func parseAdvanceConfiguration(ctx context.Context, advConfig *fabricv4.AdvanceConfiguration) (fwtypes.ObjectValueOf[AdvanceConfigurationModel], diag.Diagnostics) {
	var diags diag.Diagnostics
	advConfigModel := &AdvanceConfigurationModel{}

	md5s, diags := parseNtp(ctx, advConfig.GetNtp())
	if diags.HasError() {
		return fwtypes.NewObjectValueOfNull[AdvanceConfigurationModel](ctx), diags
	}
	advConfigModel.Ntp = md5s

	ptp := advConfig.GetPtp()
	parsedPtp, diags := parsePtp(ctx, &ptp)
	if diags.HasError() {
		return fwtypes.NewObjectValueOfNull[AdvanceConfigurationModel](ctx), diags
	}
	advConfigModel.Ptp = parsedPtp

	return fwtypes.NewObjectValueOf[AdvanceConfigurationModel](ctx, advConfigModel), nil
}

func parseNtp(ctx context.Context, ntp []fabricv4.Md5) (fwtypes.ListNestedObjectValueOf[MD5Model], diag.Diagnostics) {
	ntpModel := make([]MD5Model, len(ntp))

	for index, md5 := range ntp {
		ntpModel[index] = MD5Model{
			Type:     types.StringValue(string(md5.GetType())),
			Id:       types.StringValue(md5.GetId()),
			Password: types.StringValue(md5.GetPassword()),
		}
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, ntpModel), nil
}

func parsePtp(ctx context.Context, ptp *fabricv4.PtpAdvanceConfiguration) (fwtypes.ObjectValueOf[PTPModel], diag.Diagnostics) {
	ptpModel := &PTPModel{}

	ptpModel.TimeScale = types.StringValue(string(ptp.GetTimeScale()))
	ptpModel.Domain = types.Int64Value(int64(ptp.GetDomain()))
	ptpModel.Priority1 = types.Int64Value(int64(ptp.GetPriority1()))
	ptpModel.Priority2 = types.Int64Value(int64(ptp.GetPriority2()))
	ptpModel.LogAnnounceInterval = types.Int64Value(int64(ptp.GetLogAnnounceInterval()))
	ptpModel.LogSyncInterval = types.Int64Value(int64(ptp.GetLogSyncInterval()))
	ptpModel.LogDelayReqInterval = types.Int64Value(int64(ptp.GetLogDelayReqInterval()))
	ptpModel.TransportMode = types.StringValue(string(ptp.GetTransportMode()))
	ptpModel.GrantTime = types.Int64Value(int64(ptp.GetGrantTime()))

	return fwtypes.NewObjectValueOf[PTPModel](ctx, ptpModel), nil

}

func parseProject(ctx context.Context, project *fabricv4.Project) (fwtypes.ObjectValueOf[ProjectModel], diag.Diagnostics) {
	projectModel := &ProjectModel{}

	projectModel.ProjectId = types.StringValue(project.GetProjectId())

	return fwtypes.NewObjectValueOf[ProjectModel](ctx, projectModel), nil
}
