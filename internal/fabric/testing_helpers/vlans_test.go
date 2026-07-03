package testinghelpers

import (
	"testing"
)

func TestRandomVlanExcluding(t *testing.T) {
	t.Run("randomVlanExcluding returns only available VLAN", func(t *testing.T) {
		usedVlans := []int{}

		for i := 2; i <= 4092; i++ {
			if i == 3060 {
				continue
			}

			usedVlans = append(usedVlans, i)
		}

		out, err := randomVlanExcluding(usedVlans)

		if err != nil {
			t.Errorf(`errored with = %q`, err)
		}

		if out != 3060 {
			t.Errorf(`Output = %q, want 3060`, out)
		}
	})

	t.Run("randomVlanExcluding errors when no available VLAN", func(t *testing.T) {
		usedVlans := []int{}

		for i := 2; i <= 4092; i++ {
			usedVlans = append(usedVlans, i)
		}

		_, err := randomVlanExcluding(usedVlans)

		if err == nil {
			t.Errorf(`errored with = %q`, err)
		}
	})
}
