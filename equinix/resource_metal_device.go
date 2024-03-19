package equinix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/converters"
	"github.com/equinix/terraform-provider-equinix/internal/network"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	matchIPXEScript = regexp.MustCompile(`(?i)^#![i]?pxe`)
	ipAddressTypes  = []string{"public_ipv4", "private_ipv4", "public_ipv6"}
)

var deviceCommonIncludes = []string{"project", "metro", "facility", "hardware_reservation"}

func resourceMetalDevice() *schema.Resource {
	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		CreateContext:      resourceMetalDeviceCreate,
		ReadWithoutTimeout: resourceMetalDeviceRead,
		UpdateContext:      resourceMetalDeviceUpdate,
		DeleteContext:      resourceMetalDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "The operating system slug. To find the slug, or visit [Operating Systems API docs](https://metal.equinix.com/developers/api/operatingsystems), set your API auth token in the top of the page and see JSON from the API response.  By default, changing this attribute will cause your device to be deleted and recreated.  If `reinstall` is enabled, the device will be updated in-place instead of recreated.",
				Required:    true,
				ForceNew:    false, // Computed; see CustomizeDiff below
			},

			"deployed_facility": {
				Type:        schema.TypeString,
				Description: "The facility where the device is deployed",
				Deprecated:  "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
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
				StateFunc: converters.ToLowerIf,
			},
			"facilities": {
				Type:        schema.TypeList,
				Description: "List of facility codes with deployment preferences. Equinix Metal API will go through the list and will deploy your device to first facility with free capacity. List items must be facility codes or any (a wildcard). To find the facility code, visit [Facilities API docs](https://metal.equinix.com/developers/api/facilities/), set your API auth token in the top of the page and see JSON from the API response. Conflicts with metro",
				Deprecated:  "Use metro instead of facilities.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				ForceNew:    true,
				MinItems:    1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					fsRaw := d.Get("facilities")
					fs := converters.IfArrToStringArr(fsRaw.([]interface{}))
					df := d.Get("deployed_facility").(string)
					if slices.Contains(fs, df) {
						return true
					}
					if slices.Contains(fs, "any") && (len(df) != 0) {
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(ipAddressTypes, false),
							Description:  fmt.Sprintf("one of %s", strings.Join(ipAddressTypes, ",")),
						},
						"cidr": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "CIDR suffix for IP block assigned to this device",
						},
						"reservation_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "IDs of reservations to pick the blocks from",
							MinItems:    1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
					},
				},
				MinItems: 1,
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
				Description: "Whether the device is locked or unlocked. Locking a device prevents you from deleting or reinstalling the device or performing a firmware update on the device, and it prevents an instance with a termination time set from being reclaimed, even if the termination time was reached",
				Optional:    true,
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
				Description: "Network type of a device, used in [Layer 2 networking](https://metal.equinix.com/developers/docs/networking/layer2/). Will be one of " + network.NetworkTypeListHB,
				Computed:    true,
				Deprecated:  "You should handle Network Type with one of 'equinix_metal_port' or 'equinix_metal_device_network_type' resources. See section 'Guides' for more info",
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
				Description: "A string of the desired User Data for the device.  By default, changing this attribute will cause the provider to destroy and recreate your device.  If `reinstall` is specified or `behavior.allow_changes` includes `\"user_data\"`, the device will be updated in-place instead of recreated.",
				Optional:    true,
				Sensitive:   true,
				ForceNew:    false, // Computed; see CustomizeDiff below
			},
			"custom_data": {
				Type:        schema.TypeString,
				Description: "A string of the desired Custom Data for the device.  By default, changing this attribute will cause the provider to destroy and recreate your device.  If `reinstall` is specified or `behavior.allow_changes` includes `\"custom_data\"`, the device will be updated in-place instead of recreated.",
				Optional:    true,
				Sensitive:   true,
				ForceNew:    false, // Computed; see CustomizeDiff below
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
				Description: "Array of IDs of the project SSH keys which should be added to the device. If you specify this array, only the listed project SSH keys (and any SSH keys for the users specified in user_ssh_key_ids) will be added. If no SSH keys are specified (both user_ssh_keys_ids and project_ssh_key_ids are empty lists or omitted), all parent project keys, parent project members keys and organization members keys will be included.  Project SSH keys can be created with the [equinix_metal_project_ssh_key](equinix_metal_project_ssh_key.md) resource",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_ssh_key_ids": {
				Type:        schema.TypeList,
				Description: "Array of IDs of the users whose SSH keys should be added to the device. If you specify this array, only the listed users' SSH keys (and any project SSH keys specified in project_ssh_key_ids) will be added. If no SSH keys are specified (both user_ssh_keys_ids and project_ssh_key_ids are empty lists or omitted), all parent project keys, parent project members keys and organization members keys will be included. User SSH keys can be created with the [equinix_metal_ssh_key](equinix_metal_ssh_key.md) resource",
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
			"behavior": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_changes": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
									attribute := val.(string)
									supportedAttributes := []string{"custom_data", "user_data"}
									if !slices.Contains(supportedAttributes, attribute) {
										errs = []error{fmt.Errorf("behavior.allow_changes was given %s, but only supports %v", attribute, supportedAttributes)}
									}
									return
								},
							},
							Description: "List of attributes that are allowed to change without recreating the instance. Supported attributes: `custom_data`, `user_data`",
							Optional:    true,
						},
					},
				},
			},
			"sos_hostname": {
				Type:        schema.TypeString,
				Description: "The hostname to use for [Serial over SSH](https://deploy.equinix.com/developers/docs/metal/resilience-recovery/serial-over-ssh/) access to the device",
				Computed:    true,
			},
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIf("custom_data", reinstallDisabledAndNoChangesAllowed("custom_data")),
			customdiff.ForceNewIf("operating_system", reinstallDisabled),
			customdiff.ForceNewIf("user_data", reinstallDisabledAndNoChangesAllowed("user_data")),
		),
	}
}

// This method returns true if reinstall is disabled, and false if it is enabled.
// This is used to set ForceNew to true when reinstall is disabled
func reinstallDisabled(_ context.Context, d *schema.ResourceDiff, meta interface{}) bool {
	reinstall, ok := d.GetOk("reinstall")

	if !ok {
		// There is no reinstall attribute, so ForceNew should be true
		return true
	}

	// To reach this point, the device config had to include a `reinstall` block,
	// so we can assume all necessary parts of that block are filled in
	reinstall_list := reinstall.([]interface{})
	reinstall_config := reinstall_list[0].(map[string]interface{})

	return !reinstall_config["enabled"].(bool)
}

func reinstallDisabledAndNoChangesAllowed(attribute string) customdiff.ResourceConditionFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
		if reinstallDisabled(ctx, d, meta) {
			// If reinstall is disabled, we need to see if ForceNew
			// should be disabled due to behavior settings
			behavior, ok := d.GetOk("behavior")

			if !ok {
				// This means reinstall is disabled and there is no behavior
				// attribute, so ForceNew should be true
				return true
			}

			// To reach this point, the device config had to include a `behavior`
			// block, so we can assume all necessary parts of that block are filled in
			behavior_list := behavior.([]interface{})
			behavior_config := behavior_list[0].(map[string]interface{})

			allow_changes := converters.IfArrToStringArr(behavior_config["allow_changes"].([]interface{}))

			// This means we got a valid behavior specification, so we set ForceNew
			// to true if behavior.allow_changes includes the attribute that is changing
			return !slices.Contains(allow_changes, attribute)
		}

		// This means reinstall is enabled, so it doesn't matter what the behavior
		// says; ForceNew should not be set to true in this case
		return false
	}
}

func resourceMetalDeviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	createRequest := metalv1.CreateDeviceRequest{}

	facsRaw, facsOk := d.GetOk("facilities")
	metroRaw, metroOk := d.GetOk("metro")

	if !facsOk && !metroOk {
		return diag.Errorf("one of facilies and metro must be configured")
	}

	if facsOk {
		facilityRequest := &metalv1.DeviceCreateInFacilityInput{
			Facility: converters.IfArrToStringArr(facsRaw.([]interface{})),
		}

		diagErr := setupDeviceCreateRequest(d, facilityRequest)
		if diagErr != nil {
			return diagErr
		}

		createRequest.DeviceCreateInFacilityInput = facilityRequest
	}

	if metroOk {
		metroRequest := &metalv1.DeviceCreateInMetroInput{
			Metro: metroRaw.(string),
		}

		diagErr := setupDeviceCreateRequest(d, metroRequest)
		if diagErr != nil {
			return diagErr
		}

		createRequest.DeviceCreateInMetroInput = metroRequest
	}

	start := time.Now()
	projectID := d.Get("project_id").(string)
	newDevice, _, err := client.DevicesApi.CreateDevice(ctx, projectID).CreateDeviceRequest(createRequest).Execute()
	if err != nil {
		retErr := equinix_errors.FriendlyError(err)
		if equinix_errors.IsNotFound(retErr) {
			retErr = fmt.Errorf("%s, make sure project \"%s\" exists", retErr, projectID)
		}
		return diag.FromErr(retErr)
	}

	d.SetId(newDevice.GetId())

	createTimeout := d.Timeout(schema.TimeoutCreate) - 30*time.Second - time.Since(start)
	if err = waitForActiveDevice(ctx, d, meta, createTimeout); err != nil {
		return diag.FromErr(err)
	}

	return resourceMetalDeviceRead(ctx, d, meta)
}

func resourceMetalDeviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	device, resp, err := client.DevicesApi.FindDeviceById(ctx, d.Id()).Include(deviceCommonIncludes).Execute()
	if err != nil {
		err = equinix_errors.FriendlyErrorForMetalGo(err, resp)

		// If the device somehow already destroyed, mark as successfully gone.
		// Checking d.IsNewResource prevents the creation of a resource from failing
		// silently. Note d.IsNewResource is false in resource import operations.
		if !d.IsNewResource() && (equinix_errors.IsNotFound(err) || equinix_errors.IsForbidden(err)) {
			log.Printf("[WARN] Device (%s) not found or in failed status, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.Set("hostname", device.GetHostname())
	d.Set("plan", device.Plan.GetSlug())
	d.Set("deployed_facility", device.Facility.GetCode())
	d.Set("facilities", []string{device.Facility.GetCode()})
	if device.Metro != nil {
		d.Set("metro", device.Metro.GetCode())
	}
	d.Set("operating_system", device.OperatingSystem.GetSlug())
	d.Set("state", device.GetState())
	d.Set("billing_cycle", device.GetBillingCycle())
	d.Set("locked", device.GetLocked())
	d.Set("created", device.GetCreatedAt().Format(time.RFC3339))
	d.Set("updated", device.GetUpdatedAt().Format(time.RFC3339))
	d.Set("ipxe_script_url", device.GetIpxeScriptUrl())
	d.Set("always_pxe", device.GetAlwaysPxe())
	d.Set("root_password", device.GetRootPassword())
	d.Set("project_id", device.Project.GetId())
	d.Set("sos_hostname", device.GetSos())
	if device.Storage != nil {
		rawStorageBytes, err := json.Marshal(device.Storage)
		if err != nil {
			return diag.Errorf("[ERR] Error getting storage JSON string for device (%s): %s", d.Id(), err)
		}

		storageString, err := structure.NormalizeJsonString(string(rawStorageBytes))
		if err != nil {
			return diag.Errorf("[ERR] Error normalizing storage JSON string for device (%s): %s", d.Id(), err)
		}
		d.Set("storage", storageString)
	}
	if device.HardwareReservation != nil {
		d.Set("deployed_hardware_reservation_id", device.HardwareReservation.GetId())
	}

	networkType, err := getNetworkType(device)
	if err != nil {
		return diag.Errorf("[ERR] Error computing network type for device (%s): %s", d.Id(), err)
	}
	d.Set("network_type", networkType)

	wfrd := "wait_for_reservation_deprovision"
	if _, ok := d.GetOk(wfrd); !ok {
		d.Set(wfrd, nil)
	}
	fdv := "force_detach_volumes"
	if _, ok := d.GetOk(fdv); !ok {
		d.Set(fdv, nil)
	}
	tt := "termination_time"
	if _, ok := d.GetOk(tt); !ok {
		d.Set(tt, nil)
	}

	d.Set("tags", device.Tags)
	keyIDs := []string{}
	for _, k := range device.SshKeys {
		keyIDs = append(keyIDs, path.Base(k.Href))
	}
	d.Set("ssh_key_ids", keyIDs)
	networkInfo := getNetworkInfo(device.IpAddresses)

	sort.SliceStable(networkInfo.Networks, func(i, j int) bool {
		famI := networkInfo.Networks[i]["family"].(int32)
		famJ := networkInfo.Networks[j]["family"].(int32)
		pubI := networkInfo.Networks[i]["public"].(bool)
		pubJ := networkInfo.Networks[j]["public"].(bool)
		return getNetworkRank(int(famI), pubI) < getNetworkRank(int(famJ), pubJ)
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

func resourceMetalDeviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	ur := metalv1.DeviceUpdateInput{}

	if d.HasChange("locked") {
		ur.Locked = metalv1.PtrBool(d.Get("locked").(bool))
	}

	if d.HasChange("description") {
		dDesc := d.Get("description").(string)
		ur.Description = &dDesc
	}
	if d.HasChange("user_data") {
		dUserData := d.Get("user_data").(string)
		ur.Userdata = &dUserData
	}
	if d.HasChange("custom_data") {
		var customdata map[string]interface{}
		err := json.Unmarshal([]byte(d.Get("custom_data").(string)), &customdata)
		if err != nil {
			return diag.Errorf("error reading custom_data from state: %v", err)
		}
		ur.Customdata = customdata
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
			ur.Tags = sts
		default:
			return diag.Errorf("garbage in tags: %s", ts)
		}
	}
	if d.HasChange("ipxe_script_url") {
		dUrl := d.Get("ipxe_script_url").(string)
		ur.IpxeScriptUrl = &dUrl
	}
	if d.HasChange("always_pxe") {
		dPXE := d.Get("always_pxe").(bool)
		ur.AlwaysPxe = &dPXE
	}

	start := time.Now()
	if !reflect.DeepEqual(ur, metalv1.DeviceUpdateInput{}) {
		if _, _, err := client.DevicesApi.UpdateDevice(ctx, d.Id()).DeviceUpdateInput(ur).Execute(); err != nil {
			return diag.FromErr(equinix_errors.FriendlyError(err))
		}
	}

	if err := doReinstall(ctx, client, d, meta, start); err != nil {
		return diag.FromErr(err)
	}

	return resourceMetalDeviceRead(ctx, d, meta)
}

func doReinstall(ctx context.Context, client *metalv1.APIClient, d *schema.ResourceData, meta interface{}, start time.Time) error {
	if d.HasChange("operating_system") || d.HasChange("user_data") || d.HasChange("custom_data") {
		reinstall, ok := d.GetOk("reinstall")

		if !ok {
			// Assume we're here because behavior.allow_changes was set (not an error)
			return nil
		}

		reinstall_list := reinstall.([]interface{})
		reinstall_config := reinstall_list[0].(map[string]interface{})

		if !reinstall_config["enabled"].(bool) {
			// This means a reinstall block was provided, but reinstall was explicitly
			// disabled.  Assume we're here because behavior.allow_changes was set (not an error)
			return nil
		}

		reinstallOptions := metalv1.DeviceActionInput{
			Type:            metalv1.DEVICEACTIONINPUTTYPE_REINSTALL,
			OperatingSystem: metalv1.PtrString(d.Get("operating_system").(string)),
			PreserveData:    metalv1.PtrBool(reinstall_config["preserve_data"].(bool)),
			DeprovisionFast: metalv1.PtrBool(reinstall_config["deprovision_fast"].(bool)),
		}

		if _, err := client.DevicesApi.PerformAction(ctx, d.Id()).DeviceActionInput(reinstallOptions).Execute(); err != nil {
			return equinix_errors.FriendlyError(err)
		}

		updateTimeout := d.Timeout(schema.TimeoutUpdate) - 30*time.Second - time.Since(start)
		if err := waitForActiveDevice(ctx, d, meta, updateTimeout); err != nil {
			return err
		}
	}

	return nil
}

func resourceMetalDeviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	fdvIf, fdvOk := d.GetOk("force_detach_volumes")
	fdv := false
	if fdvOk && fdvIf.(bool) {
		fdv = true
	}

	start := time.Now()

	resp, err := client.DevicesApi.DeleteDevice(ctx, d.Id()).ForceDelete(fdv).Execute()
	if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	}

	resId, resIdOk := d.GetOk("deployed_hardware_reservation_id")
	if resIdOk {
		wfrd, wfrdOK := d.GetOk("wait_for_reservation_deprovision")
		if wfrdOK && wfrd.(bool) {
			// avoid "context: deadline exceeded"
			timeout := d.Timeout(schema.TimeoutDelete) - 30*time.Second - time.Since(start)

			err := waitUntilReservationProvisionable(ctx, client, resId.(string), d.Id(), 10*time.Second, timeout, 3*time.Second)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

func waitForActiveDevice(ctx context.Context, d *schema.ResourceData, meta interface{}, timeout time.Duration) error {
	targets := []string{"active", "failed"}
	pending := []string{"queued", "provisioning", "reinstalling"}

	stateConf := &retry.StateChangeConf{
		Pending: pending,
		Target:  targets,
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).NewMetalClientForSDK(d)

			device, _, err := client.DevicesApi.FindDeviceById(ctx, d.Id()).Include([]string{"project"}).Execute()
			if err == nil {
				retAttrVal := fmt.Sprint(device.GetState())
				return retAttrVal, retAttrVal, nil
			}
			return "error", "error", err
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	// Wait for the device so we can get the networking attributes that show up after a while.
	state, err := waitForDeviceAttribute(ctx, d, stateConf)
	if err != nil {
		d.SetId("")
		fErr := equinix_errors.FriendlyError(err)
		if equinix_errors.IsForbidden(fErr) {
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

type deviceCreateRequest interface {
	SetUserdata(string)
	GetUserdata() string
	SetCustomdata(map[string]interface{})
	SetAlwaysPxe(bool)
	SetIpxeScriptUrl(string)
	GetIpxeScriptUrl() string
	SetTerminationTime(time.Time)
	SetHardwareReservationId(string)
	SetBillingCycle(metalv1.DeviceCreateInputBillingCycle)
	GetOperatingSystem() string
	SetProjectSshKeys([]string)
	SetUserSshKeys([]string)
	SetTags([]string)
	SetStorage(metalv1.Storage)
	SetHostname(string)
	SetPlan(string)
	SetOperatingSystem(string)
	SetIpAddresses([]metalv1.IPAddress)
	SetLocked(bool)
}

func setupDeviceCreateRequest(d *schema.ResourceData, createRequest deviceCreateRequest) diag.Diagnostics {
	var addressTypesSlice []metalv1.IPAddress
	_, ok := d.GetOk("ip_address")
	if ok {
		arr := d.Get("ip_address").([]interface{})

		addressTypesSlice = getNewIPAddressSlice(arr)
	}

	if hostname, ok := d.GetOk("hostname"); ok {
		createRequest.SetHostname(hostname.(string))
	}

	createRequest.SetPlan(d.Get("plan").(string))
	createRequest.SetIpAddresses(addressTypesSlice)
	createRequest.SetOperatingSystem(d.Get("operating_system").(string))

	if rawBillingCycle, ok := d.GetOk("billing_cycle"); ok {
		billingCycle, err := metalv1.NewDeviceCreateInputBillingCycleFromValue(rawBillingCycle.(string))
		if err != nil {
			return diag.Errorf("unknown billing cycle: %v", err)
		}

		createRequest.SetBillingCycle(*billingCycle)
	}

	if attr, ok := d.GetOk("user_data"); ok {
		createRequest.SetUserdata(attr.(string))
	}

	if attr, ok := d.GetOk("custom_data"); ok {
		var customdata map[string]interface{}
		err := json.Unmarshal([]byte(attr.(string)), &customdata)
		if err != nil {
			return diag.FromErr(err)
		}
		createRequest.SetCustomdata(customdata)
	}

	if attr, ok := d.GetOk("ipxe_script_url"); ok {
		createRequest.SetIpxeScriptUrl(attr.(string))
	}

	if attr, ok := d.GetOk("termination_time"); ok {
		tt, err := time.ParseInLocation(time.RFC3339, attr.(string), time.UTC)
		if err != nil {
			return diag.FromErr(err)
		}
		createRequest.SetTerminationTime(tt)
	}

	if attr, ok := d.GetOk("hardware_reservation_id"); ok {
		createRequest.SetHardwareReservationId(attr.(string))
	} else {
		wfrd := "wait_for_reservation_deprovision"
		if d.Get(wfrd).(bool) {
			return diag.Errorf("You can't set %s when not using a hardware reservation", wfrd)
		}
	}

	if attr, ok := d.GetOk("locked"); ok {
		createRequest.SetLocked(attr.(bool))
	}

	if createRequest.GetOperatingSystem() == "custom_ipxe" {
		if createRequest.GetIpxeScriptUrl() == "" && createRequest.GetUserdata() == "" {
			return diag.Errorf("\"ipxe_script_url\" or \"user_data\"" +
				" must be provided when \"custom_ipxe\" OS is selected.")
		}

		// ipxe_script_url + user_data is OK, unless user_data is an ipxe script in
		// which case it's an error.
		if createRequest.GetIpxeScriptUrl() != "" {
			if matchIPXEScript.MatchString(createRequest.GetUserdata()) {
				return diag.Errorf("\"user_data\" should not be an iPXE " +
					"script when \"ipxe_script_url\" is also provided.")
			}
		}
	}

	if createRequest.GetOperatingSystem() != "custom_ipxe" && createRequest.GetIpxeScriptUrl() != "" {
		return diag.Errorf("\"ipxe_script_url\" argument provided, but" +
			" OS is not \"custom_ipxe\". Please verify and fix device arguments.")
	}

	if attr, ok := d.GetOk("always_pxe"); ok {
		createRequest.SetAlwaysPxe(attr.(bool))
	}

	projectKeys := d.Get("project_ssh_key_ids.#").(int)
	if projectKeys > 0 {
		createRequest.SetProjectSshKeys(converters.IfArrToStringArr(d.Get("project_ssh_key_ids").([]interface{})))
	}

	userKeys := d.Get("user_ssh_key_ids.#").(int)
	if userKeys > 0 {
		createRequest.SetUserSshKeys(converters.IfArrToStringArr(d.Get("user_ssh_key_ids").([]interface{})))
	}

	tags := d.Get("tags.#").(int)
	if tags > 0 {
		createRequest.SetTags(converters.IfArrToStringArr(d.Get("tags").([]interface{})))
	}

	if attr, ok := d.GetOk("storage"); ok {
		s, err := structure.NormalizeJsonString(attr.(string))
		if err != nil {
			return diag.Errorf("storage param contains invalid JSON: %s", err)
		}
		var storage metalv1.Storage
		err = json.Unmarshal([]byte(s), &storage)
		if err != nil {
			return diag.Errorf("error parsing Storage string: %s", err)
		}
		createRequest.SetStorage(storage)
	}

	return nil
}

func getNewIPAddressSlice(arr []interface{}) []metalv1.IPAddress {
	addressTypesSlice := make([]metalv1.IPAddress, len(arr))

	for i, m := range arr {
		addressTypesSlice[i] = ifToIPCreateRequest(m)
	}
	return addressTypesSlice
}

func ifToIPCreateRequest(m interface{}) metalv1.IPAddress {
	iacr := metalv1.IPAddress{}
	ia := m.(map[string]interface{})
	at := ia["type"].(string)
	switch at {
	case "public_ipv4":
		iacr.SetAddressFamily(4)
		iacr.SetPublic(true)
	case "private_ipv4":
		iacr.SetAddressFamily(4)
		iacr.SetPublic(false)
	case "public_ipv6":
		iacr.SetAddressFamily(6)
		iacr.SetPublic(true)
	}
	if cidr := ia["cidr"].(int); cidr > 0 {
		iacr.SetCidr(int32(cidr))
	}
	iacr.SetIpReservations(converters.IfArrToStringArr(ia["reservation_ids"].([]interface{})))
	return iacr
}
