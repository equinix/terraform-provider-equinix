package equinix

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func fabricServiceProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service Profile URI response attribute",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Service profile type - L2_PROFILE, L3_PROFILE, ECIA_PROFILE, ECMC_PROFILE",
		},
		"visibility": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Service profile visibility - PUBLIC, PRIVATE",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Customer-assigned service profile name",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix assigned service profile identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User-provided service description",
		},
		"notifications": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Preferences for notifications on connection configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.NotificationSch(),
			},
		},
		"access_point_type_configs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Access point config information",
			Elem: &schema.Resource{
				Schema: createSPAccessPointTypeConfigSch(),
			},
		},
		"custom_fields": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Custom Fields",
			Elem: &schema.Resource{
				Schema: createCustomFieldSch(),
			},
		},
		"marketing_info": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Marketing Info",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createMarketingInfoSch(),
			},
		},
		"ports": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ports",
			Elem: &schema.Resource{
				Schema: createServiceProfileAccessPointColo(),
			},
		},
		"virtual_devices": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Virtual Devices",
			Elem: &schema.Resource{
				Schema: createServiceProfileAccessPointVd(),
			},
		},
		"allowed_emails": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Array of contact emails",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:        schema.TypeList,
			Description: "Tags attached to the connection",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"metros": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Access point config information",
			Elem: &schema.Resource{
				Schema: createServiceMetroSch(),
			},
		},
		"self_profile": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Self Profile indicating if the profile is created for customer's  self use",
		},
		"state": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Service profile state - ACTIVE, PENDING_APPROVAL, DELETED, REJECTED",
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Service Profile Owner Account Information",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.AccountSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Project information",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ProjectSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures connection lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ChangeLogSch(),
			},
		},
		"view_point": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Flips view between buyer and seller representation. Available values : aSide, zSide. Default value : aSide",
			ValidateFunc: validation.StringInSlice([]string{"aSide", "zSide"}, false),
			Default:      "aSide",
		},
	}
}

func createCustomFieldSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Label",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description",
		},
		"required": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Required field",
		},
		"data_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Data type",
		},
		"options": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Options",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"capture_in_email": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Required field",
		},
	}
}

func createMarketingInfoSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"logo": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Logo",
		},
		"promotion": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Promotion",
		},
		"process_step": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Process Step",
			Elem: &schema.Resource{
				Schema: createProcessStepSch(),
			},
		},
	}
}

func createProcessStepSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Title",
		},
		"sub_title": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Sub Title",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description",
		},
	}
}

func createServiceProfileAccessPointColo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Colo/Port Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Colo/Port Uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Colo/Port Location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.LocationSch(),
			},
		},
		"seller_region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Seller Region",
		},
		"seller_region_description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Seller Region details",
		},
		"cross_connect_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cross Connect Id",
		},
	}
}

func createServiceProfileAccessPointVd() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Virtual Device Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Virtual Device Uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Device Location",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.LocationSch(),
			},
		},
		"interface_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Device Interface Uuid",
		},
	}
}

func createServiceMetroSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Metro Code - Example SV",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Metro Name",
		},
		"ibxs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "IBX- Equinix International Business Exchange list",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"in_trail": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "In Trail",
		},
		"display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Display Name",
		},
		"seller_regions": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Seller Regions",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func createSPAccessPointTypeConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Type of access point type config - VD, COLO",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Colo/Port Uuid",
		},
		"connection_redundancy_required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Mandate redundant connections",
		},
		"allow_bandwidth_auto_approval": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting to enable or disable the ability of the buyer to change connection bandwidth without approval of the seller",
		},
		"allow_remote_connections": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting to allow or prohibit remote connections to the service profile",
		},
		"allow_bandwidth_upgrade": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Availability of a bandwidth upgrade. The default is false",
		},
		"connection_label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom name for Connection",
		},
		"enable_auto_generate_service_key": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable auto generate service key",
		},
		"bandwidth_alert_threshold": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "Percentage of port bandwidth at which an allocation alert is generated",
		},
		"allow_custom_bandwidth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting to enable or disable the ability of the buyer to customize the bandwidth",
		},
		"api_config": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Api configuration details",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createApiConfigSch(),
			},
		},
		"authentication_key": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Authentication key details",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createAuthenticationKeySch(),
			},
		},
		"link_protocol_config": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Link protocol configuration details",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createLinkProtocolConfigSch(),
			},
		},
		"supported_bandwidths": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Supported bandwidths",
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
	}
}

func createApiConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_available": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates if it's possible to establish connections based on the given service profile using the Equinix Fabric API.",
		},
		"equinix_managed_vlan": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting indicating that the VLAN is managed by Equinix (true) or not (false)",
		},
		"allow_over_subscription": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting showing that oversubscription support is available (true) or not (false). The default is false",
		},
		"over_subscription_limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Port bandwidth multiplier that determines the total bandwidth that can be allocated to users creating connections to your services. For example, a 10 Gbps port combined with an overSubscriptionLimit parameter value of 10 allows your subscribers to create connections with a total bandwidth of 100 Gbps.",
		},
		"bandwidth_from_api": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates if the connection bandwidth can be obtained directly from the cloud service provider.",
		},
		"integration_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A unique identifier issued during onboarding and used to integrate the customer's service profile with the Equinix Fabric API.",
		},
		"equinix_managed_port": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Setting indicating that the port is managed by Equinix (true) or not (false)",
		},
	}
}

func createAuthenticationKeySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Requirement to configure an authentication key.",
		},
		"label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the parameter that must be provided to authorize the connection.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of authorization key",
		},
	}
}

func createLinkProtocolConfigSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"encapsulation_strategy": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Additional tagging information required by the seller profile.",
		},
		"reuse_vlan_s_tag": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Automatically accept subsequent DOT1Q to QINQ connections that use the same authentication key. These connections will have the same VLAN S-tag assigned as the initial connection.",
		},
		"encapsulation": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Data frames encapsulation standard.UNTAGGED - Untagged encapsulation for EPL connections. DOT1Q - DOT1Q encapsulation standard. QINQ - QINQ encapsulation standard.",
		},
	}
}

func resourceFabricServiceProfile() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
		},
		ReadContext:   resourceFabricServiceProfileRead,
		CreateContext: resourceFabricServiceProfileCreate,
		UpdateContext: resourceFabricServiceProfileUpdate,
		DeleteContext: resourceFabricServiceProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      fabricServiceProfileSchema(),
		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric Service Profile",
	}
}

func resourceFabricServiceProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	serviceProfile, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, d.Id()).Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(serviceProfile.GetUuid())
	return setFabricServiceProfileMap(d, serviceProfile)
}

func resourceFabricServiceProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)

	createRequest := getServiceProfileRequestPayload(d)
	sp, _, err := client.ServiceProfilesApi.CreateServiceProfile(ctx).ServiceProfileRequest(createRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(sp.GetUuid())
	return resourceFabricServiceProfileRead(ctx, d, meta)
}

func getServiceProfileRequestPayload(d *schema.ResourceData) fabricv4.ServiceProfileRequest {
	spType, _ := fabricv4.NewServiceProfileTypeEnumFromValue(d.Get("type").(string))

	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsTerraformToGo(schemaNotifications)

	var tags []string
	if d.Get("tags") != nil {
		schemaTags := d.Get("tags").([]interface{})
		tags = converters.IfArrToStringArr(schemaTags)
	}

	spVisibility, _ := fabricv4.NewServiceProfileVisibilityEnumFromValue(d.Get("visibility").(string))

	var spAllowedEmails []string
	if d.Get("allowed_emails") != nil {
		schemaAllowedEmails := d.Get("allowed_emails").([]interface{})
		spAllowedEmails = converters.IfArrToStringArr(schemaAllowedEmails)
	}

	schemaAccessPointTypeConfigs := d.Get("access_point_type_configs").([]interface{})
	spAccessPointTypeConfigs := accessPointTypeConfigsTerraformToGo(schemaAccessPointTypeConfigs)

	schemaCustomFields := d.Get("custom_fields").([]interface{})
	spCustomFields := customFieldsTerraformToGo(schemaCustomFields)

	schemaMarketingInfo := d.Get("marketing_info").(*schema.Set).List()
	spMarketingInfo := marketingInfoTerraformToGo(schemaMarketingInfo)

	schemaPorts := d.Get("ports").([]interface{})
	spPorts := portsTerraformToGo(schemaPorts)

	schemaVirtualDevices := d.Get("virtual_devices").([]interface{})
	spVirtualDevices := virtualDevicesTerraformToGo(schemaVirtualDevices)

	schemaMetros := d.Get("metros").([]interface{})
	spMetros := metrosTerraformToGo(schemaMetros)

	createRequest := fabricv4.ServiceProfileRequest{
		Name:                   d.Get("name").(string),
		Type:                   *spType,
		Description:            d.Get("description").(string),
		Notifications:          notifications,
		Tags:                   tags,
		Visibility:             spVisibility,
		AllowedEmails:          spAllowedEmails,
		AccessPointTypeConfigs: spAccessPointTypeConfigs,
		CustomFields:           spCustomFields,
		MarketingInfo:          spMarketingInfo,
		Ports:                  spPorts,
		VirtualDevices:         spVirtualDevices,
		Metros:                 spMetros,
		SelfProfile:            d.Get("self_profile").(*bool),
	}
	return createRequest
}

func resourceFabricServiceProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	uuid := d.Id()
	updateRequest := getServiceProfileRequestPayload(d)

	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	var err error
	var eTag int64 = 0
	_, err, eTag = waitForActiveServiceProfileAndPopulateETag(uuid, meta, d, ctx, updateTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("Either timed out or errored out while fetching service profile for uuid %s and error %v", uuid, err)
	}
	_, _, err = client.ServiceProfilesApi.PutServiceProfileByUuid(ctx, uuid).IfMatch(strconv.FormatInt(eTag, 10)).ServiceProfileRequest(updateRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	var updatedServiceProfile *fabricv4.ServiceProfile
	updatedServiceProfile, err = waitForServiceProfileUpdateCompletion(uuid, meta, d, ctx, updateTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful service profile update, error %v", err))
	}
	d.SetId(updatedServiceProfile.GetUuid())
	return setFabricServiceProfileMap(d, updatedServiceProfile)
}

func waitForServiceProfileUpdateCompletion(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.ServiceProfile, error) {
	log.Printf("Waiting for service profile update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbServiceProfile, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			updatableState := "COMPLETED"
			return dbServiceProfile, updatableState, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	var dbSp *fabricv4.ServiceProfile

	if err == nil {
		dbSp = inter.(*fabricv4.ServiceProfile)
	}
	return dbSp, err
}

func waitForActiveServiceProfileAndPopulateETag(uuid string, meta interface{}, d *schema.ResourceData, ctx context.Context, timeout time.Duration) (*fabricv4.ServiceProfile, error, int64) {
	log.Printf("Waiting for service profile to be in active state, uuid %s", uuid)
	var eTag int64 = 0
	stateConf := &retry.StateChangeConf{
		Target: []string{string(fabricv4.SERVICEPROFILESTATEENUM_ACTIVE)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewFabricClientForSDK(d)
			dbServiceProfile, res, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid).Execute()
			if err != nil {
				return nil, "", equinix_errors.FormatFabricError(err)
			}

			eTagStr := res.Header.Get("ETag")
			eTag, err = strconv.ParseInt(strings.Trim(eTagStr, "\""), 10, 64)
			if err != nil {
				return nil, "", err
			}

			updatableState := ""
			if dbServiceProfile.GetState() == fabricv4.SERVICEPROFILESTATEENUM_ACTIVE {
				updatableState = string(dbServiceProfile.GetState())
			}
			return dbServiceProfile, updatableState, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	inter, err := stateConf.WaitForStateContext(ctx)
	var dbServiceProfile *fabricv4.ServiceProfile
	if err == nil {
		dbServiceProfile = inter.(*fabricv4.ServiceProfile)
	}
	return dbServiceProfile, err, eTag
}

func resourceFabricServiceProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	uuid := d.Id()
	if uuid == "" {
		return diag.Errorf("No uuid found for Service Profile Deletion %v ", uuid)
	}

	start := time.Now()
	_, _, err := client.ServiceProfilesApi.DeleteServiceProfileByUuid(ctx, uuid).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	deleteTimeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)
	waitErr := WaitAndCheckServiceProfileDeleted(uuid, client, ctx, deleteTimeout)
	if waitErr != nil {
		return diag.Errorf("Error while waiting for Service Profile deletion: %v", waitErr)
	}

	return diags
}

func WaitAndCheckServiceProfileDeleted(uuid string, client *fabricv4.APIClient, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for service profile to be in deleted, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{string(fabricv4.SERVICEPROFILESTATEENUM_DELETED)},
		Refresh: func() (interface{}, string, error) {
			dbConn, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid).Execute()
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			updatableState := ""
			if dbConn.GetState() == fabricv4.SERVICEPROFILESTATEENUM_DELETED {
				updatableState = string(dbConn.GetState())
			}
			return dbConn, updatableState, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func setFabricServiceProfilesListMap(d *schema.ResourceData, spl *fabricv4.ServiceProfiles) diag.Diagnostics {
	diags := diag.Diagnostics{}
	mappedServiceProfiles := make([]map[string]interface{}, len(spl.Data))
	if spl.Data != nil {
		for index, serviceProfile := range spl.Data {
			mappedServiceProfiles[index] = fabricServiceProfileMap(&serviceProfile)
		}
	} else {
		mappedServiceProfiles = nil
	}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"data": mappedServiceProfiles,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setFabricServiceProfileMap(d *schema.ResourceData, sp *fabricv4.ServiceProfile) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, fabricServiceProfileMap(sp))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func fabricServiceProfileMap(serviceProfile *fabricv4.ServiceProfile) map[string]interface{} {
	if serviceProfile == nil {
		return nil
	}

	marketingInfo := serviceProfile.GetMarketingInfo()
	account := serviceProfile.GetAccount()
	project := serviceProfile.GetProject()
	changeLog := serviceProfile.GetChangeLog()
	return map[string]interface{}{
		"href":                      serviceProfile.GetHref(),
		"type":                      serviceProfile.GetType(),
		"name":                      serviceProfile.GetName(),
		"uuid":                      serviceProfile.GetUuid(),
		"description":               serviceProfile.GetDescription(),
		"notifications":             equinix_fabric_schema.NotificationsGoToTerraform(serviceProfile.GetNotifications()),
		"tags":                      tagsGoToTerraform(serviceProfile.GetTags()),
		"visibility":                serviceProfile.GetVisibility(),
		"access_point_type_configs": accessPointTypeConfigGoToTerraform(serviceProfile.GetAccessPointTypeConfigs()),
		"custom_fields":             customFieldGoToTerraform(serviceProfile.GetCustomFields()),
		"marketing_info":            marketingInfoGoToTerraform(&marketingInfo),
		"ports":                     serviceProfileAccessPointColoGoToTerraform(serviceProfile.GetPorts()),
		"allowed_emails":            allowedEmailsGoToTerraform(serviceProfile.GetAllowedEmails()),
		"metros":                    serviceMetroGoToTerraform(serviceProfile.GetMetros()),
		"self_profile":              serviceProfile.GetSelfProfile(),
		"state":                     serviceProfile.GetState(),
		"account":                   equinix_fabric_schema.AccountGoToTerraform(&account),
		"project":                   equinix_fabric_schema.ProjectGoToTerraform(&project),
		"change_log":                equinix_fabric_schema.ChangeLogGoToTerraform(&changeLog),
	}
}

func resourceServiceProfilesSearchRequest(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewFabricClientForSDK(d)
	schemaFilter := d.Get("filter").(*schema.Set).List()
	filter := serviceProfilesSearchFilterRequestTerraformToGo(schemaFilter)
	schemaSort := d.Get("sort").([]interface{})
	sort := serviceProfilesSearchSortRequestTerraformToGo(schemaSort)
	viewPoint, _ := fabricv4.NewGetServiceProfilesViewPointParameterFromValue(d.Get("view_point").(string))

	createServiceProfilesSearchRequest := fabricv4.ServiceProfileSearchRequest{
		Filter: &fabricv4.ServiceProfileFilter{
			ServiceProfileSimpleExpression: filter,
		},
		Sort: sort,
	}
	serviceProfiles, _, err := client.ServiceProfilesApi.SearchServiceProfiles(ctx).ViewPoint(*viewPoint).ServiceProfileSearchRequest(createServiceProfilesSearchRequest).Execute()
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	if len(serviceProfiles.Data) != 1 {
		error := fmt.Errorf("incorrect # of records are found for the service profile search criteria - %d , please change the criteria", len(serviceProfiles.Data))
		return diag.FromErr(error)
	}
	d.SetId(serviceProfiles.Data[0].GetUuid())
	return setFabricServiceProfilesListMap(d, serviceProfiles)
}

func customFieldGoToTerraform(customFields []fabricv4.CustomField) []interface{} {
	if customFields == nil {
		return nil
	}
	mappedCustomFields := make([]interface{}, len(customFields))
	for index, customField := range customFields {
		mappedCustomFields[index] = map[string]interface{}{
			"label":       customField.GetLabel(),
			"description": customField.GetDescription(),
			"required":    customField.GetRequired(),
			"data_type":   customField.GetDataType(),
			"options":     customField.GetOptions(),
		}
	}
	return mappedCustomFields
}

func processStepsGoToTerraform(processSteps []fabricv4.ProcessStep) []interface{} {
	if processSteps == nil {
		return nil
	}
	mappedProcessSteps := make([]interface{}, len(processSteps))
	for index, processStep := range processSteps {
		mappedProcessSteps[index] = map[string]interface{}{
			"title":       processStep.Title,
			"sub_title":   processStep.SubTitle,
			"description": processStep.Description,
		}
	}
	return mappedProcessSteps
}

func marketingInfoGoToTerraform(mkinfo *fabricv4.MarketingInfo) *schema.Set {
	if mkinfo == nil {
		return nil
	}
	mappedMkInfo := make(map[string]interface{})
	mappedMkInfo["logo"] = mkinfo.GetLogo()
	mappedMkInfo["promotion"] = mkinfo.GetPromotion()
	processSteps := processStepsGoToTerraform(mkinfo.GetProcessSteps())
	if processSteps != nil {
		mappedMkInfo["process_step"] = processSteps
	}
	marketingInfoSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createMarketingInfoSch()}),
		[]interface{}{mappedMkInfo},
	)
	return marketingInfoSet
}

func serviceProfileAccessPointColoGoToTerraform(accessPointColos []fabricv4.ServiceProfileAccessPointCOLO) []interface{} {
	if accessPointColos == nil {
		return nil
	}
	mappedAccessPointColos := make([]interface{}, len(accessPointColos))
	for index, accessPointColo := range accessPointColos {
		location := accessPointColo.GetLocation()
		mappedAccessPointColos[index] = map[string]interface{}{
			"type":                      accessPointColo.GetType(),
			"uuid":                      accessPointColo.GetUuid(),
			"location":                  equinix_fabric_schema.LocationGoToTerraform(&location),
			"seller_region":             accessPointColo.GetSellerRegion(),
			"seller_region_description": accessPointColo.GetSellerRegionDescription(),
			"cross_connect_id":          accessPointColo.GetCrossConnectId(),
		}
	}
	return mappedAccessPointColos
}

func serviceMetroGoToTerraform(serviceMetros []fabricv4.ServiceMetro) []interface{} {
	if serviceMetros == nil {
		return nil
	}
	mappedServiceMetros := make([]interface{}, len(serviceMetros))
	for index, serviceMetro := range serviceMetros {
		mappedServiceMetros[index] = map[string]interface{}{
			"code":           serviceMetro.GetCode(),
			"name":           serviceMetro.GetName(),
			"ibxs":           serviceMetro.GetIbxs(),
			"in_trail":       serviceMetro.GetInTrail(),
			"display_name":   serviceMetro.GetDisplayName(),
			"seller_regions": serviceMetro.GetSellerRegions(),
		}
	}
	return mappedServiceMetros
}

func tagsGoToTerraform(tags []string) []interface{} {
	if tags == nil {
		return nil
	}
	mappedTags := make([]interface{}, len(tags))
	for index, tag := range tags {
		mappedTags[index] = tag
	}
	return mappedTags
}

func allowedEmailsGoToTerraform(allowedEmails []string) []interface{} {
	if allowedEmails == nil {
		return nil
	}
	mappedEmails := make([]interface{}, len(allowedEmails))
	for index, email := range allowedEmails {
		mappedEmails[index] = email
	}
	return mappedEmails
}

func accessPointTypeConfigsTerraformToGo(schemaAccessPointTypeConfigs []interface{}) []fabricv4.ServiceProfileAccessPointType {
	if schemaAccessPointTypeConfigs == nil {
		return nil
	}
	var accessPointTypeConfigs []fabricv4.ServiceProfileAccessPointType
	for _, accessPoint := range schemaAccessPointTypeConfigs {
		apMap := accessPoint.(map[string]interface{})
		spType := fabricv4.ServiceProfileAccessPointTypeEnum(apMap["type"].(string))
		spConnectionRedundancyRequired := apMap["connection_redundancy_required"].(*bool)
		spAllowBandwidthAutoApproval := apMap["allow_bandwidth_auto_approval"].(*bool)
		spAllowRemoteConnections := apMap["allow_remote_connections"].(*bool)
		spConnectionLabel := apMap["connection_label"].(*string)
		spEnableAutoGenerateServiceKey := apMap["enable_auto_generate_service_key"].(*bool)
		spBandwidthAlertThreshold := apMap["bandwidth_alert_threshold"].(*float32)
		spAllowCustomBandwidth := apMap["allow_custom_bandwidth"].(*bool)

		var spApiConfig *fabricv4.ApiConfig
		if apMap["api_config"] != nil {
			apiConfig := apMap["api_config"].(*schema.Set).List()
			spApiConfig = apiConfigTerraformToGo(apiConfig)
		}

		var spAuthenticationKey *fabricv4.AuthenticationKey
		if apMap["authentication_key"] != nil {
			authenticationKey := apMap["authentication_key"].(*schema.Set).List()
			spAuthenticationKey = authenticationKeyTerraformToGo(authenticationKey)
		}

		supportedBandwidthsRaw := apMap["supported_bandwidths"].([]interface{})
		spSupportedBandwidths := converters.ListToInt32List(supportedBandwidthsRaw)

		accessPointTypeConfigs = append(accessPointTypeConfigs, fabricv4.ServiceProfileAccessPointType{
			ServiceProfileAccessPointTypeCOLO: &fabricv4.ServiceProfileAccessPointTypeCOLO{
				Type:                         spType,
				ConnectionRedundancyRequired: spConnectionRedundancyRequired,
				AllowBandwidthAutoApproval:   spAllowBandwidthAutoApproval,
				AllowRemoteConnections:       spAllowRemoteConnections,
				ConnectionLabel:              spConnectionLabel,
				EnableAutoGenerateServiceKey: spEnableAutoGenerateServiceKey,
				BandwidthAlertThreshold:      spBandwidthAlertThreshold,
				AllowCustomBandwidth:         spAllowCustomBandwidth,
				ApiConfig:                    spApiConfig,
				AuthenticationKey:            spAuthenticationKey,
				SupportedBandwidths:          spSupportedBandwidths,
			},
		})
	}
	return accessPointTypeConfigs
}

func accessPointTypeConfigGoToTerraform(spAccessPointTypes []fabricv4.ServiceProfileAccessPointType) []interface{} {
	mappedSpAccessPointTypes := make([]interface{}, len(spAccessPointTypes))
	for index, spAccessPointType := range spAccessPointTypes {
		spAccessPointType := spAccessPointType.GetActualInstance().(fabricv4.ServiceProfileAccessPointTypeCOLO)
		apiConfig := spAccessPointType.GetApiConfig()
		authKey := spAccessPointType.GetAuthenticationKey()
		mappedSpAccessPointTypes[index] = map[string]interface{}{
			"type":                             string(spAccessPointType.Type),
			"uuid":                             spAccessPointType.GetUuid(),
			"allow_remote_connections":         spAccessPointType.GetAllowRemoteConnections(),
			"allow_custom_bandwidth":           spAccessPointType.GetAllowCustomBandwidth(),
			"allow_bandwidth_auto_approval":    spAccessPointType.GetAllowBandwidthAutoApproval(),
			"enable_auto_generate_service_key": spAccessPointType.GetEnableAutoGenerateServiceKey(),
			"connection_redundancy_required":   spAccessPointType.GetConnectionRedundancyRequired(),
			"connection_label":                 spAccessPointType.GetConnectionLabel(),
			"api_config":                       apiConfigGoToTerraform(&apiConfig),
			"authentication_key":               authenticationKeyGoToTerraform(&authKey),
			"supported_bandwidths":             supportedBandwidthsGoToTerraform(spAccessPointType.GetSupportedBandwidths()),
		}
	}

	return mappedSpAccessPointTypes
}

func apiConfigGoToTerraform(apiConfig *fabricv4.ApiConfig) *schema.Set {

	mappedApiConfig := make(map[string]interface{})
	mappedApiConfig["api_available"] = apiConfig.ApiAvailable
	mappedApiConfig["equinix_managed_vlan"] = apiConfig.EquinixManagedVlan
	mappedApiConfig["bandwidth_from_api"] = apiConfig.BandwidthFromApi
	mappedApiConfig["integration_id"] = apiConfig.IntegrationId
	mappedApiConfig["equinix_managed_port"] = apiConfig.EquinixManagedPort

	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createApiConfigSch()}),
		[]interface{}{mappedApiConfig})
	return apiConfigSet
}

func authenticationKeyGoToTerraform(authenticationKey *fabricv4.AuthenticationKey) *schema.Set {
	mappedAuthenticationKey := make(map[string]interface{})
	mappedAuthenticationKey["required"] = authenticationKey.Required
	mappedAuthenticationKey["label"] = authenticationKey.Label
	mappedAuthenticationKey["description"] = authenticationKey.Description

	apiConfigSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createAuthenticationKeySch()}),
		[]interface{}{mappedAuthenticationKey})
	return apiConfigSet
}

func supportedBandwidthsGoToTerraform(supportedBandwidths []int32) []interface{} {
	if supportedBandwidths == nil {
		return nil
	}
	mappedSupportedBandwidths := make([]interface{}, len(supportedBandwidths))
	for index, bandwidth := range supportedBandwidths {
		mappedSupportedBandwidths[index] = int(bandwidth)
	}
	return mappedSupportedBandwidths
}

func apiConfigTerraformToGo(apiConfigs []interface{}) *fabricv4.ApiConfig {
	if apiConfigs == nil {
		return nil
	}
	var apiConfig *fabricv4.ApiConfig
	apiConfigMap := apiConfigs[0].(map[string]interface{})
	psApiAvailable := apiConfigMap["api_available"].(*bool)
	psEquinixManagedVlan := apiConfigMap["equinix_managed_vlan"].(*bool)
	psBandwidthFromApi := apiConfigMap["bandwidth_from_api"].(*bool)
	psIntegrationId := apiConfigMap["integration_id"].(*string)
	psEquinixManagedPort := apiConfigMap["equinix_managed_port"].(*bool)
	apiConfig = &fabricv4.ApiConfig{
		ApiAvailable:       psApiAvailable,
		EquinixManagedVlan: psEquinixManagedVlan,
		BandwidthFromApi:   psBandwidthFromApi,
		IntegrationId:      psIntegrationId,
		EquinixManagedPort: psEquinixManagedPort,
	}

	return apiConfig
}

func authenticationKeyTerraformToGo(authenticationKeys []interface{}) *fabricv4.AuthenticationKey {
	if authenticationKeys == nil {
		return nil
	}
	var authenticationKeyRes *fabricv4.AuthenticationKey
	authKeyMap := authenticationKeys[0].(map[string]interface{})
	psRequired := authKeyMap["required"].(*bool)
	psLabel := authKeyMap["label"].(*string)
	psDescription := authKeyMap["description"].(*string)
	authenticationKeyRes = &fabricv4.AuthenticationKey{
		Required:    psRequired,
		Label:       psLabel,
		Description: psDescription,
	}

	return authenticationKeyRes
}

func customFieldsTerraformToGo(schemaCustomField []interface{}) []fabricv4.CustomField {
	if schemaCustomField == nil {
		return nil
	}
	customFields := make([]fabricv4.CustomField, len(schemaCustomField))
	for _, customField := range schemaCustomField {
		cfMap := customField.(map[string]interface{})
		cfLabel := cfMap["label"].(string)
		cfDescription := cfMap["description"].(string)
		cfRequired := cfMap["required"].(*bool)
		cfDataType, _ := fabricv4.NewCustomFieldDataTypeFromValue(cfMap["data_type"].(string))
		optionsRaw := cfMap["options"].([]interface{})
		cfOptions := converters.IfArrToStringArr(optionsRaw)
		cfCaptureInEmail := cfMap["capture_in_email"].(*bool)
		customFields = append(customFields, fabricv4.CustomField{
			Label:          cfLabel,
			Description:    cfDescription,
			Required:       cfRequired,
			DataType:       *cfDataType,
			Options:        cfOptions,
			CaptureInEmail: cfCaptureInEmail,
		})
	}
	return customFields
}

func marketingInfoTerraformToGo(schemaMarketingInfos []interface{}) *fabricv4.MarketingInfo {
	if schemaMarketingInfos == nil {
		return nil
	}
	var marketingInfo *fabricv4.MarketingInfo
	marketingInfoMap := schemaMarketingInfos[0].(map[string]interface{})
	miLogo := marketingInfoMap["logo"].(*string)
	miPromotion := marketingInfoMap["promotion"].(*bool)

	var miProcessSteps []fabricv4.ProcessStep
	if marketingInfoMap["process_steps"] != nil {
		processStepsList := marketingInfoMap["process_steps"].([]interface{})
		miProcessSteps = processStepTerraformToGo(processStepsList)
	}

	marketingInfo = &fabricv4.MarketingInfo{
		Logo:         miLogo,
		Promotion:    miPromotion,
		ProcessSteps: miProcessSteps,
	}

	return marketingInfo
}

func processStepTerraformToGo(processSteps []interface{}) []fabricv4.ProcessStep {
	if processSteps == nil {
		return nil
	}
	processStepRes := make([]fabricv4.ProcessStep, len(processSteps))
	for index, processStep := range processSteps {
		processStepMap := processStep.(map[string]interface{})
		psTitle := processStepMap["title"].(interface{}).(*string)
		psSubTitle := processStepMap["sub_title"].(interface{}).(*string)
		psDescription := processStepMap["description"].(interface{}).(*string)
		processStepRes[index] = fabricv4.ProcessStep{
			Title:       psTitle,
			SubTitle:    psSubTitle,
			Description: psDescription,
		}
	}
	return processStepRes
}

func portsTerraformToGo(schemaPorts []interface{}) []fabricv4.ServiceProfileAccessPointCOLO {
	if schemaPorts == nil {
		return nil
	}
	serviceProfileAccessPointColos := make([]fabricv4.ServiceProfileAccessPointCOLO, len(schemaPorts))
	for index, schemaPort := range schemaPorts {
		portMap := schemaPort.(map[string]interface{})
		pType, _ := fabricv4.NewServiceProfileAccessPointCOLOTypeFromValue(portMap["type"].(string))
		pUuid := portMap["uuid"].(string)
		locationList := portMap["location"].(interface{}).(*schema.Set).List()
		var pLocation *fabricv4.SimplifiedLocation
		if len(locationList) != 0 {
			pLocation = equinix_fabric_schema.LocationTerraformToGo(locationList)
		}
		pSellerRegion := portMap["seller_region"].(*string)
		pSellerRegionDescription := portMap["seller_region_description"].(*string)
		pCrossConnectId := portMap["cross_connect_id"].(*string)
		serviceProfileAccessPointColos[index] = fabricv4.ServiceProfileAccessPointCOLO{
			Type:                    *pType,
			Uuid:                    pUuid,
			Location:                pLocation,
			SellerRegion:            pSellerRegion,
			SellerRegionDescription: pSellerRegionDescription,
			CrossConnectId:          pCrossConnectId,
		}
	}
	return serviceProfileAccessPointColos
}

func virtualDevicesTerraformToGo(schemaVirtualDevices []interface{}) []fabricv4.ServiceProfileAccessPointVD {
	if schemaVirtualDevices == nil {
		return nil
	}
	virtualDevices := make([]fabricv4.ServiceProfileAccessPointVD, len(schemaVirtualDevices))
	for index, virtualDevice := range schemaVirtualDevices {
		vdMap := virtualDevice.(map[string]interface{})
		vType, _ := fabricv4.NewServiceProfileAccessPointVDTypeFromValue(vdMap["type"].(string))
		vUuid := vdMap["uuid"].(string)
		locationList := vdMap["location"].(interface{}).(*schema.Set).List()
		var vLocation *fabricv4.SimplifiedLocation
		if len(locationList) != 0 {
			vLocation = equinix_fabric_schema.LocationTerraformToGo(locationList)
		}
		pInterfaceUuid := vdMap["interface_uuid"].(*string)
		virtualDevices[index] = fabricv4.ServiceProfileAccessPointVD{
			Type:          *vType,
			Uuid:          vUuid,
			Location:      vLocation,
			InterfaceUuid: pInterfaceUuid,
		}
	}
	return virtualDevices
}

func metrosTerraformToGo(schemaMetros []interface{}) []fabricv4.ServiceMetro {
	if schemaMetros == nil {
		return nil
	}
	var metros []fabricv4.ServiceMetro
	for index, metro := range schemaMetros {
		metroMap := metro.(map[string]interface{})
		mCode := metroMap["code"].(*string)
		mName := metroMap["name"].(*string)
		ibxsRaw := metroMap["ibxs"].([]interface{})
		mIbxs := converters.IfArrToStringArr(ibxsRaw)
		mInTrail := metroMap["in_trail"].(*bool)
		mDisplayName := metroMap["display_name"].(*string)
		mSellerRegions := metroMap["seller_regions"].(*map[string]string)
		metros[index] = fabricv4.ServiceMetro{
			Code:          mCode,
			Name:          mName,
			Ibxs:          mIbxs,
			InTrail:       mInTrail,
			DisplayName:   mDisplayName,
			SellerRegions: mSellerRegions,
		}
	}
	return metros
}

func serviceProfilesSearchFilterRequestTerraformToGo(schemaServiceProfileFilterRequest []interface{}) *fabricv4.ServiceProfileSimpleExpression {
	if schemaServiceProfileFilterRequest == nil {
		return nil
	}
	var mappedFilter *fabricv4.ServiceProfileSimpleExpression
	simpleExpressionMap := schemaServiceProfileFilterRequest[0].(map[string]interface{})
	sProperty := simpleExpressionMap["property"].(*string)
	operator := simpleExpressionMap["operator"].(*string)
	valuesRaw := simpleExpressionMap["values"].([]interface{})
	values := converters.IfArrToStringArr(valuesRaw)
	mappedFilter = &fabricv4.ServiceProfileSimpleExpression{
		Property: sProperty,
		Operator: operator,
		Values:   values,
	}

	return mappedFilter
}

func serviceProfilesSearchSortRequestTerraformToGo(schemaServiceProfilesSearchSortRequest []interface{}) []fabricv4.ServiceProfileSortCriteria {
	if schemaServiceProfilesSearchSortRequest == nil {
		return nil
	}
	spSortCriteria := make([]fabricv4.ServiceProfileSortCriteria, len(schemaServiceProfilesSearchSortRequest))
	for index, sp := range schemaServiceProfilesSearchSortRequest {
		serviceProfileSortCriteriaMap := sp.(map[string]interface{})
		direction := serviceProfileSortCriteriaMap["direction"]
		directionCont, _ := fabricv4.NewServiceProfileSortDirectionFromValue(direction.(string))
		property := serviceProfileSortCriteriaMap["property"]
		propertyCont, _ := fabricv4.NewServiceProfileSortByFromValue(property.(string))
		spSortCriteria[index] = fabricv4.ServiceProfileSortCriteria{
			Direction: directionCont,
			Property:  propertyCont,
		}

	}
	return spSortCriteria
}
