package precision_time

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"net/http"
	"reflect"
	"slices"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_precision_time",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = resourceSchema(ctx)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {

	var plan PrecisionTimeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	createRequest, diags := buildCreateRequest(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	ept, _, err := client.PrecisionTimeApi.CreateTimeServices(ctx).PrecisionTimeServiceRequest(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Precision Time Service",
			equinix_errors.FormatFabricError(err).Error(),
		)
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
			fmt.Sprintf("Failed to create Precision Time Service %s", ept.GetUuid()), err.Error())
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, eptChecked.(*fabricv4.PrecisionTimeServiceCreateResponse))...)
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
	var state PrecisionTimeModel
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
	var state, plan PrecisionTimeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := plan.ID.ValueString()

	// Prepare update request based on the changes
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
		packageModel := PackageModel{}
		diags := plan.Package.As(ctx, &packageModel, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		updateRequest = append(updateRequest, fabricv4.PrecisionTimeChangeOperation{
			Op:    fabricv4.PRECISIONTIMECHANGEOPERATIONOP_REPLACE,
			Path:  fabricv4.PRECISIONTIMECHANGEOPERATIONPATH_PACKAGE,
			Value: packageModel.Code.ValueString(),
		})
	}

	if len(updateRequest) > 1 {
		resp.Diagnostics.AddError("Error updating Precision Time Service",
			"This resource only accepts one attribute change at a time; please reduce changes and try again")
		return
	}
	for _, update := range updateRequest {
		if !reflect.DeepEqual(updateRequest, fabricv4.PrecisionTimeChangeOperation{}) {
			_, _, err := client.PrecisionTimeApi.UpdateTimeServicesById(ctx, id).
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

	updateWaiter := getCreateUpdateWaiter(ctx, client, id, updateTimeout)
	ept, err := updateWaiter.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update Precision Time Service %s", id), err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, ept.(*fabricv4.PrecisionTimeServiceCreateResponse))...)
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
	var state PrecisionTimeModel
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
	}
}

func buildCreateRequest(ctx context.Context, plan PrecisionTimeModel) (fabricv4.PrecisionTimeServiceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.PrecisionTimeServiceRequest{
		Type: fabricv4.PrecisionTimeServiceRequestType(plan.Type.ValueString()),
		Name: plan.Name.ValueString(),
	}

	if plan.Description.ValueString() != "" {
		request.SetDescription(plan.Description.ValueString())
	}

	packageModel := PackageModel{}
	diags = plan.Package.As(ctx, &packageModel, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	if diags.HasError() {
		return fabricv4.PrecisionTimeServiceRequest{}, diags
	}

	package_ := fabricv4.PrecisionTimePackageRequest{
		Code: fabricv4.GetTimeServicesPackageByCodePackageCodeParameter(packageModel.Code.ValueString()),
	}

	if href := packageModel.Href.ValueString(); href != "" {
		package_.SetHref(href)
	}
	request.SetPackage(package_)

	connectionsModels, diags := plan.Connections.ToSlice(ctx)
	if diags.HasError() {
		return fabricv4.PrecisionTimeServiceRequest{}, diags
	}

	connections := make([]fabricv4.FabricConnectionUuid, len(connectionsModels))
	for index, connection := range connectionsModels {
		connections[index].SetUuid(connection.Uuid.ValueString())

		if type_ := connection.Type.ValueString(); type_ != "" {
			connections[index].SetType(type_)
		}

		if href := connection.Href.ValueString(); href != "" {
			connections[index].SetHref(href)
		}
	}
	request.SetConnections(connections)

	ipv4Model := Ipv4Model{}
	diags = plan.Ipv4.As(ctx, &ipv4Model, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	if diags.HasError() {
		return fabricv4.PrecisionTimeServiceRequest{}, diags
	}

	ipv4 := fabricv4.Ipv4{}
	ipv4.SetPrimary(ipv4Model.Primary.ValueString())
	ipv4.SetSecondary(ipv4Model.Secondary.ValueString())
	ipv4.SetNetworkMask(ipv4Model.NetworkMask.ValueString())
	ipv4.SetDefaultGateway(ipv4Model.DefaultGateway.ValueString())
	request.SetIpv4(ipv4)

	advConfigList := make([]AdvanceConfigurationModel, 0)
	diags = plan.AdvanceConfiguration.ElementsAs(ctx, &advConfigList, true)
	if diags.HasError() {
		return fabricv4.PrecisionTimeServiceRequest{}, diags
	}

	if len(advConfigList) != 0 {
		advConfigModel := advConfigList[0]
		ptpList := make([]PTPModel, 0)
		diags = advConfigModel.Ptp.ElementsAs(ctx, &ptpList, true)
		if diags.HasError() {
			return fabricv4.PrecisionTimeServiceRequest{}, diags
		}

		ptp := fabricv4.PtpAdvanceConfiguration{}
		if len(ptpList) != 0 {
			ptpModel := ptpList[0]
			if timeScale := ptpModel.TimeScale.ValueString(); timeScale != "" {
				ptp.SetTimeScale(fabricv4.PtpAdvanceConfigurationTimeScale(timeScale))
			}
			if domain := ptpModel.Domain.ValueInt64(); domain != 0 {
				ptp.SetPriority1(int32(domain))
			}
			if priority1 := ptpModel.Priority1.ValueInt64(); priority1 != 0 {
				ptp.SetPriority1(int32(priority1))
			}
			if priority2 := ptpModel.Priority2.ValueInt64(); priority2 != 0 {
				ptp.SetPriority2(int32(priority2))
			}
			if logAnnounceInterval := ptpModel.LogAnnounceInterval.ValueInt64(); logAnnounceInterval != 0 {
				ptp.SetLogAnnounceInterval(int32(logAnnounceInterval))
			}
			if logSyncInterval := ptpModel.LogSyncInterval.ValueInt64(); logSyncInterval != 0 {
				ptp.SetLogSyncInterval(int32(logSyncInterval))
			}
			if logDelayReqInterval := ptpModel.LogDelayReqInterval.ValueInt64(); logDelayReqInterval != 0 {
				ptp.SetLogDelayReqInterval(int32(logDelayReqInterval))
			}
			if transportMode := ptpModel.TransportMode.ValueString(); transportMode != "" {
				ptp.SetTransportMode(fabricv4.PtpAdvanceConfigurationTransportMode(transportMode))
			}
			if grantTime := ptpModel.GrantTime.ValueInt64(); grantTime != 0 {
				ptp.SetGrantTime(int32(grantTime))
			}
		}

		ntpModels, diags := advConfigModel.Ntp.ToSlice(ctx)
		if diags.HasError() {
			return fabricv4.PrecisionTimeServiceRequest{}, diags
		}

		ntps := make([]fabricv4.Md5, len(ntpModels))
		for index, md5 := range ntpModels {
			ntps[index] = fabricv4.Md5{}
			if type_ := md5.Type.ValueString(); type_ != "" {
				ntps[index].SetType(fabricv4.Md5Type(type_))
			}
			if id := md5.Id.ValueString(); id != "" {
				ntps[index].SetId(id)
			}
			if password := md5.Password.ValueString(); password != "" {
				ntps[index].SetPassword(password)
			}
		}

		advConfig := fabricv4.AdvanceConfiguration{}
		if len(ntps) > 0 {
			advConfig.SetNtp(ntps)
		}
		if !reflect.DeepEqual(ptp, fabricv4.PtpAdvanceConfiguration{}) {
			advConfig.SetPtp(ptp)
		}
		if !reflect.DeepEqual(advConfig, fabricv4.AdvanceConfiguration{}) {
			request.SetAdvanceConfiguration(advConfig)
		}
	}

	if plan.ProjectId.ValueString() != "" {
		project := fabricv4.Project{}
		project.SetProjectId(plan.ProjectId.ValueString())
		request.SetProject(project)
	}

	return request, diags
}

func getCreateUpdateWaiter(ctx context.Context, client *fabricv4.APIClient, id string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Target: []string{
			string(fabricv4.PRECISIONTIMESERVICECREATERESPONSESTATE_PROVISIONED),
			string(fabricv4.PRECISIONTIMESERVICECREATERESPONSESTATE_PENDING_CONFIGURATION),
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
			string(fabricv4.PRECISIONTIMESERVICECREATERESPONSESTATE_DEPROVISIONED),
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
