package services

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func FabricResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

func FabricDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
