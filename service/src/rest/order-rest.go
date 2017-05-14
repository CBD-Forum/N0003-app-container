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
package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"io/ioutil"
	"net/http"

	"cbdforum/app-container/service/src/common"
	"cbdforum/app-container/service/src/user"
)

type Goods struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Measurement int    `json:"measurement"`
	GrossWeight int    `json:"grossweight"`
}

type Order struct {
	ObjectType string `json:"docType"` //field for CouchDB
	Id         string `json:"id"`
	OrderNo    string `json:"orderno"`

	ConsigningForm ConsigningForm `json:"consigningform"`
	BookingForm    BookingForm    `json:"bookingform"`
	CarryingForm   CarryingForm   `json:"carryingform"`

	State        string   `json:"state"`
	DeletedByWho []string `json:"deletedbywho"`
	Remark       string   `json:"remark"` // 用来记录订单被拒绝的原因

	CreatedAt string `json:"createdat"` // 订单创建的时间
}

type ConsigningForm struct {
	ConsigningFormNo string `json:"consigningformno"`
	ClientId         string `json:"clientid"`
	CargoAgentId     string `json:"cargoagentid"`

	GoodsList            []Goods `json:"goodslist"`
	DeliveryAddress      string  `json:"deliveryaddress"`
	ShippingAddress      string  `json:"shippingaddress"`
	Consignee            string  `json:"consignee"`
	ConsigneePhone       string  `json:"consigneePhone"`
	Consignor            string  `json:"consignor"`
	ConsignorPhone       string  `json:"consignorPhone"`
	ExpectedDeliveryDate string  `json:"expecteddeliverydate"`

	DateForConfirmReceipt string `json:"dateforconfirmreceipt"` // 确认签收时间, 由 client 写入
	DateForFinish         string `json:"dateforfinish"`         // 订单最终完成时间，由cargo agent 写入

}

type BookingForm struct {
	BookingFormNo       string           `json:"bookingformno"`
	ShipperId           string           `json:"shipperid"`
	Voyage              ShippingSchedule `json:"voyage"`
	BerthNo             []string         `json:"berthno"`
	Containers          []Container      `json:"containers"`
	DateForLoading      string           `json:"dateforloading"`      // 货物装船时间
	DateForDeparture    string           `json:"datefordeparture"`    // 起航时间
	DateForArrival      string           `json:"dateforarrival"`      // 到达时间
	DateForDeliverGoods string           `json:"datefordelivergoods"` // 送货时间
}

type CarryingForm struct {
	CarryingFormNo  string      `json:"carryingformno"`
	CarrierId       string      `json:"carrierid"`
	TransportTaskId string      `json:"transporttaskid"`
	Vehicles        []Vehicle   `json:"vehicles"`
	PackingList     PackingList `json:"packinglist"`
	Status          string      `json:"status"`
	DateForReceiver string      `json:"dateforreceiver"` // 设备移入时间
	DateForDeliver  string      `json:"datefordeliver"`  // 设备移出时间
}

type PackingList struct {
	Items               []PackingListItem `json:"items"`
	DateForPackingGoods string            `json:"dateforpackinggoods"` //装箱时间
}

type PackingListItem struct {
	ContainerId string `json:"containerid"`
	Goods       Goods  `json:"goods"`
}

// ==================================================================================
// GetOrderById -
// route - /order/{orderId}
// method - GET
// ==================================================================================
func (s *ServerContainerREST) GetOrderById(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")
	id := req.PathParams["id"]

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session"})
		restLogger.Errorf("Invalid session")
		return
	}

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_order", id)

	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress+"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to connect to fabric peer: %v", err)})
		restLogger.Errorf("Error trying to connect to fabric peer: %v", err)
		return
	}

	data, _ := ioutil.ReadAll(ccResp.Body)
	defer ccResp.Body.Close()
	if ccResp.StatusCode != http.StatusOK {
		rw.WriteHeader(ccResp.StatusCode)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read order %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read order %s from the ledger: %s", id, bytes.NewBuffer(data).String())
		return
	}
	var rpcResponse CCResponse
	err = json.Unmarshal(data, &rpcResponse)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Failed to unmarshal response from peer: %v", err)})
		restLogger.Errorf("Failed to unmarshal response from peer: %v", err)
		return
	}

	restLogger.Infof("Read order by id %s from the ledger: %s", id, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}

// ==================================================================================
// DeleteOrder -
// route - /order/{orderId}
// method - DELETE
// ==================================================================================
func (s *ServerContainerREST) DeleteOrder(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")
	id := req.PathParams["id"]

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session"})
		restLogger.Errorf("Invalid session")
		return
	}

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "delete_order", id, userid)

	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress+"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to connect to fabric peer: %v", err)})
		restLogger.Errorf("Error trying to connect to fabric peer: %v", err)
		return
	}

	data, _ := ioutil.ReadAll(ccResp.Body)
	defer ccResp.Body.Close()
	if ccResp.StatusCode != http.StatusOK {
		rw.WriteHeader(ccResp.StatusCode)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error delete order %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error delete order %s from the ledger: %s", id, bytes.NewBuffer(data).String())
		return
	}
	restLogger.Infof("Delte order %s from the ledger: %s", id, bytes.NewBuffer(data).String())

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: bytes.NewBuffer(data).String()})
}

// ==================================================================================
// FindOrdersByUserId -
// route - /order/findByUserId
// method - GET
// ==================================================================================
func (s *ServerContainerREST) FindOrdersByUserId(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session"})
		restLogger.Errorf("Invalid session")
		return
	}

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_orders_by_userid", userid)

	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress+"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to connect to fabric peer: %v", err)})
		restLogger.Errorf("Error trying to connect to fabric peer: %v", err)
		return
	}

	data, _ := ioutil.ReadAll(ccResp.Body)
	defer ccResp.Body.Close()
	if ccResp.StatusCode != http.StatusOK {
		rw.WriteHeader(ccResp.StatusCode)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read orders for user %s from the ledger: %s", userid, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read orders for user %s from the ledger: %s", userid, bytes.NewBuffer(data).String())
		return
	}
	var rpcResponse CCResponse
	err = json.Unmarshal(data, &rpcResponse)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Failed to unmarshal response from peer: %v", err)})
		restLogger.Errorf("Failed to unmarshal response from peer: %v", err)
		return
	}

	restLogger.Infof("Read orders by userid %s from the ledger: %s", userid, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}

// ==================================================================================
// CreateOrder -
// route - /order/create
// method - POST
// ==================================================================================
func (s *ServerContainerREST) CreateOrder(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	// Parse out the user enrollment ID
	userid := req.Header.Get("userid")
	token := req.Header.Get("token")
	sessionid := req.Header.Get("sessionid")

	orderid := common.GenerateUUID()

	encoder := json.NewEncoder(rw)
	//db, err := sql.Open(viper.GetString("database.name"), viper.GetString("database.dsn"))

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session, userid, sessoinid and token not set in the header or session has expired"})
		restLogger.Errorf("Invalid session")
		return
	}

	//insert the vehicle into the ledger
	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil || data == nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: "Error trying to create an empty order into the ledger"})
		restLogger.Error("Error trying to create an empty order into the ledger")
		return
	}

	argsBuffer := bytes.NewBuffer(data)
	restLogger.Infof("Info reqire for %s %s: %v", req.RoutePath(), req.Method, argsBuffer.String())

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "create_order", argsBuffer.String(), orderid)

	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress+"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to connect to fabric peer: %v", err)})
		restLogger.Errorf("Error trying to connect to fabric peer: %v", err)
		return
	}
	defer ccResp.Body.Close()

	data, _ = ioutil.ReadAll(ccResp.Body)
	if ccResp.StatusCode != http.StatusOK {
		rw.WriteHeader(ccResp.StatusCode)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error inserting container into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error inserting container into the ledger: %s", bytes.NewBuffer(data).String())
		return
	}

	type AuthResult struct {
		UserId  string `json:"userid"`
		OrderId string `json:"orderid"`
		Message string `json:"message"`
	}

	auth := AuthResult{UserId: userid, OrderId: orderid, Message: bytes.NewBuffer(data).String()}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("Create order %v, response %s", orderid, buffer.String())
}

// ==================================================================================
// CheckOrder -
// route - /order/check
// method - POST
// ==================================================================================
func (s *ServerContainerREST) CheckOrder(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "check_order")
}

// ==================================================================================
// BookSpace -
// route - /order/bookspace
// method - POST
// ==================================================================================
func (s *ServerContainerREST) BookSpace(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "book_space")
}

// ==================================================================================
// BookVehicle -
// route - /order/bookvehicle
// method - POST
// ==================================================================================
func (s *ServerContainerREST) BookVehicle(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "book_vehicle")
}

// ==================================================================================
// FetchEmptyContainers -
// route - /order/fetchemptycontainers
// method - POST
// ==================================================================================
func (s *ServerContainerREST) FetchEmptyContainers(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "fetch_empty_containers")
}

// ==================================================================================
// PackGoods -
// route - /order/packgoods
// method - POST
// ==================================================================================
func (s *ServerContainerREST) PackGoods(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "pack_goods")
}

// ==================================================================================
// ArriveYard -
// route - /order/arriveyard
// method - POST
// ==================================================================================
func (s *ServerContainerREST) ArriveYard(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "arrive_yard")
}

// ==================================================================================
// LoadGoods -
// route - /order/loadgoods
// method - POST
// ==================================================================================
func (s *ServerContainerREST) LoadGoods(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "load_goods")
}

// ==================================================================================
// Departure -
// route - /order/departure
// method - POST
// ==================================================================================
func (s *ServerContainerREST) Departure(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "departure")
}

// ==================================================================================
// ArriveDestinationPort -
// route - /order/arrivedestinationport
// method - POST
// ==================================================================================
func (s *ServerContainerREST) ArriveDestinationPort(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "arrive_destination_port")
}

// ==================================================================================
// DeliverGoods -
// route - /order/delivergoods
// method - POST
// ==================================================================================
func (s *ServerContainerREST) DeliverGoods(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "deliver_goods")
}

// ==================================================================================
// ConfirmReceipt -
// route - /order/confirmreceipt
// method - POST
// ==================================================================================
func (s *ServerContainerREST) ConfirmReceipt(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "confirm_receipt")
}

// ==================================================================================
// FinishOrder -
// route - /order/finish
// method - POST
// ==================================================================================
func (s *ServerContainerREST) FinishOrder(rw web.ResponseWriter, req *web.Request) {
	handleOrder(rw, req, "finish_order")
}

func handleOrder(rw web.ResponseWriter, req *web.Request, functionName string) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	// Parse out the user enrollment ID
	userid := req.Header.Get("userid")
	token := req.Header.Get("token")
	sessionid := req.Header.Get("sessionid")

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session, userid, sessoinid and token not set in the header or session has expired"})
		restLogger.Errorf("Invalid session")
		return
	}

	//insert the vehicle into the ledger
	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil || data == nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to %s,  empty request", functionName)})
		restLogger.Errorf("Error trying to %s, empty request", functionName)
		return
	}

	argsBuffer := bytes.NewBuffer(data)
	restLogger.Infof("Info require for %s %s: %v", req.RoutePath(), req.Method, argsBuffer.String())
	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, functionName, argsBuffer.String())

	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress+"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to connect to fabric peer: %v", err)})
		restLogger.Errorf("Error trying to connect to fabric peer: %v", err)
		return
	}
	defer ccResp.Body.Close()

	data, _ = ioutil.ReadAll(ccResp.Body)
	if ccResp.StatusCode != http.StatusOK {
		rw.WriteHeader(ccResp.StatusCode)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error %s: %s", functionName, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error %s: %s", functionName, bytes.NewBuffer(data).String())
		return
	}

	type AuthResult struct {
		UserId  string `json:"userid"`
		Message string `json:"message"`
	}

	auth := AuthResult{UserId: userid, Message: bytes.NewBuffer(data).String()}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("%s, response %s", functionName, buffer.String())
}
