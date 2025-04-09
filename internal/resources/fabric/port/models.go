package port

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/fabric"
	int_fw "github.com/equinix/terraform-provider-equinix/internal/framework"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	basePortModel
}

type basePortModel struct {
	Type                   types.String                                         `tfsdk:"type"`
	Name                   types.String                                         `tfsdk:"name"`
	ConnectivitySourceType types.String                                         `tfsdk:"connectivity_source_type"`
	Location               fwtypes.ObjectValueOf[locationModel]                 `tfsdk:"location"`
	Settings               fwtypes.ObjectValueOf[settingsModel]                 `tfsdk:"settings"`
	Encapsulation          fwtypes.ObjectValueOf[encapsulationModel]            `tfsdk:"encapsulation"`
	Account                fwtypes.ObjectValueOf[accountModel]                  `tfsdk:"account"`
	Project                fwtypes.ObjectValueOf[projectModel]                  `tfsdk:"project"`
	Redundancy             fwtypes.ObjectValueOf[redundancyModel]               `tfsdk:"redundancy"`
	LagEnabled             types.Bool                                           `tfsdk:"lag_enabled"`
	PhysicalPorts          fwtypes.ListNestedObjectValueOf[physicalPortModel]   `tfsdk:"physical_ports"`
	PhysicalPortsSpeed     types.Int32                                          `tfsdk:"physical_ports_speed"`
	PhysicalPortsType      types.String                                         `tfsdk:"physical_ports_type"`
	PhysicalPortsCount     types.Int32                                          `tfsdk:"physical_ports_count"`
	Order                  fwtypes.ObjectValueOf[orderModel]                    `tfsdk:"order"`
	Notifications          fwtypes.ListNestedObjectValueOf[notificationModel]   `tfsdk:"notifications"`
	AdditionalInfo         fwtypes.ListNestedObjectValueOf[additionalInfoModel] `tfsdk:"additional_info"`
	ChangeLog              fwtypes.ObjectValueOf[changeLogModel]                `tfsdk:"change_log"`
	Href                   types.String                                         `tfsdk:"href"`
	UUID                   types.String                                         `tfsdk:"uuid"`
	State                  types.String                                         `tfsdk:"state"`
}

type locationModel struct {
	MetroCode types.String `tfsdk:"metro_code"`
}

type settingsModel struct {
	SharedPortType types.Bool `tfsdk:"shared_port_type"`
}

type encapsulationModel struct {
	Type          types.String `tfsdk:"type"`
	TagProtocolID types.String `tfsdk:"tag_protocol_id"`
}

type accountModel struct {
	AccountNumber types.Int64  `tfsdk:"account_number"`
	AccountName   types.String `tfsdk:"account_name"`
	UcmID         types.String `tfsdk:"ucm_id"`
}

type projectModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

type redundancyModel struct {
	Priority types.String `tfsdk:"priority"`
}

type physicalPortModel struct {
	Type             types.String                                 `tfsdk:"type"`
	DemarcationPoint fwtypes.ObjectValueOf[demarcationPointModel] `tfsdk:"demarcation_point"`
}

type demarcationPointModel struct {
	Ibx                  types.String `tfsdk:"ibx"`
	CageUniqueSpaceId    types.String `tfsdk:"cage_unique_space_id"`
	CabinetUniqueSpaceId types.String `tfsdk:"cabinet_unique_space_id"`
	PatchPanel           types.String `tfsdk:"patch_panel"`
	ConnectorType        types.String `tfsdk:"connector_type"`
}

type orderModel struct {
	PurchaseOrder       fwtypes.ObjectValueOf[purchaseOrderModel] `tfsdk:"purchase_order"`
	OrderNumber         types.String                              `tfsdk:"order_number"`
	OrderId             types.String                              `tfsdk:"order_id"`
	UUID                types.String                              `tfsdk:"uuid"`
	CustomerReferenceId types.String                              `tfsdk:"customer_reference_id"`
	Signature           fwtypes.ObjectValueOf[signatureModel]     `tfsdk:"signature"`
}

type purchaseOrderModel struct {
	Number       types.String `tfsdk:"number"`
	Amount       types.String `tfsdk:"amount"`
	AttachmentId types.String `tfsdk:"attachment_id"`
	Type         types.String `tfsdk:"type"`
	StartDate    types.String `tfsdk:"start_date"`
	EndDate      types.String `tfsdk:"end_date"`
}

type signatureModel struct {
	Signatory types.String                         `tfsdk:"signatory"`
	Delegate  fwtypes.ObjectValueOf[delegateModel] `tfsdk:"delegate"`
}

type delegateModel struct {
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Email     types.String `tfsdk:"email"`
}

type notificationModel struct {
	Type            types.String                      `tfsdk:"type"`
	RegisteredUsers fwtypes.ListValueOf[types.String] `tfsdk:"registered_users"`
}

type additionalInfoModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
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

func (m *resourceModel) parse(ctx context.Context, port *fabricv4.Port) diag.Diagnostics {
	m.ID = types.StringValue(port.GetUuid())
	diags := m.basePortModel.parse(ctx, port)
	if diags.HasError() {
		return diags
	}

	return diags
}

func (m *basePortModel) parse(ctx context.Context, port *fabricv4.Port) diag.Diagnostics {
	var mDiags diag.Diagnostics
	m.Type = types.StringValue(string(port.GetType()))
	m.Name = types.StringValue(port.GetName())
	m.ConnectivitySourceType = types.StringValue(string(port.GetConnectivitySourceType()))
	m.LagEnabled = types.BoolValue(port.GetLagEnabled())
	m.PhysicalPortsSpeed = types.Int32Value(port.GetPhysicalPortsSpeed())
	m.PhysicalPortsType = types.StringValue(string(port.GetPhysicalPortsType()))
	m.PhysicalPortsCount = types.Int32Value(port.GetPhysicalPortsCount())
	m.Href = types.StringValue(port.GetHref())
	m.UUID = types.StringValue(port.GetUuid())
	m.State = types.StringValue(string(port.GetState()))

	portLocation := port.GetLocation()
	location := locationModel{
		MetroCode: types.StringValue(portLocation.GetMetroCode()),
	}
	m.Location = fwtypes.NewObjectValueOf[locationModel](ctx, &location)

	portSettings := port.GetSettings()
	settings := settingsModel{
		SharedPortType: types.BoolValue(portSettings.GetSharedPortType()),
	}
	m.Settings = fwtypes.NewObjectValueOf[settingsModel](ctx, &settings)

	portEncapsulation := port.GetEncapsulation()
	encapsulation := encapsulationModel{
		Type:          types.StringValue(string(portEncapsulation.GetType())),
		TagProtocolID: types.StringValue(portEncapsulation.GetTagProtocolId()),
	}
	m.Encapsulation = fwtypes.NewObjectValueOf[encapsulationModel](ctx, &encapsulation)

	portAccount := port.GetAccount()
	account := accountModel{
		AccountNumber: types.Int64Value(portAccount.GetAccountNumber()),
		AccountName:   types.StringValue(portAccount.GetAccountName()),
		UcmID:         types.StringValue(portAccount.GetUcmId()),
	}
	m.Account = fwtypes.NewObjectValueOf[accountModel](ctx, &account)

	portProject := port.GetProject()
	project := projectModel{
		ProjectID: types.StringValue(portProject.GetProjectId()),
	}
	m.Project = fwtypes.NewObjectValueOf[projectModel](ctx, &project)

	portRedundancy := port.GetRedundancy()
	redundancy := redundancyModel{
		Priority: types.StringValue(string(portRedundancy.GetPriority())),
	}
	m.Redundancy = fwtypes.NewObjectValueOf[redundancyModel](ctx, &redundancy)

	m.PhysicalPorts = parsePhysicalPorts(ctx, port.GetPhysicalPorts())

	m.Order = parseOrder(ctx, port.GetOrder())

	notifications, diags := parseNotifications(ctx, port.GetNotifications())
	if diags.HasError() {
		mDiags.Append(diags...)
		return mDiags
	}
	m.Notifications = notifications

	m.AdditionalInfo = parseAdditionalInfo(ctx, port.GetAdditionalInfo())

	portChangeLog := port.GetChangeLog()
	changeLog := changeLogModel{
		CreatedBy:         types.StringValue(portChangeLog.GetCreatedBy()),
		CreatedByFullName: types.StringValue(portChangeLog.GetCreatedByFullName()),
		CreatedByEmail:    types.StringValue(portChangeLog.GetCreatedByEmail()),
		CreatedDateTime:   types.StringValue(portChangeLog.GetCreatedDateTime().Format(fabric.TimeFormat)),
		UpdatedBy:         types.StringValue(portChangeLog.GetUpdatedBy()),
		UpdatedByFullName: types.StringValue(portChangeLog.GetUpdatedByFullName()),
		UpdatedByEmail:    types.StringValue(portChangeLog.GetUpdatedByEmail()),
		UpdatedDateTime:   types.StringValue(portChangeLog.GetUpdatedDateTime().Format(fabric.TimeFormat)),
		DeletedBy:         types.StringValue(portChangeLog.GetDeletedBy()),
		DeletedByFullName: types.StringValue(portChangeLog.GetDeletedByFullName()),
		DeletedByEmail:    types.StringValue(portChangeLog.GetDeletedByEmail()),
		DeletedDateTime:   types.StringValue(portChangeLog.GetDeletedDateTime().Format(fabric.TimeFormat)),
	}
	m.ChangeLog = fwtypes.NewObjectValueOf[changeLogModel](ctx, &changeLog)

	return diags
}

func parsePhysicalPorts(ctx context.Context, portPhysicalPorts []fabricv4.PhysicalPort) fwtypes.ListNestedObjectValueOf[physicalPortModel] {
	physicalPorts := make([]physicalPortModel, len(portPhysicalPorts))
	for i, portPhysicalPort := range portPhysicalPorts {
		physicalPort := physicalPortModel{
			Type:             types.StringValue(string(portPhysicalPort.GetType())),
			DemarcationPoint: parseDemarcationPoint(ctx, portPhysicalPort.GetDemarcationPoint()),
		}
		physicalPorts[i] = physicalPort
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice[physicalPortModel](ctx, physicalPorts)
}

func parseDemarcationPoint(ctx context.Context, demPoint fabricv4.PortDemarcationPoint) fwtypes.ObjectValueOf[demarcationPointModel] {
	demarcationPoint := demarcationPointModel{
		Ibx:                  types.StringValue(demPoint.GetIbx()),
		CageUniqueSpaceId:    types.StringValue(demPoint.GetCageUniqueSpaceId()),
		CabinetUniqueSpaceId: types.StringValue(demPoint.GetCabinetUniqueSpaceId()),
		PatchPanel:           types.StringValue(demPoint.GetPatchPanel()),
		ConnectorType:        types.StringValue(demPoint.GetConnectorType()),
	}

	return fwtypes.NewObjectValueOf[demarcationPointModel](ctx, &demarcationPoint)
}

func parseOrder(ctx context.Context, portOrder fabricv4.PortOrder) fwtypes.ObjectValueOf[orderModel] {
	order := orderModel{
		OrderNumber:         types.StringValue(portOrder.GetOrderNumber()),
		OrderId:             types.StringValue(portOrder.GetOrderId()),
		UUID:                types.StringValue(portOrder.GetUuid()),
		CustomerReferenceId: types.StringValue(portOrder.GetCustomerReferenceId()),
	}

	purchaseOrder := portOrder.GetPurchaseOrder()
	order.PurchaseOrder = fwtypes.NewObjectValueOf[purchaseOrderModel](ctx, &purchaseOrderModel{
		Number:       types.StringValue(purchaseOrder.GetNumber()),
		Amount:       types.StringValue(purchaseOrder.GetAmount()),
		AttachmentId: types.StringValue(purchaseOrder.GetAttachmentId()),
		Type:         types.StringValue(string(purchaseOrder.GetType())),
		StartDate:    types.StringValue(purchaseOrder.GetStartDate()),
		EndDate:      types.StringValue(purchaseOrder.GetEndDate()),
	})

	signature := portOrder.GetSignature()
	signatureDelegate := signature.GetDelegate()
	delegate := fwtypes.NewObjectValueOf[delegateModel](ctx, &delegateModel{
		FirstName: types.StringValue(signatureDelegate.GetFirstName()),
		LastName:  types.StringValue(signatureDelegate.GetLastName()),
		Email:     types.StringValue(signatureDelegate.GetEmail()),
	})
	order.Signature = fwtypes.NewObjectValueOf[signatureModel](ctx, &signatureModel{
		Signatory: types.StringValue(string(signature.GetSignatory())),
		Delegate:  delegate,
	})

	return fwtypes.NewObjectValueOf[orderModel](ctx, &order)
}

func parseNotifications(ctx context.Context, portNotifications []fabricv4.PortNotification) (fwtypes.ListNestedObjectValueOf[notificationModel], diag.Diagnostics) {
	notifications := make([]notificationModel, len(portNotifications))
	for i, notification := range portNotifications {
		notificationRegisteredUsers := int_fw.StringSliceToAttrValue(notification.GetRegisteredUsers())
		registeredUsers, diags := fwtypes.NewListValueOf[types.String](ctx, notificationRegisteredUsers)
		if diags.HasError() {
			return fwtypes.NewListNestedObjectValueOfNull[notificationModel](ctx), diags
		}
		notifications[i] = notificationModel{
			Type:            types.StringValue(string(notification.GetType())),
			RegisteredUsers: registeredUsers,
		}
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice[notificationModel](ctx, notifications), diag.Diagnostics{}
}

func parseAdditionalInfo(ctx context.Context, portAdditionalInfo []fabricv4.PortAdditionalInfo) fwtypes.ListNestedObjectValueOf[additionalInfoModel] {
	additionalInfo := make([]additionalInfoModel, len(portAdditionalInfo))
	for i, addInfo := range portAdditionalInfo {
		additionalInfo[i] = additionalInfoModel{
			Key:   types.StringValue(addInfo.GetKey()),
			Value: types.StringValue(addInfo.GetValue()),
		}
	}
	return fwtypes.NewListNestedObjectValueOfValueSlice[additionalInfoModel](ctx, additionalInfo)
}
