package connection_route_aggregation

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DataSourceByIdModel struct {
	ID types.String `tfsdk:"id"`
	BaseConnectionRouteAggregationModel
}

type DatsSourceAllConnectionRouteAggregationModel struct {
	ID           types.String                                                         `tfsdk:"id"`
	ConnectionId types.String                                                         `tfsdk:"connection_id"`
	Data         fwtypes.ListNestedObjectValueOf[BaseConnectionRouteAggregationModel] `tfsdk:"data"`
	Pagination   fwtypes.ObjectValueOf[PaginationModel]                               `tfsdk:"pagination"`
}

type PaginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type ResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	BaseConnectionRouteAggregationModel
}

type BaseConnectionRouteAggregationModel struct {
	RouteAggregationId types.String `tfsdk:"route_aggregation_id"`
	ConnectionId       types.String `tfsdk:"connection_id"`
	Href               types.String `tfsdk:"href"`
	Type               types.String `tfsdk:"type"`
	Uuid               types.String `tfsdk:"uuid"`
	AttachmentStatus   types.String `tfsdk:"attachment_status"`
}

func (m *DataSourceByIdModel) parse(ctx context.Context, connectionRouteAggregation *fabricv4.ConnectionRouteAggregationData) diag.Diagnostics {
	m.ID = types.StringValue(connectionRouteAggregation.GetUuid())

	diags := parseConnectionRouteAggregation(ctx, connectionRouteAggregation,
		&m.Href,
		&m.Type,
		&m.Uuid,
		&m.AttachmentStatus)
	if diags.HasError() {
		return diags
	}
	return diags
}

func (m *DatsSourceAllConnectionRouteAggregationModel) parse(ctx context.Context, connectionRouteAggregationsResponse *fabricv4.GetAllConnectionRouteAggregationsResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(connectionRouteAggregationsResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by connection route aggrgeations data source", "either the account does not have any connection route aggregation data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]BaseConnectionRouteAggregationModel, len(connectionRouteAggregationsResponse.GetData()))
	connectionRouteAggregations := connectionRouteAggregationsResponse.GetData()
	for index, routeAggregationRule := range connectionRouteAggregations {
		var connectionRouteAggregationModel BaseConnectionRouteAggregationModel
		diags = connectionRouteAggregationModel.parse(ctx, &routeAggregationRule)
		if diags.HasError() {
			return diags
		}
		data[index] = connectionRouteAggregationModel
	}
	responsePagination := connectionRouteAggregationsResponse.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}
	m.ID = types.StringValue(data[0].Uuid.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)

	dataPtr := make([]*BaseConnectionRouteAggregationModel, len(data))
	for i := range data {
		dataPtr[i] = &data[i]
	}
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[BaseConnectionRouteAggregationModel](ctx, dataPtr)

	return diags
}

func (m *ResourceModel) parse(ctx context.Context, connectionRouteAggregation *fabricv4.ConnectionRouteAggregationData) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(connectionRouteAggregation.GetUuid())

	diags = parseConnectionRouteAggregation(ctx, connectionRouteAggregation,
		&m.Href,
		&m.Type,
		&m.Uuid,
		&m.AttachmentStatus)
	if diags.HasError() {
		return diags
	}
	return diags
}

func (m *BaseConnectionRouteAggregationModel) parse(ctx context.Context, connectionRouteAggregation *fabricv4.ConnectionRouteAggregationData) diag.Diagnostics {
	var diags diag.Diagnostics = parseConnectionRouteAggregation(ctx, connectionRouteAggregation,
		&m.Href,
		&m.Type,
		&m.Uuid,
		&m.AttachmentStatus)
	if diags.HasError() {
		return diags
	}
	return diags
}

func parseConnectionRouteAggregation(ctx context.Context, connectionRouteAggregation *fabricv4.ConnectionRouteAggregationData,
	href, type_, uuid, attachmentStatus *basetypes.StringValue) diag.Diagnostics {
	var diag diag.Diagnostics

	*href = types.StringValue(connectionRouteAggregation.GetHref())
	*type_ = types.StringValue(string(connectionRouteAggregation.GetType()))
	*uuid = types.StringValue(connectionRouteAggregation.GetUuid())
	*attachmentStatus = types.StringValue(string(connectionRouteAggregation.GetAttachmentStatus()))

	return diag
}
