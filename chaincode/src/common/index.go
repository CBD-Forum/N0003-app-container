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
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func GetStringList(stub shim.ChaincodeStubInterface, key string) ([]string, error) {
	var value []string
	valueAsBytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if len(valueAsBytes) == 0 {
		return nil, nil
	}
	err = json.Unmarshal(valueAsBytes, &value)
	if err != nil {
		return nil, err
	}
	return value, err
}

func PutStringList(stub shim.ChaincodeStubInterface, key string, value []string) (error) {
	valueAsBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = stub.PutState(key, valueAsBytes)
	if err != nil {
		return err
	}
	return nil
}

