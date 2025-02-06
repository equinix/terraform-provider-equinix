package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/metro"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"
	streamsubscription "github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream_subscription"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{
		stream.NewResource,
		streamsubscription.NewResource,
	}
}

func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		metro.NewDataSourceMetroCode,
		metro.NewDataSourceMetros,
		stream.NewDataSourceAllStreams,
		stream.NewDataSourceByStreamID,
		streamsubscription.NewDataSourceAllStreamSubscriptions,
		streamsubscription.NewDataSourceByIDs,
	}
}
