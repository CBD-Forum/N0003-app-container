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
type OrderConfirmReceiptRequest struct {
	ClientId              string `json:"clientid"`
	OrderId               string `json:"orderid"`
	DateForConfirmReceipt string `json:"dateforconfirmreceipt"` // 货物确认签收时间
}

// =============================================================================
// ConfirmReceipt - to confirm that the client has received the goods, proposed by the client
// Inputs - Array of strings
// 	0,
// 	OrderConfirmReceipt,
//
// Returns -
//	0,
//	orderId
// =============================================================================
func ConfirmReceipt(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Info("Starting confirmReceipt")
	logger.Infof("Receive %d arguments for confirmReceipt: %v\n", len(args), args)

	var err error
	err = CheckArguments(args, 1)
	if err != nil {
		return nil, err
	}
	request := new(OrderConfirmReceiptRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Error(NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil, NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}
	logger.Infof("Confirm receipt, unmarshal request: %+v",request)

	order := new(Order)
	// check orderID exist or not; check shipperId matched or not
	if isValid, err := request.isValid(stub, order); !isValid {
		return nil, err
	}

	// change the state
	handleFSM := NewOrderHandleFSM(stub, request.ClientId, request.OrderId, order.State)
	err = handleFSM.FSM.Event(EVENT_CONFIRM_RECEIPT)
	if err != nil {
		return nil, NewShimError(ERROR_FSM, "Failed to confirmReceipt: %v", err)
	}
	order.State = handleFSM.FSM.Current()

	// update consigning form which the client and the carrier is in charge of
	order.ConsigningForm.DateForConfirmReceipt = request.DateForConfirmReceipt

	// write order back into the ledger
	if err = order.PutOrder(stub); err != nil {
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}

	logger.Infof("Confirmed receipt, order is %v", order)
	// todo: create message to notify the the cargo agent
	err = SendMessage(stub, order.ConsigningForm.CargoAgentId, "Order %s has been processed by client, with goods delivery received by client. Further details are displayed in the order platform.", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}

	fmt.Println("- end confirmReceipt")
	return []byte(order.Id), nil
}

func (request *OrderConfirmReceiptRequest) isValid(stub shim.ChaincodeStubInterface, order *Order) (bool, error) {
	if !IsOrderExist(stub, request.OrderId, order){
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Order %s not exist", request.OrderId)
	}
	if order.ConsigningForm.ClientId != request.ClientId {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Client %s can't modify order %v", request.ClientId, order)
	}
	return true, nil
}

func (request *OrderConfirmReceiptRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
