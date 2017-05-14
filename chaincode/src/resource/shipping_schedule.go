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

type ShippingSchedule struct {
	ObjectType      string `json:"docType"` //field for couchdb
	Id              string `json:"id"`
	OwnerId         string `json:"ownerid"`
	VoyNo           string `json:"voyno"`
	Vessel          Vessel `json:"vessel"`
	PortOfLoading   string `json:"portofloading"`
	PortOfDischarge string `json:"portofdischarge"`
	PlaceOfDelivery string `json:"placeofdelivery"`
	DepartureDate   string `json:"departuredate"`
	ArrivalDate     string `json:"arrivaldate"`
	Status          string `json:"status"`
}

type Space struct {
	TotalNum    int      `json:"totalnum"`    //舱位总数
	RestNum     int      `json:"restnum"`     //舱位剩余数目
	SpaceUsed   []string `json:"spaceused"`   //已使用的舱位号列表
	SpaceUnUsed []string `json:"spaceunused"` //未使用的舱位号列表
}

type Vessel struct {
	VesselNo string `json:"vesselno"`
	Name     string `json:"name"`
	OwnerId  string `json:"ownerid"`
	Space    Space  `json:"space"` // 舱位
}

// ============================================================================================================================
// ReadShippingSchedule - Read a voyage from the ledger
// Input - Array of strings
// 	0,
// 	shippingScheduleId,
// ============================================================================================================================
func ReadShippingSchedule(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("starting readShippingSchedule")

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
	var value *ShippingSchedule = new(ShippingSchedule)
	valueAsBytes, err := stub.GetState(value.BuildKey(id))
	if err != nil { //this seems to always succeed, even if key didn't exist
		fmt.Printf("Failed to find shipping schedule %s\n", id)
		return nil, err
	}

	fmt.Println("- end readShippingSchedule")
	return valueAsBytes, nil
}

// ============================================================================================================================
// ReadAllShippingSchedulesByOwnerId -
// Input - Array of strings
// 	0,
//	ownerId
// Returns - Array of ShippingSchedule
// [
// 	{ shippingSchedule }, ...
// ]
// ============================================================================================================================
func ReadAllShippingSchedulesByOwnerId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var shippingSchedules []ShippingSchedule

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}
	ownerId := args[0]

	// ---- Get All Containers ---- //
	indexName := common.INDEX_SHIPPING_SCHEDULE
	idList, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}

	for i, id := range idList {
		if i == 0 {
			continue
		}
		var value ShippingSchedule
		queryValAsBytes, err := stub.GetState(value.BuildKey(id))
		if err != nil {
			//todo: value that have been deleted should update index
			continue
		}
		fmt.Println("on shipping schedule id - ", id)
		json.Unmarshal(queryValAsBytes, &value) //un stringify it aka JSON.parse()
		if value.OwnerId == ownerId {
			shippingSchedules = append(shippingSchedules, value)  //add this container to the list
		}

	}
	fmt.Println("shipping schedule array - ", shippingSchedules)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(shippingSchedules) //convert to array of bytes
	return everythingAsBytes, nil

}

// ============================================================================================================================
// InsertShippingSchedule - Insert a new shipping schedule into the ledger
// inputs - Array of Strings
// 	0
// 	ShippingSchedule
// ============================================================================================================================
func InsertShippingSchedule(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting insertShippingSchedule")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	err = common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	var value *ShippingSchedule = new(ShippingSchedule)
	err = json.Unmarshal([]byte(args[0]), value)
	if err != nil {
		return nil, errors.New("Invalid type of arguments. Expecting Shipping Schedule")
	}
	if len(value.Id) == 0 {
		return nil, errors.New("Container Id should not be empty")
	}

	var owner *auth.User = new(auth.User)
	if !auth.IsUserExist(stub, value.OwnerId, owner) {
		return nil, errors.New("Owner of the shipping schedule not exists.")
	}

	value.SetObjectType()

	if err = value.PutShippingSchedule(stub); err != nil {
		return nil, err
	}

	indexName := common.INDEX_SHIPPING_SCHEDULE
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
// UpdateShippingSchedule - Update an existing shipping schedule in the ledger
// Inputs - Array of strings
// 	0,
// 	shippingSchedule
// ============================================================================================================================
func UpdateShippingSchedule(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting updateContainer")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	if err = common.SanitizeArguments(args); err != nil {
		return nil, err
	}

	var t *ShippingSchedule = new(ShippingSchedule)
	if err = json.Unmarshal([]byte(args[0]), t); err != nil {
		return nil, errors.New("Invalid type of arguments. Expecting ShippingSchedule")
	}

	var old *ShippingSchedule = new(ShippingSchedule)
	if err = old.GetShippingSchedule(stub, t.Id); err != nil {
		fmt.Println(t)
		return nil, errors.New("This shipping schedule does not exist - " + t.Id)
	}

	if err = t.PutShippingSchedule(stub); err != nil {
		fmt.Println(t)
		return nil, err
	}

	fmt.Println("- end updateShippingSchedule")
	return nil, nil
}

// ============================================================================================================================
// DeleteShippingSchedule - Delete a shipping schedule from the ledger
// Inputs - Array of strings
// 	0, 			1
// 	shippingScheduleId, 	ownerId
// ============================================================================================================================
func DeleteShippingSchedule(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	fmt.Println("starting deleteShippingSchedule")

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
	var value *ShippingSchedule = new(ShippingSchedule)
	if err := value.GetShippingSchedule(stub,id); err != nil {
		fmt.Printf("Failed to find shipping schedule %s\n", id)
		return nil, err
	}
	if value.OwnerId != ownerId {
		fmt.Printf("User %s is not authorized to delete shipping schedule %v\n", ownerId, value)
		return nil, errors.New(fmt.Sprintf("User %s is not authorized to delete shipping schedule %v\n", ownerId, value))
	}

	// remove the container
	err = stub.DelState(value.BuildKey(id)) //remove the key from chaincode state
	if err != nil {
		return nil, err
	}

	fmt.Println("- end deleteShippingSchedule")
	return nil, nil
}


// ============================================================================================================================
// Get ShippingSchedule - get shipping schedule asset from ledger
// ============================================================================================================================
func (t *ShippingSchedule) GetShippingSchedule(stub shim.ChaincodeStubInterface, id string) error {
	containerAsBytes, err := stub.GetState(t.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                                   //this seems to always succeed, even if key didn't exist
		return errors.New("Failed to find container - " + id)
	}
	err = json.Unmarshal(containerAsBytes, t) //un stringify it aka JSON.parse()
	if err != nil {
		return err
	}
	return nil
}

// ============================================================================================================================
// Put ShippingSchedule - put a shipping schedule asset into ledger
// ============================================================================================================================
func (t *ShippingSchedule) PutShippingSchedule(stub shim.ChaincodeStubInterface) error {
	valueAsBytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("ShippingSchedule: %#v", t)
		return errors.New("Failed to marshal this container: " + err.Error())
	}

	err = stub.PutState(t.BuildKey(t.Id), valueAsBytes)
	if err != nil {
		return errors.New("Failed to put this container into state: " + err.Error())
	}

	return nil
}

// ============================================================================================================================
// Build Key for a given Shipping Schedule Id
// ============================================================================================================================
func (t *ShippingSchedule) BuildKey(id string) (key string) {
	return common.ShippingScheduleKeyPrefix + id
}

func (t *ShippingSchedule) SetObjectType() {
	t.ObjectType = common.ObjectTypeShippingSchedule
}


//
func (t *ShippingSchedule) HasSpace(reqNum int) bool {
	//return t.Vessel.Space.RestNum >= reqNum
	return len(t.Vessel.Space.SpaceUnUsed) >= reqNum
}

func (t *ShippingSchedule) AllocateSpace(reqNum int) []string {
	fmt.Println("Start to allocate shipping space -")

	if !t.HasSpace(reqNum) {
		return nil
	}
	spaces := t.Vessel.Space.SpaceUnUsed[0:reqNum]
	t.Vessel.Space.SpaceUsed = append(t.Vessel.Space.SpaceUsed, spaces...)
	if len(t.Vessel.Space.SpaceUnUsed) == reqNum {
		t.Vessel.Space.SpaceUnUsed = nil
	}else {
		t.Vessel.Space.SpaceUnUsed = t.Vessel.Space.SpaceUnUsed[reqNum:]
	}
	t.Vessel.Space.RestNum = t.Vessel.Space.RestNum - reqNum
	fmt.Println("- end allocating shipping space")
	return spaces
}



func IsShippingScheduleExist(stub shim.ChaincodeStubInterface, shippingScheduleId string, voyage *ShippingSchedule) bool {
	if voyage == nil {
		voyage = new(ShippingSchedule)
	}
	if err := voyage.GetShippingSchedule(stub, shippingScheduleId); err != nil {
		return false
	}
	return true
}
