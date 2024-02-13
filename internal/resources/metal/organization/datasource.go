package organization

import (
	"context"
	"fmt"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/packethost/packngo"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_metal_organizations",
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
	r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
	client := r.Meta.Metal

	// Retrieve values from plan
	var data DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	id := data.organization_id.ValueString()

	var org *packngo.Organization

	// orgs, _, err := client.Organizations.List(&packngo.GetOptions{Includes: []string{"address"}})
	org, _, err := client.Organizations.Get(id, &packngo.GetOptions{Includes: []string{"address"}})
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		resp.Diagnostics.AddError(
			"Error listing Organizations",
			err.Error(),
		)
		return
	}

	if org.ID == "" {
		// Not Found
		resp.Diagnostics.AddError(
			"Error listing organizations ", fmt.Errorf("organizations was not found").Error(),
		)
		return
	}

	// Set state to fully populated data
	data.parse(ctx, org)

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
