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
	"strings"
)

type Container struct {
	ObjectType  string `json:"docType"` //field for couchdb
	Id          string `json:"id"`
	ContainerNo string `json:"containerno"`
	Type        string `json:"type"`
	MaxWeight   uint64 `json:"maxweight"`
	TareWeight  uint64 `json:"tareweight"`
	Measurement uint64 `json:"measurement"`
	Location    string `json:"location"`
	Status      string `json:"status"`
	OwnerId     string `json:"ownerid"`
}

// ==================================================================================
// InsertContainer -
// route - /resource/container
// method - POST
// ==================================================================================
func (s *ServerContainerREST) InsertContainer(rw web.ResponseWriter, req *web.Request) {
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
		encoder.Encode(restResult{Error: "Error trying to add an empty container into the ledger"})
		restLogger.Error("Error trying to add an empty container into the ledger")
		return
	}

	var value Container
	err = json.Unmarshal(data, &value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to unmarshal input %s into container: %v", bytes.NewBuffer(data).String(), err)})
		restLogger.Errorf("Error trying to unmarshal input %s into container: %v", err)
		return
	}
	// generate container id
	value.Id = common.GenerateUUID()
	data, err = json.Marshal(value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to marshal container %v into json bytes: %v", value, err)})
		restLogger.Errorf("Error trying to marshal container%v into json bytes", value, err)
		return
	}

	argsBuffer := bytes.NewBuffer(data)
	restLogger.Infof("Info reqire for /container POST: %v", argsBuffer.String())
	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "insert_container", argsBuffer.String())

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error inserting container into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error inserting container into the ledger: %s", bytes.NewBuffer(data).String())
		return
	}


	type AuthResult struct {
		UserId    string `json:"userid"`
		ContainerId string `json:"containerid"`
		Message     string `json:"message"`
	}

	auth := AuthResult{UserId: userid, ContainerId: value.Id, Message: bytes.NewBuffer(data).String()}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("Insert container %v, response %s", value, buffer.String())
}

// ==================================================================================
// UpdateContainer -
// route - /resource/container
// method - PUT
// ==================================================================================
func (s *ServerContainerREST) UpdateContainer(rw web.ResponseWriter, req *web.Request) {
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
		encoder.Encode(restResult{Error: "Error trying to update an empty container into the ledger"})
		restLogger.Error("Error trying to update an empty container into the ledger")
		return
	}

	var value Container
	err = json.Unmarshal(data, &value)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Error trying to unmarshal input into container"})
		restLogger.Error("Error trying to unmarshal input into container")
		return
	}
	if value.OwnerId != userid {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: fmt.Sprintf("Unauthorized operation, container %v not owned by user %s", value, userid)})
		restLogger.Errorf(fmt.Sprintf("Unauthorized operation, container %v not owned by user %s", value, userid))
		return
	}

	argsBuffer := bytes.NewBuffer(data)
	restLogger.Infof("Info reqire for %s %s: %v",req.RoutePath(), req.Method, argsBuffer.String())

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "update_container", argsBuffer.String())

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error inserting container into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error inserting container into the ledger: %s", bytes.NewBuffer(data).String())
		return
	}


	type AuthResult struct {
		UserId    string `json:"userid"`
		ContainerId string `json:"containerid"`
		Message     string `json:"message"`
	}

	auth := AuthResult{UserId: userid, ContainerId: value.Id, Message: bytes.NewBuffer(data).String()}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("Update container %v, response %s", value, buffer.String())
}

// ==================================================================================
// GetContainerById -
// route - /resource/container/{containerId}
// method - GET
// ==================================================================================
func (s *ServerContainerREST) GetContainerById(rw web.ResponseWriter, req *web.Request) {
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

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_container", id)

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read container %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read container %s from the ledger: %s", id, bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read container by id %s from the ledger: %s", id, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}

// ==================================================================================
// DeleteContainer -
// route - /resource/vehicle/{containerId}
// method - DELETE
// ==================================================================================
func (s *ServerContainerREST) DeleteContainer(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")
	id := req.PathParams["id"]

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	/*if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to connect %s: %s", common.DBName, err)})
		restLogger.Errorf("Error trying to connect %s: %s", common.DBName, err)
		return
	}*/
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session"})
		restLogger.Errorf("Invalid session")
		return
	}

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "delete_container", id, userid)

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error delete container %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error delete container %s from the ledger: %s", id, bytes.NewBuffer(data).String())
		return
	}
	restLogger.Infof("Delte container %s from the ledger: %s", id, bytes.NewBuffer(data).String())

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: bytes.NewBuffer(data).String()})
}

// ==================================================================================
// GetContainersByOwnerId -
// route - /resource/container/findByOwnerId
// method - GET
// ==================================================================================
func (s *ServerContainerREST) FindContainersByOwnerId(rw web.ResponseWriter, req *web.Request) {
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
	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_containers_by_ownerid", id)

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read container by ownerid %s from the ledger: %s", id, bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error read containers by ownerid %s from the ledger: %s", id, bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read containers by ownerid %s from the ledger: %s", id, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}



// ==================================================================================
// GetContainersHistory -
// route - /resource/container/track
// method - POST
// ==================================================================================
func (s *ServerContainerREST) TrackContainers(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	req.ParseForm()
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

	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil || data == nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: "Error trying to read container list"})
		restLogger.Error("Error trying to read container list")
		return
	}

	var containerNoList []string
	err = json.Unmarshal(data, &containerNoList)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Error trying to unmarshal input into containerNo list"})
		restLogger.Error("Error trying to unmarshal input into containerNo list")
		return
	}

	var containerTrackInfos []common.ContainerTrackInfo

	for _, item := range containerNoList {
		func() {
			req, _ := http.NewRequest(http.MethodGet, iotApiUrl + "/containers/" + item, strings.NewReader(""))
			req.Header.Add("Accept", "application/json")
			req.Header.Add("Authorization", "Bearer c6cce6c495d2a17043d972d25066381d")
			client := new(http.Client)
			iotResp, err := client.Do(req)
			if err != nil {
				restLogger.Warningf("Error trying to connect to iot to get container status: %v", err)
				return
			}

			data, err := ioutil.ReadAll(iotResp.Body)
			if err != nil || data == nil {
				restLogger.Warningf("Error trying to read response from iot: %v", err)
				return
			}
			defer iotResp.Body.Close()

			var value []common.ContainerTrackInfo = make([]common.ContainerTrackInfo, 0, 3000)
			err = json.Unmarshal(data, &value)
			if err != nil {
				restLogger.Warningf("Error read container status info by container number %s from the iot: %s", item, bytes.NewBuffer(data).String())
				return
			}
			if len(value) != 0 {
				containerTrackInfos = append(containerTrackInfos, value[0])
			}

		}()
	}

	valueAsBytes, err := json.Marshal(containerTrackInfos)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Failed to marshal response from iot: %v", err)})
		restLogger.Errorf("Failed to marshal response from iot: %v", err)
		return
	}

	restLogger.Infof("Track containers by containerno %v from the iot: %v", containerNoList, containerTrackInfos)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: bytes.NewBuffer(valueAsBytes).String()})
}
