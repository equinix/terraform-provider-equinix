package equinix

import (
	"context"
	"fmt"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var neDeviceStates = []string{ // Not sure if other states should be included
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
					neDeviceClusterSchemaNames["ClusterId"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceClusterDescriptions["ClusterId"],
					},
					neDeviceClusterSchemaNames["ClusterName"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceClusterDescriptions["ClusterName"],
					},
					neDeviceClusterSchemaNames["NumOfNodes"]: {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: neDeviceClusterDescriptions["NumOfNodes"],
					},
					neDeviceClusterSchemaNames["Node0"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createDataSourceClusterNodeDetailSchema(),
						},
						Description: neDeviceClusterDescriptions["Node0"],
					},
					neDeviceClusterSchemaNames["Node1"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createDataSourceClusterNodeDetailSchema(),
						},
						Description: neDeviceClusterDescriptions["Node1"],
					},
				},
			},
		},
	}
}

func createDataSourceNetworkDeviceInterfaceSchema() map[string]*schema.Schema {
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

func createDataSourceNetworkDeviceUserKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceUserKeySchemaNames["Username"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceUserKeyDescriptions["Username"],
		},
		neDeviceUserKeySchemaNames["KeyName"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceUserKeyDescriptions["KeyName"],
		},
	}
}

func createDataSourceClusterNodeDetailSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceClusterNodeSchemaNames["LicenseFileId"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: neDeviceClusterNodeDescriptions["LicenseFileId"],
		},
		neDeviceClusterNodeSchemaNames["LicenseToken"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceClusterNodeDescriptions["LicenseToken"],
		},
		neDeviceClusterNodeSchemaNames["VendorConfiguration"]: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createDataSourceVendorConfigurationSchema(),
			},
			Description: neDeviceClusterNodeDescriptions["VendorConfiguration"],
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

func createDataSourceVendorConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceVendorConfigSchemaNames["Hostname"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceVendorConfigDescriptions["Hostname"],
		},
		neDeviceVendorConfigSchemaNames["AdminPassword"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: neDeviceVendorConfigDescriptions["AdminPassword"],
		},
		neDeviceVendorConfigSchemaNames["Controller1"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceVendorConfigDescriptions["Controller1"],
		},
		neDeviceVendorConfigSchemaNames["ActivationKey"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: neDeviceVendorConfigDescriptions["ActivationKey"],
		},
		neDeviceVendorConfigSchemaNames["ControllerFqdn"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceVendorConfigDescriptions["ControllerFqdn"],
		},
		neDeviceVendorConfigSchemaNames["RootPassword"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: neDeviceVendorConfigDescriptions["RootPassword"],
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
	devices, err = conf.ne.GetDevices(neDeviceStates)
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
