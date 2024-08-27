package equinix

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
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

func fabricMarketplaceSubscriptionSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Marketplace Subscription type like; AWS_MARKETPLACE_SUBSCRIPTION",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Equinix-assigned Marketplace Subscription identifier",
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
			Description: "Customer resource hierarchy project information. Applicable to customers onboarded to Equinix Identity and Access Management. For more information see Identity and Access Management: Projects",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: fabricCloudRouterProjectSch(),
			},
		},
		"marketplace_subscription": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Equinix Fabric Entity for Marketplace Subscription",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: fabricMarketplaceSubscriptionSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Description: "Customer account information that is associated with this Fabric Cloud Router",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: fabricCloudRouterAccountSch(),
			},
		},
		"order": {
			Type:        schema.TypeSet,
			Computed:    true,
			Optional:    true,
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

func accountCloudRouterTerraformToGo(accountList []interface{}) fabricv4.SimplifiedAccount {
	if accountList == nil || len(accountList) == 0 {
		return fabricv4.SimplifiedAccount{}
	}
	simplifiedAccount := fabricv4.SimplifiedAccount{}
	accountMap := accountList[0].(map[string]interface{})
	accountNumber := int64(accountMap["account_number"].(int))
	simplifiedAccount.SetAccountNumber(accountNumber)

	return simplifiedAccount
}

func packageCloudRouterTerraformToGo(packageList []interface{}) fabricv4.CloudRouterPostRequestPackage {
	if packageList == nil || len(packageList) == 0 {
		return fabricv4.CloudRouterPostRequestPackage{}
	}

	package_ := fabricv4.CloudRouterPostRequestPackage{}
	packageMap := packageList[0].(map[string]interface{})
	code := fabricv4.CloudRouterPostRequestPackageCode(packageMap["code"].(string))
	package_.SetCode(code)

	return package_
}
func projectCloudRouterTerraformToGo(projectTerraform []interface{}) fabricv4.Project {
	if projectTerraform == nil || len(projectTerraform) == 0 {
		return fabricv4.Project{}
	}
	project := fabricv4.Project{}
	projectMap := projectTerraform[0].(map[string]interface{})
	projectId := projectMap["project_id"].(string)
	project.SetProjectId(projectId)

	return project
}
func marketplaceSubscriptionCloudRouterTerraformToGo(marketplaceSubscriptionTerraform []interface{}) fabricv4.MarketplaceSubscription {
	if marketplaceSubscriptionTerraform == nil || len(marketplaceSubscriptionTerraform) == 0 {
		return fabricv4.MarketplaceSubscription{}
	}
	marketplaceSubscription := fabricv4.MarketplaceSubscription{}
	marketplaceSubscriptionMap := marketplaceSubscriptionTerraform[0].(map[string]interface{})
	subscriptionUUID := marketplaceSubscriptionMap["uuid"].(string)
	marketplaceSubscription.SetUuid(subscriptionUUID)

	return marketplaceSubscription
}
func resourceFabricCloudRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)

	createCloudRouterRequest := fabricv4.CloudRouterPostRequest{}

	createCloudRouterRequest.SetName(d.Get("name").(string))

	type_ := fabricv4.CloudRouterPostRequestType(d.Get("type").(string))
	createCloudRouterRequest.SetType(type_)

	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsTerraformToGo(schemaNotifications)
	createCloudRouterRequest.SetNotifications(notifications)

	schemaAccount := d.Get("account").(*schema.Set).List()
	account := accountCloudRouterTerraformToGo(schemaAccount)
	createCloudRouterRequest.SetAccount(account)

	schemaLocation := d.Get("location").(*schema.Set).List()
	location := equinix_fabric_schema.LocationWithoutIBXTerraformToGo(schemaLocation)
	createCloudRouterRequest.SetLocation(location)

	schemaProject := d.Get("project").(*schema.Set).List()
	project := projectCloudRouterTerraformToGo(schemaProject)
	createCloudRouterRequest.SetProject(project)

	schemaPackage := d.Get("package").(*schema.Set).List()
	package_ := packageCloudRouterTerraformToGo(schemaPackage)
	createCloudRouterRequest.SetPackage(package_)

	schemaMarketplaceSubscription := d.Get("marketplace_subscription").(*schema.Set).List()
	marketplaceSubscription := marketplaceSubscriptionCloudRouterTerraformToGo(schemaMarketplaceSubscription)
	createCloudRouterRequest.SetMarketplaceSubscription(marketplaceSubscription)

	if orderTerraform, ok := d.GetOk("order"); ok {
		order := equinix_fabric_schema.OrderTerraformToGo(orderTerraform.(*schema.Set).List())
		createCloudRouterRequest.SetOrder(order)
	}

	start := time.Now()
	fcr, _, err := client.CloudRoutersApi.CreateCloudRouter(ctx).CloudRouterPostRequest(createCloudRouterRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(fcr.GetUuid())

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if _, err = waitUntilCloudRouterIsProvisioned(d.Id(), meta, d, ctx, createTimeout); err != nil {
		return diag.Errorf("error waiting for Cloud Router (%s) to be created: %s", d.Id(), err)
	}

	return resourceFabricCloudRouterRead(ctx, d, meta)
}

func resourceFabricCloudRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	cloudRouter, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, d.Id()).Execute()
	if err != nil {
		log.Printf("[WARN] Fabric Cloud Router %s not found , error %s", d.Id(), err)
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(cloudRouter.GetUuid())
	return setCloudRouterMap(d, cloudRouter)
}

func fabricCloudRouterMap(fcr *fabricv4.CloudRouter) map[string]interface{} {
	package_ := fcr.GetPackage()
	location := fcr.GetLocation()
	changeLog := fcr.GetChangeLog()
	account := fcr.GetAccount()
	notifications := fcr.GetNotifications()
	project := fcr.GetProject()
	order := fcr.GetOrder()
	marketplaceSubscription := fcr.GetMarketplaceSubscription()
	return map[string]interface{}{
		"name":                         fcr.GetName(),
		"uuid":                         fcr.GetUuid(),
		"href":                         fcr.GetHref(),
		"type":                         string(fcr.GetType()),
		"state":                        string(fcr.GetState()),
		"package":                      packageCloudRouterGoToTerraform(&package_),
		"location":                     equinix_fabric_schema.LocationWithoutIBXGoToTerraform(&location),
		"change_log":                   equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog),
		"account":                      accountCloudRouterGoToTerraform(&account),
		"notifications":                equinix_fabric_schema.NotificationsGoToTerraform(notifications),
		"project":                      equinix_fabric_schema.ProjectGoToTerraform(&project),
		"equinix_asn":                  fcr.GetEquinixAsn(),
		"bgp_ipv4_routes_count":        fcr.GetBgpIpv4RoutesCount(),
		"bgp_ipv6_routes_count":        fcr.GetBgpIpv6RoutesCount(),
		"distinct_ipv4_prefixes_count": fcr.GetDistinctIpv4PrefixesCount(),
		"distinct_ipv6_prefixes_count": fcr.GetDistinctIpv6PrefixesCount(),
		"connections_count":            fcr.GetConnectionsCount(),
		"order":                        equinix_fabric_schema.OrderGoToTerraform(&order),
		"marketplace_subscription":     marketplaceSubscriptionCloudRouterGoToTerraform(&marketplaceSubscription),
	}
}

func setCloudRouterMap(d *schema.ResourceData, fcr *fabricv4.CloudRouter) diag.Diagnostics {
	diags := diag.Diagnostics{}
	cloudRouterMap := fabricCloudRouterMap(fcr)
	err := equinix_schema.SetMap(d, cloudRouterMap)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func accountCloudRouterGoToTerraform(account *fabricv4.SimplifiedAccount) *schema.Set {
	if account == nil {
		return nil
	}

	mappedAccount := map[string]interface{}{
		"account_number": int(account.GetAccountNumber()),
	}

	accountSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: equinix_fabric_schema.AccountSch()}),
		[]interface{}{mappedAccount},
	)

	return accountSet
}
func packageCloudRouterGoToTerraform(packageType *fabricv4.CloudRouterPostRequestPackage) *schema.Set {
	mappedPackage := map[string]interface{}{
		"code": string(packageType.GetCode()),
	}
	packageSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricCloudRouterPackageSch()}),
		[]interface{}{mappedPackage},
	)
	return packageSet
}
func marketplaceSubscriptionCloudRouterGoToTerraform(subscription *fabricv4.MarketplaceSubscription) *schema.Set {
	if subscription == nil {
		return nil
	}
	mappedSubscription := make(map[string]interface{})
	mappedSubscription["type"] = string(subscription.GetType())
	mappedSubscription["uuid"] = subscription.GetUuid()

	subscriptionSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: fabricMarketplaceSubscriptionSch()}),
		[]interface{}{mappedSubscription})
	return subscriptionSet
}
func getCloudRouterUpdateRequest(conn *fabricv4.CloudRouter, d *schema.ResourceData) (fabricv4.CloudRouterChangeOperation, error) {
	changeOps := fabricv4.CloudRouterChangeOperation{}
	existingName := conn.GetName()
	existingPackage := conn.GetPackage()
	updateNameVal := d.Get("name").(string)
	updatePackageVal := d.Get("package.0.code").(string)

	log.Printf("[INFO] existing name %s, existing package code %s, new name %s, new package code %s ",
		existingName, existingPackage.GetCode(), updateNameVal, updatePackageVal)

	if existingName != updateNameVal {
		changeOps = fabricv4.CloudRouterChangeOperation{Op: "replace", Path: "/name", Value: updateNameVal}
	} else if string(existingPackage.GetCode()) != updatePackageVal {
		changeOps = fabricv4.CloudRouterChangeOperation{Op: "replace", Path: "/package/code", Value: updatePackageVal}
	} else {
		return changeOps, fmt.Errorf("nothing to update for the connection %s", existingName)
	}
	return changeOps, nil
}

func resourceFabricCloudRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbConn, err := waitUntilCloudRouterIsProvisioned(d.Id(), meta, d, ctx, updateTimeout)
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
	updates := []fabricv4.CloudRouterChangeOperation{update}
	_, _, err = client.CloudRoutersApi.UpdateCloudRouterByUuid(ctx, d.Id()).CloudRouterChangeOperation(updates).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	updateCloudRouter, err := waitForCloudRouterUpdateCompletion(d.Id(), meta, d, ctx, updateTimeout)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful Fabric Cloud Router update, error %v", err))
	}

	d.SetId(updateCloudRouter.GetUuid())
	return setCloudRouterMap(d, updateCloudRouter)
}

func waitForCloudRouterUpdateCompletion(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.CloudRouter, error) {
	log.Printf("Waiting for Cloud Router update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{string(fabricv4.CLOUDROUTERACCESSPOINTSTATE_PROVISIONED)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid).Execute()
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
	var dbConn *fabricv4.CloudRouter

	if err == nil {
		dbConn = inter.(*fabricv4.CloudRouter)
	}
	return dbConn, err
}

func waitUntilCloudRouterIsProvisioned(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.CloudRouter, error) {
	log.Printf("Waiting for Cloud Router to be provisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CLOUDROUTERACCESSPOINTSTATE_PROVISIONING),
		},
		Target: []string{
			string(fabricv4.CLOUDROUTERACCESSPOINTSTATE_PROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid).Execute()
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
	var dbConn *fabricv4.CloudRouter

	if err == nil {
		dbConn = inter.(*fabricv4.CloudRouter)
	}
	return dbConn, err
}

func resourceFabricCloudRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	start := time.Now()
	_, err := client.CloudRoutersApi.DeleteCloudRouterByUuid(ctx, d.Id()).Execute()
	if err != nil {
		if genericError, ok := err.(*fabricv4.GenericOpenAPIError); ok {
			if fabricErrs, ok := genericError.Model().([]fabricv4.Error); ok {
				// EQ-3040055 = There is an existing update in REQUESTED state
				if equinix_errors.HasErrorCode(fabricErrs, "EQ-3040055") {
					return diags
				}
			}
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	err = WaitUntilCloudRouterDeprovisioned(d.Id(), meta, d, ctx, deleteTimeout)
	if err != nil {
		return diag.FromErr(fmt.Errorf("API call failed while waiting for resource deletion. Error %v", err))
	}
	return diags
}

func WaitUntilCloudRouterDeprovisioned(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for Fabric Cloud Router to be deprovisioned, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(fabricv4.CLOUDROUTERACCESSPOINTSTATE_DEPROVISIONING),
		},
		Target: []string{
			string(fabricv4.CLOUDROUTERACCESSPOINTSTATE_DEPROVISIONED),
		},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbConn, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid).Execute()
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
