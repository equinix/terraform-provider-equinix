package asn

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewDataSourceAllAsns() datasource.DataSource {
	return &DataSourceAllAsns{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_asns",
			},
		),
	}
}

type DataSourceAllAsns struct {
	framework.BaseDataSource
}

func (r *DataSourceAllAsns) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllAsnsSchema(ctx)
}

func (r *DataSourceAllAsns) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceAllAsnsModel

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

	property := tffilter.Property.ValueString()
	operator := fabricv4.SearchExpressionOperator(tffilter.Operator.ValueString())
	filter := fabricv4.SearchExpression{
		And: []fabricv4.SearchExpression{
			{
				Property: &property,
				Operator: &operator,
				Values:   values,
			},
		},
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

	searchRequest := fabricv4.AsnSearchRequestBody{
		Filter: &filter,
		Pagination: &fabricv4.PaginationRequest{
			Offset: &offset,
			Limit:  &limit,
		},
	}

	var tfsort SortModel
	if diags := data.Sort.As(ctx, &tfsort, basetypes.ObjectAsOptions{}); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	if !tfsort.Direction.IsNull() && !tfsort.Direction.IsUnknown() &&
		!tfsort.Property.IsNull() && !tfsort.Property.IsUnknown() {
		sortDirection := fabricv4.AsnSearchSortCriteriaDirection(tfsort.Direction.ValueString())
		sortProperty := tfsort.Property.ValueString()
		searchRequest.Sort = []fabricv4.AsnSearchSortCriteria{
			{
				Direction: &sortDirection,
				Property:  &sortProperty,
			},
		}
	}

	asns, _, err := client.ASNsApi.SearchAsns(ctx).
		AsnSearchRequestBody(searchRequest).
		Execute()
	if err != nil {
		response.Diagnostics.AddError("api error retrieving asns", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, asns)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
