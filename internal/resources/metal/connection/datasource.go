package connection

import (
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
				Name: "equinix_metal_connection",
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
	client := r.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	// Retrieve values from plan
	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := data.ConnectionID.ValueString()

	// Use API client to get the current state of the resource
	conn, _, err := client.InterconnectionsApi.GetInterconnection(ctx, id).
		Include([]string{"service_tokens", "organization", "organization.address", "organization.billing_address", "facility", "metro", "project"}).
		Execute()

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
