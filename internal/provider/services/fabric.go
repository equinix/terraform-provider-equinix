package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/metro"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/routeaggregation"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"
	streamattachment "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream_attachment"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{
		routeaggregation.NewResource,
		stream.NewResource,
		streamattachment.NewResource,
	}
}

func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		metro.NewDataSourceMetroCode,
		metro.NewDataSourceMetros,
		routeaggregation.NewDataSourceByRouteAggregationID,
		routeaggregation.NewDataSourceAllRouteAggregation,
		stream.NewDataSourceByStreamID,
		stream.NewDataSourceAllStreams,
		streamattachment.NewDataSourceAllStreamAttachments,
		streamattachment.NewDataSourceByIDs,
	}
}
