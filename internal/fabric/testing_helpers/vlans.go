package testinghelpers

import (
	"context"
	"errors"
	"os"
	"testing"

	"math/rand"
)

func RandomVlan(t testing.TB, portUUID string) (int, error) {
	if os.Getenv("TF_ACC") == "" {
		t.Log("Skip fetching VLAN because TF_ACC is not set")
		return 0, nil
	}

	usedVlans, err := getUsedVlans(t, portUUID)
	if err != nil {
		return 0, err
	}

	targetVlan, err := randomVlanExcluding(usedVlans)
	if err != nil {
		return 0, err
	}
	t.Logf("[DEBUG] Fetched VLAN: %d", targetVlan)

	return targetVlan, nil
}

func getUsedVlans(t testing.TB, portUUID string) ([]int, error) {
	t.Log("[DEBUG] Fetching vlans")
	ctx := context.Background()
	meta, err := GetConfigForFabric()
	if err != nil {
		return nil, err
	}

	if err := meta.Load(ctx); err != nil {
		return nil, err
	}

	fabric := meta.NewFabricClientForTesting(ctx)

	vlans, _, err := fabric.PortsApi.GetVlans(ctx, portUUID).Execute()
	if err != nil {
		return nil, err
	}

	usedVlans := make([]int, len(vlans.Data))
	for i, v := range vlans.Data {
		usedVlans[i] = int(v.GetVlanTag())
	}

	return usedVlans, nil
}

func randomVlanExcluding(excluded []int) (int, error) {
	const minVlan = 2
	const maxVlan = 4092

	excludedSet := make(map[int]any)
	for _, n := range excluded {
		excludedSet[n] = struct{}{}
	}

	var valid []int
	for i := minVlan; i <= maxVlan; i++ {
		if _, found := excludedSet[i]; !found {
			valid = append(valid, i)
		}
	}

	if len(valid) == 0 {
		return 0, errors.New("no available numbers")
	}

	return valid[rand.Intn(len(valid))], nil
}
