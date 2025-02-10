package streamattachment

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func NewDataSourceAllStreamAttachments() datasource.DataSource {
	return &DataSourceAllStreamAttachments{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_stream_attachments",
			},
		),
	}
}

type DataSourceAllStreamAttachments struct {
	framework.BaseDataSource
}

func (r *DataSourceAllStreamAttachments) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceAllStreamAttachmentsSchema(ctx)
}

func (r *DataSourceAllStreamAttachments) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	// Retrieve values from plan
	var data DataSourceAllAssetsModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	searchRequest, diags := buildSearchRequest(ctx, data)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	// Use API client to get the current state of the resource
	assets, _, err := client.StreamsApi.GetStreamsAssets(ctx).StreamAssetSearchRequest(searchRequest).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving stream assets data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(data.parse(ctx, assets)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func buildSearchRequest(ctx context.Context, plan DataSourceAllAssetsModel) (fabricv4.StreamAssetSearchRequest, diag.Diagnostics) {
	var pagination PaginationModel
	diags := plan.Pagination.As(ctx, &pagination, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return fabricv4.StreamAssetSearchRequest{}, diags
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

	searchRequest := fabricv4.StreamAssetSearchRequest{}
	searchRequest.SetPagination(paginationRequest)

	if !plan.Filters.IsNull() && !plan.Filters.IsUnknown() {
		filterModels := make([]FilterModel, len(plan.Filters.Elements()))
		diags = plan.Filters.ElementsAs(ctx, &filterModels, false)
		if diags.HasError() {
			return fabricv4.StreamAssetSearchRequest{}, diags
		}
		var assetFilter fabricv4.StreamAssetFilters
		var filters []fabricv4.StreamAssetFilter
		var orFilter fabricv4.StreamAssetOrFilter
		for _, filter := range filterModels {
			var expression fabricv4.StreamAssetSimpleExpression
			expression.SetOperator(filter.Operator.ValueString())
			expression.SetProperty(filter.Property.ValueString())
			var values []string
			diags = filter.Values.ElementsAs(ctx, &values, false)
			if diags.HasError() {
				return fabricv4.StreamAssetSearchRequest{}, diags
			}
			expression.SetValues(values)
			if filter.Or.ValueBool() {
				orFilter.SetOr(append(orFilter.GetOr(), expression))
			} else {
				filters = append(filters, fabricv4.StreamAssetFilter{
					StreamAssetSimpleExpression: &expression,
				})
			}
		}

		if len(orFilter.GetOr()) > 0 {
			filters = append(filters, fabricv4.StreamAssetFilter{
				StreamAssetOrFilter: &orFilter,
			})
		}
		assetFilter.SetAnd(filters)
		searchRequest.SetFilter(assetFilter)
	}

	if !plan.Sort.IsNull() && plan.Sort.IsUnknown() {
		sortModels := make([]SortModel, len(plan.Sort.Elements()))
		diags = plan.Sort.ElementsAs(ctx, &sortModels, false)
		if diags.HasError() {
			return fabricv4.StreamAssetSearchRequest{}, diags
		}
		assetSort := make([]fabricv4.StreamAssetSortCriteria, len(sortModels))
		for i, criteria := range sortModels {
			sort := fabricv4.StreamAssetSortCriteria{}
			sort.SetDirection(fabricv4.StreamAssetSortDirection(criteria.Direction.ValueString()))
			sort.SetProperty(fabricv4.StreamAssetSortBy(criteria.Property.ValueString()))
			assetSort[i] = sort
		}
		searchRequest.SetSort(assetSort)
	}

	return searchRequest, diags
}
