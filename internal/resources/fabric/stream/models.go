package stream

import (
	"context"

	"github.com/equinix/terraform-provider-equinix/internal/fabric"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DataSourceByIDModel struct {
	StreamID types.String `tfsdk:"stream_id"`
	ID       types.String `tfsdk:"id"`
	BaseStreamModel
}

type DataSourceAllStreamsModel struct {
	ID         types.String                                     `tfsdk:"id"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]           `tfsdk:"pagination"`
	Data       fwtypes.ListNestedObjectValueOf[BaseStreamModel] `tfsdk:"data"`
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
	BaseStreamModel
}

type BaseStreamModel struct {
	Type                     types.String                          `tfsdk:"type"`
	Name                     types.String                          `tfsdk:"name"`
	Description              types.String                          `tfsdk:"description"`
	Href                     types.String                          `tfsdk:"href"`
	UUID                     types.String                          `tfsdk:"uuid"`
	State                    types.String                          `tfsdk:"state"`
	AssetsCount              types.Int32                           `tfsdk:"assets_count"`
	StreamSubscriptionsCount types.Int32                           `tfsdk:"stream_subscriptions_count"`
	Project                  fwtypes.ObjectValueOf[ProjectModel]   `tfsdk:"project"`    // Object of ProjectModel
	ChangeLog                fwtypes.ObjectValueOf[ChangeLogModel] `tfsdk:"change_log"` // Object of ChangeLogModel
}

type ProjectModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

type ChangeLogModel struct {
	CreatedBy         types.String `tfsdk:"created_by"`
	CreatedByFullName types.String `tfsdk:"created_by_full_name"`
	CreatedByEmail    types.String `tfsdk:"created_by_email"`
	CreatedDateTime   types.String `tfsdk:"created_date_time"`
	UpdatedBy         types.String `tfsdk:"updated_by"`
	UpdatedByFullName types.String `tfsdk:"updated_by_full_name"`
	UpdatedByEmail    types.String `tfsdk:"updated_by_email"`
	UpdatedDateTime   types.String `tfsdk:"updated_date_time"`
	DeletedBy         types.String `tfsdk:"deleted_by"`
	DeletedByFullName types.String `tfsdk:"deleted_by_full_name"`
	DeletedByEmail    types.String `tfsdk:"deleted_by_email"`
	DeletedDateTime   types.String `tfsdk:"deleted_date_time"`
}

func (m *DataSourceByIDModel) parse(ctx context.Context, stream *fabricv4.Stream) diag.Diagnostics {

	m.StreamID = types.StringValue(stream.GetUuid())
	m.ID = types.StringValue(stream.GetUuid())

	diags := parseStream(ctx, stream,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.UUID,
		&m.State,
		&m.AssetsCount,
		&m.StreamSubscriptionsCount,
		&m.Project,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}

	return diags
}

func (m *DataSourceAllStreamsModel) parse(ctx context.Context, streamsResponse *fabricv4.GetAllStreamResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(streamsResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by streams data source",
			"either the account does not have any streams data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]BaseStreamModel, len(streamsResponse.GetData()))
	streams := streamsResponse.GetData()
	for index, stream := range streams {
		var streamModel BaseStreamModel
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
	m.Data = fwtypes.NewListNestedObjectValueOfValueSlice[BaseStreamModel](ctx, data)

	return diags
}

func (m *ResourceModel) parse(ctx context.Context, stream *fabricv4.Stream) diag.Diagnostics {
	m.ID = types.StringValue(stream.GetUuid())

	diags := parseStream(ctx, stream,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.UUID,
		&m.State,
		&m.AssetsCount,
		&m.StreamSubscriptionsCount,
		&m.Project,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}

	return diags
}

func (m *BaseStreamModel) parse(ctx context.Context, stream *fabricv4.Stream) diag.Diagnostics {
	diags := parseStream(ctx, stream,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.UUID,
		&m.State,
		&m.AssetsCount,
		&m.StreamSubscriptionsCount,
		&m.Project,
		&m.ChangeLog)
	if diags.HasError() {
		return diags
	}

	return diags
}

func parseStream(ctx context.Context, stream *fabricv4.Stream,
	streamType, name, description, href, uuid, state *basetypes.StringValue,
	assetsCount, streamSubscriptionCount *basetypes.Int32Value,
	project *fwtypes.ObjectValueOf[ProjectModel],
	changeLog *fwtypes.ObjectValueOf[ChangeLogModel]) diag.Diagnostics {

	var diag diag.Diagnostics

	*streamType = types.StringValue(string(stream.GetType()))
	*name = types.StringValue(stream.GetName())
	*description = types.StringValue(stream.GetDescription())
	*href = types.StringValue(stream.GetHref())
	*uuid = types.StringValue(stream.GetUuid())
	*state = types.StringValue(stream.GetState())
	*assetsCount = types.Int32Value(stream.GetAssetsCount())
	*streamSubscriptionCount = types.Int32Value(stream.GetStreamSubscriptionsCount())

	streamProject := stream.GetProject()
	projectModel := ProjectModel{
		ProjectID: types.StringValue(streamProject.GetProjectId()),
	}
	*project = fwtypes.NewObjectValueOf[ProjectModel](ctx, &projectModel)

	streamChangeLog := stream.GetChangeLog()
	changeLogModel := ChangeLogModel{
		CreatedBy:         types.StringValue(streamChangeLog.GetCreatedBy()),
		CreatedByFullName: types.StringValue(streamChangeLog.GetCreatedByFullName()),
		CreatedByEmail:    types.StringValue(streamChangeLog.GetCreatedByEmail()),
		CreatedDateTime:   types.StringValue(streamChangeLog.GetCreatedDateTime().Format(fabric.TimeFormat)),
		UpdatedBy:         types.StringValue(streamChangeLog.GetUpdatedBy()),
		UpdatedByFullName: types.StringValue(streamChangeLog.GetUpdatedByFullName()),
		UpdatedByEmail:    types.StringValue(streamChangeLog.GetUpdatedByEmail()),
		UpdatedDateTime:   types.StringValue(streamChangeLog.GetUpdatedDateTime().Format(fabric.TimeFormat)),
		DeletedBy:         types.StringValue(streamChangeLog.GetDeletedBy()),
		DeletedByFullName: types.StringValue(streamChangeLog.GetDeletedByFullName()),
		DeletedByEmail:    types.StringValue(streamChangeLog.GetDeletedByEmail()),
		DeletedDateTime:   types.StringValue(streamChangeLog.GetDeletedDateTime().Format(fabric.TimeFormat)),
	}
	*changeLog = fwtypes.NewObjectValueOf[ChangeLogModel](ctx, &changeLogModel)
	return diag
}
