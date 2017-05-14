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
	"cbdforum/app-container/chaincode/src/resource"
	"cbdforum/app-container/chaincode/src/common"
)

// OrderLoadGoodsRequest - load goods request
type OrderDeliverGoodsRequest struct {
	ShipperId           string `json:"shipperId"`
	OrderId             string `json:"orderid"`
	DateForDeliverGoods string `json:"datefordelivergoods"` // 货物装船时间
}

// =============================================================================
// Departure -
// arguments:
// 	0,
// 	OrderDeliverGoods,
//
// response:
//	0,
//	orderId
// =============================================================================

func DeliverGoods(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Info("Starting deliverGoods")
	logger.Infof("Receive %d arguments for deliver_goods: %v\n", len(args), args)

	var err error
	err = CheckArguments(args, 1)
	if err != nil {
		return nil, err
	}
	request := new(OrderDeliverGoodsRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Error(NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil, NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}
	logger.Infof("Deliver goods, unmarshal request: %+v",request)


	order := new(Order)
	if isValid, err := request.isValid(stub, order); !isValid {
		return nil, err
	}

	// change the state
	handleFSM := NewOrderHandleFSM(stub, request.ShipperId, request.OrderId, order.State)
	err = handleFSM.FSM.Event(EVENT_DELIVER_GOODS)
	if err != nil {
		return nil, NewShimError(ERROR_FSM, "Failed to deliverGoods: %v", err)
	}
	order.State = handleFSM.FSM.Current()

	// update booking form which the carrier is in charge of
	order.BookingForm.DateForDeliverGoods = request.DateForDeliverGoods

	// reset container location (set as destination port) and status
	for _, item := range order.BookingForm.Containers {
		var container *resource.Container = new(resource.Container)
		container = &item
		container.Location = order.BookingForm.Voyage.PortOfDischarge
		container.Status = common.ResourceStatusFree
		err = container.PutContainer(stub)
		if err != nil {
			return nil, err
		}
	}

	// write order back into the ledger
	if err = order.PutOrder(stub); err != nil {
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}

	logger.Infof("Delivered goods, order is %v", order)

	// todo: create message to notify the client and the cargo agent
	err = SendMessage(stub, order.ConsigningForm.ClientId, "Order %s has been processed by shipper, with goods delivered to consigner. Further details are displayed in the order platform.", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}

	fmt.Println("- end departure")
	return []byte(order.Id), nil
}

func (request *OrderDeliverGoodsRequest) isValid(stub shim.ChaincodeStubInterface, order *Order) (bool, error) {
	if !IsOrderExist(stub, request.OrderId, order) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Order %s not exist", request.OrderId)
	}
	if order.BookingForm.ShipperId != request.ShipperId {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Shipper %s can't modify order %v", request.ShipperId, order)
	}
	return true, nil
}

func (request *OrderDeliverGoodsRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
