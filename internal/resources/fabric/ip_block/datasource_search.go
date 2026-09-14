package ipblock

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewDataSourceAllIpBlocks creates a new data source for searching IP blocks
func NewDataSourceAllIpBlocks() datasource.DataSource {
	return &DataSourceAllIpBlocks{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_ip_blocks",
			},
		),
	}
}

// DataSourceAllIpBlocks datasource for searching IP blocks
type DataSourceAllIpBlocks struct {
	framework.BaseDataSource
}

// Schema returns the data source schema
func (r *DataSourceAllIpBlocks) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllIpBlocksSchema(ctx)
}

// Read searches for IP blocks using the provided filter, sort, and pagination
func (r *DataSourceAllIpBlocks) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceAllIpBlocksModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	var tffilter FilterModel
	if diags := data.Filter.As(ctx, &tffilter, basetypes.ObjectAsOptions{}); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	values := make([]string, 0, len(tffilter.Values))
	for _, v := range tffilter.Values {
		if !v.IsNull() && !v.IsUnknown() {
			values = append(values, v.ValueString())
		}
	}

	filterItem := fabricv4.IpBlockAndQuery{
		Property: tffilter.Property.ValueString(),
		Operator: fabricv4.ExchangeServicePropertyExpressionOperator(tffilter.Operator.ValueString()),
		Values:   values,
	}

	filter := fabricv4.IpBlockFilter{
		And: []fabricv4.IpBlockAndQuery{filterItem},
	}

	var tfpagination PaginationModel
	if diags := data.Pagination.As(ctx, &tfpagination, basetypes.ObjectAsOptions{}); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	offset := tfpagination.Offset.ValueInt32()
	limit := tfpagination.Limit.ValueInt32()
	if limit == 0 {
		limit = 20
	}

	pagination := fabricv4.PaginationRequest{
		Offset: &offset,
		Limit:  &limit,
	}

	searchRequest := fabricv4.IpBlocksSearchRequestBody{
		Filter:     &filter,
		Pagination: &pagination,
	}

	var tfsort SortModel
	if diags := data.Sort.As(ctx, &tfsort, basetypes.ObjectAsOptions{}); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	if !tfsort.Direction.IsNull() && !tfsort.Direction.IsUnknown() &&
		!tfsort.Property.IsNull() && !tfsort.Property.IsUnknown() {
		searchRequest.Sort = []fabricv4.SearchSortItem{
			{
				Direction: fabricv4.SearchSortItemDirection(tfsort.Direction.ValueString()),
				Property:  tfsort.Property.ValueString(),
			},
		}
	}

	ipBlocks, _, err := client.IPBlocksApi.SearchIpBlocks(ctx).
		IpBlocksSearchRequestBody(searchRequest).
		Execute()
	if err != nil {
		response.Diagnostics.AddError("api error retrieving ip blocks", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, ipBlocks)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
