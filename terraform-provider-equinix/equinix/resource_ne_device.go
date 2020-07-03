package equinix

import (
	"fmt"
	"ne-go"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var neDeviceSchemaNames = map[string]string{
	"AccountNumber":       "account_number",
	"ACL":                 "acls",
	"AdditionalBandwidth": "additional_bandwidth",
	"DeviceTypeCode":      "device_type",
	"HostName":            "hostname",
	"LicenseFileID":       "license_file_id",
	"LicenseKey":          "license_key",
	"LicenseSecret":       "license_secret",
	"LicenseToken":        "license_token",
	"LicenseType":         "license_mode",
	"MetroCode":           "metro_code",
	"Name":                "name",
	"Notifications":       "notifications",
	"PackageCode":         "package_code",
	"PurchaseOrderNumber": "order_number",
	"TermLength":          "term_length",
	"Throughput":          "throughput",
	"ThroughputUnit":      "throughput_unit",
	"VendorConfig":        "vendor_configuration",
	"Secondary":           "secondary",
	//computed
	"DeviceSerialNo": "serial_number",
	"LicenseStatus":  "license_status",
	"PrimaryDNSName": "dns_name",
	"RedundancyType": "redundancy_type",
	"RedundantUUID":  "redundant_uuid",
	"SSHIPAddress":   "ssh_address",
	"SSHIPFqdn":      "ssh_fqdn",
	"Status":         "status",
	"UUID":           "uuid",
	"Version":        "version",
}

var neDeviceVendorSchemaNames = map[string]string{
	"SiteID":          "site_id",
	"SystemIPAddress": "system_ip_address",
}

func resourceNeDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceNeDeviceCreate,
		Read:   resourceNeDeviceRead,
		Update: resourceNeDeviceUpdate,
		Delete: resourceNeDeviceDelete,
		Schema: createNeDeviceResourceSchema(),
	}
}

func createNeDeviceResourceSchema() map[string]*schema.Schema {
	device := createNeDeviceSchema()
	device[neDeviceSchemaNames["Secondary"]] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: createNeDeviceSecondarySchema(),
		},
	}
	return device
}

func createNeDeviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceSchemaNames["AccountNumber"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["ACL"]: {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		neDeviceSchemaNames["AdditionalBandwidth"]: {
			Type:     schema.TypeInt,
			Optional: true,
		},
		neDeviceSchemaNames["DeviceTypeCode"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["HostName"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["LicenseFileID"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["LicenseKey"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["LicenseSecret"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["LicenseToken"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["LicenseType"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["MetroCode"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["Name"]: {
			Type:     schema.TypeString,
			Required: true,
		},
		neDeviceSchemaNames["Notifications"]: {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		neDeviceSchemaNames["PackageCode"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["PurchaseOrderNumber"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["TermLength"]: {
			Type:     schema.TypeInt,
			Required: true,
		},
		neDeviceSchemaNames["Throughput"]: {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["ThroughputUnit"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		neDeviceSchemaNames["VendorConfig"]: {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: createNeDeviceVendorConfigSchema(),
			},
		},
		neDeviceSchemaNames["Version"]: {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		//Computed
		neDeviceSchemaNames["DeviceSerialNo"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["LicenseStatus"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["PrimaryDNSName"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["RedundancyType"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["RedundantUUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["SSHIPAddress"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["SSHIPFqdn"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["Status"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
		neDeviceSchemaNames["UUID"]: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func createNeDeviceSecondarySchema() map[string]*schema.Schema {
	secondary := createNeDeviceSchema()
	delete(secondary, neDeviceSchemaNames["DeviceTypeCode"])
	delete(secondary, neDeviceSchemaNames["PackageCode"])
	delete(secondary, neDeviceSchemaNames["Throughput"])
	delete(secondary, neDeviceSchemaNames["ThroughputUnit"])
	delete(secondary, neDeviceSchemaNames["LicenseType"])
	delete(secondary, neDeviceSchemaNames["TermLength"])
	for _, v := range secondary {
		v.ForceNew = false
	}
	return secondary
}

func createNeDeviceVendorConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		neDeviceVendorSchemaNames["SiteID"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		neDeviceVendorSchemaNames["SystemIPAddress"]: {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
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
			return secondaryUpdate(secUpdateReq, a.(*schema.Set), b.(*schema.Set))
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
	primary := &ne.Device{}
	var secondary *ne.Device
	if v, ok := d.GetOk(neDeviceSchemaNames["AccountNumber"]); ok {
		primary.AccountNumber = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ACL"]); ok {
		primary.ACL = expandSetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["AdditionalBandwidth"]); ok {
		primary.AdditionalBandwidth = v.(int)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["DeviceTypeCode"]); ok {
		primary.DeviceTypeCode = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["HostName"]); ok {
		primary.HostName = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseFileID"]); ok {
		primary.LicenseFileID = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseKey"]); ok {
		primary.LicenseKey = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseSecret"]); ok {
		primary.LicenseSecret = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseToken"]); ok {
		primary.LicenseToken = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseType"]); ok {
		primary.LicenseType = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["MetroCode"]); ok {
		primary.MetroCode = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Name"]); ok {
		primary.Name = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Notifications"]); ok {
		primary.Notifications = expandSetToStringList(v.(*schema.Set))
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["PackageCode"]); ok {
		primary.PackageCode = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["PurchaseOrderNumber"]); ok {
		primary.PurchaseOrderNumber = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["TermLength"]); ok {
		primary.TermLength = v.(int)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Throughput"]); ok {
		primary.Throughput = v.(int)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["ThroughputUnit"]); ok {
		primary.ThroughputUnit = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["VendorConfig"]); ok {
		confSet := v.(*schema.Set)
		if confSet.Len() > 0 {
			primary.VendorConfig = &expandNeDeviceVendorConfig(confSet)[0]
		}
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Version"]); ok {
		primary.Version = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Secondary"]); ok {
		secSet := v.(*schema.Set)
		if secSet.Len() > 0 {
			secondary = &expandNeDeviceSecondary(secSet)[0]
		}
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["DeviceSerialNo"]); ok {
		primary.DeviceSerialNo = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["LicenseStatus"]); ok {
		primary.LicenseStatus = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["PrimaryDNSName"]); ok {
		primary.PrimaryDNSName = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["RedundancyType"]); ok {
		primary.RedundancyType = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["RedundantUUID"]); ok {
		primary.RedundantUUID = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["SSHIPAddress"]); ok {
		primary.SSHIPAddress = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["SSHIPFqdn"]); ok {
		primary.SSHIPFqdn = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["Status"]); ok {
		primary.Status = v.(string)
	}
	if v, ok := d.GetOk(neDeviceSchemaNames["UUID"]); ok {
		primary.UUID = v.(string)
	}
	return primary, secondary
}

func updateNeDeviceResource(primary *ne.Device, secondary *ne.Device, d *schema.ResourceData) error {
	if err := d.Set(neDeviceSchemaNames["AccountNumber"], primary.AccountNumber); err != nil {
		return fmt.Errorf("error reading AccountNumber: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["ACL"], primary.ACL); err != nil {
		return fmt.Errorf("error reading ACL: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["AdditionalBandwidth"], primary.AdditionalBandwidth); err != nil {
		return fmt.Errorf("error reading AdditionalBandwidth: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["DeviceTypeCode"], primary.DeviceTypeCode); err != nil {
		return fmt.Errorf("error reading DeviceTypeCode: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["HostName"], primary.HostName); err != nil {
		return fmt.Errorf("error reading HostName: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["LicenseFileID"], primary.LicenseFileID); err != nil {
		return fmt.Errorf("error reading LicenseFileID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["LicenseKey"], primary.LicenseKey); err != nil {
		return fmt.Errorf("error reading LicenseKey: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["LicenseSecret"], primary.LicenseSecret); err != nil {
		return fmt.Errorf("error reading LicenseSecret: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["LicenseToken"], primary.LicenseToken); err != nil {
		return fmt.Errorf("error reading LicenseToken: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["LicenseType"], primary.LicenseType); err != nil {
		return fmt.Errorf("error reading LicenseType: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["MetroCode"], primary.MetroCode); err != nil {
		return fmt.Errorf("error reading MetroCode: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Name"], primary.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Notifications"], primary.Notifications); err != nil {
		return fmt.Errorf("error reading Notifications: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["PackageCode"], primary.PackageCode); err != nil {
		return fmt.Errorf("error reading PackageCode: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["PurchaseOrderNumber"], primary.PurchaseOrderNumber); err != nil {
		return fmt.Errorf("error reading PurchaseOrderNumber: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["TermLength"], primary.TermLength); err != nil {
		return fmt.Errorf("error reading TermLength: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Throughput"], primary.Throughput); err != nil {
		return fmt.Errorf("error reading Throughput: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["ThroughputUnit"], primary.ThroughputUnit); err != nil {
		return fmt.Errorf("error reading ThroughputUnit: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["VendorConfig"], flattenNeDeviceVendorConfig(primary.VendorConfig)); err != nil {
		return fmt.Errorf("error reading VendorConfig: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Version"], primary.Version); err != nil {
		return fmt.Errorf("error reading Version: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["DeviceSerialNo"], primary.DeviceSerialNo); err != nil {
		return fmt.Errorf("error reading DeviceSerialNo: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["LicenseStatus"], primary.LicenseStatus); err != nil {
		return fmt.Errorf("error reading LicenseStatus: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["PrimaryDNSName"], primary.PrimaryDNSName); err != nil {
		return fmt.Errorf("error reading PrimaryDNSName: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["RedundancyType"], primary.RedundancyType); err != nil {
		return fmt.Errorf("error reading RedundancyType: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["RedundantUUID"], primary.RedundantUUID); err != nil {
		return fmt.Errorf("error reading RedundantUUID: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["SSHIPAddress"], primary.SSHIPAddress); err != nil {
		return fmt.Errorf("error reading SSHIPAddress: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["SSHIPFqdn"], primary.SSHIPFqdn); err != nil {
		return fmt.Errorf("error reading SSHIPFqdn: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["Status"], primary.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(neDeviceSchemaNames["UUID"], primary.UUID); err != nil {
		return fmt.Errorf("error reading UUID: %s", err)
	}
	if secondary != nil {
		if err := d.Set(neDeviceSchemaNames["Secondary"], flattenNeDeviceSecondary(*secondary)); err != nil {
			return fmt.Errorf("error reading Secondary: %s", err)
		}
	}
	return nil
}

func flattenNeDeviceVendorConfig(vendorConfig *ne.DeviceVendorConfig) interface{} {
	if vendorConfig != nil {
		transformed := make(map[string]interface{})
		transformed[neDeviceVendorSchemaNames["SiteID"]] = vendorConfig.SiteID
		transformed[neDeviceVendorSchemaNames["SystemIPAddress"]] = vendorConfig.SystemIPAddress
		return []map[string]interface{}{transformed}
	}
	return nil
}

func flattenNeDeviceSecondary(device ne.Device) interface{} {
	transformed := make(map[string]interface{})
	transformed[neDeviceSchemaNames["AccountNumber"]] = device.AccountNumber
	transformed[neDeviceSchemaNames["ACL"]] = device.ACL
	transformed[neDeviceSchemaNames["AdditionalBandwidth"]] = device.AdditionalBandwidth
	transformed[neDeviceSchemaNames["HostName"]] = device.HostName
	transformed[neDeviceSchemaNames["LicenseFileID"]] = device.LicenseFileID
	transformed[neDeviceSchemaNames["LicenseKey"]] = device.LicenseKey
	transformed[neDeviceSchemaNames["LicenseSecret"]] = device.LicenseSecret
	transformed[neDeviceSchemaNames["LicenseToken"]] = device.LicenseToken
	transformed[neDeviceSchemaNames["MetroCode"]] = device.MetroCode
	transformed[neDeviceSchemaNames["Name"]] = device.Name
	transformed[neDeviceSchemaNames["Notifications"]] = device.Notifications
	transformed[neDeviceSchemaNames["PurchaseOrderNumber"]] = device.PurchaseOrderNumber
	transformed[neDeviceSchemaNames["VendorConfig"]] = flattenNeDeviceVendorConfig(device.VendorConfig)
	transformed[neDeviceSchemaNames["DeviceSerialNo"]] = device.DeviceSerialNo
	transformed[neDeviceSchemaNames["LicenseStatus"]] = device.LicenseStatus
	transformed[neDeviceSchemaNames["PrimaryDNSName"]] = device.PrimaryDNSName
	transformed[neDeviceSchemaNames["RedundancyType"]] = device.RedundancyType
	transformed[neDeviceSchemaNames["RedundantUUID"]] = device.RedundantUUID
	transformed[neDeviceSchemaNames["SSHIPAddress"]] = device.SSHIPAddress
	transformed[neDeviceSchemaNames["SSHIPFqdn"]] = device.SSHIPFqdn
	transformed[neDeviceSchemaNames["Status"]] = device.Status
	transformed[neDeviceSchemaNames["UUID"]] = device.UUID
	return []map[string]interface{}{transformed}
}

func expandNeDeviceVendorConfig(features *schema.Set) []ne.DeviceVendorConfig {
	transformed := make([]ne.DeviceVendorConfig, 0, features.Len())
	for _, feature := range features.List() {
		confMap := feature.(map[string]interface{})
		transformed = append(transformed, ne.DeviceVendorConfig{
			SiteID:          confMap[neDeviceVendorSchemaNames["SiteID"]].(string),
			SystemIPAddress: confMap[neDeviceVendorSchemaNames["SystemIPAddress"]].(string),
		})
	}
	return transformed
}

func expandNeDeviceSecondary(devices *schema.Set) []ne.Device {
	transformed := make([]ne.Device, 0, devices.Len())
	for _, device := range devices.List() {
		devMap := device.(map[string]interface{})
		dev := ne.Device{}
		if v, ok := devMap[neDeviceSchemaNames["AccountNumber"]]; ok {
			dev.AccountNumber = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["ACL"]]; ok {
			dev.ACL = expandSetToStringList(v.(*schema.Set))
		}
		if v, ok := devMap[neDeviceSchemaNames["AdditionalBandwidth"]]; ok {
			dev.AdditionalBandwidth = v.(int)
		}
		if v, ok := devMap[neDeviceSchemaNames["HostName"]]; ok {
			dev.HostName = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["LicenseFileID"]]; ok {
			dev.LicenseFileID = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["LicenseKey"]]; ok {
			dev.LicenseKey = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["LicenseSecret"]]; ok {
			dev.LicenseSecret = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["LicenseToken"]]; ok {
			dev.LicenseToken = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["MetroCode"]]; ok {
			dev.MetroCode = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["Name"]]; ok {
			dev.Name = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["Notifications"]]; ok {
			dev.Notifications = expandSetToStringList(v.(*schema.Set))
		}
		if v, ok := devMap[neDeviceSchemaNames["PurchaseOrderNumber"]]; ok {
			dev.PurchaseOrderNumber = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["VendorConfig"]]; ok {
			conf := v.(*schema.Set)
			if conf.Len() > 0 {
				dev.VendorConfig = &expandNeDeviceVendorConfig(conf)[0]
			}
		}
		if v, ok := devMap[neDeviceSchemaNames["DeviceSerialNo"]]; ok {
			dev.DeviceSerialNo = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["LicenseStatus"]]; ok {
			dev.LicenseStatus = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["PrimaryDNSName"]]; ok {
			dev.PrimaryDNSName = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["RedundancyType"]]; ok {
			dev.RedundancyType = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["RedundantUUID"]]; ok {
			dev.RedundantUUID = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["SSHIPAddress"]]; ok {
			dev.SSHIPAddress = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["SSHIPFqdn"]]; ok {
			dev.SSHIPFqdn = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["Status"]]; ok {
			dev.Status = v.(string)
		}
		if v, ok := devMap[neDeviceSchemaNames["UUID"]]; ok {
			dev.UUID = v.(string)
		}
		transformed = append(transformed, dev)
	}
	return transformed
}

func secondaryUpdate(req ne.DeviceUpdateRequest, a, b *schema.Set) error {
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
