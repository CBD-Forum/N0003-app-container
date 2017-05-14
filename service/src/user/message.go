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
	"cbdforum/app-container/service/src/common"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sort"
)

const (
	MessageStatusRead = "read"
	MessageStatusUnread = "unread"

	SQL_AddMessage = "INSERT INTO message(userid, messageid, message, status, createdat) VALUES(?, ?, ?, ?, ?)"
	SQL_GetMessagesByUserId = "SELECT userid, messageid, message, status, createdat FROM message WHERE userid = ? and deleted = 0"
	SQL_GetMessageById = "SELECT userid, messageid, message, status, createdat FROM message WHERE messageid = ? and deleted = 0"
	SQL_DeleteMessage = "UPDATE message SET deleted = 1 WHERE messageid = ? and deleted = 0"
	SQL_UpdateMessage = "UPDATE message SET status = ? WHERE messageid = ? and deleted = 0"
)

type Message struct {
	UserId    string	`json:"userid"`
	MessageId string	`json:"messageid"`
	Message   string	`json:"message"`
	Status    string	`json:"status"`
	CreatedAt string	`json:"createdat"`
}

type MessageList []Message

func IsMessageIdExist(db *sql.DB, messageid string) bool {
	var msg *Message = new(Message)
	err := GetMessageById(db, messageid, msg)
	if err != nil {
		return false
	}
	return true
}



func AddMessage(db *sql.DB, userid, messageid, message, createdAt string) error {
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatalf("addMessage: %v", err)
		return common.NewErrorMessageDB("addMessage: %v", err)
	}

	stmt, err = db.Prepare(SQL_AddMessage)
	if err != nil {
		logger.Errorf("addMessage: %v", err)
		return common.NewErrorMessageDB("addMessage: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userid, messageid, message, MessageStatusUnread, createdAt)
	if err != nil {
		logger.Errorf("addMessage:  %v", err)
		return common.NewErrorMessageDB("addMessage: %v", err)
	}

	return nil

}


func UpdateMessage(db *sql.DB, messageid, status string)(error){
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatal(err.Error())
		return common.NewErrorMessageDB("updateMessage: %v", err)
	}

	stmt, err = db.Prepare(SQL_UpdateMessage)
	if err != nil {
		logger.Errorf("updateMessage: %v", err)
		return common.NewErrorMessageDB("updateMessage: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, messageid)
	if err != nil {
		logger.Errorf("Failed executing statement:  %v", err)
		return common.NewErrorMessageDB("udpateMessage: %v", err)
	}
	return nil
}


// Notice: sort messages by order unread, latest > unread, older > read
func GetMessagesByUserId(db *sql.DB, userId string) (error,  []Message) {
	var err error
	var rows *sql.Rows
	var messages []Message

	if err = db.Ping(); err != nil {
		logger.Fatalf("getMessagesByUserId: %v", err)
		return common.NewErrorMessageDB("getMessagesByUserId: %v", err), nil
	}

	rows, err = db.Query(SQL_GetMessagesByUserId, userId)
	if err != nil {
		logger.Errorf("getMessagesByUserId: %v", err)
		return common.NewErrorMessageDB("getMessagesByUserId: %v", err), nil
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		err = rows.Scan(&(msg.UserId), &(msg.MessageId), &(msg.Message), &(msg.Status), &(msg.CreatedAt))
		if err == nil {
			messages = append(messages, msg)
		}
	}

	err = rows.Err()
	if err != nil {
		return common.NewErrorMessageDB("getMessagesByUserId: %v", err), nil
	}

	sort.Sort(MessageList(messages))

	logger.Debugf("Get messages by userid %s: \n%+v", userId, messages)
	return nil, messages
}

func GetMessageById(db *sql.DB, messageId string, msg *Message) (error){
	var err error
	var stmt *sql.Stmt
	if err = db.Ping(); err != nil {
		logger.Fatalf("getMessageById: %v", err)
		return common.NewErrorMessageDB("getMessageById: %v", err)
	}

	stmt, err = db.Prepare(SQL_GetMessageById)
	if err != nil {
		logger.Errorf("getMessageById: %v", err)
		return common.NewErrorMessageDB("getMessageById: %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(messageId).Scan(&(msg.UserId), &(msg.MessageId), &(msg.Message), &(msg.Status), &(msg.CreatedAt)); err != nil {
		logger.Errorf("getMessageById: failed to  get message %s, error %v", messageId, err)
		return common.NewErrorMessageDB("getMessageById: failed to get message %s, error %v", messageId, err)
	}
	logger.Debugf("getMessageById: message %+v", *msg)

	return nil
}

func DeleteMessage(db *sql.DB, messageId string) error {
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		logger.Fatalf("deleteMessage: %v", err)
		return common.NewErrorMessageDB("deleteMessage: %v", err)
	}

	stmt, err = db.Prepare(SQL_DeleteMessage)
	if err != nil {
		logger.Errorf("deleteMessage: %v", err)
		return common.NewErrorMessageDB("deleteMessage: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(messageId)
	if err != nil {
		logger.Errorf("deleteMessage: %v", err)
		return common.NewErrorMessageDB("deleteMessage: %v", err)
	}

	return nil
}


/*
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}
*/
func (m MessageList) Len() int {
	return len(m)
}
func (m MessageList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m MessageList) Less(i, j int) bool {
	const UNREAD = 0
	const READ = 1
	var statusI int
	var statusJ int
	if m[i].Status == MessageStatusUnread {
		statusI = UNREAD
	} else {
		statusI = READ
	}

	if m[j].Status == MessageStatusUnread {
		statusJ = UNREAD
	} else {
		statusJ = READ
	}

	if statusI != statusJ {
		return statusI < statusJ
	} else {
		return !common.IsTimeBefore(m[i].CreatedAt, m[j].CreatedAt)
	}
}
