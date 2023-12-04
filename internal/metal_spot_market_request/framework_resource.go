package metal_spot_market_request

import (
	"context"
	"fmt"
    "regexp"
    "time"

	"github.com/equinix/terraform-provider-equinix/internal/helper"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/packethost/packngo"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var (
	matchIPXEScript = regexp.MustCompile(`(?i)^#![i]?pxe`)
)

func NewResource(ctx context.Context) resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_metal_device_network_type",
				Schema: metalSpotMarketRequestResourceSchema(ctx),
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan MetalSpotMarketRequestResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Prepare the data for API request
	params := packngo.SpotMarketRequestInstanceParameters{
		Hostname:        plan.InstanceParams.Hostname.ValueString(),
		BillingCycle:    plan.InstanceParams.BillingCycle.ValueString(),
		Plan:            plan.InstanceParams.Plan.ValueString(),
		OperatingSystem: plan.InstanceParams.OperatingSystem.ValueString(),
	}

    if !plan.InstanceParams.IPXEScriptURL.IsNull() {
        params.IPXEScriptURL = plan.InstanceParams.IPXEScriptURL.ValueString()
    }
    if !plan.InstanceParams.Userdata.IsNull() {
        params.UserData = plan.InstanceParams.Userdata.ValueString()
    }
    if params.OperatingSystem == "custom_ipxe" {
        if params.IPXEScriptURL == "" && params.UserData == "" {
            resp.Diagnostics.AddError(
                "Error creating Metal Spot Market Request",
                "\"ipxe_script_url\" or \"user_data\" must be provided when \"custom_ipxe\" OS is selected.",
            )
            return
		}

        // ipxe_script_url + user_data is OK, unless user_data is an ipxe script in
		// which case it's an error.
		if params.IPXEScriptURL != "" {
			if matchIPXEScript.MatchString(params.UserData) {
                resp.Diagnostics.AddError(
                    "Error creating Metal Spot Market Request",
                    "\"user_data\" should not be an iPXE script when \"ipxe_script_url\" is also provided.",
                )
                return
			}
		}
    }
    if params.OperatingSystem != "custom_ipxe" && params.IPXEScriptURL != "" {
        resp.Diagnostics.AddError(
            "Error creating Metal Spot Market Request",
            "\"ipxe_script_url\" argument provided, but OS is not \"custom_ipxe\". Please verify and fix device arguments.",
        )
        return
	}
    if !plan.InstanceParams.Customdata.IsNull() {
        params.CustomData = plan.InstanceParams.Customdata.ValueString()
    }
    if !plan.InstanceParams.AlwaysPXE.IsNull() {
        params.AlwaysPXE = plan.InstanceParams.AlwaysPXE.ValueBool()
    }
    if !plan.InstanceParams.Description.IsNull() {
        params.Description = plan.InstanceParams.Description.ValueString()
    }
    if !plan.InstanceParams.Locked.IsNull() {
        params.Locked = plan.InstanceParams.Locked.ValueBool()
    }
    if !plan.InstanceParams.Features.IsNull() {
        resp.Diagnostics.Append(plan.InstanceParams.Features.ElementsAs(ctx, &params.Features, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
    }
    if !plan.InstanceParams.Tags.IsNull() {
        resp.Diagnostics.Append(plan.InstanceParams.Tags.ElementsAs(ctx, &params.Tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
    }
    if !plan.InstanceParams.ProjectSSHKeys.IsNull() {
        resp.Diagnostics.Append(plan.InstanceParams.ProjectSSHKeys.ElementsAs(ctx, &params.ProjectSSHKeys, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
    }
    if !plan.InstanceParams.UserSSHKeys.IsNull() {
        resp.Diagnostics.Append(plan.InstanceParams.UserSSHKeys.ElementsAs(ctx, &params.UserSSHKeys, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
    }

    smrc := &packngo.SpotMarketRequestCreateRequest{
		DevicesMax:  int(plan.DevicesMax.ValueInt64()),
		DevicesMin:  int(plan.DevicesMin.ValueInt64()),
		MaxBidPrice: plan.MaxBidPrice.ValueFloat64(),
		Parameters:  params,
	}
    if !plan.Facilities.IsNull() {
        resp.Diagnostics.Append(plan.Facilities.ElementsAs(ctx, &smrc.FacilityIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
    }
    if !plan.Metro.IsNull() {
        smrc.Metro = plan.Metro.ValueString()
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Making an API call to configure the resource
    waitForDevices := plan.WaitForDevices.ValueBool()
    start := time.Now()
	smr, _, err := client.SpotMarketRequests.Create(smrc, plan.ProjectID.String())
	if err != nil {
        err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error creating Metal Spot Market Request",
            "Could not create Spot Market Request: " + err.Error(),        )
        return
	}

    if waitForDevices {
        createTimeout, diags := plan.Timeouts.Create(ctx, 30*time.Minute)
        if diags.HasError() {
            resp.Diagnostics.Append(diags...)
            return
        }
        createTimeout = createTimeout - time.Since(start) - time.Second*10 // reduce 30s to avoid context deadline
		stateConf := &retry.StateChangeConf{
			Pending:        []string{"not_done"},
			Target:         []string{"done"},
			Refresh:        resourceStateRefreshFunc(client, smr.ID),
			Timeout:        createTimeout,
			MinTimeout:     5 * time.Second,
			Delay:          3 * time.Second, // Wait 10 secs before starting
			NotFoundChecks: 600,             // Setting high number, to support long timeouts
		}

		smrRespItf, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
            err = helper.FriendlyError(err)
			resp.Diagnostics.AddError(
                "Error waiting for creation of Metal Spot Market Request",
                fmt.Sprintf("error waiting for Spot Market Request (%s) to become 'done': %s", smr.ID, err),
            )
            return
		}

        var ok bool
        smr, ok = smrRespItf.(*packngo.SpotMarketRequest)
        if !ok {
            resp.Diagnostics.AddError(
                "Error parsing IP Reservation response",
                "Unexpected response type from API",
            )
            return
        }
	}

    // Map the created resource data back to the Terraform state
    stateDiags := plan.parse(ctx, smr)
    resp.Diagnostics.Append(stateDiags...)
    if stateDiags.HasError() {
        return
    }
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state MetalSpotMarketRequestResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the resource from the state
    id := state.ID.ValueString()

    // Retrieve the resource from the API
	smr, _, err := client.SpotMarketRequests.Get(id, &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
	if err != nil {
		err = helper.FriendlyError(err)

         // Check if the Device no longer exists
         if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Metal Spot Market Request",
				fmt.Sprintf("[WARN] Spot Market Request (%s) not found, removing from state", id),
			)
            resp.State.RemoveResource(ctx)
            return
        }

		resp.Diagnostics.AddError(
            "Error reading Device",
            "Could not read Device with ID " + id + ": " + err.Error(),
        )
        return
	}

    // Parse the API response into the Terraform state
    parseDiags := state.parse(ctx, smr)
    resp.Diagnostics.Append(parseDiags...)
    if parseDiags.HasError() {
        return
    }

    // Update the Terraform state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // This resource does not support update
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state MetalSpotMarketRequestResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the organization from the state
    id := state.ID.ValueString()

    waitForDevices := state.WaitForDevices.ValueBool()
    if waitForDevices {
		smr, _, err := client.SpotMarketRequests.Get(id, &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
		if err != nil {
            resp.Diagnostics.AddWarning(
				"Metal Spot Market Request",
				fmt.Sprintf("[WARN] Spot Market Request (%s) not accessible for deletion, removing from state", id),
			)
			return
		}

        deleteTimeout, diags := state.Timeouts.Delete(ctx, 30*time.Minute)
        if diags.HasError() {
            resp.Diagnostics.Append(diags...)
            return
        }
        deleteTimeout = deleteTimeout - time.Second*30 // reduce 30s to avoid context deadline
		stateConf := &retry.StateChangeConf{
			Pending:        []string{"not_done"},
			Target:         []string{"done"},
			Refresh:        resourceStateRefreshFunc(client, id),
			Timeout:        deleteTimeout,
			MinTimeout:     5 * time.Second,
			Delay:          3 * time.Second, // Wait 10 secs before starting
			NotFoundChecks: 600,             // Setting high number, to support long timeouts
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
            err = helper.FriendlyError(err)
            resp.Diagnostics.AddError(
                "Failed to delete Metal Spot Market Request",
                fmt.Sprintf("error waiting for Spot Market Request (%s) to become 'done' before proceed with deletion: %s", id, err),
            )
			return
		}

		for _, d := range smr.Devices {
			deleteResp, err := client.Devices.Delete(d.ID, true)
            if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
                err = helper.FriendlyError(err)
                resp.Diagnostics.AddError(
                    fmt.Sprintf("Failed to delete Metal Spot Market Request %s", id),
                    fmt.Sprintf("error waiting for Spot Market Request (%s) to be deleted: %s", id, err),
                )
                return
            }
		}
	}

}

func resourceStateRefreshFunc(client *packngo.Client, requestId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		smr, _, err := client.SpotMarketRequests.Get(requestId, &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
		if err != nil {
			return nil, "", fmt.Errorf("failed to fetch Spot market request with following error: %s", err.Error())
		}
		var finished bool

		for _, d := range smr.Devices {

			dev, _, err := client.Devices.Get(d.ID, nil)
			if err != nil {
				return nil, "", fmt.Errorf("failed to fetch Device with following error: %s", err.Error())
			}
			if dev.State != "active" {
				break
			} else {
				finished = true
			}
		}
		if finished {
			return smr, "done", nil
		}
		return nil, "not_done", nil
	}
}