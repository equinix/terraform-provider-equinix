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
				Name: "equinix_metal_organization",
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
	var data dataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the state
	orgId := data.OrganizationID
	name := data.Name

	var (
		orgOk *packngo.Organization
	)

	if !name.IsNull() {
		orgList, _, err := client.Organizations.List(&packngo.GetOptions{Includes: []string{"address"}})
		if err != nil {
			err = equinix_errors.FriendlyError(err)
			resp.Diagnostics.AddError(
				"Error listing Organizations",
				err.Error(),
			)
			return
		}
		orgOk, err = findOrgByName(orgList, name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing organizations ", fmt.Errorf("organizations was not found").Error(),
			)
			return
		}

	} else {
		orgID := orgId.ValueString()
		org, _, err := client.Organizations.Get(orgID, &packngo.GetOptions{Includes: []string{"address"}})
		if err != nil {
			err = equinix_errors.FriendlyError(err)
			resp.Diagnostics.AddError(
				"Error getting Organization",
				err.Error(),
			)
			return
		}
		orgOk = org
	}

	if orgOk.ID == "" {
		// Not Found
		resp.Diagnostics.AddError(
			"Error listing organizations ", fmt.Errorf("organizations was not found").Error(),
		)
		return
	}

	// Set state to fully populated data
	data.parse(ctx, orgOk)

	// Update the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findOrgByName(os []packngo.Organization, name string) (*packngo.Organization, error) {
	results := make([]packngo.Organization, 0)
	for _, o := range os {
		if o.Name == name {
			results = append(results, o)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no organization found with name %s", name)
	}
	return nil, fmt.Errorf("too many organizations found with name %s (found %d, expected 1)", name, len(results))
}
