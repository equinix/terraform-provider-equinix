package port

import (
	"context"
	"fmt"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"net/http"
	"slices"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_port",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = resourceSchema(ctx)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {

	var plan resourceModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	createRequest, diags := buildCreateRequest(ctx, plan)
	if diags.HasError() {
		return
	}

	port, _, err := client.PortsApi.CreatePort(ctx).PortRequest(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed creating port", equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, port.GetUuid(), createTimeout)
	portChecked, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed creating port %s", port.GetUuid()), err.Error())
		return
	}

	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, portChecked.(*fabricv4.Port))...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state resourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	port, _, err := client.PortsApi.GetPortByUuid(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed retrieving port %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, port)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Retrieve values from plan
	var state, plan resourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	if plan.Name.ValueString() == state.Name.ValueString() {
		resp.Diagnostics.AddWarning("No updatable fields have changed",
			"Terraform detected a config change, but it is for a field that isn't updatable for the port resource. Please revert to prior config")
		return
	}

	updateRequest := []fabricv4.PortChangeOperation{{
		Op:    "replace",
		Path:  "/name",
		Value: plan.Name.ValueString(),
	}}

	_, _, err := client.PortsApi.UpdatePortByUuid(ctx, id).PortChangeOperation(updateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed updating port %s", id), equinix_errors.FormatFabricError(err).Error())
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updateWaiter := getCreateUpdateWaiter(ctx, client, id, updateTimeout)
	portChecked, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed updating port %s", id), err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, portChecked.(*fabricv4.Port))...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve the API client
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	// Retrieve the current state
	var state resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, deleteResp, err := client.PortsApi.DeletePort(ctx, id).Execute()
	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed deleting port %s", id), equinix_errors.FormatFabricError(err).Error())
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	deleteWaiter := getDeleteWaiter(ctx, client, id, deleteTimeout)
	_, err = deleteWaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed deleting port %s", id), err.Error())
		return
	}

}

func buildCreateRequest(ctx context.Context, plan resourceModel) (fabricv4.PortRequest, diag.Diagnostics) {
	var mDiags diag.Diagnostics
	request := fabricv4.PortRequest{}

	request.SetType(fabricv4.PortType(plan.Type.ValueString()))
	request.SetConnectivitySourceType(fabricv4.PortConnectivitySourceType(plan.ConnectivitySourceType.ValueString()))
	request.SetLagEnabled(plan.LagEnabled.ValueBool())
	request.SetPhysicalPortsSpeed(plan.PhysicalPortsSpeed.ValueInt32())
	request.SetPhysicalPortsType(fabricv4.PortPhysicalPortsType(plan.PhysicalPortsType.ValueString()))
	request.SetPhysicalPortsCount(plan.PhysicalPortsCount.ValueInt32())
	request.SetHref(plan.Href.ValueString())
	request.SetUuid(plan.UUID.ValueString())

	if !plan.Location.IsNull() && !plan.Location.IsUnknown() {
		var location locationModel
		diags := plan.Location.As(ctx, &location, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		simplifiedLocation := fabricv4.SimplifiedLocation{}
		simplifiedLocation.SetMetroCode(location.MetroCode.ValueString())
		request.SetLocation(simplifiedLocation)
	}

	if !plan.Settings.IsNull() && !plan.Settings.IsUnknown() {
		var settings settingsModel
		diags := plan.Settings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		portSettings := fabricv4.PortSettings{}
		portSettings.SetSharedPortType(settings.SharedPortType.ValueBool())
		request.SetSettings(portSettings)
	}

	if !plan.Encapsulation.IsNull() && !plan.Encapsulation.IsUnknown() {
		var encapsulation encapsulationModel
		diags := plan.Encapsulation.As(ctx, &encapsulation, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		portEncapsulation := fabricv4.PortEncapsulation{}
		portEncapsulation.SetType(fabricv4.PortEncapsulationType(encapsulation.Type.ValueString()))
		portEncapsulation.SetTagProtocolId(encapsulation.TagProtocolID.ValueString())
		request.SetEncapsulation(fabricv4.PortEncapsulation{})
	}

	if !plan.Account.IsNull() && !plan.Account.IsUnknown() {
		var account accountModel
		diags := plan.Account.As(ctx, &account, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		simplifiedAccount := fabricv4.SimplifiedAccount{}
		simplifiedAccount.SetAccountNumber(account.AccountNumber.ValueInt64())
		simplifiedAccount.SetAccountName(account.AccountName.ValueString())
		simplifiedAccount.SetUcmId(account.UcmID.ValueString())
		request.SetAccount(simplifiedAccount)
	}

	if !plan.Project.IsNull() && !plan.Project.IsUnknown() {
		var project projectModel
		diags := plan.Project.As(ctx, &project, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		request.SetProject(fabricv4.Project{ProjectId: project.ProjectID.ValueString()})
	}

	if !plan.Redundancy.IsNull() && !plan.Redundancy.IsUnknown() {
		var redundancy redundancyModel
		diags := plan.Redundancy.As(ctx, &redundancy, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		portRedundancy := fabricv4.PortRedundancy{}
		portRedundancy.SetPriority(fabricv4.PortPriority(redundancy.Priority.ValueString()))
		request.SetRedundancy(portRedundancy)
	}

	if !plan.Redundancy.IsNull() && !plan.Redundancy.IsUnknown() {
		var redundancy redundancyModel
		diags := plan.Redundancy.As(ctx, &redundancy, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		portRedundancy := fabricv4.PortRedundancy{}
		portRedundancy.SetPriority(fabricv4.PortPriority(redundancy.Priority.ValueString()))
		request.SetRedundancy(portRedundancy)
	}

	if !plan.PhysicalPorts.IsNull() && !plan.PhysicalPorts.IsUnknown() {
		physicalPorts, diags := buildPhysicalPorts(ctx, plan.PhysicalPorts)
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		request.SetPhysicalPorts(physicalPorts)
	}

	if !plan.Order.IsNull() && !plan.Order.IsUnknown() {
		order, diags := buildOrder(ctx, plan.Order)
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		request.SetOrder(order)
	}

	if !plan.Notifications.IsNull() && !plan.Notifications.IsUnknown() {
		notifications := make([]notificationModel, len(plan.Notifications.Elements()))
		diags := plan.Notifications.ElementsAs(ctx, notifications, false)
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		portNotifications := make([]fabricv4.PortNotification, len(notifications))
		for i, v := range notifications {
			portNotifications[i].SetType(fabricv4.PortNotificationType(v.Type.ValueString()))
			registeredUsers := make([]string, len(v.RegisteredUsers.Elements()))
			diags = v.RegisteredUsers.ElementsAs(ctx, registeredUsers, false)
			if diags.HasError() {
				mDiags.Append(diags...)
				return fabricv4.PortRequest{}, mDiags
			}
			portNotifications[i].SetRegisteredUsers(registeredUsers)
		}
		request.SetNotifications(portNotifications)
	}

	if !plan.AdditionalInfo.IsNull() && !plan.AdditionalInfo.IsUnknown() {
		additionalInfo := make([]additionalInfoModel, len(plan.AdditionalInfo.Elements()))
		diags := plan.AdditionalInfo.ElementsAs(ctx, additionalInfo, false)
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PortRequest{}, mDiags
		}
		portAdditionalInfo := make([]fabricv4.PortAdditionalInfo, len(additionalInfo))
		for i, v := range additionalInfo {
			portAdditionalInfo[i].SetKey(v.Key.ValueString())
			portAdditionalInfo[i].SetValue(v.Value.ValueString())
		}
		request.SetAdditionalInfo(portAdditionalInfo)
	}

	return request, mDiags
}

func buildPhysicalPorts(ctx context.Context, physicalPortsObject fwtypes.ListNestedObjectValueOf[physicalPortModel]) ([]fabricv4.PhysicalPort, diag.Diagnostics) {
	physicalPortModels := make([]physicalPortModel, len(physicalPortsObject.Elements()))
	mDiags := physicalPortsObject.ElementsAs(ctx, physicalPortModels, false)
	if mDiags.HasError() {
		return []fabricv4.PhysicalPort{}, mDiags
	}
	physicalPorts := make([]fabricv4.PhysicalPort, len(physicalPortModels))
	for i, v := range physicalPortModels {
		physicalPorts[i].SetType(fabricv4.PhysicalPortType(v.Type.ValueString()))
		var demarcationPoint demarcationPointModel
		diags := v.DemarcationPoint.As(ctx, &demarcationPoint, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return []fabricv4.PhysicalPort{}, mDiags
		}
		fabricDemarcationPoint := fabricv4.PortDemarcationPoint{}
		fabricDemarcationPoint.SetIbx(demarcationPoint.Ibx.ValueString())
		fabricDemarcationPoint.SetCageUniqueSpaceId(demarcationPoint.CageUniqueSpaceId.ValueString())
		fabricDemarcationPoint.SetCabinetUniqueSpaceId(demarcationPoint.CabinetUniqueSpaceId.ValueString())
		fabricDemarcationPoint.SetConnectorType(demarcationPoint.ConnectorType.ValueString())
		fabricDemarcationPoint.SetPatchPanel(demarcationPoint.PatchPanel.ValueString())
		physicalPorts[i].SetDemarcationPoint(fabricDemarcationPoint)
	}

	return physicalPorts, mDiags
}

func buildOrder(ctx context.Context, orderObject fwtypes.ObjectValueOf[orderModel]) (fabricv4.PortOrder, diag.Diagnostics) {
	var order orderModel
	mDiags := orderObject.As(ctx, &order, basetypes.ObjectAsOptions{})
	if mDiags.HasError() {
		return fabricv4.PortOrder{}, mDiags
	}

	portOrder := fabricv4.PortOrder{}
	portOrder.SetOrderNumber(order.OrderNumber.ValueString())
	portOrder.SetOrderId(order.OrderId.ValueString())
	portOrder.SetCustomerReferenceId(order.CustomerReferenceId.ValueString())

	var purchaseOrder purchaseOrderModel
	diags := order.PurchaseOrder.As(ctx, &purchaseOrder, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		mDiags.Append(diags...)
		return fabricv4.PortOrder{}, mDiags
	}
	portPurchaseOrder := fabricv4.PortOrderPurchaseOrder{}
	portPurchaseOrder.SetType(fabricv4.PortOrderPurchaseOrderType(purchaseOrder.Type.ValueString()))
	portPurchaseOrder.SetAmount(purchaseOrder.Amount.ValueString())
	portPurchaseOrder.SetAttachmentId(purchaseOrder.AttachmentId.ValueString())
	portPurchaseOrder.SetNumber(purchaseOrder.Number.ValueString())
	portPurchaseOrder.SetStartDate(purchaseOrder.StartDate.ValueString())
	portPurchaseOrder.SetEndDate(purchaseOrder.EndDate.ValueString())

	portOrder.SetPurchaseOrder(portPurchaseOrder)

	var signature signatureModel
	diags = order.Signature.As(ctx, &signature, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		mDiags.Append(diags...)
		return fabricv4.PortOrder{}, mDiags
	}
	portSignature := fabricv4.PortOrderSignature{}
	portSignature.SetSignatory(fabricv4.PortOrderSignatureSignatory(signature.Signatory.ValueString()))
	var delegate delegateModel
	diags = signature.Delegate.As(ctx, &delegate, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		mDiags.Append(diags...)
		return fabricv4.PortOrder{}, mDiags
	}
	signatureDelegate := fabricv4.PortOrderSignatureDelegate{}
	signatureDelegate.SetFirstName(delegate.FirstName.ValueString())
	signatureDelegate.SetLastName(delegate.LastName.ValueString())
	signatureDelegate.SetEmail(delegate.Email.ValueString())
	portSignature.SetDelegate(signatureDelegate)
	portOrder.SetSignature(portSignature)

	return portOrder, mDiags
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.PORTSTATE_PROVISIONING),
			string(fabricv4.PORTSTATE_PENDING),
		},
		Target: []string{
			string(fabricv4.PORTSTATE_PROVISIONED),
			string(fabricv4.PORTSTATE_ADDED),
			string(fabricv4.PORTSTATE_ACTIVE),
		},
		Refresh: func() (interface{}, string, error) {
			port, _, err := client.PortsApi.GetPortByUuid(ctx, id).Execute()
			if err != nil {
				return 0, "", err
			}
			return port, string(port.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the connection appears to be deleted successfully based on
	// status code
	deletedMarker := "tf-marker-for-deletion"
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.PORTSTATE_PROVISIONED),
			string(fabricv4.PORTSTATE_ADDED),
			string(fabricv4.PORTSTATE_ACTIVE),
		},
		Target: []string{
			deletedMarker,
			string(fabricv4.PORTSTATE_DELETED),
			string(fabricv4.PORTSTATE_TO_BE_DELETED),
			string(fabricv4.PORTSTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			port, resp, err := client.PortsApi.GetPortByUuid(ctx, id).Execute()
			if err != nil {
				if resp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return port, deletedMarker, nil
				}
				return 0, "", err
			}
			return port, string(port.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
