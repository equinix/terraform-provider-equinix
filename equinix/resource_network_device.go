package equinix

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/equinix/ne-go"
	"github.com/equinix/rest-go"
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
	"ACLTemplateUUID":     "acl_template_id",
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
	"Interfaces":          "interface",
	"VendorConfiguration": "vendor_configuration",
	"UserPublicKey":       "ssh_key",
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

var neDeviceUserKeySchemaNames = map[string]string{
	"Username": "username",
	"KeyName":  "key_name",
}

func resourceNetworkDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkDeviceCreate,
		Read:   resourceNetworkDeviceRead,
		Update: resourceNetworkDeviceUpdate,
		Delete: resourceNetworkDeviceDelete,
		Schema: createNetworkDeviceSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
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
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		networkDeviceSchemaNames["ThroughputUnit"]: {
			Type:         schema.TypeString,
			Optional:     true,
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
		networkDeviceSchemaNames["ACLTemplateUUID"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
		networkDeviceSchemaNames["UserPublicKey"]: {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: createNetworkDeviceUserKeySchema(),
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
					networkDeviceSchemaNames["ACLTemplateUUID"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
					networkDeviceSchemaNames["UserPublicKey"]: {
						Type:     schema.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 1,
						//RequiredWith: []string{networkDeviceSchemaNames["UserPublicKey"]},
						Elem: &schema.Resource{
							Schema: createNetworkDeviceUserKeySchema(),
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

func createNetworkDeviceUserKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceUserKeySchemaNames["Username"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		neDeviceUserKeySchemaNames["KeyName"]: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
	if !primary.IsBYOL {
		licenseStateConf := *createStateConf
		licenseStateConf.Pending = []string{ne.DeviceLicenseStateApplying}
		licenseStateConf.Target = []string{ne.DeviceLicenseStateRegistered}
		licenseStateConf.Refresh = func() (interface{}, string, error) {
			resp, err := conf.ne.GetDevice(d.Id())
			if err != nil {
				return nil, "", err
			}
			return resp, resp.LicenseStatus, nil
		}
		if _, err := licenseStateConf.WaitForState(); err != nil {
			return fmt.Errorf("error waiting for network device (%s) to register license: %s", d.Id(), err)
		}
	}
	return resourceNetworkDeviceRead(d, m)
}

func resourceNetworkDeviceRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	var err error
	var primary, secondary *ne.Device
	primary, err = conf.ne.GetDevice(d.Id())
	if err != nil {
		return fmt.Errorf("cannot fetch primary network device due to %v", err)
	}
	if primary.Status == ne.DeviceStateDeprovisioning || primary.Status == ne.DeviceStateDeprovisioned {
		d.SetId("")
		return nil
	}
	if primary.RedundantUUID != "" {
		secondary, err = conf.ne.GetDevice(primary.RedundantUUID)
		if err != nil {
			return fmt.Errorf("cannot fetch secondary network device due to %v", err)
		}
	}
	if err = updateNetworkDeviceResource(primary, secondary, d); err != nil {
		return err
	}
	return nil
}

func resourceNetworkDeviceUpdate(d *schema.ResourceData, m interface{}) error {
	supportedChanges := []string{"Name", "TermLength", "Notifications", "AdditionalBandwidth", "ACLTemplateUUID"}
	conf := m.(*Config)
	updateReq := conf.ne.NewDeviceUpdateRequest(d.Id())
	primaryChanges := getNetworkDeviceChanges(supportedChanges, d)
	if err := fillNetworkDeviceUpdateRequest(updateReq, primaryChanges).Execute(); err != nil {
		return err
	}
	var secondaryChanges map[string]interface{}
	if v, ok := d.GetOk(networkDeviceSchemaNames["RedundantUUID"]); ok {
		secondaryUpdateReq := conf.ne.NewDeviceUpdateRequest(v.(string))
		secondaryChanges = getNetworkDeviceChangesSecondary(supportedChanges, d)
		if err := fillNetworkDeviceUpdateRequest(secondaryUpdateReq, secondaryChanges).Execute(); err != nil {
			return err
		}
	}
	for _, stateChangeConf := range getNetworkDeviceStateChangeConfigs(conf.ne, d, primaryChanges) {
		if _, err := stateChangeConf.WaitForState(); err != nil {
			return fmt.Errorf("error waiting for network device %q to be updated: %s", d.Id(), err)
		}
	}
	for _, stateChangeConf := range getNetworkDeviceStateChangeConfigs(conf.ne, d, secondaryChanges) {
		if _, err := stateChangeConf.WaitForState(); err != nil {
			return fmt.Errorf("error waiting for network device %q to be updated: %s", d.Get(networkDeviceSchemaNames["RedundantUUID"]), err)
		}
	}
	return resourceNetworkDeviceRead(d, m)
}

func resourceNetworkDeviceDelete(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	if v, ok := d.GetOk(networkDeviceSchemaNames["ACLTemplateUUID"]); ok {
		if err := conf.ne.NewDeviceUpdateRequest(d.Id()).WithACLTemplate("").Execute(); err != nil {
			log.Printf("[WARN] could not unassign ACL template %q from device %q: %s", v, d.Id(), err)
		}
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Secondary"]); ok {
		secondarySet := v.(*schema.Set)
		if secondarySet.Len() > 0 {
			secondaries := expandNetworkDeviceSecondary(secondarySet)
			secondary := &secondaries[0]
			if secondary.ACLTemplateUUID != "" {
				if err := conf.ne.NewDeviceUpdateRequest(secondary.UUID).WithACLTemplate("").Execute(); err != nil {
					log.Printf("[WARN] could not unassign ACL template %q from device %q: %s", v, secondary.UUID, err)
				}
			}
		}
	}
	if err := conf.ne.DeleteDevice(d.Id()); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			for _, detailedErr := range restErr.ApplicationErrors {
				if detailedErr.Code == ne.ErrorCodeDeviceRemoved {
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
	if v, ok := d.GetOk(networkDeviceSchemaNames["ACLTemplateUUID"]); ok {
		primary.ACLTemplateUUID = v.(string)
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
	if v, ok := d.GetOk(networkDeviceSchemaNames["UserPublicKey"]); ok {
		userKeys := expandNetworkDeviceUserKeys(v.(*schema.Set))
		if len(userKeys) > 0 {
			primary.UserPublicKey = userKeys[0]
		}
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
	if secondary != nil {
		if err := d.Set(networkDeviceSchemaNames["Secondary"], flattenNetworkDeviceSecondary(*secondary)); err != nil {
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
	transformed[networkDeviceSchemaNames["ACLTemplateUUID"]] = device.ACLTemplateUUID
	transformed[networkDeviceSchemaNames["SSHIPAddress"]] = device.SSHIPAddress
	transformed[networkDeviceSchemaNames["SSHIPFqdn"]] = device.SSHIPFqdn
	transformed[networkDeviceSchemaNames["AccountNumber"]] = device.AccountNumber
	transformed[networkDeviceSchemaNames["Notifications"]] = device.Notifications
	transformed[networkDeviceSchemaNames["RedundancyType"]] = device.RedundancyType
	transformed[networkDeviceSchemaNames["RedundantUUID"]] = device.RedundantUUID
	transformed[networkDeviceSchemaNames["AdditionalBandwidth"]] = device.AdditionalBandwidth
	transformed[networkDeviceSchemaNames["Interfaces"]] = flattenNetworkDeviceInterfaces(device.Interfaces)
	transformed[networkDeviceSchemaNames["VendorConfiguration"]] = device.VendorConfiguration
	transformed[networkDeviceSchemaNames["UserPublicKey"]] = flattenNetworkDeviceUserKeys([]*ne.DeviceUserPublicKey{device.UserPublicKey})
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
		if v, ok := devMap[networkDeviceSchemaNames["ACLTemplateUUID"]]; ok {
			dev.ACLTemplateUUID = v.(string)
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
		if v, ok := devMap[networkDeviceSchemaNames["UserPublicKey"]]; ok {
			userKeys := expandNetworkDeviceUserKeys(v.(*schema.Set))
			if len(userKeys) > 0 {
				dev.UserPublicKey = userKeys[0]
			}
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
			Username: userKeyMap[neDeviceUserKeySchemaNames["Username"]].(string),
			KeyName:  userKeyMap[neDeviceUserKeySchemaNames["KeyName"]].(string),
		}
	}
	return transformed
}

func getNetworkDeviceChanges(keys []string, d *schema.ResourceData) map[string]interface{} {
	changed := make(map[string]interface{})
	for _, key := range keys {
		if schemaKey, ok := networkDeviceSchemaNames[key]; ok {
			if v := d.Get(schemaKey); v != nil && d.HasChange(schemaKey) {
				changed[key] = v
			}
		}
	}
	return changed
}

func getNetworkDeviceChangesSecondary(keys []string, d *schema.ResourceData) map[string]interface{} {
	changed := make(map[string]interface{})
	if !d.HasChange(networkDeviceSchemaNames["Secondary"]) {
		return changed
	}
	a, b := d.GetChange(networkDeviceSchemaNames["Secondary"])
	aSet := a.(*schema.Set)
	bSet := b.(*schema.Set)
	if aSet.Len() < 1 || bSet.Len() < 1 {
		return changed
	}
	aMap := aSet.List()[0].(map[string]interface{})
	bMap := bSet.List()[0].(map[string]interface{})
	for _, key := range keys {
		if schemaKey, ok := networkDeviceSchemaNames[key]; ok {
			if !reflect.DeepEqual(aMap[schemaKey], bMap[schemaKey]) {
				changed[key] = bMap[schemaKey]
			}
		}
	}
	return changed
}

func fillNetworkDeviceUpdateRequest(updateReq ne.DeviceUpdateRequest, changes map[string]interface{}) ne.DeviceUpdateRequest {
	for change, changeValue := range changes {
		switch change {
		case "Name":
			updateReq.WithDeviceName(changeValue.(string))
		case "TermLength":
			updateReq.WithTermLength(changeValue.(int))
		case "Notifications":
			updateReq.WithNotifications((expandSetToStringList(changeValue.(*schema.Set))))
		case "AdditionalBandwidth":
			updateReq.WithAdditionalBandwidth(changeValue.(int))
		case "ACLTemplateUUID":
			updateReq.WithACLTemplate(changeValue.(string))
		}
	}
	return updateReq
}

func getNetworkDeviceStateChangeConfigs(c ne.Client, d *schema.ResourceData, changes map[string]interface{}) []*resource.StateChangeConf {
	configs := make([]*resource.StateChangeConf, 0, len(changes))
	for change, changeValue := range changes {
		switch change {
		case "ACLTemplateUUID":
			aclTempID, ok := changeValue.(string)
			if !ok || aclTempID == "" {
				break
			}
			configs = append(configs, &resource.StateChangeConf{
				Pending: []string{
					ne.ACLDeviceStatusProvisioning,
				},
				Target: []string{
					ne.ACLDeviceStatusProvisioned,
				},
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      1 * time.Second,
				MinTimeout: 1 * time.Second,
				Refresh: func() (interface{}, string, error) {
					resp, err := c.GetACLTemplate(aclTempID)
					if err != nil {
						return nil, "", err
					}
					return resp, resp.DeviceACLStatus, nil
				},
			})
		}
	}
	return configs
}
