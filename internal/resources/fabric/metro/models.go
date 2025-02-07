package metro

import (
	"context"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Model struct {
	Href                types.String                                         `tfsdk:"href"`
	Type                types.String                                         `tfsdk:"type"`
	Code                types.String                                         `tfsdk:"code"`
	Region              types.String                                         `tfsdk:"region"`
	Name                types.String                                         `tfsdk:"name"`
	EquinixASN          types.Int64                                          `tfsdk:"equinix_asn"`
	LocalVCBandwidthMax types.Int64                                          `tfsdk:"local_vc_bandwidth_max"`
	GeoCoordinates      fwtypes.ObjectValueOf[GeoCoordinatesModel]           `tfsdk:"geo_coordinates"`
	ConnectedMetros     fwtypes.ListNestedObjectValueOf[ConnectedMetroModel] `tfsdk:"connected_metros"`
	GeoScopes           fwtypes.ListValueOf[types.String]                    `tfsdk:"geo_scopes"`
}

type ConnectedMetroModel struct {
	Href                 types.String  `tfsdk:"href"`
	Code                 types.String  `tfsdk:"code"`
	AvgLatency           types.Float32 `tfsdk:"avg_latency"`
	RemoteVCBandwidthMax types.Int64   `tfsdk:"remote_vc_bandwidth_max"`
}

type GeoCoordinatesModel struct {
	Latitude  types.Float64 `tfsdk:"latitude"`
	Longitude types.Float64 `tfsdk:"longitude"`
}

type PaginationModel struct {
	Offset   types.Int32  `tfsdk:"offset"`
	Limit    types.Int32  `tfsdk:"limit"`
	Total    types.Int32  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type DataSourceByCodeModel struct {
	ID        types.String `tfsdk:"id"`
	MetroCode types.String `tfsdk:"metro_code"`
	Model
}

type DataSourceAllMetrosModel struct {
	ID         types.String                           `tfsdk:"id"`
	Presence   types.String                           `tfsdk:"presence"`
	Data       fwtypes.ListNestedObjectValueOf[Model] `tfsdk:"data"`
	Pagination fwtypes.ObjectValueOf[PaginationModel] `tfsdk:"pagination"`
}

func (a *DataSourceAllMetrosModel) parse(ctx context.Context, metroResponse *fabricv4.MetroResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(metroResponse.GetData()) < 1 {
		diags.AddError("no data retrieved by metros data source",
			"either the account does not have any metros data to pull or the combination of limit and offset needs to be updated")
		return diags
	}

	data := make([]Model, len(metroResponse.GetData()))
	metros := metroResponse.GetData()
	for i, metro := range metros {
		var metroModel Model
		diags := metroModel.parse(ctx, &metro)
		if diags.HasError() {
			return diags
		}
		data[i] = metroModel
	}
	responsePagination := metroResponse.GetPagination()
	pagination := PaginationModel{
		Offset:   types.Int32Value(responsePagination.GetOffset()),
		Limit:    types.Int32Value(responsePagination.GetLimit()),
		Total:    types.Int32Value(responsePagination.GetTotal()),
		Next:     types.StringValue(responsePagination.GetNext()),
		Previous: types.StringValue(responsePagination.GetPrevious()),
	}

	a.ID = types.StringValue(data[0].Code.ValueString())
	a.Pagination = fwtypes.NewObjectValueOf[PaginationModel](ctx, &pagination)
	a.Data = fwtypes.NewListNestedObjectValueOfValueSlice[Model](ctx, data)

	return diags
}

func (m *DataSourceByCodeModel) parse(ctx context.Context, metro *fabricv4.Metro) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(metro.GetCode())
	m.MetroCode = types.StringValue(metro.GetCode())

	var metroModel Model

	diags = metroModel.parse(ctx, metro)
	if diags.HasError() {
		return diags
	}

	m.Model = metroModel

	return diags
}

func (m *Model) parse(ctx context.Context, metro *fabricv4.Metro) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Href = types.StringValue(metro.GetHref())
	m.Type = types.StringValue(metro.GetType())
	m.Code = types.StringValue(metro.GetCode())
	m.Region = types.StringValue(metro.GetRegion())
	m.Name = types.StringValue(metro.GetName())
	m.EquinixASN = types.Int64Value(metro.GetEquinixAsn())
	m.LocalVCBandwidthMax = types.Int64Value(metro.GetLocalVCBandwidthMax())

	m.GeoCoordinates, diags = parseGeoCoordinates(ctx, metro.GetGeoCoordinates())
	if diags.HasError() {
		return diags
	}

	m.ConnectedMetros, diags = parseConnectedMetros(ctx, metro.GetConnectedMetros())
	if diags.HasError() {
		return diags
	}

	m.GeoScopes, diags = parseGeoScopes(ctx, metro.GetGeoScopes())
	if diags.HasError() {
		return diags
	}

	return diags
}

func parseGeoScopes(ctx context.Context, scopes []fabricv4.GeoScopeType) (fwtypes.ListValueOf[types.String], diag.Diagnostics) {
	var diags diag.Diagnostics
	geoScopeTypeList := make([]attr.Value, len(scopes))

	for i, scope := range scopes {
		geoScopeTypeList[i] = types.StringValue(string(scope))
	}
	geoScopeValue, diags := fwtypes.NewListValueOf[types.String](ctx, geoScopeTypeList)

	if diags.HasError() {
		return fwtypes.NewListValueOfNull[types.String](ctx), diags
	}
	return geoScopeValue, diags
}

func parseGeoCoordinates(ctx context.Context, coordinates fabricv4.GeoCoordinates) (fwtypes.ObjectValueOf[GeoCoordinatesModel], diag.Diagnostics) {
	diags := diag.Diagnostics{}

	result := GeoCoordinatesModel{}
	if coordinates.Latitude != nil {
		result.Latitude = types.Float64Value(coordinates.GetLatitude())
	}

	if coordinates.Longitude != nil {
		result.Longitude = types.Float64Value(coordinates.GetLongitude())
	}
	return fwtypes.NewObjectValueOf[GeoCoordinatesModel](ctx, &result), diags
}

func parseConnectedMetros(ctx context.Context, connectedMetros []fabricv4.ConnectedMetro) (fwtypes.ListNestedObjectValueOf[ConnectedMetroModel], diag.Diagnostics) {
	connMetros := make([]ConnectedMetroModel, len(connectedMetros))
	for i, metro := range connectedMetros {
		connMetros[i] = ConnectedMetroModel{}
		if metro.Href != nil {
			connMetros[i].Href = types.StringValue(metro.GetHref())
		}
		if metro.Code != nil {
			connMetros[i].Code = types.StringValue(metro.GetCode())
		}
		if metro.AvgLatency != nil {
			connMetros[i].AvgLatency = types.Float32Value(metro.GetAvgLatency())
		}
		if metro.RemoteVCBandwidthMax != nil {
			connMetros[i].RemoteVCBandwidthMax = types.Int64Value(metro.GetRemoteVCBandwidthMax())
		}
	}
	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, connMetros), nil
}
