package schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func OrderSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"purchase_order_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Purchase order number",
		},
		"billing_tier": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Billing tier for connection bandwidth",
		},
		"order_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Order Identification",
		},
		"order_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Order Reference Number",
		},
	}
}

func NotificationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Notification Type - ALL,CONNECTION_APPROVAL,SALES_REP_NOTIFICATIONS, NOTIFICATIONS",
		},
		"send_interval": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Send interval",
		},
		"emails": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Array of contact emails",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func AccountSch() map[string]*schema.Schema {
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
		//"ucm_id": {
		//	Type:        schema.TypeString,
		//	Computed:    true,
		//	Description: "Account ucmId",
		//},
		"global_cust_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Global Customer organization identifier",
		},
		//"reseller_account_number": {
		//	Type:        schema.TypeInt,
		//	Computed:    true,
		//	Description: "Reseller account number",
		//},
		//"reseller_account_name": {
		//	Type:        schema.TypeString,
		//	Computed:    true,
		//	Description: "Reseller account name",
		//},
		//"reseller_ucm_id": {
		//	Type:        schema.TypeString,
		//	Computed:    true,
		//	Description: "Reseller account ucmId",
		//},
		//"reseller_org_id": {
		//	Type:        schema.TypeInt,
		//	Computed:    true,
		//	Description: "Reseller customer organization identifier",
		//},
	}
}

func LocationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Access point region",
		},
		"metro_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Access point metro name",
		},
		"metro_code": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Access point metro code",
		},
		"ibx": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "IBX Code",
		},
	}
}

func ProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Project Id",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource URL",
		},
	}
}

func ChangeLogSch() map[string]*schema.Schema {
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

func ErrorSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"error_code": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Error  code",
		},
		"error_message": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Error Message",
		},
		"correlation_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "CorrelationId",
		},
		"details": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Details",
		},
		"help": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Help",
		},
		"additional_info": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Pricing error additional Info",
			Elem: &schema.Resource{
				Schema: ErrorAdditionalInfoSch(),
			},
		},
	}
}

func ErrorAdditionalInfoSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"property": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Property at which the error potentially occurred",
		},
		"reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Reason for the error",
		},
	}
}
