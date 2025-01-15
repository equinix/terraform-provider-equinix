package metros

import (
	"context"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	fwtypes "github.com/equinix/terraform-provider-equinix/internal/framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type MetroModel struct {
	ID                  types.String                                         `tfsdk:"id"`
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
	Offset   types.Int64  `tfsdk:"offset"`
	Limit    types.Int64  `tfsdk:"limit"`
	Total    types.Int64  `tfsdk:"total"`
	Next     types.String `tfsdk:"next"`
	Previous types.String `tfsdk:"previous"`
}

type AllMetrosModel struct {
	ID         types.String                                `tfsdk:"id"`
	Presence   types.String                                `tfsdk:"presence"`
	Pagination fwtypes.ObjectValueOf[PaginationModel]      `tfsdk:"pagination"`
	Data       fwtypes.ListNestedObjectValueOf[MetroModel] `tfsdk:"metro_model"`
}

func (m *MetroModel) parseDataSourceByMetroCode(ctx context.Context, metro *fabricv4.Metro) diag.Diagnostics {
	diags := parseMetros(ctx, metro,
		&m.ID, &m.Type, &m.Href, &m.Code, &m.Region, &m.Name,
		&m.EquinixASN, &m.LocalVCBandwidthMax, &m.GeoCoordinates, &m.ConnectedMetros, &m.GeoScopes)

	return diags
}

func parseMetros(ctx context.Context, metro *fabricv4.Metro, id, tp, href, code, region, name *basetypes.StringValue, equinixAsn, localBandwidthMax *basetypes.Int64Value, geoCoordinates *fwtypes.ObjectValueOf[GeoCoordinatesModel], connectedMetros *fwtypes.ListNestedObjectValueOf[ConnectedMetroModel], gScopes *fwtypes.ListValueOf[types.String]) diag.Diagnostics {

	var diags diag.Diagnostics
	*href = types.StringValue(metro.GetHref())
	*tp = types.StringValue(metro.GetType())
	*code = types.StringValue(metro.GetCode())
	*region = types.StringValue(metro.GetRegion())
	if metro.GetName() != "" {
		*name = types.StringValue(metro.GetName())
	}

	if equinixAsn != nil {
		*equinixAsn = types.Int64Value(metro.GetEquinixAsn())
	}

	if localBandwidthMax != nil {
		*localBandwidthMax = types.Int64Value(metro.GetLocalVCBandwidthMax())
	}

	geoCoord, diags := parseGeoCoordinates(ctx, metro.GetGeoCoordinates())

	if diags.HasError() {
		return diags
	}

	*geoCoordinates = geoCoord

	connMetros, diags := parseconnectedMetros(ctx, metro.GetConnectedMetros())
	if diags.HasError() {
		return diags
	}
	*connectedMetros = connMetros
	geoScopes, diags := parseGeoScopes(ctx, metro.GetGeoScopes())
	if diags.HasError() {
		return diags
	}
	*gScopes = geoScopes
	return diags
}

func parseconnectedMetros(ctx context.Context, connectedMetros []fabricv4.ConnectedMetro) (fwtypes.ListNestedObjectValueOf[ConnectedMetroModel], diag.Diagnostics) {
	connMetros := make([]ConnectedMetroModel, len(connectedMetros))
	for i, metro := range connectedMetros {
		connMetros[i] = ConnectedMetroModel{
			Href:                 types.StringValue(*metro.Href),
			Code:                 types.StringValue(*metro.Code),
			AvgLatency:           types.Float32Value(*metro.AvgLatency),
			RemoteVCBandwidthMax: types.Int64Value(*metro.RemoteVCBandwidthMax),
		}
	}
	return fwtypes.NewListNestedObjectValueOfValueSlice(ctx, connMetros), nil
}

func parseGeoCoordinates(ctx context.Context, coordinates fabricv4.GeoCoordinates) (fwtypes.ObjectValueOf[GeoCoordinatesModel], diag.Diagnostics) {
	var diags diag.Diagnostics
	//if coordinates == nil {
	//	diags.AddError("Invalid Input", "Coordinates should not be nil")
	//	return nil, diags
	//}
	//
	//if coordinates.Latitude == nil || coordinates.Longitude == nil {
	//	diags.AddError("Invalid Input", "Latitude and Longitude should not be nil")
	//	return nil, diags
	//}

	result := GeoCoordinatesModel{
		Latitude:  types.Float64Value(*coordinates.Latitude),
		Longitude: types.Float64Value(*coordinates.Longitude),
	}

	return fwtypes.NewObjectValueOf[GeoCoordinatesModel](ctx, &result), diags
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
