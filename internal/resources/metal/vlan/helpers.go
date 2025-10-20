package vlan

import (
	"strings"

	"github.com/packethost/packngo"
)

// isVlanEqual compares an ID to a VLAN object, returning if the ID
// matches. This is handy when you might not have the full object, and
// only have the light reference to the object, meaning only the
// `href` key is populated.
func isVlanEqual(id string, vlan packngo.VirtualNetwork) bool {
	if vlan.ID != "" {
		return vlan.ID == id
	}

	return strings.Contains(vlan.Href, id)
}
