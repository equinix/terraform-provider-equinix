package equinix

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFabricCreateAzureConnection(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateEPLConnectionConfig(50),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", fmt.Sprint("fabric_tf_acc_test_CCEPL")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", fmt.Sprint("50")),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricCreateEPLConnectionConfig(100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", fmt.Sprint("fabric_tf_acc_test_CCEPL")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", fmt.Sprint("100")),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func checkConnectionDelete(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).fabricClient
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, testAccProvider.Meta().(*Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection" {
			continue
		}
		_, err := waitUntilConnectionDeprovisioned(rs.Primary.ID, client, ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}

func waitUntilConnectionDeprovisioned(uuid string, client *v4.APIClient, ctx context.Context) (v4.Connection, error) {
	log.Printf("Waiting for connection to be in deprovisioned, uuid %s", uuid)
	stateConf := &resource.StateChangeConf{
		Target: []string{"DEPROVISIONED"},
		Refresh: func() (interface{}, string, error) {
			dbConn, _, err := client.ConnectionsApi.GetConnectionByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", err
			}
			updatableState := ""
			if "DEPROVISIONED" == *dbConn.State {
				updatableState = string(*dbConn.State)
			}
			return dbConn, updatableState, nil
		},
		Timeout:    3 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.Connection{}

	if err == nil {
		dbConn = inter.(v4.Connection)
	}
	return dbConn, err
}

func testAccFabricCreateEPLConnectionConfig(bandwidth int32) string {
	return fmt.Sprintf(`resource "equinix_fabric_connection" "test" {
		type = "EVPL_VC"
		name = "fabric_tf_acc_test_CCEPL"
		notifications{
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		order {
		purchase_order_number = "1-129105284100"
			}
		bandwidth = %d
		a_side {
		access_point {
			type = "COLO"
				port {
				uuid = "c4d9350e-77cd-7cdd-1ce0-306a5c00a600"
				}
				link_protocol {
					type= "QINQ"
					vlan_s_tag= 1231
				}
			}
		}
		z_side {
		access_point {
			type = "COLO"
				port{
				uuid = "c4d9350e-783c-83cd-1ce0-306a5c00a600"
				}
				link_protocol {
					type= "QINQ"
					vlan_s_tag= 457
				}
			location {
        		metro_code= "SV"
      		}
			}
		}
	}`, bandwidth)
}

func testAccFabricCreateAzureConnectionConfig(bandwidth int32) string {
	return fmt.Sprintf(`resource "equinix_fabric_connection" "test" {
	name = "fabric_tf_acc_CSAZURE"
	description = "Test Connection"
	type = "EVPL_VC"
	notifications{
		type="ALL" 
		emails=["example@equinix.com"]
	} 
	bandwidth = %d
	redundancy {priority= "PRIMARY"}
	order {
    	purchase_order_number= "1-323292"
  	}
  	a_side {
    	access_point {
      		type= "COLO"
      		port {
        		uuid= "c4d9350e-783c-83cd-1ce0-306a5c00a600"
      		}
      	link_protocol {
        	type= "QINQ"
        	vlan_s_tag= "2019"
      	}
    }
  }
  	z_side {
    	access_point {
      		type= "SP"
			authentication_key= "a38565b9-5d32-45ba-bb01-0649e2735753"
      		profile {
        		type= "L2_PROFILE"
        		uuid= "bfb74121-7e2c-4f74-99b3-69cdafb03b41"
      		}
      		location {
        		metro_code= "SV"
      		}
    	}
  	}
}
`, bandwidth)
}

func testAccFabricUpdateConnectionConfig(bandwidth int32) string {
	return fmt.Sprintf(`resource "equinix_fabric_connection" "test" {
	uuid = equinix_fabric_connection.test.uuid
	name = "fabric_tf_acc_CSAZURE"
	description = "Test Connection"
	type = "EVPL_VC"
	notifications{
		type="ALL" 
		emails=["example@equinix.com"]
	} 
	bandwidth = %d
	redundancy {priority= "PRIMARY"}
	order {
    	purchase_order_number= "1-323292"
  	}
  	a_side {
    	access_point {
      		type= "COLO"
      		port {
        		uuid= "c4d9350e-783c-83cd-1ce0-306a5c00a600"
      		}
      	link_protocol {
        	type= "QINQ"
        	vlan_s_tag= "2019"
      	}
    }
  }
  	z_side {
    	access_point {
      		type= "SP"
			authentication_key= "a38565b9-5d32-45ba-bb01-0649e2735753"
      		profile {
        		type= "L2_PROFILE"
        		uuid= "bfb74121-7e2c-4f74-99b3-69cdafb03b41"
      		}
      		location {
        		metro_code= "SV"
      		}
    	}
  	}
}
`, bandwidth)
}

func TestAccFabricReadConnection(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadConnectionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", fmt.Sprint("fabric_tf_acc_test")),
				),
			},
			{
				ResourceName:      "equinix_fabric_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFabricReadConnectionConfig() string {
	return fmt.Sprint(`data "equinix_fabric_connection" "test" {
	uuid = "3e91216d-526a-45d2-9029-0c8c8ba48b60"
	}`)
}
