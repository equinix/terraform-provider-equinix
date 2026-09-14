// Package internetaccess implements the data source for searching Equinix Internet Access (EIA) services
package internetaccess

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// NewDataSourceAllInternetAccessServices creates a new data source for searching Internet Access services
func NewDataSourceAllInternetAccessServices() datasource.DataSource {
	return &DataSourceAllInternetAccessServices{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_internet_access_services",
			},
		),
	}
}

// DataSourceAllInternetAccessServices datasource for searching EIA services
type DataSourceAllInternetAccessServices struct {
	framework.BaseDataSource
}

// Schema returns the data source schema
func (r *DataSourceAllInternetAccessServices) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllInternetAccessServicesSchema(ctx)
}

// Read searches for Internet Access services using the provided filter, sort, and pagination
func (r *DataSourceAllInternetAccessServices) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data DataSourceAllInternetAccessServicesModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	var tffilters []FilterModel
	if diags := data.Filter.ElementsAs(ctx, &tffilters, false); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	andConditions := make([]fabricv4.SearchExpression, 0, len(tffilters))
	for _, f := range tffilters {
		values := make([]string, 0, len(f.Values))
		for _, v := range f.Values {
			if !v.IsNull() && !v.IsUnknown() {
				values = append(values, v.ValueString())
			}
		}
		prop := f.Property.ValueString()
		op := fabricv4.SearchExpressionOperator(f.Operator.ValueString())
		andConditions = append(andConditions, fabricv4.SearchExpression{
			Property: &prop,
			Operator: &op,
			Values:   values,
		})
	}

	filterExpr := fabricv4.SearchExpression{
		And: andConditions,
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

	searchRequest := fabricv4.InternetAccessSearchRequest{
		Filter:     &filterExpr,
		Pagination: &pagination,
	}

	if !data.Sort.IsNull() && !data.Sort.IsUnknown() {
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
	}

	services, err := searchInternetAccessServicesRaw(ctx, client, searchRequest)
	if err != nil {
		response.Diagnostics.AddError("api error retrieving internet access services", err.Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, services)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
