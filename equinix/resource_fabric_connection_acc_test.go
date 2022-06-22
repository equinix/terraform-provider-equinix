package equinix

import (
	"context"
	"fmt"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"testing"
	"time"
)

func TestAccFabricCreateAzureConnection(t *testing.T) {
	log.Printf(" inside fabric create connection test ")
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
	conn := v4.Connection{}
	var err error
	counter := 0
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection" {
			continue
		}

		//TODO Use retry terraform helpers - Move to delete Func
		for conn.State == nil || "DEPROVISIONED" != *conn.State {
			time.Sleep(30 * time.Second)
			conn, _, err = client.ConnectionsApi.GetConnectionByUuid(ctx, rs.Primary.ID, nil)
			if err != nil {
				return fmt.Errorf("API call failed for the resource ")
			}
			if counter >= 4 {
				break
			}
			counter++
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
				uuid = "c4d9350e-77cd-7cdd-1ce0-306a5c00a600"
				}
				link_protocol {
					type= "QINQ"
					vlan_s_tag= 2345
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
					vlan_s_tag= 1927
				}
			location {
        		metro_code= "SV"
      		}
			}
		}
	}`, bandwidth)
}

func testAccFabricCreateAzureConnectionConfig(bandwidth int32) string {
	log.Printf(" inside create connection test config ")
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
	log.Printf(" inside create connection test config ")
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

/*func testAccFabricCreateConnectionConfig(conectionType string) string {
	fmt.Println(" inside create connection test config ")
	var config string

	if conectionType == "COLO2COLO_EPL" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
		type = "EVPL_VC"
		name = "fabric_tf_acc_test_CCEPL"
		notifications{
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		order {
		purchase_order_number = "1-129105284100"
			}
		bandwidth = 100
		a_side {
		access_point {
			type = "COLO"
				port {
				uuid = "c4d9350e-77cd-7cdd-1ce0-306a5c00a600"
				}
				link_protocol {
					type= "QINQ"
					vlan_s_tag= 1234
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
					vlan_s_tag= "1001"
				}
			location {
        		metro_code= "SV"
      		}
			}
		}
	}`)
	} else if conectionType == "COLO2COLO_ACCESS_EPL" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "ACCESS_EPL_VC",
			name= "fabric_tf_acc_test_CCAEPL",
			order{
				purchase_order_number= "1-129105284100"
			},
			bandwidth= "100",
			a_side{
			access_point {
				type= "COLO",
					port{
					uuid= "c4d9350e-77cd-7cdd-1ce0-306a5c00a600"
				},
				link_protocol {
					type= "QINQ",
					vlan_s_tag= 1234
					}
				}
			},
			z_side {
			access_point {
				type = "COLO",
					port {
					uuid= "20d32a80-0d61-4333-bc03-707b591ae2f4"
				}
			}
		},
			notifications{
			type= "ALL",
			emails= ["test@equinix.com","test1@equinix.com"]
			}
	}`)
	} else if conectionType == "COLO2COLO_DD" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "EVPL_VC",
			name= "fabric_tf_acc_test_CCDD",
			bandwidth= "1000",
			redundancy {
			priority= "PRIMARY"
		},
			a_side{
			access_point{
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag"= 1001
				}
			}
		},
			order {
			"purchase_order_number= "po1234"
		},
			z_side {
			access_point {
				type= "COLO",
					port {
					uuid= "b067f685-49b0-1a09-6fe0-360a5d00afdg"
				},
				link_protocol {
					type= "DOT1Q",
						vlan_tag= 1001
				}
			}
		},
			notifications{
				"type": "ALL",
				"emails": ["test@test.com"]
			}
	}`)
	} else if conectionType == "COLO2COLO_DQ" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type="EVPL_VC",
			name= "fabric_tf_acc_test_CCDQ",
			bandwidth= "1000",
			redundancy {
				priority= "PRIMARY"
			},
			a_side {
			access_point {
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
  					vlan_tag= "1001"
				}
			}
		},
			order{
			purchase_order_number= "po1234"
		},
			z_side {
			access_point {
				type= "COLO",
					port {
					uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				link_protocol {
					type= "QINQ",
					vlan_s_tag= "1001",
					vlan_c_tag= "1002"
				}
			}
		},
			notifications {
				type= "ALL",
				emails= ["test@test.com"]
			}
	}`)
	} else if conectionType == "COLO2COLO_QQ" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "EVPL_VC",
			name= "fabric_tf_acc_test_CCQQ",
			bandwidth= 1000,
			redundancy {
			priority= "PRIMARY"
			},
			a_side {
			access_point {
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "QINQ",
					vlan_s_tag= "1001"
				}
			}
		},
			order {
			purchase_order_number= "po1234"
		},
			z_side {
			access_point {
				type= "COLO",
					port{
					uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				link_protocol {
					type= "QINQ",
					vlan_s_tag= "1001"
				}
			}
		},
			notifications{
				type= "ALL",
				emails= ["test@test.com"]
		}
	}`)
	} else if conectionType == "COLO2COLO_QD" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {"type": "EVPL_VC",
			"name": "fabric_tf_acc_test_CCQD,
			"bandwidth": 1000,
			"redundancy": {
			"priority": "PRIMARY"
		},
			a_side {
			access_point {
				type = "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "QINQ",
					vlan_s_tag= "1001",
					vlan_c_tag= "1125"
				}
			}
		},
			order {
			purchase_order_number= "po1234"
		},
			z_side {
			access_point {
				type= "COLO",
					port {
					uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag= "1001"
				}
			}
		},
			notifications{
				type= "ALL",
				emails= ["test@test.com"]
			}
	}`)
	} else if conectionType == "COLO2COLO_RESELLER" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "EVPL_VC",
			name= "fabric_tf_acc_test_RESELLER",
			order {
				purchase_order_number= "1-129105284100"
			},
			bandwidth= "100",
			a_side {
			access_point {
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag= "1001"
				}
			}
		},
			z_side {
			access_point {
				type= "COLO",
					port {
					uuid= "20d32a80-0d61-4333-bc03-707b591ae2f4"
				},
				link_protocol {
					type= "QINQ",
					vlan_s_tag= "1002",
					vlan_c__tag= "1001"
				}
			}
		},
			notifications{
				type= "ALL",
				emails= ["test@equinix.com","test1@equinix.com"]
			}
	}`)
	} else if conectionType == "COLO2SP_AZURE" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
	name = "fabric_tf_acc_CSAZURE"
	description = "Test Connection"
	type = "EVPL_VC"
	notifications{
		type="ALL"
		emails=["example@equinix.com"]
	}
	bandwidth = 50
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
`)
	} else if conectionType == "COLO2SP_AWS" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "EVPL_VC",
			name= "fabric_tf_acc_test_CSAWS",
			bandwidth= "1000",
			redundancy {
			priority= "PRIMARY"
		},
			a_side {
			access_point {
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag= "1001"
				}
			}
		},
			order {
			purchase_order_number= "po1234"
		},
			z_side {
			access_point {
				type= "SP",
					profile {
						type= "L2_PROFILE",
						uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				location {
					metro_code= "DC"
				},
				seller_region= "us-east-1",
				authentication_key= "357848976964"
			}
		},
			notifications={
				type= "ALL",
				emails= ["test@test.com"]
			}
	}`)
	} else if conectionType == "COLO2SP_GOOG" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "EVPL_VC",
			name= "fabric_tf_acc_test_CSGOO",
			bandwidth= 1000,
			redundancy {
				priority= "PRIMARY"
			},
			a_side {
			access_point {
				type= "COLO",
					port{
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag= "1001"
					}
				}
			},
			order{
				purchase_order_number= "po1234"
			},
			z_side {
			access_point {
				type= "SP",
					profile {
					type= L2_PROFILE,
					uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				location {
					metro_code= "DC"
				},
				authentication_key= "6b5596e3-ee7f-4a74-b5ff-2ac28e2b961d/us-west2/2"
			}
		},
			notifications{
			type= "ALL",
			emails= ["test@test.com"]
		}
	}`)
	} else if conectionType == "COLO2SP_GEN" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "EVPL_VC",
			name= "fabric_tf_acc_test_CSGEN",
			bandwidth= "1000",
			redundancy {
				priority= "PRIMARY"
			},
			a_side {
			access_point {
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag= 1001
				}
			}
		},
			order{
			purchase_order_number= "po1234"
		},
			z_side {
			access_point {
				type= "SP",
					profile {
					type= "L2_PROFILE",
						uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				location {
					metro_code= "DC"
				}
			}
		},
			notifications{
		type= "ALL",
		emails= ["test@test.com"]
		}
	}`)
	} else if conectionType == "COLO2SP_IBM" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type="EVPL_VC",
			name= "fabric_tf_acc_test_CSIBM",
			bandwidth= 1000,
			redundancy {
			priority= "PRIMARY"
		},
			a_side {
			access_point {
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag= 1001
				}
			}
		},
			order {
				purchase_order_number= "po1234"
			},
			z_side {
			access_point {
				type= "SP",
					profile {
					type="L2_PROFILE",
						uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				location {
					metro_code="DC"
				},
				authentication_key="5bf92b31d921499f963592cd816f6be7",
					seller_region= "San Jose 2"
			}
		},
			notifications={
				"type": "ALL",
				"emails": ["test@test.com"]
			}
		additional_info= [
		{
		key= "ASN",
		value= "1234"
		},
		{
		key= "Global",
		value= false
		},
		{
		key= "BGP_IBM_CIDR",
		value= "172.16.0.18/30"
		},
		{
		key= "BGP_CER_CIDR",
		value= "172.16.0.19/30"
		}
	]
	}`)
	} else if conectionType == "COLO2SP_ORA" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type="EVPL_VC",
			name= "fabric_tf_acc_test_CSORA",
			bandwidth= "1000",
			redundancy {
			priority= "PRIMARY"
		},
			a_side {
			access_point {
				type= "COLO",
					port= {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
						vlan_tag= 1001
				}
			}
		},
			order= {
			purchase_order_number= "po1234"
		},
			z_side {
			access_point {
				type= "SP",
					profile {
					type= "L2_PROFILE",
						uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				location {
					metro_code= "DC"
				},
				seller_region= "us-ashburn-1",
					authentication_key= "ocid1.virtualcircuit.oc1.iad.aaaaaaaanqtbkv4nvi6elbx2d4tprmtzgodmj6imm46j5ltmpootwz2vgcda"
			}
		},
			notifications{
			"type": "ALL",
			"emails": ["test@test.com"]
		}
	}`)
	} else if conectionType == "COLO2SP_ALIBABA" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type="EVPL_VC",
			name= "fabric_tf_acc_test_CSALI",
			bandwidth= 1000,
			redundancy {
			priority= "PRIMARY"
		},
			a_side {
			access_point {
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol= {
					type= "DOT1Q",
					vlan_tag= 1001
				}
			}
		},
			order{
			purchase_order_number= "po1234"
		},
			z_side {
			access_point {
				type= "SP",
					profile {
					type= "L2_PROFILE",
						uuid= "22d4e853-ef33-4ff0-b5b2-a2b1d5dfa50c"
				},
				location{
					"metro_code": "SV"
				},
				seller_region= "San Jose 2",
				authentication_key= 1956030
			}
		},
			notifications{
				"type": "ALL",
				"emails": ["test@test.com"]
			}
	}`)
	} else if conectionType == "COLO2ST" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type="EVPL_VC",
			name= "fabric_tf_acc_test_CST",
			order {
				purchase_order_number": "1-129105284100"
			},
			bandwidth= 100,
			redundancy {
				group= "m167f685-41b0-1b07-6de0-320a5c00abeu",
				priority= "SECONDARY"
			},
			a_side {
			access_point {
				type= "COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag= "1001"
				}
			}
		},
			z_side {
			service_token {
				uuid= "20d32a80-0d61-4333-bc03-707b591ae2f5"
			}
		},
			notifications{
				type= "ALL",
				emails= ["test@equinix.com","test1@equinix.com"]
		}
	}`)
	} else if conectionType == "COLO2GW" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "VC_GW",
			name= "fabric_tf_acc_test_CGW",
			order {
			purchase_order_number= "1-129105284100"
		}
			bandwidth= "100",
			redundancy{
			group= "m167f685-41b0-1b07-6de0-320a5c00abeu",
			priority= "SECONDARY"
		},
			a_side {
			access_point"{
				type="COLO",
					port {
					uuid= "a867f685-41b0-1b07-6de0-320a5c00abdd"
				},
				link_protocol{
					type= "DOT1Q",
					vlan_tag= "1001"
				}
			}
		},
			z_side {
			service_token{
				uuid= "20d32a80-0d61-4333-bc03-707b591ae2f5"
			}
		},
			notifications={
				type": "ALL",
				emails": ["test@equinix.com","test1@equinix.com"]
		}
	}`)
	} else if conectionType == "VD2SP" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type= "EVPL_VC",
			bandwidth= "50",
			name="fabric_tf_acc_test_VS",
			a_side {
			access_point {
				type ="VD",
					virtual_device {
					type="EDGE",
						uuid="20d32a80-0d61-4333-bc03-707b591ae2f4"
				},
				interface {
					type= "NETWORK",
					id= 45645
				}
			}
		},
			z_side {
			access_point {
				type= "SP",
					profile {
					type= "L2_PROFILE",
					uuid= "95542b34-cf1c-41aa-89f7-590946f9df53"
				},
				authentication_key= "9ec79dca-3729-4cbb-b372-a4e72b615e12",
				seller_region= "us-west-1"
			}
		},
			order {
				purchase_order_number= "1-323292"
		},
			notifications{
			type= "ALL",
			emails= ["test@equinix.com","test1@equinix.com"]
		}
	}`)

	} else if conectionType == "VD2COLO" {
		config = fmt.Sprint(`resource "equinix_fabric_connection" "test" {
			type="EVPL_VC",
			bandwidth= "1000",
			name= "fabric_tf_acc_test_VCOLO",
			a_side {
			access_point {
				type= "VD",
				virtual_device {
					type= "EDGE",
					uuid= "20d32a80-0d61-4333-bc03-707b591ae2f4"
				},
				interface {
					type= "NETWORK",
					id= 45645
				}
			}
		},
			z_side {
			access_point {
				type="COLO",
					port {
					uuid= "20d32a80-0d61-4333-bc03-707b591ae2f4"
				},
				link_protocol {
					type= "DOT1Q",
					vlan_tag= "300"
				}
			}
		},
			order {
				purchase_order_number= "1-129105284100"
			},
			notifications{
			type="ALL",
			emails= ["test@equinix.com","test1@equinix.com"]
		}
	}
	}`)
	}
	return config
}*/

func TestAccFabricReadConnection(t *testing.T) {
	log.Printf(" inside fabric read connection test ")
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
