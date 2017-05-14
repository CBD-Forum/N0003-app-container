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
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"time"
	"bytes"
)

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func SanitizeArguments(strs []string) error {
	for i, val := range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
	}
	return nil
}

func GetUint64Value(stub shim.ChaincodeStubInterface, key_value string) (uint64, error) {
	var t_val uint64
	t_bytes, err := stub.GetState(key_value)
	t_val = 0
	if err == nil {
		if t_bytes != nil {
			t_val, err = strconv.ParseUint(string(t_bytes), 10, 0)
			fmt.Println("in GetUint64Value:", t_val)
			return t_val, nil
		}
	}
	return t_val, err
}

func SetUint64Value(stub shim.ChaincodeStubInterface, key string, value uint64) error {
	var err error
	var bvalue string
	bvalue = strconv.FormatUint(value, 10)

	err = stub.PutState(key, []byte(bvalue))

	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func GetCurrentTime() string {
	return time.Now().Format(GOLANG_TIME_FMT_STRING)
}


func GetMessageServerAddress(stub shim.ChaincodeStubInterface) string {
	valueAsBytes, err := stub.GetState(MessageServer)
	if (err != nil) {
		return "0.0.0.0:9090"
	}
	return bytes.NewBuffer(valueAsBytes).String()
}
