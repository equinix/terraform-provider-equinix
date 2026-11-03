package asn

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// NewDataSourceByAsnID creates a new data source for fetching a single ASN by UUID
func NewDataSourceByAsnID() datasource.DataSource {
	return &DataSourceByAsnID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_asn",
			},
		),
	}
}

// DataSourceByAsnID datasource for fetching a single ASN
type DataSourceByAsnID struct {
	framework.BaseDataSource
}

// Schema returns the data source schema
func (r *DataSourceByAsnID) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceSingleAsnSchema(ctx)
}

// Read fetches an ASN by its UUID
func (r *DataSourceByAsnID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceByIDModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	asn, _, err := client.ASNsApi.GetAsnById(ctx, data.AsnID.ValueString()).Execute()
	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving asn", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, asn)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
