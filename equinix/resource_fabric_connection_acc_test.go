package equinix

import (
	"context"
	"fmt"
	"testing"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFabricCreateConnection(t *testing.T) {
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

func TestAccFabricCreateFGConnection(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateFG2portConnectionConfig("fabric_tf_acc_FG2port1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", fmt.Sprint("fabric_tf_acc_FG2port1")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", fmt.Sprint("100")),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricCreateFG2portConnectionConfig("fabric_tf_acc_FG2port2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", fmt.Sprint("fabric_tf_acc_test_FG2port2")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", fmt.Sprint("100")),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func checkConnectionDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, testAccProvider.Meta().(*Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection" {
			continue
		}
		err := waitUntilConnectionDeprovisioned(rs.Primary.ID, testAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
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
				uuid = "eb92632a-3747-7478-b5e0-306a5c00aecd"
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= 2397
				}
			}
		}
		z_side {
		access_point {
			type = "COLO"
				port{
				uuid = "3d7c1d97-2833-46fd-b1b1-ca619263eeb9"
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= 2398
				}
			location {
        		metro_code= "CH"
      		}
			}
		}
	}`, bandwidth)
}

func testAccFabricCreateAzureConnectionConfig(bandwidth int32) string {
	return fmt.Sprintf(`resource "equinix_fabric_connection" "test" {
	name = "fabric_tf_acc_CSAZURE"
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
        	vlan_s_tag= "2231"
      	}
    }
  }
  	z_side {
    	access_point {
      		type= "SP"
			authentication_key= "7244f849-8665-493e-8877-a4a0abb2a07e"
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

func testAccFabricCreateGenericConfig(bandwidth int32) string {
	return fmt.Sprintf(`resource "equinix_fabric_connection" "test" {
	name = "fabric_tf_acc_Generic"
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
      		profile {
        		type= "L2_PROFILE"
        		uuid= "7a278326-cfd3-46a6-92d0-e10ed0d7af50"
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

func testAccFabricCreateFG2portConnectionConfig(name string) string {
	return fmt.Sprintf(`resource "equinix_fabric_connection" "test" {
		type = "IP_VC"
		name = "%s"
		notifications{
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		order {
		purchase_order_number = "1-129105284100"
			}
		bandwidth = 100
		redundancy {
			priority= "PRIMARY"
		}
		a_side {
		access_point {
			type = "GW"
      		gateway {
        		uuid = "4f543d31-88f7-4eaf-b378-6b6a08e31e94"
      		}
			}
		}
		project{
		   project_id = "776847000642406"
		}
		z_side {
		access_point {
			type = "COLO"
				port{
					uuid = "3d7c1d97-2833-46fd-b1b1-ca619263eeb9"
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= 2325
				}
			}
		}
	}`, name)
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
						"data.equinix_fabric_connection.test", "name", fmt.Sprint("fabric_tf_acc_test")),
				),
			},
		},
	})
}

func testAccFabricReadConnectionConfig() string {
	return fmt.Sprint(`data "equinix_fabric_connection" "test" {
	uuid = "3e91216d-526a-45d2-9029-0c8c8ba48b60"
	}`)
}
