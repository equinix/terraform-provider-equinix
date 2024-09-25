package device

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/converters"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

type NetworkInfo struct {
	Networks       []map[string]interface{}
	IPv4SubnetSize int
	Host           string
	PublicIPv4     string
	PublicIPv6     string
	PrivateIPv4    string
}

func getNetworkInfo(ips []metalv1.IPAssignment) NetworkInfo {
	ni := NetworkInfo{Networks: make([]map[string]interface{}, 0, 1)}
	for _, ip := range ips {
		network := map[string]interface{}{
			"address": ip.GetAddress(),
			"gateway": ip.GetGateway(),
			"family":  ip.GetAddressFamily(),
			"cidr":    ip.GetCidr(),
			"public":  ip.GetPublic(),
		}
		ni.Networks = append(ni.Networks, network)

		// Initial device IPs are fixed and marked as "Management"
		if ip.GetManagement() {
			if ip.GetAddressFamily() == 4 {
				if ip.GetPublic() {
					ni.Host = ip.GetAddress()
					ni.IPv4SubnetSize = int(ip.GetCidr())
					ni.PublicIPv4 = ip.GetAddress()
				} else {
					ni.PrivateIPv4 = ip.GetAddress()
				}
			} else {
				ni.PublicIPv6 = ip.GetAddress()
			}
		}
	}
	return ni
}

// Deprecated: this exists for backwards compatibility with
// packngo-based resources.  It relies on the deprecated packngo
// SDK and either the logic from packngo should be pulled in to
// the provider or this functionality should be removed entirely
func getNetworkType(device *metalv1.Device) (*string, error) {
	pgDevice := packngo.Device{}
	res, err := device.MarshalJSON()
	if err == nil {
		if err = json.Unmarshal(res, &pgDevice); err == nil {
			networkType := pgDevice.GetNetworkType()
			return &networkType, nil
		}
	}
	return nil, err
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

func getPorts(ps []metalv1.Port) []map[string]interface{} {
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

func hwReservationStateRefreshFunc(ctx context.Context, client *metalv1.APIClient, reservationId, instanceId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, _, err := client.HardwareReservationsApi.FindHardwareReservationById(ctx, reservationId).Include([]string{"device"}).Execute()
		state := deprovisioning
		switch {
		case err != nil:
			err = equinix_errors.FriendlyError(err)
			state = errstate
		case r != nil && r.GetProvisionable():
			state = provisionable
		case r != nil && r.Device != nil && (r.Device.GetId() != "" && r.Device.GetId() != instanceId):
			log.Printf("[WARN] Equinix Metal device instance %s (reservation %s) was reprovisioned to a another instance (%s)", instanceId, reservationId, r.Device.GetId())
			state = reprovisioned
		default:
			log.Printf("[DEBUG] Equinix Metal device instance %s (reservation %s) is still deprovisioning", instanceId, reservationId)
		}

		return r, state, err
	}
}

func WaitUntilReservationProvisionable(ctx context.Context, client *metalv1.APIClient, reservationId, instanceId string, delay, timeout, minTimeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{deprovisioning},
		Target:     []string{provisionable, reprovisioned},
		Refresh:    hwReservationStateRefreshFunc(ctx, client, reservationId, instanceId),
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

func getDeviceMap(device metalv1.Device) map[string]interface{} {
	networkInfo := getNetworkInfo(device.IpAddresses)
	sort.SliceStable(networkInfo.Networks, func(i, j int) bool {
		famI := int(networkInfo.Networks[i]["family"].(int32))
		famJ := int(networkInfo.Networks[j]["family"].(int32))
		pubI := networkInfo.Networks[i]["public"].(bool)
		pubJ := networkInfo.Networks[j]["public"].(bool)
		return getNetworkRank(famI, pubI) < getNetworkRank(famJ, pubJ)
	})
	keyIDs := []string{}
	for _, k := range device.SshKeys {
		keyIDs = append(keyIDs, path.Base(k.GetHref()))
	}
	ports := getPorts(device.NetworkPorts)

	return map[string]interface{}{
		"hostname":            device.GetHostname(),
		"project_id":          device.Project.GetId(),
		"description":         device.GetDescription(),
		"device_id":           device.GetId(),
		"facility":            device.Facility.GetCode(),
		"metro":               device.Metro.GetCode(),
		"plan":                device.Plan.GetSlug(),
		"operating_system":    device.OperatingSystem.GetSlug(),
		"state":               device.GetState(),
		"billing_cycle":       device.GetBillingCycle(),
		"ipxe_script_url":     device.GetIpxeScriptUrl(),
		"always_pxe":          device.GetAlwaysPxe(),
		"root_password":       device.GetRootPassword(),
		"tags":                converters.StringArrToIfArr(device.GetTags()),
		"access_public_ipv6":  networkInfo.PublicIPv6,
		"access_public_ipv4":  networkInfo.PublicIPv4,
		"access_private_ipv4": networkInfo.PrivateIPv4,
		"network":             networkInfo.Networks,
		"ssh_key_ids":         keyIDs,
		"ports":               ports,
		"sos_hostname":        device.GetSos(),
	}
}
