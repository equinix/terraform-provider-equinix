package device

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var neDeviceStateMap = map[string]string{
	"deprovisioned":                     ne.DeviceStateDeprovisioned,
	"failed":                            ne.DeviceStateFailed,
	"initializing":                      ne.DeviceStateInitializing,
	"provisioned":                       ne.DeviceStateProvisioned,
	"provisioning":                      ne.DeviceStateProvisioning,
	"cluster_setup_in_progress":         ne.DeviceStateClusterSetUpInProgress,
	"waiting_for_replica_cluster_nodes": ne.DeviceStateWaitingClusterNodes,
	"waiting_for_primary":               ne.DeviceStateWaitingPrimary,
	"waiting_for_secondary":             ne.DeviceStateWaitingSecondary,
	"resource_upgrade_in_progress":      ne.DeviceStateResourceUpgradeInProgress,
	"resource_upgrade_failed":           ne.DeviceStateResourceUpgradeFailed,
}

func getNeDeviceStatusList(deviceStateText string) (*[]string, error) {
	if deviceStateText == "" {
		return &[]string{}, nil
	}
	deviceStateTextSlice := strings.Split(strings.ToLower(deviceStateText), ",")
	invalidItems := []string{}
	validItems := []string{}
	for _, textItem := range deviceStateTextSlice {
		textItem = strings.TrimSpace(textItem)
		val, ok := neDeviceStateMap[textItem]
		if !ok {
			invalidItems = append(invalidItems, textItem)
		}
		validItems = append(validItems, val)
	}
	if len(invalidItems) > 0 {
		return nil, fmt.Errorf("Invalid Items: %v", invalidItems)
	}
	return &validItems, nil
}

func stringIsValidDeviceStateList(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}
	if _, err := getNeDeviceStatusList(v); err != nil {
		return nil, []error{fmt.Errorf("value of %v is not a valid list of device states, got error: %v", k, err)}
	}
	return nil, nil
}

func createDataSourceNetworkDeviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceSchemaNames["UUID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  neDeviceDescriptions["UUID"],
			ExactlyOneOf: []string{neDeviceSchemaNames["UUID"], neDeviceSchemaNames["Name"]},
		},
		neDeviceSchemaNames["Name"]: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return old == new+"-Node0"
			},
			ValidateFunc: validation.StringLenBetween(3, 50),
			Description:  neDeviceDescriptions["Name"],
			ExactlyOneOf: []string{neDeviceSchemaNames["UUID"], neDeviceSchemaNames["Name"]},
		},
		neDeviceSchemaNames["ValidStatusList"]: {
			Type:          schema.TypeString,
			Optional:      true,
			Default:       "Provisioned",
			Description:   neDeviceDescriptions["ValidStatusList"],
			ValidateFunc:  stringIsValidDeviceStateList,
			ConflictsWith: []string{neDeviceSchemaNames["UUID"]},
		},
		neDeviceSchemaNames["Status"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["Status"],
		},
		neDeviceSchemaNames["TypeCode"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["TypeCode"],
		},
		neDeviceSchemaNames["LicenseStatus"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["LicenseStatus"],
		},
		neDeviceSchemaNames["MetroCode"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["MetroCode"],
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
			Type:        schema.TypeInt,
			Computed:    true,
			Description: neDeviceDescriptions["Throughput"],
		},
		neDeviceSchemaNames["ThroughputUnit"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["ThroughputUnit"],
		},
		neDeviceSchemaNames["HostName"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["HostName"],
		},
		neDeviceSchemaNames["PackageCode"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["PackageCode"],
		},
		neDeviceSchemaNames["Version"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["Version"],
		},
		neDeviceSchemaNames["IsBYOL"]: {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: neDeviceDescriptions["IsBYOL"],
		},
		neDeviceSchemaNames["LicenseToken"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["LicenseToken"],
		},
		neDeviceSchemaNames["LicenseFile"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["LicenseFile"],
		},
		neDeviceSchemaNames["LicenseFileID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["LicenseFileID"],
		},
		neDeviceSchemaNames["ACLTemplateUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["ACLTemplateUUID"],
		},
		neDeviceSchemaNames["MgmtAclTemplateUuid"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["MgmtAclTemplateUuid"],
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
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["AccountNumber"],
		},
		neDeviceSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: neDeviceDescriptions["Notifications"],
		},
		neDeviceSchemaNames["PurchaseOrderNumber"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["PurchaseOrderNumber"],
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
			Type:        schema.TypeInt,
			Computed:    true,
			Description: neDeviceDescriptions["TermLength"],
		},
		neDeviceSchemaNames["AdditionalBandwidth"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: neDeviceDescriptions["AdditionalBandwidth"],
		},
		neDeviceSchemaNames["OrderReference"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["OrderReference"],
		},
		neDeviceSchemaNames["InterfaceCount"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: neDeviceDescriptions["InterfaceCount"],
		},
		neDeviceSchemaNames["CoreCount"]: {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: neDeviceDescriptions["CoreCount"],
		},
		neDeviceSchemaNames["IsSelfManaged"]: {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: neDeviceDescriptions["IsSelfManaged"],
		},
		neDeviceSchemaNames["WanInterfaceId"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["WanInterfaceId"],
		},
		neDeviceSchemaNames["Interfaces"]: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createDataSourceNetworkDeviceInterfaceSchema(),
			},
			Description: neDeviceDescriptions["Interfaces"],
		},
		neDeviceSchemaNames["VendorConfiguration"]: {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: neDeviceDescriptions["VendorConfiguration"],
		},
		neDeviceSchemaNames["UserPublicKey"]: {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createDataSourceNetworkDeviceUserKeySchema(),
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
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["Connectivity"],
		},
		neDeviceSchemaNames["ProjectID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["ProjectID"],
		},
		neDeviceSchemaNames["DiverseFromDeviceUUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["DiverseFromDeviceUUID"],
		},
		neDeviceSchemaNames["DiverseFromDeviceName"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: neDeviceDescriptions["DiverseFromDeviceName"],
		},
		neDeviceSchemaNames["Secondary"]: {
			Type:        schema.TypeList,
			Computed:    true,
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
						Type:        schema.TypeString,
						Required:    true,
						Description: neDeviceDescriptions["Name"],
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
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["MetroCode"],
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
						Computed:    true,
						Description: neDeviceDescriptions["HostName"],
					},
					neDeviceSchemaNames["LicenseToken"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["LicenseToken"],
					},
					neDeviceSchemaNames["LicenseFile"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["LicenseFile"],
					},
					neDeviceSchemaNames["LicenseFileID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["LicenseFileID"],
					},
					neDeviceSchemaNames["CloudInitFileID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["CloudInitFileID"],
					},
					neDeviceSchemaNames["ACLTemplateUUID"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["ACLTemplateUUID"],
					},
					neDeviceSchemaNames["MgmtAclTemplateUuid"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["MgmtAclTemplateUuid"],
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
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["AccountNumber"],
					},
					neDeviceSchemaNames["Notifications"]: {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
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
						Computed:    true,
						Description: neDeviceDescriptions["AdditionalBandwidth"],
					},
					neDeviceSchemaNames["WanInterfaceId"]: {
						Type:        schema.TypeString,
						Computed:    true,
						Description: neDeviceDescriptions["WanInterfaceId"],
					},
					neDeviceSchemaNames["Interfaces"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createDataSourceNetworkDeviceInterfaceSchema(),
						},
						Description: neDeviceDescriptions["Interfaces"],
					},
					neDeviceSchemaNames["VendorConfiguration"]: {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						Description: neDeviceDescriptions["VendorConfiguration"],
					},
					neDeviceSchemaNames["UserPublicKey"]: {
						Type:     schema.TypeSet,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createDataSourceNetworkDeviceUserKeySchema(),
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
			Computed:    true,
			Description: neDeviceDescriptions["ClusterDetails"],
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

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkDeviceRead,
		Description: "Use this data source to get details of Equinix Network Edge network device with a given Name or UUID",
		Schema:      createDataSourceNetworkDeviceSchema(),
	}
}

func getDeviceByName(deviceName string, conf *config.Config, validDeviceStateList *[]string) (*ne.Device, error) {
	var devices []ne.Device
	err := error(nil)
	devices, err = conf.Ne.GetDevices(*validDeviceStateList)
	if err != nil {
		return nil, fmt.Errorf("'devices: %v'", devices)
	}
	for _, device := range devices {
		if ne.StringValue(device.Name) == deviceName {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("device %s not found", deviceName)
}

func dataSourceNetworkDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*config.Config)
	var diags diag.Diagnostics
	var err error
	var primary, secondary *ne.Device

	// exactly one of uuid & name is guaranteed to be present by schema
	nameIf, nameExists := d.GetOk(neDeviceSchemaNames["Name"])
	name := nameIf.(string)
	uuidIf := d.Get(neDeviceSchemaNames["UUID"])
	uuid := uuidIf.(string)

	validDeviceStatusList, err := getNeDeviceStatusList(d.Get(neDeviceSchemaNames["ValidStatusList"]).(string))
	if err != nil {
		return diag.Errorf("cannot get network device status list due to '%v'", err)
	}

	if nameExists {
		primary, err = getDeviceByName(name, conf, validDeviceStatusList)
	} else {
		primary, err = conf.Ne.GetDevice(uuid)
	}

	if err != nil {
		return diag.Errorf("cannot fetch primary network device due to '%v'", err)
	}

	if slices.Contains([]string{ne.DeviceStateDeprovisioning, ne.DeviceStateDeprovisioned}, ne.StringValue(primary.Status)) {
		d.SetId("")
		return diags
	}
	if ne.StringValue(primary.RedundantUUID) != "" {

		secondary, err = conf.Ne.GetDevice(ne.StringValue(primary.RedundantUUID))
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
	if err := d.Set(neDeviceSchemaNames["UUID"], primary.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Name"], primary.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
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
	if err := d.Set(neDeviceSchemaNames["ProjectID"], primary.ProjectID); err != nil {
		return fmt.Errorf("error reading ProjectID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["DiverseFromDeviceUUID"], primary.DiverseFromDeviceUUID); err != nil {
		return fmt.Errorf("error reading DiverseFromDeviceUUID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["DiverseFromDeviceName"], primary.DiverseFromDeviceName); err != nil {
		return fmt.Errorf("error reading DiverseFromDeviceName: %s", err)
	}
	if secondary != nil {
		if v, ok := d.GetOk(neDeviceSchemaNames["Secondary"]); ok {
			secondaryFromSchema := expandNetworkDeviceSecondary(v.([]interface{}))
			secondary.LicenseFile = secondaryFromSchema.LicenseFile
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
