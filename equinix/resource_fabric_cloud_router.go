package equinix

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func fabricCloudRouterPackageSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Fabric Cloud Router package code",
		},
	}
}
func fabricCloudRouterAccountSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_number": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "Account Number",
		},
	}
}
func fabricCloudRouterProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Project Id",
		},
		"href": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Unique Resource URL",
		},
	}
}

func fabricCloudRouterResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Equinix-assigned Fabric Cloud Router identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Fabric Cloud Router URI information",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 24),
			Description:  "Fabric Cloud Router name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Customer-provided Fabric Cloud Router description",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Fabric Cloud Router overall state",
		},
		"equinix_asn": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix ASN",
		},
		"package": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Fabric Cloud Router Package Type",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: fabricCloudRouterPackageSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures Fabric Cloud Router lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
			},
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"XF_ROUTER"}, true),
			Description:  "Defines the FCR type like; XF_ROUTER",
		},
		"location": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Fabric Cloud Router location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.LocationSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Customer resource hierarchy project information.Applicable to customers onboarded to Equinix Identity and Access Management. For more information see Identity and Access Management: Projects",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: fabricCloudRouterProjectSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Customer account information that is associated with this Fabric Cloud Router",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: fabricCloudRouterAccountSch(),
			},
		},
		"order": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Order information related to this Fabric Cloud Router",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.OrderSch(),
			},
		},
		"notifications": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Preferences for notifications on Fabric Cloud Router configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.NotificationSch(),
			},
		},
		"bgp_ipv4_routes_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of IPv4 BGP routes in use (including non-distinct prefixes)",
		},
		"bgp_ipv6_routes_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of IPv6 BGP routes in use (including non-distinct prefixes)",
		},
		"distinct_ipv4_prefixes_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of distinct IPv4 routes",
		},
		"distinct_ipv6_prefixes_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of distinct IPv6 routes",
		},
		"connections_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of connections associated with this Fabric Cloud Router instance",
		},
	}
}

func resourceFabricCloudRouter() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
		},
		ReadContext:   resourceFabricCloudRouterRead,
		CreateContext: resourceFabricCloudRouterCreate,
		UpdateContext: resourceFabricCloudRouterUpdate,
		DeleteContext: resourceFabricCloudRouterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: fabricCloudRouterResourceSchema(),

		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric Cloud Router",
	}
}

func accountCloudRouterTerraToGo(accountList []interface{}) v4.SimplifiedAccount {
	sa := v4.SimplifiedAccount{}
	for _, ll := range accountList {
		llMap := ll.(map[string]interface{})
		ac := llMap["account_number"].(int)
		sa = v4.SimplifiedAccount{AccountNumber: int64(ac)}
	}
	return sa
}
func locationCloudRouterTerraToGo(locationList []interface{}) v4.SimplifiedLocationWithoutIbx {
	sl := v4.SimplifiedLocationWithoutIbx{}
	for _, ll := range locationList {
		llMap := ll.(map[string]interface{})
		mc := llMap["metro_code"].(string)
		sl = v4.SimplifiedLocationWithoutIbx{MetroCode: mc}
	}
	return sl
}
func packageCloudRouterTerraToGo(packageList []interface{}) v4.CloudRouterPackageType {
	p := v4.CloudRouterPackageType{}
	for _, pl := range packageList {
		plMap := pl.(map[string]interface{})
		code := plMap["code"].(string)
		p = v4.CloudRouterPackageType{Code: code}
	}
	return p
}
func projectCloudRouterTerraToGo(projectRequest []interface{}) v4.Project {
	if projectRequest == nil {
		return v4.Project{}
	}
	mappedPr := v4.Project{}
	for _, pr := range projectRequest {
		prMap := pr.(map[string]interface{})
		projectId := prMap["project_id"].(string)
		mappedPr = v4.Project{ProjectId: projectId}
	}
	return mappedPr
}
func resourceFabricCloudRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsToFabric(schemaNotifications)
	schemaAccount := d.Get("account").(*schema.Set).List()
	account := accountCloudRouterTerraToGo(schemaAccount)
	schemaLocation := d.Get("location").(*schema.Set).List()
	location := locationCloudRouterTerraToGo(schemaLocation)
	project := v4.Project{}
	schemaProject := d.Get("project").(*schema.Set).List()
	if len(schemaProject) != 0 {
		project = projectCloudRouterTerraToGo(schemaProject)
	}
	schemaPackage := d.Get("package").(*schema.Set).List()
	packages := packageCloudRouterTerraToGo(schemaPackage)

	createRequest := v4.CloudRouterPostRequest{
		Name:          d.Get("name").(string),
		Type_:         d.Get("type").(string),
		Location:      &location,
		Notifications: notifications,
		Package_:      &packages,
		Account:       &account,
		Project:       &project,
	}

	if v, ok := d.GetOk("order"); ok {
		order := equinix_fabric_schema.OrderToFabric(v.(*schema.Set).List())
		createRequest.Order = &order
	}

	start := time.Now()
	fcr, _, err := client.CloudRoutersApi.CreateCloudRouter(ctx, createRequest)
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(fcr.Uuid)

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if _, err = waitUntilCloudRouterIsProvisioned(d.Id(), meta, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for Cloud Router (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricCloudRouterRead(ctx, d, meta)
}

func resourceFabricCloudRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	CloudRouter, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Fabric Cloud Router %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(CloudRouter.Uuid)
	return setCloudRouterMap(d, CloudRouter)
}

func setCloudRouterMap(d *schema.ResourceData, fcr v4.CloudRouter) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"name":                         fcr.Name,
		"href":                         fcr.Href,
		"type":                         fcr.Type_,
		"state":                        fcr.State,
		"package":                      packageCloudRouterGoToTerra(fcr.Package_),
		"location":                     equinix_fabric_schema.LocationWithoutIBXToTerra(fcr.Location),
		"change_log":                   equinix_fabric_schema.ChangeLogToTerra(fcr.ChangeLog),
		"account":                      accountCloudRouterToTerra(fcr.Account),
		"notifications":                equinix_fabric_schema.NotificationsToTerra(fcr.Notifications),
		"project":                      equinix_fabric_schema.ProjectToTerra(fcr.Project),
		"equinix_asn":                  fcr.EquinixAsn,
		"bgp_ipv4_routes_count":        fcr.BgpIpv4RoutesCount,
		"bgp_ipv6_routes_count":        fcr.BgpIpv6RoutesCount,
		"distinct_ipv4_prefixes_count": fcr.DistinctIpv4PrefixesCount,
		"distinct_ipv6_prefixes_count": fcr.DistinctIpv6PrefixesCount,
		"connections_count":            fcr.ConnectionsCount,
		"order":                        equinix_fabric_schema.OrderToTerra(fcr.Order),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func accountCloudRouterToTerra(account *v4.SimplifiedAccount) *schema.Set {
	if account == nil {
		return nil
	}
	accounts := []*v4.SimplifiedAccount{account}
	mappedAccounts := make([]interface{}, len(accounts))
	for i, account := range accounts {
		mappedAccounts[i] = map[string]interface{}{
			"account_number": int(account.AccountNumber),
		}
	}
	accountSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: equinix_fabric_schema.AccountSch()}),
		mappedAccounts,
	)

	return accountSet
}
func packageCloudRouterGoToTerra(packageType *v4.CloudRouterPackageType) *schema.Set {
	packageTypes := []*v4.CloudRouterPackageType{packageType}
	mappedPackages := make([]interface{}, len(packageTypes))
	for i, packageType := range packageTypes {
		mappedPackages[i] = map[string]interface{}{
			"code": packageType.Code,
		}
	}
	packageSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricCloudRouterPackageSch()}),
		mappedPackages,
	)
	return packageSet
}
func getCloudRouterUpdateRequest(conn v4.CloudRouter, d *schema.ResourceData) (v4.CloudRouterChangeOperation, error) {
	changeOps := v4.CloudRouterChangeOperation{}
	existingName := conn.Name
	existingPackage := conn.Package_.Code
	updateNameVal := d.Get("name")
	updatePackageVal := d.Get("conn.Package_.Code")

	log.Printf("existing name %s, existing Package %s, Update Name Request %s, Update Package Request %s ",
		existingName, existingPackage, updateNameVal, updatePackageVal)

	if existingName != updateNameVal {
		changeOps = v4.CloudRouterChangeOperation{Op: "replace", Path: "/name", Value: &updateNameVal}
	} else if existingPackage != updatePackageVal {
		changeOps = v4.CloudRouterChangeOperation{Op: "replace", Path: "/package", Value: &updatePackageVal}
	} else {
		return changeOps, fmt.Errorf("nothing to update for the connection %s", existingName)
	}
	return changeOps, nil
}

func resourceFabricCloudRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbConn, err := waitUntilCloudRouterIsProvisioned(d.Id(), meta, ctx, updateTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("either timed out or errored out while fetching Fabric Cloud Router for uuid %s and error %v", d.Id(), err)
	}
	// TO-DO
	update, err := getCloudRouterUpdateRequest(dbConn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updates := []v4.CloudRouterChangeOperation{update}
	_, _, err = client.CloudRoutersApi.UpdateCloudRouterByUuid(ctx, updates, d.Id())
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	updateFg := v4.CloudRouter{}
	updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	updateFg, err = waitForCloudRouterUpdateCompletion(d.Id(), meta, ctx, updateTimeout)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful Fabric Cloud Router update, error %v", err))
	}

	d.SetId(updateFg.Uuid)
	return setCloudRouterMap(d, updateFg)
}

func waitForCloudRouterUpdateCompletion(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.CloudRouter, error) {
	log.Printf("Waiting for Cloud Router update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{string(v4.PROVISIONED_CloudRouterAccessPointState)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.CloudRouter{}

	if err == nil {
		dbConn = inter.(v4.CloudRouter)
	}
	return dbConn, err
}

func waitUntilCloudRouterIsProvisioned(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.CloudRouter, error) {
	log.Printf("Waiting for Cloud Router to be provisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.PROVISIONING_CloudRouterAccessPointState),
		},
		Target: []string{
			string(v4.PROVISIONED_CloudRouterAccessPointState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.CloudRouter{}

	if err == nil {
		dbConn = inter.(v4.CloudRouter)
	}
	return dbConn, err
}

func resourceFabricCloudRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	start := time.Now()
	_, err := client.CloudRoutersApi.DeleteCloudRouterByUuid(ctx, d.Id())
	if err != nil {
		errors, ok := err.(v4.GenericSwaggerError).Model().([]v4.ModelError)
		if ok {
			// EQ-3040055 = There is an existing update in REQUESTED state
			if equinix_errors.HasModelErrorCode(errors, "EQ-3040055") {
				return diags
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilCloudRouterDeprovisioned(d.Id(), meta, ctx, deleteTimeout)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func WaitUntilCloudRouterDeprovisioned(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for Fabric Cloud Router to be deprovisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(v4.DEPROVISIONING_CloudRouterAccessPointState),
		},
		Target: []string{
			string(v4.DEPROVISIONED_CloudRouterAccessPointState),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbConn, string(*dbConn.State), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
