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
package order

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"cbdforum/app-container/chaincode/src/common"
)

func SendMessage(stub shim.ChaincodeStubInterface, userId string, fmtString string, a ...interface{}) error {
	// create message to notify the client
	var req MessageRequest = NewMessageRequest(userId, fmtString, a...)
	serverAddress := common.GetMessageServerAddress(stub)
	valueAsBytes, _ := json.Marshal(&req)
	resp, err := http.Post(serverAddress+"/message", "application/json", bytes.NewBuffer(valueAsBytes))
	if err != nil {
		logger.Errorf("Error trying to connect to message server %s: %v", serverAddress, err)
		return NewShimError(ERROR_INTERNAL, "Error trying to connect to message server %s: %v", serverAddress, err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Error notifying user with message %v, response is %s", req, bytes.NewBuffer(data).String())
		return NewOrderErrorMessage(ERROR_INTERNAL, "Error notifying user with message %v, response is %s", req, bytes.NewBuffer(data).String())
	}
	return nil
}
