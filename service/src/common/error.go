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

import "fmt"

const (
	ERROR_DB = "Error database"
	ERROR_CHAINCODE = "Error chiancode"
	ERROR_INTERNAL = "Error internal"
)

type ErrorMessage struct {
	Type string
	Detail string
}

func (e ErrorMessage) String() string {
	if len(e.Detail) != 0 {
		return e.Type + " - " + e.Detail
	}
	return e.Type
}

func (e ErrorMessage) Error() string {
	return e.Type + " - " + e.Detail
}

func NewErrorMessage(errType string, fmtString string, a ...interface{}) error{
	return ErrorMessage{
		Type: errType,
		Detail: fmt.Sprintf(fmtString, a...),
	}
}

func NewErrorMessageDB(fmtString string, a ...interface{}) error {
	return ErrorMessage{
		Type: ERROR_DB,
		Detail: fmt.Sprintf(fmtString, a...),
	}
}

