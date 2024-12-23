package service_token

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
		},
		ReadContext:   resourceRead,
		CreateContext: resourceCreate,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      resourceSchema(),
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Service Token`,
	}
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	serviceToken, _, err := client.ServiceTokensApi.GetServiceTokenByUuid(ctx, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Service Token %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(serviceToken.GetUuid())
	return setServiceTokenMap(d, serviceToken)
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	createRequest := buildCreateRequest(d)

	start := time.Now()
	serviceToken, _, err := client.ServiceTokensApi.CreateServiceToken(ctx).ServiceToken(createRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(serviceToken.GetUuid())
	notificationsMap := equinix_fabric_schema.NotificationsGoToTerraform(createRequest.GetNotifications())
	if err = d.Set("notifications", notificationsMap); err != nil {
		return diag.Errorf("error setting notifications config to state: %s", err)
	}

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if _, err = waitForStability(d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for service token (%s) to be created: %s", d.Id(), err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbToken, err := waitForStability(d.Id(), meta, d, ctx, updateTimeout)
	if err != nil {
		return diag.Errorf("either timed out or errored out while fetching Fabric Service Token for uuid %s and error %v", d.Id(), err)
	}
	diags := diag.Diagnostics{}
	updates, err := buildUpdateRequest(d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{Severity: 1, Summary: err.Error()})
		return diags
	}
	for _, update := range updates {
		_, _, err = client.ServiceTokensApi.UpdateServiceTokenByUuid(ctx, d.Id()).ServiceTokenChangeOperation(update).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("service token property update request error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", equinix_errors.FormatFabricError(err), update)})
			continue
		}
		updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
		updateServiceToken, err := waitForStability(d.Id(), meta, d, ctx, updateTimeout)
		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("service token property update completion timeout error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", equinix_errors.FormatFabricError(err), update)})
		} else {
			dbToken = updateServiceToken
		}

	}
	d.SetId(dbToken.GetUuid())
	return append(diags, setServiceTokenMap(d, dbToken)...)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(d)

	start := time.Now()
	_, _, err := client.ServiceTokensApi.DeleteServiceTokenByUuid(ctx, d.Id()).Execute()
	if err != nil {
		if genericError, ok := err.(*fabricv4.GenericOpenAPIError); ok {
			if fabricErrs, ok := genericError.Model().([]fabricv4.Error); ok {
				// EQ-3034019 = Service Token already deleted
				if equinix_errors.HasErrorCode(fabricErrs, "EQ-3034019") {
					return diags
				}
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	if err = WaitForDeletion(d.Id(), meta, d, ctx, deleteTimeout); err != nil {
		return diag.Errorf("error waiting for service token (%s) to be deleted: %s", d.Id(), err)
	}
	return diags
}

func waitForStability(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.ServiceToken, error) {
	log.Printf("Waiting for service token to be created, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{
			string(fabricv4.SERVICETOKENSTATE_INACTIVE),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			serviceToken, _, err := client.ServiceTokensApi.GetServiceTokenByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			currentState := string(serviceToken.GetState())
			return serviceToken, currentState, nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var serviceToken *fabricv4.ServiceToken

	if err != nil {
		log.Printf("[ERROR] Error while waiting for service token to go to INACTIVE state: %v", err)
		return nil, err
	}

	serviceToken, ok := inter.(*fabricv4.ServiceToken)
	if !ok {
		return nil, fmt.Errorf("expected *fabricv4.ServiceToken, but got %T", inter)
	}
	return serviceToken, err
}

func WaitForDeletion(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for service token to be deleted, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.SERVICETOKENSTATE_INACTIVE),
		},
		Target: []string{
			string(fabricv4.SERVICETOKENSTATE_DELETED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			serviceToken, body, err := client.ServiceTokensApi.GetServiceTokenByUuid(ctx, uuid).Execute()
			if err != nil {
				if body.StatusCode >= 400 && body.StatusCode <= 499 {
					// Already deleted resource
					return serviceToken, string(fabricv4.SERVICETOKENSTATE_DELETED), nil
				}
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return serviceToken, string(serviceToken.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
