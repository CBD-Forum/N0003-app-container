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
	"cbdforum/app-container/service/src/common"
	"cbdforum/app-container/service/src/user"

)



// ==================================================================================
// LogIn -
// route - /session/login
// method - POST
// ==================================================================================
func (s *ServerContainerREST) LogIn(rw web.ResponseWriter, req *web.Request) {
	// Parse out the user enrollment ID
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	encoder := json.NewEncoder(rw)

	if len(username)==0 || len(password)==0 {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error providing empty username or password")})
		restLogger.Error("Error providing empty username or password")
		return
	}

	db := common.GetDBInstance()
	if !user.IsUserNameExist(db, username) {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to login an unexisting user %s", username)})
		restLogger.Errorf("Error trying to login an unexisting user %s", username)
		return
	}

	var account *user.User = new(user.User)
	err := user.GetUserByName(db, username, account)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to connect %s: %s", common.GetDatabaseName(), err)})
		restLogger.Errorf("Error trying to connect %s: %s", common.GetDatabaseName(), err)
		return
	}
	if account.Password != password {
		rw.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(restResult{Error: "Error incorrect passworrd"})
		restLogger.Error("Error incorrect password")
		return
	}

	// Generate a new session
	sessionid := common.GenerateUUID()
	expiredAt := common.GenSessionExpireTime()
	err = user.AddSession(db, account.UserId, sessionid, expiredAt)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: "Error trying to add a new session"})
		return
	}

	token := common.ComputeSessionToken(account.UserId, sessionid, password)

	type AuthResult struct {
		UserId    string `json:"userid"`
		SessionId string `json:"sessionid"`
		Token     string `json:"token"`
		ExpiredAt string `json:"expiredat"`
	}

	auth := AuthResult{UserId: account.UserId, SessionId: sessionid, Token: token, ExpiredAt: expiredAt}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)
	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK: buffer.String()})
}

// ==================================================================================
// LogOut -
// route - /session/logout
// method - POST
// ==================================================================================
func (s *ServerContainerREST) LogOut(rw web.ResponseWriter, req *web.Request) {
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

	err := user.DeleteSession(db, sessionid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to logout, can't delte the session %s", sessionid)})
		return

	}

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK:fmt.Sprintf("User %s log out the system)", userid)})
}

// ==================================================================================
// Refresh -
// route - /session/refresh
// method - POST
// ==================================================================================
func (s *ServerContainerREST) Refresh(rw web.ResponseWriter, req *web.Request) {
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
	refreshedExpiredTime := common.GenSessionExpireTime()
	err := user.UpdateSession(db, sessionid, refreshedExpiredTime)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(restResult{Error: fmt.Sprintf("Error trying to refresh session, can't delte the session %s", sessionid)})
		return

	}

	type AuthResult struct {
		UserId    string `json:"userid"`
		SessionId string `json:"sessionid"`
		Token     string `json:"token"`
		ExpiredAt string `json:"refreshedexpiredat"`
	}

	auth := AuthResult{UserId: userid, SessionId: sessionid, Token: token, ExpiredAt: refreshedExpiredTime}
	authAsBytes, _ := json.Marshal(&auth)
	buffer := bytes.NewBuffer(authAsBytes)

	rw.WriteHeader(http.StatusOK)
	encoder.Encode(restResult{OK:buffer.String()})
}


