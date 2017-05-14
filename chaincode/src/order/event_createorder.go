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
	auth "cbdforum/app-container/chaincode/src/user"
)

type OrderCreateRequest struct {
	ConsigningForm ConsigningForm `json:"consigningform"`
}


// =========================================================================================
// CreateOrder - The client post a request to create a order
// Inputs - Array of strings
// 	0, 			1
//	OrderCreateRequest, 	orderid
//
//
// CreateOrder:
// 	- proposed by the regular client,
// 	- mainly used to create a new order
//	- optionally notify the cargo agent that he/she should handle the order immediately
// =========================================================================================


func CreateOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	var err error
	logger.Info("Starting createOrder")
	logger.Infof("Receive %d arguments for create_order: %v", len(args), args)

	err = CheckArguments(args, 2)
	if err != nil {
		logger.Errorf("checkArguments: Arguments %v not valid: %v", args, err)
		return nil, err
	}

	request := new(OrderCreateRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Errorf("%v", NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil, NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}

	logger.Infof("Info create order, unmarshal request: %+v", request)

	var client *auth.User = new(auth.User)
	if isValid, err := request.isValid(stub, client); !isValid {
		logger.Errorf("Request is not valid: %v", err)
		return nil, err
	}

	logger.Infof("Info create order, request %+v is valid", request)

	order := new(Order)
	order = NewOrder(request)
	order.Id = args[1]
	order.CreatedAt = common.GetCurrentTime()

	// change order state
	logger.Infof("Before transition, current order state is %s", order.State)
	handleFSM := NewOrderHandleFSM(stub, request.ConsigningForm.ClientId, order.Id, order.State)
	err = handleFSM.FSM.Event(EVENT_CREATE_ORDER)
	if err != nil {
		logger.Errorf(ERROR_FSM, "Failed to create order")
		return nil, NewShimError(ERROR_FSM, "Failed to create order")
	}
	order.State = handleFSM.FSM.Current()
	logger.Infof("After transition, current order state is %s", order.State)


	// write order into state
	if err = order.PutOrder(stub); err != nil {
		logger.Error(NewShimError(ERROR_INTERNAL, err.Error()))
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}
	logger.Infof("Create order, put order %+v into the ledger", *order)

	// add this order into the index
	indexName := common.INDEX_ORDER
	values, err := common.GetStringList(stub, indexName)
	if err != nil {
		logger.Errorf("Error trying to get index %s", indexName)
		return nil, err
	}
	values = append(values, order.Id)
	err = common.PutStringList(stub, indexName, values)
	if err != nil {
		logger.Errorf("Failed to put index % back into ledger", indexName)
		return nil, err
	}

	logger.Info("- end createOrder")
	return []byte(order.Id), nil
}

func (request *OrderCreateRequest) isValid(stub shim.ChaincodeStubInterface, client *auth.User) (bool, error) {
	// check whether user exists; check whether user role is valid
	if !auth.IsUserExist(stub, request.ConsigningForm.ClientId, client) {
		logger.Error(NewOrderErrorMessage(ERROR_REQUEST, "User %s not exist", request.ConsigningForm.ClientId))
		return false, NewOrderErrorMessage(ERROR_REQUEST, "User %s not exist", request.ConsigningForm.ClientId)
	}

	logger.Infof("User %+v exist", *client)
	if !client.IsRegularClient() {
		logger.Error(NewOrderErrorMessage(ERROR_REQUEST, "User %v is not regular client", *client))
		return false, NewOrderErrorMessage(ERROR_REQUEST, "User %v is not regular client", *client)
	}

	logger.Infof("OrderCreateRequest %+v post by user %+v is valid",request, client)

	return true, nil
}



func (request *OrderCreateRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
