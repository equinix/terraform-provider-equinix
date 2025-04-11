package precisiontime

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_precision_time_service",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
		resp.Diagnostics.Append(diags...)
		return
	}

	ept, _, err := client.PrecisionTimeApi.CreateTimeServices(ctx).PrecisionTimeServiceRequest(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Precision Time Service", equinix_errors.FormatFabricError(err).Error())
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createWaiter := getCreateUpdateWaiter(ctx, client, ept.GetUuid(), createTimeout)
	eptChecked, err := createWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to creating Precision Time Service %s", ept.GetUuid()), err.Error())
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, eptChecked.(*fabricv4.PrecisionTimeServiceResponse))...)
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

	ept, _, err := client.PrecisionTimeApi.GetTimeServicesById(ctx, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed retrieving Precision Time Service %s", id), err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(state.parse(ctx, ept)...)
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

	// Extract the ID of the resource from the state
	serviceID := state.ID.ValueString()
	updateRequest, diags := buildUpdateRequest(ctx, state, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if len(updateRequest) > 1 {
		resp.Diagnostics.AddError("Error updating Precision Time Service",
			"This resource only accepts one attribute change at a time; please reduce changes and try again")
		return
	}

	for _, update := range updateRequest {
		if !reflect.DeepEqual(updateRequest, fabricv4.PrecisionTimeChangeOperation{}) {
			_, _, err := client.PrecisionTimeApi.UpdateTimeServicesById(ctx, serviceID).
				PrecisionTimeChangeOperation([]fabricv4.PrecisionTimeChangeOperation{update}).
				Execute()

			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating Precision Time Service",
					equinix_errors.FormatFabricError(err).Error(),
				)
			}
		}
	}

	updateTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updateWaiter := getCreateUpdateWaiter(ctx, client, serviceID, updateTimeout)
	ept, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update Precision Time Service %s", serviceID), err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, ept.(*fabricv4.PrecisionTimeServiceResponse))...)
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
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to delete the Precision Time Service
	_, deleteResp, err := client.PrecisionTimeApi.DeleteTimeServiceById(ctx, id).Execute()
	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete Precision Time Service %s", id), err.Error())
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Create(ctx, 10*time.Minute)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	deleteWaiter := getDeleteWaiter(ctx, client, id, deleteTimeout)
	_, err = deleteWaiter.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Precision Time Service %s", id), err.Error())
		return
	}
}

func buildUpdateRequest(ctx context.Context, state resourceModel, plan resourceModel) ([]fabricv4.PrecisionTimeChangeOperation, diag.Diagnostics) {
	var mDiags diag.Diagnostics

	updateRequest := make([]fabricv4.PrecisionTimeChangeOperation, 0)
	if !state.Name.Equal(plan.Name) {
		op := fabricv4.PRECISIONTIMECHANGEOPERATIONOP_REPLACE
		if plan.Name.ValueString() != "" && state.Name.ValueString() == "" {
			op = fabricv4.PRECISIONTIMECHANGEOPERATIONOP_ADD
		} else if plan.Name.ValueString() == "" && state.Name.ValueString() != "" {
			op = fabricv4.PRECISIONTIMECHANGEOPERATIONOP_REMOVE
		}
		updateRequest = append(updateRequest, fabricv4.PrecisionTimeChangeOperation{
			Op:    op,
			Path:  fabricv4.PRECISIONTIMECHANGEOPERATIONPATH_NAME,
			Value: plan.Name.ValueString(),
		})
	}
	if !state.Package.Equal(plan.Package) {
		planPackageModel := packageModel{}
		statePackageModel := packageModel{}

		planDiags := plan.Package.As(ctx, &planPackageModel, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		stateDiags := state.Package.As(ctx, &statePackageModel, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if planDiags.HasError() || stateDiags.HasError() {
			mDiags.Append(planDiags...)
			mDiags.Append(stateDiags...)
			return []fabricv4.PrecisionTimeChangeOperation{}, mDiags
		}

		if statePackageModel.Code.ValueString() != planPackageModel.Code.ValueString() {
			updateRequest = append(updateRequest, fabricv4.PrecisionTimeChangeOperation{
				Op:    fabricv4.PRECISIONTIMECHANGEOPERATIONOP_REPLACE,
				Path:  fabricv4.PRECISIONTIMECHANGEOPERATIONPATH_PACKAGE_CODE,
				Value: planPackageModel.Code.ValueString(),
			})
		}
	}

	if !state.Ipv4.Equal(plan.Ipv4) {
		planIpv4Model := ipv4Model{}
		stateIpv4Model := ipv4Model{}

		planDiags := plan.Ipv4.As(ctx, &planIpv4Model, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		stateDiags := state.Ipv4.As(ctx, &stateIpv4Model, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if planDiags.HasError() || stateDiags.HasError() {
			mDiags.Append(planDiags...)
			mDiags.Append(stateDiags...)
			return []fabricv4.PrecisionTimeChangeOperation{}, mDiags
		}

		ipv4Value := map[string]string{}

		if !planIpv4Model.Primary.IsNull() {
			ipv4Value["primary"] = planIpv4Model.Primary.ValueString()
		}
		if !planIpv4Model.Secondary.IsNull() {
			ipv4Value["secondary"] = planIpv4Model.Secondary.ValueString()
		}
		if !planIpv4Model.NetworkMask.IsNull() {
			ipv4Value["networkMask"] = planIpv4Model.NetworkMask.ValueString()
		}
		if !planIpv4Model.DefaultGateway.IsNull() {
			ipv4Value["defaultGateway"] = planIpv4Model.DefaultGateway.ValueString()
		}

		updateRequest = append(updateRequest, fabricv4.PrecisionTimeChangeOperation{
			Op:    fabricv4.PRECISIONTIMECHANGEOPERATIONOP_REPLACE,
			Path:  fabricv4.PRECISIONTIMECHANGEOPERATIONPATH_IPV4,
			Value: ipv4Value,
		})
	}

	if !state.NtpAdvanceConfiguration.Equal(plan.NtpAdvanceConfiguration) {
		var planNtpModel []ntpAdvanceConfigurationModel
		var stateNtpModel []ntpAdvanceConfigurationModel
		planDiags := plan.NtpAdvanceConfiguration.ElementsAs(ctx, &planNtpModel, false)
		stateDiags := state.NtpAdvanceConfiguration.ElementsAs(ctx, &stateNtpModel, false)

		if planDiags.HasError() || stateDiags.HasError() {
			mDiags.Append(planDiags...)
			mDiags.Append(stateDiags...)
			return []fabricv4.PrecisionTimeChangeOperation{}, mDiags
		}

		var ntpList []map[string]interface{}
		for _, ntp := range planNtpModel {
			ntpList = append(ntpList, map[string]interface{}{
				"type":      ntp.Type.ValueString(),
				"keyNumber": ntp.KeyNumber.ValueInt32(),
				"key":       ntp.Key.ValueString(),
			})
		}
		updateRequest = append(updateRequest, fabricv4.PrecisionTimeChangeOperation{
			Op:    fabricv4.PRECISIONTIMECHANGEOPERATIONOP_REPLACE,
			Path:  fabricv4.PRECISIONTIMECHANGEOPERATIONPATH_NTP_ADVANCED_CONFIGURATION,
			Value: ntpList,
		})
	}

	if !state.PtpAdvanceConfiguration.Equal(plan.PtpAdvanceConfiguration) {
		planPtpModel := ptpAdvanceConfigurationModel{}
		statePtpModel := ptpAdvanceConfigurationModel{}
		planDiags := plan.PtpAdvanceConfiguration.As(ctx, &planPtpModel, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		stateDiags := state.PtpAdvanceConfiguration.As(ctx, &statePtpModel, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if planDiags.HasError() || stateDiags.HasError() {
			mDiags.Append(planDiags...)
			mDiags.Append(stateDiags...)
			return []fabricv4.PrecisionTimeChangeOperation{}, mDiags
		}

		ptpAdvancedConfigurationValue := map[string]string{}
		if !planPtpModel.TimeScale.IsNull() {
			ptpAdvancedConfigurationValue["timeScale"] = planPtpModel.TimeScale.ValueString()
		}
		if !planPtpModel.Domain.IsNull() {
			ptpAdvancedConfigurationValue["domain"] = string(planPtpModel.Domain.ValueInt32())
		}
		if !planPtpModel.Priority1.IsNull() {
			ptpAdvancedConfigurationValue["priority1"] = string(planPtpModel.Priority1.ValueInt32())
		}
		if !planPtpModel.Priority2.IsNull() {
			ptpAdvancedConfigurationValue["priority2"] = string(planPtpModel.Priority2.ValueInt32())
		}
		if !planPtpModel.LogAnnounceInterval.IsNull() {
			ptpAdvancedConfigurationValue["logAnnounceInterval"] = string(planPtpModel.LogAnnounceInterval.ValueInt32())
		}
		if !planPtpModel.LogSyncInterval.IsNull() {
			ptpAdvancedConfigurationValue["logSyncInterval"] = string(planPtpModel.LogSyncInterval.ValueInt32())
		}
		if !planPtpModel.LogDelayReqInterval.IsNull() {
			ptpAdvancedConfigurationValue["logDelayReqInterval"] = string(planPtpModel.LogDelayReqInterval.ValueInt32())
		}
		if !planPtpModel.TransportMode.IsNull() {
			ptpAdvancedConfigurationValue["transportMode"] = planPtpModel.TransportMode.ValueString()
		}
		if !planPtpModel.GrantTime.IsNull() {
			ptpAdvancedConfigurationValue["grantTime"] = string(planPtpModel.LogDelayReqInterval.ValueInt32())
		}

		updateRequest = append(updateRequest, fabricv4.PrecisionTimeChangeOperation{
			Op:    fabricv4.PRECISIONTIMECHANGEOPERATIONOP_REPLACE,
			Path:  fabricv4.PRECISIONTIMECHANGEOPERATIONPATH_PTP_ADVANCED_CONFIGURATION,
			Value: ptpAdvancedConfigurationValue,
		})
	}

	return updateRequest, mDiags
}
func buildCreateRequest(ctx context.Context, plan resourceModel) (fabricv4.PrecisionTimeServiceRequest, diag.Diagnostics) {
	var mDiags diag.Diagnostics

	request := fabricv4.PrecisionTimeServiceRequest{}

	request.SetName(plan.Name.ValueString())
	request.SetType(fabricv4.PrecisionTimeServiceRequestType(plan.Type.ValueString()))

	var eptPackage packageModel

	if !plan.Package.IsNull() && !plan.Package.IsUnknown() {
		diags := plan.Package.As(ctx, &eptPackage, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PrecisionTimeServiceRequest{}, mDiags
		}
		request.SetPackage(fabricv4.PrecisionTimePackageRequest{Code: fabricv4.PrecisionTimePackageRequestCode(eptPackage.Code.ValueString())})
	}

	if !plan.Connections.IsNull() && !plan.Connections.IsUnknown() {
		connectionModels := make([]connectionModel, len(plan.Connections.Elements()))
		diags := plan.Connections.ElementsAs(ctx, &connectionModels, false)
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PrecisionTimeServiceRequest{}, mDiags
		}
		connections := make([]fabricv4.VirtualConnectionUuid, len(connectionModels))
		for index, connection := range connectionModels {
			connections[index].SetUuid(connection.UUID.ValueString())

			if connType := connection.Type.ValueString(); connType != "" {
				connections[index].SetType(connType)
			}

			if href := connection.Href.ValueString(); href != "" {
				connections[index].SetHref(href)
			}
		}
		request.SetConnections(connections)
	}

	var ipv4 ipv4Model
	if !plan.Ipv4.IsNull() && !plan.Ipv4.IsUnknown() {
		diags := plan.Ipv4.As(ctx, &ipv4, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PrecisionTimeServiceRequest{}, mDiags
		}
		request.SetIpv4(fabricv4.Ipv4{
			Primary:        ipv4.Primary.ValueString(),
			Secondary:      ipv4.Secondary.ValueString(),
			NetworkMask:    ipv4.NetworkMask.ValueString(),
			DefaultGateway: ipv4.DefaultGateway.ValueStringPointer(),
		})
	}

	if !plan.PtpAdvanceConfiguration.IsNull() && !plan.PtpAdvanceConfiguration.IsUnknown() {
		var ptpAdvanceConfiguration ptpAdvanceConfigurationModel
		fmt.Println(plan.PtpAdvanceConfiguration)

		diags := plan.PtpAdvanceConfiguration.As(ctx, &ptpAdvanceConfiguration, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PrecisionTimeServiceRequest{}, mDiags
		}
		timeScaleValue, _ := fabricv4.NewPtpAdvanceConfigurationTimeScaleFromValue(ptpAdvanceConfiguration.TimeScale.ValueString())

		ptpConfig := fabricv4.NewPtpAdvanceConfiguration()
		ptpConfig.SetTimeScale(*timeScaleValue)
		ptpConfig.SetDomain(ptpAdvanceConfiguration.Domain.ValueInt32())
		ptpConfig.SetPriority1(ptpAdvanceConfiguration.Priority1.ValueInt32())
		ptpConfig.SetPriority2(ptpAdvanceConfiguration.Priority2.ValueInt32())
		ptpConfig.SetLogAnnounceInterval(fabricv4.PtpAdvanceConfigurationLogAnnounceInterval(ptpAdvanceConfiguration.LogAnnounceInterval.ValueInt32()))
		ptpConfig.SetLogSyncInterval(fabricv4.PtpAdvanceConfigurationLogSyncInterval(ptpAdvanceConfiguration.LogSyncInterval.ValueInt32()))
		ptpConfig.SetTransportMode(fabricv4.PtpAdvanceConfigurationTransportMode(ptpAdvanceConfiguration.TransportMode.ValueString()))
		ptpConfig.SetGrantTime(ptpAdvanceConfiguration.GrantTime.ValueInt32())

		request.SetPtpAdvancedConfiguration(*ptpConfig)
	}

	if !plan.NtpAdvanceConfiguration.IsNull() && !plan.NtpAdvanceConfiguration.IsUnknown() {
		var ntpConfigs []ntpAdvanceConfigurationModel
		diags := plan.NtpAdvanceConfiguration.ElementsAs(ctx, &ntpConfigs, false)
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PrecisionTimeServiceRequest{}, mDiags
		}
		var convertedConfigs []fabricv4.Md5
		for _, config := range ntpConfigs {
			md5Type, _ := fabricv4.NewMd5TypeFromValue(config.Type.ValueString())
			convertedConfigs = append(convertedConfigs, fabricv4.Md5{
				Type:      md5Type,
				KeyNumber: config.KeyNumber.ValueInt32Pointer(),
				Key:       config.Key.ValueStringPointer(),
			})
		}
		request.NtpAdvancedConfiguration = convertedConfigs
	}

	var project projectModel
	if !plan.Project.IsNull() && !plan.Project.IsUnknown() {
		diags := plan.Project.As(ctx, &project, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			mDiags.Append(diags...)
			return fabricv4.PrecisionTimeServiceRequest{}, mDiags
		}
		request.SetProject(fabricv4.Project{ProjectId: project.ProjectID.ValueString()})
	}

	return request, mDiags
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Target: []string{
			string(fabricv4.PRECISIONTIMESERVICERESPONSESTATE_PROVISIONED),
			string(fabricv4.PRECISIONTIMESERVICERESPONSESTATE_CONFIGURING),
		},
		Refresh: func() (interface{}, string, error) {
			ept, _, err := client.PrecisionTimeApi.GetTimeServicesById(ctx, id).Execute()
			if err != nil {
				return 0, "", err
			}
			return ept, string(ept.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	// deletedMarker is a terraform-provider-only value that is used by the waiter
	// to indicate that the Precision Time Service appears to be deleted successfully based on
	// status code
	deletedMarker := "tf-marker-for-deleted-precision-time-service"
	return &retry.StateChangeConf{
		Target: []string{
			deletedMarker,
			string(fabricv4.PRECISIONTIMESERVICERESPONSESTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			ept, resp, err := client.PrecisionTimeApi.GetTimeServicesById(ctx, id).Execute()
			if err != nil {
				if resp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, resp.StatusCode) {
					return ept, deletedMarker, nil
				}
				return 0, "", err
			}
			return ept, string(ept.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
}
