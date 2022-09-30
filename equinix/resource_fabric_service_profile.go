package equinix

import (
	"context"
	"fmt"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4" //TODO: Update to ..equinix-lab/fabric-go project before Production merge
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
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
		Description: "Resource allows creation and management of Equinix Fabric	Service Profiles",
	}
}

func resourceFabricServiceProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	serviceProfile, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, d.Id(), nil)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			error := v4.ModelError{}
			d.SetId("")
			log.Printf("Error Status Message: %s", error.ErrorMessage)
		}
		return diag.FromErr(err)
	}
	d.SetId(serviceProfile.Uuid)
	return setFabricServiceProfileMap(d, serviceProfile)
}

func resourceFabricServiceProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)

	createRequest := getServiceProfileRequestPayload(d)
	sp, _, err := client.ServiceProfilesApi.CreateServiceProfile(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(sp.Uuid)
	return resourceFabricServiceProfileRead(ctx, d, meta)
}

func getServiceProfileRequestPayload(d *schema.ResourceData) v4.ServiceProfileRequest {
	spType := v4.ServiceProfileTypeEnum(d.Get("type").(string))

	schemaNotifications := d.Get("notifications").([]interface{})
	notifications := notificationToFabric(schemaNotifications)

	var tags []string
	if d.Get("tags") != nil {
		schemaTags := d.Get("tags").([]interface{})
		tags = expandListToStringList(schemaTags)
	}

	spVisibility := v4.ServiceProfileVisibilityEnum(d.Get("visibility").(string))

	var spAllowedEmails []string
	if d.Get("allowed_emails") != nil {
		schemaAllowedEmails := d.Get("allowed_emails").([]interface{})
		spAllowedEmails = expandListToStringList(schemaAllowedEmails)
	}

	schemaAccessPointTypeConfigs := d.Get("access_point_type_configs").([]interface{})
	spAccessPointTypeConfigs := accessPointTypeConfigsToFabric(schemaAccessPointTypeConfigs)

	schemaCustomFields := d.Get("custom_fields").([]interface{})
	spCustomFields := customFieldsToFabric(schemaCustomFields)

	schemaMarketingInfo := d.Get("marketing_info").(interface{}).(*schema.Set).List()
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
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	uuid := d.Id()
	//TODO Why we need the below check?
	if uuid == "" {
		return diag.Errorf("No service profile found for the value uuid %v ", uuid)
	}

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

	_, res, err := client.ServiceProfilesApi.PutServiceProfileByUuid(ctx, updateRequest, strconv.FormatInt(eTag, 10), uuid)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error response for the service profile update, response %v, error %v", res, err))
	}
	updatedServiceProfile := v4.ServiceProfile{}
	updatedServiceProfile, err = waitForServiceProfileUpdateCompletion(uuid, meta, ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			d.SetId("")
		}
		return diag.FromErr(fmt.Errorf("Errored while waiting for successful service profile update, response %v, error %v", res, err))
	}
	d.SetId(updatedServiceProfile.Uuid)
	return setFabricServiceProfileMap(d, updatedServiceProfile)
}

func waitForServiceProfileUpdateCompletion(uuid string, meta interface{}, ctx context.Context) (v4.ServiceProfile, error) {
	log.Printf("Waiting for service profile update to complete, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Target: []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbServiceProfile, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", err
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
	stateConf := &resource.StateChangeConf{
		Target: []string{"ACTIVE"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*Config).fabricClient
			dbServiceProfile, res, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid, nil)
			if err != nil {
				return nil, "", err
			}

			eTagStr := res.Header.Get("ETag")
			eTag, err = strconv.ParseInt(strings.Trim(eTagStr, "\""), 10, 64)

			updatableState := ""
			if "ACTIVE" == *dbServiceProfile.State {
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
	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	uuid := d.Id()
	if uuid == "" {
		return diag.Errorf("No uuid found %v ", uuid)
	}
	_, resp, err := client.ServiceProfilesApi.DeleteServiceProfileByUuid(ctx, uuid)
	if err != nil {
		fmt.Errorf("Error response for the Service Profile delete error %v and response %v", err, resp)
		return diag.FromErr(err)
	}
	return diags
}

func setFabricServiceProfileMap(d *schema.ResourceData, sp v4.ServiceProfile) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := setMap(d, map[string]interface{}{
		"href":                      sp.Href,
		"type":                      sp.Type_,
		"name":                      sp.Name,
		"uuid":                      sp.Uuid,
		"description":               sp.Description,
		"notifications":             notificationToTerra(sp.Notifications),
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
		"project":                   projectToTerra(sp.Project),
		"change_log":                allOfServiceProfileChangeLogToTerra(sp.ChangeLog),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setFabricServiceProfilesListMap(d *schema.ResourceData, spl v4.ServiceProfiles) diag.Diagnostics {
	diags := diag.Diagnostics{}
	err := setMap(d, map[string]interface{}{
		"data": fabricServiceProfilesListToTerra(spl),
	})
	if err != nil {

		return diag.FromErr(err)

	}
	return diags
}

func resourceServiceProfilesSearchRequest(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*Config).fabricClient
	ctx = context.WithValue(ctx, v4.ContextAccessToken, meta.(*Config).FabricAuthToken)
	schemaFilter := d.Get("filter").(interface{}).(*schema.Set).List()
	filter := serviceProfilesSearchFilterRequestToFabric(schemaFilter)
	var serviceProfileFlt v4.ServiceProfileFilter //Cast ServiceProfile search expression struct type to interface
	serviceProfileFlt = filter
	schemaSort := d.Get("sort").([]interface{})
	sort := serviceProfilesSearchSortRequestToFabric(schemaSort)
	createServiceProfilesSearchRequest := v4.ServiceProfileSearchRequest{
		Filter: &serviceProfileFlt,
		Sort:   sort,
	}
	serviceProfiles, _, err := client.ServiceProfilesApi.SearchServiceProfiles(ctx, createServiceProfilesSearchRequest)

	if err != nil {
		if !strings.Contains(err.Error(), "500") {
			error := v4.ModelError{}
			d.SetId("")
			log.Printf("Error Status Message: %s", error.ErrorMessage)
		}
		return diag.FromErr(err)
	}

	if len(serviceProfiles.Data) != 1 {
		error := fmt.Errorf("incorrect # of records are found for the service profile search criteria - %d , please change the criteria", len(serviceProfiles.Data))
		return diag.FromErr(error)
	}
	d.SetId(serviceProfiles.Data[0].Uuid)
	return setFabricServiceProfilesListMap(d, serviceProfiles)
}
