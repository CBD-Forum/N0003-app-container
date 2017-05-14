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
package user

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"cbdforum/app-container/service/src/common"
)

type Session struct {
	UserId    string
	SessionId string
	expiredAt string
}

const (
	sqlAddSession = "INSERT INTO session(userid, sessionid, expiredAt) VALUES(?, ?, ?)"
	sqlGetSession = "SELECT userid, sessionid, expiredAt FROM session WHERE userid = ? and sessionid = ? and deleted = 0"
	sqlUpdateSession = "UPDATE session SET expiredat = ? WHERE sessionid = ? and deleted = 0"
	sqlDeleteSession = "UPDATE session SET deleted = 1 WHERE sessionid = ? and deleted = 0"

)

func (s *Session) IsExpired() bool {
	return common.IsTimeBefore(s.expiredAt, common.GetCurrentTime())
}

// AddSession: insert a new session into table session
func AddSession(db *sql.DB, userId, sessionid, expiredAt string)(error) {
	var err error
	var stmt *sql.Stmt
	if err := db.Ping(); err != nil {
		logger.Fatal(err.Error())
		return common.NewErrorMessageDB("addSession: %v", err)
	}
	stmt, err = db.Prepare(sqlAddSession)
	if err != nil {
		logger.Errorf("addSession: %v", err)
		return common.NewErrorMessageDB("addSession: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, sessionid, expiredAt)
	if err != nil {
		logger.Errorf("addSessioin: %v", err)
		return common.NewErrorMessageDB("addSession: %v", err)
	}
	return err
}

// GetSession - query userSession table by userID and sessionUUID
// Inputs -
// Returns -
func GetSession(db *sql.DB, userId string, sessionId string, session *Session)(error){
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatal(err.Error())
		return common.NewErrorMessageDB("getSession: %v", err)
	}

	stmt, err = db.Prepare(sqlGetSession)
	if err != nil {
		logger.Errorf("getSession: %v", err)
		return common.NewErrorMessageDB("getSession: %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(userId, sessionId).Scan(&session.UserId, &session.SessionId, &session.expiredAt); err != nil {
		logger.Errorf("getSession: %v", err)
		return common.NewErrorMessageDB("getSession: %v", err)
	}
	return nil
}

func UpdateSession(db *sql.DB, sessionid, expiredAt string)(error){
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatal(err.Error())
		return common.NewErrorMessageDB("updateSession: %v", err)
	}

	stmt, err = db.Prepare(sqlUpdateSession)
	if err != nil {
		logger.Errorf("updateSession: %v", err)
		return common.NewErrorMessageDB("updateSession: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(expiredAt, sessionid)
	if err != nil {
		logger.Errorf("Failed executing statement:  %v", err)
		return common.NewErrorMessageDB("updateSession: %v", err)
	}
	return nil
}

func DeleteSession(db *sql.DB, sessionId string)(error){
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatal("deleteSession: %v", err)
		return common.NewErrorMessageDB("deleteSession: %v", err)
	}

	stmt, err = db.Prepare(sqlDeleteSession)
	if err != nil {
		logger.Errorf("deleteSession: %v", err)
		return common.NewErrorMessageDB("deleteSession: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(sessionId)
	if err != nil {
		logger.Errorf("deleteSession:  %v", err)
		return common.NewErrorMessageDB("deleteSession: %v", err)
	}

	return nil
}


func IsSessionValid(db *sql.DB, userid, sessionid, token string) bool {
	if len(userid)==0 || len(sessionid)==0 || len(token)==0 {
		return false
	}

	var user *User = new(User)
	err := GetUserById(db, userid, user)
	if err != nil {
		return false
	}
	var session *Session = new(Session)
	err = GetSession(db, userid, sessionid, session)
	if err != nil {
		return false
	}
	if session.IsExpired() {
		return false
	}
	if token != common.ComputeSessionToken(userid, sessionid, user.Password) {
		return false
	}
	return true
}
