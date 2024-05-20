package equinix

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"

	"github.com/equinix/ecx-go/v2"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ecxL2ConnectionSchemaNames = map[string]string{
	"UUID":                "uuid",
	"Name":                "name",
	"ProfileUUID":         "profile_uuid",
	"Speed":               "speed",
	"SpeedUnit":           "speed_unit",
	"Status":              "status",
	"ProviderStatus":      "provider_status",
	"Notifications":       "notifications",
	"PurchaseOrderNumber": "purchase_order_number",
	"PortUUID":            "port_uuid",
	"DeviceUUID":          "device_uuid",
	"DeviceInterfaceID":   "device_interface_id",
	"VlanSTag":            "vlan_stag",
	"VlanCTag":            "vlan_ctag",
	"NamedTag":            "named_tag",
	"AdditionalInfo":      "additional_info",
	"ZSidePortUUID":       "zside_port_uuid",
	"ZSideServiceToken":   "zside_service_token",
	"ZSideVlanSTag":       "zside_vlan_stag",
	"ZSideVlanCTag":       "zside_vlan_ctag",
	"SellerRegion":        "seller_region",
	"SellerMetroCode":     "seller_metro_code",
	"AuthorizationKey":    "authorization_key",
	"RedundantUUID":       "redundant_uuid",
	"RedundancyType":      "redundancy_type",
	"RedundancyGroup":     "redundancy_group",
	"SecondaryConnection": "secondary_connection",
	"Actions":             "actions",
	"ServiceToken":        "service_token",
	"VendorToken":         "vendor_token",
}

var ecxL2ConnectionDescriptions = map[string]string{
	"UUID":                "Unique identifier of the connection",
	"Name":                "Connection name. An alpha-numeric 24 characters string which can include only hyphens and underscores",
	"ProfileUUID":         "Unique identifier of the service provider's service profile",
	"Speed":               "Speed/Bandwidth to be allocated to the connection",
	"SpeedUnit":           "Unit of the speed/bandwidth to be allocated to the connection",
	"Status":              "Connection provisioning status on Equinix Fabric side",
	"ProviderStatus":      "Connection provisioning status on service provider's side",
	"Notifications":       "A list of email addresses used for sending connection update notifications",
	"PurchaseOrderNumber": "Connection's purchase order number to reflect on the invoice",
	"PortUUID":            "Unique identifier of the buyer's port from which the connection would originate",
	"DeviceUUID":          "Unique identifier of the Network Edge virtual device from which the connection would originate",
	"DeviceInterfaceID":   "Identifier of network interface on a given device, used for a connection. If not specified then first available interface will be selected",
	"VlanSTag":            "S-Tag/Outer-Tag of the connection, a numeric character ranging from 2 - 4094",
	"VlanCTag":            "C-Tag/Inner-Tag of the connection, a numeric character ranging from 2 - 4094",
	"NamedTag":            "The type of peering to set up in case when connecting to Azure Express Route. One of PRIVATE, MICROSOFT, MANUAL, PUBLIC (MANUAL and PUBLIC are deprecated and not available for new connections)",
	"AdditionalInfo":      "One or more additional information key-value objects",
	"ZSidePortUUID":       "Unique identifier of the port on the remote side (z-side)",
	"ZSideServiceToken":   "Unique Equinix Fabric key given by a provider that grants you authorization to enable connectivity to a shared multi-tenant port (z-side)",
	"ZSideVlanSTag":       "S-Tag/Outer-Tag of the connection on the remote side (z-side)",
	"ZSideVlanCTag":       "C-Tag/Inner-Tag of the connection on the remote side (z-side)",
	"SellerRegion":        "The region in which the seller port resides",
	"SellerMetroCode":     "The metro code that denotes the connection's remote side (z-side)",
	"AuthorizationKey":    "Text field used to authorize connection on the provider side. Value depends on a provider service profile used for connection",
	"RedundantUUID":       "Unique identifier of the redundant connection, applicable for HA connections",
	"RedundancyType":      "Connection redundancy type, applicable for HA connections. Either primary or secondary",
	"RedundancyGroup":     "Unique identifier of group containing a primary and secondary connection",
	"SecondaryConnection": "Definition of secondary connection for redundant, HA connectivity",
	"Actions":             "One or more pending actions to complete connection provisioning",
	"ServiceToken":        "Unique Equinix Fabric key given by a provider that grants you authorization to enable connectivity from a shared multi-tenant port (a-side)",
	"VendorToken":         "The Equinix Fabric Token the connection was created with. Applicable if the connection was created with a ServiceToken (a-side) or ZSideServiceToken (z-side)",
}

var ecxL2ConnectionAdditionalInfoSchemaNames = map[string]string{
	"Name":  "name",
	"Value": "value",
}

var ecxL2ConnectionAdditionalInfoDescriptions = map[string]string{
	"Name":  "Additional information key",
	"Value": "Additional information value",
}

var ecxL2ConnectionActionsSchemaNames = map[string]string{
	"Type":         "type",
	"OperationID":  "operation_id",
	"Message":      "message",
	"RequiredData": "required_data",
}

var ecxL2ConnectionActionsDescriptions = map[string]string{
	"Type":         "Action type",
	"OperationID":  "Action identifier",
	"Message":      "Action information",
	"RequiredData": "Action list of required data",
}

var ecxL2ConnectionActionDataSchemaNames = map[string]string{
	"Key":               "key",
	"Label":             "label",
	"Value":             "value",
	"IsEditable":        "editable",
	"ValidationPattern": "validation_pattern",
}

var ecxL2ConnectionActionDataDescriptions = map[string]string{
	"Key":               "Action data key",
	"Label":             "Action data label",
	"Value":             "Action data value",
	"IsEditable":        "Action data is editable",
	"ValidationPattern": "Action data pattern",
}

type (
	getL2Connection func(uuid string) (*ecx.L2Connection, error)
)

func resourceECXL2Connection() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated. End of Life will be June 30th, 2024. Use equinix_fabric_connection instead.",
		CreateContext:      resourceECXL2ConnectionCreate,
		ReadContext:        resourceECXL2ConnectionRead,
		UpdateContext:      resourceECXL2ConnectionUpdate,
		DeleteContext:      resourceECXL2ConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// The expected ID to import redundant connections is '(primaryID):(secondaryID)', e.g.,
				//   terraform import equinix_ecx_l2_connection.example 1111:2222
				ids := strings.Split(d.Id(), ":")
				d.SetId(ids[0])
				if len(ids) > 1 {
					d.Set(ecxL2ConnectionSchemaNames["RedundantUUID"], ids[1])
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: createECXL2ConnectionResourceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Description: "Resource allows creation and management of Equinix Fabric	layer 2 connections",
	}
}

func createECXL2ConnectionResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxL2ConnectionSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["UUID"],
		},
		ecxL2ConnectionSchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 24),
			Description:  ecxL2ConnectionDescriptions["Name"],
		},
		ecxL2ConnectionSchemaNames["ProfileUUID"]: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			AtLeastOneOf: []string{
				ecxL2ConnectionSchemaNames["ProfileUUID"],
				ecxL2ConnectionSchemaNames["ZSidePortUUID"],
				ecxL2ConnectionSchemaNames["ZSideServiceToken"],
			},
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionDescriptions["ProfileUUID"],
		},
		ecxL2ConnectionSchemaNames["Speed"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  ecxL2ConnectionDescriptions["Speed"],
		},
		ecxL2ConnectionSchemaNames["SpeedUnit"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"MB", "GB"}, false),
			Description:  ecxL2ConnectionDescriptions["SpeedUnit"],
		},
		ecxL2ConnectionSchemaNames["Status"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["Status"],
		},
		ecxL2ConnectionSchemaNames["ProviderStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["ProviderStatus"],
		},
		ecxL2ConnectionSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: equinix_validation.StringIsEmailAddress,
			},
			Description: ecxL2ConnectionDescriptions["Notifications"],
		},
		ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 30),
			Description:  ecxL2ConnectionDescriptions["PurchaseOrderNumber"],
		},
		ecxL2ConnectionSchemaNames["PortUUID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{
				ecxL2ConnectionSchemaNames["PortUUID"],
				ecxL2ConnectionSchemaNames["DeviceUUID"],
				ecxL2ConnectionSchemaNames["ServiceToken"],
			},
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"], ecxL2ConnectionSchemaNames["ServiceToken"]},
			Description:   ecxL2ConnectionDescriptions["PortUUID"],
		},
		ecxL2ConnectionSchemaNames["DeviceUUID"]: {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["PortUUID"], ecxL2ConnectionSchemaNames["ServiceToken"]},
			Description:   ecxL2ConnectionDescriptions["DeviceUUID"],
		},
		ecxL2ConnectionSchemaNames["DeviceInterfaceID"]: {
			Type:          schema.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["PortUUID"], ecxL2ConnectionSchemaNames["ServiceToken"]},
			Description:   ecxL2ConnectionDescriptions["DeviceInterfaceID"],
		},
		ecxL2ConnectionSchemaNames["VlanSTag"]: {
			Type:          schema.TypeInt,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ValidateFunc:  validation.IntBetween(2, 4092),
			RequiredWith:  []string{ecxL2ConnectionSchemaNames["PortUUID"]},
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"], ecxL2ConnectionSchemaNames["ServiceToken"]},
			Description:   ecxL2ConnectionDescriptions["VlanSTag"],
		},
		ecxL2ConnectionSchemaNames["VlanCTag"]: {
			Type:          schema.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.IntBetween(2, 4092),
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["DeviceUUID"]},
			Description:   ecxL2ConnectionDescriptions["VlanCTag"],
		},
		ecxL2ConnectionSchemaNames["NamedTag"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "MICROSOFT", "MANUAL", "PUBLIC"}, true),
			Description:  ecxL2ConnectionDescriptions["NamedTag"],
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return strings.EqualFold(old, new)
			},
		},
		ecxL2ConnectionSchemaNames["AdditionalInfo"]: {
			Type:        schema.TypeSet,
			Optional:    true,
			ForceNew:    true,
			MinItems:    1,
			Description: ecxL2ConnectionDescriptions["AdditionalInfo"],
			Elem: &schema.Resource{
				Schema: createECXL2ConnectionAdditionalInfoResourceSchema(),
			},
		},
		ecxL2ConnectionSchemaNames["ZSidePortUUID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionDescriptions["ZSidePortUUID"],
		},
		ecxL2ConnectionSchemaNames["ZSideVlanSTag"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(2, 4092),
			Description:  ecxL2ConnectionDescriptions["ZSideVlanSTag"],
		},
		ecxL2ConnectionSchemaNames["ZSideVlanCTag"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(2, 4092),
			Description:  ecxL2ConnectionDescriptions["ZSideVlanCTag"],
		},
		ecxL2ConnectionSchemaNames["SellerRegion"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionDescriptions["SellerRegion"],
		},
		ecxL2ConnectionSchemaNames["SellerMetroCode"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: equinix_validation.StringIsMetroCode,
			Description:  ecxL2ConnectionDescriptions["SellerMetroCode"],
		},
		ecxL2ConnectionSchemaNames["AuthorizationKey"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionDescriptions["AuthorizationKey"],
		},
		ecxL2ConnectionSchemaNames["RedundantUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["RedundantUUID"],
		},
		ecxL2ConnectionSchemaNames["RedundancyType"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["RedundancyType"],
		},
		ecxL2ConnectionSchemaNames["RedundancyGroup"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["RedundancyGroup"],
		},
		ecxL2ConnectionSchemaNames["Actions"]: {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["Actions"],
			Elem: &schema.Resource{
				Schema: createECXL2ConnectionActionsSchema(),
			},
		},
		ecxL2ConnectionSchemaNames["SecondaryConnection"]: {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: ecxL2ConnectionDescriptions["SecondaryConnection"],
			Elem: &schema.Resource{
				Schema: createECXL2ConnectionSecondaryResourceSchema(),
			},
		},
		ecxL2ConnectionSchemaNames["ServiceToken"]: {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["PortUUID"], ecxL2ConnectionSchemaNames["DeviceUUID"]},
			Description:   ecxL2ConnectionDescriptions["ServiceToken"],
		},
		ecxL2ConnectionSchemaNames["ZSideServiceToken"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ConflictsWith: []string{
				ecxL2ConnectionSchemaNames["ServiceToken"],
				ecxL2ConnectionSchemaNames["ProfileUUID"],
				ecxL2ConnectionSchemaNames["ZSidePortUUID"],
				ecxL2ConnectionSchemaNames["AuthorizationKey"],
				ecxL2ConnectionSchemaNames["SecondaryConnection"],
			},
			Description: ecxL2ConnectionDescriptions["ZSideServiceToken"],
		},
		ecxL2ConnectionSchemaNames["VendorToken"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["VendorToken"],
		},
	}
}

func createECXL2ConnectionAdditionalInfoResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxL2ConnectionAdditionalInfoSchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionAdditionalInfoDescriptions["Name"],
		},
		ecxL2ConnectionAdditionalInfoSchemaNames["Value"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionAdditionalInfoDescriptions["Value"],
		},
	}
}

func createECXL2ConnectionActionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxL2ConnectionActionsSchemaNames["Type"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionActionsDescriptions["Type"],
		},
		ecxL2ConnectionActionsSchemaNames["OperationID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionActionsDescriptions["OperationID"],
		},
		ecxL2ConnectionActionsSchemaNames["Message"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionActionsDescriptions["Message"],
		},
		ecxL2ConnectionActionsSchemaNames["RequiredData"]: {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: ecxL2ConnectionActionsDescriptions["RequiredData"],
			Elem: &schema.Resource{
				Schema: createECXL2ConnectionActionsRequiredDataSchema(),
			},
		},
	}
}

func createECXL2ConnectionSecondaryResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxL2ConnectionSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["UUID"],
		},
		ecxL2ConnectionSchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 24),
			Description:  ecxL2ConnectionDescriptions["Name"],
		},
		ecxL2ConnectionSchemaNames["ProfileUUID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionDescriptions["ProfileUUID"],
		},
		ecxL2ConnectionSchemaNames["Speed"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  ecxL2ConnectionDescriptions["Speed"],
		},
		ecxL2ConnectionSchemaNames["SpeedUnit"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"MB", "GB"}, false),
			RequiredWith: []string{ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["Speed"]},
			Description:  ecxL2ConnectionDescriptions["SpeedUnit"],
		},
		ecxL2ConnectionSchemaNames["Status"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["Status"],
		},
		ecxL2ConnectionSchemaNames["ProviderStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["ProviderStatus"],
		},
		ecxL2ConnectionSchemaNames["PortUUID"]: {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["PortUUID"],
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["DeviceUUID"],
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["ServiceToken"],
			},
			ConflictsWith: []string{
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["DeviceUUID"],
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["ServiceToken"],
			},
			Description: ecxL2ConnectionDescriptions["PortUUID"],
		},
		ecxL2ConnectionSchemaNames["DeviceUUID"]: {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,

			Description: ecxL2ConnectionDescriptions["DeviceUUID"],
		},
		ecxL2ConnectionSchemaNames["DeviceInterfaceID"]: {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ConflictsWith: []string{
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["PortUUID"],
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["ServiceToken"],
			},
			Description: ecxL2ConnectionDescriptions["DeviceInterfaceID"],
		},
		ecxL2ConnectionSchemaNames["VlanSTag"]: {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(2, 4092),
			RequiredWith: []string{ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["PortUUID"]},
			ConflictsWith: []string{
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["DeviceUUID"],
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["ServiceToken"],
			},
			Description: ecxL2ConnectionDescriptions["VlanSTag"],
		},
		ecxL2ConnectionSchemaNames["VlanCTag"]: {
			Type:          schema.TypeInt,
			ForceNew:      true,
			Optional:      true,
			ValidateFunc:  validation.IntBetween(2, 4092),
			ConflictsWith: []string{ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["DeviceUUID"]},
			Description:   ecxL2ConnectionDescriptions["VlanCTag"],
		},
		ecxL2ConnectionSchemaNames["ZSidePortUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["ZSidePortUUID"],
		},
		ecxL2ConnectionSchemaNames["ZSideVlanSTag"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["ZSideVlanSTag"],
		},
		ecxL2ConnectionSchemaNames["ZSideVlanCTag"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["ZSideVlanCTag"],
		},
		ecxL2ConnectionSchemaNames["SellerRegion"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionDescriptions["SellerRegion"],
		},
		ecxL2ConnectionSchemaNames["SellerMetroCode"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: equinix_validation.StringIsMetroCode,
			Description:  ecxL2ConnectionDescriptions["SellerMetroCode"],
		},
		ecxL2ConnectionSchemaNames["AuthorizationKey"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  ecxL2ConnectionDescriptions["AuthorizationKey"],
		},
		ecxL2ConnectionSchemaNames["RedundantUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["RedundantUUID"],
			Deprecated:  "SecondaryConnection.0.RedundantUUID will not be returned. Use UUID instead",
		},
		ecxL2ConnectionSchemaNames["RedundancyType"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["RedundancyType"],
		},
		ecxL2ConnectionSchemaNames["RedundancyGroup"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["RedundancyGroup"],
		},
		ecxL2ConnectionSchemaNames["Actions"]: {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["Actions"],
			Elem: &schema.Resource{
				Schema: createECXL2ConnectionActionsSchema(),
			},
		},
		ecxL2ConnectionSchemaNames["ServiceToken"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ConflictsWith: []string{
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["PortUUID"],
				ecxL2ConnectionSchemaNames["SecondaryConnection"] + ".0." + ecxL2ConnectionSchemaNames["DeviceUUID"],
			},
			Description: ecxL2ConnectionDescriptions["ServiceToken"],
		},
		ecxL2ConnectionSchemaNames["VendorToken"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionDescriptions["VendorToken"],
		},
	}
}

func createECXL2ConnectionActionsRequiredDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxL2ConnectionActionDataSchemaNames["Key"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionActionDataDescriptions["Key"],
		},
		ecxL2ConnectionActionDataSchemaNames["Label"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionActionDataDescriptions["Label"],
		},
		ecxL2ConnectionActionDataSchemaNames["Value"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionActionDataDescriptions["Value"],
		},
		ecxL2ConnectionActionDataSchemaNames["IsEditable"]: {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: ecxL2ConnectionActionDataDescriptions["IsEditable"],
		},
		ecxL2ConnectionActionDataSchemaNames["ValidationPattern"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxL2ConnectionActionDataDescriptions["ValidationPattern"],
		},
	}
}

func resourceECXL2ConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ecx
	m.(*config.Config).AddModuleToECXUserAgent(&client, d)

	var diags diag.Diagnostics
	primary, secondary := createECXL2Connections(d)
	var primaryID, secondaryID *string
	var err error
	if secondary != nil {
		primaryID, secondaryID, err = client.CreateL2RedundantConnection(*primary, *secondary)
	} else {
		primaryID, err = client.CreateL2Connection(*primary)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ecx.StringValue(primaryID))
	waitConfigs := []*retry.StateChangeConf{
		createConnectionStatusProvisioningWaitConfiguration(client.GetL2Connection, d.Id(), 5*time.Second, d.Timeout(schema.TimeoutCreate)),
	}
	if ecx.StringValue(secondaryID) != "" {
		d.Set(ecxL2ConnectionSchemaNames["RedundantUUID"], secondaryID)
		waitConfigs = append(waitConfigs,
			createConnectionStatusProvisioningWaitConfiguration(client.GetL2Connection, ecx.StringValue(secondaryID), 2*time.Second, d.Timeout(schema.TimeoutCreate)),
		)
	}
	for _, config := range waitConfigs {
		if config == nil {
			continue
		}
		if _, err := config.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error waiting for connection (%s) to be created: %s", d.Id(), err)
		}
	}
	diags = append(diags, resourceECXL2ConnectionRead(ctx, d, m)...)
	return diags
}

func resourceECXL2ConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ecx
	m.(*config.Config).AddModuleToECXUserAgent(&client, d)
	var diags diag.Diagnostics
	var err error
	var primary *ecx.L2Connection
	var secondary *ecx.L2Connection

	primary, err = client.GetL2Connection(d.Id())
	if err != nil {
		return diag.Errorf("cannot fetch primary connection due to %v", err)
	}
	if slices.Contains([]string{
		ecx.ConnectionStatusPendingDelete,
		ecx.ConnectionStatusDeprovisioning,
		ecx.ConnectionStatusDeprovisioned,
		ecx.ConnectionStatusDeleted,
	}, ecx.StringValue(primary.Status)) {
		d.SetId("")
		return diags
	}

	// RedundantUUID value is set in CreateContext/Importer functions
	// Implementing a l2_connection datasource will require search for secondary connection before using
	// resourceECXL2ConnectionRead or explicitly request the names or identifiers of each connection
	if redID, ok := d.GetOk(ecxL2ConnectionSchemaNames["RedundantUUID"]); ok {
		secondary, err = client.GetL2Connection(redID.(string))
		if err != nil {
			return diag.Errorf("cannot fetch secondary connection due to %v", err)
		}
		if ecx.StringValue(primary.RedundancyGroup) != ecx.StringValue(secondary.RedundancyGroup) || !strings.EqualFold(ecx.StringValue(secondary.RedundancyType), "secondary") {
			return diag.Errorf("connection '%s' (%s) was found but is not the redundant connection for '%s' (%s)",
				ecx.StringValue(secondary.Name),
				ecx.StringValue(secondary.UUID),
				ecx.StringValue(primary.Name),
				ecx.StringValue(primary.UUID))
		}
	}

	if err := updateECXL2ConnectionResource(primary, secondary, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceECXL2ConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ecx
	m.(*config.Config).AddModuleToECXUserAgent(&client, d)
	var diags diag.Diagnostics
	supportedChanges := []string{
		ecxL2ConnectionSchemaNames["Name"],
		ecxL2ConnectionSchemaNames["Speed"],
		ecxL2ConnectionSchemaNames["SpeedUnit"],
	}
	primaryChanges := equinix_schema.GetResourceDataChangedKeys(supportedChanges, d)
	primaryUpdateReq := client.NewL2ConnectionUpdateRequest(d.Id())
	if err := fillFabricL2ConnectionUpdateRequest(primaryUpdateReq, primaryChanges).Execute(); err != nil {
		return diag.FromErr(err)
	}
	if redID, ok := d.GetOk(ecxL2ConnectionSchemaNames["RedundantUUID"]); ok {
		secondaryChanges := equinix_schema.GetResourceDataListElementChanges(supportedChanges, ecxL2ConnectionSchemaNames["SecondaryConnection"], 0, d)
		secondaryUpdateReq := client.NewL2ConnectionUpdateRequest(redID.(string))
		if err := fillFabricL2ConnectionUpdateRequest(secondaryUpdateReq, secondaryChanges).Execute(); err != nil {
			return diag.FromErr(err)
		}
	}
	diags = append(diags, resourceECXL2ConnectionRead(ctx, d, m)...)
	return diags
}

func resourceECXL2ConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ecx
	m.(*config.Config).AddModuleToECXUserAgent(&client, d)

	var diags diag.Diagnostics
	if err := client.DeleteL2Connection(d.Id()); err != nil {
		restErr, ok := err.(rest.Error)
		if ok {
			// IC-LAYER2-4021 = Connection already deleted
			if equinix_errors.HasApplicationErrorCode(restErr.ApplicationErrors, "IC-LAYER2-4021") {
				return diags
			}
		}
		return diag.FromErr(err)
	}
	waitConfigs := []*retry.StateChangeConf{
		createConnectionStatusDeleteWaitConfiguration(client.GetL2Connection, d.Id(), 5*time.Second, d.Timeout(schema.TimeoutDelete)),
	}
	if redID, ok := d.GetOk(ecxL2ConnectionSchemaNames["RedundantUUID"]); ok {
		if err := client.DeleteL2Connection(redID.(string)); err != nil {
			restErr, ok := err.(rest.Error)
			if ok {
				// IC-LAYER2-4021 = Connection already deleted
				if equinix_errors.HasApplicationErrorCode(restErr.ApplicationErrors, "IC-LAYER2-4021") {
					return diags
				}
			}
			return diag.FromErr(err)
		}
		waitConfigs = append(waitConfigs,
			createConnectionStatusDeleteWaitConfiguration(client.GetL2Connection, redID.(string), 2*time.Second, d.Timeout(schema.TimeoutDelete)),
		)
	}
	for _, config := range waitConfigs {
		if _, err := config.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error waiting for connection (%s) to be removed: %s", d.Id(), err)
		}
	}
	return diags
}

func createECXL2Connections(d *schema.ResourceData) (*ecx.L2Connection, *ecx.L2Connection) {
	var primary, secondary *ecx.L2Connection
	primary = &ecx.L2Connection{}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["Name"]); ok {
		primary.Name = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ProfileUUID"]); ok {
		primary.ProfileUUID = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["Speed"]); ok {
		primary.Speed = ecx.Int(v.(int))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SpeedUnit"]); ok {
		primary.SpeedUnit = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["Notifications"]); ok {
		primary.Notifications = converters.SetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"]); ok {
		primary.PurchaseOrderNumber = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["PortUUID"]); ok {
		primary.PortUUID = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["DeviceUUID"]); ok {
		primary.DeviceUUID = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["DeviceInterfaceID"]); ok {
		primary.DeviceInterfaceID = ecx.Int(v.(int))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["VlanSTag"]); ok {
		primary.VlanSTag = ecx.Int(v.(int))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["VlanCTag"]); ok {
		primary.VlanCTag = ecx.Int(v.(int))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["NamedTag"]); ok {
		primary.NamedTag = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["AdditionalInfo"]); ok {
		primary.AdditionalInfo = expandECXL2ConnectionAdditionalInfo(v.(*schema.Set))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSidePortUUID"]); ok {
		primary.ZSidePortUUID = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSideVlanSTag"]); ok {
		primary.ZSideVlanSTag = ecx.Int(v.(int))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSideVlanCTag"]); ok {
		primary.ZSideVlanCTag = ecx.Int(v.(int))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SellerRegion"]); ok {
		primary.SellerRegion = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SellerMetroCode"]); ok {
		primary.SellerMetroCode = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ServiceToken"]); ok {
		primary.ServiceToken = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSideServiceToken"]); ok {
		primary.ZSideServiceToken = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["AuthorizationKey"]); ok {
		primary.AuthorizationKey = ecx.String(v.(string))
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SecondaryConnection"]); ok {
		secondary = expandECXL2ConnectionSecondary(v.([]interface{}))
	}
	return primary, secondary
}

func updateECXL2ConnectionResource(primary *ecx.L2Connection, secondary *ecx.L2Connection, d *schema.ResourceData) error {
	if err := d.Set(ecxL2ConnectionSchemaNames["UUID"], primary.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Name"], primary.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ProfileUUID"], primary.ProfileUUID); err != nil {
		return fmt.Errorf("error reading ProfileUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Speed"], primary.Speed); err != nil {
		return fmt.Errorf("error reading Speed: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SpeedUnit"], primary.SpeedUnit); err != nil {
		return fmt.Errorf("error reading SpeedUnit: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Status"], primary.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ProviderStatus"], primary.ProviderStatus); err != nil {
		return fmt.Errorf("error reading ProviderStatus: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Notifications"], primary.Notifications); err != nil {
		return fmt.Errorf("error reading Notifications: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["PurchaseOrderNumber"], primary.PurchaseOrderNumber); err != nil {
		return fmt.Errorf("error reading PurchaseOrderNumber: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["PortUUID"], primary.PortUUID); err != nil {
		return fmt.Errorf("error reading PortUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["DeviceUUID"], primary.DeviceUUID); err != nil {
		return fmt.Errorf("error reading DeviceUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["VlanSTag"], primary.VlanSTag); err != nil {
		return fmt.Errorf("error reading VlanSTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["VlanCTag"], primary.VlanCTag); err != nil {
		return fmt.Errorf("error reading VlanCTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["NamedTag"], primary.NamedTag); err != nil {
		return fmt.Errorf("error reading NamedTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["AdditionalInfo"], flattenECXL2ConnectionAdditionalInfo(primary.AdditionalInfo)); err != nil {
		return fmt.Errorf("error reading AdditionalInfo: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSidePortUUID"], primary.ZSidePortUUID); err != nil {
		return fmt.Errorf("error reading ZSidePortUUID: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSideVlanSTag"], primary.ZSideVlanSTag); err != nil {
		return fmt.Errorf("error reading ZSideVlanSTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["ZSideVlanCTag"], primary.ZSideVlanCTag); err != nil {
		return fmt.Errorf("error reading ZSideVlanCTag: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SellerRegion"], primary.SellerRegion); err != nil {
		return fmt.Errorf("error reading SellerRegion: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["SellerMetroCode"], primary.SellerMetroCode); err != nil {
		return fmt.Errorf("error reading SellerMetroCode: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["AuthorizationKey"], primary.AuthorizationKey); err != nil {
		return fmt.Errorf("error reading AuthorizationKey: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["RedundancyType"], primary.RedundancyType); err != nil {
		return fmt.Errorf("error reading RedundancyType: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["RedundancyGroup"], primary.RedundancyGroup); err != nil {
		return fmt.Errorf("error reading RedundancyGroup: %s", err)
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["VendorToken"], primary.VendorToken); err != nil {
		return fmt.Errorf("error reading VendorToken: %s", err)
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ServiceToken"]); ok {
		if ecx.StringValue(primary.VendorToken) != v.(string) {
			if err := d.Set(ecxL2ConnectionSchemaNames["ServiceToken"], primary.VendorToken); err != nil {
				return fmt.Errorf("error reading ServiceToken: %s", err)
			}
		}
	}
	if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["ZSideServiceToken"]); ok {
		if ecx.StringValue(primary.VendorToken) != v.(string) {
			if err := d.Set(ecxL2ConnectionSchemaNames["ZSideServiceToken"], primary.VendorToken); err != nil {
				return fmt.Errorf("error reading ZSideServiceToken: %s", err)
			}
		}
	}
	if err := d.Set(ecxL2ConnectionSchemaNames["Actions"], flattenECXL2ConnectionActions(primary.Actions)); err != nil {
		return fmt.Errorf("error reading Actions: %s", err)
	}
	if secondary != nil {
		var prevSecondary *ecx.L2Connection
		if v, ok := d.GetOk(ecxL2ConnectionSchemaNames["SecondaryConnection"]); ok {
			prevSecondary = expandECXL2ConnectionSecondary(v.([]interface{}))
		}
		if err := d.Set(ecxL2ConnectionSchemaNames["SecondaryConnection"], flattenECXL2ConnectionSecondary(prevSecondary, secondary)); err != nil {
			return fmt.Errorf("error reading SecondaryConnection: %s", err)
		}
	}
	return nil
}

func flattenECXL2ConnectionSecondary(previous, conn *ecx.L2Connection) interface{} {
	transformed := make(map[string]interface{})
	transformed[ecxL2ConnectionSchemaNames["UUID"]] = conn.UUID
	transformed[ecxL2ConnectionSchemaNames["Name"]] = conn.Name
	transformed[ecxL2ConnectionSchemaNames["ProfileUUID"]] = conn.ProfileUUID
	transformed[ecxL2ConnectionSchemaNames["Speed"]] = conn.Speed
	transformed[ecxL2ConnectionSchemaNames["SpeedUnit"]] = conn.SpeedUnit
	transformed[ecxL2ConnectionSchemaNames["Status"]] = conn.Status
	transformed[ecxL2ConnectionSchemaNames["ProviderStatus"]] = conn.ProviderStatus
	transformed[ecxL2ConnectionSchemaNames["PortUUID"]] = conn.PortUUID
	transformed[ecxL2ConnectionSchemaNames["DeviceUUID"]] = conn.DeviceUUID
	transformed[ecxL2ConnectionSchemaNames["DeviceInterfaceID"]] = conn.DeviceInterfaceID
	transformed[ecxL2ConnectionSchemaNames["VlanSTag"]] = conn.VlanSTag
	transformed[ecxL2ConnectionSchemaNames["VlanCTag"]] = conn.VlanCTag
	transformed[ecxL2ConnectionSchemaNames["ZSidePortUUID"]] = conn.ZSidePortUUID
	transformed[ecxL2ConnectionSchemaNames["ZSideVlanSTag"]] = conn.ZSideVlanSTag
	transformed[ecxL2ConnectionSchemaNames["ZSideVlanCTag"]] = conn.ZSideVlanCTag
	transformed[ecxL2ConnectionSchemaNames["SellerRegion"]] = conn.SellerRegion
	transformed[ecxL2ConnectionSchemaNames["SellerMetroCode"]] = conn.SellerMetroCode
	transformed[ecxL2ConnectionSchemaNames["AuthorizationKey"]] = conn.AuthorizationKey
	transformed[ecxL2ConnectionSchemaNames["RedundancyType"]] = conn.RedundancyType
	transformed[ecxL2ConnectionSchemaNames["RedundancyGroup"]] = conn.RedundancyGroup
	transformed[ecxL2ConnectionSchemaNames["Actions"]] = flattenECXL2ConnectionActions(conn.Actions)
	transformed[ecxL2ConnectionSchemaNames["VendorToken"]] = conn.VendorToken
	if previous != nil {
		transformed[ecxL2ConnectionSchemaNames["DeviceInterfaceID"]] = previous.DeviceInterfaceID
		transformed[ecxL2ConnectionSchemaNames["ServiceToken"]] = previous.ServiceToken
		prevSToken := ecx.StringValue(previous.ServiceToken)
		if prevSToken != "" && ecx.StringValue(conn.VendorToken) != prevSToken {
			transformed[ecxL2ConnectionSchemaNames["ServiceToken"]] = conn.VendorToken
		}
	}
	return []interface{}{transformed}
}

func expandECXL2ConnectionSecondary(conns []interface{}) *ecx.L2Connection {
	if len(conns) < 1 {
		log.Printf("[WARN] resource_ecx_l2_connection expanding empty secondary connection collection")
		return nil
	}
	conn := conns[0].(map[string]interface{})
	transformed := ecx.L2Connection{}
	if v, ok := conn[ecxL2ConnectionSchemaNames["Name"]]; ok {
		transformed.Name = ecx.String(v.(string))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["ProfileUUID"]]; ok && !isEmpty(v) {
		transformed.ProfileUUID = ecx.String(v.(string))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["Speed"]]; ok && !isEmpty(v) {
		transformed.Speed = ecx.Int(v.(int))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["SpeedUnit"]]; ok && !isEmpty(v) {
		transformed.SpeedUnit = ecx.String(v.(string))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["PortUUID"]]; ok && !isEmpty(v) {
		transformed.PortUUID = ecx.String(v.(string))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["DeviceUUID"]]; ok && !isEmpty(v) {
		transformed.DeviceUUID = ecx.String(v.(string))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["DeviceInterfaceID"]]; ok && !isEmpty(v) {
		transformed.DeviceInterfaceID = ecx.Int(v.(int))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["VlanSTag"]]; ok && !isEmpty(v) {
		transformed.VlanSTag = ecx.Int(v.(int))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["VlanCTag"]]; ok && !isEmpty(v) {
		transformed.VlanCTag = ecx.Int(v.(int))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["SellerRegion"]]; ok && !isEmpty(v) {
		transformed.SellerRegion = ecx.String(v.(string))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["SellerMetroCode"]]; ok && !isEmpty(v) {
		transformed.SellerMetroCode = ecx.String(v.(string))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["AuthorizationKey"]]; ok && !isEmpty(v) {
		transformed.AuthorizationKey = ecx.String(v.(string))
	}
	if v, ok := conn[ecxL2ConnectionSchemaNames["ServiceToken"]]; ok && !isEmpty(v) {
		transformed.ServiceToken = ecx.String(v.(string))
	}
	return &transformed
}

func flattenECXL2ConnectionAdditionalInfo(infos []ecx.L2ConnectionAdditionalInfo) interface{} {
	transformed := make([]interface{}, 0, len(infos))
	for _, info := range infos {
		transformed = append(transformed, map[string]interface{}{
			ecxL2ConnectionAdditionalInfoSchemaNames["Name"]:  info.Name,
			ecxL2ConnectionAdditionalInfoSchemaNames["Value"]: info.Value,
		})
	}
	return transformed
}

func flattenECXL2ConnectionActions(actions []ecx.L2ConnectionAction) []interface{} {
	transformed := make([]interface{}, 0, len(actions))
	for _, action := range actions {
		transformedAction := make(map[string]interface{})
		transformedAction[ecxL2ConnectionActionsSchemaNames["Type"]] = action.Type
		transformedAction[ecxL2ConnectionActionsSchemaNames["OperationID"]] = action.OperationID
		transformedAction[ecxL2ConnectionActionsSchemaNames["Message"]] = action.Message
		if v := action.RequiredData; v != nil {
			transformedAction[ecxL2ConnectionActionsSchemaNames["RequiredData"]] = flattenECXL2ConnectionActionData(v)
		}
		transformed = append(transformed, transformedAction)
	}

	return transformed
}

func flattenECXL2ConnectionActionData(actionData []ecx.L2ConnectionActionData) []interface{} {
	transformed := make([]interface{}, 0, len(actionData))
	for _, data := range actionData {
		transformed = append(transformed, map[string]interface{}{
			ecxL2ConnectionActionDataSchemaNames["Key"]:               data.Key,
			ecxL2ConnectionActionDataSchemaNames["Label"]:             data.Label,
			ecxL2ConnectionActionDataSchemaNames["Value"]:             data.Value,
			ecxL2ConnectionActionDataSchemaNames["IsEditable"]:        data.IsEditable,
			ecxL2ConnectionActionDataSchemaNames["ValidationPattern"]: data.ValidationPattern,
		})
	}

	return transformed
}

func expandECXL2ConnectionAdditionalInfo(infos *schema.Set) []ecx.L2ConnectionAdditionalInfo {
	transformed := make([]ecx.L2ConnectionAdditionalInfo, 0, infos.Len())
	for _, info := range infos.List() {
		infoMap := info.(map[string]interface{})
		transformed = append(transformed, ecx.L2ConnectionAdditionalInfo{
			Name:  ecx.String(infoMap[ecxL2ConnectionAdditionalInfoSchemaNames["Name"]].(string)),
			Value: ecx.String(infoMap[ecxL2ConnectionAdditionalInfoSchemaNames["Value"]].(string)),
		})
	}
	return transformed
}

func fillFabricL2ConnectionUpdateRequest(updateReq ecx.L2ConnectionUpdateRequest, changes map[string]interface{}) ecx.L2ConnectionUpdateRequest {
	for change, changeValue := range changes {
		switch change {
		case ecxL2ConnectionSchemaNames["Name"]:
			updateReq.WithName(changeValue.(string))
		case ecxL2ConnectionSchemaNames["Speed"]:
			updateReq.WithSpeed(changeValue.(int))
		case ecxL2ConnectionSchemaNames["SpeedUnit"]:
			updateReq.WithSpeedUnit(changeValue.(string))
		}
	}
	return updateReq
}

func createConnectionStatusProvisioningWaitConfiguration(fetchFunc getL2Connection, id string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	pending := []string{
		ecx.ConnectionStatusProvisioning,
		ecx.ConnectionStatusPendingAutoApproval,
	}
	target := []string{
		ecx.ConnectionStatusProvisioned,
		ecx.ConnectionStatusPendingApproval,
		ecx.ConnectionStatusPendingBGPPeering,
		ecx.ConnectionStatusPendingProviderVlan,
	}
	return createConnectionStatusWaitConfiguration(fetchFunc, id, delay, timeout, target, pending)
}

func createConnectionStatusDeleteWaitConfiguration(fetchFunc getL2Connection, id string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	pending := []string{
		ecx.ConnectionStatusDeprovisioning,
	}
	target := []string{
		ecx.ConnectionStatusPendingDelete,
		ecx.ConnectionStatusDeprovisioned,
		ecx.ConnectionStatusDeleted,
	}
	return createConnectionStatusWaitConfiguration(fetchFunc, id, delay, timeout, target, pending)
}

func createConnectionStatusWaitConfiguration(fetchFunc getL2Connection, id string, delay time.Duration, timeout time.Duration, target []string, pending []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: delay,
		Refresh: func() (interface{}, string, error) {
			resp, err := fetchFunc(id)
			if err != nil {
				return nil, "", err
			}
			return resp, ecx.StringValue(resp.Status), nil
		},
	}
}
