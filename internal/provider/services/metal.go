package services

import (
	metalconnection "github.com/equinix/terraform-provider-equinix/internal/resources/metal/connection"
	metalgateway "github.com/equinix/terraform-provider-equinix/internal/resources/metal/gateway"
	metalorganization "github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization"
	metalorganizationmember "github.com/equinix/terraform-provider-equinix/internal/resources/metal/organization_member"
	metalport "github.com/equinix/terraform-provider-equinix/internal/resources/metal/port"
	metalproject "github.com/equinix/terraform-provider-equinix/internal/resources/metal/project"
	metalprojectsshkey "github.com/equinix/terraform-provider-equinix/internal/resources/metal/project_ssh_key"
	metalsshkey "github.com/equinix/terraform-provider-equinix/internal/resources/metal/ssh_key"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/vlan"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// MetalResources exposes any resources defined under the terraform-plugin-framework.
func MetalResources() []func() resource.Resource {
	return []func() resource.Resource{
		metalgateway.NewResource,
		metalport.NewResource,
		metalproject.NewResource,
		metalprojectsshkey.NewResource,
		metalsshkey.NewResource,
		metalconnection.NewResource,
		metalorganization.NewResource,
		metalorganizationmember.NewResource,
		vlan.NewResource,
	}
}

// MetalDatasources exposes any datasources defined under the terraform-plugin-framework.
func MetalDatasources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		metalgateway.NewDataSource,
		metalport.NewDataSource,
		metalproject.NewDataSource,
		metalprojectsshkey.NewDataSource,
		metalconnection.NewDataSource,
		metalorganization.NewDataSource,
		vlan.NewDataSource,
	}
}
