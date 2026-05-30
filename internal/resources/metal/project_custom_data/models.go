package project_custom_data

import (
	"encoding/json"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID         types.String `tfsdk:"id"`
	ProjectID  types.String `tfsdk:"project_id"`
	CustomData types.String `tfsdk:"custom_data"`
}

func (m *ResourceModel) parse(projectID string, project *metalv1.Project) diag.Diagnostics {
	m.ID = types.StringValue(projectID)
	m.ProjectID = types.StringValue(projectID)

	// Persist compact canonical JSON to keep state stable.
	b, err := json.Marshal(project.GetCustomdata())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Failed to marshal project custom data", err.Error())}
	}
	m.CustomData = types.StringValue(string(b))
	return nil
}
