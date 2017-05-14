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
package resource

import (
	"cbdforum/app-container/chaincode/src/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type TransportTask struct {
	ObjectType   string    `json:"docType"` //field for couchdb
	Id           string    `json:"id"`
	ClientId     string    `json:"clientid"`
	CargoAgentId string    `json:"cargoagentid"`
	CarrierId    string    `json:"carrierid"`
	Vehicles     []Vehicle `json:"vehicles"`
	OrderId      string    `json:"orderid"`
	StartAt      string    `json:"startat"`
	EndAt        string    `json:"endat"`
	Status       string    `json:"status"`
}

// ============================================================================================================================
// ReadTransportTask - Read a transport task from the ledger
// Input - Array of strings
// 	0,
// 	taskId,
// ============================================================================================================================
func ReadTransportTask(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("starting readTransportTask")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	id := args[0]

	// get the vehicle
	var value *TransportTask = new(TransportTask)
	valueAsBytes, err := stub.GetState(value.BuildKey(id))
	if err != nil { //this seems to always succeed, even if key didn't exist
		fmt.Printf("Failed to find transport task %s\n", id)
		return nil, err
	}


	fmt.Println("- end readTransportTask")
	return valueAsBytes, nil
}

// ============================================================================================================================
// ReadAllTransportTasksByOwnerId -
// Input - Array of strings
// 	0,
//	ownerId
// Returns - Array of TransportTask
// [
// 	{ transportTask }, ...
// ]
// ============================================================================================================================

func ReadAllTransportTasksByOwnerId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var transportTasks []TransportTask

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}
	ownerId := args[0]

	// ---- Get All Containers ---- //
	indexName := common.INDEX_TRANSPORT_TASK
	idList, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}

	for i, id := range idList {
		if i == 0 {
			continue
		}
		var value TransportTask
		queryValAsBytes, err := stub.GetState(value.BuildKey(id))
		if err != nil {
			//todo: value that have been deleted should update index
			continue
		}
		fmt.Println("on transport task id - ", id)
		json.Unmarshal(queryValAsBytes, &value) //un stringify it aka JSON.parse()
		if value.CarrierId == ownerId {
			transportTasks = append(transportTasks, value)  //add this container to the list
		}

	}
	fmt.Println("transport task array - ", transportTasks)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(transportTasks) //convert to array of bytes
	return everythingAsBytes, nil
}

// ============================================================================================================================
// DeleteTransportTask -
// Inputs - Array of strings
// 	0, 			1
// 	transportTaskId,	ownerId
// ============================================================================================================================
func DeleteTransportTask(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	fmt.Println("starting deleteTransportTask")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err := common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	id := args[0]
	ownerId := args[1]

	// get the transport task
	var value *TransportTask = new(TransportTask)
	if err := value.GetTransportTask(stub, id); err != nil {
		fmt.Println("Failed to find transport task by id " + id)
		return nil, err
	}
	if value.CarrierId != ownerId {
		fmt.Printf("Owner %s is not authorized to delete transport task %v", ownerId, value)
		return nil, errors.New(fmt.Sprintf("Owner %s is not authorized to delete transport task %v", ownerId, value))
	}

	// remove the container
	err = stub.DelState(value.BuildKey(id)) //remove the key from chaincode state
	if err != nil {
		return nil, err
	}

	// todo: revise the index list common.INDEX_TRANSPORT_TASK

	fmt.Println("- end deleteTransportTask")
	return nil, nil
}

// ============================================================================================================================
// Get TransportTask - get a container asset from ledger
// ============================================================================================================================
func (t *TransportTask) GetTransportTask(stub shim.ChaincodeStubInterface, id string) error {
	valueAsBytes, err := stub.GetState(t.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                             //this seems to always succeed, even if key didn't exist
		return errors.New("Failed to find transportTask - " + id)
	}
	err = json.Unmarshal(valueAsBytes, t) //un stringify it aka JSON.parse()

	if err != nil {
		return err
	}

	return nil
}

// ============================================================================================================================
// Put TransportTask - put a transportTask asset into ledger
// ============================================================================================================================
func (t *TransportTask) PutTransportTask(stub shim.ChaincodeStubInterface) error {
	valueAsBytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("TransportTask: %#v", t)
		return errors.New("Failed to marshal this transportTask: " + err.Error())
	}

	err = stub.PutState(t.BuildKey(t.Id), valueAsBytes)
	if err != nil {
		return errors.New("Failed to put this transportTask into state: " + err.Error())
	}

	return nil
}

// ============================================================================================================================
// Build Key for a given Task Id
// ============================================================================================================================
func (t *TransportTask) BuildKey(id string) (key string) {
	return common.TransportTaskKeyPrefix + id
}

func (t *TransportTask) IsFinished() bool {
	return t.Status == common.StatusFailed || t.Status == common.StatusFinished
}


func (t *TransportTask) SetObjectType() {
	t.ObjectType = common.ObjectTypeTransportTask
}


/*type TransportTask struct {
	ObjectType string	`json:"docType"`	//field for couchdb
	Id string 		`json:"id"`
	ClientId string		`json:"clientid"`
	CargoAgentId string	`json:"cargoagentid"`
	CarrierId string	`json:"carrierid"`
	Vehicles []Vehicle	`json:"vehicles"`
	OrderId string		`json:"orderid"`
	StartAt string		`json:"startat"`
	EndAt string		`json:"endat"`
	Status string		`json:"status"`
}*/

func NewTransportTask(stub shim.ChaincodeStubInterface, carrierId, clientId, cargoAgentId, orderId string, vehicles []Vehicle, initStatus string) (*TransportTask, error) {
	var err error
	var value *TransportTask = new(TransportTask)
	fmt.Println("starting NewTransportTask")

	value.CarrierId = carrierId
	value.ClientId = clientId
	value.CargoAgentId = cargoAgentId
	value.OrderId = orderId

	value.Id = common.GenerateUUID()
	value.Vehicles = vehicles
	value.Status = initStatus
	value.SetObjectType()

	if err = value.PutTransportTask(stub); err != nil {
		return nil, err
	}
	fmt.Printf("Put the transport task %v into the ledger\n", value)

	indexName := common.INDEX_TRANSPORT_TASK
	values, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}
	values = append(values, value.Id)
	err = common.PutStringList(stub, indexName, values)
	if err != nil {
		return nil, err
	}


	fmt.Println("- end NewTransportTask")
	return value, nil
}

/*
// ============================================================================================================================
// arguments:
// 0,
// TransportTaskId, CurrentTransportTaskStatus, NextTransportTaskStatus, UserId
// ============================================================================================================================
func UpdateTransportTaskStatus(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting updateContainer")

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	if err = common.SanitizeArguments(args); err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	currentStatus := args[1]
	nextStatus := args[2]
	ownerId := args[3]


	// check whether transportTask exists
	var transportTask *TransportTask
	if err = transportTask.GetTransportTask(stub, ownerId, id); err != nil {
		return shim.Error(fmt.Sprintf("The transportTask %s owned by %s does not exist - ", ownerId, id))  //all stop a transportTask by this id not exists
	}


	// check current status
	if strings.Compare(transportTask.Status, currentStatus) != 0 {
		return shim.Error(fmt.Sprintf("The status of transportTask has changed, current status is %s, not %s", transportTask.Status, currentStatus))
	}

	transportTask.Status = nextStatus
	if err = transportTask.PutTransportTask(stub); err != nil {
		fmt.Println(transportTask)
		return shim.Error(err.Error())
	}

	fmt.Println("- end updateTransportTaskStatus")
	return shim.Success(nil)
}



// ============================================================================================================================
// Get history of asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  	0,			1
//  	id,			ownerId
// 	"01490985296352SjAyM"
// ============================================================================================================================
func getHistoryForTransportTask(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   TransportTask`json:"value"`
	}
	var history []AuditHistory;
	var transportTask TransportTask

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	transportTaskId := args[0]
	ownerId := args[1]
	fmt.Printf("- start getHistoryForTransportTask: %s\n", transportTaskId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(transportTask.BuildKey(ownerId, transportTaskId))
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		txID, historicValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = txID                             //copy transaction id over
		json.Unmarshal(historicValue, &transportTask)     //un stringify it aka JSON.parse()
		if historicValue == nil {                  //container has been deleted
			var emptyTransportTask TransportTask
			tx.Value = emptyTransportTask                 //copy nil container
		} else {
			json.Unmarshal(historicValue, &transportTask) //un stringify it aka JSON.parse()
			tx.Value = transportTask                      //copy container over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForTransportTask returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}
*/
