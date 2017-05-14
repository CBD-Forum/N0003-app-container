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
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/looplab/fsm"
)

// OrderHandleFSM
type OrderHandleFSM struct {
	stub       shim.ChaincodeStubInterface
	OrderId    string
	UserId     string
	FSM        *fsm.FSM
}

// NewOrderHandleFSM creates and returns a OrderUpdateFSM
func NewOrderHandleFSM(stub shim.ChaincodeStubInterface, userId string, orderId string, currentState string) *OrderHandleFSM {
	d := &OrderHandleFSM{
		stub:       stub,
		UserId:     userId,
		OrderId:    orderId,
	}

	d.FSM = fsm.NewFSM(
		currentState,
		fsm.Events{
			{Name: EVENT_CREATE_ORDER, Src: []string{STATE_INIT}, Dst: STATE_CREATED},
			{Name: EVENT_CHECK_ORDER, Src: []string{STATE_CREATED}, Dst: STATE_CHECKED},
			{Name: EVENT_BOOK_SPACE, Src: []string{STATE_CHECKED}, Dst: STATE_SPACE_BOOKED},
			{Name: EVENT_BOOK_VEHICLE, Src: []string{STATE_SPACE_BOOKED}, Dst: STATE_VEHICLE_BOOKED},
			{Name: EVENT_FETCH_EMPTY_CONTAINERS, Src: []string{STATE_VEHICLE_BOOKED}, Dst: STATE_EMPTY_CONTAINER_FETCHED},
			{Name: EVENT_PACK_GOODS, Src: []string{STATE_EMPTY_CONTAINER_FETCHED}, Dst: STATE_GOODS_PACKED},
			{Name: EVENT_ARRIVE_YARD, Src: []string{STATE_GOODS_PACKED}, Dst: STATE_YARD_ARRIVED},
			{Name: EVENT_LOAD_GOODS, Src: []string{STATE_YARD_ARRIVED}, Dst: STATE_GOODS_LOADED},
			{Name: EVENT_DEPARTURE, Src: []string{STATE_GOODS_LOADED}, Dst: STATE_GOODS_SHIPPING},
			{Name: EVENT_ARRIVE_DESTINATION_PORT, Src: []string{STATE_GOODS_SHIPPING}, Dst: STATE_GOODS_ARRIVED},
			{Name: EVENT_DELIVER_GOODS, Src: []string{STATE_GOODS_ARRIVED}, Dst: STATE_GOODS_DELIVERED},
			{Name: EVENT_CONFIRM_RECEIPT, Src: []string{STATE_GOODS_DELIVERED}, Dst: STATE_GOODS_RECEIVED},
			{Name: EVENT_FINISH_ORDER, Src: []string{STATE_GOODS_RECEIVED}, Dst: STATE_FINISHED},
			{Name: EVENT_DENY_ORDER, Src: []string{STATE_CREATED}, Dst: STATE_FAILED},
		},
		fsm.Callbacks{
			"enter_state":                             func(e *fsm.Event) { d.enterState(e) },
			"after-events":                            func(e *fsm.Event) { d.afterEvents(e) },
			"before_" + EVENT_CREATE_ORDER:            func(e *fsm.Event) { d.beforeCreateOrder(e) },
			"before_" + EVENT_CHECK_ORDER:             func(e *fsm.Event) { d.beforeCheckOrder(e) },
			"before_" + EVENT_BOOK_SPACE:              func(e *fsm.Event) { d.beforeBookSpace(e) },
			"before_" + EVENT_BOOK_VEHICLE:            func(e *fsm.Event) { d.beforeFetchEmptyContainer(e) },
			"before_" + EVENT_PACK_GOODS:              func(e *fsm.Event) { d.beforePackGoods(e) },
			"before_" + EVENT_ARRIVE_YARD:             func(e *fsm.Event) { d.beforeArriveYard(e) },
			"before_" + EVENT_LOAD_GOODS:              func(e *fsm.Event) { d.beforeLoadGoods(e) },
			"before_" + EVENT_ARRIVE_DESTINATION_PORT: func(e *fsm.Event) { d.beforeArriveDestinationPort(e) },
			"before_" + EVENT_DELIVER_GOODS:           func(e *fsm.Event) { d.beforeDeliverGoods(e) },
			"before_" + EVENT_CONFIRM_RECEIPT:         func(e *fsm.Event) { d.beforeConfirmReceipt(e) },
			"before_" + EVENT_FINISH_ORDER:            func(e *fsm.Event) { d.beforeFinishOrder(e) },
			"before_" + EVENT_DENY_ORDER:              func(e *fsm.Event) { d.beforeDenyOrder(e) },
		},
	)

	return d
}

func (d *OrderHandleFSM) enterState(e *fsm.Event) {
	logger.Debugf("The bi-directional stream to %s is %s, from event %s\n", d.FSM.Current(), e.Dst, e.Event)
}

// called after all events: to put order back into the state
func (d *OrderHandleFSM) afterEvents(e *fsm.Event) {
	logger.Debugf("Leave state from event %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeCreateOrder(e *fsm.Event) {
	logger.Debugf("Before reception of %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())

}

func (d *OrderHandleFSM) beforeCheckOrder(e *fsm.Event) {
	logger.Debugf("After reception of %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())

}

func (d *OrderHandleFSM) afterCheckOrder(e *fsm.Event) {
	logger.Debugf("After reception of %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())

}

func (d *OrderHandleFSM) beforeBookSpace(e *fsm.Event) {
	logger.Debugf("After reception of %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) afterBookSpace(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeBookVehicle(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) afterBookVehicle(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeFetchEmptyContainer(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) afterFetchEmptyContainer(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforePackGoods(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) afterPackGoods(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeArriveYard(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) afterArriveYard(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeLoadGoods(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) afterLoadGoods(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeArriveDestinationPort(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) afterArriveDestinationPort(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeDeliverGoods(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

// todo: add messages to notify the client and the cargo agent
func (d *OrderHandleFSM) afterDeliverGoods(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeConfirmReceipt(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

// todo: add messages to notify the cargo agent
func (d *OrderHandleFSM) afterConfirmReceipt(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeFinishOrder(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

// todo: add messages to
func (d *OrderHandleFSM) afterFinishOrder(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) beforeDenyOrder(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}

func (d *OrderHandleFSM) afterDenyOrder(e *fsm.Event) {
	logger.Debugf("Before %s, dest is %s, current is %s", e.Event, e.Dst, d.FSM.Current())
}
