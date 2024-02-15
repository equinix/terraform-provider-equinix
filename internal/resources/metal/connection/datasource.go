package connection

import (
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/packethost/packngo"

	"context"
	"fmt"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name:   "equinix_metal_connection",
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
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve values from plan
	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := data.ConnectionID.ValueString()

	// Use API client to get the current state of the resource
	getOpts := &packngo.GetOptions{Includes: []string{"service_tokens", "organization", "facility", "metro", "project"}}
	conn, _, err := client.Connections.Get(id, getOpts)
	if err != nil {
		// If the Metal Connection is not found, remove it from the state
		if equinix_errors.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Metal Connection",
				fmt.Sprintf("[WARN] Connection (%s) not found, removing from state", id),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading Metal Connection",
			"Could not read Metal Connection with ID "+id+": "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(data.parse(ctx, conn)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
