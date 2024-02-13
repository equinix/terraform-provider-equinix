package project_ssh_key

import (
	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"context"
	"fmt"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name:   "equinix_metal_project_ssh_key",
				Schema: &dataSourceSchema,
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
	r.Meta.AddFwModuleToMetalGoUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metalgo

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
		key metalv1.SSHKey
	)

	// Use API client to list SSH keys
	keysList, _, err := client.SSHKeysApi.FindProjectSSHKeys(context.Background(), projectID).Query(search).Execute()
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error listing project ssh keys",
			err.Error(),
		)
		return
	}

	keys := keysList.GetSshKeys()
	for i := range keys {
		// use the first match for searches
		if search != "" {
			key = keys[i]
			break
		}

		// otherwise find the matching ID
		if keys[i].GetId() == id {
			key = keys[i]
			break
		}
	}

	if key.GetId() == "" {
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
