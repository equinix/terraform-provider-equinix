package connectionrouteaggregation

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceByIDModel struct {
	ID types.String `tfsdk:"id"`
	baseConnectionRouteAggregationModel
}

type datsSourceAllConnectionRouteAggregationModel struct {
	ID           types.String                                                         `tfsdk:"id"`
	ConnectionID types.String                                                         `tfsdk:"connection_id"`
	Data         fwtypes.ListNestedObjectValueOf[baseConnectionRouteAggregationModel] `tfsdk:"data"`
	Pagination   fwtypes.ObjectValueOf[paginationModel]                               `tfsdk:"pagination"`
}

type paginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type resourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	baseConnectionRouteAggregationModel
}

type baseConnectionRouteAggregationModel struct {
	RouteAggregationID types.String `tfsdk:"route_aggregation_id"`
	ConnectionID       types.String `tfsdk:"connection_id"`
	Href               types.String `tfsdk:"href"`
	Type               types.String `tfsdk:"type"`
	UUID               types.String `tfsdk:"uuid"`
	AttachmentStatus   types.String `tfsdk:"attachment_status"`
}

func (m *dataSourceByIDModel) parse(ctx context.Context, connectionRouteAggregation *fabricv4.ConnectionRouteAggregationData) diag.Diagnostics {
	m.ID = types.StringValue(connectionRouteAggregation.GetUuid())
	diags := m.baseConnectionRouteAggregationModel.parse(ctx, connectionRouteAggregation)
	return diags
}

func (m *datsSourceAllConnectionRouteAggregationModel) parse(ctx context.Context, connectionRouteAggregationsResponse *fabricv4.GetAllConnectionRouteAggregationsResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(connectionRouteAggregationsResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by connection route aggrgeations data source", "either the account does not have any connection route aggregation data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]baseConnectionRouteAggregationModel, len(connectionRouteAggregationsResponse.GetData()))
	connectionRouteAggregations := connectionRouteAggregationsResponse.GetData()
	for index, routeAggregationRule := range connectionRouteAggregations {
		var connectionRouteAggregationModel baseConnectionRouteAggregationModel
		diags = connectionRouteAggregationModel.parse(ctx, &routeAggregationRule)
		if diags.HasError() {
			return diags
		}
		data[index] = connectionRouteAggregationModel
	}
	responsePagination := connectionRouteAggregationsResponse.GetPagination()
	pagination := paginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}
	m.ID = types.StringValue(data[0].UUID.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[paginationModel](ctx, &pagination)

	dataPtr := make([]*baseConnectionRouteAggregationModel, len(data))
	for i := range data {
		dataPtr[i] = &data[i]
	}
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[baseConnectionRouteAggregationModel](ctx, dataPtr)

	return diags
}

func (m *resourceModel) parse(ctx context.Context, connectionRouteAggregation *fabricv4.ConnectionRouteAggregationData) diag.Diagnostics {
	m.ID = types.StringValue(connectionRouteAggregation.GetUuid())
	diags := m.baseConnectionRouteAggregationModel.parse(ctx, connectionRouteAggregation)
	return diags
}

func (m *baseConnectionRouteAggregationModel) parse(_ context.Context, connectionRouteAggregation *fabricv4.ConnectionRouteAggregationData) diag.Diagnostics {
	var diag diag.Diagnostics

	m.Href = types.StringValue(connectionRouteAggregation.GetHref())
	m.Type = types.StringValue(string(connectionRouteAggregation.GetType()))
	m.UUID = types.StringValue(connectionRouteAggregation.GetUuid())
	m.AttachmentStatus = types.StringValue(string(connectionRouteAggregation.GetAttachmentStatus()))

	return diag
}
