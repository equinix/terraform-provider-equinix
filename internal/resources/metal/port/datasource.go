package port

import (
	"context"
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type dataSource struct {
	framework.BaseDataSource
}

// NewDataSource returns the TF resource representing device network ports.
func NewDataSource() datasource.DataSource {
	return &dataSource{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_metal_port",
			},
		),
	}
}

// Read implements datasource.DataSource.
func (d *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Meta.NewMetalClientForFramework(ctx, req.ProviderMeta)

	var port *metalv1.Port
	var method string
	var err error
	if !data.PortID.IsNull() {
		method = "port ID"
		port, _, err = client.PortsApi.FindPortById(ctx, data.PortID.ValueString()).Include([]string{"virtual_networks"}).Execute()
	} else if !data.DeviceID.IsNull() && !data.Name.IsNull() {
		method = "Device ID and Port Name"
		device, _, err := client.DevicesApi.FindDeviceById(ctx, data.DeviceID.ValueString()).Include([]string{"virtual_networks"}).Execute()
		if err != nil {

			resp.Diagnostics.AddError(
				"Failed to locate port",
				fmt.Sprintf("Error fetching device [%s] informations: %s", data.DeviceID.ValueString(), err),
			)
			return
		}

		for _, p := range device.NetworkPorts {
			if p.GetName() == data.Name.ValueString() {
				port = &p
				break
			}
		}
	} else {
		resp.Diagnostics.AddError(
			"Invalid datasource specification",
			"Port datasources require either 'port_id' or both 'device_id' and 'name') to be set",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to fetch port information",
			fmt.Sprintf("Failed finding port by %s: %s", method, err),
		)
		return
	}

	resp.Diagnostics.Append(data.parse(ctx, port)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Schema implements datasource.DataSource.
// Subtle: this method shadows the method (BaseDataSource).Schema of DataSource.BaseDataSource.
func (d *dataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema(ctx)
}
