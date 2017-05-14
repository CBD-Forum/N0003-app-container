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
)

// OrderArriveYardRequest - arrive yard request
type OrderArriveYardRequest struct {
	CarrierId       string `json:"carrierId"`
	OrderId         string `json:"orderid"`
	DateForReceiver string `json:"dateforreceiver"` // 设备进场时间
}

// =============================================================================
// ArriveYard - containers and goods arrive yard, waiting for being loaded into the vessel
// arguments:
// 	0,
// 	OrderArriveYardRequest,
//
// response:
//	0,
//	orderId
// =============================================================================

func ArriveYard(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting packGoods")

	err = CheckArguments(args, 1)
	if err != nil {
		return nil, err
	}
	request := new(OrderArriveYardRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Error(NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil, NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}
	logger.Infof("Arrive yard, unmarshal request: %+v",request)


	order := new(Order)
	if isValid, err := request.isValid(stub, order); !isValid {
		return nil, err
	}

	// change the state
	handleFSM := NewOrderHandleFSM(stub, request.CarrierId, request.OrderId, order.State)
	err = handleFSM.FSM.Event(EVENT_ARRIVE_YARD)
	if err != nil {
		return nil, NewShimError(ERROR_FSM, "Failed to arriveYard: %v", err)
	}
	order.State = handleFSM.FSM.Current()

	// todo: update transport task status as finished, set the status of vehicles as free
	order.CarryingForm.DateForReceiver = request.DateForReceiver
	order.CarryingForm.Status = handleFSM.FSM.Current()

	// write order back into the ledger
	if err = order.PutOrder(stub); err != nil {
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}
	logger.Infof("Arrived yard, order is %v", order)

	// todo: create message to notify the client and the cargo agent
	//create message to notify the client
	err = SendMessage(stub, order.ConsigningForm.ClientId, "Order %s has been processed by carrier, with goods tranported to yards. Further details are displayed in the order platform.", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}

	fmt.Println("- end arriveYard")
	return []byte(order.Id), nil
}

func (request *OrderArriveYardRequest) isValid(stub shim.ChaincodeStubInterface, order *Order) (bool, error) {
	var err error
	if !IsOrderExist(stub, request.OrderId, order) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Order %s not exist: %v", request.OrderId, err)
	}
	if order.CarryingForm.CarrierId != request.CarrierId {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Carrier %s can't modify order %v", request.CarrierId, order)
	}
	return true, nil
}

func (request *OrderArriveYardRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
