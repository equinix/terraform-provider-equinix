package project

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_metal_project",
			},
		),
	}
}

type DataSource struct {
	framework.BaseDataSource
}

func (r *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSchema(ctx)
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Retrieve values from plan
	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the API client from the provider metadata
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Use API client to get the current state of the resource
	var project *metalv1.Project
	if !data.Name.IsNull() {
		name := data.Name.ValueString()
		projects, err := client.ProjectsApi.FindProjects(ctx).Name(name).ExecuteWithPagination()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading Metal Project",
				"Could not read Metal Connection with Name "+name+": "+err.Error(),
			)
			return
		}
		if len(projects.Projects) > 1 {
			resp.Diagnostics.AddError(
				"Error reading Metal Project",
				fmt.Sprintf("too many projects found with name %s (found %d, expected 1)", name, len(projects.Projects)),
			)
			return
		}
		project = &projects.Projects[0]
		// project, err = findProjectByName(projects, name)
		// if err != nil {
		// 	resp.Diagnostics.AddError(
		// 		"Error reading Metal Project",
		// 		"Could not read Metal Connection with ID "+id+": "+err.Error(),
		// 	)
		// 	return
		// }
	} else {
		id := data.ProjectID.ValueString()
		var err error
		project, _, err = client.ProjectsApi.FindProjectById(ctx, id).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading Metal Project",
				"Could not read Metal Project with ID "+id+": "+err.Error(),
			)
			return
		}
	}

	bgpConf, _, err := client.BGPApi.FindBgpConfigByProject(ctx, project.GetId()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Metal Project",
			"Could not read BGP Config for Metal Project with ID "+project.GetId()+": "+err.Error(),
		)
	}
	// Set state to fully populated data
	resp.Diagnostics.Append(data.parse(ctx, project, bgpConf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// func findProjectByName(ps *metalv1.ProjectList, name string) (*metalv1.Project, error) {
// 	results := make([]metalv1.Project, 0)
// 	for _, p := range ps.Projects {
// 		if p.GetName() == name {
// 			results = append(results, p)
// 		}
// 	}
// 	if len(results) == 1 {
// 		return &results[0], nil
// 	}
// 	if len(results) == 0 {
// 		return nil, fmt.Errorf("no project found with name %s", name)
// 	}
// 	return nil, fmt.Errorf("too many projects found with name %s (found %d, expected 1)", name, len(results))
// }
