package service_token

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"
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
		Schema: resourceSchema(),
		CustomizeDiff: customdiff.All(
			customdiff.ValidateChange("service_token_connection", func(ctx context.Context, oldValue, newValue, meta interface{}) error {
				var oldZsideBandwidth, newZsideBandwidth []int
				log.Printf("service token old %v", oldValue)
				log.Printf("service token new %v", newValue)

				for _, connection := range oldValue.(*schema.Set).List() {
					olSupportedBandwidthMap := connection.(map[string]interface{})

					if bandwidth, ok := olSupportedBandwidthMap["supported_bandwidths"]; ok {
						oldSupportedBandwidth := bandwidth.([]interface{})
						log.Printf("old value %v", oldSupportedBandwidth)
						if len(oldSupportedBandwidth) > 0 {
							sort.Ints(oldZsideBandwidth)

						}
					}
				}

				for _, connection := range newValue.(*schema.Set).List() {
					newSupportedBandwidthMap := connection.(map[string]interface{})

					if bandwidth, ok := newSupportedBandwidthMap["supported_bandwidths"]; ok {
						newSupportedBandwidth := bandwidth.([]interface{})
						log.Printf("new value %v", newSupportedBandwidth)
						if len(newSupportedBandwidth) > 0 {
							sort.Ints(newZsideBandwidth)

						}
					}
				}
				log.Printf("!! old value %v", oldZsideBandwidth)
				log.Printf("!! old value %v", oldZsideBandwidth)

				if reflect.DeepEqual(oldZsideBandwidth, newZsideBandwidth) {
					return nil // No diff if equivalent after sorting
				}

				return fmt.Errorf("supported_bandwidths values are different: %v", newZsideBandwidth)

			}),
		),
		Description: `Fabric V4 API compatible resource allows creation and management of Equinix Fabric Service Token`,
	}
}

func extractIntValues(value interface{}) []int {
	var intSlice []int

	switch v := value.(type) {
	case *schema.Set:
		for _, item := range v.List() {
			if val, ok := item.(int); ok {
				intSlice = append(intSlice, val)
			}
		}
	case []interface{}:
		for _, item := range v {
			if val, ok := item.(int); ok {
				intSlice = append(intSlice, val)
			}
		}
	}

	return intSlice
}

//func convertToIntSlice(input interface{}) ([]int, bool) {
//	if input == nil {
//		return []int{}, true
//	}
//
//	// Type assert to a slice of interfaces, if applicable
//	if slice, ok := input.([]interface{}); ok {
//		intSlice := make([]int, 0, len(slice))
//		for _, v := range slice {
//			if intVal, ok := v.(int); ok {
//				intSlice = append(intSlice, intVal)
//			} else {
//				return []int{}, false // Conversion failed
//			}
//		}
//		return intSlice, true
//	}
//	return []int{}, false // Not a valid slice type
//}

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
	if err = waitForStability(d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for service token (%s) to be created: %s", d.Id(), err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	updateRequest := buildUpdateRequest(d)

	start := time.Now()
	serviceToken, _, err := client.ServiceTokensApi.UpdateServiceTokenByUuid(ctx, d.Id()).ServiceTokenChangeOperation(updateRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	if err = waitForStability(d.Id(), meta, d, ctx, updateTimeout); err != nil {
		return diag.Errorf("error waiting for service token (%s) to be updated: %s", d.Id(), err)
	}

	return setServiceTokenMap(d, serviceToken)
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

func waitForStability(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, tieout time.Duration) error {
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
			return serviceToken, string(serviceToken.GetState()), nil
		},
		Timeout:    tieout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
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
