package stream

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DataSourceByIdModel struct {
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
	Type                     types.String `tfsdk:"type"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	Href                     types.String `tfsdk:"href"`
	Uuid                     types.String `tfsdk:"uuid"`
	State                    types.String `tfsdk:"state"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	AssetsCount              types.Int32  `tfsdk:"assets_count"`
	StreamSubscriptionsCount types.Int32  `tfsdk:"stream_subscriptions_count"`
	Project                  types.Object `tfsdk:"project"`    // Object of ProjectModel
	ChangeLog                types.Object `tfsdk:"change_log"` // Object of ChangeLogModel
}

type ProjectModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

func (m ProjectModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project_id": types.StringType,
	}
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

func (m ChangeLogModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_by":           types.StringType,
		"created_by_full_name": types.StringType,
		"created_by_email":     types.StringType,
		"created_date_time":    types.StringType,
		"updated_by":           types.StringType,
		"updated_by_full_name": types.StringType,
		"updated_by_email":     types.StringType,
		"updated_date_time":    types.StringType,
		"deleted_by":           types.StringType,
		"deleted_by_full_name": types.StringType,
		"deleted_by_email":     types.StringType,
		"deleted_date_time":    types.StringType,
	}
}

func (m *DataSourceByIdModel) parse(ctx context.Context, stream *fabricv4.Stream) diag.Diagnostics {
	var diags diag.Diagnostics

	m.StreamID = types.StringValue(stream.GetUuid())
	m.ID = types.StringValue(stream.GetUuid())

	diags = parseStream(ctx, stream,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.Enabled,
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

	data := make([]BaseStreamModel, len(streamsResponse.GetData()))
	streams := streamsResponse.GetData()
	for _, stream := range streams {
		var streamModel BaseStreamModel
		diags = streamModel.parse(ctx, &stream)
		if diags.HasError() {
			return diags
		}
		data = append(data, streamModel)
	}
	responsePagination := streamsResponse.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}

	m.ID = types.StringValue(data[0].Uuid.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)
	m.Data = fwtypes.NewListNestedObjectValueOfValueSlice[BaseStreamModel](ctx, data)

	return diags
}

func (m *ResourceModel) parse(ctx context.Context, stream *fabricv4.Stream) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(stream.GetUuid())

	diags = parseStream(ctx, stream,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.Enabled,
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
	var diags diag.Diagnostics

	diags = parseStream(ctx, stream,
		&m.Type,
		&m.Name,
		&m.Description,
		&m.Href,
		&m.Uuid,
		&m.State,
		&m.Enabled,
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
	type_, name, description, href, uuid, state *basetypes.StringValue,
	enabled *basetypes.BoolValue,
	assetsCount, streamSubscriptionCount *basetypes.Int32Value,
	project, changeLog *basetypes.ObjectValue) diag.Diagnostics {

	var diag diag.Diagnostics

	*type_ = types.StringValue(string(stream.GetType()))
	*name = types.StringValue(stream.GetName())
	*description = types.StringValue(stream.GetDescription())
	*href = types.StringValue(stream.GetHref())
	*uuid = types.StringValue(stream.GetUuid())
	*state = types.StringValue(stream.GetState())
	*enabled = types.BoolValue(stream.GetEnabled())
	*assetsCount = types.Int32Value(stream.GetAssetsCount())
	*streamSubscriptionCount = types.Int32Value(stream.GetStreamSubscriptionsCount())

	streamProject := stream.GetProject()
	projectModel := ProjectModel{
		ProjectID: types.StringValue(streamProject.GetProjectId()),
	}
	terraformProject, diags := types.ObjectValueFrom(ctx, projectModel.AttributeTypes(), projectModel)
	if diags.HasError() {
		return diags
	}
	*project = terraformProject

	const TIMEFORMAT = "2006-01-02T15:04:05.000Z"
	streamChangeLog := stream.GetChangelog()
	changeLogModel := ChangeLogModel{
		CreatedBy:         types.StringValue(streamChangeLog.GetCreatedBy()),
		CreatedByFullName: types.StringValue(streamChangeLog.GetCreatedByFullName()),
		CreatedByEmail:    types.StringValue(streamChangeLog.GetCreatedByEmail()),
		CreatedDateTime:   types.StringValue(streamChangeLog.GetCreatedDateTime().Format(TIMEFORMAT)),
		UpdatedBy:         types.StringValue(streamChangeLog.GetUpdatedBy()),
		UpdatedByFullName: types.StringValue(streamChangeLog.GetUpdatedByFullName()),
		UpdatedByEmail:    types.StringValue(streamChangeLog.GetUpdatedByEmail()),
		UpdatedDateTime:   types.StringValue(streamChangeLog.GetUpdatedDateTime().Format(TIMEFORMAT)),
		DeletedBy:         types.StringValue(streamChangeLog.GetDeletedBy()),
		DeletedByFullName: types.StringValue(streamChangeLog.GetDeletedByFullName()),
		DeletedByEmail:    types.StringValue(streamChangeLog.GetDeletedByEmail()),
		DeletedDateTime:   types.StringValue(streamChangeLog.GetDeletedDateTime().Format(TIMEFORMAT)),
	}
	terraformChangeLog, diags := types.ObjectValueFrom(ctx, changeLogModel.AttributeTypes(), changeLogModel)
	if diags.HasError() {
		return diags
	}
	*changeLog = terraformChangeLog
	return diag
}
