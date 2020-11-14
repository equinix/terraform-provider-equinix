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

	if os.Getenv("PACKET_AUTH_TOKEN") == "" {
		return nil, fmt.Errorf("you must set PACKET_AUTH_TOKEN")
	}

	config := Config{
		AuthToken: os.Getenv("PACKET_AUTH_TOKEN"),
	}

	return config.Client(), nil
}
