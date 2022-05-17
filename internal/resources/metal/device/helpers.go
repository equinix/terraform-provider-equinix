package device

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

var (
	wgMap   = map[string]*sync.WaitGroup{}
	wgMutex = sync.Mutex{}

	DeviceNetworkTypes   = []string{"layer3", "hybrid", "layer2-individual", "layer2-bonded"}
	DeviceNetworkTypesHB = []string{"layer3", "hybrid", "hybrid-bonded", "layer2-individual", "layer2-bonded"}
	NetworkTypeList      = strings.Join(DeviceNetworkTypes, ", ")
	NetworkTypeListHB    = strings.Join(DeviceNetworkTypesHB, ", ")
)

const (
	deprovisioning = "deprovisioning"
	provisionable  = "provisionable"
	reprovisioned  = "reprovisioned"
	errstate       = "error"
)

type NetworkInfo struct {
	Networks       []map[string]interface{}
	IPv4SubnetSize int
	Host           string
	PublicIPv4     string
	PublicIPv6     string
	PrivateIPv4    string
}

func ifToIPCreateRequest(m interface{}) packngo.IPAddressCreateRequest {
	iacr := packngo.IPAddressCreateRequest{}
	ia := m.(map[string]interface{})
	at := ia["type"].(string)
	switch at {
	case "public_ipv4":
		iacr.AddressFamily = 4
		iacr.Public = true
	case "private_ipv4":
		iacr.AddressFamily = 4
		iacr.Public = false
	case "public_ipv6":
		iacr.AddressFamily = 6
		iacr.Public = true
	}
	iacr.CIDR = ia["cidr"].(int)
	iacr.Reservations = converters.ConvertStringArr(ia["reservation_ids"].([]interface{}))
	return iacr
}

func getNewIPAddressSlice(arr []interface{}) []packngo.IPAddressCreateRequest {
	addressTypesSlice := make([]packngo.IPAddressCreateRequest, len(arr))

	for i, m := range arr {
		addressTypesSlice[i] = ifToIPCreateRequest(m)
	}
	return addressTypesSlice
}

func getPorts(ps []packngo.Port) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 1)
	for _, p := range ps {
		port := map[string]interface{}{
			"name":   p.Name,
			"id":     p.ID,
			"type":   p.Type,
			"mac":    p.Data.MAC,
			"bonded": p.Data.Bonded,
		}
		ret = append(ret, port)
	}
	return ret
}

func ipAddressSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func waitForDeviceAttribute(d *schema.ResourceData, targets []string, pending []string, attribute string, meta interface{}) (string, error) {
	wg := getWaitForDeviceLock(d.Id())
	wg.Wait()

	wgMutex.Lock()
	wg.Add(1)
	wgMutex.Unlock()

	defer func() {
		wgMutex.Lock()
		wg.Done()
		wgMutex.Unlock()
	}()

	if attribute != "state" && attribute != "network_type" {
		return "", fmt.Errorf("unsupported attr to wait for: %s", attribute)
	}

	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  targets,
		Refresh: func() (interface{}, string, error) {
			client := meta.(*config.Config).MetalClient
			device, _, err := client.Devices.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project"}})
			if err == nil {
				retAttrVal := device.State
				if attribute == "network_type" {
					networkType := device.GetNetworkType()
					retAttrVal = networkType
				}
				return retAttrVal, retAttrVal, nil
			}
			return "error", "error", err
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	attrValRaw, err := stateConf.WaitForState()

	if v, ok := attrValRaw.(string); ok {
		return v, err
	}

	return "", err
}

func waitUntilReservationProvisionable(client *packngo.Client, reservationId, instanceId string, delay, timeout, minTimeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{deprovisioning},
		Target:     []string{provisionable, reprovisioned},
		Refresh:    hwReservationStateRefreshFunc(client, reservationId, instanceId),
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: minTimeout,
	}
	_, err := stateConf.WaitForState()
	return err
}

func getWaitForDeviceLock(deviceID string) *sync.WaitGroup {
	wgMutex.Lock()
	defer wgMutex.Unlock()
	wg, ok := wgMap[deviceID]
	if !ok {
		wg = &sync.WaitGroup{}
		wgMap[deviceID] = wg
	}
	return wg
}

func hwReservationStateRefreshFunc(client *packngo.Client, reservationId, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, _, err := client.HardwareReservations.Get(reservationId, &packngo.GetOptions{Includes: []string{"device"}})
		state := deprovisioning
		switch {
		case err != nil:
			err = errors.FriendlyError(err)
			state = errstate
		case r != nil && r.Provisionable:
			state = provisionable
		case r != nil && r.Device != nil && (r.Device.ID != "" && r.Device.ID != instanceId):
			log.Printf("[WARN] Equinix Metal device instance %s (reservation %s) was reprovisioned to a another instance (%s)", instanceId, reservationId, r.Device.ID)
			state = reprovisioned
		default:
			log.Printf("[DEBUG] Equinix Metal device instance %s (reservation %s) is still deprovisioning", instanceId, reservationId)
		}

		return r, state, err
	}
}

func getNetworkRank(family int, public bool) int {
	switch {
	case family == 4 && public:
		return 0
	case family == 6:
		return 1
	case family == 4 && public:
		return 2
	}
	return 3
}

func getNetworkInfo(ips []*packngo.IPAddressAssignment) NetworkInfo {
	ni := NetworkInfo{Networks: make([]map[string]interface{}, 0, 1)}
	for _, ip := range ips {
		network := map[string]interface{}{
			"address": ip.Address,
			"gateway": ip.Gateway,
			"family":  ip.AddressFamily,
			"cidr":    ip.CIDR,
			"public":  ip.Public,
		}
		ni.Networks = append(ni.Networks, network)

		// Initial device IPs are fixed and marked as "Management"
		if ip.Management {
			if ip.AddressFamily == 4 {
				if ip.Public {
					ni.Host = ip.Address
					ni.IPv4SubnetSize = ip.CIDR
					ni.PublicIPv4 = ip.Address
				} else {
					ni.PrivateIPv4 = ip.Address
				}
			} else {
				ni.PublicIPv6 = ip.Address
			}
		}
	}
	return ni
}
