package vlan

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/packethost/packngo"
	"github.com/stretchr/testify/assert"
)

func TestIsVlanEqual(t *testing.T) {

	vlanA := uuid.NewString()
	vlanB := uuid.NewString()

	hrefBuilder := func(id string) string {
		return fmt.Sprintf("/metal/v1/virtual-networks/%s", id)
	}

	cases := []struct {
		name           string
		input          string
		network        packngo.VirtualNetwork
		expectedResult bool
	}{
		{
			"ID equals ID",
			vlanA,
			packngo.VirtualNetwork{
				ID: vlanA,
			},
			true,
		},
		{
			"ID does not equal ID",
			vlanB,
			packngo.VirtualNetwork{
				ID: vlanA,
			},
			false,
		},
		{
			"ID contained within Href",
			vlanA,
			packngo.VirtualNetwork{
				Href: hrefBuilder(vlanA),
			},
			true,
		},
		{
			"ID not contained within Href",
			vlanB,
			packngo.VirtualNetwork{
				Href: hrefBuilder(vlanA),
			},
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			result := isVlanEqual(c.input, c.network)
			assert.Equal(tt, c.expectedResult, result)
		})
	}

}
