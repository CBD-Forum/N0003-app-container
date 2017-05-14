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

// OrderConfirmReceiptRequest -
type OrderFinishRequest struct {
	CargoAgentId  string `json:"cargoagentid"`
	OrderId       string `json:"orderid"`
	DateForFinish string `json:"dateforfinish"` // 订单完成关闭时间
}

// =============================================================================
// Departure -
// arguments:
// 	0,
// 	OrderFinishRequest,
//
// response:
//	0,
//	orderId
// =============================================================================

func FinishOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Info("Starting finishOrder")
	logger.Infof("Receive %d arguments for finish_order: %v\n", len(args), args)

	var err error
	err = CheckArguments(args, 1)
	if err != nil {
		return nil, err
	}
	request := new(OrderFinishRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Error(NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil, NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}
	logger.Infof("Finish order, unmarshal request: %+v",request)

	var order *Order = new(Order)
	if isValid, err := request.isValid(stub, order); !isValid {
		return nil, err
	}

	// change the state
	handleFSM := NewOrderHandleFSM(stub, request.CargoAgentId, request.OrderId, order.State)
	err = handleFSM.FSM.Event(EVENT_FINISH_ORDER)
	if err != nil {
		return nil, NewShimError(ERROR_FSM, "Failed to finishOrder: %v", err)
	}
	order.State = handleFSM.FSM.Current()

	// update consigning form which the client and the carrier is in charge of
	order.ConsigningForm.DateForFinish = request.DateForFinish

	// write order back into the ledger
	if err = order.PutOrder(stub); err != nil {
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}
	logger.Infof("Finished order, order is %v", order)
	// todo: create message to notify the client that the order has finished.
	err = SendMessage(stub, order.ConsigningForm.ClientId, "Order %s has been processed by cargoagent, with order finished. Further details are displayed in the order platform.", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}

	fmt.Println("- end finishOrder")
	return []byte(order.Id), nil
}

func (request *OrderFinishRequest) isValid(stub shim.ChaincodeStubInterface, order *Order) (bool, error) {
	if !IsOrderExist(stub, request.OrderId, order) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Order %s not exist.", request.OrderId)
	}

	if order.ConsigningForm.CargoAgentId != request.CargoAgentId {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Cargo agent %s can't modify order %v", request.CargoAgentId, order)
	}
	return true, nil
}

func (request *OrderFinishRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
