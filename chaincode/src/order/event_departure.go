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

// OrderLoadGoodsRequest - load goods request
type OrderDepartureRequest struct {
	ShipperId        string `json:"shipperId"`
	OrderId          string `json:"orderid"`
	DateForDeparture string `json:"datefordeparture"` // 货物装船时间
}

// =============================================================================
// Departure -
// arguments:
// 	0,
// 	OrderDepartureRequest,
//
// response:
//	0,
//	orderId
// =============================================================================

func Departure(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	logger.Info("Starting departure")
	logger.Infof("Receive %d arguments for departure: %v\n", len(args), args)

	var err error
	err = CheckArguments(args, 1)
	if err != nil {
		return nil, err
	}
	request := new(OrderDepartureRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Error(NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil,NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}
	logger.Infof("Departure, unmarshal request: %+v",request)

	order := new(Order)
	if isValid, err := request.isValid(stub, order); !isValid {
		return nil, err
	}

	// change the state
	handleFSM := NewOrderHandleFSM(stub, request.ShipperId, request.OrderId, order.State)
	err = handleFSM.FSM.Event(EVENT_DEPARTURE)
	if err != nil {
		return nil, NewShimError(ERROR_FSM, "Failed to loadGoods: %v", err)
	}
	order.State = handleFSM.FSM.Current()

	// update booking form which the carrier is in charge of
	order.BookingForm.DateForDeparture = request.DateForDeparture

	// write order back into the ledger
	if err = order.PutOrder(stub); err != nil {
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}
	logger.Infof("Depatured, order is %v", order)

	// todo: create message to notify the client and the cargo agent
	//create message to notify the client
	err = SendMessage(stub, order.ConsigningForm.ClientId, "Order %s has been processed by shipper, with ship departured. Further details are displayed in the order platform.", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}
	fmt.Println("- end departure")
	return []byte(order.Id), nil
}

func (request *OrderDepartureRequest) isValid(stub shim.ChaincodeStubInterface, order *Order) (bool, error) {

	if !IsOrderExist(stub, request.OrderId, order) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Order %s not exist", request.OrderId)
	}
	if order.BookingForm.ShipperId != request.ShipperId {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Shipper %s can't modify order %v", request.ShipperId, order)
	}
	return true, nil
}

func (request *OrderDepartureRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
