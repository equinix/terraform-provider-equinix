package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"
)

func OrderSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"purchase_order_number": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Purchase order number. Short name/number to identify this order on the invoice",
		},
		"billing_tier": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Billing tier for connection bandwidth",
		},
		"order_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Order Identification",
		},
		"order_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Order Reference Number",
		},
	}
}

func NotificationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: equinix_validation.StringInEnumSlice(fabricv4.AllowedSimplifiedNotificationTypeEnumValues, false),
			Description: fmt.Sprintf("Notification type. One of %v", fabricv4.AllowedSimplifiedNotificationTypeEnumValues),
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
			Description: "Equinix-assigned account number.",
		},
		"account_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Legal name of the accountholder.",
		},
		"org_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Equinix-assigned ID of the subscriber's organization.",
		},
		"organization_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned name of the subscriber's organization.",
		},
		"global_org_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned ID of the subscriber's parent organization.",
		},
		"global_organization_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned name of the subscriber's parent organization.",
		},
		"global_cust_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Equinix-assigned ID of the subscriber's parent organization.",
		},
		"ucm_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Enterprise datastore id",
		},
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

func LocationSchWithoutIbx() map[string]*schema.Schema {
	schemaMap := LocationSch()
	delete(schemaMap, "ibx")
	return schemaMap
}

func ProjectSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Project Id",
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

func OperationSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"provider_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection provider readiness status",
		},
		"equinix_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Connection status",
		},
		"errors": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Errors occurred",
			Elem: &schema.Resource{
				Schema: ErrorSch(),
			},
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

func PortSch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Equinix-assigned Port identifier",
		},
		"href": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Unique Resource Identifier",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port name",
		},
		"redundancy": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "Redundancy Information",
			Elem: &schema.Resource{
				Schema: PortRedundancySch(),
			},
		},
	}
}

func PortRedundancySch() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Access point redundancy",
		},
		"group": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Port redundancy group",
		},
		"priority": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Priority type - Primary or Secondary",
		},
	}
}