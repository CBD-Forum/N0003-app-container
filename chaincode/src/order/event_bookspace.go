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

// OrderBookSpaceRequest - book space request
type OrderBookSpaceRequest struct {
	CargoAgentId  string `json:"cargoagentid"`
	ShipperId     string `json:"shipperid"`
	OrderId       string `json:"orderid"`
	VoyageId      string `json:"voyageid"` // 航次号
	ContainerType string `json:"containertype"`
	ContainerNum  int    `json:"containernum"` // 集装箱数
	SpaceNum      int    `json:"spacenum"`     // 舱位数
}

// =============================================================================
// BookSpace - book space
// Inputs - Array of strings
// 	0,
// 	BookSpaceRequest,
// Returns -
//	0,
//	orderId
//
// BookSpace:
// 	- proposed by the cargo agent,
// 	- mainly used to create a new message to notify the client the result
// =============================================================================
func BookSpace(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	var err error
	logger.Info("starting bookSpace")
	logger.Infof("receive %d arguments for book_space: %v\n", len(args), args)

	err = CheckArguments(args, 1)
	if err != nil {
		return nil, err
	}

	request := new(OrderBookSpaceRequest)
	err = json.Unmarshal([]byte(args[0]), request)
	if err != nil {
		logger.Error(NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType()))
		return nil, NewOrderErrorMessage(ERROR_ARGUMENTS, "Incorrect type, expecting %s", request.GetType())
	}
	logger.Infof("Book space, unmarshal request: %+v",request)


	// check orderID exist or not; check cargoagentID matched or not
	order := new(Order)
	shipper := new(auth.User)
	voyage := new(resource.ShippingSchedule)
	if isValid, err := request.isValid(stub, order, shipper, voyage); !isValid {
		logger.Errorf("Book space request %v is not valid: %v", request, err)
		return nil, err
	}
	logger.Infof("Info book space, request %+v is valid", request)


	handleFSM := NewOrderHandleFSM(stub, request.CargoAgentId, request.OrderId, order.State)

	// 分配舱位、集装箱，若失败，则状态不改变
	if !voyage.HasSpace(request.SpaceNum) {
		// if resource not enough, keep the status
		//send message to the cargo agent, notifying him that the shipping schedule he selected does not have enough space
		err = SendMessage(stub, order.ConsigningForm.CargoAgentId, "Order %s failed to book space, not enough shipping space. Further details are displayed in the order platform.", order.OrderNo)
		if err != nil {
			logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
		}

		logger.Errorf("%v", NewOrderErrorMessage(ERROR_RESOURCE, "Space not enough"))
		return nil, NewOrderErrorMessage(ERROR_RESOURCE, "Space not enough")
	}

	logger.Infof("Voyage has space !")

	// update booking form
	var bookingForm BookingForm
	bookingForm.ShipperId = request.ShipperId
	bookingForm.BerthNo = voyage.AllocateSpace(request.SpaceNum)
	bookingForm.BookingFormNo = common.GenerateUUID()
	bookingForm.Voyage = *voyage
	bookingForm.Containers = resource.AllocateContainers(stub, request.ContainerNum, request.ContainerType, request.ShipperId)
	if bookingForm.Containers == nil {
		// if resource not enough, keep the status
		//send message to the cargo agent, notifying him that the shipping schedule he selected does not have enough space
		err = SendMessage(stub, order.ConsigningForm.CargoAgentId, "Order %s failed to book space, not enough containers. Further details are displayed in the order platform.", order.OrderNo)
		if err != nil {
			logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
		}
		logger.Warningf("Failed to allocate %d containers, reason is %v", request.ContainerNum, NewShimError(ERROR_RESOURCE, "Containers not enough"))
		return nil, NewShimError(ERROR_RESOURCE, "Containers not enough")
	}
	order.BookingForm = bookingForm
	logger.Infof("Allocated space and containers, bookingform is %+v", bookingForm)

	// update voyage and containers
	err = bookingForm.Voyage.PutShippingSchedule(stub)
	if err != nil {
		logger.Errorf("%v", NewShimError(ERROR_INTERNAL, "Failed to update voyage: %v", err))
		return nil, NewShimError(ERROR_INTERNAL, "Failed to update voyage: %v", err)
	}
	logger.Infof("Updated voyage %v in book space!", bookingForm.Voyage)

	for _, item := range bookingForm.Containers {
		err = item.PutContainer(stub)
		if err != nil {
			logger.Errorf("%v", NewShimError(ERROR_INTERNAL, "Failed to update container %v: %v", item, err))
			return nil, NewShimError(ERROR_INTERNAL, "Failed to update container %v: %v", item, err)
		}
	}
	logger.Infof("Update containers %v in book space!", bookingForm.Containers)

	// update order state
	err = handleFSM.FSM.Event(EVENT_BOOK_SPACE)
	if err != nil {
		logger.Errorf("%v", NewShimError(ERROR_FSM, "Failed to book space: %v", err))
		return nil, NewShimError(ERROR_FSM, "Failed to book space: %v", err)
	}
	order.State = handleFSM.FSM.Current()

	// write order back into the ledger
	if err = order.PutOrder(stub); err != nil {
		logger.Errorf("%v", NewShimError(ERROR_INTERNAL, err.Error()))
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}
	logger.Infof("! Book space successed, order is %+v", order)

	// create message to notify the cargo agent
	err = SendMessage(stub, order.ConsigningForm.CargoAgentId, "Order %s has been processed by shipper, with shipping space booked. Further details are displayed in the order platform.", order.OrderNo)
	if err != nil {
		logger.Warningf("Failed to send message to user %s: %v", order.ConsigningForm.ClientId, err)
	}

	fmt.Println("- end bookSpace")
	return []byte(order.Id), nil
}

func (request *OrderBookSpaceRequest) isValid(stub shim.ChaincodeStubInterface, order *Order, shipper *auth.User, voyage *resource.ShippingSchedule) (bool, error) {
	var err error
	// check orderID exist or not; check cargoagentID matched or not
	if !IsOrderExist(stub, request.OrderId, order) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Order %s not exist: %v", request.OrderId, err)
	}
	if order.ConsigningForm.CargoAgentId != request.CargoAgentId {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Cargo agent %s can't modify order %v", request.CargoAgentId, order)
	}

	// check shipper exist or not; check voyage exist or not; check voyage owned by shipper or not
	if !auth.IsUserExist(stub, request.ShipperId, shipper) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Shipper %s not exist: %v", request.ShipperId, err)
	}

	if !resource.IsShippingScheduleExist(stub, request.VoyageId, voyage) {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Voyage %s not exist: %v", request.VoyageId, err)
	}

	if voyage.OwnerId != request.ShipperId {
		return false, NewOrderErrorMessage(ERROR_REQUEST, "Voyage %s not belong to shipper %s", request.VoyageId, request.ShipperId)
	}
	return true, nil
}

func (request *OrderBookSpaceRequest) GetType() string {
	return fmt.Sprintf("%T", *request)
}
