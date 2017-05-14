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

	"cbdforum/app-container/service/src/user"
	"cbdforum/app-container/service/src/common"
)


// ==================================================================================
// GetTransportTaskById -
// route - /resource/transporttask/{taskId}
// method - GET
// ==================================================================================
func (s *ServerContainerREST) GetTransportTaskById(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_transporttask", id)


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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read transporttask %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read transporttask %s from the ledger: %s", id, bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read transport task by id %s from the ledger: %s", id, rpcResponse.Result.Message)


	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}

// ==================================================================================
// DeleteTransportTask -
// route - /resource/transporttask/{taskId}
// method - DELETE
// ==================================================================================
func (s *ServerContainerREST) DeleteTransportTask(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "delete_transporttask", id, userid)
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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error delete transporttask %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error delete transporttask %s from the ledger: %s", id, bytes.NewBuffer(data).String())
		return
	}
	restLogger.Infof("Delte transporttask %s from the ledger: %s", id, bytes.NewBuffer(data).String())

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: bytes.NewBuffer(data).String()})
}

// ==================================================================================
// FindTransportTaskByOwnerId -
// route - /resource/transporttask/findByOwnerId
// method - GET
// ==================================================================================
func (s *ServerContainerREST) FindTransportTasksByOwnerId(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_transporttasks_by_ownerid", id)

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read transporttasks by ownerid %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read transporttasks by ownerid %s from the ledger: %s", id, bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read transporttasks by ownerid %s from the ledger: %s", id, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}
