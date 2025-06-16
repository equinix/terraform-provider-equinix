package stream_alert_rule

import (
	"context"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// NewDataSourceByStreamAlertRuleIDs creates a new data source for stream alert rule by IDs
func NewDataSourceByStreamAlertRuleIDs() datasource.DataSource {
	return &DataSourceByStreamAlertRuleID{
		BaseDataSource: framework.NewBaseDataSource(
			framework.BaseDataSourceConfig{
				Name: "equinix_fabric_stream_alert_rule",
			},
		),
	}
}

// DataSourceByStreamAlertRuleID datasource represents stream alert rule by IDs
type DataSourceByStreamAlertRuleID struct {
	framework.BaseDataSource
}

// Schema returns the datasource schema
func (r *DataSourceByStreamAlertRuleID) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = dataSourceStreamAlertRuleByID(ctx)
}

func (r *DataSourceByStreamAlertRuleID) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, request.ProviderMeta)

	// Retrieve values from plan
	var data dataSourceStreamAlertRuleByIDsModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Use API client to get the current state of the resource
	streamSubscription, _, err := client.StreamAlertRulesApi.GetStreamAlertRuleByUuid(ctx, data.StreamID.ValueString(), data.AlertRuleID.ValueString()).Execute()

	if err != nil {
		response.State.RemoveResource(ctx)
		response.Diagnostics.AddError("api error retrieving stream subscription data", equinix_errors.FormatFabricError(err).Error())
		return
	}

	// Set state to fully populated data
	response.Diagnostics.Append(data.parse(ctx, streamSubscription)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Update the Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
