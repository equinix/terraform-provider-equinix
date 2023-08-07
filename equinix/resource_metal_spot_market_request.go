package equinix

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalSpotMarketRequest() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalSpotMarketRequestCreate,
		Read:   resourceMetalSpotMarketRequestRead,
		Delete: resourceMetalSpotMarketRequestDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"devices_min": {
				Type:        schema.TypeInt,
				Description: "Miniumum number devices to be created",
				Required:    true,
				ForceNew:    true,
			},
			"devices_max": {
				Type:        schema.TypeInt,
				Description: "Maximum number devices to be created",
				Required:    true,
				ForceNew:    true,
			},
			"max_bid_price": {
				Type:        schema.TypeFloat,
				Description: "Maximum price user is willing to pay per hour per device",
				Required:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldF, err := strconv.ParseFloat(old, 64)
					if err != nil {
						return false
					}
					newF, err := strconv.ParseFloat(new, 64)
					if err != nil {
						return false
					}
					// suppress diff if the difference between existing and new bid price
					// is less than 2%
					diffThreshold := .02
					priceDiff := oldF / newF

					if diffThreshold < priceDiff {
						return true
					}
					return false
				},
			},
			"facilities": {
				Type:          schema.TypeList,
				Description:   "Facility IDs where devices should be created",
				Deprecated:    "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"metro"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldData, newData := d.GetChange("facilities")

					// If this function is called and oldData or newData is nil,
					// then the attribute changed
					if oldData == nil || newData == nil {
						return false
					}

					oldArray := oldData.([]interface{})
					newArray := newData.([]interface{})

					// If the number of items in the list is different,
					// then the attribute changed
					if len(oldArray) != len(newArray) {
						return false
					}

					// Convert data to string arrays
					oldFacilities := make([]string, len(oldArray))
					newFacilities := make([]string, len(newArray))
					for i, oldFacility := range oldArray {
						oldFacilities[i] = fmt.Sprint(oldFacility)
					}
					for j, newFacility := range newArray {
						newFacilities[j] = fmt.Sprint(newFacility)
					}
					// Sort the old and new arrays so that we don't show a diff
					// if the facilities are the same but the order is different
					sort.Strings(oldFacilities)
					sort.Strings(newFacilities)
					return reflect.DeepEqual(oldFacilities, newFacilities)
				},
			},
			"metro": {
				Type:          schema.TypeString,
				Description:   "Metro where devices should be created",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"facilities"},
				StateFunc:     toLower,
			},
			"instance_parameters": {
				Type:        schema.TypeList,
				Description: "Parameters for devices provisioned from this request. You can find the parameter description from the [equinix_metal_device doc](device.md)",
				Required:    true,
				MaxItems:    1,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"billing_cycle": {
							Type:     schema.TypeString,
							Required: true,
						},
						"plan": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operating_system": {
							Type:     schema.TypeString,
							Required: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"termintation_time": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "Use instance_parameters.termination_time instead",
						},
						"termination_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"always_pxe": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"features": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"locked": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"project_ssh_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"user_ssh_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"userdata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"customdata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipxe_script_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "Project ID",
				Required:    true,
				ForceNew:    true,
			},
			"wait_for_devices": {
				Type:        schema.TypeBool,
				Description: "On resource creation - wait until all desired devices are active, on resource destruction - wait until devices are removed",
				Optional:    true,
				ForceNew:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceMetalSpotMarketRequestCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal
	var waitForDevices bool

	metro := d.Get("metro").(string)

	facilitiesRaw := d.Get("facilities").([]interface{})
	facilities := []string{}

	for _, f := range facilitiesRaw {
		facilities = append(facilities, f.(string))
	}

	params := packngo.SpotMarketRequestInstanceParameters{
		Hostname:        d.Get("instance_parameters.0.hostname").(string),
		BillingCycle:    d.Get("instance_parameters.0.billing_cycle").(string),
		Plan:            d.Get("instance_parameters.0.plan").(string),
		OperatingSystem: d.Get("instance_parameters.0.operating_system").(string),
	}

	if val, ok := d.GetOk("instance_parameters.0.userdata"); ok {
		params.UserData = val.(string)
	}

	if val, ok := d.GetOk("instance_parameters.0.customdata"); ok {
		params.CustomData = val.(string)
	}

	if val, ok := d.GetOk("instance_parameters.0.ipxe_script_url"); ok {
		params.IPXEScriptURL = val.(string)
	}

	if val, ok := d.GetOk("instance_parameters.0.always_pxe"); ok {
		params.AlwaysPXE = val.(bool)
	}

	if params.OperatingSystem == "custom_ipxe" {
		if params.IPXEScriptURL == "" && params.UserData == "" {
			return fmt.Errorf("\"ipxe_script_url\" or \"user_data\"" +
				" must be provided when \"custom_ipxe\" OS is selected.")
		}

		// ipxe_script_url + user_data is OK, unless user_data is an ipxe script in
		// which case it's an error.
		if params.IPXEScriptURL != "" {
			if matchIPXEScript.MatchString(params.UserData) {
				return fmt.Errorf("\"user_data\" should not be an iPXE " +
					"script when \"ipxe_script_url\" is also provided.")
			}
		}
	}

	if params.OperatingSystem != "custom_ipxe" && params.IPXEScriptURL != "" {
		return fmt.Errorf("\"ipxe_script_url\" argument provided, but" +
			" OS is not \"custom_ipxe\". Please verify and fix device arguments.")
	}

	if val, ok := d.GetOk("instance_parameters.0.description"); ok {
		params.Description = val.(string)
	}

	if val, ok := d.GetOk("instance_parameters.0.features"); ok {
		temp := val.([]interface{})
		for _, i := range temp {
			if i != nil {
				params.Features = append(params.Features, i.(string))
			}
		}
	}

	if val, ok := d.GetOk("wait_for_devices"); ok {
		waitForDevices = val.(bool)
	}

	if val, ok := d.GetOk("instance_parameters.0.locked"); ok {
		params.Locked = val.(bool)
	}

	if val, ok := d.GetOk("instance_parameters.0.project_ssh_keys"); ok {
		temp := val.([]interface{})
		for _, i := range temp {
			if i != nil {
				params.ProjectSSHKeys = append(params.ProjectSSHKeys, i.(string))
			}
		}
	}

	if val, ok := d.GetOk("instance_parameters.0.tags"); ok {
		temp := val.([]interface{})
		for _, i := range temp {
			if i != nil {
				params.Tags = append(params.Tags, i.(string))
			}
		}
	}

	if val, ok := d.GetOk("instance_parameters.0.user_ssh_keys"); ok {
		temp := val.([]interface{})
		for _, i := range temp {
			if i != nil {
				params.UserSSHKeys = append(params.UserSSHKeys, i.(string))
			}
		}
	}

	smrc := &packngo.SpotMarketRequestCreateRequest{
		DevicesMax:  d.Get("devices_max").(int),
		DevicesMin:  d.Get("devices_min").(int),
		MaxBidPrice: d.Get("max_bid_price").(float64),
		FacilityIDs: facilities,
		Metro:       metro,
		Parameters:  params,
	}

	smr, _, err := client.SpotMarketRequests.Create(smrc, d.Get("project_id").(string))
	if err != nil {
		return err
	}

	d.SetId(smr.ID)

	if waitForDevices {
		stateConf := &retry.StateChangeConf{
			Pending:        []string{"not_done"},
			Target:         []string{"done"},
			Refresh:        resourceStateRefreshFunc(d, meta),
			Timeout:        d.Timeout(schema.TimeoutCreate),
			MinTimeout:     5 * time.Second,
			Delay:          3 * time.Second, // Wait 10 secs before starting
			NotFoundChecks: 600,             // Setting high number, to support long timeouts
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}
	}

	return resourceMetalSpotMarketRequestRead(d, meta)
}

func resourceMetalSpotMarketRequestRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal

	smr, _, err := client.SpotMarketRequests.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			log.Printf("[WARN] SpotMarketRequest (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	metro := ""
	if smr.Metro != nil {
		metro = smr.Metro.Code
	}

	return setMap(d, map[string]interface{}{
		"metro":         metro,
		"project_id":    smr.Project.ID,
		"devices_min":   smr.DevicesMin,
		"devices_max":   smr.DevicesMax,
		"max_bid_price": smr.MaxBidPrice,
		"facilities": func(d *schema.ResourceData, k string) error {
			facilityIDs := make([]string, len(smr.Facilities))
			facilityCodes := make([]string, len(smr.Facilities))
			if len(smr.Facilities) > 0 {
				for i, f := range smr.Facilities {
					facilityIDs[i] = f.ID
					facilityCodes[i] = f.Code
				}
			}
			return d.Set(k, facilityCodes)
		},
	})
}

func resourceMetalSpotMarketRequestDelete(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal
	var waitForDevices bool

	if val, ok := d.GetOk("wait_for_devices"); ok {
		waitForDevices = val.(bool)
	}
	if waitForDevices {
		smr, _, err := client.SpotMarketRequests.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
		if err != nil {
			return nil
		}

		stateConf := &retry.StateChangeConf{
			Pending:        []string{"not_done"},
			Target:         []string{"done"},
			Refresh:        resourceStateRefreshFunc(d, meta),
			Timeout:        d.Timeout(schema.TimeoutDelete),
			MinTimeout:     5 * time.Second,
			Delay:          3 * time.Second, // Wait 10 secs before starting
			NotFoundChecks: 600,             // Setting high number, to support long timeouts
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}

		for _, d := range smr.Devices {
			resp, err := client.Devices.Delete(d.ID, true)
			if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
				return err
			}
		}
	}
	resp, err := client.SpotMarketRequests.Delete(d.Id(), true)
	return ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err)
}

type InstanceParams []map[string]interface{}

func getInstanceParams(params *packngo.SpotMarketRequestInstanceParameters) InstanceParams {
	p := InstanceParams{{
		"billing_cycle":    params.BillingCycle,
		"plan":             params.Plan,
		"operating_system": params.OperatingSystem,
		"hostname":         params.Hostname,
		"always_pxe":       params.AlwaysPXE,
		"description":      params.Description,
		"features":         params.Features,
		"locked":           params.Locked,
		"project_ssh_keys": params.ProjectSSHKeys,
		"user_ssh_keys":    params.UserSSHKeys,
		"userdata":         params.UserData,
		"customdata":       params.CustomData,
		"ipxe_script_url":  params.IPXEScriptURL,
		"tags":             params.Tags,
	}}
	if params.TerminationTime != nil {
		p[0]["termintation_time"] = params.TerminationTime.Time
		p[0]["termination_time"] = params.TerminationTime.Time
	}
	return p
}

func resourceStateRefreshFunc(d *schema.ResourceData, meta interface{}) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		meta.(*Config).addModuleToMetalUserAgent(d)
		client := meta.(*Config).metal

		smr, _, err := client.SpotMarketRequests.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
		if err != nil {
			return nil, "", fmt.Errorf("Failed to fetch Spot market request with following error: %s", err.Error())
		}
		var finished bool

		for _, d := range smr.Devices {

			dev, _, err := client.Devices.Get(d.ID, nil)
			if err != nil {
				return nil, "", fmt.Errorf("Failed to fetch Device with following error: %s", err.Error())
			}
			if dev.State != "active" {
				break
			} else {
				finished = true
			}
		}
		if finished {
			return smr, "done", nil
		}
		return nil, "not_done", nil
	}
}
