package internetaccess

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/equinix/terraform-provider-equinix/internal/framework"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const deletedMarker = "tf-marker-for-deleted-internet-access-service"

// polling intervals; variables so that unit tests can shorten them
var (
	waitDelay       = 10 * time.Second
	waitMinInterval = 10 * time.Second
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: framework.NewBaseResource(
			framework.BaseResourceConfig{
				Name: "equinix_fabric_internet_access_service",
			},
		),
	}
}

type Resource struct {
	framework.BaseResource
}

func (r Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema(ctx)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest, err := buildCreateRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("failed building Internet Access service create request", err.Error())
		return
	}

	svc, _, apiErr := client.InternetAccessServicesApi.CreateEiaService(ctx).InternetAccessPostRequest(createRequest).Execute()
	if apiErr != nil {
		resp.Diagnostics.AddError("api error creating Internet Access service", equinix_errors.FormatFabricError(apiErr).Error())
		return
	}

	// The API acknowledges the request with 201 and provisions asynchronously (the service can stay
	// in PROVISIONING for a long time or end up FAILED), so we do not wait for a particular state.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), svc.GetUuid())...)

	resp.Diagnostics.Append(plan.parse(ctx, svc)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	svc, httpResp, err := getService(ctx, client, state.ID.ValueString())
	if err != nil {
		if httpResp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, httpResp.StatusCode) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed retrieving Internet Access service %s", state.ID.ValueString()),
			equinix_errors.FormatFabricError(err).Error(),
		)
		return
	}

	// a deprovisioned service no longer exists from the user's point of view
	if svc.GetState() == fabricv4.INTERNETACCESSSERVICESTATE_DEPROVISIONED {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(state.parse(ctx, svc)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is not implemented yet: only name, bandwidth, bandwidth_commit, customer_routes and
// order.purchase_order_number can change in place and they will be sent as a PATCH.
func (r *Resource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"in-place update of Internet Access service is not supported yet",
		"Recreate the resource (terraform apply -replace) to change it.",
	)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client := r.Meta.NewFabricClientForFramework(ctx, req.ProviderMeta)

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, deleteResp, err := client.InternetAccessServicesApi.DeleteEiaService(ctx, id).Execute()
	if err != nil {
		if deleteResp == nil || !slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed deleting Internet Access service %s", id),
				equinix_errors.FormatFabricError(err).Error(),
			)
			return
		}
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 30*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, waitErr := getDeleteWaiter(ctx, client, id, deleteTimeout).WaitForStateContext(ctx); waitErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed waiting for Internet Access service %s to be deleted", id), waitErr.Error())
	}
}

func getDeleteWaiter(ctx context.Context, client *fabricv4.APIClient, uuid string, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.INTERNETACCESSSERVICESTATE_PROVISIONED),
			string(fabricv4.INTERNETACCESSSERVICESTATE_PROVISIONING),
			string(fabricv4.INTERNETACCESSSERVICESTATE_DEPROVISIONING),
		},
		// FAILED is final: the API keeps failed services in that state even after a DELETE
		Target: []string{deletedMarker, string(fabricv4.INTERNETACCESSSERVICESTATE_DEPROVISIONED), string(fabricv4.INTERNETACCESSSERVICESTATE_FAILED)},
		Refresh: func() (any, string, error) {
			svc, httpResp, err := getService(ctx, client, uuid)
			if err != nil {
				if httpResp != nil && slices.Contains([]int{http.StatusForbidden, http.StatusNotFound}, httpResp.StatusCode) {
					return svc, deletedMarker, nil
				}
				return nil, "", err
			}
			return svc, string(svc.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      waitDelay,
		MinTimeout: waitMinInterval,
	}
}
