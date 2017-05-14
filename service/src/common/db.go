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
package common

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var db *sql.DB
var once sync.Once
func GetDBInstance() (*sql.DB) {
	once.Do(func(){
		var err error
		db, err = sql.Open(databaseName, databaseDsn)
		if err != nil {
			panic(fmt.Sprintf("Failed to connect database %s in %s: %v", databaseName, databaseDsn, err))
		}
	})
	return db
}


func GetDatabaseName() string {
	return databaseName
}

func GetDatabaseDsn() string {
	return databaseDsn
}
