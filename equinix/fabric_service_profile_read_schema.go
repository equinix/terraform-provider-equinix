package equinix

import (
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readFabricServiceProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service Profile URI response attribute",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service profile type - L2_PROFILE, L3_PROFILE, ECIA_PROFILE, ECMC_PROFILE",
		},
		"visibility": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service profile visibility - PUBLIC, PRIVATE",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer-assigned service profile name",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix assigned service profile identifier",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "User-provided service description",
		},
		"notifications": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Preferences for notifications on connection configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_schema.NotificationSch(),
			},
		},
		"access_point_type_configs": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Access point config information",
			Elem: &schema.Resource{
				Schema: AccessPointTypeConfigSch(),
			},
		},

		"custom_fields": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Custom Fields",
			Elem: &schema.Resource{
				Schema: readCustomFieldSch(),
			},
		},
		"marketing_info": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Marketing Info",
			Elem: &schema.Resource{
				Schema: readMarketingInfoSch(),
			},
		},
		"ports": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Ports",
			Elem: &schema.Resource{
				Schema: readServiceProfileAccessPointColo(),
			},
		},
		"allowed_emails": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Array of contact emails",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "Tags attached to the connection",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"metros": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Access point config information",
			Elem: &schema.Resource{
				Schema: readServiceMetroSch(),
			},
		},
		"self_profile": {
			Type:        schema.TypeBool,
			Computed:    true,
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
			Description: "Account",
			Elem: &schema.Resource{
				Schema: readSPAccountSch(),
			},
		},
		"project": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Project information",
			Elem: &schema.Resource{
				Schema: equinix_schema.ProjectSch(),
			},
		},
		"change_log": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Captures connection lifecycle change information",
			Elem: &schema.Resource{
				Schema: equinix_schema.ChangeLogSch(),
			},
		},
	}
}

var readSpAccountRes = &schema.Resource{
	Schema: readSPAccountSch(),
}

func readSPAccountSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_number": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Account Number",
		},
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Account Name",
		},
		"org_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Customer organization identifier",
		},
		"organization_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Customer organization name",
		},
		"global_org_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global organization identifier",
		},
		"global_organization_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global organization name",
		},
		"global_cust_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global Customer organization identifier",
		},
		"ucm_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Enterprise datastore id",
		},
	}
}

func readCustomFieldSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Label",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description",
		},
		"required": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Required field",
		},
		"data_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Data type",
		},
		"options": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Options",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"capture_in_email": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Required field",
		},
	}
}

func readMarketingInfoSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"logo": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Logo",
		},
		"promotion": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Promotion",
		},
		"process_step": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of process steps",
			Elem: &schema.Resource{
				Schema: readProcessStepSch(),
			},
		},
	}
}

func readProcessStepSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Title",
		},
		"sub_title": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Sub Title",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description",
		},
	}
}

func readServiceProfileAccessPointColo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Colo/Port Type",
		},
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Colo/Port Uuid",
		},
		"location": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Colo/Port Location",
			Elem: &schema.Resource{
				Schema: equinix_schema.LocationSch(),
			},
		},
		"seller_region": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Seller Region",
		},
		"seller_region_description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Seller Region details",
		},
		"cross_connect_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cross Connect Id",
		},
	}
}

func readServiceMetroSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Metro Code - Example SV",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Metro Name",
		},
		"ibxs": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "IBX- Equinix International Business Exchange list",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"in_trail": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "In Trail",
		},
		"display_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Display Name",
		},
		"seller_regions": {
			Type:        schema.TypeMap,
			Computed:    true,
			Description: "Seller Regions",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

var readAllOfServiceProfileChangeLogRes = &schema.Resource{
	Schema: readAllOfServiceProfileChangeLogSch(),
}

func readAllOfServiceProfileChangeLogSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created by User Key",
		},
		"created_by_full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created by User Full Name",
		},
		"created_by_email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created by User Email Address",
		},
		"created_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created by Date and Time",
		},
		"updated_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Updated by User Key",
		},
		"updated_by_full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Updated by User Full Name",
		},
		"updated_by_email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Updated by User Email Address",
		},
		"updated_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Updated by Date and Time",
		},
		"deleted_by": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deleted by User Key",
		},
		"deleted_by_full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deleted by User Full Name",
		},
		"deleted_by_email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deleted by User Email Address",
		},
		"deleted_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Deleted by Date and Time",
		},
	}
}

func readFabricServiceProfileSchemaUpdated() map[string]*schema.Schema {
	sch := readFabricServiceProfileSchema()
	sch["uuid"].Optional = true
	sch["uuid"].Required = false
	return sch
}

func readFabricServiceProfilesSearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Service Profiles",
			Elem: &schema.Resource{
				Schema: readFabricServiceProfileSchemaUpdated(),
			},
		},
		"view_point": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "flips view between buyer and seller representation. Available values : aSide, zSide. Default value : aSide",
		},
		"filter": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Service Profile Search Filter",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: createServiceProfilesSearchFilterSch(),
			},
		},
		"sort": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Service Profile Sort criteria for Search Request response payload",
			Elem: &schema.Resource{
				Schema: createServiceProfilesSearchSortCriteriaSch(),
			},
		},
	}
}
