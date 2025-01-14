package metros

import (
	"context"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_metros",
			},
		),
	}
}

type DataSource struct {
	framework.BaseDataSource
}

func (r *DataSource) AllMetrosSchema(ctx context.Context, response *datasource.SchemaResponse) {
	response.Schema = dataSourceAllMetroSchema(ctx)
}

func (r *DataSource) MetroCodeSchema(ctx context.Context, response *datasource.SchemaResponse) {
	response.Schema = dataSourceSingleMetroSchema(ctx)
}

// READ function for GET Metro Code data source
func (r *DataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	// Retrieve values from plan
	var data MetroModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Extract the ID of the resource from the config
	//id := data.ID.ValueString()

	metroCode := data.Code.ValueString()

	// Use API client to get the current state of the resource
	metroByCode, _, err := client.MetrosApi.GetMetroByCode(ctx, metroCode).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		diag.FromErr(equinix_errors.FormatFabricError(err))
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(data.parseDataSourceByMetroCode(ctx, metroByCode)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

//READ function for the Get Metros data source also, use the Get Metros API

//func (r *DataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
//	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)
//
//	// Retrieve values from plan
//	var data AllMetrosModel
//	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//
//	// Extract the ID of the resource from the config
//	id := data.ID.ValueString()
//
//	// Use API client to get the current state of the resource
//	metros, _, err := client.MetrosApi.GetMetros(ctx).Execute()
//
//	if err != nil {
//		// If Metros is not found, remove it from the state
//		if equinix_errors.IsNotFound(err) {
//			response.Diagnostics.AddWarning(
//				"Get All Metros",
//				fmt.Sprintf("[WARN] Metros (%s) not found, removing from state", id),
//			)
//			response.State.RemoveResource(ctx)
//			return
//		}
//		response.Diagnostics.AddError(
//			"Error reading - Get All Metros",
//			"Could not read Metros with ID:"+id+": "+err.Error(),
//		)
//		return
//	}
//
//	// Set state to fully populated data
//	response.Diagnostics.Append(data.parse(ctx, metros)...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//
//	// Update the Terraform state
//	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
//}
