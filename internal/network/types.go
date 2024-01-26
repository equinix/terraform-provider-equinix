package network

import "strings"

var (
	DeviceNetworkTypes   = []string{"layer3", "hybrid", "layer2-individual", "layer2-bonded"}
	DeviceNetworkTypesHB = []string{"layer3", "hybrid", "hybrid-bonded", "layer2-individual", "layer2-bonded"}
	NetworkTypeList      = strings.Join(DeviceNetworkTypes, ", ")
	NetworkTypeListHB    = strings.Join(DeviceNetworkTypesHB, ", ")
)
