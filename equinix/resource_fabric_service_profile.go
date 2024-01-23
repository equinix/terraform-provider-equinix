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
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"

	"github.com/antihax/optional"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFabricServiceProfile() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
			Read:   schema.DefaultTimeout(6 * time.Minute),
		},
		ReadContext:   resourceFabricServiceProfileRead,
		CreateContext: resourceFabricServiceProfileCreate,
		UpdateContext: resourceFabricServiceProfileUpdate,
		DeleteContext: resourceFabricServiceProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      createFabricServiceProfileSchema(),
		Description: "Fabric V4 API compatible resource allows creation and management of Equinix Fabric Service Profile\n\n~> **Note** Equinix Fabric v4 resources and datasources are currently in Beta. The interfaces related to `equinix_fabric_` resources and datasources may change ahead of general availability. Please, do not hesitate to report any problems that you experience by opening a new [issue](https://github.com/equinix/terraform-provider-equinix/issues/new?template=bug.md)",
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
	notifications := equinix_schema.NotificationsToFabric(schemaNotifications)

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
		MarketingInfo:          &spMarketingInfo,
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

	var err error
	var eTag int64 = 0
	_, err, eTag = waitForActiveServiceProfileAndPopulateETag(uuid, meta, ctx)
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
	updatedServiceProfile := v4.ServiceProfile{}
	updatedServiceProfile, err = waitForServiceProfileUpdateCompletion(uuid, meta, ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("errored while waiting for successful service profile update, error %v", err))
	}
	d.SetId(updatedServiceProfile.Uuid)
	return setFabricServiceProfileMap(d, updatedServiceProfile)
}

func waitForServiceProfileUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.ServiceProfile, error) {
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
		Timeout:    1 * time.Minute,
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

func waitForActiveServiceProfileAndPopulateETag(uuid string, meta interface{}, ctx context.Context) (v4.ServiceProfile, error, int64) {
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
		Timeout:    1 * time.Minute,
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
		return diag.Errorf("No uuid found %v ", uuid)
	}
	_, _, err := client.ServiceProfilesApi.DeleteServiceProfileByUuid(ctx, uuid)
	if err != nil {
		return diag.FromErr(equinix_errors.FormatFabricError(err))
	}
	return diags
}

func setFabricServiceProfileMap(d *schema.ResourceData, sp v4.ServiceProfile) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"href":                      sp.Href,
		"type":                      sp.Type_,
		"name":                      sp.Name,
		"uuid":                      sp.Uuid,
		"description":               sp.Description,
		"notifications":             equinix_schema.NotificationsToTerra(sp.Notifications),
		"tags":                      tagsFabricSpToTerra(sp.Tags),
		"visibility":                sp.Visibility,
		"access_point_type_configs": accessPointTypeConfigToTerra(sp.AccessPointTypeConfigs),
		"custom_fields":             customFieldFabricSpToTerra(sp.CustomFields),
		"marketing_info":            marketingInfoMappingToTerra(sp.MarketingInfo),
		"ports":                     accessPointColoFabricSpToTerra(sp.Ports),
		"allowed_emails":            allowedEmailsFabricSpToTerra(sp.AllowedEmails),
		"metros":                    serviceMetroFabricSpToTerra(sp.Metros),
		"self_profile":              sp.SelfProfile,
		"state":                     sp.State,
		"account":                   serviceProfileAccountFabricSpToTerra(sp.Account),
		"project":                   equinix_schema.ProjectToTerra(sp.Project),
		"change_log":                allOfServiceProfileChangeLogToTerra(sp.ChangeLog),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setFabricServiceProfilesListMap(d *schema.ResourceData, spl v4.ServiceProfiles) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := equinix_schema.SetMap(d, map[string]interface{}{
		"data": fabricServiceProfilesListToTerra(spl),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
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
