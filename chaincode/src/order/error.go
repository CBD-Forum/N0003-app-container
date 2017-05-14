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
	"fmt"
)

type OrderErrorMessage struct {
	errorType string
	detail    string
}

const (
	ERROR_REQUEST   string = "Error request"
	ERROR_RESOURCE         = "Error resource"
	ERROR_INTERNAL         = "Error internal"
	ERROR_FSM              = "Error finite state machine"
	ERROR_ARGUMENTS        = "Error arguments"
)

func NewShimError(errorType string, fmtString string, a ...interface{}) error {
	return NewOrderErrorMessage(errorType, fmtString, a...)
}

func NewOrderErrorMessage(errorType string, fmtString string, a ...interface{}) OrderErrorMessage {
	var t OrderErrorMessage
	t.errorType = errorType
	t.detail = fmt.Sprintf(fmtString, a...)
	return t
}
func (t OrderErrorMessage) String() string {
	return t.errorType + ": " + t.detail
}

func (t OrderErrorMessage) Error() string {
	return t.errorType + ": " + t.detail
}

func (t OrderErrorMessage) Sprintf(fmtString string, a ...interface{}) string {
	t.detail = fmt.Sprintf(fmtString, a...)

	return t.errorType + ": " + t.detail
}
