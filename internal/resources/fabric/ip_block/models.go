package ipblock

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const timeFormat = "2006-01-02T15:04:05.000Z"

// ---- datasource by UUID ----

type DataSourceByIDModel struct {
	IpBlockID types.String `tfsdk:"ip_block_id"`
	ID        types.String `tfsdk:"id"`
	BaseIpBlockModel
}

// ---- datasource search ----

type DataSourceAllIpBlocksModel struct {
	ID         types.String                                      `tfsdk:"id"`
	Data       fwtypes.ListNestedObjectValueOf[BaseIpBlockModel] `tfsdk:"data"`
	Filter     types.Object                                      `tfsdk:"filter"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]            `tfsdk:"pagination"`
	Sort       fwtypes.ObjectValueOf[SortModel]                  `tfsdk:"sort"`
}

type FilterModel struct {
	Property types.String   `tfsdk:"property"`
	Operator types.String   `tfsdk:"operator"`
	Values   []types.String `tfsdk:"values"`
}

type PaginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type SortModel struct {
	Direction types.String `tfsdk:"direction"`
	Property  types.String `tfsdk:"property"`
}

// ---- base model ----

type BaseIpBlockModel struct {
	UUID         types.String                                `tfsdk:"uuid"`
	Href         types.String                                `tfsdk:"href"`
	Type         types.String                                `tfsdk:"type"`
	State        types.String                                `tfsdk:"state"`
	Ownership    types.String                                `tfsdk:"ownership"`
	PrefixLength types.Int32                                 `tfsdk:"prefix_length"`
	Prefix       types.String                                `tfsdk:"prefix"`
	Location     fwtypes.ObjectValueOf[LocationModel]        `tfsdk:"location"`
	Order        fwtypes.ObjectValueOf[OrderModel]           `tfsdk:"order"`
	Account      fwtypes.ObjectValueOf[AccountModel]         `tfsdk:"account"`
	Project      fwtypes.ObjectValueOf[ProjectModel]         `tfsdk:"project"`
	Assets       fwtypes.ListNestedObjectValueOf[AssetModel] `tfsdk:"assets"`
	Change       fwtypes.ObjectValueOf[ChangeModel]          `tfsdk:"change"`
	ChangeLog    fwtypes.ObjectValueOf[ChangeLogModel]       `tfsdk:"change_log"`
}

type LocationModel struct {
	MetroHref types.String `tfsdk:"metro_href"`
	MetroCode types.String `tfsdk:"metro_code"`
}

type OrderModel struct {
	Href        types.String `tfsdk:"href"`
	OrderNumber types.String `tfsdk:"order_number"`
}

type AccountModel struct {
	AccountNumber types.String `tfsdk:"account_number"`
}

type ProjectModel struct {
	Href      types.String `tfsdk:"href"`
	ProjectID types.String `tfsdk:"project_id"`
}

type AssetModel struct {
	Type types.String `tfsdk:"type"`
	UUID types.String `tfsdk:"uuid"`
	Href types.String `tfsdk:"href"`
}

type ChangeModel struct {
	Href types.String `tfsdk:"href"`
}

type ChangeLogModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	UpdatedDateTime types.String `tfsdk:"updated_date_time"`
	DeletedDateTime types.String `tfsdk:"deleted_date_time"`
}

// ---- parse functions ----

func (m *DataSourceByIDModel) parse(ctx context.Context, ipBlock *fabricv4.IpBlock) diag.Diagnostics {
	m.IpBlockID = types.StringValue(ipBlock.GetUuid())
	m.ID = types.StringValue(ipBlock.GetUuid())
	return m.BaseIpBlockModel.parse(ctx, ipBlock)
}

func (m *DataSourceAllIpBlocksModel) parse(ctx context.Context, resp *fabricv4.IpBlockSearchResponseBody) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(resp.GetData()) < 1 {
		diags.AddError("no data retrieved by ip blocks data source",
			"either the account does not have any ip blocks or the combination of filter, limit and offset needs to be updated")
		return diags
	}

	blocks := resp.GetData()
	data := make([]*BaseIpBlockModel, len(blocks))
	for i, block := range blocks {
		var model BaseIpBlockModel
		diags = model.parse(ctx, &block)
		if diags.HasError() {
			return diags
		}
		data[i] = &model
	}

	responsePagination := resp.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}

	m.ID = types.StringValue(data[0].UUID.ValueString())
	m.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[BaseIpBlockModel](ctx, data)

	return diags
}

func (m *BaseIpBlockModel) parse(ctx context.Context, ipBlock *fabricv4.IpBlock) diag.Diagnostics {
	var diags diag.Diagnostics

	m.UUID = types.StringValue(ipBlock.GetUuid())
	m.Href = types.StringValue(ipBlock.GetHref())
	m.Type = types.StringValue(string(ipBlock.GetType()))
	m.State = types.StringValue(string(ipBlock.GetState()))
	m.Ownership = types.StringValue(string(ipBlock.GetOwnership()))
	m.PrefixLength = types.Int32Value(ipBlock.GetPrefixLength())
	m.Prefix = types.StringValue(ipBlock.GetPrefix())

	if loc, ok := ipBlock.GetLocationOk(); ok && loc != nil {
		m.Location = fwtypes.NewObjectValueOf[LocationModel](ctx, &LocationModel{
			MetroHref: types.StringValue(loc.GetMetroHref()),
			MetroCode: types.StringValue(loc.GetMetroCode()),
		})
	} else {
		m.Location = fwtypes.NewObjectValueOfNull[LocationModel](ctx)
	}

	if order, ok := ipBlock.GetOrderOk(); ok && order != nil {
		m.Order = fwtypes.NewObjectValueOf[OrderModel](ctx, &OrderModel{
			Href:        types.StringValue(order.GetHref()),
			OrderNumber: types.StringValue(order.GetOrderNumber()),
		})
	} else {
		m.Order = fwtypes.NewObjectValueOfNull[OrderModel](ctx)
	}

	if account, ok := ipBlock.GetAccountOk(); ok && account != nil {
		m.Account = fwtypes.NewObjectValueOf[AccountModel](ctx, &AccountModel{
			AccountNumber: types.StringValue(account.GetAccountNumber()),
		})
	} else {
		m.Account = fwtypes.NewObjectValueOfNull[AccountModel](ctx)
	}

	project := ipBlock.GetProject()
	m.Project = fwtypes.NewObjectValueOf[ProjectModel](ctx, &ProjectModel{
		Href:      types.StringValue(project.GetHref()),
		ProjectID: types.StringValue(project.GetProjectId()),
	})

	assets := ipBlock.GetAssets()
	assetModels := make([]*AssetModel, len(assets))
	for i, asset := range assets {
		assetModels[i] = &AssetModel{
			Type: types.StringValue(string(asset.GetType())),
			UUID: types.StringValue(asset.GetUuid()),
			Href: types.StringValue(asset.GetHref()),
		}
	}
	m.Assets = fwtypes.NewListNestedObjectValueOfSlice[AssetModel](ctx, assetModels)

	change := ipBlock.GetChange()
	m.Change = fwtypes.NewObjectValueOf[ChangeModel](ctx, &ChangeModel{
		Href: types.StringValue(change.GetHref()),
	})

	cl := ipBlock.GetChangeLog()
	clModel := ChangeLogModel{
		CreatedDateTime: types.StringValue(cl.GetCreatedDateTime().Format(timeFormat)),
	}
	if updatedAt, ok := cl.GetUpdatedDateTimeOk(); ok && updatedAt != nil {
		clModel.UpdatedDateTime = types.StringValue(updatedAt.Format(timeFormat))
	} else {
		clModel.UpdatedDateTime = types.StringNull()
	}
	if deletedAt, ok := cl.GetDeletedDateTimeOk(); ok && deletedAt != nil {
		clModel.DeletedDateTime = types.StringValue(deletedAt.Format(timeFormat))
	} else {
		clModel.DeletedDateTime = types.StringNull()
	}
	m.ChangeLog = fwtypes.NewObjectValueOf[ChangeLogModel](ctx, &clModel)

	return diags
}
