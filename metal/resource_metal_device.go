package metal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path"
	"reflect"
	"regexp"
	"sort"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

var matchIPXEScript = regexp.MustCompile(`(?i)^#![i]?pxe`)
var ipAddressTypes = []string{"public_ipv4", "private_ipv4", "public_ipv6"}

var deviceCommonIncludes = []string{"project", "metro", "facility", "hardware_reservation"}
var deviceReadOptions = &packngo.GetOptions{Includes: deviceCommonIncludes}

func resourceMetalDevice() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Create: resourceMetalDeviceCreate,
		Read:   resourceMetalDeviceRead,
		Update: resourceMetalDeviceUpdate,
		Delete: resourceMetalDeviceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "The ID of the project in which to create the device",
				Required:    true,
				ForceNew:    true,
			},

			"hostname": {
				Type:        schema.TypeString,
				Description: "The device hostname used in deployments taking advantage of Layer3 DHCP or metadata service configuration.",
				Optional:    true,
				Computed:    true,
			},

			"description": {
				Type:        schema.TypeString,
				Description: "Description string for the device",
				Optional:    true,
			},

			"operating_system": {
				Type:        schema.TypeString,
				Description: "The operating system slug. To find the slug, or visit [Operating Systems API docs](https://metal.equinix.com/developers/api/operatingsystems), set your API auth token in the top of the page and see JSON from the API response",
				Required:    true,
				ForceNew:    true,
			},

			"deployed_facility": {
				Type:        schema.TypeString,
				Description: "The facility where the device is deployed",
				Computed:    true,
			},

			"metro": {
				Type:          schema.TypeString,
				Description:   "Metro area for the new device. Conflicts with facilities",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"facilities"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if len(old) > 0 && new == "" {
						// here it would be good to also test if the "old" metro
						// contains the device facility. If yes, we'd suppress diff
						// and if it's a different metro, we would re-create.
						// Not sure if this is possible.
						return true
					}
					return old == new
				},
				StateFunc: toLower,
			},

			"facilities": {
				Type:        schema.TypeList,
				Description: "List of facility codes with deployment preferences. Equinix Metal API will go through the list and will deploy your device to first facility with free capacity. List items must be facility codes or any (a wildcard). To find the facility code, visit [Facilities API docs](https://metal.equinix.com/developers/api/facilities/), set your API auth token in the top of the page and see JSON from the API response. Conflicts with metro",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
				MinItems:    1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					fsRaw := d.Get("facilities")
					fs := convertStringArr(fsRaw.([]interface{}))
					df := d.Get("deployed_facility").(string)
					if contains(fs, df) {
						return true
					}
					if contains(fs, "any") && (len(df) != 0) {
						return true
					}
					return false
				},
				ConflictsWith: []string{"metro"},
			},
			"ip_address": {
				Type:        schema.TypeList,
				Description: "A list of IP address types for the device (structure is documented below)",
				Optional:    true,
				Elem:        ipAddressSchema(),
				MinItems:    1,
			},

			"plan": {
				Type:        schema.TypeString,
				Description: "The device plan slug. To find the plan slug, visit [Device plans API docs](https://metal.equinix.com/developers/api/plans), set your auth token in the top of the page and see JSON from the API response",
				Required:    true,
				ForceNew:    true,
			},

			"billing_cycle": {
				Type:        schema.TypeString,
				Description: "monthly or hourly",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "The status of the device",
				Computed:    true,
			},

			"root_password": {
				Type:        schema.TypeString,
				Description: "Root password to the server (disabled after 24 hours)",
				Computed:    true,
				Sensitive:   true,
			},

			"locked": {
				Type:        schema.TypeBool,
				Description: "Whether the device is locked",
				Computed:    true,
			},

			"access_public_ipv6": {
				Type:        schema.TypeString,
				Description: "The ipv6 maintenance IP assigned to the device",
				Computed:    true,
			},

			"access_public_ipv4": {
				Type:        schema.TypeString,
				Description: "The ipv4 maintenance IP assigned to the device",
				Computed:    true,
			},

			"access_private_ipv4": {
				Type:        schema.TypeString,
				Description: "The ipv4 private IP assigned to the device",
				Computed:    true,
			},
			"network_type": {
				Type:        schema.TypeString,
				Description: "Network type of a device, used in [Layer 2 networking](https://metal.equinix.com/developers/docs/networking/layer2/). Will be one of " + NetworkTypeListHB,
				Computed:    true,
				Deprecated:  "You should handle Network Type with the new metal_device_network_type resource.",
			},

			"ports": {
				Type:        schema.TypeList,
				Description: "Ports assigned to the device",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the port (e.g. eth0, or bond0)",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "The ID of the device",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "One of [private_ipv4, public_ipv4, public_ipv6]",
							Computed:    true,
						},
						"mac": {
							Type:        schema.TypeString,
							Description: "MAC address assigned to the port",
							Computed:    true,
						},
						"bonded": {
							Type:        schema.TypeBool,
							Description: "Whether this port is part of a bond in bonded network setup",
							Computed:    true,
						},
					},
				},
			},

			"network": {
				Type:        schema.TypeList,
				Description: "The device's private and public IP (v4 and v6) network details. When a device is run without any special network configuration, it will have 3 addresses: public ipv4, private ipv4 and ipv6",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Description: "IPv4 or IPv6 address string",
							Computed:    true,
						},

						"gateway": {
							Type:        schema.TypeString,
							Description: "Address of router",
							Computed:    true,
						},

						"family": {
							Type:        schema.TypeInt,
							Description: "IP version - \"4\" or \"6\"",
							Computed:    true,
						},

						"cidr": {
							Type:        schema.TypeInt,
							Description: "CIDR suffix for IP address block to be assigned, i.e. amount of addresses",
							Computed:    true,
						},

						"public": {
							Type:        schema.TypeBool,
							Description: "Whether the address is routable from the Internet",
							Computed:    true,
						},
					},
				},
			},

			"created": {
				Type:        schema.TypeString,
				Description: "The timestamp for when the device was created",
				Computed:    true,
			},

			"updated": {
				Type:        schema.TypeString,
				Description: "The timestamp for the last time the device was updated",
				Computed:    true,
			},

			"user_data": {
				Type:        schema.TypeString,
				Description: "A string of the desired User Data for the device",
				Optional:    true,
				Sensitive:   true,
				ForceNew:    false,
			},

			"custom_data": {
				Type:        schema.TypeString,
				Description: "A string of the desired Custom Data for the device",
				Optional:    true,
				Sensitive:   true,
				ForceNew:    false,
			},

			"ipxe_script_url": {
				Type:        schema.TypeString,
				Description: "URL pointing to a hosted iPXE script. More",
				Optional:    true,
			},

			"always_pxe": {
				Type:        schema.TypeBool,
				Description: "If true, a device with OS custom_ipxe will",
				Optional:    true,
				Default:     false,
			},

			"deployed_hardware_reservation_id": {
				Type:        schema.TypeString,
				Description: "ID of hardware reservation where this device was deployed. It is useful when using the next-available hardware reservation",
				Computed:    true,
			},

			"hardware_reservation_id": {
				Type:        schema.TypeString,
				Description: "The UUID of the hardware reservation where you want this device deployed, or next-available if you want to pick your next available reservation automatically",
				Optional:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					dhwr, ok := d.GetOk("deployed_hardware_reservation_id")
					return ok && dhwr == new
				},
			},

			"tags": {
				Type:        schema.TypeList,
				Description: "Tags attached to the device",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"storage": {
				Type:        schema.TypeString,
				Description: "JSON for custom partitioning. Only usable on reserved hardware. More information in in the [Custom Partitioning and RAID](https://metal.equinix.com/developers/docs/servers/custom-partitioning-raid/) doc",
				Optional:    true,
				ForceNew:    true,
				StateFunc: func(v interface{}) string {
					s, _ := structure.NormalizeJsonString(v)
					return s
				},
				ValidateFunc: validation.StringIsJSON,
			},
			"project_ssh_key_ids": {
				Type:        schema.TypeList,
				Description: "Array of IDs of the project SSH keys which should be added to the device. If you omit this, SSH keys of all the members of the parent project will be added to the device. If you specify this array, only the listed project SSH keys will be added. Project SSH keys can be created with the [metal_project_ssh_key](project_ssh_key.md) resource",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ssh_key_ids": {
				Type:        schema.TypeList,
				Description: "List of IDs of SSH keys deployed in the device, can be both user and project SSH keys",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"wait_for_reservation_deprovision": {
				Type:        schema.TypeBool,
				Description: "Only used for devices in reserved hardware. If set, the deletion of this device will block until the hardware reservation is marked provisionable (about 4 minutes in August 2019)",
				Optional:    true,
				Default:     false,
				ForceNew:    false,
			},
			"force_detach_volumes": {
				Type:        schema.TypeBool,
				Description: "Delete device even if it has volumes attached. Only applies for destroy action",
				Optional:    true,
				Default:     false,
				ForceNew:    false,
			},
			"termination_time": {
				Type:        schema.TypeString,
				Description: "Timestamp for device termination. For example \"2021-09-03T16:32:00+03:00\". If you don't supply timezone info, timestamp is assumed to be in UTC.",
				Optional:    true,
				ForceNew:    false,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					_, err := time.ParseInLocation(time.RFC3339, val.(string), time.UTC)
					if err != nil {
						errs = []error{err}
					}
					return
				},
			},

			"reinstall": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Whether the device should be reinstalled instead of destroyed",
							Optional:    true,
							Default:     false,
						},
						"deprovision_fast": {
							Type:        schema.TypeBool,
							Description: "Whether the OS disk should be filled with `00h` bytes before reinstall",
							Optional:    true,
							Default:     false,
						},
						"preserve_data": {
							Type:        schema.TypeBool,
							Description: "Whether the non-OS disks should be kept or wiped during reinstall",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIf("custom_data", shouldReinstall),
			customdiff.ForceNewIf("operating_system", shouldReinstall),
			customdiff.ForceNewIf("user_data", shouldReinstall),
		),
	}
}

func shouldReinstall(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
	reinstall, ok := d.GetOk("reinstall")

	// Prior behaviour was to always destroy and create,
	// so in the event we can't get the reinstall config; let's
	// continue to do so.
	if !ok {
		return true
	}

	// We didn't get a reinstall configuration
	reinstall_list, ok := reinstall.([]interface{})
	if !ok {
		return true
	}

	reinstall_config, ok := reinstall_list[0].(map[string]interface{})

	// We didn't get a reinstall configuration
	if !ok {
		return true
	}

	return !reinstall_config["enabled"].(bool)
}

func resourceMetalDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	var addressTypesSlice []packngo.IPAddressCreateRequest
	_, ok := d.GetOk("ip_address")
	if ok {
		arr := d.Get("ip_address").([]interface{})
		addressTypesSlice = getNewIPAddressSlice(arr)
	}

	createRequest := &packngo.DeviceCreateRequest{
		Hostname:     d.Get("hostname").(string),
		Plan:         d.Get("plan").(string),
		IPAddresses:  addressTypesSlice,
		OS:           d.Get("operating_system").(string),
		BillingCycle: d.Get("billing_cycle").(string),
		ProjectID:    d.Get("project_id").(string),
	}

	facsRaw, facsOk := d.GetOk("facilities")
	metroRaw, metroOk := d.GetOk("metro")

	if !facsOk && !metroOk {
		return friendlyError(errors.New("one of facilies and metro must be configured"))
	}

	if facsOk {
		createRequest.Facility = convertStringArr(facsRaw.([]interface{}))
	}

	if metroOk {
		createRequest.Metro = metroRaw.(string)
	}

	if attr, ok := d.GetOk("user_data"); ok {
		createRequest.UserData = attr.(string)
	}

	if attr, ok := d.GetOk("custom_data"); ok {
		createRequest.CustomData = attr.(string)
	}

	if attr, ok := d.GetOk("ipxe_script_url"); ok {
		createRequest.IPXEScriptURL = attr.(string)
	}

	if attr, ok := d.GetOk("termination_time"); ok {
		tt, err := time.ParseInLocation(time.RFC3339, attr.(string), time.UTC)
		if err != nil {
			return err
		}
		createRequest.TerminationTime = &packngo.Timestamp{Time: tt}
	}

	if attr, ok := d.GetOk("hardware_reservation_id"); ok {
		createRequest.HardwareReservationID = attr.(string)
	} else {
		wfrd := "wait_for_reservation_deprovision"
		if d.Get(wfrd).(bool) {
			return friendlyError(fmt.Errorf("You can't set %s when not using a hardware reservation", wfrd))
		}
	}

	if createRequest.OS == "custom_ipxe" {
		if createRequest.IPXEScriptURL == "" && createRequest.UserData == "" {
			return friendlyError(errors.New("\"ipxe_script_url\" or \"user_data\"" +
				" must be provided when \"custom_ipxe\" OS is selected."))
		}

		// ipxe_script_url + user_data is OK, unless user_data is an ipxe script in
		// which case it's an error.
		if createRequest.IPXEScriptURL != "" {
			if matchIPXEScript.MatchString(createRequest.UserData) {
				return friendlyError(errors.New("\"user_data\" should not be an iPXE " +
					"script when \"ipxe_script_url\" is also provided."))
			}
		}
	}

	if createRequest.OS != "custom_ipxe" && createRequest.IPXEScriptURL != "" {
		return friendlyError(errors.New("\"ipxe_script_url\" argument provided, but" +
			" OS is not \"custom_ipxe\". Please verify and fix device arguments."))
	}

	if attr, ok := d.GetOk("always_pxe"); ok {
		createRequest.AlwaysPXE = attr.(bool)
	}

	projectKeys := d.Get("project_ssh_key_ids.#").(int)
	if projectKeys > 0 {
		createRequest.ProjectSSHKeys = convertStringArr(d.Get("project_ssh_key_ids").([]interface{}))
	}

	tags := d.Get("tags.#").(int)
	if tags > 0 {
		createRequest.Tags = convertStringArr(d.Get("tags").([]interface{}))
	}

	if attr, ok := d.GetOk("storage"); ok {
		s, err := structure.NormalizeJsonString(attr.(string))
		if err != nil {
			return errwrap.Wrapf("storage param contains invalid JSON: {{err}}", err)
		}
		var cpr packngo.CPR
		err = json.Unmarshal([]byte(s), &cpr)
		if err != nil {
			return errwrap.Wrapf("Error parsing Storage string: {{err}}", err)
		}
		createRequest.Storage = &cpr
	}

	newDevice, _, err := client.Devices.Create(createRequest)
	if err != nil {
		retErr := friendlyError(err)
		if isNotFound(retErr) {
			retErr = fmt.Errorf("%s, make sure project \"%s\" exists", retErr, createRequest.ProjectID)
		}
		return retErr
	}

	d.SetId(newDevice.ID)

	if err = waitForActiveDevice(d, meta); err != nil {
		return err
	}

	return resourceMetalDeviceRead(d, meta)
}

func resourceMetalDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	device, _, err := client.Devices.Get(d.Id(), deviceReadOptions)
	if err != nil {
		err = friendlyError(err)

		// If the device somehow already destroyed, mark as successfully gone.
		if isNotFound(err) {
			log.Printf("[WARN] Device (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("hostname", device.Hostname)
	d.Set("plan", device.Plan.Slug)
	d.Set("deployed_facility", device.Facility.Code)
	d.Set("facilities", []string{device.Facility.Code})
	if device.Metro != nil {
		d.Set("metro", device.Metro.Code)
	}
	d.Set("operating_system", device.OS.Slug)
	d.Set("state", device.State)
	d.Set("billing_cycle", device.BillingCycle)
	d.Set("locked", device.Locked)
	d.Set("created", device.Created)
	d.Set("updated", device.Updated)
	d.Set("ipxe_script_url", device.IPXEScriptURL)
	d.Set("always_pxe", device.AlwaysPXE)
	d.Set("root_password", device.RootPassword)
	d.Set("project_id", device.Project.ID)
	if device.Storage != nil {
		rawStorageBytes, err := json.Marshal(device.Storage)
		if err != nil {
			return fmt.Errorf("[ERR] Error getting storage JSON string for device (%s): %s", d.Id(), err)
		}

		storageString, err := structure.NormalizeJsonString(string(rawStorageBytes))
		if err != nil {
			return fmt.Errorf("[ERR] Error normalizing storage JSON string for device (%s): %s", d.Id(), err)
		}
		d.Set("storage", storageString)
	}
	if device.HardwareReservation != nil {
		d.Set("deployed_hardware_reservation_id", device.HardwareReservation.ID)
	}

	networkType := device.GetNetworkType()
	d.Set("network_type", networkType)

	wfrd := "wait_for_reservation_deprovision"
	if _, ok := d.GetOk(wfrd); !ok {
		d.Set(wfrd, nil)
	}
	fdv := "force_detach_volumes"
	if _, ok := d.GetOk(fdv); !ok {
		d.Set(fdv, nil)

		tt := "termination_time"
		if _, ok := d.GetOk(tt); !ok {
			d.Set(tt, nil)
		}
	}

	d.Set("tags", device.Tags)
	keyIDs := []string{}
	for _, k := range device.SSHKeys {
		keyIDs = append(keyIDs, path.Base(k.URL))
	}
	d.Set("ssh_key_ids", keyIDs)
	networkInfo := getNetworkInfo(device.Network)

	sort.SliceStable(networkInfo.Networks, func(i, j int) bool {
		famI := networkInfo.Networks[i]["family"].(int)
		famJ := networkInfo.Networks[j]["family"].(int)
		pubI := networkInfo.Networks[i]["public"].(bool)
		pubJ := networkInfo.Networks[j]["public"].(bool)
		return getNetworkRank(famI, pubI) < getNetworkRank(famJ, pubJ)
	})

	d.Set("network", networkInfo.Networks)
	d.Set("access_public_ipv4", networkInfo.PublicIPv4)
	d.Set("access_private_ipv4", networkInfo.PrivateIPv4)
	d.Set("access_public_ipv6", networkInfo.PublicIPv6)

	ports := getPorts(device.NetworkPorts)
	d.Set("ports", ports)

	if networkInfo.Host != "" {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": networkInfo.Host,
		})
	}

	return nil
}

func resourceMetalDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	if d.HasChange("locked") {
		var action func(string) (*packngo.Response, error)
		if d.Get("locked").(bool) {
			action = client.Devices.Lock
		} else {
			action = client.Devices.Unlock
		}
		if _, err := action(d.Id()); err != nil {
			return friendlyError(err)
		}
	}
	ur := packngo.DeviceUpdateRequest{}

	if d.HasChange("description") {
		dDesc := d.Get("description").(string)
		ur.Description = &dDesc
	}
	if d.HasChange("user_data") {
		dUserData := d.Get("user_data").(string)
		ur.UserData = &dUserData
	}
	if d.HasChange("custom_data") {
		dCustomData := d.Get("custom_data").(string)
		ur.CustomData = &dCustomData
	}
	if d.HasChange("hostname") {
		dHostname := d.Get("hostname").(string)
		ur.Hostname = &dHostname
	}
	if d.HasChange("tags") {
		ts := d.Get("tags")
		sts := []string{}

		switch ts.(type) {
		case []interface{}:
			for _, v := range ts.([]interface{}) {
				sts = append(sts, v.(string))
			}
			ur.Tags = &sts
		default:
			return friendlyError(fmt.Errorf("garbage in tags: %s", ts))
		}
	}
	if d.HasChange("ipxe_script_url") {
		dUrl := d.Get("ipxe_script_url").(string)
		ur.IPXEScriptURL = &dUrl
	}
	if d.HasChange("always_pxe") {
		dPXE := d.Get("always_pxe").(bool)
		ur.AlwaysPXE = &dPXE
	}
	if !reflect.DeepEqual(ur, packngo.DeviceUpdateRequest{}) {
		if _, _, err := client.Devices.Update(d.Id(), &ur); err != nil {
			return friendlyError(err)
		}

	}

	if d.HasChange("operating_system") || d.HasChange("user_data") || d.HasChange("custom_data") {
		reinstallOptions, err := getReinstallOptions(d)

		if err != nil {
			return friendlyError(err)
		}

		if _, err := client.Devices.Reinstall(d.Id(), &reinstallOptions); err != nil {
			return friendlyError(err)
		}

		if err = waitForActiveDevice(d, meta); err != nil {
			return err
		}
	}

	return resourceMetalDeviceRead(d, meta)
}

func getReinstallOptions(d *schema.ResourceData) (packngo.DeviceReinstallFields, error) {
	reinstall_list, ok := d.Get("reinstall").([]interface{})

	if !ok {
		return packngo.DeviceReinstallFields{}, fmt.Errorf("expected reinstall configuration and none available")
	}

	if len(reinstall_list) == 0 {
		return packngo.DeviceReinstallFields{}, fmt.Errorf("expected reinstall configuration and none available")
	}

	reinstall_config, ok := reinstall_list[0].(map[string]interface{})

	if !ok {
		return packngo.DeviceReinstallFields{}, fmt.Errorf("expected reinstall configuration and none available")
	}

	return packngo.DeviceReinstallFields{
		PreserveData:    reinstall_config["preserve_data"].(bool),
		DeprovisionFast: reinstall_config["deprovision_fast"].(bool),
	}, nil
}

func resourceMetalDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	fdvIf, fdvOk := d.GetOk("force_detach_volumes")
	fdv := false
	if fdvOk && fdvIf.(bool) {
		fdv = true
	}

	resp, err := client.Devices.Delete(d.Id(), fdv)
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}

	resId, resIdOk := d.GetOk("hardware_reservation_id")
	if resIdOk {
		wfrd, wfrdOK := d.GetOk("wait_for_reservation_deprovision")
		if wfrdOK && wfrd.(bool) {
			err := waitUntilReservationProvisionable(resId.(string), meta)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func waitForActiveDevice(d *schema.ResourceData, meta interface{}) error {
	// Wait for the device so we can get the networking attributes that show up after a while.
	state, err := waitForDeviceAttribute(d, []string{"active", "failed"}, []string{"queued", "provisioning", "reinstalling"}, "state", meta)
	if err != nil {
		d.SetId("")
		fErr := friendlyError(err)
		if isForbidden(fErr) {
			// If the device doesn't get to the active state, we can't recover it from here.

			return errors.New("provisioning time limit exceeded; the Equinix Metal team will investigate")
		}
		return fErr
	}

	if state != "active" {
		d.SetId("")
		return fmt.Errorf("Device in non-active state \"%s\"", state)
	}

	return nil
}
