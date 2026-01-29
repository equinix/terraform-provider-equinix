package received_route

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type receivedRoutesBaseModel struct {
	Type            types.String                           `tfsdk:"type"`
	ProtocolType    types.String                           `tfsdk:"protocol_type"`
	State           types.String                           `tfsdk:"state"`
	Prefix          types.String                           `tfsdk:"prefix"`
	NextHop         types.String                           `tfsdk:"next_hop"`
	MED             types.Int32                            `tfsdk:"med"`
	LocalPreference types.Int32                            `tfsdk:"local_preference"`
	AsPath          fwtypes.ListValueOf[types.String]      `tfsdk:"as_path"`
	Connection      fwtypes.ObjectValueOf[connectionModel] `tfsdk:"connection"`
	Changelog       fwtypes.ObjectValueOf[changeLogModel]  `tfsdk:"change_log"`
}

type connectionModel struct {
	Href types.String `tfsdk:"href"`
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

type changeLogModel struct {
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

type paginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type sortModel struct {
	Direction types.String `tfsdk:"direction"`
	Property  types.String `tfsdk:"property"`
}

type FilterModel struct {
	Property types.String   `tfsdk:"property"`
	Operator types.String   `tfsdk:"operator"`
	Values   []types.String `tfsdk:"values"`
}

type dataSourceSearchReceivedRoutesModel struct {
	ID           types.String                                             `tfsdk:"id"`
	ConnectionID types.String                                             `tfsdk:"connection_id"`
	Filter       types.Object                                             `tfsdk:"filter"`
	Data         fwtypes.ListNestedObjectValueOf[receivedRoutesBaseModel] `tfsdk:"data"`
	Pagination   fwtypes.ObjectValueOf[paginationModel]                   `tfsdk:"pagination"`
	Sort         fwtypes.ObjectValueOf[sortModel]                         `tfsdk:"sort"`
}

func (a *dataSourceSearchReceivedRoutesModel) parse(ctx context.Context, receivedRoutesResponse *fabricv4.ConnectionRouteTableEntrySearchResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(receivedRoutesResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by received routes data source",
			"either the connection does not have any received routes data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]receivedRoutesBaseModel, len(receivedRoutesResponse.GetData()))
	receivedRoutes := receivedRoutesResponse.GetData()
	for i, receivedRoute := range receivedRoutes {
		var receivedRoutesModel receivedRoutesBaseModel
		diags := receivedRoutesModel.parse(ctx, &receivedRoute)
		if diags.HasError() {
			return diags
		}
		data[i] = receivedRoutesModel
	}
	responsePagination := receivedRoutesResponse.GetPagination()
	pagination := paginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}

	// a.ID = types.StringValue(data[0].ValueString()) // correct?
	a.Pagination = fwtypes.NewObjectValueOf[paginationModel](ctx, &pagination)
	a.Data = fwtypes.NewListNestedObjectValueOfValueSlice[receivedRoutesBaseModel](ctx, data)

	return diags
}

func (a *receivedRoutesBaseModel) parse(ctx context.Context, receivedRoute *fabricv4.ConnectionRouteTableEntry) diag.Diagnostics {
	var diags diag.Diagnostics
	a.Type = types.StringValue(string(receivedRoute.GetType()))
	a.ProtocolType = types.StringValue(string(receivedRoute.GetProtocolType()))
	a.State = types.StringValue(string(receivedRoute.GetState()))
	a.Prefix = types.StringValue(receivedRoute.GetPrefix())
	a.NextHop = types.StringValue(receivedRoute.GetNextHop())
	a.MED = types.Int32Value(receivedRoute.GetMED())
	a.LocalPreference = types.Int32Value(receivedRoute.GetLocalPreference())
	a.AsPath, diags = parseAsPaths(ctx, receivedRoute.GetAsPath())
	if diags.HasError() {
		return diags
	}
	a.Connection, diags = parseConnection(ctx, receivedRoute.GetConnection())
	if diags.HasError() {
		return diags
	}

	a.Changelog, diags = parseChangelog(ctx, receivedRoute.GetChangeLog())
	if diags.HasError() {
		return diags
	}

	return diags
}

func parseAsPaths(ctx context.Context, asPaths []string) (fwtypes.ListValueOf[types.String], diag.Diagnostics) {
	var diags diag.Diagnostics
	asPathTypeList := make([]attr.Value, len(asPaths))

	for i, asPath := range asPaths {
		asPathTypeList[i] = types.StringValue(string(asPath))
	}
	asPathValue, diags := fwtypes.NewListValueOf[types.String](ctx, asPathTypeList)

	if diags.HasError() {
		return fwtypes.NewListValueOfNull[types.String](ctx), diags
	}
	return asPathValue, diags
}

func parseConnection(ctx context.Context, connection fabricv4.ConnectionRouteTableEntryConnection) (fwtypes.ObjectValueOf[connectionModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}

	result := connectionModel{}
	if connection.Href != nil {
		result.Href = types.StringValue(connection.GetHref())
	}

	if connection.Name != nil {
		result.Name = types.StringValue(connection.GetName())
	}

	if connection.Uuid != nil {
		result.UUID = types.StringValue(connection.GetUuid())
	}
	return fwtypes.NewObjectValueOf[connectionModel](ctx, &result), diags
}

func parseChangelog(ctx context.Context, changeLog fabricv4.Changelog) (fwtypes.ObjectValueOf[changeLogModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}

	result := changeLogModel{}
	if changeLog.CreatedBy != nil {
		result.CreatedBy = types.StringValue(changeLog.GetCreatedBy()) //Get functon not loading
	}

	if changeLog.CreatedByFullName != nil {
		result.CreatedByFullName = types.StringValue(changeLog.GetCreatedByFullName())
	}

	if changeLog.CreatedByEmail != nil {
		result.CreatedByEmail = types.StringValue(changeLog.GetCreatedByEmail())
	}

	if changeLog.CreatedDateTime != nil {
		result.CreatedDateTime = types.StringValue(changeLog.GetCreatedDateTime().String())
	}

	if changeLog.UpdatedBy != nil {
		result.UpdatedBy = types.StringValue(changeLog.GetUpdatedBy())
	}

	if changeLog.UpdatedByFullName != nil {
		result.UpdatedByFullName = types.StringValue(changeLog.GetUpdatedByFullName())
	}
	if changeLog.UpdatedByEmail != nil {
		result.UpdatedByEmail = types.StringValue(changeLog.GetUpdatedByEmail())
	}

	if changeLog.UpdatedDateTime != nil {
		result.UpdatedDateTime = types.StringValue(changeLog.GetUpdatedDateTime().String())
	}

	if changeLog.DeletedBy != nil {
		result.DeletedBy = types.StringValue(changeLog.GetDeletedBy())
	}

	if changeLog.DeletedByFullName != nil {
		result.DeletedByFullName = types.StringValue(changeLog.GetDeletedByFullName())
	}

	if changeLog.DeletedByEmail != nil {
		result.DeletedByEmail = types.StringValue(changeLog.GetDeletedByEmail())
	}

	if changeLog.DeletedDateTime != nil {
		result.DeletedDateTime = types.StringValue(changeLog.GetDeletedDateTime().String())
	}
	return fwtypes.NewObjectValueOf[changeLogModel](ctx, &result), diags
}
