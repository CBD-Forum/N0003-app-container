// Copyright [2016] [Cuiting Shi ]
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// 
package main

import (
	"fmt"
	"errors"
	"cbdforum/app-container/chaincode/src/common"
	"cbdforum/app-container/chaincode/src/order"
	"cbdforum/app-container/chaincode/src/resource"
	auth "cbdforum/app-container/chaincode/src/user"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type ContainerChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	shim.SetLoggingLevel(shim.LogDebug)
	err := shim.Start(new(ContainerChaincode))
	if err != nil {
		fmt.Printf("Error starting Container chaincode - %s", err)
	}
}

// ============================================================================================================================
// Init - initialize the chaincode
// ============================================================================================================================
func (t *ContainerChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Container Is Starting Up")
	//_, args := stub.GetFunctionAndParameters()
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// message server address
	messageServer := args[0]
	idList := []string{"0"}
	err = stub.PutState(common.MessageServer, []byte(messageServer))
	if err != nil {
		return nil, err
	}

	// store compaitible container application version
	err = stub.PutState(common.ContainerUIVersion, []byte("1.0.0"))
	if err != nil {
		return nil, err
	}

	// store the uppper value for container id, range is [0, containerIdMax)
	indexNames := []string {
		common.INDEX_USER,
		common.INDEX_CONTAINER,
		common.INDEX_VEHICLE,
		common.INDEX_TRANSPORT_TASK,
		common.INDEX_SHIPPING_SCHEDULE,
		common.INDEX_ORDER,
	}
	for _, indexName := range indexNames {
		err = common.PutStringList(stub, indexName, idList)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error trying to put index %s",indexName))
		}

	}

	fmt.Println(" - ready for action")
	return nil, nil
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *ContainerChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "insert_container" { //deletes a container from state
		return resource.InsertContainer(stub, args)
	} else if function == "update_container" { //create a new container
		return resource.UpdateContainer(stub, args)
	} else if function == "delete_container" {
		return resource.DeleteContainer(stub, args)

	} else if function == "insert_vehicle" {
		return resource.InsertVehicle(stub, args)
	} else if function == "update_vehicle" {
		return resource.UpdateVehicle(stub, args)
	} else if function == "delete_vehicle" {
		return resource.DeleteVehicle(stub, args)

	} else if function == "insert_shippingschedule" {
		return resource.InsertShippingSchedule(stub, args)
	} else if function == "update_shippingschedule" {
		return resource.UpdateShippingSchedule(stub, args)
	} else if function == "delete_shippingschedule" {
		return resource.DeleteShippingSchedule(stub, args)

	} else if function == "delete_transporttask" {
		return resource.DeleteTransportTask(stub, args)

	} else if function == "delete_order" {
		return order.DeleteOrder(stub, args)

	} else if function == "create_order" {
		return order.CreateOrder(stub, args)
	} else if function == "check_order" {
		return order.CheckOrder(stub, args)
	} else if function == "book_space" {
		return order.BookSpace(stub, args)
	} else if function == "book_vehicle" {
		return order.BookVehicle(stub, args)
	} else if function == "fetch_empty_containers" {
		return order.FetchEmptyContainer(stub, args)
	} else if function == "pack_goods" {
		return order.PackGoods(stub, args)
	} else if function == "arrive_yard" {
		return order.ArriveYard(stub, args)
	} else if function == "load_goods" {
		return order.LoadGoods(stub, args)
	} else if function == "departure" {
		return order.Departure(stub, args)
	} else if function == "arrive_destination_port" {
		return order.ArriveDestinationPort(stub, args)
	} else if function == "deliver_goods" {
		return order.DeliverGoods(stub, args)
	} else if function == "confirm_receipt" {
		return order.ConfirmReceipt(stub, args)
	} else if function == "finish_order" {
		return order.FinishOrder(stub, args)

	} else if function == "register_user" {
		//create a new container owner
		return auth.RegisterUser(stub, args)
	}
	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return nil, errors.New("Received unknown invoke function name - '" + function + "'")
}

// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *ContainerChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Handle different functions
	if function == "read_containers_by_ownerid" {
		return resource.ReadAllContainersByOwnerId(stub, args)
	} else if function == "read_container" {
		return resource.ReadContainer(stub, args)

	} else if function == "read_vehicles_by_ownerid" {
		return resource.ReadAllVehiclesByOwnerId(stub, args)
	} else if function == "read_vehicle" {
		return resource.ReadVehicle(stub, args)

	} else if function == "read_shippingschedules_by_ownerid" {
		return resource.ReadAllShippingSchedulesByOwnerId(stub, args)
	} else if function == "read_shippingschedule" {
		return resource.ReadShippingSchedule(stub, args)

	} else if function == "read_transporttasks_by_ownerid" {
		return resource.ReadAllTransportTasksByOwnerId(stub, args)
	} else if function == "read_transporttask" {
		return resource.ReadTransportTask(stub, args)


	} else if function == "read_orders_by_userid" {
		return order.ReadAllOrdersByUserId(stub, args)
	} else if function == "read_order" {
		return order.ReadOrder(stub, args)


	} else if function == "read_user" {
		return auth.ReadUserById(stub, args)
	} else if function == "find_users_by_role_type" {
		switch args[0] {
		case auth.ROLE_CARGO_AGENT:
			return auth.ReadAllCargoAgents(stub)
		case auth.ROLE_CARRIER:
			return auth.ReadAllCarriers(stub)
		case auth.ROLE_SHIPPER:
			return auth.ReadAllShippers(stub)
		default:
			return nil, errors.New(fmt.Sprintf("Incorrect role type, only supported query on [%s, %s, %s]", auth.ROLE_CARGO_AGENT, auth.ROLE_CARRIER, auth.ROLE_SHIPPER))
		}
	}

	// error out
	fmt.Println("Received unknown query function name - " + function)
	return nil, errors.New("Received unknown invoke function name - '" + function + "'")
}

// ====CHAINCODE EXECUTION SAMPLES (CLI) ==================

// ==== Invoke marbles ====
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["initMarble","marble1","blue","35","tom"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["initMarble","marble2","red","50","tom"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["initMarble","marble3","blue","70","tom"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["transferMarble","marble2","jerry"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["transferMarblesBasedOnColor","blue","jerry"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["delete","marble1"]}'

// ==== Query marbles ====
// peer chaincode query -C myc1 -n marbles -c '{"Args":["readMarble","marble1"]}'
// peer chaincode query -C myc1 -n marbles -c '{"Args":["getMarblesByRange","marble1","marble3"]}'
// peer chaincode query -C myc1 -n marbles -c '{"Args":["getHistoryForMarble","marble1"]}'

// Rich Query (Only supported if CouchDB is used as state database):
//   peer chaincode query -C myc1 -n marbles -c '{"Args":["queryMarblesByOwner","tom"]}'
//   peer chaincode query -C myc1 -n marbles -c '{"Args":["queryMarbles","{\"selector\":{\"owner\":\"tom\"}}"]}'
