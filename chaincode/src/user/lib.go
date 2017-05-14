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
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"cbdforum/app-container/chaincode/src/common"
	"strings"
)

func getUsersByType(stub shim.ChaincodeStubInterface, role string) ([]User, error) {
	var err error
	var users []User

	// ---- Get All Users ---- //
	indexName := common.INDEX_USER
	valueList, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}

	for i, id := range valueList {
		if i == 0 {
			continue
		}
		var user User
		queryValAsBytes, err := stub.GetState(user.BuildKey(id))
		if err != nil {
			return nil, err
		}

		fmt.Println("on user id - ", id)
		json.Unmarshal(queryValAsBytes, &user) //un stringify it aka JSON.parse()
		if strings.Compare(user.Role, role) == 0 {
			users = append(users, user) //add this user to the list
		}
	}
	return users, nil
}
