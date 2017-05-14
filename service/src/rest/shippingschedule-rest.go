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
	"github.com/gocraft/web"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"cbdforum/app-container/service/src/common"
	"cbdforum/app-container/service/src/user"
)

type ShippingSchedule struct {
	ObjectType      string `json:"docType"` //field for couchdb
	Id              string `json:"id"`
	OwnerId         string `json:"ownerid"`
	VoyNo           string `json:"voyno"`
	Vessel          Vessel `json:"vessel"`
	PortOfLoading   string `json:"portofloading"`
	PortOfDischarge string `json:"portofdischarge"`
	PlaceOfDelivery string `json:"placeofdelivery"`
	DepartureDate   string `json:"departuredate"`
	ArrivalDate     string `json:"arrivaldate"`
	Status          string `json:"status"`
}

type Space struct {
	TotalNum    int      `json:"totalnum"`    //舱位总数
	RestNum     int      `json:"restnum"`     //舱位剩余数目
	SpaceUsed   []string `json:"spaceused"`   //已使用的舱位号列表
	SpaceUnUsed []string `json:"spaceunused"` //未使用的舱位号列表
}

type Vessel struct {
	VesselNo string `json:"vesselno"`
	Name     string `json:"name"`
	OwnerId  string `json:"ownerid"`
	Space    Space  `json:"space"` // 舱位
}

// ==================================================================================
// InsertShippingSchedule -
// route - /resource/shippingschedule
// method - POST
// ==================================================================================
func (s *ServerContainerREST) InsertShippingSchedule(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	// Parse out the user enrollment ID
	userid := req.Header.Get("userid")
	token := req.Header.Get("token")
	sessionid := req.Header.Get("sessionid")

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
		encoder.Encode(restResult{Error: "Error trying to add an empty shhipping schedule into the ledger"})
		restLogger.Error("Error trying to register an empty shipping scheduleinto the ledger")
		return
	}

	var value ShippingSchedule
	err = json.Unmarshal(data, &value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to unmarshal input %s into shipping schedule: %v", bytes.NewBuffer(data).String(), err)})
		restLogger.Errorf("Error trying to unmarshal input %s into shipping schedule: %v",bytes.NewBuffer(data).String(), err )
		return
	}
	// generate vehicle id
	value.Id = common.GenerateUUID()
	data, err = json.Marshal(value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to marshal shipping schedule %v into bytes: %v",value, err )})
		restLogger.Errorf("Error trying to marshal shipping schedule %v into bytes: %v", value, err)
		return
	}

	argsBuffer := bytes.NewBuffer(data)
	restLogger.Infof("Info reqire for /vehicle POST: %v", argsBuffer.String())
	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "insert_shippingschedule", argsBuffer.String())


	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress +"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error inserting shipping schedule %v into the ledger", value)})
		restLogger.Errorf("Error inserting shipping schedule %v into the ledger", value)
		return
	}


	type AuthResult struct {
		UserId    string `json:"userid"`
		ShippingScheduleId string `json:"shippingscheduleid"`
		Message     string `json:"message"`
	}

	auth := AuthResult{UserId: userid, ShippingScheduleId: value.Id, Message: bytes.NewBuffer(data).String()}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("Insert shipping schedule %v, response %s", value, buffer.String())
}

// ==================================================================================
// UpdateShippingSchedule -
// route - /resource/shippingschedule
// method - PUT
// ==================================================================================
func (s *ServerContainerREST) UpdateShippingSchedule(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	// Parse out the user enrollment ID
	userid := req.Header.Get("userid")
	token := req.Header.Get("token")
	sessionid := req.Header.Get("sessionid")

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
		encoder.Encode(restResult{Error: "Error trying to update an empty shipping schedule into the ledger"})
		restLogger.Error("Error trying to update an empty shipping schedule into the ledger")
		return
	}

	var value ShippingSchedule
	err = json.Unmarshal(data, &value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Error trying to unmarshal input into shipping schedule"})
		restLogger.Error("Error trying to unmarshal input into shipping schedule")
		return
	}
	if value.OwnerId != userid {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: fmt.Sprintf("Unauthorized operation, shipping schedule %v not owned by user %s", value, userid)})
		restLogger.Errorf(fmt.Sprintf("Unauthorized operation, shipping schedule %v not owned by user %s", value, userid))
		return
	}

	argsBuffer := bytes.NewBuffer(data)
	restLogger.Infof("Info reqire for %s %s: %v",req.RoutePath(), req.Method, argsBuffer.String())

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "update_shippingschedule", argsBuffer.String())


	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress +"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error inserting shipping schedule into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error inserting container into the ledger: %s", bytes.NewBuffer(data).String())
		return
	}


	type AuthResult struct {
		UserId    string `json:"userid"`
		ShippingScheduleId string `json:"shippingscheduleid"`
		Message     string `json:"message"`
	}

	auth := AuthResult{UserId: userid, ShippingScheduleId: value.Id, Message: bytes.NewBuffer(data).String()}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("Update shipping schedule %v, response %s", value, buffer.String())
}

// ==================================================================================
// GetShippingScheduleById -
// route - /resource/shippingschedule/{shippingScheduleId}
// method - GET
// ==================================================================================
func (s *ServerContainerREST) GetShippingScheduleById(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_shippingschedule", id)


	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress +"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read shipping schedule %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read shipping schedule %s from the ledger: %s", id, bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read shipping schedule by id %s from the ledger: %s", id, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}

// ==================================================================================
// DeleteShippingSchedule -
// route - /resource/shippingschedule/{shippingScheduleId}
// method - DELETE
// ==================================================================================
func (s *ServerContainerREST) DeleteShippingSchedule(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "delete_shippingschedule", id, userid)


	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress +"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error delete shipping schedule %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error delete shipping schedule %s from the ledger: %s", id, bytes.NewBuffer(data).String())
		return
	}
	restLogger.Infof("Delte shipping schedule %s from the ledger: %s", id, bytes.NewBuffer(data).String())

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: bytes.NewBuffer(data).String()})
}

// ==================================================================================
// GetVehiclesByOwnerId -
// route - /resource/shippingschedule/findByOwnerId
// method - GET
// ==================================================================================
func (s *ServerContainerREST) FindShippingSchedulesByOwnerId(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	req.ParseForm()
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")
	id := req.FormValue("ownerid")

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session"})
		restLogger.Errorf("Invalid session")
		return
	}
	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_shippingschedules_by_ownerid", id)


	valueAsBytes, _ := json.Marshal(&ccReq)
	ccResp, err := http.Post(fabricPeerAddress +"/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read shippingschedules by ownerid %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read shipping schedules by ownerid %s from the ledger: %s", id, bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read shipping schedules by ownerid %s from the ledger: %s", id, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}
