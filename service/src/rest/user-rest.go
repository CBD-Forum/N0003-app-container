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
	"cbdforum/app-container/service/src/common"
	"cbdforum/app-container/service/src/user"
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"io/ioutil"
	"net/http"
)

type Person struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhoneNum string `json:"phonenum"`
	Address  string `json:"fax"`
	Fax      string
}

type Company struct {
	OrgId      string `json:"orgid"`
	TaxId      string `json:"taxid"`
	CreditId   string `json:"creditid"`
	BusinessId string `json:"businessid"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Homepage   string `json:"homepage"`
}

type RoleType string

const (
	ROLE_UNDEFINED      RoleType = "undefined"
	ROLE_REGULAR_CLIENT          = "regularclient"
	ROLE_CARGO_AGENT             = "cargoagent"
	ROLE_CARRIER                 = "carrier"
	ROLE_SHIPPER                 = "shipper"
)

type CCUser struct {
	ObjectType   string   `json:"docType"`
	Id           string   `json:"id"`
	UserName     string   `json:"username"`
	Password     string    `json:"password"`
	PersonalInfo Person   `json:"personalinfo"`
	Company      Company  `json:"company"`
	Role         RoleType `json:"role"`
	CreatedAt    string   `json:"createdat"`
	UpdatedAt    string   `json:"updatedat"`
	DeletedAt    string   `json:"deltedat"`
}

// ==================================================================================
// RegisterUser -
// route - /user
// method - POST
// ==================================================================================
func (s *ServerContainerREST) RegisterUser(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	encoder := json.NewEncoder(rw)

	//register the user into the ledger
	data, err := ioutil.ReadAll(req.Body)
	if err != nil || data == nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: "Error trying to register an empty user into the ledger"})
		restLogger.Error("Error trying to register an empty user into the ledger")
		return
	}
	defer req.Body.Close()


	// rewrite the request
	var ccuser CCUser
	err = json.Unmarshal(data, &ccuser)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to unmarshal request: %v", err)})
		restLogger.Errorf("Error trying to unmarshal request: %v", err)
		return
	}

	restLogger.Infof("Request for regsitering a user: %s", bytes.NewBuffer(data).String())

	if len(ccuser.UserName) == 0 || len(ccuser.Password) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to register a user with empty username or password: %+v", ccuser)})
		restLogger.Errorf("Error trying to register a user with empty username or password: %v", ccuser)
		return
	}

	db := common.GetDBInstance()
	if user.IsUserNameExist(db, ccuser.UserName) {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to register an existing user %s", ccuser.UserName)})
		restLogger.Errorf("Error trying to register an existing user %v", ccuser)
		return
	}

	ccuser.Id = common.GenerateUUID()
	valueAsBytes, _ := json.Marshal(ccuser)
	argsBuffer := bytes.NewBuffer(valueAsBytes)

	restLogger.Infof("Info require for %s %s: %v", req.RoutePath(), req.Method, argsBuffer.String())

	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "register_user", argsBuffer.String())

	valueAsBytes, _ = json.Marshal(&ccReq)
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
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error registering user into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error registering user into the ledger: %s", bytes.NewBuffer(data).String())
		return
	}

	// add user account into the mysql
	err = user.AddUser(db, ccuser.Id, ccuser.UserName, ccuser.Password)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to register user %s : %v", ccuser.UserName, err)})
		restLogger.Errorf("Error trying to register user %s : %v", ccuser.UserName, err)
		return
	}

	// Get a new session
	sessionid := common.GenerateUUID()
	expiredAt := common.AddTime(common.GetCurrentTime(), common.GetLocalSessionDuration())
	err = user.AddSession(db, ccuser.Id, sessionid, expiredAt)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Error trying to add a new session"})
		return
	}

	token := common.ComputeSessionToken(ccuser.Id, sessionid, ccuser.Password)

	type AuthResult struct {
		UserId    string `json:"userid"`
		SessionId string `json:"sessionid"`
		Token     string `json:"token"`
		ExpiredAt string `json:"expiredat"`
	}

	auth := AuthResult{UserId: ccuser.Id, SessionId: sessionid, Token: token, ExpiredAt: expiredAt}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
	restLogger.Infof("Register user %s, response %s", ccuser.UserName, buffer.String())
}

// todo: UpdateUser in the ledger
// ==================================================================================
// UpdateUser -
// route - /user
// method - PUT
// ==================================================================================
func (s *ServerContainerREST) UpdateUser(rw web.ResponseWriter, req *web.Request) {
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

	//todo: update user info in the ledger
	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil || data == nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: "Error trying to update an empty user into the ledger"})
		restLogger.Error("Error trying to update an empty user into the ledger")
		return
	}
	argsBuffer := bytes.NewBuffer(data)
	var ccReq CCRequest = NewCCRequest("invoke", fabricChaincodeName, "update_user", argsBuffer.String())

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error registering user into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error registering user into the ledger: %s", bytes.NewBuffer(data).String())
		return
	}

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: fmt.Sprintf("Update user info")})
}

// todo: GetUserById
// ==================================================================================
// GetUserById -
// route - /user/{userId}
// method - GET
// Returns:
// REST successfully query chaincode: {
// "jsonrpc":"2.0",
// "result":{
// 	"status":"OK",
// 	"message": "{
// 		\"docType\":\"user\",
// 		\"id\":\"415821d5-d3b0-444d-bd67-6418a0574ede\",
// 		\"username\":\"string\",
// 		\"personalinfo\":{
// 			\"name\":\"test\",
// 			\"email\":\"test@mail\",
// 			\"phonenum\":\"1234\",
// 			\"fax\":\"sdf\",
// 			\"Fax\":\"\"
// 		},
// 		\"company\":{
// 			\"orgid\":\"\",
// 			\"taxid\":\"\",
// 			\"creditid\":\"\",
// 			\"businessid\":\"\",
// 			\"name\":\"\",
//			\"address\":\"\",
// 			\"homepage\":\"\"
// 		},
// 		\"role\":\"regularclient\",
// 		\"createdat\":\"\",
// 		\"updatedat\":\"\",
// 		\"deltedat\":\"\"
// 		}"
// 	},
// 	"id":null
// }
// ==================================================================================
func (s *ServerContainerREST) GetUserById(rw web.ResponseWriter, req *web.Request) {
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

	// get user by id from fabric
	if userid != id {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Error trying to get other user's info"})
		restLogger.Error("Error trying to get other user's info")
		return
	}


	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "read_user", id)
	valueAsBytes, _ := json.Marshal(&ccReq)

	ccResp, err := http.Post(fabricPeerAddress + "/chaincode", "application/json", bytes.NewBuffer(valueAsBytes))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to connect to fabric peer: %v", err)})
		restLogger.Errorf("Error trying to connect to fabric peer: %v", err)
		return
	}
	defer ccResp.Body.Close()

	ccRespBody, err := ioutil.ReadAll(ccResp.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Internal Error when reading response body from peer: %v", err)})
		restLogger.Errorf("Internal Error when reading response body from peer: %v", err)
		return
	}

	if string(ccRespBody) == "" {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Response body from peer is empty"})
		restLogger.Error("Response body from peer is empty")
		return
	}
	var rpcResponse CCResponse
	err = json.Unmarshal(ccRespBody, &rpcResponse)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Failed to unmarshal response %s from peer: %v",string(ccRespBody), err)})
		restLogger.Errorf("Failed to unmarshal response %s from peer: %v",string(ccRespBody), err)
		return
	}

	if ccResp.StatusCode != http.StatusOK {
		rw.WriteHeader(ccResp.StatusCode)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error read user %s from the ledger: %s", id, rpcResponse.Result.Message)})
		restLogger.Errorf("Error read user %s from the ledger: %s", id, rpcResponse.Result.Message)
		return
	}
	restLogger.Errorf("Read user %s from the ledger: %s", id, rpcResponse.Result.Message)
/*
	var rpcResponse CCResponse
	err = json.Unmarshal(data, &rpcResponse)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Failed to unmarshal response from peer: %v", err)})
		restLogger.Errorf("Failed to unmarshal response from peer: %v", err)
		return
	}
*/
	restLogger.Infof("Info receive read_user response from fabric: %v", rpcResponse)

	restLogger.Infof("Read user by id %s from the ledger: %s", id, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}

// todo: FindUsersByType
// ==================================================================================
// FindUsersByRoleType -
// route - /user/finduserbytype/{userRoleType}
// method - GET
// ==================================================================================
func (s *ServerContainerREST) FindUsersByType(rw web.ResponseWriter, req *web.Request) {
	restLogger.Infof("Router: %s, method: %s", req.RoutePath(), req.Method)
	req.ParseForm()
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")
	roleType := req.FormValue("roletype")

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session"})
		restLogger.Errorf("Invalid session")
		return
	}

	var ccReq CCRequest = NewCCRequest("query", fabricChaincodeName, "find_users_by_role_type", roleType)

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
		encoder.Encode(restResult{Error: fmt.Sprintf("Error registering user into the ledger: %s", bytes.NewBuffer(data).String())})
		restLogger.Errorf("Error registering user into the ledger: %s", bytes.NewBuffer(data).String())
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

	restLogger.Infof("Read users by role type %s from the ledger: %s", roleType, rpcResponse.Result.Message)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: rpcResponse.Result.Message})
}
