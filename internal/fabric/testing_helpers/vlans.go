package testinghelpers

import (
	"context"
	"errors"
	"log"

	"math/rand"

	"github.com/equinix/terraform-provider-equinix/internal/sweep"
)

func RandomVlan(portUUID string) (int, error) {
	usedVlans, err := getUsedVlans(portUUID)
	if err != nil {
		return 0, err
	}

	targetVlan, err := randomVlanExcluding(usedVlans)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] Fetched VLAN: %d", targetVlan)

	return targetVlan, nil
}

func getUsedVlans(portUUID string) ([]int, error) {
	log.Printf("[DEBUG] Fetching vlans")
	ctx := context.Background()
	meta, err := sweep.GetConfigForFabric()
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
