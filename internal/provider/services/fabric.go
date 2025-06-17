// Package services for Fabric resources and data sources
package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connectionrouteaggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/metro"
	precisiontime "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/precision_time"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/routeaggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/routeaggregationrule"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"
	streamattachment "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream_attachment"
	streamsubscription "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream_subscription"
	streamalertrule "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/streamalertrule"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// FabricResources represents fabric resources
func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{
		connectionrouteaggregation.NewResource,
		precisiontime.NewResource,
		routeaggregation.NewResource,
		routeaggregationrule.NewResource,
		stream.NewResource,
		streamattachment.NewResource,
		streamsubscription.NewResource,
		streamalertrule.NewResource,
	}
}

// FabricDatasources represents fabric data source
func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		connectionrouteaggregation.NewDataSourceByConnectionRouteAggregationID,
		connectionrouteaggregation.NewDataSourceAllConnectionRouteAggregations,
		metro.NewDataSourceMetroCode,
		metro.NewDataSourceMetros,
		precisiontime.NewDataSourceByEptServiceID,
		precisiontime.NewDataSourceAllEptServices,
		routeaggregation.NewDataSourceByRouteAggregationID,
		routeaggregation.NewDataSourceAllRouteAggregation,
		routeaggregationrule.NewDataSourceByRouteAggregationRuleID,
		routeaggregationrule.NewDataSourceAllRouteAggregationRule,
		stream.NewDataSourceByStreamID,
		stream.NewDataSourceAllStreams,
		streamattachment.NewDataSourceAllStreamAttachments,
		streamattachment.NewDataSourceByIDs,
		streamsubscription.NewDataSourceAllStreamSubscriptions,
		streamsubscription.NewDataSourceByIDs,
		streamalertrule.NewDataSourceAllStreamAlertRules,
		streamalertrule.NewDataSourceByStreamAlertRuleIDs,
	}
}
