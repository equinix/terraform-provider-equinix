package equinix

import (
	"fmt"
	"reflect"
	"time"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var networkDeviceSchemaNames = map[string]string{
	"UUID":                "uuid",
	"Name":                "name",
	"TypeCode":            "type_code",
	"Status":              "status",
	"LicenseStatus":       "license_status",
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
	"ACLs":                "acls",
	"ACLStatus":           "acls_status",
	"SSHIPAddress":        "ssh_ip_address",
	"SSHIPFqdn":           "ssh_ip_fqdn",
	"AccountNumber":       "account_number",
	"Notifications":       "notifications",
	"PurchaseOrderNumber": "purchase_order_number",
	"RedundancyType":      "redundancy_type",
	"RedundantUUID":       "redundant_uuid",
	"TermLength":          "term_length",
	"AdditionalBandwidth": "additional_bandwidth",
	"OrderReference":      "order_reference",
	"InterfaceCount":      "interface_count",
	"CoreCount":           "core_count",
	"IsSelfManaged":       "self_managed",
	"Interfaces":          "interface",
	"VendorConfiguration": "vendor_configuration",
	"Secondary":           "secondary_device",
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

func resourceNetworkDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkDeviceCreate,
		Read:   resourceNetworkDeviceRead,
		Update: resourceNetworkDeviceUpdate,
		Delete: resourceNetworkDeviceDelete,
		Schema: createNetworkDeviceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func createNetworkDeviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkDeviceSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 50),
		},
		networkDeviceSchemaNames["TypeCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		networkDeviceSchemaNames["Status"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["LicenseStatus"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["MetroCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: stringIsMetroCode(),
		},
		networkDeviceSchemaNames["IBX"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["Region"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["Throughput"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		networkDeviceSchemaNames["ThroughputUnit"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"Mbps", "Gbps"}, false),
		},
		networkDeviceSchemaNames["HostName"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(2, 10),
		},
		networkDeviceSchemaNames["PackageCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		networkDeviceSchemaNames["Version"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		networkDeviceSchemaNames["IsBYOL"]: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
		networkDeviceSchemaNames["LicenseToken"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		networkDeviceSchemaNames["ACLs"]: {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},
		networkDeviceSchemaNames["ACLStatus"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["SSHIPAddress"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["SSHIPFqdn"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["AccountNumber"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		networkDeviceSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: stringIsEmailAddress(),
			},
		},
		networkDeviceSchemaNames["PurchaseOrderNumber"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 30),
		},
		networkDeviceSchemaNames["RedundancyType"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["RedundantUUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		networkDeviceSchemaNames["TermLength"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntInSlice([]int{1, 12, 24, 36}),
		},
		networkDeviceSchemaNames["AdditionalBandwidth"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		networkDeviceSchemaNames["OrderReference"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 100),
		},
		networkDeviceSchemaNames["InterfaceCount"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		networkDeviceSchemaNames["CoreCount"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		networkDeviceSchemaNames["IsSelfManaged"]: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
		networkDeviceSchemaNames["Interfaces"]: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createNetworkDeviceInterfaceSchema(),
			},
		},
		networkDeviceSchemaNames["VendorConfiguration"]: {
			Type:     schema.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		networkDeviceSchemaNames["Secondary"]: {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					networkDeviceSchemaNames["UUID"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["Name"]: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(3, 50),
					},
					networkDeviceSchemaNames["Status"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["LicenseStatus"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["MetroCode"]: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: stringIsMetroCode(),
					},
					networkDeviceSchemaNames["IBX"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["Region"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["HostName"]: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(2, 15),
					},
					networkDeviceSchemaNames["LicenseToken"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					networkDeviceSchemaNames["ACLs"]: {
						Type:     schema.TypeSet,
						Optional: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.IsCIDR,
						},
					},
					networkDeviceSchemaNames["ACLStatus"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["SSHIPAddress"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["SSHIPFqdn"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["AccountNumber"]: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					networkDeviceSchemaNames["Notifications"]: {
						Type:     schema.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: stringIsEmailAddress(),
						},
					},
					networkDeviceSchemaNames["RedundancyType"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["RedundantUUID"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					networkDeviceSchemaNames["AdditionalBandwidth"]: {
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					networkDeviceSchemaNames["Interfaces"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: createNetworkDeviceInterfaceSchema(),
						},
					},
					networkDeviceSchemaNames["VendorConfiguration"]: {
						Type:     schema.TypeMap,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func createNetworkDeviceInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceInterfaceSchemaNames["ID"]: {
			Type:     schema.TypeInt,
			Computed: true,
		},
		neDeviceInterfaceSchemaNames["Name"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceInterfaceSchemaNames["Status"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceInterfaceSchemaNames["OperationalStatus"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceInterfaceSchemaNames["MACAddress"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceInterfaceSchemaNames["IPAddress"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceInterfaceSchemaNames["AssignedType"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceInterfaceSchemaNames["Type"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceNetworkDeviceCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	primary, secondary := createNetworkDevices(d)
	var uuid string
	var err error
	if secondary != nil {
		uuid, _, err = conf.ne.CreateRedundantDevice(*primary, *secondary)
	} else {
		uuid, err = conf.ne.CreateDevice(*primary)
	}
	if err != nil {
		return err
	}
	d.SetId(uuid)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			ne.DeviceStateInitializing,
			ne.DeviceStateProvisioning,
			ne.DeviceStateWaitingSecondary,
		},
		Target: []string{
			ne.DeviceStateProvisioned,
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := conf.ne.GetDevice(d.Id())
			if err != nil {
				return nil, "", err
			}
			return resp, resp.Status, nil
		},
	}
	if _, err := createStateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for network device (%s) to be created: %s", d.Id(), err)
	}
	return resourceNetworkDeviceRead(d, m)
}

func resourceNetworkDeviceRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	var err error
	var primary, secondary *ne.Device
	var primaryACLs, secondaryACLs *ne.DeviceACLs
	primary, err = conf.ne.GetDevice(d.Id())
	if err != nil {
		return fmt.Errorf("cannot fetch primary network device due to %v", err)
	}
	if primary.Status == ne.DeviceStateDeprovisioning || primary.Status == ne.DeviceStateDeprovisioned {
		d.SetId("")
		return nil
	}
	primaryACLs, err = conf.ne.GetDeviceACLs(d.Id())
	if err != nil {
		return fmt.Errorf("cannot fetch primary network device ACLs due to %v", err)
	}
	if primary.RedundantUUID != "" {
		secondary, err = conf.ne.GetDevice(primary.RedundantUUID)
		if err != nil {
			return fmt.Errorf("cannot fetch secondary network device due to %v", err)
		}
		secondaryACLs, err = conf.ne.GetDeviceACLs(primary.RedundantUUID)
		if err != nil {
			return fmt.Errorf("cannot fetch secondary network device ACLs due to %v", err)
		}
	}
	if err = updateNetworkDeviceResource(primary, secondary, d); err != nil {
		return err
	}
	if err = updateNetworkDeviceResourceACLs(primaryACLs, secondaryACLs, d); err != nil {
		return err
	}
	return nil
}

func resourceNetworkDeviceUpdate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	updateReq := conf.ne.NewDeviceUpdateRequest(d.Id())
	if v, ok := d.GetOk(networkDeviceSchemaNames["Name"]); ok && d.HasChange(networkDeviceSchemaNames["Name"]) {
		updateReq.WithDeviceName(v.(string))
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["TermLength"]); ok && d.HasChange(networkDeviceSchemaNames["TermLength"]) {
		updateReq.WithTermLength(v.(int))
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Notifications"]); ok && d.HasChange(networkDeviceSchemaNames["Notifications"]) {
		updateReq.WithNotifications(expandSetToStringList(v.(*schema.Set)))
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["AdditionalBandwidth"]); ok && d.HasChange(networkDeviceSchemaNames["AdditionalBandwidth"]) {
		updateReq.WithAdditionalBandwidth(v.(int))
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["ACLs"]); ok && d.HasChange(networkDeviceSchemaNames["ACLs"]) {
		updateReq.WithACLs(expandSetToStringList(v.(*schema.Set)))
	}
	if err := updateReq.Execute(); err != nil {
		return err
	}
	if a, b := d.GetChange(networkDeviceSchemaNames["Secondary"]); d.HasChange(networkDeviceSchemaNames["Secondary"]) {
		if v, ok := d.GetOk(networkDeviceSchemaNames["RedundantUUID"]); ok {
			secUpdateReq := conf.ne.NewDeviceUpdateRequest(v.(string))
			if err := networkDeviceSecondaryUpdate(secUpdateReq, a.(*schema.Set), b.(*schema.Set)); err != nil {
				return err
			}
		}
	}
	return resourceNetworkDeviceRead(d, m)
}

func resourceNetworkDeviceDelete(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	if err := conf.ne.DeleteDevice(d.Id()); err != nil {
		if neRestErr, ok := err.(ne.RestError); ok {
			for _, detailedErr := range neRestErr.Errors {
				if detailedErr.ErrorCode == ne.ErrorCodeDeviceRemoved {
					return nil
				}
			}
		}
		return err
	}
	return nil
}

func createNetworkDevices(d *schema.ResourceData) (*ne.Device, *ne.Device) {
	var primary *ne.Device = &ne.Device{}
	var secondary *ne.Device
	if v, ok := d.GetOk(networkDeviceSchemaNames["UUID"]); ok {
		primary.UUID = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Name"]); ok {
		primary.Name = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["TypeCode"]); ok {
		primary.TypeCode = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Status"]); ok {
		primary.Status = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["LicenseStatus"]); ok {
		primary.LicenseStatus = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["MetroCode"]); ok {
		primary.MetroCode = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["IBX"]); ok {
		primary.IBX = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Region"]); ok {
		primary.Region = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Throughput"]); ok {
		primary.Throughput = v.(int)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["ThroughputUnit"]); ok {
		primary.ThroughputUnit = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["HostName"]); ok {
		primary.HostName = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["PackageCode"]); ok {
		primary.PackageCode = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Version"]); ok {
		primary.Version = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["IsBYOL"]); ok {
		primary.IsBYOL = v.(bool)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["LicenseToken"]); ok {
		primary.LicenseToken = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["ACLs"]); ok {
		primary.ACLs = expandSetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["SSHIPAddress"]); ok {
		primary.SSHIPAddress = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["SSHIPFqdn"]); ok {
		primary.SSHIPFqdn = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["AccountNumber"]); ok {
		primary.AccountNumber = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Notifications"]); ok {
		primary.Notifications = expandSetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["PurchaseOrderNumber"]); ok {
		primary.PurchaseOrderNumber = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["RedundancyType"]); ok {
		primary.RedundancyType = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["RedundantUUID"]); ok {
		primary.RedundantUUID = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["TermLength"]); ok {
		primary.TermLength = v.(int)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["AdditionalBandwidth"]); ok {
		primary.AdditionalBandwidth = v.(int)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["OrderReference"]); ok {
		primary.OrderReference = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["InterfaceCount"]); ok {
		primary.InterfaceCount = v.(int)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["CoreCount"]); ok {
		primary.CoreCount = v.(int)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["IsSelfManaged"]); ok {
		primary.IsSelfManaged = v.(bool)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["VendorConfiguration"]); ok {
		primary.VendorConfiguration = expandInterfaceMapToStringMap(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Secondary"]); ok {
		secondarySet := v.(*schema.Set)
		if secondarySet.Len() > 0 {
			secondaries := expandNetworkDeviceSecondary(secondarySet)
			secondary = &secondaries[0]
		}
	}
	return primary, secondary
}

func updateNetworkDeviceResource(primary *ne.Device, secondary *ne.Device, d *schema.ResourceData) error {
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
	if err := d.Set(networkDeviceSchemaNames["LicenseToken"], primary.LicenseToken); err != nil {
		return fmt.Errorf("error reading LicenseToken: %s", err)
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
	if secondary != nil {
		if err := d.Set(networkDeviceSchemaNames["Secondary"], flattenNetworkDeviceSecondary(*secondary)); err != nil {
			return fmt.Errorf("error reading Secondary: %s", err)
		}
	}
	return nil
}

func updateNetworkDeviceResourceACLs(primaryACLs, secondaryACLs *ne.DeviceACLs, d *schema.ResourceData) error {
	if err := d.Set(networkDeviceSchemaNames["ACLs"], primaryACLs.ACLs); err != nil {
		return fmt.Errorf("error reading ACLs: %s", err)
	}
	if err := d.Set(networkDeviceSchemaNames["ACLStatus"], primaryACLs.Status); err != nil {
		return fmt.Errorf("error reading ACLStatus: %s", err)
	}
	if secondaryACLs != nil {
		secondarySet := d.Get(networkDeviceSchemaNames["Secondary"]).(*schema.Set)
		if secondarySet.Len() != 1 {
			return fmt.Errorf("cannot update secondary device ACLs: secondary set size is not equal to 1")
		}
		secondary := secondarySet.List()[0].(map[string]interface{})
		secondary[networkDeviceSchemaNames["ACLs"]] = secondaryACLs.ACLs
		secondary[networkDeviceSchemaNames["ACLStatus"]] = secondaryACLs.Status
		if err := d.Set(networkDeviceSchemaNames["Secondary"], []map[string]interface{}{secondary}); err != nil {
			return fmt.Errorf("error reading Secondary: %s", err)
		}
	}
	return nil
}

func flattenNetworkDeviceSecondary(device ne.Device) interface{} {
	transformed := make(map[string]interface{})
	transformed[networkDeviceSchemaNames["UUID"]] = device.UUID
	transformed[networkDeviceSchemaNames["Name"]] = device.Name
	transformed[networkDeviceSchemaNames["Status"]] = device.Status
	transformed[networkDeviceSchemaNames["LicenseStatus"]] = device.LicenseStatus
	transformed[networkDeviceSchemaNames["MetroCode"]] = device.MetroCode
	transformed[networkDeviceSchemaNames["IBX"]] = device.IBX
	transformed[networkDeviceSchemaNames["Region"]] = device.Region
	transformed[networkDeviceSchemaNames["HostName"]] = device.HostName
	transformed[networkDeviceSchemaNames["LicenseToken"]] = device.LicenseToken
	transformed[networkDeviceSchemaNames["SSHIPAddress"]] = device.SSHIPAddress
	transformed[networkDeviceSchemaNames["SSHIPFqdn"]] = device.SSHIPFqdn
	transformed[networkDeviceSchemaNames["AccountNumber"]] = device.AccountNumber
	transformed[networkDeviceSchemaNames["Notifications"]] = device.Notifications
	transformed[networkDeviceSchemaNames["RedundancyType"]] = device.RedundancyType
	transformed[networkDeviceSchemaNames["RedundantUUID"]] = device.RedundantUUID
	transformed[networkDeviceSchemaNames["AdditionalBandwidth"]] = device.AdditionalBandwidth
	transformed[networkDeviceSchemaNames["Interfaces"]] = flattenNetworkDeviceInterfaces(device.Interfaces)
	transformed[networkDeviceSchemaNames["VendorConfiguration"]] = device.VendorConfiguration
	return []map[string]interface{}{transformed}
}

func expandNetworkDeviceSecondary(devices *schema.Set) []ne.Device {
	transformed := make([]ne.Device, 0, devices.Len())
	for _, device := range devices.List() {
		devMap := device.(map[string]interface{})
		dev := ne.Device{}
		if v, ok := devMap[networkDeviceSchemaNames["UUID"]]; ok {
			dev.UUID = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["Name"]]; ok {
			dev.Name = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["Status"]]; ok {
			dev.Status = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["LicenseStatus"]]; ok {
			dev.LicenseStatus = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["MetroCode"]]; ok {
			dev.MetroCode = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["IBX"]]; ok {
			dev.IBX = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["Region"]]; ok {
			dev.Region = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["HostName"]]; ok {
			dev.HostName = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["LicenseToken"]]; ok {
			dev.LicenseToken = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["ACLs"]]; ok {
			dev.ACLs = expandSetToStringList(v.(*schema.Set))
		}
		if v, ok := devMap[networkDeviceSchemaNames["SSHIPAddress"]]; ok {
			dev.SSHIPAddress = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["SSHIPFqdn"]]; ok {
			dev.SSHIPFqdn = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["AccountNumber"]]; ok {
			dev.AccountNumber = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["Notifications"]]; ok {
			dev.Notifications = expandSetToStringList(v.(*schema.Set))
		}
		if v, ok := devMap[networkDeviceSchemaNames["RedundancyType"]]; ok {
			dev.RedundancyType = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["RedundantUUID"]]; ok {
			dev.RedundantUUID = v.(string)
		}
		if v, ok := devMap[networkDeviceSchemaNames["AdditionalBandwidth"]]; ok {
			dev.AdditionalBandwidth = v.(int)
		}
		if v, ok := devMap[networkDeviceSchemaNames["VendorConfiguration"]]; ok {
			dev.VendorConfiguration = expandInterfaceMapToStringMap(v.(map[string]interface{}))
		}
		transformed = append(transformed, dev)
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

func expandNetworkDeviceInterfaces(interfaces *schema.Set) []ne.DeviceInterface {
	interfacesList := interfaces.List()
	transformed := make([]ne.DeviceInterface, len(interfacesList))
	for i := range interfacesList {
		interfaceMap := interfacesList[i].(map[string]interface{})
		transformed[i] = ne.DeviceInterface{
			ID:                interfaceMap[neDeviceInterfaceSchemaNames["ID"]].(int),
			Name:              interfaceMap[neDeviceInterfaceSchemaNames["Name"]].(string),
			Status:            interfaceMap[neDeviceInterfaceSchemaNames["Status"]].(string),
			OperationalStatus: interfaceMap[neDeviceInterfaceSchemaNames["OperationalStatus"]].(string),
			MACAddress:        interfaceMap[neDeviceInterfaceSchemaNames["MACAddress"]].(string),
			IPAddress:         interfaceMap[neDeviceInterfaceSchemaNames["IPAddress"]].(string),
			AssignedType:      interfaceMap[neDeviceInterfaceSchemaNames["AssignedType"]].(string),
			Type:              interfaceMap[neDeviceInterfaceSchemaNames["Type"]].(string),
		}
	}
	return transformed
}

func networkDeviceSecondaryUpdate(req ne.DeviceUpdateRequest, a, b *schema.Set) error {
	if req == nil || a.Len() < 0 || b.Len() < 0 {
		return nil
	}
	aMap := a.List()[0].(map[string]interface{})
	bMap := b.List()[0].(map[string]interface{})
	if !reflect.DeepEqual(aMap[networkDeviceSchemaNames["Name"]], bMap[networkDeviceSchemaNames["Name"]]) {
		req.WithDeviceName(bMap[networkDeviceSchemaNames["Name"]].(string))
	}
	if !reflect.DeepEqual(aMap[networkDeviceSchemaNames["Notifications"]], bMap[networkDeviceSchemaNames["Notifications"]]) {
		req.WithNotifications(expandSetToStringList(bMap[networkDeviceSchemaNames["Notifications"]].(*schema.Set)))
	}
	if !reflect.DeepEqual(aMap[networkDeviceSchemaNames["AdditionalBandwidth"]], bMap[networkDeviceSchemaNames["AdditionalBandwidth"]]) {
		req.WithAdditionalBandwidth(bMap[networkDeviceSchemaNames["AdditionalBandwidth"]].(int))
	}
	if !reflect.DeepEqual(aMap[networkDeviceSchemaNames["ACLs"]], bMap[networkDeviceSchemaNames["ACLs"]]) {
		req.WithACLs(expandSetToStringList(bMap[networkDeviceSchemaNames["ACLs"]].(*schema.Set)))
	}
	return req.Execute()
}
