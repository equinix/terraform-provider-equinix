package services

import (
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/stream"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{
		stream.NewResource,
	}
}

func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		stream.NewDataSourceByStreamID,
		stream.NewDataSourceAllStreams,
	}
}
