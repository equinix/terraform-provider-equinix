package services

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NetworkEdgeResources() []func() resource.Resource {
	return []func() resource.Resource{}
}

func NetworkEdgeDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
