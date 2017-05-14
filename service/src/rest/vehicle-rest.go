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


type Driver struct {
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	DrivingLicense string `json:"drivinglicense"`
}

type Vehicle struct {
	ObjectType string `json:"docType"`
	Id        string `json:"id"`
	VehicleNo string `json:"vehicleno"`
	Driver    Driver `json:"driver"`
	Status    string `json:"status"`
	OwnerId   string `json:"ownerid"`
}


// ==================================================================================
// InsertVehicle -
// route - /resource/vehicle
// method - POST
// ==================================================================================
func (s *ServerContainerREST) InsertVehicle(rw web.ResponseWriter, req *web.Request) {
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
		encoder.Encode(restResult{Error: "Error trying to add an empty vehicle into the ledger"})
		restLogger.Error("Error trying to add an empty vehicle into the ledger")
		return
	}

	var value Vehicle
	err = json.Unmarshal(data, &value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Error trying to unmarshal input into vehicle"})
		restLogger.Error("Error trying to unmarshal input into vehicle")
		return
	}
	// generate vehicle id
	value.Id = common.GenerateUUID()
	data, err = json.Marshal(value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Error trying to unmarshal input into vehicle"})
		restLogger.Error("Error trying to unmarshal input into vehicle")
		return
	}

	argsBuffer := bytes.NewBuffer(data)
	restLogger.Infof("Info reqire for /vehicle POST: %v", argsBuffer.String())

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "insert_vehicle", argsBuffer.String())



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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error inserting vehicle into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error inserting vehicle into the ledger: %s", bytes.NewBuffer(data).String())
		return
	}


	type AuthResult struct {
		UserId    string `json:"userid"`
		VehicleId string `json:"vehicleid"`
		Message     string `json:"message"`
	}

	auth := AuthResult{UserId: userid, VehicleId: value.Id, Message: bytes.NewBuffer(data).String()}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("Insert vehicle %v, response %s", value, buffer.String())
}

// ==================================================================================
// UpdateVehicle -
// route - /resource/vehicle
// method - PUT
// ==================================================================================
func (s *ServerContainerREST) UpdateVehicle(rw web.ResponseWriter, req *web.Request) {
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
		encoder.Encode(restResult{Error: "Error trying to update an empty vehicle into the ledger"})
		restLogger.Error("Error trying to update an empty vehicle into the ledger")
		return
	}

	var value Vehicle
	err = json.Unmarshal(data, &value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Error trying to unmarshal input into vehicle"})
		restLogger.Error("Error trying to unmarshal input into vehicle")
		return
	}
	if value.OwnerId != userid {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: fmt.Sprintf("Unauthorized operation, vehicle %v not owned by user %s", value, userid)})
		restLogger.Errorf(fmt.Sprintf("Unauthorized operation, vehicle %v not owned by user %s", value, userid))
		return
	}

	argsBuffer := bytes.NewBuffer(data)
	restLogger.Infof("Info reqire for %s %s: %v",req.RoutePath(), req.Method, argsBuffer.String())

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "update_vehicle", argsBuffer.String())


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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error inserting vehicle into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error inserting vehicle into the ledger: %s", bytes.NewBuffer(data).String())
		return
	}


	type AuthResult struct {
		UserId    string `json:"userid"`
		VehicleId string `json:"vehicleid"`
		Message     string `json:"message"`
	}

	auth := AuthResult{UserId: userid, VehicleId: value.Id, Message: bytes.NewBuffer(data).String()}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("Update vehicle %v, response %s", value, buffer.String())
}

// ==================================================================================
// GetVehicleById -
// route - /resource/vehicle/{vehicleId}
// method - GET
// ==================================================================================
func (s *ServerContainerREST) GetVehicleById(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_vehicle", id)

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read vehicle %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read vehicle %s from the ledger: %s", id, bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read vehicle by id %s from the ledger: %s", id, rpcResponse.Result.Message)


	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}

// ==================================================================================
// DeleteVehicle -
// route - /resource/vehicle/{vehicleId}
// method - DELETE
// ==================================================================================
func (s *ServerContainerREST) DeleteVehicle(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "delete_vehicle", id, userid)


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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error delete vehicle %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error delete vehicle %s from the ledger: %s", id, bytes.NewBuffer(data).String())
		return
	}
	restLogger.Infof("Delte vehicle %s from the ledger: %s", id, bytes.NewBuffer(data).String())

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: bytes.NewBuffer(data).String()})
}

// ==================================================================================
// GetVehiclesByOwnerId -
// route - /resource/vehicle/findByOwnerId
// method - GET
// ==================================================================================
func (s *ServerContainerREST) FindVehiclesByOwnerId(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_vehicles_by_ownerid", id)

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read vehicles by ownerid %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read vehicles by ownerid %s from the ledger: %s", id, bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read vehicles by ownerid %s from the ledger: %s", id, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}
