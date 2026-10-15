package connection

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
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

// Resource returns the schema.Resource for managing Equinix Fabric connections.
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
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
			if d.Id() == "" {
				return nil
			}

			if !d.HasChange("a_side") {
				return nil
			}

			_oldASide, _newASide := d.GetChange("a_side")
			oldAside := connectionSideTerraformToGo(_oldASide.(*schema.Set).List())
			newAside := connectionSideTerraformToGo(_newASide.(*schema.Set).List())

			oldLinkProtocol := oldAside.GetAccessPoint().LinkProtocol
			newLinkProtocol := newAside.GetAccessPoint().LinkProtocol

			if oldLinkProtocol == nil || newLinkProtocol == nil {
				return nil
			}

			allowedTypesForVlanChange := []string{string(fabricv4.CONNECTIONTYPE_EVPL_VC), string(fabricv4.CONNECTIONTYPE_EIA_VC)}

			connType := d.Get("type").(string)
			if oldLinkProtocol.VlanTag != nil && newLinkProtocol.VlanTag != nil && *oldLinkProtocol.VlanTag != *newLinkProtocol.VlanTag {
				if !slices.Contains(allowedTypesForVlanChange, connType) {
					return fmt.Errorf(
						"vlan update not allowed for connection of type %s",
						connType,
					)
				}

				if newLinkProtocol.Type == nil || oldLinkProtocol.Type == nil {
					return fmt.Errorf("invalid link protocol state")
				}

				if *oldLinkProtocol.Type != *newLinkProtocol.Type {
					return fmt.Errorf("link protocol type update not allowed")
				}

				if *newLinkProtocol.Type != fabricv4.LINKPROTOCOLTYPE_DOT1_Q {
					return fmt.Errorf(
						"vlan update not allowed for link protocol of type %s",
						*newLinkProtocol.Type,
					)
				}
			}

			return nil
		},

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric connection",
	}
}

func resourceFabricConnectionCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)

	createConnectionRequest := fabricv4.ConnectionPostRequest{}

	name := d.Get("name").(string)
	createConnectionRequest.SetName(name)

	conType := d.Get("type").(string)
	createConnectionRequest.SetType(fabricv4.ConnectionType(conType))

	if orderSchema, ok := d.GetOk("order"); ok {
		order := equinix_fabric_schema.OrderTerraformToGo(orderSchema.(*schema.Set).List())
		createConnectionRequest.SetOrder(order)
	}

	schemaNotifications := d.Get("notifications").([]any)
	notifications := equinix_fabric_schema.NotificationsTerraformToGo(schemaNotifications)
	createConnectionRequest.SetNotifications(notifications)

	bandwidth := d.Get("bandwidth").(int)
	createConnectionRequest.SetBandwidth(int32(bandwidth))

	geoScope := d.Get("geo_scope").(string)
	if geoScope != "" {
		createConnectionRequest.SetGeoScope(fabricv4.GeoScopeType(geoScope))
	}

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
			additionalInfo := additionalInfoTerraformToGo(additionalInfoTerraConfig.([]any))
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
	if err = waitUntilConnectionIsCreated(ctx, d.Id(), meta, d, createTimeout); err != nil {
		return diag.Errorf("error waiting for connection (%s) to be created: %s", d.Id(), err)
	}

	awsSecrets, hasAWSSecrets := additionalInfoContainsAWSSecrets(additionalInfoTerraConfig.([]any))
	if hasAWSSecrets {
		patchChangeOperation := []fabricv4.ConnectionChangeOperation{
			{
				Op:    "add",
				Path:  "",
				Value: map[string]any{"additionalInfo": awsSecrets},
			},
		}

		_, _, patchErr := client.ConnectionsApi.UpdateConnectionByUuid(ctx, *conn.Uuid).ConnectionChangeOperation(patchChangeOperation).Execute()
		if patchErr != nil {
			return diag.FromErr(equinix_errors.FormatFabricError(patchErr))
		}

		createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
		if _, statusChangeErr := waitForConnectionProviderStatusChange(ctx, d.Id(), meta, d, createTimeout); statusChangeErr != nil {
			return diag.Errorf("error waiting for AWS Approval for connection %s: %v", d.Id(), statusChangeErr)
		}
	}

	return resourceFabricConnectionRead(ctx, d, meta)
}

func resourceFabricConnectionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
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

func resourceFabricConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbConn, err := verifyConnectionCreated(ctx, d.Id(), meta, d, updateTimeout)
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

		var waitFunction func(ctx context.Context, uuid string, meta any, d *schema.ResourceData, timeout time.Duration) (*fabricv4.Connection, error)
		switch op := update[0].Op; op {
		case "replace":
			// Update type is either name or bandwidth
			waitFunction = waitForConnectionUpdateCompletion
		case "add":
			// Update type is aws secret additionalInfo
			waitFunction = waitForConnectionProviderStatusChange
		default:
			continue
		}

		updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
		conn, err := waitFunction(ctx, d.Id(), meta, d, updateTimeout)

		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("connection property update completion timeout error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", err, update)})
		} else {
			updatedConn = conn
		}
	}

	d.SetId(updatedConn.GetUuid())
	return append(diags, setFabricMap(d, updatedConn)...)
}

// connectionWait holds the knobs that differ between connection waiters.
type connectionWait struct {
	status        func(*fabricv4.Connection) string
	notFoundState string
	pollInterval  time.Duration
}

// connectionWaitOption customizes waitForConnection.
type connectionWaitOption func(*connectionWait)

// withConnectionStatus reads the polled status from somewhere other than the
// connection state (e.g. operation.providerStatus or change.status).
func withConnectionStatus(status func(*fabricv4.Connection) string) connectionWaitOption {
	return func(w *connectionWait) { w.status = status }
}

// withNotFoundAs reports a 404 as the given status instead of failing the wait.
func withNotFoundAs(status string) connectionWaitOption {
	return func(w *connectionWait) { w.notFoundState = status }
}

// withPollInterval overrides the initial delay and minimum poll interval.
func withPollInterval(interval time.Duration) connectionWaitOption {
	return func(w *connectionWait) { w.pollInterval = interval }
}

// waitForConnection polls the connection until its status (the connection
// state, unless overridden) reaches one of target. An empty pending list
// keeps waiting through any non-target status.
func waitForConnection(ctx context.Context, client *fabricv4.APIClient, uuid string, pending, target []string, timeout time.Duration, opts ...connectionWaitOption) (*fabricv4.Connection, error) {
	w := connectionWait{
		status:       func(c *fabricv4.Connection) string { return string(c.GetState()) },
		pollInterval: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(&w)
	}

	stateConf := &retry.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (any, string, error) {
			dbConn, resp, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid).Execute()
			if err != nil {
				if w.notFoundState != "" && resp != nil && resp.StatusCode == http.StatusNotFound {
					return &fabricv4.Connection{}, w.notFoundState, nil
				}
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, w.status(dbConn), nil
		},
		Timeout:    timeout,
		Delay:      w.pollInterval,
		MinTimeout: w.pollInterval,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	return inter.(*fabricv4.Connection), nil
}

func waitForConnectionUpdateCompletion(ctx context.Context, uuid string, meta any, d *schema.ResourceData, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("[DEBUG] Waiting for connection update to complete, uuid %s", uuid)
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	return waitForConnection(ctx, client, uuid,
		nil,
		[]string{"COMPLETED", "SUBMITTED_FOR_APPROVAL"},
		timeout,
		withConnectionStatus(func(c *fabricv4.Connection) string {
			change := c.GetChange()
			if status := string(change.GetStatus()); status == "COMPLETED" {
				return status
			}
			return ""
		}),
	)
}

func waitUntilConnectionIsCreated(ctx context.Context, uuid string, meta any, d *schema.ResourceData, timeout time.Duration) error {
	log.Printf("Waiting for connection to be created, uuid %s", uuid)
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	_, err := waitForConnection(ctx, client, uuid,
		[]string{
			string(fabricv4.CONNECTIONSTATE_PROVISIONING),
		},
		[]string{
			string(fabricv4.CONNECTIONSTATE_PENDING),
			string(fabricv4.CONNECTIONSTATE_PROVISIONED),
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
		},
		timeout,
	)
	return err
}

func waitForConnectionProviderStatusChange(ctx context.Context, uuid string, meta any, d *schema.ResourceData, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("DEBUG: wating for provider status to update. Connection uuid: %s", uuid)
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	return waitForConnection(ctx, client, uuid,
		[]string{
			string(fabricv4.PROVIDERSTATUS_PENDING_APPROVAL),
			string(fabricv4.PROVIDERSTATUS_PROVISIONING),
		},
		[]string{
			string(fabricv4.PROVIDERSTATUS_PROVISIONED),
		},
		timeout,
		withConnectionStatus(func(c *fabricv4.Connection) string {
			operation := c.GetOperation()
			return string(operation.GetProviderStatus())
		}),
	)
}

func verifyConnectionCreated(ctx context.Context, uuid string, meta any, d *schema.ResourceData, timeout time.Duration) (*fabricv4.Connection, error) {
	log.Printf("Waiting for the connection to be in created state, uuid %s", uuid)
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	return waitForConnection(ctx, client, uuid,
		nil,
		[]string{
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
			string(fabricv4.CONNECTIONSTATE_PROVISIONED),
			string(fabricv4.CONNECTIONSTATE_PENDING),
		},
		timeout,
	)
}

func resourceFabricConnectionDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	start := time.Now()
	_, _, err := client.ConnectionsApi.DeleteConnectionByUuid(ctx, d.Id()).Execute()
	if err != nil {
		if isConnectionAlreadyDeleted(err) {
			return diags
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilConnectionDeprovisioned(ctx, d.Id(), meta, d, deleteTimeout)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for connection deletion. ID: %s, Error %v", d.Id(), err))
	}
	return diags
}

// WaitUntilConnectionDeprovisioned waits until the connection is in DEPROVISIONED state, which indicates that the connection has been deleted successfully. This is required as the API allows deletion of the resource, but the actual resource gets deleted only after it is in DEPROVISIONED state.
func WaitUntilConnectionDeprovisioned(ctx context.Context, uuid string, meta any, d *schema.ResourceData, timeout time.Duration) error {
	log.Printf("Waiting for connection to be deprovisioned, uuid %s", uuid)
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	_, err := waitForConnection(ctx, client, uuid,
		[]string{
			string(fabricv4.CONNECTIONSTATE_DEPROVISIONING),
			string(fabricv4.CONNECTIONSTATE_ACTIVE),
			string(fabricv4.CONNECTIONSTATE_PROVISIONED),
			string(fabricv4.CONNECTIONSTATE_PENDING),
		},
		[]string{
			string(fabricv4.CONNECTIONSTATE_DEPROVISIONED),
		},
		timeout,
	)
	return err
}

// isConnectionAlreadyDeleted reports whether err is the API's
// EQ-3142509 "Connection already deleted" validation error.
func isConnectionAlreadyDeleted(err error) bool {
	return hasFabricErrorCode(err, "EQ-3142509")
}

// isConnectionInTransientState reports whether err is the API's
// EQ-3142510 "Connection is in transient state" validation error, returned
// when a delete is attempted while another operation is still in flight.
func isConnectionInTransientState(err error) bool {
	return hasFabricErrorCode(err, "EQ-3142510")
}

func hasFabricErrorCode(err error, code string) bool {
	var genericError *fabricv4.GenericOpenAPIError
	if !errors.As(err, &genericError) {
		return false
	}
	fabricErrs, ok := genericError.Model().([]fabricv4.Error)
	return ok && equinix_errors.HasErrorCode(fabricErrs, code)
}
