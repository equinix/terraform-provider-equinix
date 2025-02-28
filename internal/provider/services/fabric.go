package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connectionrouteaggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/metro"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/routeaggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/routeaggregationrule"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"
	streamattachment "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream_attachment"
	streamsubscription "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream_subscription"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{
		connectionrouteaggregation.NewResource,
		routeaggregation.NewResource,
		routeaggregationrule.NewResource,
		stream.NewResource,
		streamattachment.NewResource,
		streamsubscription.NewResource,
	}
}

func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		connectionrouteaggregation.NewDataSourceByConnectionRouteAggregationID,
		connectionrouteaggregation.NewDataSourceAllConnectionRouteAggregations,
		metro.NewDataSourceMetroCode,
		metro.NewDataSourceMetros,
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
	}
}
