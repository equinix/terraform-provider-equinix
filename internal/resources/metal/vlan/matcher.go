package vlan

import (
	"fmt"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"
	"github.com/packethost/packngo"
)

func MatchingVlan(vlans []packngo.VirtualNetwork, vxlan int, projectID, facility, metro string) (*packngo.VirtualNetwork, error) {
	matches := []packngo.VirtualNetwork{}
	for _, v := range vlans {
		if vxlan != 0 && v.VXLAN != vxlan {
			continue
		}
		if facility != "" && v.FacilityCode != facility {
			continue
		}
		if metro != "" && v.MetroCode != metro {
			continue
		}
		matches = append(matches, v)
	}
	if len(matches) > 1 {
		return nil, equinix_errors.FriendlyError(fmt.Errorf("Project %s has more than one matching VLAN", projectID))
	}

	if len(matches) == 0 {
		return nil, equinix_errors.FriendlyError(fmt.Errorf("Project %s does not have matching VLANs", projectID))
	}
	return &matches[0], nil
}
