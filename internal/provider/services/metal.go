package services

import (
	metalconnection "github.com/equinix/terraform-provider-equinix/internal/resources/metal/connection"
	metalgateway "github.com/equinix/terraform-provider-equinix/internal/resources/metal/gateway"
	metalorganization "github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization"
	metalorganizationmember "github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization_member"
	metalproject "github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	metalprojectsshkey "github.com/equinix/terraform-provider-equinix/internal/resources/metal/project_ssh_key"
	metalsshkey "github.com/equinix/terraform-provider-equinix/internal/resources/metal/ssh_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func MetalResources() []func() resource.Resource {
	return []func() resource.Resource{
		metalgateway.NewResource,
		metalproject.NewResource,
		metalprojectsshkey.NewResource,
		metalsshkey.NewResource,
		metalconnection.NewResource,
		metalorganization.NewResource,
		metalorganizationmember.NewResource,
		vlan.NewResource,
	}
}

func MetalDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		metalgateway.NewDataSource,
		metalproject.NewDataSource,
		metalprojectsshkey.NewDataSource,
		metalconnection.NewDataSource,
		metalorganization.NewDataSource,
		vlan.NewDataSource,
	}
}
