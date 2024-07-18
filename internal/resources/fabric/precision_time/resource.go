package precision_time

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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

	var plan ResourceModel
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

	ept, diags = getEpt(ctx, client, &resp.State, ept.GetUuid())
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, ept)...)
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

}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {

}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {

}

func buildCreateRequest(ctx context.Context, plan ResourceModel) (fabricv4.PrecisionTimeServiceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	request := fabricv4.PrecisionTimeServiceRequest{
		Type: fabricv4.PrecisionTimeServiceRequestType(plan.Type.ValueString()),
		Name: plan.Name.ValueString(),
	}

	packageModel := PackageModel{}

	diags = plan.Package.As(ctx, &packageModel, basetypes.ObjectAsOptions{})
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

	diags = plan.Ipv4.As(ctx, &ipv4Model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.PrecisionTimeServiceRequest{}, diags
	}

	ipv4 := fabricv4.Ipv4{}
	ipv4.SetPrimary(ipv4Model.Primary.ValueString())
	ipv4.SetSecondary(ipv4Model.Secondary.ValueString())
	ipv4.SetNetworkMask(ipv4Model.NetworkMask.ValueString())
	ipv4.SetDefaultGateway(ipv4Model.DefaultGateway.ValueString())
	request.SetIpv4(ipv4)

	advConfigModel := AdvanceConfigurationModel{}

	diags = plan.AdvanceConfiguration.As(ctx, &advConfigModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.PrecisionTimeServiceRequest{}, diags
	}

	ptpModel := PTPModel{}
	diags = advConfigModel.Ptp.As(ctx, ptpModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.PrecisionTimeServiceRequest{}, diags
	}

	ptp := fabricv4.PtpAdvanceConfiguration{}
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
		ptp.SetPriority1(int32(priority2))
	}
	if logAnnounceInterval := ptpModel.LogAnnounceInterval.ValueInt64(); logAnnounceInterval != 0 {
		ptp.SetPriority1(int32(logAnnounceInterval))
	}
	if logSyncInterval := ptpModel.LogSyncInterval.ValueInt64(); logSyncInterval != 0 {
		ptp.SetPriority1(int32(logSyncInterval))
	}
	if logDelayReqInterval := ptpModel.LogDelayReqInterval.ValueInt64(); logDelayReqInterval != 0 {
		ptp.SetPriority1(int32(logDelayReqInterval))
	}
	if transportMode := ptpModel.TransportMode.ValueString(); transportMode != "" {
		ptp.SetPriority1(fabricv4.PtpAdvanceConfigurationTransportMode(transportMode))
	}
	if grantTime := ptpModel.GrantTime.ValueInt64(); grantTime != 0 {
		ptp.SetPriority1(int32(grantTime))
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
	advConfig.SetNtp(ntps)
	advConfig.SetPtp(ptp)
	request.SetAdvanceConfiguration(advConfig)

	projectModel := ProjectModel{}

	diags = plan.Project.As(ctx, &projectModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.PrecisionTimeServiceRequest{}, diags
	}

	project := fabricv4.Project{}
	if projectId := projectModel.ProjectId.ValueString(); projectId != "" {
		project.SetProjectId(projectId)
	}
	request.SetProject(project)

	return request, diags
}

func getEpt(ctx context.Context, client *fabricv4.APIClient, state *tfsdk.State, id string) (*fabricv4.PrecisionTimeServiceCreateResponse, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Use API client to get the current state of the resource
	ept, _, err := client.PrecisionTimeApi.GetTimeServicesById(ctx, id).Execute()

	if err != nil {
		// If the Precision Time Service is not found, remove it from the state
		if equinix_errors.IsNotFound(err) {
			diags.AddWarning(
				"Precision Time Service",
				fmt.Sprintf("[WARN] Precision Time Service (%s) not found, removing from state", id),
			)
			state.RemoveResource(ctx)
			return nil, diags
		}

		diags.AddError(
			"Error reading Precision Time Service",
			equinix_errors.FormatFabricError(err).Error(),
		)
		return nil, diags
	}
	return ept, diags
}
