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


type User struct {
	UserId   string
	Password string
	Username string
}


const (
	sqlAddUser = "INSERT INTO user(userid, username, password) VALUES(?, ?, ?)"
	sqlGetUserByName = "SELECT userid, username, password FROM user WHERE username = ? and deleted = 0"
	sqlGetUserById = "SELECT userid, username, password FROM user WHERE userid = ? and deleted = 0"
	sqlDeleteUser = "UPDATE user SET deleted = 1 WHERE userid = ? and deleted = 0"
)


func IsUserNameExist(db *sql.DB, name string) bool {
	var user *User = new(User)
	err := GetUserByName(db, name, user)
	if err != nil {
		return false
	}
	return true
}

func IsUserIdExist(db *sql.DB, userid string) bool {
	var user *User = new(User)
	err := GetUserById(db, userid, user)
	if err != nil {
		return false
	}
	return true
}

func AddUser(db *sql.DB, userId, username, password string) (error) {
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatalf("addUser: %v", err)
		return common.NewErrorMessageDB("addUser: %v", err)
	}

	stmt, err = db.Prepare(sqlAddUser)
	if err != nil {
		logger.Errorf("addUser: %v", err)
		return common.NewErrorMessageDB("addUser: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, username, password)
	if err != nil {
		logger.Errorf("addUser:  %v", err)
		return common.NewErrorMessageDB("addUser: %v", err)
	}

	return nil

}

func GetUserByName(db *sql.DB, name string, user *User) (error) {
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatalf("getUserByName: %v", err)
		return common.NewErrorMessageDB("getUserByName: %v", err)
	}

	stmt, err = db.Prepare(sqlGetUserByName)
	if err != nil {
		logger.Errorf("getUserByName: %v", err)
		return common.NewErrorMessageDB("getUserByName: %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(name).Scan(&(user.UserId), &(user.Username), &(user.Password)); err != nil {
		logger.Errorf("getUserByName: failed to get user %s, error %v", name, err)
		return common.NewErrorMessageDB("getUserByName: failed to get user %s, error %v", name, err)
	}
	logger.Debugf("Get user by name %s: \n%+v", name, *user)

	return nil
}

func GetUserById(db *sql.DB, userId string, user *User) (error) {
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatalf("getUserById: %v", err)
		return common.NewErrorMessageDB("getUserById: %v", err)
	}

	stmt, err = db.Prepare(sqlGetUserById)
	if err != nil {
		logger.Errorf("getUserById: %v", err)
		return common.NewErrorMessageDB("getUserById: %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(userId).Scan(&(user.UserId), &(user.Username), &(user.Password)); err != nil {
		logger.Errorf("getUserById: failed to get user %s, error %v", userId, err)
		return common.NewErrorMessageDB("getUserById: failed to get user %s, error %v", userId, err)
	}
	logger.Debugf("Get user %s: \n%+v", userId, *user)

	return nil
}

func DeleteUser(db *sql.DB, userId string) (error) {
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatalf("deleteUser: %v", err)
		return common.NewErrorMessageDB("deleteUser: %v", err)
	}

	stmt, err = db.Prepare(sqlDeleteUser)
	if err != nil {
		logger.Errorf("deleteUser: %v", err)
		return common.NewErrorMessageDB("deleteUser: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		logger.Errorf("deleteUser: %v", err)
		return common.NewErrorMessageDB("deleteUser: %v", err)
	}

	return nil
}
