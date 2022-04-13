package equinix

import (
	"context"
	"fmt"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var deviceStates = []string{ // Not sure if other states should be included
	// ne.DeviceStateClusterSetUpInProgress,
	// ne.DeviceStateDeprovisioned,
	// ne.DeviceStateDeprovisioning,
	// ne.DeviceStateFailed,
	// ne.DeviceStateInitializing,
	ne.DeviceStateProvisioned,
	// ne.DeviceStateProvisioning,
	// ne.DeviceStateWaitingClusterNodes,
	// ne.DeviceStateWaitingPrimary,
	// ne.DeviceStateWaitingSecondary,
}

var ecxNetworkDeviceSchemaNames = map[string]string{
	"UUID":                "uuid",
	"Name":                "name",
	"TypeCode":            "type_code",
	"Status":              "status",
	"MetroCode":           "metro_code",
	"IBX":                 "ibx",
	"Region":              "region",
	"Throughput":          "throughput",
	"ThroughputUnit":      "throughput_unit",
	"HostName":            "hostname",
	"PackageCode":         "package_code",
	"Version":             "version",
	"IsBYOL":              "byol",
	"LicenseToken":        "license_token",
	"LicenseFile":         "license_file",
	"LicenseFileID":       "license_file_id",
	"LicenseStatus":       "license_status",
	"ACLTemplateUUID":     "acl_template_id",
	"MgmtAclTemplateUuid": "mgmt_acl_template_uuid",
	"SSHIPAddress":        "ssh_ip_address",
	"SSHIPFqdn":           "ssh_ip_fqdn",
	"AccountNumber":       "account_number",
	"Notifications":       "notifications",
	"PurchaseOrderNumber": "purchase_order_number",
	"RedundancyType":      "redundancy_type",
	"RedundantUUID":       "redundant_id",
	"TermLength":          "term_length",
	"AdditionalBandwidth": "additional_bandwidth",
	"OrderReference":      "order_reference",
	"InterfaceCount":      "interface_count",
	"CoreCount":           "core_count",
	"IsSelfManaged":       "self_managed",
	"WanInterfaceId":      "wan_interface_id",
	"Interfaces":          "interface",
	"VendorConfiguration": "vendor_configuration",
	"UserPublicKey":       "ssh_key",
	"ASN":                 "asn",
	"ZoneCode":            "zone_code",
	"Secondary":           "secondary_device",
	"ClusterDetails":      "cluster_details",
}

var ecxNetworkDeviceDescriptions = map[string]string{
	"UUID":                "Device unique identifier",
	"Name":                "Device name",
	"TypeCode":            "Device type code",
	"Status":              "Device provisioning status",
	"MetroCode":           "Device location metro code",
	"IBX":                 "Device location Equinix Business Exchange name",
	"Region":              "Device location region",
	"Throughput":          "Device license throughput",
	"ThroughputUnit":      "Device license throughput unit (Mbps or Gbps)",
	"HostName":            "Device hostname prefix",
	"PackageCode":         "Device software package code",
	"Version":             "Device software software version",
	"IsBYOL":              "Boolean value that determines device licensing mode: bring your own license or subscription (default)",
	"LicenseToken":        "License Token applicable for some device types in BYOL licensing mode",
	"LicenseFile":         "Path to the license file that will be uploaded and applied on a device, applicable for some device types in BYOL licensing mode",
	"LicenseFileID":       "Unique identifier of applied license file",
	"LicenseStatus":       "Device license registration status",
	"ACLTemplateUUID":     "Unique identifier of applied ACL template",
	"MgmtAclTemplateUuid": "Unique identifier of applied MGMT ACL template",
	"SSHIPAddress":        "IP address of SSH enabled interface on the device",
	"SSHIPFqdn":           "FQDN of SSH enabled interface on the device",
	"AccountNumber":       "Device billing account number",
	"Notifications":       "List of email addresses that will receive device status notifications",
	"PurchaseOrderNumber": "Purchase order number associated with a device order",
	"RedundancyType":      "Device redundancy type applicable for HA devices, either primary or secondary",
	"RedundantUUID":       "Unique identifier for a redundant device, applicable for HA device",
	"TermLength":          "Device term length",
	"AdditionalBandwidth": "Additional Internet bandwidth, in Mbps, that will be allocated to the device",
	"OrderReference":      "Name/number used to identify device order on the invoice",
	"InterfaceCount":      "Number of network interfaces on a device. If not specified, default number for a given device type will be used",
	"CoreCount":           "Number of CPU cores used by device",
	"IsSelfManaged":       "Boolean value that determines device management mode: self-managed or subscription (default)",
	"WanInterfaceId":      "device interface id picked for WAN",
	"Interfaces":          "List of device interfaces",
	"VendorConfiguration": "Map of vendor specific configuration parameters for a device (controller1, activationKey, managementType, siteId, systemIpAddress)",
	"UserPublicKey":       "Definition of SSH key that will be provisioned on a device",
	"ASN":                 "Autonomous system number",
	"ZoneCode":            "Device location zone code",
	"Secondary":           "Definition of secondary device applicable for HA setup",
	"ClusterDetails":      "An object that has the cluster details",
}

var ecxDeviceInterfaceSchemaNames = map[string]string{
	"ID":                "id",
	"Name":              "name",
	"Status":            "status",
	"OperationalStatus": "operational_status",
	"MACAddress":        "mac_address",
	"IPAddress":         "ip_address",
	"AssignedType":      "assigned_type",
	"Type":              "type",
}

var ecxDeviceInterfaceDescriptions = map[string]string{
	"ID":                "Interface identifier",
	"Name":              "Interface name",
	"Status":            "Interface status (AVAILABLE, RESERVED, ASSIGNED)",
	"OperationalStatus": "Interface operational status (up or down)",
	"MACAddress":        "Interface MAC addres",
	"IPAddress":         "interface IP address",
	"AssignedType":      "Interface management type (Equinix Managed or empty)",
	"Type":              "Interface type",
}

var ecxDeviceUserKeySchemaNames = map[string]string{
	"Username": "username",
	"KeyName":  "key_name",
}

var ecxDeviceUserKeyDescriptions = map[string]string{
	"Username": "Username associated with given key",
	"KeyName":  "Reference by name to previously provisioned public SSH key",
}

var ecxDeviceClusterSchemaNames = map[string]string{
	"ClusterId":   "cluster_id",
	"ClusterName": "cluster_name",
	"NumOfNodes":  "num_of_nodes",
	"Node0":       "node0",
	"Node1":       "node1",
}

var ecxDeviceClusterDescriptions = map[string]string{
	"ClusterId":   "The id of the cluster",
	"ClusterName": "The name of the cluster device",
	"NumOfNodes":  "The number of nodes in the cluster",
	"Node0":       "An object that has node0 details",
	"Node1":       "An object that has node1 details",
}

var ecxDeviceClusterNodeSchemaNames = map[string]string{
	"VendorConfiguration": "vendor_configuration",
	"LicenseFileId":       "license_file_id",
	"LicenseToken":        "license_token",
	"UUID":                "uuid",
	"Name":                "name",
}

var ecxDeviceClusterNodeDescriptions = map[string]string{
	"VendorConfiguration": "An object that has fields relevant to the vendor of the cluster device",
	"LicenseFileId":       "License file id. This is necessary for Fortinet and Juniper clusters",
	"LicenseToken":        "License token. This is necessary for Palo Alto clusters",
	"UUID":                "The unique id of the node",
	"Name":                "The name of the node",
}

var ecxDeviceVendorConfigSchemaNames = map[string]string{
	"Hostname":       "hostname",
	"AdminPassword":  "admin_password",
	"Controller1":    "controller1",
	"ActivationKey":  "activation_key",
	"ControllerFqdn": "controller_fqdn",
	"RootPassword":   "root_password",
}

var ecxDeviceVendorConfigDescriptions = map[string]string{
	"Hostname":       "Hostname. This is necessary for Palo Alto, Juniper, and Fortinet clusters",
	"AdminPassword":  "The administrative password of the device. You can use it to log in to the console. This field is not available for all device types",
	"Controller1":    "System IP Address. Mandatory for the Fortinet SDWAN cluster device",
	"ActivationKey":  "Activation key. This is required for Velocloud clusters",
	"ControllerFqdn": "Controller fqdn. This is required for Velocloud clusters",
	"RootPassword":   "The CLI password of the device. This field is relevant only for the Velocloud SDWAN cluster",
}

func createDataSourceNetworkDeviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkDeviceSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: networkDeviceDescriptions["UUID"],
		},
		networkDeviceSchemaNames["Name"]: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return old == new+"-Node0"
			},
			ValidateFunc: validation.StringLenBetween(3, 50),
			Description:  networkDeviceDescriptions["Name"],
		},
		networkDeviceSchemaNames["TypeCode"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["TypeCode"],
		},
		networkDeviceSchemaNames["Status"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["Status"],
		},
		networkDeviceSchemaNames["LicenseStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["LicenseStatus"],
		},
		networkDeviceSchemaNames["MetroCode"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["MetroCode"],
		},
		networkDeviceSchemaNames["IBX"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["IBX"],
		},
		networkDeviceSchemaNames["Region"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["Region"],
		},
		networkDeviceSchemaNames["Throughput"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: networkDeviceDescriptions["Throughput"],
		},
		networkDeviceSchemaNames["ThroughputUnit"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["ThroughputUnit"],
		},
		networkDeviceSchemaNames["HostName"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["HostName"],
		},
		networkDeviceSchemaNames["PackageCode"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["PackageCode"],
		},
		networkDeviceSchemaNames["Version"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["Version"],
		},
		networkDeviceSchemaNames["IsBYOL"]: {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: networkDeviceDescriptions["IsBYOL"],
		},
		networkDeviceSchemaNames["LicenseToken"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["LicenseToken"],
		},
		networkDeviceSchemaNames["LicenseFile"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["LicenseFile"],
		},
		networkDeviceSchemaNames["LicenseFileID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["LicenseFileID"],
		},
		networkDeviceSchemaNames["ACLTemplateUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["ACLTemplateUUID"],
		},
		networkDeviceSchemaNames["MgmtAclTemplateUuid"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["MgmtAclTemplateUuid"],
		},
		networkDeviceSchemaNames["SSHIPAddress"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["SSHIPAddress"],
		},
		networkDeviceSchemaNames["SSHIPFqdn"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["SSHIPFqdn"],
		},
		networkDeviceSchemaNames["AccountNumber"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["AccountNumber"],
		},
		networkDeviceSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: networkDeviceDescriptions["Notifications"],
		},
		networkDeviceSchemaNames["PurchaseOrderNumber"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["PurchaseOrderNumber"],
		},
		networkDeviceSchemaNames["RedundancyType"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["RedundancyType"],
		},
		networkDeviceSchemaNames["RedundantUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["RedundantUUID"],
		},
		networkDeviceSchemaNames["TermLength"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: networkDeviceDescriptions["TermLength"],
		},
		networkDeviceSchemaNames["AdditionalBandwidth"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: networkDeviceDescriptions["AdditionalBandwidth"],
		},
		networkDeviceSchemaNames["OrderReference"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["OrderReference"],
		},
		networkDeviceSchemaNames["InterfaceCount"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: networkDeviceDescriptions["InterfaceCount"],
		},
		networkDeviceSchemaNames["CoreCount"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: networkDeviceDescriptions["CoreCount"],
		},
		networkDeviceSchemaNames["IsSelfManaged"]: {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: networkDeviceDescriptions["IsSelfManaged"],
		},
		networkDeviceSchemaNames["WanInterfaceId"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["WanInterfaceId"],
		},
		networkDeviceSchemaNames["Interfaces"]: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createDataSourceNetworkDeviceInterfaceSchema(),
			},
			Description: networkDeviceDescriptions["Interfaces"],
		},
		networkDeviceSchemaNames["VendorConfiguration"]: {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: networkDeviceDescriptions["VendorConfiguration"],
		},
		networkDeviceSchemaNames["UserPublicKey"]: {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createDataSourceNetworkDeviceUserKeySchema(),
			},
			Description: networkDeviceDescriptions["UserPublicKey"],
		},
		networkDeviceSchemaNames["ASN"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: networkDeviceDescriptions["ASN"],
		},
		networkDeviceSchemaNames["ZoneCode"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkDeviceDescriptions["ZoneCode"],
		},
		networkDeviceSchemaNames["Secondary"]: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: networkDeviceDescriptions["Secondary"],
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					networkDeviceSchemaNames["UUID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["UUID"],
					},
					networkDeviceSchemaNames["Name"]: {
						Type:        schema.TypeString,
						Required:    true,
						Description: networkDeviceDescriptions["Name"],
					},
					networkDeviceSchemaNames["Status"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["Status"],
					},
					networkDeviceSchemaNames["LicenseStatus"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["LicenseStatus"],
					},
					networkDeviceSchemaNames["MetroCode"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["MetroCode"],
					},
					networkDeviceSchemaNames["IBX"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["IBX"],
					},
					networkDeviceSchemaNames["Region"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["Region"],
					},
					networkDeviceSchemaNames["HostName"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["HostName"],
					},
					networkDeviceSchemaNames["LicenseToken"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["LicenseToken"],
					},
					networkDeviceSchemaNames["LicenseFile"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["LicenseFile"],
					},
					networkDeviceSchemaNames["LicenseFileID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["LicenseFileID"],
					},
					networkDeviceSchemaNames["ACLTemplateUUID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["ACLTemplateUUID"],
					},
					networkDeviceSchemaNames["MgmtAclTemplateUuid"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["MgmtAclTemplateUuid"],
					},
					networkDeviceSchemaNames["SSHIPAddress"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["SSHIPAddress"],
					},
					networkDeviceSchemaNames["SSHIPFqdn"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["SSHIPFqdn"],
					},
					networkDeviceSchemaNames["AccountNumber"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["AccountNumber"],
					},
					networkDeviceSchemaNames["Notifications"]: {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Description: networkDeviceDescriptions["Notifications"],
					},
					networkDeviceSchemaNames["RedundancyType"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["RedundancyType"],
					},
					networkDeviceSchemaNames["RedundantUUID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["RedundantUUID"],
					},
					networkDeviceSchemaNames["AdditionalBandwidth"]: {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: networkDeviceDescriptions["AdditionalBandwidth"],
					},
					networkDeviceSchemaNames["WanInterfaceId"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["WanInterfaceId"],
					},
					networkDeviceSchemaNames["Interfaces"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createDataSourceNetworkDeviceInterfaceSchema(),
						},
						Description: networkDeviceDescriptions["Interfaces"],
					},
					networkDeviceSchemaNames["VendorConfiguration"]: {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						Description: networkDeviceDescriptions["VendorConfiguration"],
					},
					networkDeviceSchemaNames["UserPublicKey"]: {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createDataSourceNetworkDeviceUserKeySchema(),
						},
						Description: networkDeviceDescriptions["UserPublicKey"],
					},
					networkDeviceSchemaNames["ASN"]: {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: networkDeviceDescriptions["ASN"],
					},
					networkDeviceSchemaNames["ZoneCode"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: networkDeviceDescriptions["ZoneCode"],
					},
				},
			},
		},
		networkDeviceSchemaNames["ClusterDetails"]: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: networkDeviceDescriptions["ClusterDetails"],
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					ecxDeviceClusterSchemaNames["ClusterId"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: ecxDeviceClusterDescriptions["ClusterId"],
					},
					ecxDeviceClusterSchemaNames["ClusterName"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: ecxDeviceClusterDescriptions["ClusterName"],
					},
					ecxDeviceClusterSchemaNames["NumOfNodes"]: {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: ecxDeviceClusterDescriptions["NumOfNodes"],
					},
					ecxDeviceClusterSchemaNames["Node0"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createDataSourceClusterNodeDetailSchema(),
						},
						Description: ecxDeviceClusterDescriptions["Node0"],
					},
					ecxDeviceClusterSchemaNames["Node1"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createDataSourceClusterNodeDetailSchema(),
						},
						Description: ecxDeviceClusterDescriptions["Node1"],
					},
				},
			},
		},
	}
}

func createDataSourceNetworkDeviceInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxDeviceInterfaceSchemaNames["ID"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: ecxDeviceInterfaceDescriptions["ID"],
		},
		ecxDeviceInterfaceSchemaNames["Name"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceInterfaceDescriptions["Name"],
		},
		ecxDeviceInterfaceSchemaNames["Status"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceInterfaceDescriptions["Status"],
		},
		ecxDeviceInterfaceSchemaNames["OperationalStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceInterfaceDescriptions["OperationalStatus"],
		},
		ecxDeviceInterfaceSchemaNames["MACAddress"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceInterfaceDescriptions["MACAddress"],
		},
		ecxDeviceInterfaceSchemaNames["IPAddress"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceInterfaceDescriptions["IPAddress"],
		},
		ecxDeviceInterfaceSchemaNames["AssignedType"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceInterfaceDescriptions["AssignedType"],
		},
		ecxDeviceInterfaceSchemaNames["Type"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceInterfaceDescriptions["Type"],
		},
	}
}

func createDataSourceNetworkDeviceUserKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxDeviceUserKeySchemaNames["Username"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceUserKeyDescriptions["Username"],
		},
		ecxDeviceUserKeySchemaNames["KeyName"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceUserKeyDescriptions["KeyName"],
		},
	}
}

func createDataSourceClusterNodeDetailSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxDeviceClusterNodeSchemaNames["LicenseFileId"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: ecxDeviceClusterNodeDescriptions["LicenseFileId"],
		},
		ecxDeviceClusterNodeSchemaNames["LicenseToken"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceClusterNodeDescriptions["LicenseToken"],
		},
		ecxDeviceClusterNodeSchemaNames["VendorConfiguration"]: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createDataSourceVendorConfigurationSchema(),
			},
			Description: ecxDeviceClusterNodeDescriptions["VendorConfiguration"],
		},
		ecxDeviceClusterNodeSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceClusterNodeDescriptions["UUID"],
		},
		ecxDeviceClusterNodeSchemaNames["Name"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceClusterNodeDescriptions["Name"],
		},
	}
}

func createDataSourceVendorConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ecxDeviceVendorConfigSchemaNames["Hostname"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceVendorConfigDescriptions["Hostname"],
		},
		ecxDeviceVendorConfigSchemaNames["AdminPassword"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: ecxDeviceVendorConfigDescriptions["AdminPassword"],
		},
		ecxDeviceVendorConfigSchemaNames["Controller1"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceVendorConfigDescriptions["Controller1"],
		},
		ecxDeviceVendorConfigSchemaNames["ActivationKey"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: ecxDeviceVendorConfigDescriptions["ActivationKey"],
		},
		ecxDeviceVendorConfigSchemaNames["ControllerFqdn"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ecxDeviceVendorConfigDescriptions["ControllerFqdn"],
		},
		ecxDeviceVendorConfigSchemaNames["RootPassword"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: ecxDeviceVendorConfigDescriptions["RootPassword"],
		},
	}
}
func dataSourceNetworkDevice() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkDeviceRead,
		Description: "Use this data source to get details of Equinix Network Edge network device with a given Name or UUID",
		Schema:      createDataSourceNetworkDeviceSchema(),
	}
}

func getDeviceByName(deviceName string, conf *Config) (*ne.Device, error) {
	var devices []ne.Device
	err := error(nil)
	devices, err = conf.ne.GetDevices(deviceStates)
	if err != nil {
		return nil, fmt.Errorf("'devices: %v'", devices)
	}
	for _, device := range devices {
		// return nil, fmt.Errorf("device name is %s", *device.Name)
		if ne.StringValue(device.Name) == deviceName {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("device %s not found", deviceName)
}

func dataSourceNetworkDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	var err error
	var primary, secondary *ne.Device

	nameIf, nameExists := d.GetOk("name")
	name := nameIf.(string)
	uuidIf, uuidExists := d.GetOk("uuid")
	uuid := uuidIf.(string)

	if nameExists && uuidExists {
		return diag.Errorf("name and uuid arguments can't be used together")
	}

	if !nameExists && !uuidExists {
		return diag.Errorf("either name or uuid must be set")
	}

	if nameExists {
		primary, err = getDeviceByName(name, conf)
	} else {
		primary, err = conf.ne.GetDevice(uuid)
	}

	if err != nil {
		return diag.Errorf("cannot fetch primary network device due to '%v'", err)
	}

	if isStringInSlice(ne.StringValue(primary.Status), []string{ne.DeviceStateDeprovisioning, ne.DeviceStateDeprovisioned}) {
		d.SetId("")
		return diags
	}
	if ne.StringValue(primary.RedundantUUID) != "" {

		secondary, err = conf.ne.GetDevice(ne.StringValue(primary.RedundantUUID))
		if err != nil {
			return diag.Errorf("cannot fetch secondary network device due to '%v'", err)
		}
	}
	if err = updateDataSourceNetworkDeviceResource(primary, secondary, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateDataSourceNetworkDeviceResource(primary *ne.Device, secondary *ne.Device, d *schema.ResourceData) error {
	d.SetId(ne.StringValue(primary.UUID))
	if err := d.Set(networkDeviceSchemaNames["UUID"], primary.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["Name"], primary.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}

	if err := d.Set(networkDeviceSchemaNames["TypeCode"], primary.TypeCode); err != nil {
		return fmt.Errorf("error reading TypeCode: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["Status"], primary.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["LicenseStatus"], primary.LicenseStatus); err != nil {
		return fmt.Errorf("error reading LicenseStatus: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["MetroCode"], primary.MetroCode); err != nil {
		return fmt.Errorf("error reading MetroCode: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["IBX"], primary.IBX); err != nil {
		return fmt.Errorf("error reading IBX: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["Region"], primary.Region); err != nil {
		return fmt.Errorf("error reading Region: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["Throughput"], primary.Throughput); err != nil {
		return fmt.Errorf("error reading Throughput: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["ThroughputUnit"], primary.ThroughputUnit); err != nil {
		return fmt.Errorf("error reading ThroughputUnit: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["HostName"], primary.HostName); err != nil {
		return fmt.Errorf("error reading HostName: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["PackageCode"], primary.PackageCode); err != nil {
		return fmt.Errorf("error reading PackageCode: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["Version"], primary.Version); err != nil {
		return fmt.Errorf("error reading Version: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["IsBYOL"], primary.IsBYOL); err != nil {
		return fmt.Errorf("error reading IsBYOL: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["LicenseFileID"], primary.LicenseFileID); err != nil {
		return fmt.Errorf("error reading LicenseFileID: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["ACLTemplateUUID"], primary.ACLTemplateUUID); err != nil {
		return fmt.Errorf("error reading ACLTemplateUUID: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["SSHIPAddress"], primary.SSHIPAddress); err != nil {
		return fmt.Errorf("error reading SSHIPAddress: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["SSHIPFqdn"], primary.SSHIPFqdn); err != nil {
		return fmt.Errorf("error reading SSHIPFqdn: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["AccountNumber"], primary.AccountNumber); err != nil {
		return fmt.Errorf("error reading AccountNumber: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["Notifications"], primary.Notifications); err != nil {
		return fmt.Errorf("error reading Notifications: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["PurchaseOrderNumber"], primary.PurchaseOrderNumber); err != nil {
		return fmt.Errorf("error reading PurchaseOrderNumber: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["RedundancyType"], primary.RedundancyType); err != nil {
		return fmt.Errorf("error reading RedundancyType: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["RedundantUUID"], primary.RedundantUUID); err != nil {
		return fmt.Errorf("error reading RedundantUUID: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["TermLength"], primary.TermLength); err != nil {
		return fmt.Errorf("error reading TermLength: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["AdditionalBandwidth"], primary.AdditionalBandwidth); err != nil {
		return fmt.Errorf("error reading AdditionalBandwidth: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["OrderReference"], primary.OrderReference); err != nil {
		return fmt.Errorf("error reading OrderReference: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["InterfaceCount"], primary.InterfaceCount); err != nil {
		return fmt.Errorf("error reading InterfaceCount: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["CoreCount"], primary.CoreCount); err != nil {
		return fmt.Errorf("error reading CoreCount: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["IsSelfManaged"], primary.IsSelfManaged); err != nil {
		return fmt.Errorf("error reading IsSelfManaged: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["Interfaces"], flattenNetworkDeviceInterfaces(primary.Interfaces)); err != nil {
		return fmt.Errorf("error reading Interfaces: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["VendorConfiguration"], primary.VendorConfiguration); err != nil {
		return fmt.Errorf("error reading VendorConfiguration: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["UserPublicKey"], flattenNetworkDeviceUserKeys([]*ne.DeviceUserPublicKey{primary.UserPublicKey})); err != nil {
		return fmt.Errorf("error reading VendorConfiguration: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["ASN"], primary.ASN); err != nil {
		return fmt.Errorf("error reading ASN: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["ZoneCode"], primary.ZoneCode); err != nil {
		return fmt.Errorf("error reading ZoneCode: %s", err)
	}
	if secondary != nil {
		if v, ok := d.GetOk(networkDeviceSchemaNames["Secondary"]); ok {
			secondaryFromSchema := expandNetworkDeviceSecondary(v.([]interface{}))
			secondary.LicenseFile = secondaryFromSchema.LicenseFile
		}
		if err := d.Set(networkDeviceSchemaNames["Secondary"], flattenNetworkDeviceSecondary(secondary)); err != nil {
			return fmt.Errorf("error reading Secondary: %s", err)
		}
	}
	if primary.ClusterDetails != nil {
		if v, ok := d.GetOk(networkDeviceSchemaNames["ClusterDetails"]); ok {
			clusterDetailsFromSchema := expandNetworkDeviceClusterDetails(v.([]interface{}))
			primary.ClusterDetails.Node0.LicenseFileId = clusterDetailsFromSchema.Node0.LicenseFileId
			primary.ClusterDetails.Node0.LicenseToken = clusterDetailsFromSchema.Node0.LicenseToken
			primary.ClusterDetails.Node1.LicenseFileId = clusterDetailsFromSchema.Node1.LicenseFileId
			primary.ClusterDetails.Node1.LicenseToken = clusterDetailsFromSchema.Node1.LicenseToken
		}
		if err := d.Set(networkDeviceSchemaNames["ClusterDetails"], flattenNetworkDeviceClusterDetails(primary.ClusterDetails)); err != nil {
			return fmt.Errorf("error reading ClusterDetails: %s", err)
		}
	}
	return nil
}
