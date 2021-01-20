package equinix

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
			Update: schema.DefaultTimeout(10 * time.Minute),
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
			Optional:     true,
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
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{networkDeviceSchemaNames["LicenseFile"]},
		},
		networkDeviceSchemaNames["LicenseFile"]: {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		networkDeviceSchemaNames["LicenseFileID"]: {
			Type:     schema.TypeString,
			Computed: true,
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
			ForceNew:     true,
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
			Type:     schema.TypeList,
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
					networkDeviceSchemaNames["HostName"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(2, 15),
					},
					networkDeviceSchemaNames["LicenseToken"]: {
						Type:          schema.TypeString,
						Optional:      true,
						ForceNew:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{networkDeviceSchemaNames["Secondary"] + ".0." + networkDeviceSchemaNames["LicenseFile"]},
					},
					networkDeviceSchemaNames["LicenseFile"]: {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					networkDeviceSchemaNames["LicenseFileID"]: {
						Type:     schema.TypeString,
						Computed: true,
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
	var primaryID, secondaryID string
	var err error
	if err := uploadDeviceLicenseFile(conf.ne, primary.TypeCode, primary); err != nil {
		return fmt.Errorf("could not upload primary device license file due to %s", err)
	}
	if err := uploadDeviceLicenseFile(conf.ne, primary.TypeCode, secondary); err != nil {
		return fmt.Errorf("could not upload secondary device license file due to %s", err)
	}
	if secondary != nil {
		primaryID, secondaryID, err = conf.ne.CreateRedundantDevice(*primary, *secondary)
	} else {
		primaryID, err = conf.ne.CreateDevice(*primary)
	}
	if err != nil {
		return err
	}
	d.SetId(primaryID)
	provWaitConfigs := []*resource.StateChangeConf{createNetworkDeviceProvisioningWaitConfiguration(conf.ne, d.Timeout(schema.TimeoutCreate), primaryID)}
	licWaitConfigs := []*resource.StateChangeConf{createNetworkDeviceLicenseWaitConfiguration(conf.ne, d.Timeout(schema.TimeoutCreate), primaryID)}
	if secondary != nil {
		provWaitConfigs = append(provWaitConfigs, createNetworkDeviceProvisioningWaitConfiguration(conf.ne, d.Timeout(schema.TimeoutCreate), secondaryID))
		licWaitConfigs = append(licWaitConfigs, createNetworkDeviceLicenseWaitConfiguration(conf.ne, d.Timeout(schema.TimeoutCreate), secondaryID))
	}
	for i := range provWaitConfigs {
		if _, err := provWaitConfigs[i].WaitForState(); err != nil {
			return fmt.Errorf("error waiting for network device (%s) to be created: %s", primaryID, err)
		}
	}
	for i := range licWaitConfigs {
		if _, err := licWaitConfigs[i].WaitForState(); err != nil {
			return fmt.Errorf("error waiting for network device (%s) license to be applied: %s", primaryID, err)
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
	conf := m.(*Config)
	supportedChanges := []string{networkDeviceSchemaNames["Name"], networkDeviceSchemaNames["TermLength"],
		networkDeviceSchemaNames["Notifications"], networkDeviceSchemaNames["AdditionalBandwidth"],
		networkDeviceSchemaNames["ACLTemplateUUID"]}
	updateReq := conf.ne.NewDeviceUpdateRequest(d.Id())
	primaryChanges := getResourceDataChangedKeys(supportedChanges, d)
	if err := fillNetworkDeviceUpdateRequest(updateReq, primaryChanges).Execute(); err != nil {
		return err
	}
	var secondaryChanges map[string]interface{}
	if v, ok := d.GetOk(networkDeviceSchemaNames["RedundantUUID"]); ok {
		secondaryChanges = getResourceDataListElementChanges(supportedChanges, networkDeviceSchemaNames["Secondary"], 0, d)
		secondaryUpdateReq := conf.ne.NewDeviceUpdateRequest(v.(string))
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
		if secondary := expandNetworkDeviceSecondary(v.([]interface{})); secondary != nil {
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
	var primary, secondary *ne.Device
	primary = &ne.Device{}
	if v, ok := d.GetOk(networkDeviceSchemaNames["Name"]); ok {
		primary.Name = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["TypeCode"]); ok {
		primary.TypeCode = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["MetroCode"]); ok {
		primary.MetroCode = v.(string)
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
	if v, ok := d.GetOk(networkDeviceSchemaNames["LicenseFile"]); ok {
		primary.LicenseFile = v.(string)
	}
	if v, ok := d.GetOk(networkDeviceSchemaNames["ACLTemplateUUID"]); ok {
		primary.ACLTemplateUUID = v.(string)
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
		secondary = expandNetworkDeviceSecondary(v.([]interface{}))
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
	if secondary != nil {
		if v, ok := d.GetOk(networkDeviceSchemaNames["Secondary"]); ok {
			secondaryFromSchema := expandNetworkDeviceSecondary(v.([]interface{}))
			secondary.LicenseFile = secondaryFromSchema.LicenseFile
		}
		if err := d.Set(networkDeviceSchemaNames["Secondary"], flattenNetworkDeviceSecondary(secondary)); err != nil {
			return fmt.Errorf("error reading Secondary: %s", err)
		}
	}
	return nil
}

func flattenNetworkDeviceSecondary(device *ne.Device) interface{} {
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
	transformed[networkDeviceSchemaNames["LicenseFileID"]] = device.LicenseFileID
	transformed[networkDeviceSchemaNames["LicenseFile"]] = device.LicenseFile
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
	return []interface{}{transformed}
}

func expandNetworkDeviceSecondary(devices []interface{}) *ne.Device {
	if len(devices) < 1 {
		log.Printf("[WARN] resource_network_device expanding empty secondary device collection")
		return nil
	}
	device := devices[0].(map[string]interface{})
	transformed := &ne.Device{}
	if v, ok := device[networkDeviceSchemaNames["Name"]]; ok {
		transformed.Name = v.(string)
	}
	if v, ok := device[networkDeviceSchemaNames["MetroCode"]]; ok {
		transformed.MetroCode = v.(string)
	}
	if v, ok := device[networkDeviceSchemaNames["HostName"]]; ok {
		transformed.HostName = v.(string)
	}
	if v, ok := device[networkDeviceSchemaNames["LicenseToken"]]; ok {
		transformed.LicenseToken = v.(string)
	}
	if v, ok := device[networkDeviceSchemaNames["LicenseFile"]]; ok {
		transformed.LicenseFile = v.(string)
	}
	if v, ok := device[networkDeviceSchemaNames["ACLTemplateUUID"]]; ok {
		transformed.ACLTemplateUUID = v.(string)
	}
	if v, ok := device[networkDeviceSchemaNames["AccountNumber"]]; ok {
		transformed.AccountNumber = v.(string)
	}
	if v, ok := device[networkDeviceSchemaNames["Notifications"]]; ok {
		transformed.Notifications = expandSetToStringList(v.(*schema.Set))
	}
	if v, ok := device[networkDeviceSchemaNames["AdditionalBandwidth"]]; ok {
		transformed.AdditionalBandwidth = v.(int)
	}
	if v, ok := device[networkDeviceSchemaNames["VendorConfiguration"]]; ok {
		transformed.VendorConfiguration = expandInterfaceMapToStringMap(v.(map[string]interface{}))
	}
	if v, ok := device[networkDeviceSchemaNames["UserPublicKey"]]; ok {
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
			Username: userKeyMap[neDeviceUserKeySchemaNames["Username"]].(string),
			KeyName:  userKeyMap[neDeviceUserKeySchemaNames["KeyName"]].(string),
		}
	}
	return transformed
}

func fillNetworkDeviceUpdateRequest(updateReq ne.DeviceUpdateRequest, changes map[string]interface{}) ne.DeviceUpdateRequest {
	for change, changeValue := range changes {
		switch change {
		case networkDeviceSchemaNames["Name"]:
			updateReq.WithDeviceName(changeValue.(string))
		case networkDeviceSchemaNames["TermLength"]:
			updateReq.WithTermLength(changeValue.(int))
		case networkDeviceSchemaNames["Notifications"]:
			updateReq.WithNotifications((expandSetToStringList(changeValue.(*schema.Set))))
		case networkDeviceSchemaNames["AdditionalBandwidth"]:
			updateReq.WithAdditionalBandwidth(changeValue.(int))
		case networkDeviceSchemaNames["ACLTemplateUUID"]:
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

func uploadDeviceLicenseFile(c ne.Client, typeCode string, device *ne.Device) error {
	if device == nil || device.LicenseFile == "" {
		return nil
	}
	fileName := filepath.Base(device.LicenseFile)
	file, err := os.Open(device.LicenseFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("[WARN] could not close file %q due to an error: %s", device.LicenseFile, err)
		}
	}()
	fileID, err := c.UploadLicenseFile(device.MetroCode, typeCode, ne.DeviceManagementTypeSelf, ne.DeviceLicenseModeBYOL, fileName, file)
	if err != nil {
		return err
	}
	device.LicenseFileID = fileID
	return nil
}

func createNetworkDeviceProvisioningWaitConfiguration(c ne.Client, timeout time.Duration, uuid string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{
			ne.DeviceStateInitializing,
			ne.DeviceStateProvisioning,
			ne.DeviceStateWaitingSecondary,
		},
		Target: []string{
			ne.DeviceStateProvisioned,
		},
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := c.GetDevice(uuid)
			if err != nil {
				return nil, "", err
			}
			return resp, resp.Status, nil
		},
	}
}

func createNetworkDeviceLicenseWaitConfiguration(c ne.Client, timeout time.Duration, uuid string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{
			ne.DeviceLicenseStateApplying,
			"",
		},
		Target: []string{
			ne.DeviceLicenseStateRegistered,
			ne.DeviceLicenseStateApplied,
		},
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err := c.GetDevice(uuid)
			if err != nil {
				return nil, "", err
			}
			return resp, resp.LicenseStatus, nil
		},
	}
}
