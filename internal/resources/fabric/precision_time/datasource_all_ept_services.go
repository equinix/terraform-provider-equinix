package precisiontime

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewDataSourceAllEptServices() datasource.DataSource {
	return &DataSourceAllEptServices{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_precision_time_services",
			},
		),
	}
}

type DataSourceAllEptServices struct {
	framework.BaseDataSource
}

func (r *DataSourceAllEptServices) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllEptServicesSchema(ctx)
}

func (r *DataSourceAllEptServices) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	var data dataSourceAllEptServicesModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	searchRequest, diags := buildSearchRequest(ctx, data)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	eptServices, _, err := client.PrecisionTimeApi.SearchTimeServices(ctx).TimeServicesSearchRequest(searchRequest).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving ept services data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	response.Diagnostics.Append(data.parse(ctx, eptServices)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func buildSearchRequest(ctx context.Context, plan dataSourceAllEptServicesModel) (fabricv4.TimeServicesSearchRequest, diag.Diagnostics) {
	var pagination paginationModel
	diags := plan.Pagination.As(ctx, &pagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.TimeServicesSearchRequest{}, diags
	}
	offset := pagination.Offset.ValueInt32()
	limit := pagination.Limit.ValueInt32()
	if limit == 0 {
		limit = 20
	}

	paginationRequest := fabricv4.PaginationRequest{
		Offset: &offset,
		Limit:  &limit,
	}

	searchRequest := fabricv4.TimeServicesSearchRequest{}
	searchRequest.SetPagination(paginationRequest)

	if !plan.Filter.IsNull() && !plan.Filter.IsUnknown() {
		filterModels := make([]filterModel, len(plan.Filter.Elements()))
		diags = plan.Filter.ElementsAs(ctx, &filterModels, false)
		if diags.HasError() {
			return fabricv4.TimeServicesSearchRequest{}, diags
		}
		var eptFilter fabricv4.TimeServiceFilters
		var filters []fabricv4.TimeServiceFilter
		var orFilter fabricv4.TimeServiceOrFilter
		for _, filter := range filterModels {
			var expression fabricv4.TimeServiceSimpleExpression
			expression.SetOperator(filter.Operator.ValueString())
			expression.SetProperty(filter.Property.ValueString())
			var values []string
			diags = filter.Values.ElementsAs(ctx, &values, false)
			if diags.HasError() {
				return fabricv4.TimeServicesSearchRequest{}, diags
			}
			expression.SetValues(values)
			if filter.Or.ValueBool() {
				orFilter.SetOr(append(orFilter.GetOr(), expression))
			} else {
				filters = append(filters, fabricv4.TimeServiceFilter{
					TimeServiceSimpleExpression: &expression,
				})
			}
		}

		if len(orFilter.GetOr()) > 0 {
			filters = append(filters, fabricv4.TimeServiceFilter{
				TimeServiceOrFilter: &orFilter,
			})
		}
		eptFilter.SetAnd(filters)
		searchRequest.SetFilter(eptFilter)
	}

	if !plan.Sort.IsNull() && plan.Sort.IsUnknown() {
		sortModels := make([]sortModel, len(plan.Sort.Elements()))
		diags = plan.Sort.ElementsAs(ctx, &sortModels, false)
		if diags.HasError() {
			return fabricv4.TimeServicesSearchRequest{}, diags
		}
		assetSort := make([]fabricv4.TimeServiceSortCriteria, len(sortModels))
		for i, criteria := range sortModels {
			sort := fabricv4.TimeServiceSortCriteria{}
			sort.SetDirection(fabricv4.TimeServiceSortDirection(criteria.Direction.ValueString()))
			sort.SetProperty(fabricv4.TimeServiceSortBy(criteria.Property.ValueString()))
			assetSort[i] = sort
		}
		searchRequest.SetSort(assetSort)
	}

	return searchRequest, diags
}
