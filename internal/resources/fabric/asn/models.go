package asn

import (
	"context"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const timeFormat = "2006-01-02T15:04:05.000Z"

type DataSourceByIDModel struct {
	AsnID types.String `tfsdk:"asn_id"`
	ID    types.String `tfsdk:"id"`
	BaseAsnModel
}

type DataSourceAllAsnsModel struct {
	ID         types.String                                  `tfsdk:"id"`
	Data       fwtypes.ListNestedObjectValueOf[BaseAsnModel] `tfsdk:"data"`
	Filter     types.Object                                  `tfsdk:"filter"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]        `tfsdk:"pagination"`
	Sort       fwtypes.ObjectValueOf[SortModel]              `tfsdk:"sort"`
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

type BaseAsnModel struct {
	UUID                      types.String                                   `tfsdk:"uuid"`
	Href                      types.String                                   `tfsdk:"href"`
	State                     types.String                                   `tfsdk:"state"`
	Number                    types.Int64                                    `tfsdk:"number"`
	RirRegion                 types.String                                   `tfsdk:"rir_region"`
	Range                     types.String                                   `tfsdk:"range"`
	Project                   fwtypes.ObjectValueOf[ProjectModel]            `tfsdk:"project"`
	Location                  fwtypes.ObjectValueOf[LocationModel]           `tfsdk:"location"`
	OwnershipChallengeDetails fwtypes.ObjectValueOf[OwnershipChallengeModel] `tfsdk:"ownership_challenge_details"`
	Change                    fwtypes.ObjectValueOf[ChangeModel]             `tfsdk:"change"`
	Assets                    fwtypes.ListNestedObjectValueOf[AssetModel]    `tfsdk:"assets"`
	ChangeLog                 fwtypes.ObjectValueOf[ChangeLogModel]          `tfsdk:"change_log"`
}

type ProjectModel struct {
	ProjectID types.String `tfsdk:"project_id"`
}

type LocationModel struct {
	MetroHref types.String `tfsdk:"metro_href"`
	MetroCode types.String `tfsdk:"metro_code"`
	MetroName types.String `tfsdk:"metro_name"`
	Region    types.String `tfsdk:"region"`
	Ibx       types.String `tfsdk:"ibx"`
}

type OwnershipChallengeModel struct {
	Type types.String `tfsdk:"type"`
}

type ChangeModel struct {
	Href types.String `tfsdk:"href"`
}

type AssetModel struct {
	UUID            types.String                                       `tfsdk:"uuid"`
	AnnouncingMetro types.String                                       `tfsdk:"announcing_metro"`
	Service         fwtypes.ObjectValueOf[AssetServiceModel]           `tfsdk:"service"`
	IpBlocks        fwtypes.ListNestedObjectValueOf[AssetIpBlockModel] `tfsdk:"ip_blocks"`
}

type AssetServiceModel struct {
	Href types.String `tfsdk:"href"`
	Type types.String `tfsdk:"type"`
	UUID types.String `tfsdk:"uuid"`
}

type AssetIpBlockModel struct {
	Href types.String `tfsdk:"href"`
	Type types.String `tfsdk:"type"`
	UUID types.String `tfsdk:"uuid"`
}

type ChangeLogModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	UpdatedDateTime types.String `tfsdk:"updated_date_time"`
	DeletedDateTime types.String `tfsdk:"deleted_date_time"`
}

func (m *DataSourceByIDModel) parse(ctx context.Context, asn *fabricv4.Asn) diag.Diagnostics {
	m.AsnID = types.StringValue(asn.GetUuid())
	m.ID = types.StringValue(asn.GetUuid())
	return m.BaseAsnModel.parse(ctx, asn)
}

func (m *DataSourceAllAsnsModel) parse(ctx context.Context, resp *fabricv4.AsnSearchResponseBody) diag.Diagnostics {
	var diags diag.Diagnostics

	asns := resp.GetData()
	if len(asns) < 1 {
		diags.AddError("no data retrieved by asns data source",
			"either the account does not have any ASNs or the combination of filter, limit and offset needs to be updated")
		return diags
	}

	data := make([]*BaseAsnModel, len(asns))
	for i := range asns {
		var model BaseAsnModel
		diags = model.parse(ctx, &asns[i])
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
	m.Data = fwtypes.NewListNestedObjectValueOfSlice[BaseAsnModel](ctx, data)

	return diags
}

func (m *BaseAsnModel) parse(ctx context.Context, asn *fabricv4.Asn) diag.Diagnostics {
	var diags diag.Diagnostics

	m.UUID = types.StringValue(asn.GetUuid())
	m.Href = types.StringValue(asn.GetHref())
	m.State = types.StringValue(string(asn.GetState()))
	m.RirRegion = types.StringValue(string(asn.GetRirRegion()))
	m.Range = types.StringValue(string(asn.GetRange()))

	m.Number = types.Int64Null()
	if v, ok := asn.GetNumberOk(); ok && v != nil {
		m.Number = types.Int64Value(*v)
	}

	project := asn.GetProject()
	m.Project = fwtypes.NewObjectValueOf[ProjectModel](ctx, &ProjectModel{ProjectID: types.StringValue(project.GetProjectId())})

	location := asn.GetLocation()
	m.Location = fwtypes.NewObjectValueOf[LocationModel](ctx, &LocationModel{
		MetroHref: types.StringValue(location.GetMetroHref()),
		MetroCode: types.StringValue(location.GetMetroCode()),
		MetroName: types.StringValue(location.GetMetroName()),
		Region:    types.StringValue(location.GetRegion()),
		Ibx:       types.StringValue(location.GetIbx()),
	})

	challenge := asn.GetOwnershipChallengeDetails()
	m.OwnershipChallengeDetails = fwtypes.NewObjectValueOf[OwnershipChallengeModel](ctx, &OwnershipChallengeModel{
		Type: types.StringValue(string(challenge.GetType())),
	})

	change := asn.GetChange()
	m.Change = fwtypes.NewObjectValueOf[ChangeModel](ctx, &ChangeModel{Href: types.StringValue(change.GetHref())})

	assets := asn.GetAssets()
	assetModels := make([]*AssetModel, len(assets))
	for i, asset := range assets {
		assetModels[i] = assetToModel(ctx, &asset)
	}
	m.Assets = fwtypes.NewListNestedObjectValueOfSlice[AssetModel](ctx, assetModels)

	cl := asn.GetChangeLog()
	m.ChangeLog = fwtypes.NewObjectValueOf[ChangeLogModel](ctx, &ChangeLogModel{
		CreatedDateTime: optionalTime(cl.CreatedDateTime),
		UpdatedDateTime: optionalTime(cl.UpdatedDateTime),
		DeletedDateTime: optionalTime(cl.DeletedDateTime),
	})

	return diags
}

func assetToModel(ctx context.Context, asset *fabricv4.AsnAsset) *AssetModel {
	service := asset.GetService()
	ipBlocks := asset.GetIpBlocks()
	ipBlockModels := make([]*AssetIpBlockModel, len(ipBlocks))
	for i, b := range ipBlocks {
		ipBlockModels[i] = &AssetIpBlockModel{
			Href: types.StringValue(b.GetHref()),
			Type: types.StringValue(string(b.GetType())),
			UUID: types.StringValue(b.GetUuid()),
		}
	}
	return &AssetModel{
		UUID:            types.StringValue(asset.GetUuid()),
		AnnouncingMetro: types.StringValue(asset.GetAnnouncingMetro()),
		Service: fwtypes.NewObjectValueOf[AssetServiceModel](ctx, &AssetServiceModel{
			Href: types.StringValue(service.GetHref()),
			Type: types.StringValue(string(service.GetType())),
			UUID: types.StringValue(service.GetUuid()),
		}),
		IpBlocks: fwtypes.NewListNestedObjectValueOfSlice[AssetIpBlockModel](ctx, ipBlockModels),
	}
}

func optionalTime(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.Format(timeFormat))
}

const (
	ownershipMessageKey   = "message"
	ownershipSignatureKey = "signature"
	typeKey               = "type"
)

type ResourceModel struct {
	ID                        types.String                                  `tfsdk:"id"`
	Timeouts                  timeouts.Value                                `tfsdk:"timeouts"`
	UUID                      types.String                                  `tfsdk:"uuid"`
	Href                      types.String                                  `tfsdk:"href"`
	Type                      types.String                                  `tfsdk:"type"`
	State                     types.String                                  `tfsdk:"state"`
	Number                    types.Int64                                   `tfsdk:"number"`
	RirRegion                 types.String                                  `tfsdk:"rir_region"`
	Range                     types.String                                  `tfsdk:"range"`
	Project                   fwtypes.ObjectValueOf[ProjectModel]           `tfsdk:"project"`
	Location                  fwtypes.ObjectValueOf[ResourceLocationModel]  `tfsdk:"location"`
	OwnershipChallengeDetails fwtypes.ObjectValueOf[ResourceChallengeModel] `tfsdk:"ownership_challenge_details"`
	Change                    fwtypes.ObjectValueOf[ChangeModel]            `tfsdk:"change"`
	Assets                    fwtypes.ListNestedObjectValueOf[AssetModel]   `tfsdk:"assets"`
	ChangeLog                 fwtypes.ObjectValueOf[ChangeLogModel]         `tfsdk:"change_log"`
}

type ResourceLocationModel struct {
	MetroCode types.String `tfsdk:"metro_code"`
}

type ResourceChallengeModel struct {
	Type      types.String `tfsdk:"type"`
	Message   types.String `tfsdk:"message"`
	Signature types.String `tfsdk:"signature"`
}

func (m *ResourceModel) parse(ctx context.Context, asn *fabricv4.Asn) diag.Diagnostics {
	var diags diag.Diagnostics

	var base BaseAsnModel
	diags.Append(base.parse(ctx, asn)...)
	if diags.HasError() {
		return diags
	}

	m.ID = base.UUID
	m.UUID = base.UUID
	m.Href = base.Href
	m.State = base.State
	m.Number = base.Number
	m.Project = base.Project
	m.Change = base.Change
	m.Assets = base.Assets
	m.ChangeLog = base.ChangeLog

	m.RirRegion = types.StringNull()
	if v, ok := asn.GetRirRegionOk(); ok && v != nil {
		m.RirRegion = types.StringValue(string(*v))
	}
	m.Range = types.StringNull()
	if v, ok := asn.GetRangeOk(); ok && v != nil {
		m.Range = types.StringValue(string(*v))
	}

	if t, ok := asn.AdditionalProperties[typeKey].(string); ok && t != "" {
		m.Type = types.StringValue(t)
	}

	m.Location = fwtypes.NewObjectValueOfNull[ResourceLocationModel](ctx)
	if loc, ok := asn.GetLocationOk(); ok && loc != nil && loc.GetMetroCode() != "" {
		m.Location = fwtypes.NewObjectValueOf[ResourceLocationModel](ctx, &ResourceLocationModel{
			MetroCode: types.StringValue(loc.GetMetroCode()),
		})
	}

	previous := ResourceChallengeModel{Message: types.StringNull(), Signature: types.StringNull()}
	if prev, d := m.OwnershipChallengeDetails.ToPtr(ctx); !d.HasError() && prev != nil {
		previous = *prev
	}
	m.OwnershipChallengeDetails = fwtypes.NewObjectValueOfNull[ResourceChallengeModel](ctx)
	if challenge, ok := asn.GetOwnershipChallengeDetailsOk(); ok && challenge != nil {
		model := ResourceChallengeModel{
			Type:      types.StringValue(string(challenge.GetType())),
			Message:   previous.Message,
			Signature: previous.Signature,
		}
		if v, ok := challenge.AdditionalProperties[ownershipMessageKey].(string); ok {
			model.Message = types.StringValue(v)
		}
		if v, ok := challenge.AdditionalProperties[ownershipSignatureKey].(string); ok && v != "" {
			model.Signature = types.StringValue(v)
		}
		m.OwnershipChallengeDetails = fwtypes.NewObjectValueOf[ResourceChallengeModel](ctx, &model)
	}

	return diags
}
