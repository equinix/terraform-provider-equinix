package project

import (
	"context"
	"fmt"
	"reflect"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_metal_project",
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
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Prepare the data for API request
	createRequest := metalv1.ProjectCreateFromRootInput{
		Name: plan.Name.ValueString(),
	}

	// Include optional fields if they are set
	if !plan.OrganizationID.IsNull() && !plan.OrganizationID.IsUnknown() {
		createRequest.OrganizationId = plan.OrganizationID.ValueStringPointer()
	}

	// API call to create the project
	project, createResp, err := client.ProjectsApi.CreateProject(ctx).ProjectCreateFromRootInput(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project: "+equinix_errors.FriendlyErrorForMetalGo(err, createResp).Error(),
		)
		return
	}

	// Handle BGP Config if present
	if !plan.BGPConfig.IsNull() {
		bgpCreateRequest, err := expandBGPConfig(ctx, plan.BGPConfig)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating project",
				"Could not validate BGP Config: "+err.Error(),
			)
			return
		}

		createResp, err = client.BGPApi.RequestBgpConfig(ctx, project.GetId()).BgpConfigRequestInput(*bgpCreateRequest).Execute()
		if err != nil {
			err = equinix_errors.FriendlyErrorForMetalGo(err, createResp)
			resp.Diagnostics.AddError(
				"Error creating BGP configuration",
				"Could not create BGP configuration for project: "+err.Error(),
			)
			return
		}
	}

	// Enable Backend Transfer if True
	if plan.BackendTransfer.ValueBool() {
		pur := metalv1.ProjectUpdateInput{
			BackendTransferEnabled: plan.BackendTransfer.ValueBoolPointer(),
		}
		_, updateResp, err := client.ProjectsApi.UpdateProject(ctx, project.GetId()).ProjectUpdateInput(pur).Execute()
		if err != nil {
			err = equinix_errors.FriendlyErrorForMetalGo(err, updateResp)
			resp.Diagnostics.AddError(
				"Error enabling Backend Transfer",
				"Could not enable Backend Transfer for project with ID "+project.GetId()+": "+err.Error(),
			)
			return
		}
	}

	// Use API client to get the current state of the resource
	project, diags = fetchProject(ctx, client, project.GetId())
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Fetch BGP Config if needed
	var bgpConfig *metalv1.BgpConfig
	if !plan.BGPConfig.IsNull() {
		bgpConfig, diags = fetchBGPConfig(ctx, client, project.GetId())
		diags.Append(diags...)
		if diags.HasError() {
			return
		}
	}

	// Parse API response into the Terraform state
	resp.Diagnostics.Append(plan.parse(ctx, project, bgpConfig)...)
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
	// Retrieve the current state
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// Use API client to get the current state of the resource
	project, diags := fetchProject(ctx, client, id)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Use API client to fetch BGP Config
	var bgpConfig *metalv1.BgpConfig
	bgpConfig, diags = fetchBGPConfig(ctx, client, project.GetId())
	diags.Append(diags...)
	if diags.HasError() {
		return
	}

	// Parse the API response into the Terraform state
	resp.Diagnostics.Append(state.parse(ctx, project, bgpConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func fetchProject(ctx context.Context, client *metalv1.APIClient, projectID string) (*metalv1.Project, diag.Diagnostics) {
	var diags diag.Diagnostics

	project, apiResp, err := client.ProjectsApi.FindProjectById(ctx, projectID).Execute()
	if err != nil {
		err = equinix_errors.FriendlyErrorForMetalGo(err, apiResp)

		// Check if the Project no longer exists
		if equinix_errors.IsNotFound(err) {
			diags.AddWarning(
				"Project not found",
				fmt.Sprintf("Project (%s) not found, removing from state", projectID),
			)
		} else {
			diags.AddError(
				"Error reading project",
				"Could not read project with ID "+projectID+": "+err.Error(),
			)
		}
		return nil, diags
	}

	return project, diags
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Retrieve values from plan
	var state, plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// Handle BGP Config changes
	bgpConfig, diags := handleBGPConfigChanges(ctx, client, &plan, &state, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare Project update request based on the changes
	updateRequest := metalv1.ProjectUpdateInput{}
	if state.Name != plan.Name {
		updateRequest.Name = plan.Name.ValueStringPointer()
	}
	if state.PaymentMethodID != plan.PaymentMethodID {
		updateRequest.PaymentMethodId = plan.PaymentMethodID.ValueStringPointer()
	}
	if state.BackendTransfer != plan.BackendTransfer {
		updateRequest.BackendTransferEnabled = plan.BackendTransfer.ValueBoolPointer()
	}

	var project *metalv1.Project
	// Check if any update was requested
	if !reflect.DeepEqual(updateRequest, metalv1.ProjectUpdateInput{}) {
		// API call to update the project
		_, updateResp, err := client.ProjectsApi.UpdateProject(ctx, id).ProjectUpdateInput(updateRequest).Execute()
		if err != nil {
			friendlyErr := equinix_errors.FriendlyErrorForMetalGo(err, updateResp)
			resp.Diagnostics.AddError(
				"Error updating project",
				"Could not update project with ID "+id+": "+friendlyErr.Error(),
			)
			return
		}
	}

	// Use API client to get the current state of the resource
	project, diags = fetchProject(ctx, client, id)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(plan.parse(ctx, project, bgpConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state back into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve the current state
	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Extract the ID of the resource from the state
	id := state.ID.ValueString()

	// API call to delete the project
	deleteResp, err := client.ProjectsApi.DeleteProject(ctx, id).Execute()
	if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(deleteResp, err) != nil {
		err = equinix_errors.FriendlyErrorForMetalGo(err, deleteResp)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Project %s", id),
			err.Error(),
		)
	}
}
