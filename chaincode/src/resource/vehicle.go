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

type Driver struct {
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	DrivingLicense string `json:"drivinglicense"`
}

type Vehicle struct {
	ObjectType string `json:"docType"`
	Id        string `json:"id"`
	VehicleNo string `json:"vehicleno"`
	Driver    Driver `json:"driver"`
	Status    string `json:"status"`
	OwnerId   string `json:"ownerid"`
}

// ============================================================================================================================
// ReadVehicle - Read a vehicle from the ledger
// Input - Array of strings
// 	0,
// 	vehicleId,
// ============================================================================================================================
func ReadVehicle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	fmt.Println("starting readVehicle")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	id := args[0]

	// get the vehicle
	var value *Vehicle = new(Vehicle)
	valueAsBytes, err := stub.GetState(value.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                                 //this seems to always succeed, even if key didn't exist
		fmt.Printf("Failed to find vehicle %s ", id)
		return nil, err
	}

	fmt.Println("- end readVehicle")
	return valueAsBytes, nil
}

// ============================================================================================================================
// ReadAllVehiclesByOwnerId -
// Input - Array of strings
// 	0,
//	ownerId
// ============================================================================================================================

func ReadAllVehiclesByOwnerId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	var vehicles []Vehicle


	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}
	ownerId := args[0]

	// ---- Get All Vehicles ---- //
	indexName := common.INDEX_VEHICLE
	idList, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}

	for i, id := range idList {
		if i == 0 {
			continue
		}
		var value Vehicle
		queryValAsBytes, err := stub.GetState(value.BuildKey(id))
		if err != nil {
			//todo: value that have been deleted should update index
			continue
		}
		fmt.Println("on vehicle id - ", id)
		json.Unmarshal(queryValAsBytes, &value) //un stringify it aka JSON.parse()
		if value.OwnerId == ownerId {
			vehicles = append(vehicles, value)  //add this container to the list
		}

	}
	fmt.Println("vehicle array - ", vehicles)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(vehicles) //convert to array of bytes

	return everythingAsBytes, nil
}

// ============================================================================================================================
// InsertVehicle - Insert a new vehicle into the ledger
// inputs - Array of Strings
// 	0
// 	Vehicle
// ============================================================================================================================

func InsertVehicle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting insertVehicle")


	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	err = common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	var value *Vehicle = new(Vehicle)
	err = json.Unmarshal([]byte(args[0]), value)
	if err != nil {
		return nil, errors.New("Invalid type of arguments. Expecting Vehicle")
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

	if err = value.PutVehicle(stub); err != nil {
		return nil, err
	}

	indexName := common.INDEX_VEHICLE
	values, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}
	values = append(values, value.Id)
	err = common.PutStringList(stub, indexName, values)
	if err != nil {
		return nil, err
	}


	fmt.Println("- end insertVehicle")
	return nil, nil
}

// ============================================================================================================================
// UpdateVehicle - Update a existing vehicle in the ledger
// Inputs - Array of strings
// 	0,
// 	vehicle
// ============================================================================================================================
func UpdateVehicle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting updateContainer")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	if err = common.SanitizeArguments(args); err != nil {
		return nil, err
	}

	var value *Vehicle = new(Vehicle)
	if err = json.Unmarshal([]byte(args[0]), value); err != nil {
		return nil, errors.New("Invalid type of arguments. Expecting Vehicle")
	}

	// check whether container exist
	var old *Vehicle = new(Vehicle)
	if err = old.GetVehicle(stub, value.Id); err != nil {
		fmt.Println(value)
		return nil, errors.New("This vehicle does not exist - " + value.Id) //all stop a contianer by this id not exists
	}

	if err = value.PutVehicle(stub); err != nil {
		fmt.Println(value)
		return nil, err
	}


	fmt.Println("- end updateVehicle")
	return nil, nil
}

// ============================================================================================================================
// DeleteVehicle - Delete a vehicle from the ledger
// Inputs - Array of strings
// 	0, 		1
// 	vehicleid, 	ownerId
// ============================================================================================================================
func DeleteVehicle(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	fmt.Println("starting deleteVehicle")

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

	// get the vehicle
	var value *Vehicle = new(Vehicle)
	if err := value.GetVehicle(stub, id); err != nil {
		fmt.Println("Failed to find vehicle by id " + id)
		return nil, err
	}
	if value.OwnerId != ownerId {
		fmt.Printf("Owner %s is not authorized to delete vehicle %v", ownerId, value)
		return nil, errors.New(fmt.Sprintf("Owner %s is not authorized to delete vehicle %v", ownerId, value))
	}

	// remove the container
	err = stub.DelState(value.BuildKey(id)) //remove the key from chaincode state
	if err != nil {
		return nil, err
	}

	// todo: update the index

	fmt.Println("- end deleteVehicle")
	return nil, nil
}

// ============================================================================================================================
// Get Vehicle - get a container asset from ledger
// ============================================================================================================================
func (t *Vehicle) GetVehicle(stub shim.ChaincodeStubInterface, id string) error {
	containerAsBytes, err := stub.GetState(t.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                                 //this seems to always succeed, even if key didn't exist
		return errors.New("Failed to find vehicle - " + id)
	}
	err = json.Unmarshal(containerAsBytes, t) //un stringify it aka JSON.parse()
	if err != nil {
		return err
	}

	return nil
}

// ============================================================================================================================
// Put Vehicle - put a container asset into ledger
// ============================================================================================================================
func (t *Vehicle) PutVehicle(stub shim.ChaincodeStubInterface) error {
	valueAsBytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("vehicle: %#v", t)
		return errors.New("Failed to marshal this vehicle: " + err.Error())
	}

	err = stub.PutState(t.BuildKey(t.Id), valueAsBytes)
	if err != nil {
		return errors.New("Failed to put this vehicle into state: " + err.Error())
	}

	return nil
}

// ============================================================================================================================
// Build Key for a given Container Id
// ============================================================================================================================
func (t *Vehicle) BuildKey(id string) (key string) {
	return common.VehicleKeyPrefix + id
}

func (t *Vehicle) IsVehicleFree() bool {
	return t.Status == common.ResourceStatusFree
}

func (t *Vehicle) SetObjectType() {
	t.ObjectType = common.ObjectTypeVehicle
}

// ============================================================================================================================
// getVehicles:
// ============================================================================================================================
func getVehicles(stub shim.ChaincodeStubInterface, ownerId string) ([]Vehicle, error) {
	fmt.Println("Start to getVehicles")

	var vehicles []Vehicle

	// ---- Get All Vehicles ---- //
	indexName := common.INDEX_VEHICLE
	idList, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}

	for i, id := range idList {
		if i == 0 {
			continue
		}
		var value Vehicle
		queryValAsBytes, err := stub.GetState(value.BuildKey(id))
		if err != nil {
			//todo: value that have been deleted should update index
			continue
		}
		fmt.Println("on vehicle id - ", id)
		json.Unmarshal(queryValAsBytes, &value) //un stringify it aka JSON.parse()
		fmt.Printf("vehicle %d: %v\n", i, value)

		if value.OwnerId == ownerId && value.IsVehicleFree(){
			vehicles = append(vehicles, value)  //add this container to the list
		}
		fmt.Printf("vehicle %d: %v\n", i, value)

	}
	fmt.Println("vehicle array - ", vehicles)


	//change to array of bytes
	return vehicles, nil
}

func hasVehicles(stub shim.ChaincodeStubInterface, reqNum int, ownerId string) (bool, []Vehicle) {
	fmt.Println("Start to ask hasVehicles")
	vehicles, err := getVehicles(stub, ownerId)
	if err != nil || len(vehicles) < reqNum {
		return false, nil
	}
	fmt.Println(" - end to ask hasVehicles")
	return true, vehicles
}

func AllocateVehicles(stub shim.ChaincodeStubInterface, reqNum int, ownerId string) []Vehicle {
	fmt.Println("Start to allocate vehicles -")
	hasEnoughResources, resources := hasVehicles(stub, reqNum, ownerId)
	fmt.Printf("HasVehicles say isVehiclesEnough=%v, returns containers: %+v\n", hasEnoughResources, resources)
	if !hasEnoughResources {
		return nil
	}
	for _, item := range resources[0:reqNum] {
		item.Status = common.ResourceStatusInUse
	}
	fmt.Println(" - end allocate vehicles")
	return resources[0:reqNum]
}
