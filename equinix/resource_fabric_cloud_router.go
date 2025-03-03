package equinix

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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
			Required:    true,
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
			Required:    true,
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

		Description: `Fabric V4 API compatible resource allows creation and management of [Equinix Fabric Cloud Router](https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-intro.htm#HowItWorks).

Additional documentation:
* Getting Started: https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-intro.htm#HowItWorks
* API: https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#fabric-cloud-routers`,
	}
}

func accountCloudRouterTerraformToGo(accountList []interface{}) fabricv4.SimplifiedAccount {
	if len(accountList) == 0 {
		return fabricv4.SimplifiedAccount{}
	}
	simplifiedAccount := fabricv4.SimplifiedAccount{}
	accountMap := accountList[0].(map[string]interface{})
	accountNumber := int64(accountMap["account_number"].(int))
	simplifiedAccount.SetAccountNumber(accountNumber)

	return simplifiedAccount
}

func packageCloudRouterTerraformToGo(packageList []interface{}) fabricv4.CloudRouterPostRequestPackage {
	if len(packageList) == 0 {
		return fabricv4.CloudRouterPostRequestPackage{}
	}

	package_ := fabricv4.CloudRouterPostRequestPackage{}
	packageMap := packageList[0].(map[string]interface{})
	code := fabricv4.CloudRouterPostRequestPackageCode(packageMap["code"].(string))
	package_.SetCode(code)

	return package_
}
func projectCloudRouterTerraformToGo(projectTerraform []interface{}) fabricv4.Project {
	if len(projectTerraform) == 0 {
		return fabricv4.Project{}
	}
	project := fabricv4.Project{}
	projectMap := projectTerraform[0].(map[string]interface{})
	projectId := projectMap["project_id"].(string)
	project.SetProjectId(projectId)

	return project
}
func marketplaceSubscriptionCloudRouterTerraformToGo(marketplaceSubscriptionTerraform []interface{}) fabricv4.MarketplaceSubscription {
	if len(marketplaceSubscriptionTerraform) == 0 {
		return fabricv4.MarketplaceSubscription{}
	}
	marketplaceSubscription := fabricv4.MarketplaceSubscription{}
	marketplaceSubscriptionMap := marketplaceSubscriptionTerraform[0].(map[string]interface{})
	subscriptionUUID := marketplaceSubscriptionMap["uuid"].(string)
	subscriptionType := marketplaceSubscriptionMap["type"].(string)
	if subscriptionUUID != "" {
		marketplaceSubscription.SetUuid(subscriptionUUID)
	}
	if subscriptionType != "" {
		marketplaceSubscription.SetType(fabricv4.MarketplaceSubscriptionType(subscriptionType))
	}

	return marketplaceSubscription
}
func resourceFabricCloudRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)

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

	if orderTerraform, ok := d.GetOk("order"); ok {
		order := equinix_fabric_schema.OrderTerraformToGo(orderTerraform.(*schema.Set).List())
		createCloudRouterRequest.SetOrder(order)
	}

	if marketplaceSubscriptionTerraform, ok := d.GetOk("marketplace_subscription"); ok {
		marketplaceSubscription := marketplaceSubscriptionCloudRouterTerraformToGo(marketplaceSubscriptionTerraform.(*schema.Set).List())
		createCloudRouterRequest.SetMarketplaceSubscription(marketplaceSubscription)
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
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
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
		"name":                     fcr.GetName(),
		"uuid":                     fcr.GetUuid(),
		"href":                     fcr.GetHref(),
		"type":                     string(fcr.GetType()),
		"state":                    string(fcr.GetState()),
		"package":                  packageCloudRouterGoToTerraform(&package_),
		"location":                 equinix_fabric_schema.LocationWithoutIBXGoToTerraform(&location),
		"change_log":               equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog),
		"account":                  accountCloudRouterGoToTerraform(&account),
		"notifications":            equinix_fabric_schema.NotificationsGoToTerraform(notifications),
		"project":                  equinix_fabric_schema.ProjectGoToTerraform(&project),
		"equinix_asn":              fcr.GetEquinixAsn(),
		"connections_count":        fcr.GetConnectionsCount(),
		"order":                    equinix_fabric_schema.OrderGoToTerraform(&order),
		"marketplace_subscription": marketplaceSubscriptionCloudRouterGoToTerraform(&marketplaceSubscription),
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
func getCloudRouterUpdateRequests(cr *fabricv4.CloudRouter, d *schema.ResourceData) ([][]fabricv4.CloudRouterChangeOperation, error) {
	existingName := cr.GetName()
	existingPackage := cr.GetPackage()
	existingNotifications := cr.GetNotifications()
	updateNameVal := d.Get("name").(string)

	schemaPackage := d.Get("package").(*schema.Set).List()
	package_ := packageCloudRouterTerraformToGo(schemaPackage)
	updatePackageVal := package_.GetCode()

	schemaNotifications := d.Get("notifications").([]interface{})
	updateNotificationsVal := equinix_fabric_schema.NotificationsTerraformToGo(schemaNotifications)
	prevEmails, nextEmails := make([]string, len(existingNotifications[0].GetEmails())), make([]string, len(updateNotificationsVal[0].GetEmails()))
	copy(prevEmails, existingNotifications[0].GetEmails())
	copy(nextEmails, updateNotificationsVal[0].GetEmails())
	sort.Strings(prevEmails)
	sort.Strings(nextEmails)

	notificationsNeedsUpdate := len(updateNotificationsVal) > len(existingNotifications) ||
		string(existingNotifications[0].GetType()) != string(updateNotificationsVal[0].GetType()) ||
		!reflect.DeepEqual(prevEmails, nextEmails)

	log.Printf("[INFO] existing name %s, existing package code %s, new name %s, new package code %s ",
		existingName, existingPackage.GetCode(), updateNameVal, updatePackageVal)

	var changeOps [][]fabricv4.CloudRouterChangeOperation

	if existingName != updateNameVal {
		changeOps = append(changeOps, []fabricv4.CloudRouterChangeOperation{{Op: "replace", Path: "/name", Value: updateNameVal}})
	}

	if string(existingPackage.GetCode()) != string(updatePackageVal) {
		changeOps = append(changeOps, []fabricv4.CloudRouterChangeOperation{{Op: "replace", Path: "/package/code", Value: updatePackageVal}})
	}

	if notificationsNeedsUpdate {
		changeOps = append(changeOps, []fabricv4.CloudRouterChangeOperation{{Op: "replace", Path: "/notifications", Value: updateNotificationsVal}})
	}

	if len(changeOps) == 0 {
		return changeOps, fmt.Errorf("nothing to update for the fabric cloud router %s; the value terraform is detecting a change for does not have any modification available through the api. please revert to previous value to avoid incorrect change numbers in your plan", existingName)
	}

	return changeOps, nil
}

func resourceFabricCloudRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	dbCR, err := waitUntilCloudRouterIsProvisioned(d.Id(), meta, d, ctx, updateTimeout)
	if err != nil {
		return diag.Errorf("either timed out or errored out while fetching Fabric Cloud Router for uuid %s and error %v", d.Id(), err)
	}

	diags := diag.Diagnostics{}
	updates, err := getCloudRouterUpdateRequests(dbCR, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{Severity: 1, Summary: err.Error()})
		return diags
	}
	for _, update := range updates {
		_, _, err = client.CloudRoutersApi.UpdateCloudRouterByUuid(ctx, d.Id()).CloudRouterChangeOperation(update).Execute()
		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("cloud router property update request error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", equinix_errors.FormatFabricError(err), update)})
			continue
		}

		updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
		updateCloudRouter, err := waitForCloudRouterUpdateCompletion(d.Id(), meta, d, ctx, updateTimeout)

		if err != nil {
			diags = append(diags, diag.Diagnostic{Severity: 0, Summary: fmt.Sprintf("cloud router property update completion timeout error: %v [update payload: %v] (other updates will be successful if the payload is not shown)", equinix_errors.FormatFabricError(err), update)})
		} else {
			dbCR = updateCloudRouter
		}
	}

	d.SetId(dbCR.GetUuid())
	return append(diags, setCloudRouterMap(d, dbCR)...)
}

func waitForCloudRouterUpdateCompletion(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.CloudRouter, error) {
	log.Printf("Waiting for Cloud Router update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{string(fabricv4.CLOUDROUTERACCESSPOINTSTATE_PROVISIONED)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			dbCR, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbCR, string(dbCR.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbCR *fabricv4.CloudRouter

	if err == nil {
		dbCR = inter.(*fabricv4.CloudRouter)
	}
	return dbCR, err
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
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			dbCR, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbCR, string(dbCR.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbCR *fabricv4.CloudRouter

	if err == nil {
		dbCR = inter.(*fabricv4.CloudRouter)
	}
	return dbCR, err
}

func resourceFabricCloudRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
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
			client := meta.(*config.Config).NewFabricClientForSDK(ctx, d)
			dbCR, _, err := client.CloudRoutersApi.GetCloudRouterByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			return dbCR, string(dbCR.GetState()), nil
		},
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
