package metal_ssh_key

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/equinix/helper"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name: "equinix_metal_ssh_key",
				// do we have other than str id types?
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "The name of the SSH key for identification",
			Required:    true,
		},
		"public_key": schema.StringAttribute{
			Description: "The public key",
			Required:    true,
		},
		"fingerprint": schema.StringAttribute{
			Description: "The fingerprint of the SSH key",
			Computed:    true,
		},
		"owner_id": schema.StringAttribute{
			Description: "The UUID of the Equinix Metal API User who owns this key",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "The timestamp for when the SSH key was created",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "The timestamp for the last time the SSH key was updated",
			Computed:    true,
		},
	},
}
