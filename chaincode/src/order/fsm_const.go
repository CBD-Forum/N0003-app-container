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

// state transfer, 			event,			action,		role
// init --> created: 			createOrder,				client
// created --> checked:			checkOrder,				cargoAgent
// checked --> spacebooked:		bookSpace,				cargoAgent <--> shipper
// checked --> spacedenied:		bookSpace, 				shipper --> cargoAgent

// spacebooked --> vehiclebooked: 	bookVehicle,				cargoagent, carrier
// spacebooked --> vehicledenied:	bookVehicle,				carrier --> cargoagent

// vehiclebooked --> emptycontainerfetched:  fetchEmptyContainer		cargoagent
// emptycontainerfetched --> goodsPacked:	packGoods			carrier
// goodsPacked 	--> yardArrived:		arriveYard			carrier

// yardArrived 	--> goodsLoaded:		loadGoods			shipper
// goodsLoaded  --> goodsShipping:		departure			shipper
// goodsShipping --> goodsArrived: 		arriveDestinationPort		shipper
// goodsArrived --> goodsDelivered		deliverGoods			shipper
// goodsDelivered --> confirmed			confirmReceipt			client
// confirmed --> finished			finishOrder			cargoagent

// created --> denied:			checkOrder, 		order failed	cargoAgent,

// =============================== cargo agent 重置状态，重新订舱 ============================
// spacedenied --> checked:		resetSpace,				cargoagent

// =============================== cargo agent 重置状态，重新订车 ============================
// vehicledenied --> spacebooked:	resetVehicle				cargoagent

const (
	STATE_INIT    string = "order_init" //订单一开始由client创建后, 由 init 态 转换为 created 态，下一步是 cargo agent 处理，
	STATE_CREATED            = "order_created"

	// cargo agent
	STATE_CHECKED        = "order_checked"        // cargo agent 审核订单，确定是否接单
	STATE_SPACE_BOOKED   = "order_space_booked"   // cargo agent 接单之后， 订舱
	STATE_VEHICLE_BOOKED = "order_vehicle_booked" // cargo agent 订舱成功之后，订车

	// carrier
	STATE_EMPTY_CONTAINER_FETCHED = "order_empty_container_fetched" // carrier
	STATE_GOODS_PACKED            = "order_goods_packed"
	STATE_YARD_ARRIVED            = "order_yard_arrived"

	// shipper
	STATE_GOODS_LOADED    = "order_goods_loaded"
	STATE_GOODS_SHIPPING  = "order_goods_shipping"
	STATE_GOODS_ARRIVED   = "order_goods_arrived"
	STATE_GOODS_DELIVERED = "order_goods_delivered"

	STATE_GOODS_RECEIVED = "order_goods_received"

	STATE_FINISHED = "order_finished" //cargo agent 终止订单，订单完成
	STATE_FAILED   = "order_failed"   //订单失败
)

const (
	EVENT_CREATE_ORDER            string = "create_order"
	EVENT_CHECK_ORDER                        = "check_order"
	EVENT_BOOK_SPACE                         = "book_space"
	EVENT_BOOK_VEHICLE                       = "book_vehicle"
	EVENT_FETCH_EMPTY_CONTAINERS             = "fetch_empty_containers"
	EVENT_PACK_GOODS                         = "pack_goods"
	EVENT_ARRIVE_YARD                        = "arrive_yard"
	EVENT_LOAD_GOODS                         = "load_goods"
	EVENT_DEPARTURE                          = "departure"
	EVENT_ARRIVE_DESTINATION_PORT            = "arrive_destination_port"
	EVENT_DELIVER_GOODS                      = "deliver_goods"
	EVENT_CONFIRM_RECEIPT                    = "confirm_receipt"
	EVENT_FINISH_ORDER                       = "finish_order"
	EVENT_DENY_ORDER                         = "deny_order"
)
