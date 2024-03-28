package equinix

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/antihax/optional"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
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
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	serviceProfile, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, d.Id(), nil)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(serviceProfile.Uuid)
	return setFabricServiceProfileMap(d, serviceProfile)
}

func resourceFabricServiceProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)

	createRequest := getServiceProfileRequestPayload(d)
	sp, _, err := client.ServiceProfilesApi.CreateServiceProfile(ctx, createRequest)
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	d.SetId(sp.Uuid)
	return resourceFabricServiceProfileRead(ctx, d, meta)
}

func getServiceProfileRequestPayload(d *schema.ResourceData) v4.ServiceProfileRequest {
	spType := v4.ServiceProfileTypeEnum(d.Get("type").(string))

	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := equinix_fabric_schema.NotificationsToFabric(schemaNotifications)

	var tags []string
	if d.Get("tags") != nil {
		schemaTags := d.Get("tags").([]interface{})
		tags = converters.IfArrToStringArr(schemaTags)
	}

	spVisibility := v4.ServiceProfileVisibilityEnum(d.Get("visibility").(string))

	var spAllowedEmails []string
	if d.Get("allowed_emails") != nil {
		schemaAllowedEmails := d.Get("allowed_emails").([]interface{})
		spAllowedEmails = converters.IfArrToStringArr(schemaAllowedEmails)
	}

	schemaAccessPointTypeConfigs := d.Get("access_point_type_configs").([]interface{})
	spAccessPointTypeConfigs := accessPointTypeConfigsToFabric(schemaAccessPointTypeConfigs)

	schemaCustomFields := d.Get("custom_fields").([]interface{})
	spCustomFields := customFieldsToFabric(schemaCustomFields)

	schemaMarketingInfo := d.Get("marketing_info").(*schema.Set).List()
	spMarketingInfo := marketingInfoToFabric(schemaMarketingInfo)

	schemaPorts := d.Get("ports").([]interface{})
	spPorts := portsToFabric(schemaPorts)

	schemaVirtualDevices := d.Get("virtual_devices").([]interface{})
	spVirtualDevices := virtualDevicesToFabric(schemaVirtualDevices)

	schemaMetros := d.Get("metros").([]interface{})
	spMetros := metrosToFabric(schemaMetros)

	createRequest := v4.ServiceProfileRequest{
		Name:                   d.Get("name").(string),
		Type_:                  &spType,
		Description:            d.Get("description").(string),
		Notifications:          notifications,
		Tags:                   &tags,
		Visibility:             &spVisibility,
		AllowedEmails:          spAllowedEmails,
		AccessPointTypeConfigs: spAccessPointTypeConfigs,
		CustomFields:           spCustomFields,
		MarketingInfo:          spMarketingInfo,
		Ports:                  spPorts,
		VirtualDevices:         spVirtualDevices,
		Metros:                 spMetros,
		SelfProfile:            d.Get("self_profile").(bool),
	}
	return createRequest
}

func resourceFabricServiceProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	uuid := d.Id()
	updateRequest := getServiceProfileRequestPayload(d)

	start := time.Now()
	updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	var err error
	var eTag int64 = 0
	_, err, eTag = waitForActiveServiceProfileAndPopulateETag(uuid, meta, ctx, updateTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.Errorf("Either timed out or errored out while fetching service profile for uuid %s and error %v", uuid, err)
	}

	_, _, err = client.ServiceProfilesApi.PutServiceProfileByUuid(ctx, updateRequest, strconv.FormatInt(eTag, 10), uuid)
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}

	updateTimeout = d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
	updatedServiceProfile := v4.ServiceProfile{}
	updatedServiceProfile, err = waitForServiceProfileUpdateCompletion(uuid, meta, ctx, updateTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful service profile update, error %v", err))
	}
	d.SetId(updatedServiceProfile.Uuid)
	return setFabricServiceProfileMap(d, updatedServiceProfile)
}

func waitForServiceProfileUpdateCompletion(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.ServiceProfile, error) {
	log.Printf("Waiting for service profile update to complete, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbServiceProfile, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid, nil)
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
	dbSp := v4.ServiceProfile{}

	if err == nil {
		dbSp = inter.(v4.ServiceProfile)
	}
	return dbSp, err
}

func waitForActiveServiceProfileAndPopulateETag(uuid string, meta interface{}, ctx context.Context, timeout time.Duration) (v4.ServiceProfile, error, int64) {
	log.Printf("Waiting for service profile to be in active state, uuid %s", uuid)
	var eTag int64 = 0
	stateConf := &retry.StateChangeConf{
		Target: []string{string(v4.ACTIVE_ServiceProfileStateEnum)},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).FabricClient
			dbServiceProfile, res, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid, nil)
			if err != nil {
				return nil, "", equinix_errors.FormatFabricError(err)
			}

			eTagStr := res.Header.Get("ETag")
			eTag, err = strconv.ParseInt(strings.Trim(eTagStr, "\""), 10, 64)
			if err != nil {
				return nil, "", err
			}

			updatableState := ""
			if *dbServiceProfile.State == v4.ACTIVE_ServiceProfileStateEnum {
				updatableState = string(*dbServiceProfile.State)
			}
			return dbServiceProfile, updatableState, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	inter, err := stateConf.WaitForStateContext(ctx)
	dbServiceProfile := v4.ServiceProfile{}
	if err == nil {
		dbServiceProfile = inter.(v4.ServiceProfile)
	}
	return dbServiceProfile, err, eTag
}

func resourceFabricServiceProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	uuid := d.Id()
	if uuid == "" {
		return diag.Errorf("No uuid found for Service Profile Deletion %v ", uuid)
	}
	start := time.Now()
	_, _, err := client.ServiceProfilesApi.DeleteServiceProfileByUuid(ctx, uuid)
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

func WaitAndCheckServiceProfileDeleted(uuid string, client *v4.APIClient, ctx context.Context, timeout time.Duration) error {
	log.Printf("Waiting for service profile to be in deleted, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{string(v4.DELETED_ServiceProfileStateEnum)},
		Refresh: func() (interface{}, string, error) {
			dbConn, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", equinix_errors.FormatFabricError(err)
			}
			updatableState := ""
			if *dbConn.State == v4.DELETED_ServiceProfileStateEnum {
				updatableState = string(*dbConn.State)
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

func setFabricServiceProfilesListMap(d *schema.ResourceData, spl v4.ServiceProfiles) diag.Diagnostics {
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

func setFabricServiceProfileMap(d *schema.ResourceData, sp v4.ServiceProfile) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, fabricServiceProfileMap(&sp))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func fabricServiceProfileMap(serviceProfile *v4.ServiceProfile) map[string]interface{} {
	if serviceProfile == nil {
		return nil
	}
	return map[string]interface{}{
		"href":                      serviceProfile.Href,
		"type":                      serviceProfile.Type_,
		"name":                      serviceProfile.Name,
		"uuid":                      serviceProfile.Uuid,
		"description":               serviceProfile.Description,
		"notifications":             equinix_fabric_schema.NotificationsToTerra(serviceProfile.Notifications),
		"tags":                      tagsFabricSpToTerra(serviceProfile.Tags),
		"visibility":                serviceProfile.Visibility,
		"access_point_type_configs": accessPointTypeConfigToTerra(serviceProfile.AccessPointTypeConfigs),
		"custom_fields":             customFieldFabricSpToTerra(serviceProfile.CustomFields),
		"marketing_info":            marketingInfoMappingToTerra(serviceProfile.MarketingInfo),
		"ports":                     accessPointColoFabricSpToTerra(serviceProfile.Ports),
		"allowed_emails":            allowedEmailsFabricSpToTerra(serviceProfile.AllowedEmails),
		"metros":                    serviceMetroFabricSpToTerra(serviceProfile.Metros),
		"self_profile":              serviceProfile.SelfProfile,
		"state":                     serviceProfile.State,
		"account":                   equinix_fabric_schema.AccountToTerra(serviceProfile.Account),
		"project":                   equinix_fabric_schema.ProjectToTerra(serviceProfile.Project),
		"change_log":                equinix_fabric_schema.ChangeLogToTerra(serviceProfile.ChangeLog),
	}
}

func resourceServiceProfilesSearchRequest(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).FabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*config.Config).FabricAuthToken)
	schemaFilter := d.Get("filter").(*schema.Set).List()
	filter := serviceProfilesSearchFilterRequestToFabric(schemaFilter)
	var serviceProfileFlt v4.ServiceProfileFilter // Cast ServiceProfile search expression struct type to interface
	serviceProfileFlt = filter
	schemaSort := d.Get("sort").([]interface{})
	sort := serviceProfilesSearchSortRequestToFabric(schemaSort)
	schemaViewPoint := d.Get("view_point").(string)

	if schemaViewPoint != "" && schemaViewPoint != string(v4.A_SIDE_ViewPoint) && schemaViewPoint != string(v4.Z_SIDE_ViewPoint) {
		return diag.FromErr(errors.New("view_point can only be set to aSide or zSide. Omitting it will default to aSide"))
	}

	viewPoint := &v4.ServiceProfilesApiSearchServiceProfilesOpts{
		ViewPoint: optional.NewString(schemaViewPoint),
	}

	if schemaViewPoint == "" {
		viewPoint = nil
	}

	createServiceProfilesSearchRequest := v4.ServiceProfileSearchRequest{
		Filter: &serviceProfileFlt,
		Sort:   sort,
	}
	serviceProfiles, _, err := client.ServiceProfilesApi.SearchServiceProfiles(ctx, createServiceProfilesSearchRequest, viewPoint)
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
	d.SetId(serviceProfiles.Data[0].Uuid)
	return setFabricServiceProfilesListMap(d, serviceProfiles)
}

func customFieldFabricSpToTerra(customFields []v4.CustomField) []interface{} {
	if customFields == nil {
		return nil
	}
	mappedCustomFields := make([]interface{}, len(customFields))
	for index, customField := range customFields {
		mappedCustomFields[index] = map[string]interface{}{
			"label":       customField.Label,
			"description": customField.Description,
			"required":    customField.Required,
			"data_type":   customField.DataType,
			"options":     customField.Options,
		}
	}
	return mappedCustomFields
}

func processStepFabricSpToTerra(processSteps []v4.ProcessStep) []interface{} {
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

func marketingInfoMappingToTerra(mkinfo *v4.MarketingInfo) *schema.Set {
	if mkinfo == nil {
		return nil
	}
	mappedMkInfo := make(map[string]interface{})
	mappedMkInfo["logo"] = mkinfo.Logo
	mappedMkInfo["promotion"] = mkinfo.Promotion
	processSteps := processStepFabricSpToTerra(mkinfo.ProcessSteps)
	if processSteps != nil {
		mappedMkInfo["process_step"] = processSteps
	}
	marketingInfoSet := schema.NewSet(
		schema.HashResource(&schema.Resource{Schema: createMarketingInfoSch()}),
		[]interface{}{mappedMkInfo},
	)
	return marketingInfoSet
}

func accessPointColoFabricSpToTerra(accessPointColos []v4.ServiceProfileAccessPointColo) []interface{} {
	if accessPointColos == nil {
		return nil
	}
	mappedAccessPointColos := make([]interface{}, len(accessPointColos))
	for index, accessPointColo := range accessPointColos {
		mappedAccessPointColos[index] = map[string]interface{}{
			"type":                      accessPointColo.Type_,
			"uuid":                      accessPointColo.Uuid,
			"location":                  equinix_fabric_schema.LocationToTerra(accessPointColo.Location),
			"seller_region":             accessPointColo.SellerRegion,
			"seller_region_description": accessPointColo.SellerRegionDescription,
			"cross_connect_id":          accessPointColo.CrossConnectId,
		}
	}
	return mappedAccessPointColos
}

func serviceMetroFabricSpToTerra(serviceMetros []v4.ServiceMetro) []interface{} {
	if serviceMetros == nil {
		return nil
	}
	mappedServiceMetros := make([]interface{}, len(serviceMetros))
	for index, serviceMetro := range serviceMetros {
		mappedServiceMetros[index] = map[string]interface{}{
			"code":           serviceMetro.Code,
			"name":           serviceMetro.Name,
			"ibxs":           serviceMetro.Ibxs,
			"in_trail":       serviceMetro.InTrail,
			"display_name":   serviceMetro.DisplayName,
			"seller_regions": serviceMetro.SellerRegions,
		}
	}
	return mappedServiceMetros
}

func tagsFabricSpToTerra(tags *[]string) []interface{} {
	if tags == nil {
		return nil
	}
	mappedTags := make([]interface{}, len(*tags))
	for index, tag := range *tags {
		mappedTags[index] = tag
	}
	return mappedTags
}

func allowedEmailsFabricSpToTerra(allowedemails []string) []interface{} {
	if allowedemails == nil {
		return nil
	}
	mappedEmails := make([]interface{}, len(allowedemails))
	for index, email := range allowedemails {
		mappedEmails[index] = email
	}
	return mappedEmails
}

func accessPointTypeConfigsToFabric(schemaAccessPointTypeConfigs []interface{}) []v4.ServiceProfileAccessPointType {
	if schemaAccessPointTypeConfigs == nil {
		return []v4.ServiceProfileAccessPointType{}
	}
	var accessPointTypeConfigs []v4.ServiceProfileAccessPointType
	for _, accessPoint := range schemaAccessPointTypeConfigs {
		spType := v4.ServiceProfileAccessPointTypeEnum(accessPoint.(map[string]interface{})["type"].(string))
		spConnectionRedundancyRequired := accessPoint.(map[string]interface{})["connection_redundancy_required"].(bool)
		spAllowBandwidthAutoApproval := accessPoint.(map[string]interface{})["allow_bandwidth_auto_approval"].(bool)
		spAllowRemoteConnections := accessPoint.(map[string]interface{})["allow_remote_connections"].(bool)
		spConnectionLabel := accessPoint.(map[string]interface{})["connection_label"].(string)
		spEnableAutoGenerateServiceKey := accessPoint.(map[string]interface{})["enable_auto_generate_service_key"].(bool)
		spBandwidthAlertThreshold := accessPoint.(map[string]interface{})["bandwidth_alert_threshold"].(float64)
		spAllowCustomBandwidth := accessPoint.(map[string]interface{})["allow_custom_bandwidth"].(bool)

		var spApiConfig *v4.ApiConfig
		if accessPoint.(map[string]interface{})["api_config"] != nil {
			apiConfig := accessPoint.(map[string]interface{})["api_config"].(interface{}).(*schema.Set).List()
			spApiConfig = apiConfigToFabric(apiConfig)
		}

		var spAuthenticationKey *v4.AuthenticationKey
		if accessPoint.(map[string]interface{})["authentication_key"] != nil {
			authenticationKey := accessPoint.(map[string]interface{})["authentication_key"].(interface{}).(*schema.Set).List()
			spAuthenticationKey = authenticationKeyToFabric(authenticationKey)
		}

		supportedBandwidthsRaw := accessPoint.(map[string]interface{})["supported_bandwidths"].([]interface{})
		spSupportedBandwidths := converters.ListToInt32List(supportedBandwidthsRaw)

		accessPointTypeConfigs = append(accessPointTypeConfigs, v4.ServiceProfileAccessPointType{
			Type_:                        &spType,
			ConnectionRedundancyRequired: spConnectionRedundancyRequired,
			AllowBandwidthAutoApproval:   spAllowBandwidthAutoApproval,
			AllowRemoteConnections:       spAllowRemoteConnections,
			ConnectionLabel:              spConnectionLabel,
			EnableAutoGenerateServiceKey: spEnableAutoGenerateServiceKey,
			BandwidthAlertThreshold:      spBandwidthAlertThreshold,
			AllowCustomBandwidth:         spAllowCustomBandwidth,
			ApiConfig:                    spApiConfig,
			AuthenticationKey:            spAuthenticationKey,
			SupportedBandwidths:          &spSupportedBandwidths,
		})
	}
	return accessPointTypeConfigs
}

func accessPointTypeConfigToTerra(spAccessPointTypes []v4.ServiceProfileAccessPointType) []interface{} {
	mappedSpAccessPointTypes := make([]interface{}, len(spAccessPointTypes))
	for index, spAccessPointType := range spAccessPointTypes {
		mappedSpAccessPointTypes[index] = map[string]interface{}{
			"type":                             string(*spAccessPointType.Type_),
			"uuid":                             spAccessPointType.Uuid,
			"allow_remote_connections":         spAccessPointType.AllowRemoteConnections,
			"allow_custom_bandwidth":           spAccessPointType.AllowCustomBandwidth,
			"allow_bandwidth_auto_approval":    spAccessPointType.AllowBandwidthAutoApproval,
			"enable_auto_generate_service_key": spAccessPointType.EnableAutoGenerateServiceKey,
			"connection_redundancy_required":   spAccessPointType.ConnectionRedundancyRequired,
			"connection_label":                 spAccessPointType.ConnectionLabel,
			"api_config":                       apiConfigToTerra(spAccessPointType.ApiConfig),
			"authentication_key":               authenticationKeyToTerra(spAccessPointType.AuthenticationKey),
			"supported_bandwidths":             supportedBandwidthsToTerra(spAccessPointType.SupportedBandwidths),
		}
	}

	return mappedSpAccessPointTypes
}

func apiConfigToFabric(apiConfigs []interface{}) *v4.ApiConfig {
	if apiConfigs == nil {
		return nil
	}
	var apiConfigRes *v4.ApiConfig
	for _, apiConfig := range apiConfigs {
		psApiAvailable := apiConfig.(map[string]interface{})["api_available"].(interface{}).(bool)
		psEquinixManagedVlan := apiConfig.(map[string]interface{})["equinix_managed_vlan"].(interface{}).(bool)
		psBandwidthFromApi := apiConfig.(map[string]interface{})["bandwidth_from_api"].(interface{}).(bool)
		psIntegrationId := apiConfig.(map[string]interface{})["integration_id"].(interface{}).(string)
		psEquinixManagedPort := apiConfig.(map[string]interface{})["equinix_managed_port"].(interface{}).(bool)
		apiConfigRes = &v4.ApiConfig{
			ApiAvailable:       psApiAvailable,
			EquinixManagedVlan: psEquinixManagedVlan,
			BandwidthFromApi:   psBandwidthFromApi,
			IntegrationId:      psIntegrationId,
			EquinixManagedPort: psEquinixManagedPort,
		}
	}
	return apiConfigRes
}

func authenticationKeyToFabric(authenticationKeys []interface{}) *v4.AuthenticationKey {
	if authenticationKeys == nil {
		return nil
	}
	var authenticationKeyRes *v4.AuthenticationKey
	for _, authenticationKey := range authenticationKeys {
		psRequired := authenticationKey.(map[string]interface{})["required"].(interface{}).(bool)
		psLabel := authenticationKey.(map[string]interface{})["label"].(interface{}).(string)
		psDescription := authenticationKey.(map[string]interface{})["description"].(interface{}).(string)
		authenticationKeyRes = &v4.AuthenticationKey{
			Required:    psRequired,
			Label:       psLabel,
			Description: psDescription,
		}
	}
	return authenticationKeyRes
}

func customFieldsToFabric(schemaCustomField []interface{}) []v4.CustomField {
	if schemaCustomField == nil {
		return []v4.CustomField{}
	}
	var customFields []v4.CustomField
	for _, customField := range schemaCustomField {
		cfLabel := customField.(map[string]interface{})["label"].(string)
		cfDescription := customField.(map[string]interface{})["description"].(string)
		cfRequired := customField.(map[string]interface{})["required"].(bool)
		cfDataType := customField.(map[string]interface{})["data_type"].(string)
		optionsRaw := customField.(map[string]interface{})["options"].([]interface{})
		cfOptions := converters.IfArrToStringArr(optionsRaw)
		cfCaptureInEmail := customField.(map[string]interface{})["capture_in_email"].(bool)
		customFields = append(customFields, v4.CustomField{
			Label:          cfLabel,
			Description:    cfDescription,
			Required:       cfRequired,
			DataType:       cfDataType,
			Options:        cfOptions,
			CaptureInEmail: cfCaptureInEmail,
		})
	}
	return customFields
}

func marketingInfoToFabric(schemaMarketingInfos []interface{}) *v4.MarketingInfo {
	if schemaMarketingInfos == nil {
		return nil
	}
	marketingInfoRes := v4.MarketingInfo{}
	for _, marketingInfo := range schemaMarketingInfos {
		miLogo := marketingInfo.(map[string]interface{})["logo"].(string)
		miPromotion := marketingInfo.(map[string]interface{})["promotion"].(bool)

		var miProcessSteps []v4.ProcessStep
		if marketingInfo.(map[string]interface{})["process_steps"] != nil {
			processStepsList := marketingInfo.(map[string]interface{})["process_steps"].([]interface{})
			miProcessSteps = processStepToFabric(processStepsList)
		}

		marketingInfoRes = v4.MarketingInfo{
			Logo:         miLogo,
			Promotion:    miPromotion,
			ProcessSteps: miProcessSteps,
		}
	}
	return &marketingInfoRes
}

func processStepToFabric(processSteps []interface{}) []v4.ProcessStep {
	if processSteps == nil {
		return nil
	}
	processStepRes := make([]v4.ProcessStep, len(processSteps))
	for index, processStep := range processSteps {
		psTitle := processStep.(map[string]interface{})["title"].(interface{}).(string)
		psSubTitle := processStep.(map[string]interface{})["sub_title"].(interface{}).(string)
		psDescription := processStep.(map[string]interface{})["description"].(interface{}).(string)
		processStepRes[index] = v4.ProcessStep{
			Title:       psTitle,
			SubTitle:    psSubTitle,
			Description: psDescription,
		}
	}
	return processStepRes
}

func portsToFabric(schemaPorts []interface{}) []v4.ServiceProfileAccessPointColo {
	if schemaPorts == nil {
		return nil
	}
	serviceProfileAccessPointColos := make([]v4.ServiceProfileAccessPointColo, len(schemaPorts))
	for index, schemaPort := range schemaPorts {
		pType := schemaPort.(map[string]interface{})["type"].(string)
		pUuid := schemaPort.(map[string]interface{})["uuid"].(string)
		locationList := schemaPort.(map[string]interface{})["location"].(interface{}).(*schema.Set).List()
		pLocation := v4.SimplifiedLocation{}
		if len(locationList) != 0 {
			pLocation = equinix_fabric_schema.LocationToFabric(locationList)
		}
		pSellerRegion := schemaPort.(map[string]interface{})["seller_region"].(string)
		pSellerRegionDescription := schemaPort.(map[string]interface{})["seller_region_description"].(string)
		pCrossConnectId := schemaPort.(map[string]interface{})["cross_connect_id"].(string)
		serviceProfileAccessPointColos[index] = v4.ServiceProfileAccessPointColo{
			Type_:                   pType,
			Uuid:                    pUuid,
			Location:                &pLocation,
			SellerRegion:            pSellerRegion,
			SellerRegionDescription: pSellerRegionDescription,
			CrossConnectId:          pCrossConnectId,
		}
	}
	return serviceProfileAccessPointColos
}

func virtualDevicesToFabric(schemaVirtualDevices []interface{}) []v4.ServiceProfileAccessPointVd {
	if schemaVirtualDevices == nil {
		return []v4.ServiceProfileAccessPointVd{}
	}
	var virtualDevices []v4.ServiceProfileAccessPointVd
	for _, virtualDevice := range schemaVirtualDevices {
		vType := virtualDevice.(map[string]interface{})["type"].(string)
		vUuid := virtualDevice.(map[string]interface{})["uuid"].(string)
		locationList := virtualDevice.(map[string]interface{})["location"].(interface{}).(*schema.Set).List()
		vLocation := v4.SimplifiedLocation{}
		if len(locationList) != 0 {
			vLocation = equinix_fabric_schema.LocationToFabric(locationList)
		}
		pInterfaceUuid := virtualDevice.(map[string]interface{})["interface_uuid"].(string)
		virtualDevices = append(virtualDevices, v4.ServiceProfileAccessPointVd{
			Type_:         vType,
			Uuid:          vUuid,
			Location:      &vLocation,
			InterfaceUuid: pInterfaceUuid,
		})
	}
	return virtualDevices
}

func metrosToFabric(schemaMetros []interface{}) []v4.ServiceMetro {
	if schemaMetros == nil {
		return []v4.ServiceMetro{}
	}
	var metros []v4.ServiceMetro
	for _, metro := range schemaMetros {
		mCode := metro.(map[string]interface{})["code"].(string)
		mName := metro.(map[string]interface{})["name"].(string)
		ibxsRaw := metro.(map[string]interface{})["ibxs"].([]interface{})
		mIbxs := converters.IfArrToStringArr(ibxsRaw)
		mInTrail := metro.(map[string]interface{})["in_trail"].(bool)
		mDisplayName := metro.(map[string]interface{})["display_name"].(string)
		mSellerRegions := metro.(map[string]interface{})["seller_regions"].(map[string]string)
		metros = append(metros, v4.ServiceMetro{
			Code:          mCode,
			Name:          mName,
			Ibxs:          mIbxs,
			InTrail:       mInTrail,
			DisplayName:   mDisplayName,
			SellerRegions: mSellerRegions,
		})
	}
	return metros
}

func serviceProfilesSearchFilterRequestToFabric(schemaServiceProfileFilterRequest []interface{}) v4.ServiceProfileSimpleExpression {
	if schemaServiceProfileFilterRequest == nil {
		return v4.ServiceProfileSimpleExpression{}
	}
	mappedFilter := v4.ServiceProfileSimpleExpression{}
	for _, s := range schemaServiceProfileFilterRequest {
		sProperty := s.(map[string]interface{})["property"].(string)
		operator := s.(map[string]interface{})["operator"].(string)
		valuesRaw := s.(map[string]interface{})["values"].([]interface{})
		values := converters.IfArrToStringArr(valuesRaw)
		mappedFilter = v4.ServiceProfileSimpleExpression{Property: sProperty, Operator: operator, Values: values}
	}
	return mappedFilter
}

func serviceProfilesSearchSortRequestToFabric(schemaServiceProfilesSearchSortRequest []interface{}) []v4.ServiceProfileSortCriteria {
	if schemaServiceProfilesSearchSortRequest == nil {
		return []v4.ServiceProfileSortCriteria{}
	}
	var spSortCriteria []v4.ServiceProfileSortCriteria
	for _, sp := range schemaServiceProfilesSearchSortRequest {
		serviceProfileSortCriteriaMap := sp.(map[string]interface{})
		direction := serviceProfileSortCriteriaMap["direction"]
		directionCont := v4.ServiceProfileSortDirection(direction.(string))
		property := serviceProfileSortCriteriaMap["property"]
		propertyCont := v4.ServiceProfileSortBy(property.(string))
		spSortCriteria = append(spSortCriteria, v4.ServiceProfileSortCriteria{
			Direction: &directionCont,
			Property:  &propertyCont,
		})

	}
	return spSortCriteria
}
