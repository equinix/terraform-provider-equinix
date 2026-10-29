package internetaccess

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewDataSourceAllInternetAccess creates a new data source for searching Internet Access services
func NewDataSourceAllInternetAccess() datasource.DataSource {
	return &DataSourceAllInternetAccess{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_internet_access_services",
			},
		),
	}
}

// DataSourceAllInternetAccess datasource for searching Internet Access services
type DataSourceAllInternetAccess struct {
	framework.BaseDataSource
}

// Schema returns the data source schema
func (r *DataSourceAllInternetAccess) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllInternetAccessSchema(ctx)
}

// Read searches for Internet Access services using the provided filter, sort, and pagination
func (r *DataSourceAllInternetAccess) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceAllInternetAccessModel

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

	searchRequest := fabricv4.InternetAccessSearchRequest{
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
		searchRequest.Sort = []fabricv4.SearchSortItem{
			{
				Direction: fabricv4.SearchSortItemDirection(tfsort.Direction.ValueString()),
				Property:  tfsort.Property.ValueString(),
			},
		}
	}

	services, _, err := client.InternetAccessServicesApi.SearchEiaServices(ctx).
		InternetAccessSearchRequest(searchRequest).
		Execute()
	if err != nil {
		response.Diagnostics.AddError("api error retrieving internet access services", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, services)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
