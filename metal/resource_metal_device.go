package metal

import (
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

var matchIPXEScript = regexp.MustCompile(`(?i)^#![i]?pxe`)
var ipAddressTypes = []string{"public_ipv4", "private_ipv4", "public_ipv6"}

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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"operating_system": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"deployed_facility": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metro": {
				Type:          schema.TypeString,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
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
				Optional:    true,
				Description: "Inbound rules for this security group",
				Elem:        ipAddressSchema(),
				MinItems:    1,
			},

			"plan": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"billing_cycle": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"root_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"access_public_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_public_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_private_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "You should handle Network Type with the new metal_device_network_type resource.",
			},

			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bonded": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"family": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"cidr": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},

			"custom_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},

			"ipxe_script_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"always_pxe": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"deployed_hardware_reservation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hardware_reservation_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					dhwr, ok := d.GetOk("deployed_hardware_reservation_id")
					return ok && dhwr == new
				},
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"storage": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					s, _ := structure.NormalizeJsonString(v)
					return s
				},
				ValidateFunc: validation.StringIsJSON,
			},
			"project_ssh_key_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"wait_for_reservation_deprovision": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
			"force_detach_volumes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: false,
			},
		},
	}
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

	// Wait for the device so we can get the networking attributes that show up after a while.
	state, err := waitForDeviceAttribute(d, []string{"active", "failed"}, []string{"queued", "provisioning"}, "state", meta)
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
	/*
			    Possibly wait for device network state
		    	_, err := waitForDeviceAttribute(d, []string{"layer3"}, []string{"hybrid", "layer2-bonded", "layer2-individual"}, "network_type", meta)
		        if err != nil {
					return err
				}
	*/

	return resourceMetalDeviceRead(d, meta)
}

func resourceMetalDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	device, _, err := client.Devices.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project", "metro", "facility"}})
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
			return fmt.Errorf("[ERR] Errori normalizing storage JSON string for device (%s): %s", d.Id(), err)
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
	return resourceMetalDeviceRead(d, meta)
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
