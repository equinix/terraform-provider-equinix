package sweep

import (
	"strings"
)

var (
	FabricTestResourceSuffixes = []string{"_PFCR", "_PNFV", "_PPDS"}
)

func IsSweepableFabricTestResource(resourceName string) bool {
	for _, suffix := range FabricTestResourceSuffixes {
		if strings.HasSuffix(resourceName, suffix) {
			return true
		}
	}
	return false
}
