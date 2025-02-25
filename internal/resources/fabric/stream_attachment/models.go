package streamattachment

import (
	"context"

	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceByIDModel struct {
	ID       types.String `tfsdk:"id"`
	StreamID types.String `tfsdk:"stream_id"`
	Asset    types.String `tfsdk:"asset"`
	AssetID  types.String `tfsdk:"asset_id"`
	BaseAssetModel
}

type DataSourceAllAssetsModel struct {
	ID         types.String                                    `tfsdk:"id"`
	Filters    fwtypes.ListNestedObjectValueOf[FilterModel]    `tfsdk:"filters"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]          `tfsdk:"pagination"`
	Sort       fwtypes.ListNestedObjectValueOf[SortModel]      `tfsdk:"sort"`
	Data       fwtypes.ListNestedObjectValueOf[BaseAssetModel] `tfsdk:"data"`
}

type FilterModel struct {
	Property types.String                      `tfsdk:"property"`
	Operator types.String                      `tfsdk:"operator"`
	Values   fwtypes.ListValueOf[types.String] `tfsdk:"values"`
	Or       types.Bool                        `tfsdk:"or"`
}

type SortModel struct {
	Direction types.String `tfsdk:"direction"`
	Property  types.String `tfsdk:"property"`
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
	StreamID types.String   `tfsdk:"stream_id"`
	Asset    types.String   `tfsdk:"asset"`
	AssetID  types.String   `tfsdk:"asset_id"`
	BaseAssetModel
}

type BaseAssetModel struct {
	MetricsEnabled   types.Bool   `tfsdk:"metrics_enabled"`
	Type             types.String `tfsdk:"type"`
	Href             types.String `tfsdk:"href"`
	UUID             types.String `tfsdk:"uuid"`
	AttachmentStatus types.String `tfsdk:"attachment_status"`
}

func (m *DataSourceByIDModel) parse(ctx context.Context, stream *fabricv4.StreamAsset) diag.Diagnostics {
	m.ID = types.StringValue(stream.GetUuid())

	diags := m.BaseAssetModel.parse(ctx, stream)

	return diags
}

func (m *DataSourceAllAssetsModel) parse(ctx context.Context, streamsResponse *fabricv4.GetAllStreamAssetResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(streamsResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by stream attachments data source",
			"either the account does not have any stream attachments data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]BaseAssetModel, len(streamsResponse.GetData()))
	streams := streamsResponse.GetData()
	for index, stream := range streams {
		var streamModel BaseAssetModel
		diags = streamModel.parse(ctx, &stream)
		if diags.HasError() {
			return diags
		}
		data[index] = streamModel
	}
	responsePagination := streamsResponse.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}

	m.ID = types.StringValue(data[0].UUID.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)
	m.Data = fwtypes.NewListNestedObjectValueOfValueSlice[BaseAssetModel](ctx, data)

	return diags
}

func (m *ResourceModel) parse(ctx context.Context, stream *fabricv4.StreamAsset) diag.Diagnostics {
	m.ID = types.StringValue(stream.GetUuid())

	diags := m.BaseAssetModel.parse(ctx, stream)

	return diags
}

func (m *BaseAssetModel) parse(_ context.Context, stream *fabricv4.StreamAsset) diag.Diagnostics {
	var diag diag.Diagnostics

	m.MetricsEnabled = types.BoolValue(stream.GetMetricsEnabled())
	m.Type = types.StringValue(string(stream.GetType()))
	m.Href = types.StringValue(stream.GetHref())
	m.UUID = types.StringValue(stream.GetUuid())
	m.AttachmentStatus = types.StringValue(string(stream.GetAttachmentStatus()))

	return diag
}
