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
package order

import (
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	"cbdforum/app-container/chaincode/src/common"
	"cbdforum/app-container/chaincode/src/resource"
	auth "cbdforum/app-container/chaincode/src/user"
)

// OrderBookVehicleRequest - book vehicle request
type OrderBookVehicleRequest struct {
	CargoAgentId string `json:"cargoagentid"`
	CarrierId    string `json:"carrierid"`
	OrderId      string `json:"orderid"`
	StartAt      string `json:"startat"`
	VehicleNum   int    `json:"vehiclenum"`
}

// =============================================================================
// BookVehicle - book vehicle
// Inputs -
// 	0,
// 	BookVehicleRequest,
//
// Returns -
//	0,
//	orderId
// BookVehicle:
// 	- proposed by the cargo agent,
// 	- mainly used to create a new message to notify the client the result
//

func BookVehicle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Info("Starting bookVehicle")
	logger.Infof("Receive %d arguments for book_vehicle: %v\n", len(args), args)

	var err error
	err = CheckArguments(args, 1)
	if err != nil {
		return nil, err
	}
	request := new(OrderBookVehicleRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Error(NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil, NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}
	logger.Infof("Book vehicle, unmarshal request: %+v",request)


	order := new(Order)
	carrier := new(auth.User)
	if isValid, err := request.isValid(stub, order, carrier); !isValid {
		return nil, err
	}

	handleFSM := NewOrderHandleFSM(stub, request.CargoAgentId, request.OrderId, order.State)

	// allocate vehicles declared by the request
	var carryingForm CarryingForm = NewCarryingForm(request)
	carryingForm.Vehicles = resource.AllocateVehicles(stub, request.VehicleNum, request.CarrierId)
	if carryingForm.Vehicles == nil {
		// if resource not enough, keep the status
		//send message to the cargo agent, notifying him that the carriers not have enough vehicles
		err = SendMessage(stub, order.ConsigningForm.CargoAgentId, "Order %s failed to book vehilces, not enough vehicles. Further details are displayed in the order platform.", order.OrderNo)
		if err != nil {
			logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
		}
		logger.Warningf("Failed to allocate %d vehicles, reason is %v", request.VehicleNum, NewShimError(ERROR_RESOURCE, "Vehicles not enough"))
		return nil, NewShimError(ERROR_RESOURCE, "Containers not enough")
	}

	// add a new task for vehicle; put the task into ledger
	task := new(resource.TransportTask)
	task, err = resource.NewTransportTask(stub, request.CarrierId, order.ConsigningForm.ClientId, request.CargoAgentId, request.OrderId, carryingForm.Vehicles, common.StatusInit)
	if err != nil {
		return nil, NewShimError(ERROR_INTERNAL, "Failed to create a new transportTask: %v", err)
	}
	task.StartAt = request.StartAt
	err = task.PutTransportTask(stub)
	if err != nil {
		return nil, NewShimError(ERROR_INTERNAL, "Failed to put transporttask into ledger: %v", err)
	}

	carryingForm.TransportTaskId = task.Id
	order.CarryingForm = carryingForm

	for _, item := range carryingForm.Vehicles {
		err = item.PutVehicle(stub)
		if err != nil {
			return nil, NewShimError(ERROR_INTERNAL, "Failed to update vehicle: %v", err)
		}
	}

	// update order
	err = handleFSM.FSM.Event(EVENT_BOOK_VEHICLE)
	if err != nil {
		return nil, NewShimError(ERROR_INTERNAL, "Failed to book vehicle: %v", err)
	}
	order.State = handleFSM.FSM.Current()

	// write order back into the ledger
	if err = order.PutOrder(stub); err != nil {
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}

	// create message to notify the cargo agent
	err = SendMessage(stub, order.ConsigningForm.CargoAgentId, "Order %s has been processed by carrier, with vehicle booked. Further details are displayed in the order platform.", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}

	// create message to notify the client
	err = SendMessage(stub, order.ConsigningForm.CargoAgentId, "Order %s has been processed by carrier, with vehicle booked. Further details are displayed in the order platform. Please contact carrier for packing goods into the container", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}


	fmt.Println("- end bookVehicle")
	return []byte(order.Id), nil
}

func (request *OrderBookVehicleRequest) isValid(stub shim.ChaincodeStubInterface, order *Order, carrier *auth.User) (bool, error) {
	var err error
	// check orderID exist or not; check cargoagentID and matched or not;
	if !IsOrderExist(stub, request.OrderId, order) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Order %s not exist: %v", request.OrderId, err)
	}
	if order.ConsigningForm.CargoAgentId != request.CargoAgentId {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Cargo agent %s can't modify order %v", request.CargoAgentId, order)
	}

	// check carrier exist or not
	if !auth.IsUserExist(stub, request.CargoAgentId, carrier) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Carrier %s not exist: %v", request.CarrierId, err)
	}
	return true, nil
}

func (request *OrderBookVehicleRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
