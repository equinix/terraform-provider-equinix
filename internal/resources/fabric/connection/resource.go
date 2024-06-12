package connection

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
		},
		ReadContext:   resourceFabricConnectionRead,
		CreateContext: resourceFabricConnectionCreate,
		UpdateContext: resourceFabricConnectionUpdate,
		DeleteContext: resourceFabricConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: fabricConnectionResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric connection",
	}
}

func resourceFabricConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)

	createConnectionRequest := fabricv4.ConnectionPostRequest{}

	name := d.Get("name").(string)
	createConnectionRequest.SetName(name)

	conType := d.Get("type").(string)
	createConnectionRequest.SetType(fabricv4.ConnectionType(conType))

	if orderSchema, ok := d.GetOk("order"); ok {
		order := equinix_fabric_schema.OrderTerraformToGo(orderSchema.(*schema.Set).List())
		createConnectionRequest.SetOrder(order)
	}

	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsTerraformToGo(schemaNotifications)
	createConnectionRequest.SetNotifications(notifications)

	bandwidth := d.Get("bandwidth").(int)
	createConnectionRequest.SetBandwidth(int32(bandwidth))

	if schemaRedundancy, ok := d.GetOk("redundancy"); ok {
		redundancy := connectionRedundancyTerraformToGo(schemaRedundancy.(*schema.Set).List())
		createConnectionRequest.SetRedundancy(redundancy)
	}

	if terraConfigProject, ok := d.GetOk("project"); ok {
		project := equinix_fabric_schema.ProjectTerraformToGo(terraConfigProject.(*schema.Set).List())
		createConnectionRequest.SetProject(project)
	}

	aSide := d.Get("a_side").(*schema.Set).List()
	connectionASide := connectionSideTerraformToGo(aSide)
	createConnectionRequest.SetASide(connectionASide)

	zSide := d.Get("z_side").(*schema.Set).List()
	connectionZSide := connectionSideTerraformToGo(zSide)
	createConnectionRequest.SetZSide(connectionZSide)

	additionalInfoTerraConfig, ok := d.GetOk("additional_info")
	if ok {
		zSideAccessPoint := connectionZSide.GetAccessPoint()
		zSideAccessPointServiceProfile := zSideAccessPoint.GetProfile()
		serviceProfile, _, _ := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, zSideAccessPointServiceProfile.GetUuid()).Execute()
		customFields := serviceProfile.GetCustomFields()

		if len(customFields) != 0 {
			additionalInfo := additionalInfoTerraformToGo(additionalInfoTerraConfig.([]interface{}))
			createConnectionRequest.SetAdditionalInfo(additionalInfo)
		}
	}

	start := time.Now()
	conn, _, err := client.ConnectionsApi.CreateConnection(ctx).ConnectionPostRequest(createConnectionRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(conn.GetUuid())

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitUntilConnectionIsCreated(d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for connection (%s) to be created: %s", d.Id(), err)
	}

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfoTerraConfig.([]interface{}))
	if hasAWSSecrets {
		patchChangeOperation := []fabricv4.ConnectionChangeOperation{
			{
				Op:    "add",
				Path:  "",
				Value: map[string]interface{}{"additionalInfo": awsSecrets},
			},
		}

		_, _, patchErr := client.ConnectionsApi.UpdateConnectionByUuid(ctx, *conn.Uuid).ConnectionChangeOperation(patchChangeOperation).Execute()
		if patchErr != nil {
			return diag.FromErr(equinix_errors.FormatFabricError(patchErr))
		}

		createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
		if _, statusChangeErr := waitForConnectionProviderStatusChange(d.Id(), meta, d, ctx, createTimeout); statusChangeErr != nil {
			return diag.Errorf("error waiting for AWS Approval for connection %s: %v", d.Id(), statusChangeErr)
		}
	}

	return resourceFabricConnectionRead(ctx, d, meta)
}

func resourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	conn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Connection %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(conn.GetUuid())
	return setFabricMap(d, conn)
}

func resourceFabricConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbConn, err := verifyConnectionCreated(d.Id(), meta, d, ctx, updateTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("either timed out or errored out while fetching connection for uuid %s: error -> %v", d.Id(), err)
	}

	diags := diag.Diagnostics{}
	updateRequests, err := getUpdateRequests(dbConn, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{Severity: 1, Summary: err.Error()})
		return diags
	}
	updatedConn := dbConn

	for _, update := range updateRequests {
		_, _, err := client.ConnectionsApi.UpdateConnectionByUuid(ctx, d.Id()).ConnectionChangeOperation(update).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update request error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", equinix_errors.FormatFabricError(err), update)})
			continue
		}

		var waitFunction func(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Connection, error)
		if update[0].Op == "replace" {
			// Update type is either name or bandwidth
			waitFunction = waitForConnectionUpdateCompletion
		} else if update[0].Op == "add" {
			// Update type is aws secret additionalInfo
			waitFunction = waitForConnectionProviderStatusChange
		}

		updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
		conn, err := waitFunction(d.Id(), meta, d, ctx, updateTimeout)

		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update completion timeout error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", err, update)})
		} else {
			updatedConn = conn
		}
	}

	d.SetId(updatedConn.GetUuid())
	return append(diags, setFabricMap(d, updatedConn)...)
}

func waitForConnectionUpdateCompletion(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("[DEBUG] Waiting for connection update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			updatableState := ""
			change := dbConn.GetChange()
			status := change.GetStatus()
			if string(status) == "COMPLETED" {
				updatableState = string(status)
			}
			return dbConn, updatableState, nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.Connection

	if err == nil {
		dbConn = inter.(*fabricv4.Connection)
	}
	return dbConn, err
}

func waitUntilConnectionIsCreated(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for connection to be created, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONSTATE_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONSTATE_PENDING),
			string(fabricv4.CONNECTIONSTATE_PROVISIONED),
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(dbConn.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func waitForConnectionProviderStatusChange(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("DEBUG: wating for provider status to update. Connection uuid: %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.PROVIDERSTATUS_PENDING_APPROVAL),
			string(fabricv4.PROVIDERSTATUS_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.PROVIDERSTATUS_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			operation := dbConn.GetOperation()
			providerStatus := operation.GetProviderStatus()
			return dbConn, string(providerStatus), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.Connection

	if err == nil {
		dbConn = inter.(*fabricv4.Connection)
	}
	return dbConn, err
}

func verifyConnectionCreated(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("Waiting for connection to be in created state, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
			string(fabricv4.CONNECTIONSTATE_PROVISIONED),
			string(fabricv4.CONNECTIONSTATE_PENDING),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(dbConn.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbConn *fabricv4.Connection

	if err == nil {
		dbConn = inter.(*fabricv4.Connection)
	}
	return dbConn, err
}

func resourceFabricConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	_, _, err := client.ConnectionsApi.DeleteConnectionByUuid(ctx, d.Id()).Execute()
	if err != nil {
		if genericError, ok := err.(*fabricv4.GenericOpenAPIError); ok {
			if fabricErrs, ok := genericError.Model().([]fabricv4.Error); ok {
				// EQ-3142509 = Connection already deleted
				if equinix_errors.HasErrorCode(fabricErrs, "EQ-3142509") {
					return diags
				}
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilConnectionDeprovisioned(d.Id(), meta, d, ctx, deleteTimeout)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func WaitUntilConnectionDeprovisioned(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for connection to be deprovisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CONNECTIONSTATE_DEPROVISIONING),
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
			string(fabricv4.CONNECTIONSTATE_PENDING),
		},
		Target: []string{
			string(fabricv4.CONNECTIONSTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(dbConn.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
