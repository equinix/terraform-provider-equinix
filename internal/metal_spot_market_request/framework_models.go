package metal_spot_market_request

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type MetalSpotMarketRequestResourceModel struct {
    ID             types.String        `tfsdk:"id"`
    DevicesMin     types.Int64         `tfsdk:"devices_min"`
    DevicesMax     types.Int64         `tfsdk:"devices_max"`
    MaxBidPrice    types.Float64       `tfsdk:"max_bid_price"`
    Facilities     types.List          `tfsdk:"facilities"`
    Metro          types.String        `tfsdk:"metro"`
    ProjectID      types.String        `tfsdk:"project_id"`
    WaitForDevices types.Bool          `tfsdk:"wait_for_devices"`
    InstanceParams *InstanceParameters `tfsdk:"instance_parameters"`
    Timeouts       timeouts.Value      `tfsdk:"timeouts"`
}

type InstanceParameters struct {
    BillingCycle     types.String `tfsdk:"billing_cycle"`
    Plan             types.String `tfsdk:"plan"`
    OperatingSystem  types.String `tfsdk:"operating_system"`
    Hostname         types.String `tfsdk:"hostname"`
    TermintationTime types.String `tfsdk:"termination_time"`
    TerminationTime  types.String `tfsdk:"termination_time"`
    AlwaysPXE        types.Bool   `tfsdk:"always_pxe"`
    Description      types.String `tfsdk:"description"`
    Features         types.List   `tfsdk:"features"`
    Locked           types.Bool   `tfsdk:"locked"`
    ProjectSSHKeys   types.List   `tfsdk:"project_ssh_keys"`
    UserSSHKeys      types.List   `tfsdk:"user_ssh_keys"`
    Userdata         types.String `tfsdk:"userdata"`
    Customdata       types.String `tfsdk:"customdata"`
    IPXEScriptURL    types.String `tfsdk:"ipxe_script_url"`
    Tags             types.List   `tfsdk:"tags"`
}

func (rm *MetalSpotMarketRequestResourceModel) parse(ctx context.Context, smr *packngo.SpotMarketRequest) diag.Diagnostics {
    var diags diag.Diagnostics

    // Map fields from packngo.SpotMarketRequest to MetalSpotMarketRequestResourceModel
    rm.DevicesMin = types.Int64Value(int64(smr.DevicesMin))
    rm.DevicesMax = types.Int64Value(int64(smr.DevicesMax))
    rm.MaxBidPrice = types.Float64Value(smr.MaxBidPrice)

    // Assuming smr.Facilities is a slice of string
    facilities, diags := types.ListValueFrom(ctx, types.StringType, smr.Facilities)
    if diags.HasError() {
        return diags
    }
    rm.Facilities = facilities

    rm.Metro = types.StringValue(smr.Metro.ID)
    rm.ProjectID = types.StringValue(smr.Project.ID)

    // Map instance_parameters
    params := &InstanceParameters{
        BillingCycle:     types.StringValue(smr.Parameters.BillingCycle),
        Plan:             types.StringValue(smr.Parameters.Plan),
        OperatingSystem:  types.StringValue(smr.Parameters.OperatingSystem),
        Hostname:         types.StringValue(smr.Parameters.Hostname),
        TermintationTime: types.StringValue(smr.Parameters.TerminationTime.String()),
        TerminationTime:  types.StringValue(smr.Parameters.TerminationTime.String()),
        AlwaysPXE:        types.BoolValue(smr.Parameters.AlwaysPXE),
        Description:      types.StringValue(smr.Parameters.Description),
        Locked:           types.BoolValue(smr.Parameters.Locked),
        Userdata:         types.StringValue(smr.Parameters.UserData),
        Customdata:       types.StringValue(smr.Parameters.CustomData),
        IPXEScriptURL:    types.StringValue(smr.Parameters.IPXEScriptURL),
    }
    rm.InstanceParams = params

    // Handling project ssh keys as a list
    projectKeys, diags := types.ListValueFrom(ctx, types.StringType, smr.Parameters.ProjectSSHKeys)
    if diags.HasError() {
        return diags
    }
    params.ProjectSSHKeys = projectKeys

    // Handling user ssh keys as a list
    userKeys, diags := types.ListValueFrom(ctx, types.StringType, smr.Parameters.UserSSHKeys)
    if diags.HasError() {
        return diags
    }
    params.UserSSHKeys = userKeys

    // Handling features as a list
    features, diags := types.ListValueFrom(ctx, types.StringType, smr.Parameters.Features)
    if diags.HasError() {
        return diags
    }
    params.Features = features

    // Handling tags as a list
    tags, diags := types.ListValueFrom(ctx, types.StringType, smr.Parameters.Tags)
    if diags.HasError() {
        return diags
    }
    params.Tags = tags
 

    return diags
}
