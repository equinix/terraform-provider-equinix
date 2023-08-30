package equinix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	metalv1 "github.com/equinix-labs/metal-go/metal/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

const (
	deprovisioning = "deprovisioning"
	provisionable  = "provisionable"
	reprovisioned  = "reprovisioned"
	errstate       = "error"
)

var (
	wgMap   = map[string]*sync.WaitGroup{}
	wgMutex = sync.Mutex{}
)

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
	iacr.Reservations = convertStringArr(ia["reservation_ids"].([]interface{}))
	return iacr
}

func getNewIPAddressSlice(arr []interface{}) []packngo.IPAddressCreateRequest {
	addressTypesSlice := make([]packngo.IPAddressCreateRequest, len(arr))

	for i, m := range arr {
		addressTypesSlice[i] = ifToIPCreateRequest(m)
	}
	return addressTypesSlice
}

type NetworkInfo struct {
	Networks       []map[string]interface{}
	IPv4SubnetSize int
	Host           string
	PublicIPv4     string
	PublicIPv6     string
	PrivateIPv4    string
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

func getMetalGoNetworkType(device *metalv1.Device) (*string, error) {

	pgDevice := packngo.Device{}
	res, err := device.MarshalJSON()
	if err != nil {
		json.Unmarshal(res, pgDevice)
		networkType := pgDevice.GetNetworkType()
		return &networkType, nil
	}
	return nil, err
}

func getMetalGoNetworkInfo(ips []metalv1.IPAssignment) NetworkInfo {
	ni := NetworkInfo{Networks: make([]map[string]interface{}, 0, 1)}
	for _, ip := range ips {
		network := map[string]interface{}{
			"address": *ip.Address,
			"gateway": *ip.Gateway,
			"family":  *ip.AddressFamily,
			"cidr":    *ip.Cidr,
			"public":  *ip.Public,
		}
		ni.Networks = append(ni.Networks, network)

		// Initial device IPs are fixed and marked as "Management"
		if !ip.HasManagement() || *ip.Management {
			if *ip.AddressFamily == int32(4) {
				if !ip.HasPublic() || *ip.Public {
					ni.Host = *ip.Address
					ni.IPv4SubnetSize = int(*ip.Cidr)
					ni.PublicIPv4 = *ip.Address
				} else {
					ni.PrivateIPv4 = *ip.Address
				}
			} else {
				ni.PublicIPv6 = *ip.Address
			}
		}
	}
	return ni
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

func getMetalGoPorts(ps []metalv1.Port) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 1)
	for _, p := range ps {
		port := map[string]interface{}{
			"name":   p.GetName(),
			"id":     p.GetId(),
			"type":   p.GetType(),
			"mac":    p.Data.GetMac(),
			"bonded": p.Data.GetBonded(),
		}
		ret = append(ret, port)
	}
	return ret
}

func hwReservationStateRefreshFunc(client *packngo.Client, reservationId, instanceId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, _, err := client.HardwareReservations.Get(reservationId, &packngo.GetOptions{Includes: []string{"device"}})
		state := deprovisioning
		switch {
		case err != nil:
			err = friendlyError(err)
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

func waitUntilReservationProvisionable(ctx context.Context, client *packngo.Client, reservationId, instanceId string, delay, timeout, minTimeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{deprovisioning},
		Target:     []string{provisionable, reprovisioned},
		Refresh:    hwReservationStateRefreshFunc(client, reservationId, instanceId),
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: minTimeout,
	}
	_, err := stateConf.WaitForStateContext(ctx)
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

func waitForDeviceAttribute(ctx context.Context, d *schema.ResourceData, stateConf *retry.StateChangeConf) (string, error) {
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

	if stateConf == nil || stateConf.Refresh == nil {
		return "", errors.New("invalid stateconf to wait for")
	}

	attrValRaw, err := stateConf.WaitForStateContext(ctx)

	if v, ok := attrValRaw.(string); ok {
		return v, err
	}

	return "", err
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
					ValidateFunc: validation.StringMatch(uuidRE, "must be a valid UUID"),
				},
			},
		},
	}
}
