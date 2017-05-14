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
	auth "cbdforum/app-container/chaincode/src/user"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Container struct {
	ObjectType  string `json:"docType"` //field for couchdb
	Id          string `json:"id"`
	ContainerNo string `json:"containerno"`
	Type        string `json:"type"`
	MaxWeight   uint64 `json:"maxweight"`
	TareWeight  uint64 `json:"tareweight"`
	Measurement uint64 `json:"measurement"`
	Location    string `json:"location"`
	Status      string `json:"status"`
	OwnerId     string `json:"ownerid"`
}

// ============================================================================================================================
// ReadContainer - Read a container from the ledger
// Inputs - Array of strings
// 	0
// 	containerId
// ============================================================================================================================
func ReadContainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("starting ReadContainer")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	id := args[0]

	// get the container
	var value *Container = new(Container)
	valueAsBytes, err := stub.GetState(value.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                                 //this seems to always succeed, even if key didn't exist
		fmt.Printf("Failed to find container %s ", id)
		return nil, err
	}

	fmt.Println("- end ReadContainer")
	return valueAsBytes, nil
}

// ============================================================================================================================
// ReadAllContainersByOwnerId - Read all containers by ownerId (containers)
//
// Inputs - Array of strings
//	0
//	userId
// ============================================================================================================================
func ReadAllContainersByOwnerId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var Containers []Container

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}
	ownerId := args[0]

	// ---- Get All Containers ---- //
	indexName := common.INDEX_CONTAINER
	idList, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}

	for i, id := range idList {
		if i == 0 {
			continue
		}
		var value Container
		queryValAsBytes, err := stub.GetState(value.BuildKey(id))
		if err != nil {
			//todo: value that have been deleted should update index
			continue
		}
		fmt.Println("on container id - ", id)
		json.Unmarshal(queryValAsBytes, &value) //un stringify it aka JSON.parse()
		if value.OwnerId == ownerId {
			Containers = append(Containers, value)  //add this container to the list
		}

	}
	fmt.Println("container array - ", Containers)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(Containers) //convert to array of bytes
	return everythingAsBytes, nil
}

// ============================================================================================================================
// InsertContainer - Insert a new container into the ledger
// inputs - Array of Strings
// 	0
// 	Container
// ============================================================================================================================
func InsertContainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting insertContainer")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	err = common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	var value *Container = new(Container)
	err = json.Unmarshal([]byte(args[0]), value)
	if err != nil {
		return nil, errors.New("Invalid type of arguments. Expecting Container")
	}
	if len(value.Id) == 0 {
		return nil, errors.New("Container Id should not be empty")
	}

	var owner *auth.User = new(auth.User)
	if !auth.IsUserExist(stub, value.OwnerId, owner) {
		return nil, errors.New("Owner of the container not exists.")
	}

	value.Status = common.ResourceStatusFree
	value.SetObjectType()

	if err = value.PutContainer(stub); err != nil {
		return nil, err
	}

	indexName := common.INDEX_CONTAINER
	values, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}
	values = append(values, value.Id)
	err = common.PutStringList(stub, indexName, values)
	if err != nil {
		return nil, err
	}

	fmt.Println("- end insertContainer")
	return nil, nil
}

// ============================================================================================================================
// UpdateContainer - update an existing container in the ledger
// inputs - Array of Strings
// 	0
// 	Container
// ============================================================================================================================
func UpdateContainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	var err error
	fmt.Println("starting updateContainer")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	if err = common.SanitizeArguments(args); err != nil {
		return nil, err
	}

	var value *Container = new(Container)
	if err = json.Unmarshal([]byte(args[0]), value); err != nil {
		return nil, errors.New("Invalid type of arguments. Expecting Container")
	}

	// check whether container exist
	var old *Container = new(Container)
	if err = old.GetContainer(stub, value.Id); err != nil {
		fmt.Println(value)
		return nil, errors.New("This container does not exist - " + value.Id) //all stop a contianer by this id not exists
	}

	if err = value.PutContainer(stub); err != nil {
		fmt.Println(value)
		return nil, err
	}

	fmt.Println("- end updateContainer")
	return nil, nil
}

// ============================================================================================================================
// DeleteContainer - Delete a container from the ledger
// Inputs - Array of strings
// 	0, 		1
// 	containerId, 	ownerId
// ============================================================================================================================
func DeleteContainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("starting deleteContainer")

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

	// get the container
	var value *Container = new(Container)
	if err := value.GetContainer(stub, id); err != nil {
		fmt.Println("Failed to find contianer by id " + id)
		return nil, err
	}
	if value.OwnerId != ownerId {
		fmt.Printf("Owner %s is not authorized to delete container %v", ownerId, value)
		return nil, errors.New(fmt.Sprintf("Owner %s is not authorized to delete container %v", ownerId, value))
	}

	// remove the container
	err = stub.DelState(value.BuildKey(id)) //remove the key from chaincode state
	if err != nil {
		return nil, err
	}

	fmt.Println("- end delete_container")
	return nil, nil
}

/* supported only by fabric-1.0.0-alpha, not supported by fabric-0.6
// ============================================================================================================================
// Get history of asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0,
//  id,
//  "01490985296352SjAyM",
// ============================================================================================================================
func GetHistoryForContainer(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	type AuditHistory struct {
		TxId  string    `json:"txId"`
		Value Container `json:"value"`
	}
	var history []AuditHistory
	var container Container

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	containerId := args[0]
	fmt.Printf("- start getHistoryForContainer, containerid is %s, ownerid is %s\n", containerId, ownerId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(container.BuildKey(ownerId, containerId))
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
		tx.TxId = txID                            //copy transaction id over
		json.Unmarshal(historicValue, &container) //un stringify it aka JSON.parse()
		if historicValue == nil {                 //container has been deleted
			var emptyContainer Container
			tx.Value = emptyContainer //copy nil container
		} else {
			json.Unmarshal(historicValue, &container) //un stringify it aka JSON.parse()
			tx.Value = container                      //copy container over
		}
		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistoryForContainer returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
}
*/

// ============================================================================================================================
// Get Container - get a container asset from ledger
// Inputs -
// ============================================================================================================================
func (t *Container) GetContainer(stub shim.ChaincodeStubInterface, id string) error {
	containerAsBytes, err := stub.GetState(t.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                                 //this seems to always succeed, even if key didn't exist
		return errors.New("Failed to find container - " + id)
	}
	err = json.Unmarshal(containerAsBytes, t) //un stringify it aka JSON.parse()
	if err != nil {
		return err
	}
	return nil
}

// ============================================================================================================================
// Put Container - put a container asset into ledger
// ============================================================================================================================
func (t *Container) PutContainer(stub shim.ChaincodeStubInterface) error {
	containerAsBytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("%#v", t)
		return errors.New("Failed to marshal this container: " + err.Error())
	}

	err = stub.PutState(t.BuildKey(t.Id), containerAsBytes)
	if err != nil {
		return errors.New("Failed to put this container into state: " + err.Error())
	}

	return nil
}

// ============================================================================================================================
// Build Key for a given Container Id
// ============================================================================================================================
func (t *Container) BuildKey(id string) (key string) {
	return common.ContainerKeyPrefix + id
}

func (t *Container) IsContainerFree() bool {
	return t.Status == common.ResourceStatusFree
}


func (t *Container) SetObjectType() {
	t.ObjectType = common.ObjectTypeContainer
}

// ============================================================================================================================
// GetContainers:
// ============================================================================================================================
func getContainers(stub shim.ChaincodeStubInterface, containerType, ownerId string) ([]Container, error) {
	fmt.Println("Start to getContainers")
	var containers []Container

	// ---- Get All Containers ---- //
	indexName := common.INDEX_CONTAINER
	idList, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}

	for i, id := range idList {
		if i == 0 {
			continue
		}
		var value Container
		queryValAsBytes, err := stub.GetState(value.BuildKey(id))
		if err != nil {
			return nil, err
		}
		fmt.Println("on container id - ", id)
		json.Unmarshal(queryValAsBytes, &value) //un stringify it aka JSON.parse()
		fmt.Printf("container is %+v\n", value)
		if value.OwnerId != ownerId {
			continue
		}
		if value.Type == containerType && value.IsContainerFree() {
			containers = append(containers, value) //add this container to the list
		}
		fmt.Printf("container %d: %v\n", i, value)

	}
	fmt.Println("container array - ", containers)

	//change to array of bytes
	return containers, nil
}

func HasContainers(stub shim.ChaincodeStubInterface, reqNum int, containerType, ownerId string) (bool, []Container) {
	fmt.Println("Start to ask hasContainers")
	var err error
	var containers []Container
	containers, err = getContainers(stub, containerType, ownerId)
	if err != nil || len(containers) < reqNum {
		return false, nil
	}
	fmt.Println(" - end to ask hasContainers")
	return true, containers[0:reqNum]
}

func AllocateContainers(stub shim.ChaincodeStubInterface, reqNum int, containerType, ownerId string) []Container {
	fmt.Println("Start to allocateContainers -")
	var containers []Container
	var isContainersEnough bool
	isContainersEnough, containers = HasContainers(stub, reqNum, containerType, ownerId)
	fmt.Printf("HasContainers say isContainersEnough=%v, returns containers: %+v\n", isContainersEnough, containers)
	if !isContainersEnough {
		return nil
	}
	for _, item := range containers {
		item.Status = common.ResourceStatusInUse
	}
	fmt.Println(" - end allocate Containers")
	return containers

}
