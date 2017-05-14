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
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gocraft/web"
	"bytes"
	"io/ioutil"
	"cbdforum/app-container/service/src/common"
	"cbdforum/app-container/service/src/user"

)


// ==================================================================================
// InsertMessage -
// route - /message
// method - POST
// ==================================================================================
func (s *ServerContainerREST) InsertMessage(rw web.ResponseWriter, req *web.Request) {

	encoder := json.NewEncoder(rw)
	data, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	type MessageRequest struct {
		UserId string `json:"userid"`
		Message string `json:"message"`
	}

	var messageRequest *MessageRequest = new(MessageRequest)
	err := json.Unmarshal(data, messageRequest)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Requrie message request as %T: %v", *messageRequest, err)})
		restLogger.Errorf("Require message request as %T: %v", *messageRequest, err)
		return
	}

	db := common.GetDBInstance()
	if !user.IsUserIdExist(db, messageRequest.UserId) {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("User %s not exist", messageRequest.UserId)})
		restLogger.Errorf("User %s not exist", messageRequest.UserId)
		return
	}

	messageId := common.GenerateUUID()
	createdAt := common.GetCurrentTime()
	err = user.AddMessage(db, messageRequest.UserId, messageId, messageRequest.Message, createdAt)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Failed to add a new message: %v", err)})
		restLogger.Errorf("Failed to add a new message: %v", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: messageId})
}


// ==================================================================================
// UpdateMessageStatus -
// route - /message
// method - PUT
// ==================================================================================
func (s *ServerContainerREST) UpdateMessageStatus(rw web.ResponseWriter, req *web.Request) {
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session, userid, sessoinid and token not set in the header or session has expired"})
		restLogger.Errorf("Invalid session")
		return
	}


	data, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	type MessageRequest struct {
		MessageId string `json:"messageid"`
		Status string `json:"status"`
	}

	var messageRequest *MessageRequest = new(MessageRequest)
	err := json.Unmarshal(data, messageRequest)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Requrie message request as %T: %v", *messageRequest, err)})
		restLogger.Errorf("Require message request as %T: %v", *messageRequest, err)
		return
	}



	if !user.IsMessageIdExist(db, messageRequest.MessageId) {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Message %s not exist", messageRequest.MessageId)})
		restLogger.Errorf("Message %s not exist", messageRequest.MessageId)
		return
	}

	if messageRequest.Status != user.MessageStatusRead && messageRequest.Status != user.MessageStatusUnread{
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Unsupported Status %s, only supported [%s, %s]", messageRequest.Status, user.MessageStatusUnread, user.MessageStatusRead)})
		restLogger.Errorf("Unsupported Status %s, only supported [%s, %s]", messageRequest.Status, user.MessageStatusUnread, user.MessageStatusRead)
		return
	}

	err = user.UpdateMessage(db, messageRequest.MessageId, messageRequest.Status)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Failed to update message status: %v", err)})
		restLogger.Errorf("Failed to update message status: %v", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: messageRequest.MessageId})
}


// ==================================================================================
// GetMessageById -
// route - /message/{messageId}
// method - GET
// ==================================================================================
func (s *ServerContainerREST) GetMessageById(rw web.ResponseWriter, req *web.Request) {
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")
	messageid := req.PathParams["id"]

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session, userid, sessoinid and token not set in the header or session has expired"})
		restLogger.Errorf("Invalid session")
		return
	}

	var message *user.Message = new(user.Message)
	err := user.GetMessageById(db, messageid, message)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to get message %s", messageid)})
		return

	}

	valueAsBytes, err := json.Marshal(message)
	buffer := bytes.NewBuffer(valueAsBytes)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK:buffer.String()})
}

// ==================================================================================
// DeleteMessage -
// route - /message/{messageId}
// method - DELETE
// ==================================================================================
func (s *ServerContainerREST) DeleteMessage(rw web.ResponseWriter, req *web.Request) {
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")
	messageid := req.PathParams["id"]

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session, userid, sessoinid and token not set in the header or session has expired"})
		restLogger.Errorf("Invalid session")
		return
	}

	if !user.IsMessageIdExist(db, messageid){
		rw.WriteHeader(http.StatusNotFound)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to delete unexisting message %s", messageid)})
		return

	}

	err := user.DeleteMessage(db, messageid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to delete message %s: %v", messageid, err)})
		return
	}

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK:fmt.Sprintf("User %s deleted message %s", userid, messageid)})
}

// ==================================================================================
// FindMessagesByUserId -
// route - /message/findByUserId
// method - GET
// ==================================================================================
func (s *ServerContainerREST) FindMessagesByUserId(rw web.ResponseWriter, req *web.Request) {
	userid := req.Header.Get("userid")
	sessionid := req.Header.Get("sessionid")
	token := req.Header.Get("token")

	encoder := json.NewEncoder(rw)

	db := common.GetDBInstance()
	if !user.IsSessionValid(db, userid, sessionid, token) {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Invalid session, userid, sessoinid and token not set in the header or session has expired"})
		restLogger.Errorf("Invalid session")
		return
	}

	err, messages := user.GetMessagesByUserId(db, userid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to get messages by userid %s: %v", userid, err)})
		restLogger.Errorf("Error trying to get messages by userid %s: %v", userid, err)
		return

	}

	valueAsBytes, err := json.Marshal(messages)
	buffer := bytes.NewBuffer(valueAsBytes)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK:buffer.String()})
}





