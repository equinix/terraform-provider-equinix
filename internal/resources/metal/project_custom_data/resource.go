package project_custom_data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name:   "equinix_metal_project_custom_data",
				Schema: getResourceSchema(),
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)
	projectID := plan.ProjectID.ValueString()

	customData, err := parseCustomData(plan.CustomData.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid custom_data", err.Error())
		return
	}

	if err := updateProjectCustomData(ctx, client, projectID, customData); err != nil {
		resp.Diagnostics.AddError("Failed to update project custom data", err.Error())
		return
	}

	project, diags := fetchProject(ctx, client, projectID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.parse(projectID, project)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)
	projectID := state.ProjectID.ValueString()
	if projectID == "" {
		projectID = state.ID.ValueString()
	}

	project, apiResp, err := client.ProjectsApi.FindProjectById(ctx, projectID).Execute()
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to read project",
			fmt.Sprintf("Could not read project %q: %s", projectID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(state.parse(projectID, project)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)
	projectID := plan.ProjectID.ValueString()

	customData, err := parseCustomData(plan.CustomData.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid custom_data", err.Error())
		return
	}

	if err := updateProjectCustomData(ctx, client, projectID, customData); err != nil {
		resp.Diagnostics.AddError("Failed to update project custom data", err.Error())
		return
	}

	project, diags := fetchProject(ctx, client, projectID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.parse(projectID, project)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)
	projectID := state.ProjectID.ValueString()
	if projectID == "" {
		projectID = state.ID.ValueString()
	}

	err := clearProjectCustomData(ctx, client, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to reset project custom data", err.Error())
		return
	}
}

func parseCustomData(raw string) (map[string]interface{}, error) {
	customData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(raw), &customData); err != nil {
		return nil, fmt.Errorf("custom_data must be valid JSON object: %w", err)
	}
	return customData, nil
}

func updateProjectCustomData(ctx context.Context, client *metalv1.APIClient, projectID string, customData map[string]interface{}) error {
	updateRequest := metalv1.ProjectUpdateInput{Customdata: customData}
	_, _, err := client.ProjectsApi.UpdateProject(ctx, projectID).ProjectUpdateInput(updateRequest).Execute()
	if err != nil {
		return fmt.Errorf("could not update project %q: %w", projectID, err)
	}
	return nil
}

func clearProjectCustomData(ctx context.Context, client *metalv1.APIClient, projectID string) error {
	const maxAttempts = 4

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		tflog.Debug(ctx, fmt.Sprintf("Clearing project custom data, attempt %d/%d", attempt, maxAttempts), map[string]any{"project_id": projectID})

		if err := updateProjectCustomData(ctx, client, projectID, map[string]interface{}{}); err != nil {
			return err
		}

		project, diags := fetchProject(ctx, client, projectID)
		if diags.HasError() {
			return fmt.Errorf("failed to verify project %q custom data reset: %s", projectID, diags[0].Summary())
		}

		if len(project.GetCustomdata()) == 0 {
			tflog.Debug(ctx, "Project custom data cleared successfully", map[string]any{"project_id": projectID, "attempt": attempt})
			return nil
		}

		if attempt < maxAttempts {
			tflog.Trace(ctx, "Custom data still present, waiting before retry", map[string]any{"project_id": projectID, "wait_seconds": attempt})
			select {
			case <-ctx.Done():
				return fmt.Errorf("custom data reset interrupted (context timeout): %w", ctx.Err())
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
	}

	return fmt.Errorf("project %q custom data still present after %d reset attempts", projectID, maxAttempts)
}

func fetchProject(ctx context.Context, client *metalv1.APIClient, projectID string) (*metalv1.Project, diag.Diagnostics) {
	var diags diag.Diagnostics

	project, apiResp, err := client.ProjectsApi.FindProjectById(ctx, projectID).Execute()
	if err != nil {
		if apiResp != nil && apiResp.StatusCode == http.StatusNotFound {
			diags.AddError("Project not found", fmt.Sprintf("Project %q was not found", projectID))
			return nil, diags
		}
		diags.AddError("Error reading project", fmt.Sprintf("Could not read project %q: %s", projectID, err.Error()))
		return nil, diags
	}
	return project, diags
}
