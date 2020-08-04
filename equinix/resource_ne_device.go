package equinix

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var neDeviceSchemaNames = map[string]string{
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

func resourceNeDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceNeDeviceCreate,
		Read:   resourceNeDeviceRead,
		Update: resourceNeDeviceUpdate,
		Delete: resourceNeDeviceDelete,
		Schema: createNeDeviceSchema(),
	}
}

func createNeDeviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["Name"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},
		neDeviceSchemaNames["TypeCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		neDeviceSchemaNames["Status"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["LicenseStatus"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["MetroCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile("^[A-Z]{2}$"), "MetroCode must consist of two capital letters"),
		},
		neDeviceSchemaNames["IBX"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["Region"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["Throughput"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		neDeviceSchemaNames["ThroughputUnit"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"Mbps", "Gbps"}, false),
		},
		neDeviceSchemaNames["HostName"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 15),
		},
		neDeviceSchemaNames["PackageCode"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["Version"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["IsBYOL"]: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
		neDeviceSchemaNames["LicenseToken"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["ACLs"]: {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},
		neDeviceSchemaNames["SSHIPAddress"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["SSHIPFqdn"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["AccountNumber"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^ @]+@[^ @]+$"), "Notification list can contain only valid email addresses"),
			},
		},
		neDeviceSchemaNames["PurchaseOrderNumber"]: {
			Type:     schema.TypeString,
			Optional: true,
		},
		neDeviceSchemaNames["RedundancyType"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["RedundantUUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["TermLength"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		neDeviceSchemaNames["AdditionalBandwidth"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		neDeviceSchemaNames["OrderReference"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["InterfaceCount"]: {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		neDeviceSchemaNames["CoreCount"]: {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		neDeviceSchemaNames["IsSelfManaged"]: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
		neDeviceSchemaNames["Interfaces"]: {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     createNeDeviceInterfaceSchema(),
		},
		neDeviceSchemaNames["Secondary"]: {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					neDeviceSchemaNames["UUID"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["Name"]: {
						Type:     schema.TypeString,
						Required: true,
						//ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 50),
					},
					neDeviceSchemaNames["Status"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["LicenseStatus"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["MetroCode"]: {
						Type:     schema.TypeString,
						Required: true,
						//ForceNew:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile("^[A-Z]{2}$"), "MetroCode must consist of two capital letters"),
					},
					neDeviceSchemaNames["IBX"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["Region"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["HostName"]: {
						Type:     schema.TypeString,
						Required: true,
						//ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 15),
					},
					neDeviceSchemaNames["LicenseToken"]: {
						Type:     schema.TypeString,
						Optional: true,
						//ForceNew: true,
					},
					neDeviceSchemaNames["ACLs"]: {
						Type:     schema.TypeSet,
						Optional: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.IsCIDR,
						},
					},
					neDeviceSchemaNames["SSHIPAddress"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["SSHIPFqdn"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["AccountNumber"]: {
						Type:     schema.TypeString,
						Optional: true,
						//ForceNew: true,
					},
					neDeviceSchemaNames["Notifications"]: {
						Type:     schema.TypeSet,
						Optional: true,
						//ForceNew: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^ @]+@[^ @]+$"), "Notification list can contain only valid email addresses"),
						},
					},
					neDeviceSchemaNames["RedundancyType"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["RedundantUUID"]: {
						Type:     schema.TypeString,
						Computed: true,
					},
					neDeviceSchemaNames["AdditionalBandwidth"]: {
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
					neDeviceSchemaNames["Interfaces"]: {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     createNeDeviceInterfaceSchema(),
					},
				},
			},
		},
	}
}

func createNeDeviceInterfaceSchema() map[string]*schema.Schema {
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

func resourceNeDeviceCreate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	primary, secondary := createNeDevices(d)
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
	return resourceNeDeviceRead(d, m)
}

func resourceNeDeviceRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	var err error
	var primary *ne.Device
	var secondary *ne.Device
	primary, err = conf.ne.GetDevice(d.Id())
	if err != nil {
		return fmt.Errorf("cannot fetch primary device due to %v", err)
	}
	if primary.Status == "DEPROVISIONED" {
		d.SetId("")
		return nil
	}
	if primary.RedundantUUID != "" {
		secondary, err = conf.ne.GetDevice(primary.RedundantUUID)
		if err != nil {
			return fmt.Errorf("cannot fetch secondary device due to %v", err)
		}
	}
	if err = updateNeDeviceResource(primary, secondary, d); err != nil {
		return err
	}
	return nil
}

func resourceNeDeviceUpdate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	updateReq := conf.ne.NewDeviceUpdateRequest(d.Id())
	if v, ok := d.GetOk(neDeviceSchemaNames["Name"]); ok && d.HasChange(neDeviceSchemaNames["Name"]) {
		updateReq.WithDeviceName(v.(string))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["TermLength"]); ok && d.HasChange(neDeviceSchemaNames["TermLength"]) {
		updateReq.WithTermLength(v.(int))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Notifications"]); ok && d.HasChange(neDeviceSchemaNames["Notifications"]) {
		updateReq.WithNotifications(expandSetToStringList(v.(*schema.Set)))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["AdditionalBandwidth"]); ok && d.HasChange(neDeviceSchemaNames["AdditionalBandwidth"]) {
		updateReq.WithAdditionalBandwidth(v.(int))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ACL"]); ok && d.HasChange(neDeviceSchemaNames["ACL"]) {
		updateReq.WithACLs(expandSetToStringList(v.(*schema.Set)))
	}
	if err := updateReq.Execute(); err != nil {
		return err
	}
	if a, b := d.GetChange(neDeviceSchemaNames["Secondary"]); d.HasChange(neDeviceSchemaNames["Secondary"]) {
		if v, ok := d.GetOk(neDeviceSchemaNames["RedundantUUID"]); ok {
			secUpdateReq := conf.ne.NewDeviceUpdateRequest(v.(string))
			return neDeviceSecondaryUpdate(secUpdateReq, a.(*schema.Set), b.(*schema.Set))
		}
	}
	return nil
}

func resourceNeDeviceDelete(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	if err := conf.ne.DeleteDevice(d.Id()); err != nil {
		/*ecxRestErr, ok := err.(ecx.RestError)
		if ok {
			//IC-LAYER2-4021 = Connection already deleted
			if hasECXErrorCode(ecxRestErr.Errors, "IC-LAYER2-4021") {
				return nil
			}
		}*/
		return err
	}
	return nil
}

func createNeDevices(d *schema.ResourceData) (*ne.Device, *ne.Device) {
	var primary *ne.Device = &ne.Device{}
	var secondary *ne.Device
	if v, ok := d.GetOk(neDeviceSchemaNames["UUID"]); ok {
		primary.UUID = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Name"]); ok {
		primary.Name = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["TypeCode"]); ok {
		primary.TypeCode = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Status"]); ok {
		primary.Status = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseStatus"]); ok {
		primary.LicenseStatus = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["MetroCode"]); ok {
		primary.MetroCode = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["IBX"]); ok {
		primary.IBX = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Region"]); ok {
		primary.Region = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Throughput"]); ok {
		primary.Throughput = v.(int)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ThroughputUnit"]); ok {
		primary.ThroughputUnit = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["HostName"]); ok {
		primary.HostName = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["PackageCode"]); ok {
		primary.PackageCode = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Version"]); ok {
		primary.Version = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["IsBYOL"]); ok {
		primary.IsBYOL = v.(bool)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseToken"]); ok {
		primary.LicenseToken = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ACLs"]); ok {
		primary.ACLs = expandSetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["SSHIPAddress"]); ok {
		primary.SSHIPAddress = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["SSHIPFqdn"]); ok {
		primary.SSHIPFqdn = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["AccountNumber"]); ok {
		primary.AccountNumber = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Notifications"]); ok {
		primary.Notifications = expandSetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["PurchaseOrderNumber"]); ok {
		primary.PurchaseOrderNumber = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["RedundancyType"]); ok {
		primary.RedundancyType = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["RedundantUUID"]); ok {
		primary.RedundantUUID = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["TermLength"]); ok {
		primary.TermLength = v.(int)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["AdditionalBandwidth"]); ok {
		primary.AdditionalBandwidth = v.(int)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["OrderReference"]); ok {
		primary.OrderReference = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["InterfaceCount"]); ok {
		primary.InterfaceCount = v.(int)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["CoreCount"]); ok {
		primary.CoreCount = v.(int)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["IsSelfManaged"]); ok {
		primary.IsSelfManaged = v.(bool)
	}
	//if v, ok := d.GetOk(neDeviceSchemaNames["Interfaces"]); ok {
	//	primary.Interfaces = expandNeDeviceInterfaces.(v.(*schema.Set)))
	//}
	if v, ok := d.GetOk(neDeviceSchemaNames["Secondary"]); ok {
		secondarySet := v.(*schema.Set)
		if secondarySet.Len() > 0 {
			secondaries := expandNeDeviceSecondary(secondarySet)
			secondary = &secondaries[0]
		}
	}
	return primary, secondary
}

func updateNeDeviceResource(primary *ne.Device, secondary *ne.Device, d *schema.ResourceData) error {
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
	if err := d.Set(neDeviceSchemaNames["LicenseToken"], primary.LicenseToken); err != nil {
		return fmt.Errorf("error reading LicenseToken: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["ACLs"], primary.ACLs); err != nil {
		return fmt.Errorf("error reading ACLs: %s", err)
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
	if err := d.Set(neDeviceSchemaNames["Interfaces"], flattenNeDeviceInterfaces(primary.Interfaces)); err != nil {
		return fmt.Errorf("error reading Interfaces: %s", err)
	}
	if secondary != nil {
		if err := d.Set(neDeviceSchemaNames["Secondary"], flattenNeDeviceSecondary(*secondary)); err != nil {
			return fmt.Errorf("error reading Secondary: %s", err)
		}
	}
	return nil
}

func flattenNeDeviceSecondary(device ne.Device) interface{} {
	transformed := make(map[string]interface{})
	transformed[neDeviceSchemaNames["UUID"]] = device.UUID
	transformed[neDeviceSchemaNames["Name"]] = device.Name
	transformed[neDeviceSchemaNames["Status"]] = device.Status
	transformed[neDeviceSchemaNames["LicenseStatus"]] = device.LicenseStatus
	transformed[neDeviceSchemaNames["MetroCode"]] = device.MetroCode
	transformed[neDeviceSchemaNames["IBX"]] = device.IBX
	transformed[neDeviceSchemaNames["Region"]] = device.Region
	transformed[neDeviceSchemaNames["HostName"]] = device.HostName
	transformed[neDeviceSchemaNames["LicenseToken"]] = device.LicenseToken
	transformed[neDeviceSchemaNames["ACLs"]] = device.ACLs
	transformed[neDeviceSchemaNames["SSHIPAddress"]] = device.SSHIPAddress
	transformed[neDeviceSchemaNames["SSHIPFqdn"]] = device.SSHIPFqdn
	transformed[neDeviceSchemaNames["AccountNumber"]] = device.AccountNumber
	transformed[neDeviceSchemaNames["Notifications"]] = device.Notifications
	transformed[neDeviceSchemaNames["RedundancyType"]] = device.RedundancyType
	transformed[neDeviceSchemaNames["RedundantUUID"]] = device.RedundantUUID
	transformed[neDeviceSchemaNames["AdditionalBandwidth"]] = device.AdditionalBandwidth
	return []map[string]interface{}{transformed}
}

func expandNeDeviceSecondary(devices *schema.Set) []ne.Device {
	transformed := make([]ne.Device, 0, devices.Len())
	for _, device := range devices.List() {
		devMap := device.(map[string]interface{})
		dev := ne.Device{}
		if v, ok := devMap[neDeviceSchemaNames["UUID"]]; ok {
			dev.UUID = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["Name"]]; ok {
			dev.Name = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["Status"]]; ok {
			dev.Status = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["LicenseStatus"]]; ok {
			dev.LicenseStatus = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["MetroCode"]]; ok {
			dev.MetroCode = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["IBX"]]; ok {
			dev.IBX = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["Region"]]; ok {
			dev.Region = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["HostName"]]; ok {
			dev.HostName = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["LicenseToken"]]; ok {
			dev.LicenseToken = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["ACLs"]]; ok {
			dev.ACLs = expandSetToStringList(v.(*schema.Set))
		}
		if v, ok := devMap[neDeviceSchemaNames["SSHIPAddress"]]; ok {
			dev.SSHIPAddress = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["SSHIPFqdn"]]; ok {
			dev.SSHIPFqdn = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["AccountNumber"]]; ok {
			dev.AccountNumber = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["Notifications"]]; ok {
			dev.Notifications = expandSetToStringList(v.(*schema.Set))
		}
		if v, ok := devMap[neDeviceSchemaNames["RedundancyType"]]; ok {
			dev.RedundancyType = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["RedundantUUID"]]; ok {
			dev.RedundantUUID = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["AdditionalBandwidth"]]; ok {
			dev.AdditionalBandwidth = v.(int)
		}
		transformed = append(transformed, dev)
	}
	return transformed
}

func flattenNeDeviceInterfaces(interfaces []ne.DeviceInterface) interface{} {
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

func expandDeviceInterfaces(interfaces *schema.Set) []ne.DeviceInterface {
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

func neDeviceSecondaryUpdate(req ne.DeviceUpdateRequest, a, b *schema.Set) error {
	if req == nil || a.Len() < 0 || b.Len() < 0 {
		return nil
	}
	aMap := a.List()[0].(map[string]interface{})
	bMap := b.List()[0].(map[string]interface{})
	if !reflect.DeepEqual(aMap[neDeviceSchemaNames["Name"]], bMap[neDeviceSchemaNames["Name"]]) {
		req.WithDeviceName(bMap[neDeviceSchemaNames["Name"]].(string))
	}
	if !reflect.DeepEqual(aMap[neDeviceSchemaNames["Notifications"]], bMap[neDeviceSchemaNames["Notifications"]]) {
		req.WithNotifications(expandSetToStringList(bMap[neDeviceSchemaNames["Notifications"]].(*schema.Set)))
	}
	if !reflect.DeepEqual(aMap[neDeviceSchemaNames["AdditionalBandwidth"]], bMap[neDeviceSchemaNames["AdditionalBandwidth"]]) {
		req.WithAdditionalBandwidth(bMap[neDeviceSchemaNames["AdditionalBandwidth"]].(int))
	}
	if !reflect.DeepEqual(aMap[neDeviceSchemaNames["ACL"]], bMap[neDeviceSchemaNames["ACL"]]) {
		req.WithACLs(expandSetToStringList(bMap[neDeviceSchemaNames["ACL"]].(*schema.Set)))
	}
	return req.Execute()
}
