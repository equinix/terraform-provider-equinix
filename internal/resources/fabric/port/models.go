// Package port is a Terraform resource for Equinix Fabric Port Management
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
	Device                 fwtypes.ObjectValueOf[deviceModel]                   `tfsdk:"device"`
	LagEnabled             types.Bool                                           `tfsdk:"lag_enabled"`
	PhysicalPorts          fwtypes.ListNestedObjectValueOf[physicalPortModel]   `tfsdk:"physical_ports"`
	PhysicalPortsSpeed     types.Int32                                          `tfsdk:"physical_ports_speed"`
	PhysicalPortsType      types.String                                         `tfsdk:"physical_ports_type"`
	PhysicalPortsCount     types.Int32                                          `tfsdk:"physical_ports_count"`
	DemarcationPointIbx    types.String                                         `tfsdk:"demarcation_point_ibx"`
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
	PackageType    types.String `tfsdk:"package_type"`
	SharedPortType types.Bool   `tfsdk:"shared_port_type"`
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
	Interface        fwtypes.ObjectValueOf[interfaceModel]        `tfsdk:"interface"`
	DemarcationPoint fwtypes.ObjectValueOf[demarcationPointModel] `tfsdk:"demarcation_point"`
}

type interfaceModel struct {
	Type types.String `tfsdk:"type"`
}

type deviceModel struct {
	Name       types.String                                 `tfsdk:"name"`
	Redundancy fwtypes.ObjectValueOf[deviceRedundancyModel] `tfsdk:"redundancy"`
}

type deviceRedundancyModel struct {
	Priority types.String `tfsdk:"priority"`
	Group    types.String `tfsdk:"group"`
}

type demarcationPointModel struct {
	Ibx                  types.String `tfsdk:"ibx"`
	CageUniqueSpaceID    types.String `tfsdk:"cage_unique_space_id"`
	CabinetUniqueSpaceID types.String `tfsdk:"cabinet_unique_space_id"`
	PatchPanel           types.String `tfsdk:"patch_panel"`
	ConnectorType        types.String `tfsdk:"connector_type"`
}

type orderModel struct {
	PurchaseOrder       fwtypes.ObjectValueOf[purchaseOrderModel] `tfsdk:"purchase_order"`
	OrderNumber         types.String                              `tfsdk:"order_number"`
	OrderID             types.String                              `tfsdk:"order_id"`
	UUID                types.String                              `tfsdk:"uuid"`
	CustomerReferenceID types.String                              `tfsdk:"customer_reference_id"`
	Signature           fwtypes.ObjectValueOf[signatureModel]     `tfsdk:"signature"`
}

type purchaseOrderModel struct {
	Number       types.String `tfsdk:"number"`
	Amount       types.String `tfsdk:"amount"`
	AttachmentID types.String `tfsdk:"attachment_id"`
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
	if portType := port.GetType(); portType != "" {
		m.Type = types.StringValue(string(portType))
	}
	if name := port.GetName(); name != "" {
		m.Name = types.StringValue(name)
	}
	if sourceType := port.GetConnectivitySourceType(); sourceType != "" {
		m.ConnectivitySourceType = types.StringValue(string(sourceType))
	}
	m.LagEnabled = types.BoolValue(port.GetLagEnabled())
	if speed := port.GetPhysicalPortsSpeed(); speed > 0 {
		m.PhysicalPortsSpeed = types.Int32Value(speed)
	}
	if portsType := port.GetPhysicalPortsType(); portsType != "" {
		m.PhysicalPortsType = types.StringValue(string(portsType))
	}
	if count := port.GetPhysicalPortsCount(); count > 0 {
		m.PhysicalPortsCount = types.Int32Value(count)
	}
	if demarcationPointIbx := port.GetDemarcationPointIbx(); demarcationPointIbx != "" {
		m.DemarcationPointIbx = types.StringValue(port.GetDemarcationPointIbx())
	}
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
		PackageType:    types.StringValue(string(portSettings.GetPackageType())),
		SharedPortType: types.BoolValue(portSettings.GetSharedPortType()),
	}
	m.Settings = fwtypes.NewObjectValueOf[settingsModel](ctx, &settings)

	if port.Encapsulation != nil {
		portEncapsulation := port.GetEncapsulation()
		encapsulation := encapsulationModel{
			Type:          types.StringValue(string(portEncapsulation.GetType())),
			TagProtocolID: types.StringValue(portEncapsulation.GetTagProtocolId()),
		}
		m.Encapsulation = fwtypes.NewObjectValueOf[encapsulationModel](ctx, &encapsulation)
	}

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

	//// Parse device at port level if it exists
	//if port.Device != nil {
	//	m.Device = parseDevice(ctx, port.GetDevice())
	//}

	if len(port.PhysicalPorts) > 0 {
		m.PhysicalPorts = parsePhysicalPorts(ctx, port.GetPhysicalPorts())
	}

	m.Order = parseOrder(ctx, port.GetOrder())

	notifications, diags := parseNotifications(ctx, port.GetNotifications())
	if diags.HasError() {
		mDiags.Append(diags...)
		return mDiags
	}
	if len(port.GetNotifications()) > 0 {
		m.Notifications = notifications
	}

	if len(port.GetAdditionalInfo()) > 0 {
		m.AdditionalInfo = parseAdditionalInfo(ctx, port.GetAdditionalInfo())
	}

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

		//// Parse interface for each physical port if it exists
		//if portPhysicalPort.Interface != nil {
		//	physicalPort.Interface = parseInterface(ctx, portPhysicalPort.GetInterface())
		//}

		physicalPorts[i] = physicalPort
	}

	return fwtypes.NewListNestedObjectValueOfValueSlice[physicalPortModel](ctx, physicalPorts)
}

//func parseInterface(ctx context.Context, portInterface fabricv4.PortInterface) fwtypes.ObjectValueOf[interfaceModel] {
//	interfaceType := portInterface.GetType()
//
//	interfaceObj := interfaceModel{}
//
//	// Only set Type if it has a value
//	if interfaceType != "" {
//		interfaceObj.Type = types.StringValue(interfaceType)
//	} else {
//		interfaceObj.Type = types.StringNull() // ← Handle empty string
//	}
//
//	return fwtypes.NewObjectValueOf[interfaceModel](ctx, &interfaceObj)
//}

//func parseDevice(ctx context.Context, portDevice fabricv4.PortDevice) fwtypes.ObjectValueOf[deviceModel] {
//	deviceObj := deviceModel{
//		Name: types.StringValue(portDevice.GetName()),
//	}
//
//	// Parse device redundancy if it exists
//	if portDevice.Redundancy != nil {
//		redundancy := portDevice.GetRedundancy()
//		deviceRedundancy := &deviceRedundancyModel{
//			Priority: types.StringValue(string(redundancy.GetPriority())),
//		}
//
//		// Only set Group if it has a value
//		if group := redundancy.GetGroup(); group != "" {
//			deviceRedundancy.Group = types.StringValue(group)
//		} else {
//			deviceRedundancy.Group = types.StringNull() // ← FIX: Explicitly set to null if empty
//		}
//
//		deviceObj.Redundancy = fwtypes.NewObjectValueOf[deviceRedundancyModel](ctx, deviceRedundancy)
//	}
//
//	return fwtypes.NewObjectValueOf[deviceModel](ctx, &deviceObj)
//}

func parseDemarcationPoint(ctx context.Context, demPoint fabricv4.PortDemarcationPoint) fwtypes.ObjectValueOf[demarcationPointModel] {
	demarcationPoint := demarcationPointModel{
		Ibx:                  types.StringValue(demPoint.GetIbx()),
		CageUniqueSpaceID:    types.StringValue(demPoint.GetCageUniqueSpaceId()),
		CabinetUniqueSpaceID: types.StringValue(demPoint.GetCabinetUniqueSpaceId()),
		PatchPanel:           types.StringValue(demPoint.GetPatchPanel()),
		ConnectorType:        types.StringValue(demPoint.GetConnectorType()),
	}

	return fwtypes.NewObjectValueOf[demarcationPointModel](ctx, &demarcationPoint)
}

func parseOrder(ctx context.Context, portOrder fabricv4.PortOrder) fwtypes.ObjectValueOf[orderModel] {
	order := orderModel{
		OrderNumber:         types.StringValue(portOrder.GetOrderNumber()),
		OrderID:             types.StringValue(portOrder.GetOrderId()),
		UUID:                types.StringValue(portOrder.GetUuid()),
		CustomerReferenceID: types.StringValue(portOrder.GetCustomerReferenceId()),
	}

	purchaseOrder := portOrder.GetPurchaseOrder()
	order.PurchaseOrder = fwtypes.NewObjectValueOf[purchaseOrderModel](ctx, &purchaseOrderModel{
		Number:       types.StringValue(purchaseOrder.GetNumber()),
		Amount:       types.StringValue(purchaseOrder.GetAmount()),
		AttachmentID: types.StringValue(purchaseOrder.GetAttachmentId()),
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
