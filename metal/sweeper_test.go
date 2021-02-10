package metal

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (interface{}, error) {
	token := os.Getenv("METAL_AUTH_TOKEN")

	if token == "" {
		token = os.Getenv("PACKET_AUTH_TOKEN")
	}

	if token == "" {
		return nil, fmt.Errorf("you must set METAL_AUTH_TOKEN")
	}

	config := Config{
		AuthToken: token,
	}

	return config.Client(), nil
}
