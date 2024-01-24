package metal_project_ssh_key

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/packethost/packngo"

	"fmt"
	"context"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name:   "equinix_metal_project_ssh_key",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	framework.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve values from plan
	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
    id := data.ID.ValueString()
    search := data.Search.ValueString()
    projectID := data.ProjectID.ValueString()

	var (
		key        packngo.SSHKey
		searchOpts *packngo.SearchOptions
	)

	if search != "" {
		searchOpts = &packngo.SearchOptions{Search: search}
	}

    // Use API client to list SSH keys
	keys, _, err := client.Projects.ListSSHKeys(projectID, searchOpts)
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error listing project ssh keys",
			err.Error(),
		)
		return
	}

	for i := range keys {
		// use the first match for searches
		if search != "" {
			key = keys[i]
			break
		}

		// otherwise find the matching ID
		if keys[i].ID == id {
			key = keys[i]
			break
		}
	}

	if key.ID == "" {
		// Not Found
		resp.Diagnostics.AddError(
			"Error listing project ssh keys",
			fmt.Errorf("project %q SSH Key matching %q was not found", projectID, search).Error(),
		)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(data.parse(&key)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
