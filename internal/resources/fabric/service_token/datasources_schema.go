package service_token

import (
	equinix_fabric_schema "github.com/equinix/terraform-provider-equinix/internal/fabric/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service Token Type; VC_TOKEN,EPL_TOKEN",
		},
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix-assigned service token identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "An absolute URL that is the subject of the link's context.",
		},
		"issuer_side": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Information about token side; ASIDE, ZSIDE",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the Service Token",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Optional Description to the Service Token you will be creating",
		},
		"expiration_date_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Expiration date and time of the service token; 2020-11-06T07:00:00Z",
		},
		"service_token_connection": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Service Token Connection Type Information",
			Elem:        serviceTokenConnectionSch(),
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Service token state; ACTIVE, INACTIVE, EXPIRED, DELETED",
		},
		"notifications": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Preferences for notifications on Service Token configuration or status changes",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.NotificationSch(),
			},
		},
		"account": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Customer account information that is associated with this service token",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.AccountSch(),
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
		"project": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Project information",
			Elem: &schema.Resource{
				Schema: equinix_fabric_schema.ProjectSch(),
			},
		},
	}
}

func paginationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The page offset for the pagination request. Index of the first element. Default is 0.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of elements to be requested per page. Number must be between 1 and 100. Default is 20",
			},
			"total": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Total number of elements returned.",
			},
			"next": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL relative to the last item in the response.",
			},
			"previous": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL relative to the first item in the response.",
			},
		},
	}
}

func dataSourceSearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Route Filters",
			Elem: &schema.Resource{
				Schema: dataSourceBaseSchema(),
			},
		},
		"filter": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Filters for the Data Source Search Request. Maximum of 8 total filters.",
			MaxItems:    10,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"property": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The API response property which you want to filter your request on. Can be one of the following: \"/type\", \"/name\", \"/project/projectId\", \"/uuid\", \"/state\"",
						ValidateFunc: validation.StringInSlice([]string{"/uuid", "/state", "/name", "/project/projectId"}, true),
					},
					"operator": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Possible operators to use on the filter property. Can be one of the following: [ \"=\", \"!=\", \"[NOT] LIKE\", \"[NOT] IN\", \"ILIKE\" ]",
						ValidateFunc: validation.StringInSlice([]string{"=", "!=", "[NOT] LIKE", "[NOT] IN", "ILIKE"}, true),
					},
					"values": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "The values that you want to apply the property+operator combination to in order to filter your data search",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"pagination": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Pagination details for the Data Source Search Request",
			MaxItems:    1,
			Elem:        paginationSchema(),
		},
	}
}
