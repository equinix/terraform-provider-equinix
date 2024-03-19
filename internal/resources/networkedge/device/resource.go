package device

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_schema "github.com/equinix/terraform-provider-equinix/internal/schema"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"

	"github.com/equinix/ne-go"
	"github.com/equinix/rest-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var neDeviceSchemaNames = map[string]string{
	"UUID":                  "uuid",
	"Name":                  "name",
	"TypeCode":              "type_code",
	"Status":                "status",
	"MetroCode":             "metro_code",
	"IBX":                   "ibx",
	"Region":                "region",
	"Throughput":            "throughput",
	"ThroughputUnit":        "throughput_unit",
	"HostName":              "hostname",
	"PackageCode":           "package_code",
	"Version":               "version",
	"IsBYOL":                "byol",
	"LicenseToken":          "license_token",
	"LicenseFile":           "license_file",
	"LicenseFileID":         "license_file_id",
	"CloudInitFileID":       "cloud_init_file_id",
	"LicenseStatus":         "license_status",
	"ACLTemplateUUID":       "acl_template_id",
	"MgmtAclTemplateUuid":   "mgmt_acl_template_uuid",
	"SSHIPAddress":          "ssh_ip_address",
	"SSHIPFqdn":             "ssh_ip_fqdn",
	"AccountNumber":         "account_number",
	"Notifications":         "notifications",
	"PurchaseOrderNumber":   "purchase_order_number",
	"RedundancyType":        "redundancy_type",
	"RedundantUUID":         "redundant_id",
	"ProjectID":             "project_id",
	"TermLength":            "term_length",
	"AdditionalBandwidth":   "additional_bandwidth",
	"OrderReference":        "order_reference",
	"InterfaceCount":        "interface_count",
	"CoreCount":             "core_count",
	"IsSelfManaged":         "self_managed",
	"WanInterfaceId":        "wan_interface_id",
	"Interfaces":            "interface",
	"VendorConfiguration":   "vendor_configuration",
	"UserPublicKey":         "ssh_key",
	"ASN":                   "asn",
	"ZoneCode":              "zone_code",
	"Secondary":             "secondary_device",
	"ClusterDetails":        "cluster_details",
	"ValidStatusList":       "valid_status_list",
	"Connectivity":          "connectivity",
	"DiverseFromDeviceUUID": "diverse_device_id",
	"DiverseFromDeviceName": "diverse_device_name",
}

var neDeviceDescriptions = map[string]string{
	"UUID":                  "Device unique identifier",
	"Name":                  "Device name",
	"TypeCode":              "Device type code",
	"Status":                "Device provisioning status",
	"MetroCode":             "Device location metro code",
	"IBX":                   "Device location Equinix Business Exchange name",
	"Region":                "Device location region",
	"Throughput":            "Device license throughput",
	"ThroughputUnit":        "Device license throughput unit (Mbps or Gbps)",
	"HostName":              "Device hostname prefix",
	"PackageCode":           "Device software package code",
	"Version":               "Device software software version",
	"IsBYOL":                "Boolean value that determines device licensing mode: bring your own license or subscription (default)",
	"LicenseToken":          "License Token applicable for some device types in BYOL licensing mode",
	"LicenseFile":           "Path to the license file that will be uploaded and applied on a device, applicable for some device types in BYOL licensing mode",
	"LicenseFileID":         "Unique identifier of applied license file",
	"CloudInitFileID":       "Unique identifier of applied cloud init file",
	"LicenseStatus":         "Device license registration status",
	"ACLTemplateUUID":       "Unique identifier of applied ACL template",
	"MgmtAclTemplateUuid":   "Unique identifier of applied MGMT ACL template",
	"SSHIPAddress":          "IP address of SSH enabled interface on the device",
	"SSHIPFqdn":             "FQDN of SSH enabled interface on the device",
	"AccountNumber":         "Device billing account number",
	"Notifications":         "List of email addresses that will receive device status notifications",
	"PurchaseOrderNumber":   "Purchase order number associated with a device order",
	"RedundancyType":        "Device redundancy type applicable for HA devices, either primary or secondary",
	"RedundantUUID":         "Unique identifier for a redundant device, applicable for HA device",
	"TermLength":            "Device term length",
	"AdditionalBandwidth":   "Additional Internet bandwidth, in Mbps, that will be allocated to the device",
	"OrderReference":        "Name/number used to identify device order on the invoice",
	"InterfaceCount":        "Number of network interfaces on a device. If not specified, default number for a given device type will be used",
	"CoreCount":             "Number of CPU cores used by device",
	"IsSelfManaged":         "Boolean value that determines device management mode: self-managed or subscription (default)",
	"WanInterfaceId":        "device interface id picked for WAN",
	"Interfaces":            "List of device interfaces",
	"VendorConfiguration":   "Map of vendor specific configuration parameters for a device (controller1, activationKey, managementType, siteId, systemIpAddress)",
	"UserPublicKey":         "Definition of SSH key that will be provisioned on a device",
	"ASN":                   "Autonomous system number",
	"ZoneCode":              "Device location zone code",
	"Secondary":             "Definition of secondary device applicable for HA setup",
	"ClusterDetails":        "An object that has the cluster details",
	"ValidStatusList":       "Comma Separated List of states to be considered valid when searching by name",
	"Connectivity":          "Parameter to identify internet access for device. Supported Values: INTERNET-ACCESS(default) or PRIVATE or INTERNET-ACCESS-WITH-PRVT-MGMT",
	"ProjectID":             "The unique identifier of Project Resource to which device is scoped to",
	"DiverseFromDeviceUUID": "Unique ID of an existing device",
	"DiverseFromDeviceName": "Diverse Device Name of an existing device",
}

var neDeviceInterfaceSchemaNames = map[string]string{
	"ID":                "id",
	"Name":              "name",
	"Status":            "status",
	"OperationalStatus": "operational_status",
	"MACAddress":        "mac_address",
	"IPAddress":         "ip_address",
	"AssignedType":      "assigned_type",
	"Type":              "type",
}

var neDeviceInterfaceDescriptions = map[string]string{
	"ID":                "Interface identifier",
	"Name":              "Interface name",
	"Status":            "Interface status (AVAILABLE, RESERVED, ASSIGNED)",
	"OperationalStatus": "Interface operational status (up or down)",
	"MACAddress":        "Interface MAC addres",
	"IPAddress":         "interface IP address",
	"AssignedType":      "Interface management type (Equinix Managed or empty)",
	"Type":              "Interface type",
}

var neDeviceUserKeySchemaNames = map[string]string{
	"Username": "username",
	"KeyName":  "key_name",
}

var neDeviceUserKeyDescriptions = map[string]string{
	"Username": "Username associated with given key",
	"KeyName":  "Reference by name to previously provisioned public SSH key",
}

var neDeviceClusterSchemaNames = map[string]string{
	"ClusterId":   "cluster_id",
	"ClusterName": "cluster_name",
	"NumOfNodes":  "num_of_nodes",
	"Node0":       "node0",
	"Node1":       "node1",
}

var neDeviceClusterDescriptions = map[string]string{
	"ClusterId":   "The id of the cluster",
	"ClusterName": "The name of the cluster device",
	"NumOfNodes":  "The number of nodes in the cluster",
	"Node0":       "An object that has node0 details",
	"Node1":       "An object that has node1 details",
}

var neDeviceClusterNodeSchemaNames = map[string]string{
	"VendorConfiguration": "vendor_configuration",
	"LicenseFileId":       "license_file_id",
	"LicenseToken":        "license_token",
	"UUID":                "uuid",
	"Name":                "name",
}

var neDeviceClusterNodeDescriptions = map[string]string{
	"VendorConfiguration": "An object that has fields relevant to the vendor of the cluster device",
	"LicenseFileId":       "License file id. This is necessary for Fortinet and Juniper clusters",
	"LicenseToken":        "License token. This is necessary for Palo Alto clusters",
	"UUID":                "The unique id of the node",
	"Name":                "The name of the node",
}

var neDeviceVendorConfigSchemaNames = map[string]string{
	"Hostname":       "hostname",
	"AdminPassword":  "admin_password",
	"Controller1":    "controller1",
	"ActivationKey":  "activation_key",
	"ControllerFqdn": "controller_fqdn",
	"RootPassword":   "root_password",
}

var neDeviceVendorConfigDescriptions = map[string]string{
	"Hostname":       "Hostname. This is necessary for Palo Alto, Juniper, and Fortinet clusters",
	"AdminPassword":  "The administrative password of the device. You can use it to log in to the console. This field is not available for all device types",
	"Controller1":    "System IP Address. Mandatory for the Fortinet SDWAN cluster device",
	"ActivationKey":  "Activation key. This is required for Velocloud clusters",
	"ControllerFqdn": "Controller fqdn. This is required for Velocloud clusters",
	"RootPassword":   "The CLI password of the device. This field is relevant only for the Velocloud SDWAN cluster",
}

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkDeviceCreate,
		ReadContext:   resourceNetworkDeviceRead,
		UpdateContext: resourceNetworkDeviceUpdate,
		DeleteContext: resourceNetworkDeviceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: createNetworkDeviceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Description: "Resource allows creation and management of Equinix Network Edge virtual devices",
	}
}

func createNetworkDeviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["UUID"],
		},
		neDeviceSchemaNames["Name"]: {
			Type:     schema.TypeString,
			Required: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if old == new+"-Node0" {
					return true
				}
				return false
			},
			ValidateFunc: validation.StringLenBetween(3, 50),
			Description:  neDeviceDescriptions["Name"],
		},
		neDeviceSchemaNames["TypeCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["TypeCode"],
		},
		neDeviceSchemaNames["Status"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["Status"],
		},
		neDeviceSchemaNames["LicenseStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["LicenseStatus"],
		},
		neDeviceSchemaNames["MetroCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: equinix_validation.StringIsMetroCode,
			Description:  neDeviceDescriptions["MetroCode"],
		},
		neDeviceSchemaNames["IBX"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["IBX"],
		},
		neDeviceSchemaNames["Region"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["Region"],
		},
		neDeviceSchemaNames["Throughput"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  neDeviceDescriptions["Throughput"],
		},
		neDeviceSchemaNames["ThroughputUnit"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"Mbps", "Gbps"}, false),
			RequiredWith: []string{neDeviceSchemaNames["Throughput"]},
			Description:  neDeviceDescriptions["ThroughputUnit"],
		},
		neDeviceSchemaNames["HostName"]: {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Description: neDeviceDescriptions["HostName"],
		},
		neDeviceSchemaNames["PackageCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["PackageCode"],
		},
		neDeviceSchemaNames["Version"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["Version"],
		},
		neDeviceSchemaNames["IsBYOL"]: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: neDeviceDescriptions["IsBYOL"],
		},
		neDeviceSchemaNames["LicenseToken"]: {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{neDeviceSchemaNames["LicenseFile"]},
			Description:   neDeviceDescriptions["LicenseToken"],
		},
		neDeviceSchemaNames["LicenseFile"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["LicenseFile"],
		},
		neDeviceSchemaNames["LicenseFileID"]: {
			Type:          schema.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{neDeviceSchemaNames["LicenseFile"]},
			Description:   neDeviceDescriptions["LicenseFileID"],
		},
		neDeviceSchemaNames["CloudInitFileID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["CloudInitFileID"],
		},
		neDeviceSchemaNames["ACLTemplateUUID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["ACLTemplateUUID"],
		},
		neDeviceSchemaNames["MgmtAclTemplateUuid"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["MgmtAclTemplateUuid"],
		},
		neDeviceSchemaNames["SSHIPAddress"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["SSHIPAddress"],
		},
		neDeviceSchemaNames["SSHIPFqdn"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["SSHIPFqdn"],
		},
		neDeviceSchemaNames["AccountNumber"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["AccountNumber"],
		},
		neDeviceSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: equinix_validation.StringIsEmailAddress,
			},
			Description: neDeviceDescriptions["Notifications"],
		},
		neDeviceSchemaNames["PurchaseOrderNumber"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 30),
			Description:  neDeviceDescriptions["PurchaseOrderNumber"],
		},
		neDeviceSchemaNames["RedundancyType"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["RedundancyType"],
		},
		neDeviceSchemaNames["RedundantUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["RedundantUUID"],
		},
		neDeviceSchemaNames["TermLength"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntInSlice([]int{1, 12, 24, 36}),
			Description:  neDeviceDescriptions["TermLength"],
		},
		neDeviceSchemaNames["ProjectID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.IsUUID,
			Description:  neDeviceDescriptions["ProjectID"],
		},
		neDeviceSchemaNames["DiverseFromDeviceUUID"]: {
			Type:          schema.TypeString,
			Computed:      true,
			Optional:      true,
			ValidateFunc:  validation.IsUUID,
			ConflictsWith: []string{neDeviceSchemaNames["Secondary"]},
			Description:   neDeviceDescriptions["DiverseFromDeviceUUID"],
		},
		neDeviceSchemaNames["DiverseFromDeviceName"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["DiverseFromDeviceName"],
		},
		neDeviceSchemaNames["AdditionalBandwidth"]: {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: neDeviceDescriptions["AdditionalBandwidth"],
		},
		neDeviceSchemaNames["OrderReference"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 100),
			Description:  neDeviceDescriptions["OrderReference"],
		},
		neDeviceSchemaNames["InterfaceCount"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  neDeviceDescriptions["InterfaceCount"],
		},
		neDeviceSchemaNames["CoreCount"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  neDeviceDescriptions["CoreCount"],
		},
		neDeviceSchemaNames["IsSelfManaged"]: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: neDeviceDescriptions["IsSelfManaged"],
		},
		neDeviceSchemaNames["WanInterfaceId"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceDescriptions["WanInterfaceId"],
		},
		neDeviceSchemaNames["Interfaces"]: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createNetworkDeviceInterfaceSchema(),
			},
			Description: neDeviceDescriptions["Interfaces"],
		},
		neDeviceSchemaNames["VendorConfiguration"]: {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			Description: neDeviceDescriptions["VendorConfiguration"],
		},
		neDeviceSchemaNames["UserPublicKey"]: {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: createNetworkDeviceUserKeySchema(),
			},
			Description: neDeviceDescriptions["UserPublicKey"],
		},
		neDeviceSchemaNames["ASN"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: neDeviceDescriptions["ASN"],
		},
		neDeviceSchemaNames["ZoneCode"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["ZoneCode"],
		},
		neDeviceSchemaNames["Connectivity"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "INTERNET-ACCESS",
			ValidateFunc: validation.StringInSlice([]string{"INTERNET-ACCESS", "PRIVATE", "INTERNET-ACCESS-WITH-PRVT-MGMT"}, false),
			Description:  neDeviceDescriptions["Connectivity"],
		},
		neDeviceSchemaNames["Secondary"]: {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: neDeviceDescriptions["Secondary"],
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					neDeviceSchemaNames["UUID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["UUID"],
					},
					neDeviceSchemaNames["ProjectID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["ProjectID"],
					},
					neDeviceSchemaNames["Name"]: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(3, 50),
						Description:  neDeviceDescriptions["Name"],
					},
					neDeviceSchemaNames["Status"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["Status"],
					},
					neDeviceSchemaNames["LicenseStatus"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["LicenseStatus"],
					},
					neDeviceSchemaNames["MetroCode"]: {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: equinix_validation.StringIsMetroCode,
						Description:  neDeviceDescriptions["MetroCode"],
					},
					neDeviceSchemaNames["IBX"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["IBX"],
					},
					neDeviceSchemaNames["Region"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["Region"],
					},
					neDeviceSchemaNames["HostName"]: {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: neDeviceDescriptions["HostName"],
					},
					neDeviceSchemaNames["LicenseToken"]: {
						Type:          schema.TypeString,
						Optional:      true,
						ForceNew:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{neDeviceSchemaNames["Secondary"] + ".0." + neDeviceSchemaNames["LicenseFile"]},
						Description:   neDeviceDescriptions["LicenseToken"],
					},
					neDeviceSchemaNames["LicenseFile"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  neDeviceDescriptions["LicenseFile"],
					},
					neDeviceSchemaNames["LicenseFileID"]: {
						Type:          schema.TypeString,
						Optional:      true,
						Computed:      true,
						ForceNew:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{neDeviceSchemaNames["Secondary"] + ".0." + neDeviceSchemaNames["LicenseFile"]},
						Description:   neDeviceDescriptions["LicenseFileID"],
					},
					neDeviceSchemaNames["CloudInitFileID"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  neDeviceDescriptions["CloudInitFileID"],
					},
					neDeviceSchemaNames["ACLTemplateUUID"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  neDeviceDescriptions["ACLTemplateUUID"],
					},
					neDeviceSchemaNames["MgmtAclTemplateUuid"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  neDeviceDescriptions["MgmtAclTemplateUuid"],
					},
					neDeviceSchemaNames["SSHIPAddress"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["SSHIPAddress"],
					},
					neDeviceSchemaNames["SSHIPFqdn"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["SSHIPFqdn"],
					},
					neDeviceSchemaNames["AccountNumber"]: {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  neDeviceDescriptions["AccountNumber"],
					},
					neDeviceSchemaNames["Notifications"]: {
						Type:     schema.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: equinix_validation.StringIsEmailAddress,
						},
						Description: neDeviceDescriptions["Notifications"],
					},
					neDeviceSchemaNames["RedundancyType"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["RedundancyType"],
					},
					neDeviceSchemaNames["RedundantUUID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["RedundantUUID"],
					},
					neDeviceSchemaNames["AdditionalBandwidth"]: {
						Type:        schema.TypeInt,
						Optional:    true,
						Computed:    true,
						Description: neDeviceDescriptions["AdditionalBandwidth"],
					},
					neDeviceSchemaNames["WanInterfaceId"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  neDeviceDescriptions["WanInterfaceId"],
					},
					neDeviceSchemaNames["Interfaces"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createNetworkDeviceInterfaceSchema(),
						},
						Description: neDeviceDescriptions["Interfaces"],
					},
					neDeviceSchemaNames["VendorConfiguration"]: {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						Description: neDeviceDescriptions["VendorConfiguration"],
					},
					neDeviceSchemaNames["UserPublicKey"]: {
						Type:     schema.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: createNetworkDeviceUserKeySchema(),
						},
						Description: neDeviceDescriptions["UserPublicKey"],
					},
					neDeviceSchemaNames["ASN"]: {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: neDeviceDescriptions["ASN"],
					},
					neDeviceSchemaNames["ZoneCode"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["ZoneCode"],
					},
				},
			},
		},
		neDeviceSchemaNames["ClusterDetails"]: {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: neDeviceDescriptions["ClusterDetails"],
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					neDeviceClusterSchemaNames["ClusterId"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceClusterDescriptions["ClusterId"],
					},
					neDeviceClusterSchemaNames["ClusterName"]: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(3, 50),
						Description:  neDeviceClusterDescriptions["ClusterName"],
					},
					neDeviceClusterSchemaNames["NumOfNodes"]: {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: neDeviceClusterDescriptions["NumOfNodes"],
					},
					neDeviceClusterSchemaNames["Node0"]: {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: createClusterNodeDetailSchema(),
						},
						Description: neDeviceClusterDescriptions["Node0"],
					},
					neDeviceClusterSchemaNames["Node1"]: {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: createClusterNodeDetailSchema(),
						},
						Description: neDeviceClusterDescriptions["Node1"],
					},
				},
			},
		},
	}
}

func createNetworkDeviceInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceInterfaceSchemaNames["ID"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: neDeviceInterfaceDescriptions["ID"],
		},
		neDeviceInterfaceSchemaNames["Name"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceInterfaceDescriptions["Name"],
		},
		neDeviceInterfaceSchemaNames["Status"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceInterfaceDescriptions["Status"],
		},
		neDeviceInterfaceSchemaNames["OperationalStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceInterfaceDescriptions["OperationalStatus"],
		},
		neDeviceInterfaceSchemaNames["MACAddress"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceInterfaceDescriptions["MACAddress"],
		},
		neDeviceInterfaceSchemaNames["IPAddress"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceInterfaceDescriptions["IPAddress"],
		},
		neDeviceInterfaceSchemaNames["AssignedType"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceInterfaceDescriptions["AssignedType"],
		},
		neDeviceInterfaceSchemaNames["Type"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceInterfaceDescriptions["Type"],
		},
	}
}

func createNetworkDeviceUserKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceUserKeySchemaNames["Username"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceUserKeyDescriptions["Username"],
		},
		neDeviceUserKeySchemaNames["KeyName"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  neDeviceUserKeyDescriptions["KeyName"],
		},
	}
}

func createClusterNodeDetailSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceClusterNodeSchemaNames["LicenseFileId"]: {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{neDeviceSchemaNames["LicenseFileID"]},
			Description:   neDeviceClusterNodeDescriptions["LicenseFileId"],
		},
		neDeviceClusterNodeSchemaNames["LicenseToken"]: {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			Sensitive:     true,
			ConflictsWith: []string{neDeviceSchemaNames["LicenseToken"]},
			Description:   neDeviceClusterNodeDescriptions["LicenseToken"],
		},
		neDeviceClusterNodeSchemaNames["VendorConfiguration"]: {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: createVendorConfigurationSchema(),
			},
			ConflictsWith: []string{neDeviceSchemaNames["VendorConfiguration"]},
			Description:   neDeviceClusterNodeDescriptions["VendorConfiguration"],
		},
		neDeviceClusterNodeSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceClusterNodeDescriptions["UUID"],
		},
		neDeviceClusterNodeSchemaNames["Name"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceClusterNodeDescriptions["Name"],
		},
	}
}

func createVendorConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceVendorConfigSchemaNames["Hostname"]: {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: neDeviceVendorConfigDescriptions["Hostname"],
		},
		neDeviceVendorConfigSchemaNames["AdminPassword"]: {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: neDeviceVendorConfigDescriptions["AdminPassword"],
		},
		neDeviceVendorConfigSchemaNames["Controller1"]: {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: neDeviceVendorConfigDescriptions["Controller1"],
		},
		neDeviceVendorConfigSchemaNames["ActivationKey"]: {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: neDeviceVendorConfigDescriptions["ActivationKey"],
		},
		neDeviceVendorConfigSchemaNames["ControllerFqdn"]: {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: neDeviceVendorConfigDescriptions["ControllerFqdn"],
		},
		neDeviceVendorConfigSchemaNames["RootPassword"]: {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: neDeviceVendorConfigDescriptions["RootPassword"],
		},
	}
}

func resourceNetworkDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	primary, secondary := createNetworkDevices(d)
	var err error
	if err := uploadDeviceLicenseFile(os.Open, client.UploadLicenseFile, ne.StringValue(primary.TypeCode), primary); err != nil {
		return diag.Errorf("could not upload primary device license file due to %s", err)
	}
	if err := uploadDeviceLicenseFile(os.Open, client.UploadLicenseFile, ne.StringValue(primary.TypeCode), secondary); err != nil {
		return diag.Errorf("could not upload secondary device license file due to %s", err)
	}
	if secondary != nil {
		primary.UUID, secondary.UUID, err = client.CreateRedundantDevice(*primary, *secondary)
	} else {
		primary.UUID, err = client.CreateDevice(*primary)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ne.StringValue(primary.UUID))
	waitConfigs := []*retry.StateChangeConf{
		createNetworkDeviceStatusProvisioningWaitConfiguration(client.GetDevice, ne.StringValue(primary.UUID), 5*time.Second, d.Timeout(schema.TimeoutCreate)),
		createNetworkDeviceLicenseStatusWaitConfiguration(client.GetDevice, ne.StringValue(primary.UUID), 5*time.Second, d.Timeout(schema.TimeoutCreate)),
	}
	if ne.StringValue(primary.ACLTemplateUUID) != "" || ne.StringValue(primary.MgmtAclTemplateUuid) != "" {
		waitConfigs = append(waitConfigs,
			createNetworkDeviceACLStatusWaitConfiguration(client.GetDeviceACLDetails, ne.StringValue(primary.UUID), 1*time.Second, d.Timeout(schema.TimeoutUpdate)),
		)
	}
	if secondary != nil {
		waitConfigs = append(waitConfigs,
			createNetworkDeviceStatusProvisioningWaitConfiguration(client.GetDevice, ne.StringValue(secondary.UUID), 5*time.Second, d.Timeout(schema.TimeoutCreate)),
			createNetworkDeviceLicenseStatusWaitConfiguration(client.GetDevice, ne.StringValue(secondary.UUID), 5*time.Second, d.Timeout(schema.TimeoutCreate)),
		)
		if ne.StringValue(secondary.ACLTemplateUUID) != "" || ne.StringValue(secondary.MgmtAclTemplateUuid) != "" {
			waitConfigs = append(waitConfigs,
				createNetworkDeviceACLStatusWaitConfiguration(client.GetDeviceACLDetails, ne.StringValue(secondary.UUID), 1*time.Second, d.Timeout(schema.TimeoutUpdate)),
			)
		}
	}
	for _, config := range waitConfigs {
		if config == nil {
			continue
		}
		if _, err := config.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error waiting for network device (%s) to be created: %s", ne.StringValue(primary.UUID), err)
		}
	}
	diags = append(diags, resourceNetworkDeviceRead(ctx, d, m)...)
	return diags
}

func resourceNetworkDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	var err error
	var primary, secondary *ne.Device
	primary, err = client.GetDevice(d.Id())
	if err != nil {
		return diag.Errorf("cannot fetch primary network device due to %v", err)
	}
	if slices.Contains([]string{ne.DeviceStateDeprovisioning, ne.DeviceStateDeprovisioned}, ne.StringValue(primary.Status)) {
		d.SetId("")
		return diags
	}
	if ne.StringValue(primary.RedundantUUID) != "" {
		secondary, err = client.GetDevice(ne.StringValue(primary.RedundantUUID))
		if err != nil {
			return diag.Errorf("cannot fetch secondary network device due to %v", err)
		}
	}
	if err = updateNetworkDeviceResource(primary, secondary, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceNetworkDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	supportedChanges := []string{
		neDeviceSchemaNames["Name"], neDeviceSchemaNames["TermLength"], neDeviceSchemaNames["CoreCount"],
		neDeviceSchemaNames["Notifications"], neDeviceSchemaNames["AdditionalBandwidth"],
		neDeviceSchemaNames["ACLTemplateUUID"], neDeviceSchemaNames["MgmtAclTemplateUuid"],
	}
	updateReq := client.NewDeviceUpdateRequest(d.Id())
	primaryChanges := equinix_schema.GetResourceDataChangedKeys(supportedChanges, d)
	var clusterChanges map[string]interface{}
	clusterSupportedChanges := []string{neDeviceClusterSchemaNames["ClusterName"]}
	if _, ok := d.GetOk(neDeviceSchemaNames["ClusterDetails"]); ok {
		clusterChanges = equinix_schema.GetResourceDataListElementChanges(clusterSupportedChanges, neDeviceSchemaNames["ClusterDetails"], 0, d)
		for key, value := range clusterChanges {
			primaryChanges[key] = value
		}
	}
	if err := fillNetworkDeviceUpdateRequest(updateReq, primaryChanges).Execute(); err != nil {
		return diag.FromErr(err)
	}
	var secondaryChanges map[string]interface{}
	if v, ok := d.GetOk(neDeviceSchemaNames["RedundantUUID"]); ok {
		secondaryChanges = equinix_schema.GetResourceDataListElementChanges(supportedChanges, neDeviceSchemaNames["Secondary"], 0, d)
		secondaryUpdateReq := client.NewDeviceUpdateRequest(v.(string))
		if err := fillNetworkDeviceUpdateRequest(secondaryUpdateReq, secondaryChanges).Execute(); err != nil {
			return diag.FromErr(err)
		}
	}
	for _, stateChangeConf := range getNetworkDeviceStateChangeConfigs(client, d.Id(), d.Timeout(schema.TimeoutUpdate), primaryChanges) {
		if _, err := stateChangeConf.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error waiting for network device %q to be updated: %s", d.Id(), err)
		}
	}
	for _, stateChangeConf := range getNetworkDeviceStateChangeConfigs(client, d.Get(neDeviceSchemaNames["RedundantUUID"]).(string), d.Timeout(schema.TimeoutUpdate), secondaryChanges) {
		if _, err := stateChangeConf.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error waiting for network device %q to be updated: %s", d.Get(neDeviceSchemaNames["RedundantUUID"]), err)
		}
	}
	diags = append(diags, resourceNetworkDeviceRead(ctx, d, m)...)
	return diags
}

func resourceNetworkDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	waitConfigs := []*retry.StateChangeConf{
		createNetworkDeviceStatusDeleteWaitConfiguration(client.GetDevice, d.Id(), 5*time.Second, d.Timeout(schema.TimeoutDelete)),
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Secondary"]); ok {
		if secondary := expandNetworkDeviceSecondary(v.([]interface{})); secondary != nil {
			waitConfigs = append(waitConfigs,
				createNetworkDeviceStatusDeleteWaitConfiguration(client.GetDevice, ne.StringValue(secondary.UUID), 5*time.Second, d.Timeout(schema.TimeoutDelete)),
			)
		}
	}
	if err := client.DeleteDevice(d.Id()); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			for _, detailedErr := range restErr.ApplicationErrors {
				if detailedErr.Code == ne.ErrorCodeDeviceRemoved {
					return diags
				}
			}
		}
		return diag.FromErr(err)
	}
	for _, config := range waitConfigs {
		if _, err := config.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error waiting for network device (%s) to be removed: %s", d.Id(), err)
		}
	}
	return diags
}

func createNetworkDevices(d *schema.ResourceData) (*ne.Device, *ne.Device) {
	var primary, secondary *ne.Device
	primary = &ne.Device{}
	if v, ok := d.GetOk(neDeviceSchemaNames["Name"]); ok {
		primary.Name = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["TypeCode"]); ok {
		primary.TypeCode = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["MetroCode"]); ok {
		primary.MetroCode = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ProjectID"]); ok {
		primary.ProjectID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["DiverseFromDeviceUUID"]); ok {
		primary.DiverseFromDeviceUUID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Throughput"]); ok {
		primary.Throughput = ne.Int(v.(int))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ThroughputUnit"]); ok {
		primary.ThroughputUnit = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["HostName"]); ok {
		primary.HostName = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["PackageCode"]); ok {
		primary.PackageCode = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Version"]); ok {
		primary.Version = ne.String(v.(string))
	}
	primary.IsBYOL = ne.Bool(d.Get(neDeviceSchemaNames["IsBYOL"]).(bool))
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseToken"]); ok {
		primary.LicenseToken = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseFile"]); ok {
		primary.LicenseFile = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseFileID"]); ok {
		primary.LicenseFileID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["CloudInitFileID"]); ok {
		primary.CloudInitFileID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ACLTemplateUUID"]); ok {
		primary.ACLTemplateUUID = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["MgmtAclTemplateUuid"]); ok {
		primary.MgmtAclTemplateUuid = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["AccountNumber"]); ok {
		primary.AccountNumber = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Notifications"]); ok {
		primary.Notifications = converters.SetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["PurchaseOrderNumber"]); ok {
		primary.PurchaseOrderNumber = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["TermLength"]); ok {
		primary.TermLength = ne.Int(v.(int))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["AdditionalBandwidth"]); ok {
		primary.AdditionalBandwidth = ne.Int(v.(int))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["OrderReference"]); ok {
		primary.OrderReference = ne.String(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["InterfaceCount"]); ok {
		primary.InterfaceCount = ne.Int(v.(int))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["CoreCount"]); ok {
		primary.CoreCount = ne.Int(v.(int))
	}
	primary.IsSelfManaged = ne.Bool(d.Get(neDeviceSchemaNames["IsSelfManaged"]).(bool))
	if v, ok := d.GetOk(neDeviceSchemaNames["VendorConfiguration"]); ok {
		primary.VendorConfiguration = converters.InterfaceMapToStringMap(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["WanInterfaceId"]); ok {
		primary.WanInterfaceId = ne.String(v.(string))
	}

	if v, ok := d.GetOk(neDeviceSchemaNames["UserPublicKey"]); ok {
		userKeys := expandNetworkDeviceUserKeys(v.(*schema.Set))
		if len(userKeys) > 0 {
			primary.UserPublicKey = userKeys[0]
		}
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Secondary"]); ok {
		secondary = expandNetworkDeviceSecondary(v.([]interface{}))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ClusterDetails"]); ok {
		primary.ClusterDetails = expandNetworkDeviceClusterDetails(v.([]interface{}))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Connectivity"]); ok {
		primary.Connectivity = ne.String(v.(string))
	}
	return primary, secondary
}

func updateNetworkDeviceResource(primary *ne.Device, secondary *ne.Device, d *schema.ResourceData) error {
	if err := d.Set(neDeviceSchemaNames["UUID"], primary.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Name"], primary.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["ProjectID"], primary.ProjectID); err != nil {
		return fmt.Errorf("error reading ProjectID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["DiverseFromDeviceUUID"], primary.DiverseFromDeviceUUID); err != nil {
		return fmt.Errorf("error reading DiverseFromDeviceUUID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["DiverseFromDeviceName"], primary.DiverseFromDeviceName); err != nil {
		return fmt.Errorf("error reading DiverseFromDeviceName: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["TypeCode"], primary.TypeCode); err != nil {
		return fmt.Errorf("error reading TypeCode: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Status"], primary.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["LicenseStatus"], primary.LicenseStatus); err != nil {
		return fmt.Errorf("error reading LicenseStatus: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["MetroCode"], primary.MetroCode); err != nil {
		return fmt.Errorf("error reading MetroCode: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["IBX"], primary.IBX); err != nil {
		return fmt.Errorf("error reading IBX: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Region"], primary.Region); err != nil {
		return fmt.Errorf("error reading Region: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Throughput"], primary.Throughput); err != nil {
		return fmt.Errorf("error reading Throughput: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["ThroughputUnit"], primary.ThroughputUnit); err != nil {
		return fmt.Errorf("error reading ThroughputUnit: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["HostName"], primary.HostName); err != nil {
		return fmt.Errorf("error reading HostName: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["PackageCode"], primary.PackageCode); err != nil {
		return fmt.Errorf("error reading PackageCode: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Version"], primary.Version); err != nil {
		return fmt.Errorf("error reading Version: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["IsBYOL"], primary.IsBYOL); err != nil {
		return fmt.Errorf("error reading IsBYOL: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["LicenseFileID"], primary.LicenseFileID); err != nil {
		return fmt.Errorf("error reading LicenseFileID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["ACLTemplateUUID"], primary.ACLTemplateUUID); err != nil {
		return fmt.Errorf("error reading ACLTemplateUUID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["SSHIPAddress"], primary.SSHIPAddress); err != nil {
		return fmt.Errorf("error reading SSHIPAddress: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["SSHIPFqdn"], primary.SSHIPFqdn); err != nil {
		return fmt.Errorf("error reading SSHIPFqdn: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["AccountNumber"], primary.AccountNumber); err != nil {
		return fmt.Errorf("error reading AccountNumber: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Notifications"], primary.Notifications); err != nil {
		return fmt.Errorf("error reading Notifications: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["PurchaseOrderNumber"], primary.PurchaseOrderNumber); err != nil {
		return fmt.Errorf("error reading PurchaseOrderNumber: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["RedundancyType"], primary.RedundancyType); err != nil {
		return fmt.Errorf("error reading RedundancyType: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["RedundantUUID"], primary.RedundantUUID); err != nil {
		return fmt.Errorf("error reading RedundantUUID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["TermLength"], primary.TermLength); err != nil {
		return fmt.Errorf("error reading TermLength: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["AdditionalBandwidth"], primary.AdditionalBandwidth); err != nil {
		return fmt.Errorf("error reading AdditionalBandwidth: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["OrderReference"], primary.OrderReference); err != nil {
		return fmt.Errorf("error reading OrderReference: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["InterfaceCount"], primary.InterfaceCount); err != nil {
		return fmt.Errorf("error reading InterfaceCount: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["CoreCount"], primary.CoreCount); err != nil {
		return fmt.Errorf("error reading CoreCount: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["IsSelfManaged"], primary.IsSelfManaged); err != nil {
		return fmt.Errorf("error reading IsSelfManaged: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Interfaces"], flattenNetworkDeviceInterfaces(primary.Interfaces)); err != nil {
		return fmt.Errorf("error reading Interfaces: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["VendorConfiguration"], primary.VendorConfiguration); err != nil {
		return fmt.Errorf("error reading VendorConfiguration: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["UserPublicKey"], flattenNetworkDeviceUserKeys([]*ne.DeviceUserPublicKey{primary.UserPublicKey})); err != nil {
		return fmt.Errorf("error reading VendorConfiguration: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["ASN"], primary.ASN); err != nil {
		return fmt.Errorf("error reading ASN: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["ZoneCode"], primary.ZoneCode); err != nil {
		return fmt.Errorf("error reading ZoneCode: %s", err)
	}
	if secondary != nil {
		if v, ok := d.GetOk(neDeviceSchemaNames["Secondary"]); ok {
			secondaryFromSchema := expandNetworkDeviceSecondary(v.([]interface{}))
			secondary.LicenseFile = secondaryFromSchema.LicenseFile
			secondary.LicenseToken = secondaryFromSchema.LicenseToken
			secondary.CloudInitFileID = secondaryFromSchema.CloudInitFileID
		}
		if err := d.Set(neDeviceSchemaNames["Secondary"], flattenNetworkDeviceSecondary(secondary)); err != nil {
			return fmt.Errorf("error reading Secondary: %s", err)
		}
	}
	if primary.ClusterDetails != nil {
		if v, ok := d.GetOk(neDeviceSchemaNames["ClusterDetails"]); ok {
			clusterDetailsFromSchema := expandNetworkDeviceClusterDetails(v.([]interface{}))
			primary.ClusterDetails.Node0.LicenseFileId = clusterDetailsFromSchema.Node0.LicenseFileId
			primary.ClusterDetails.Node0.LicenseToken = clusterDetailsFromSchema.Node0.LicenseToken
			primary.ClusterDetails.Node1.LicenseFileId = clusterDetailsFromSchema.Node1.LicenseFileId
			primary.ClusterDetails.Node1.LicenseToken = clusterDetailsFromSchema.Node1.LicenseToken
		}
		if err := d.Set(neDeviceSchemaNames["ClusterDetails"], flattenNetworkDeviceClusterDetails(primary.ClusterDetails)); err != nil {
			return fmt.Errorf("error reading ClusterDetails: %s", err)
		}
	}
	return nil
}

func flattenNetworkDeviceSecondary(device *ne.Device) interface{} {
	transformed := make(map[string]interface{})
	transformed[neDeviceSchemaNames["UUID"]] = device.UUID
	transformed[neDeviceSchemaNames["Name"]] = device.Name
	transformed[neDeviceSchemaNames["Status"]] = device.Status
	transformed[neDeviceSchemaNames["LicenseStatus"]] = device.LicenseStatus
	transformed[neDeviceSchemaNames["MetroCode"]] = device.MetroCode
	transformed[neDeviceSchemaNames["IBX"]] = device.IBX
	transformed[neDeviceSchemaNames["Region"]] = device.Region
	transformed[neDeviceSchemaNames["HostName"]] = device.HostName
	transformed[neDeviceSchemaNames["LicenseFileID"]] = device.LicenseFileID
	transformed[neDeviceSchemaNames["LicenseFile"]] = device.LicenseFile
	transformed[neDeviceSchemaNames["LicenseToken"]] = device.LicenseToken
	transformed[neDeviceSchemaNames["CloudInitFileID"]] = device.CloudInitFileID
	transformed[neDeviceSchemaNames["ACLTemplateUUID"]] = device.ACLTemplateUUID
	transformed[neDeviceSchemaNames["SSHIPAddress"]] = device.SSHIPAddress
	transformed[neDeviceSchemaNames["SSHIPFqdn"]] = device.SSHIPFqdn
	transformed[neDeviceSchemaNames["AccountNumber"]] = device.AccountNumber
	transformed[neDeviceSchemaNames["Notifications"]] = device.Notifications
	transformed[neDeviceSchemaNames["RedundancyType"]] = device.RedundancyType
	transformed[neDeviceSchemaNames["ProjectID"]] = device.ProjectID
	transformed[neDeviceSchemaNames["RedundantUUID"]] = device.RedundantUUID
	transformed[neDeviceSchemaNames["AdditionalBandwidth"]] = device.AdditionalBandwidth
	transformed[neDeviceSchemaNames["Interfaces"]] = flattenNetworkDeviceInterfaces(device.Interfaces)
	transformed[neDeviceSchemaNames["VendorConfiguration"]] = device.VendorConfiguration
	transformed[neDeviceSchemaNames["UserPublicKey"]] = flattenNetworkDeviceUserKeys([]*ne.DeviceUserPublicKey{device.UserPublicKey})
	transformed[neDeviceSchemaNames["ASN"]] = device.ASN
	transformed[neDeviceSchemaNames["ZoneCode"]] = device.ZoneCode
	return []interface{}{transformed}
}

func expandNetworkDeviceSecondary(devices []interface{}) *ne.Device {
	if len(devices) < 1 {
		log.Printf("[WARN] resource_network_device expanding empty secondary device collection")
		return nil
	}
	device := devices[0].(map[string]interface{})
	transformed := &ne.Device{}
	if v, ok := device[neDeviceSchemaNames["UUID"]]; ok && !comparisons.IsEmpty(v) {
		transformed.UUID = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["Name"]]; ok && !comparisons.IsEmpty(v) {
		transformed.Name = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["ProjectID"]]; ok && !comparisons.IsEmpty(v) {
		transformed.ProjectID = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["MetroCode"]]; ok && !comparisons.IsEmpty(v) {
		transformed.MetroCode = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["HostName"]]; ok && !comparisons.IsEmpty(v) {
		transformed.HostName = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["LicenseToken"]]; ok && !comparisons.IsEmpty(v) {
		transformed.LicenseToken = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["LicenseFile"]]; ok && !comparisons.IsEmpty(v) {
		transformed.LicenseFile = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["LicenseFileID"]]; ok && !comparisons.IsEmpty(v) {
		transformed.LicenseFileID = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["CloudInitFileID"]]; ok && !comparisons.IsEmpty(v) {
		transformed.CloudInitFileID = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["ACLTemplateUUID"]]; ok && !comparisons.IsEmpty(v) {
		transformed.ACLTemplateUUID = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["MgmtAclTemplateUuid"]]; ok && !comparisons.IsEmpty(v) {
		transformed.MgmtAclTemplateUuid = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["AccountNumber"]]; ok && !comparisons.IsEmpty(v) {
		transformed.AccountNumber = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["Notifications"]]; ok {
		transformed.Notifications = converters.SetToStringList(v.(*schema.Set))
	}
	if v, ok := device[neDeviceSchemaNames["AdditionalBandwidth"]]; ok && !comparisons.IsEmpty(v) {
		transformed.AdditionalBandwidth = ne.Int(v.(int))
	}
	if v, ok := device[neDeviceSchemaNames["WanInterfaceId"]]; ok && !comparisons.IsEmpty(v) {
		transformed.WanInterfaceId = ne.String(v.(string))
	}
	if v, ok := device[neDeviceSchemaNames["VendorConfiguration"]]; ok {
		transformed.VendorConfiguration = converters.InterfaceMapToStringMap(v.(map[string]interface{}))
	}
	if v, ok := device[neDeviceSchemaNames["UserPublicKey"]]; ok {
		userKeys := expandNetworkDeviceUserKeys(v.(*schema.Set))
		if len(userKeys) > 0 {
			transformed.UserPublicKey = userKeys[0]
		}
	}
	return transformed
}

func flattenNetworkDeviceInterfaces(interfaces []ne.DeviceInterface) interface{} {
	transformed := make([]interface{}, len(interfaces))
	for i := range interfaces {
		transformed[i] = map[string]interface{}{
			neDeviceInterfaceSchemaNames["ID"]:                interfaces[i].ID,
			neDeviceInterfaceSchemaNames["Name"]:              interfaces[i].Name,
			neDeviceInterfaceSchemaNames["Status"]:            interfaces[i].Status,
			neDeviceInterfaceSchemaNames["OperationalStatus"]: interfaces[i].OperationalStatus,
			neDeviceInterfaceSchemaNames["MACAddress"]:        interfaces[i].MACAddress,
			neDeviceInterfaceSchemaNames["IPAddress"]:         interfaces[i].IPAddress,
			neDeviceInterfaceSchemaNames["AssignedType"]:      interfaces[i].AssignedType,
			neDeviceInterfaceSchemaNames["Type"]:              interfaces[i].Type,
		}
	}
	return transformed
}

func flattenNetworkDeviceUserKeys(userKeys []*ne.DeviceUserPublicKey) interface{} {
	transformed := make([]interface{}, 0, len(userKeys))
	for i := range userKeys {
		if userKeys[i] != nil {
			transformed = append(transformed, map[string]interface{}{
				neDeviceUserKeySchemaNames["Username"]: userKeys[i].Username,
				neDeviceUserKeySchemaNames["KeyName"]:  userKeys[i].KeyName,
			})
		}
	}
	return transformed
}

func expandNetworkDeviceUserKeys(userKeys *schema.Set) []*ne.DeviceUserPublicKey {
	userKeysList := userKeys.List()
	transformed := make([]*ne.DeviceUserPublicKey, len(userKeysList))
	for i := range userKeysList {
		userKeyMap := userKeysList[i].(map[string]interface{})
		transformed[i] = &ne.DeviceUserPublicKey{
			Username: ne.String(userKeyMap[neDeviceUserKeySchemaNames["Username"]].(string)),
			KeyName:  ne.String(userKeyMap[neDeviceUserKeySchemaNames["KeyName"]].(string)),
		}
	}
	return transformed
}

func flattenNetworkDeviceClusterDetails(clusterDetails *ne.ClusterDetails) interface{} {
	transformed := make(map[string]interface{})
	transformed[neDeviceClusterSchemaNames["ClusterId"]] = clusterDetails.ClusterId
	transformed[neDeviceClusterSchemaNames["ClusterName"]] = clusterDetails.ClusterName
	transformed[neDeviceClusterSchemaNames["NumOfNodes"]] = clusterDetails.NumOfNodes
	transformed[neDeviceClusterSchemaNames["Node0"]] = flattenNetworkDeviceClusterNodeDetail(clusterDetails.Node0)
	transformed[neDeviceClusterSchemaNames["Node1"]] = flattenNetworkDeviceClusterNodeDetail(clusterDetails.Node1)
	return []interface{}{transformed}
}

func flattenNetworkDeviceClusterNodeDetail(clusterNodeDetail *ne.ClusterNodeDetail) interface{} {
	transformed := make(map[string]interface{})
	transformed[neDeviceClusterNodeSchemaNames["UUID"]] = clusterNodeDetail.UUID
	transformed[neDeviceClusterNodeSchemaNames["Name"]] = clusterNodeDetail.Name
	transformed[neDeviceClusterNodeSchemaNames["VendorConfiguration"]] = flattenVendorConfiguration(clusterNodeDetail.VendorConfiguration)
	transformed[neDeviceClusterNodeSchemaNames["LicenseFileId"]] = clusterNodeDetail.LicenseFileId
	transformed[neDeviceClusterNodeSchemaNames["LicenseToken"]] = clusterNodeDetail.LicenseToken
	return []interface{}{transformed}
}

func flattenVendorConfiguration(vendorConfig map[string]string) interface{} {
	transformed := make(map[string]interface{})
	if v, ok := vendorConfig["hostname"]; ok {
		transformed[neDeviceVendorConfigSchemaNames["Hostname"]] = v
	}
	if v, ok := vendorConfig["adminPassword"]; ok {
		transformed[neDeviceVendorConfigSchemaNames["AdminPassword"]] = v
	}
	if v, ok := vendorConfig["controller1"]; ok {
		transformed[neDeviceVendorConfigSchemaNames["Controller1"]] = v
	}
	if v, ok := vendorConfig["activationKey"]; ok {
		transformed[neDeviceVendorConfigSchemaNames["ActivationKey"]] = v
	}
	if v, ok := vendorConfig["controllerFqdn"]; ok {
		transformed[neDeviceVendorConfigSchemaNames["ControllerFqdn"]] = v
	}
	if v, ok := vendorConfig["rootPassword"]; ok {
		transformed[neDeviceVendorConfigSchemaNames["RootPassword"]] = v
	}
	return []interface{}{transformed}
}

func expandNetworkDeviceClusterDetails(clusterDetails []interface{}) *ne.ClusterDetails {
	if len(clusterDetails) < 1 {
		log.Printf("[WARN] resource_network_device expanding empty cluster details")
		return nil
	}
	clusterDetail := clusterDetails[0].(map[string]interface{})
	transformed := &ne.ClusterDetails{}
	if v, ok := clusterDetail[neDeviceClusterSchemaNames["ClusterName"]]; ok && !comparisons.IsEmpty(v) {
		transformed.ClusterName = ne.String(v.(string))
	}
	if v, ok := clusterDetail[neDeviceClusterSchemaNames["Node0"]]; ok {
		transformed.Node0 = expandNetworkDeviceClusterNodeDetail(v.([]interface{}))
	}
	if v, ok := clusterDetail[neDeviceClusterSchemaNames["Node1"]]; ok {
		transformed.Node1 = expandNetworkDeviceClusterNodeDetail(v.([]interface{}))
	}
	return transformed
}

func expandNetworkDeviceClusterNodeDetail(clusterNodeDetails []interface{}) *ne.ClusterNodeDetail {
	if len(clusterNodeDetails) < 1 {
		log.Printf("[WARN] resource_network_device expanding empty cluster node details")
		return nil
	}
	clusterNodeDetail := clusterNodeDetails[0].(map[string]interface{})
	transformed := &ne.ClusterNodeDetail{}
	if v, ok := clusterNodeDetail[neDeviceClusterNodeSchemaNames["VendorConfiguration"]]; ok {
		transformed.VendorConfiguration = expandVendorConfiguration(v.([]interface{}))
	}
	if v, ok := clusterNodeDetail[neDeviceClusterNodeSchemaNames["LicenseFileId"]]; ok && !comparisons.IsEmpty(v) {
		transformed.LicenseFileId = ne.String(v.(string))
	}
	if v, ok := clusterNodeDetail[neDeviceClusterNodeSchemaNames["LicenseToken"]]; ok && !comparisons.IsEmpty(v) {
		transformed.LicenseToken = ne.String(v.(string))
	}
	return transformed
}

func expandVendorConfiguration(vendorConfigs []interface{}) map[string]string {
	if len(vendorConfigs) < 1 {
		log.Printf("[WARN] resource_network_device expanding empty vendor configurations")
		return nil
	}
	vendorConfig := vendorConfigs[0].(map[string]interface{})
	transformed := make(map[string]string)
	if v, ok := vendorConfig[neDeviceVendorConfigSchemaNames["Hostname"]]; ok && !comparisons.IsEmpty(v) {
		transformed["hostname"] = v.(string)
	}
	if v, ok := vendorConfig[neDeviceVendorConfigSchemaNames["AdminPassword"]]; ok && !comparisons.IsEmpty(v) {
		transformed["adminPassword"] = v.(string)
	}
	if v, ok := vendorConfig[neDeviceVendorConfigSchemaNames["Controller1"]]; ok && !comparisons.IsEmpty(v) {
		transformed["controller1"] = v.(string)
	}
	if v, ok := vendorConfig[neDeviceVendorConfigSchemaNames["ActivationKey"]]; ok && !comparisons.IsEmpty(v) {
		transformed["activationKey"] = v.(string)
	}
	if v, ok := vendorConfig[neDeviceVendorConfigSchemaNames["ControllerFqdn"]]; ok && !comparisons.IsEmpty(v) {
		transformed["controllerFqdn"] = v.(string)
	}
	if v, ok := vendorConfig[neDeviceVendorConfigSchemaNames["RootPassword"]]; ok && !comparisons.IsEmpty(v) {
		transformed["rootPassword"] = v.(string)
	}
	return transformed
}

func fillNetworkDeviceUpdateRequest(updateReq ne.DeviceUpdateRequest, changes map[string]interface{}) ne.DeviceUpdateRequest {
	for change, changeValue := range changes {
		switch change {
		case neDeviceSchemaNames["Name"]:
			updateReq.WithDeviceName(changeValue.(string))
		case neDeviceSchemaNames["TermLength"]:
			updateReq.WithTermLength(changeValue.(int))
		case neDeviceSchemaNames["Notifications"]:
			updateReq.WithNotifications(converters.SetToStringList(changeValue.(*schema.Set)))
		case neDeviceSchemaNames["CoreCount"]:
			updateReq.WithCore(changeValue.(int))
		case neDeviceSchemaNames["AdditionalBandwidth"]:
			updateReq.WithAdditionalBandwidth(changeValue.(int))
		case neDeviceSchemaNames["ACLTemplateUUID"]:
			updateReq.WithACLTemplate(changeValue.(string))
		case neDeviceSchemaNames["MgmtAclTemplateUuid"]:
			updateReq.WithMgmtAclTemplate(changeValue.(string))
		case neDeviceClusterSchemaNames["ClusterName"]:
			updateReq.WithClusterName(changeValue.(string))
		}
	}
	return updateReq
}

func getNetworkDeviceStateChangeConfigs(c ne.Client, deviceID string, timeout time.Duration, changes map[string]interface{}) []*retry.StateChangeConf {
	configs := make([]*retry.StateChangeConf, 0, len(changes))
	if changeValue, found := changes[neDeviceSchemaNames["ACLTemplateUUID"]]; found {
		aclTemplateUuid, ok := changeValue.(string)
		if ok && aclTemplateUuid != "" {
			configs = append(configs,
				createNetworkDeviceACLStatusWaitConfiguration(c.GetDeviceACLDetails, deviceID, 1*time.Second, timeout),
			)
		}
	} else if changeValue, found := changes[neDeviceSchemaNames["MgmtAclTemplateUuid"]]; found {
		mgmtAclTemplateUuid, ok := changeValue.(string)
		if ok && mgmtAclTemplateUuid != "" {
			configs = append(configs,
				createNetworkDeviceACLStatusWaitConfiguration(c.GetDeviceACLDetails, deviceID, 1*time.Second, timeout),
			)
		}
	}
	if _, found := changes[neDeviceSchemaNames["AdditionalBandwidth"]]; found {
		configs = append(configs,
			createNetworkDeviceAdditionalBandwidthStatusWaitConfiguration(c.GetDeviceAdditionalBandwidthDetails, deviceID, 1*time.Second, timeout),
		)
	}
	if _, found := changes[neDeviceSchemaNames["CoreCount"]]; found {
		configs = append(configs,
			createNetworkDeviceStatusResourceUpgradeWaitConfiguration(c.GetDevice, deviceID, 5*time.Second, timeout),
		)
	}
	return configs
}

type (
	openFile          func(name string) (*os.File, error)
	uploadLicenseFile func(metroCode, deviceTypeCode, deviceManagementMode, licenseMode, fileName string, reader io.Reader) (*string, error)
)

func uploadDeviceLicenseFile(openFunc openFile, uploadFunc uploadLicenseFile, typeCode string, device *ne.Device) error {
	if device == nil || ne.StringValue(device.LicenseFile) == "" {
		return nil
	}
	fileName := filepath.Base(ne.StringValue(device.LicenseFile))
	file, err := openFunc(ne.StringValue(device.LicenseFile))
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("[WARN] could not close file %q due to an error: %s", ne.StringValue(device.LicenseFile), err)
		}
	}()
	fileID, err := uploadFunc(ne.StringValue(device.MetroCode), typeCode, ne.DeviceManagementTypeSelf, ne.DeviceLicenseModeBYOL, fileName, file)
	if err != nil {
		return err
	}
	device.LicenseFileID = fileID
	return nil
}

type (
	getDevice                     func(uuid string) (*ne.Device, error)
	getACL                        func(uuid string) (*ne.DeviceACLDetails, error)
	getAdditionalBandwidthDetails func(uuid string) (*ne.DeviceAdditionalBandwidthDetails, error)
)

func createNetworkDeviceStatusProvisioningWaitConfiguration(fetchFunc getDevice, id string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	pending := []string{
		ne.DeviceStateInitializing,
		ne.DeviceStateProvisioning,
		ne.DeviceStateWaitingSecondary,
		ne.DeviceStateWaitingClusterNodes,
		ne.DeviceStateClusterSetUpInProgress,
	}
	target := []string{
		ne.DeviceStateProvisioned,
	}
	return createNetworkDeviceStatusWaitConfiguration(fetchFunc, id, delay, timeout, target, pending)
}

func createNetworkDeviceStatusDeleteWaitConfiguration(fetchFunc getDevice, id string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	pending := []string{
		ne.DeviceStateDeprovisioning,
	}
	target := []string{
		ne.DeviceStateDeprovisioned,
	}
	return createNetworkDeviceStatusWaitConfiguration(fetchFunc, id, delay, timeout, target, pending)
}

func createNetworkDeviceStatusResourceUpgradeWaitConfiguration(fetchFunc getDevice, id string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	pending := []string{
		ne.DeviceStateResourceUpgradeInProgress,
		ne.DeviceStateWaitingPrimary,
		ne.DeviceStateWaitingSecondary,
		ne.DeviceStateWaitingClusterNodes,
	}
	target := []string{
		ne.DeviceStateProvisioned,
	}
	return createNetworkDeviceStatusWaitConfiguration(fetchFunc, id, delay, timeout, target, pending)
}

func createNetworkDeviceStatusWaitConfiguration(fetchFunc getDevice, id string, delay time.Duration, timeout time.Duration, target []string, pending []string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: delay,
		Refresh: func() (interface{}, string, error) {
			resp, err := fetchFunc(id)
			if err != nil {
				return nil, "", err
			}
			return resp, ne.StringValue(resp.Status), nil
		},
	}
}

func createNetworkDeviceLicenseStatusWaitConfiguration(fetchFunc getDevice, id string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	pending := []string{
		ne.DeviceLicenseStateApplying,
		ne.DeviceLicenseStateWaitingClusterSetUp,
		"",
	}
	target := []string{
		ne.DeviceLicenseStateRegistered,
		ne.DeviceLicenseStateApplied,
		ne.DeviceLicenseStateNA,
	}
	return &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: delay,
		Refresh: func() (interface{}, string, error) {
			resp, err := fetchFunc(id)
			if err != nil {
				return nil, "", err
			}
			return resp, ne.StringValue(resp.LicenseStatus), nil
		},
	}
}

func createNetworkDeviceACLStatusWaitConfiguration(fetchFunc getACL, id string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			ne.ACLDeviceStatusProvisioning,
		},
		Target: []string{
			ne.ACLDeviceStatusProvisioned,
		},
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: delay,
		Refresh: func() (interface{}, string, error) {
			resp, err := fetchFunc(id)
			if err != nil {
				return nil, "", err
			}
			return resp, ne.StringValue(resp.Status), nil
		},
	}
}

func createNetworkDeviceAdditionalBandwidthStatusWaitConfiguration(fetchFunc getAdditionalBandwidthDetails, deviceID string, delay time.Duration, timeout time.Duration) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			ne.DeviceAdditionalBandwidthStatusProvisioning,
		},
		Target: []string{
			ne.DeviceAdditionalBandwidthStatusProvisioned,
		},
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: delay,
		Refresh: func() (interface{}, string, error) {
			resp, err := fetchFunc(deviceID)
			if err != nil {
				return nil, "", err
			}
			return resp, ne.StringValue(resp.Status), nil
		},
	}
}
