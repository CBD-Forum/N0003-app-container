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

// OrderCheckRequest - check order request
type OrderCheckRequest struct {
	CargoAgentId  string `json:"cargoagentid"`
	OrderId       string `json:"orderid"`
	IsOrderAccept bool   `json:"isorderaccept"` // 是否接收该订单
	Remark        string `json:"remark"`        // 拒绝的原因
}

// =============================================================================
// CheckOrder - check order
// Inputs -
// 	0,
// 	OrderCheckRequest,
//
// Returns -
//	0,
//	orderId
//
// CheckOrder:
// 	- proposed by the cargo agent,
// 	- mainly used to create a new message to notify the client the result
//

func CheckOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Info("Starting checkOrder")
	logger.Infof("Receive %d arguments for checkOrder: %v", len(args), args)

	err := CheckArguments(args, 1)
	if err != nil {
		logger.Errorf("CheckArguments: arguments %v not valid: %v", args, err)
		return nil, err
	}

	request := new(OrderCheckRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Error(NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil, NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}
	logger.Infof("Info check order, unmarshal request: %+v", request)

	order := new(Order)
	if isValid, err := request.isValid(stub, order); !isValid {
		return nil, err
	}
	logger.Infof("Info check order, request %+v is valid", request)

	// change the state
	handleFSM := NewOrderHandleFSM(stub, request.CargoAgentId, request.OrderId, order.State)

	logger.Infof("Before transition, current order state is %s", order.State)
	if request.IsOrderAccept {
		err = handleFSM.FSM.Event(EVENT_CHECK_ORDER)
		if err != nil {
			logger.Errorf("%v", NewShimError(ERROR_FSM, "Failed to checkOrder: %v", err))
			return nil, NewShimError(ERROR_FSM, "Failed to checkOrder: %v", err)
		}
		order.State = handleFSM.FSM.Current()
	} else {
		err = handleFSM.FSM.Event(EVENT_DENY_ORDER)
		if err != nil {
			return nil, NewShimError(ERROR_FSM, "Failed to deny order in checkOrder: %v", err)
		}
		order.State = handleFSM.FSM.Current()
		order.Remark = request.Remark
	}
	logger.Infof("After transition, current order state is %s", order.State)

	// write order back into the ledger
	if err = order.PutOrder(stub); err != nil {
		logger.Errorf("%v",NewShimError(ERROR_INTERNAL, err.Error()))
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}

	// create message to notify the client
	err = SendMessage(stub, order.ConsigningForm.ClientId, "Order %s has been processed by cargo agent. Further details are displayed in the order platform.", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}

	fmt.Println("- end checkOrder")
	return []byte(order.Id), nil
}

func (request *OrderCheckRequest) isValid(stub shim.ChaincodeStubInterface, order *Order) (bool, error) {
	var err error
	if !IsOrderExist(stub, request.OrderId, order) {
		logger.Errorf(ERROR_REQUEST, "Order %s not exist: %v", request.OrderId, err)
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Order %s not exist: %v", request.OrderId, err)
	}
	if order.ConsigningForm.CargoAgentId != request.CargoAgentId {
		logger.Errorf(ERROR_REQUEST, "Cargo agent %s can't modify order %v", request.CargoAgentId, order)
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Cargo agent %s can't modify order %v", request.CargoAgentId, order)
	}
	logger.Infof("OrderCheckRequest %+v on order %+v is valid",request, order)
	return true, nil
}

func (request *OrderCheckRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
