package vlan

import (
	"fmt"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
)

func MatchingVlan(vlans []metalv1.VirtualNetwork, vxlan int, projectID, facility, metro string) (*metalv1.VirtualNetwork, error) {
	matches := []metalv1.VirtualNetwork{}
	for _, v := range vlans {
		if vxlan != 0 && int(v.GetVxlan()) != vxlan {
			continue
		}
		/*if facility != "" && v.FacilityCode != facility {
			continue
		}*/
		if metro != "" && v.GetMetroCode() != metro {
			continue
		}
		matches = append(matches, v)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("Project %s has more than one matching VLAN", projectID)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("Project %s does not have matching VLANs", projectID)
	}
	return &matches[0], nil
}
