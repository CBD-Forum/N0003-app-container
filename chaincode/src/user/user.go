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
	"cbdforum/app-container/chaincode/src/common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"

)

type Person struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhoneNum string `json:"phonenum"`
	Address  string `json:"fax"`
	Fax      string
}

type Company struct {
	OrgId      string `json:"orgid"`
	TaxId      string `json:"taxid"`
	CreditId   string `json:"creditid"`
	BusinessId string `json:"businessid"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Homepage   string `json:"homepage"`
}

const (
	ROLE_UNDEFINED      string = "undefined"
	ROLE_REGULAR_CLIENT          = "regularclient"
	ROLE_CARGO_AGENT             = "cargoagent"
	ROLE_CARRIER                 = "carrier"
	ROLE_SHIPPER                 = "shipper"
)

type User struct {
	ObjectType   string   `json:"docType"`
	Id           string   `json:"id"`
	UserName     string   `json:"username"`
	PersonalInfo Person   `json:"personalinfo"`
	Company      Company  `json:"company"`
	Role         string   `json:"role"`
	CreatedAt    string   `json:"createdat"`
	UpdatedAt    string   `json:"updatedat"`
	DeletedAt    string   `json:"deltedat"`
}

// ============================================================================================================================
// ReadUserById - Get a user by id
// Inputs -
//	0
//	userId
// Returns - Array of User
// ============================================================================================================================
func ReadUserById(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var user *User = new(User)
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	userId := args[0]

	//input sanitation
	err = common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	if !IsUserExist(stub, userId, user) {
		return nil, errors.New("Failed to find user by Id " + userId)
	}
	fmt.Printf("get user %s: %v\n", userId, user)

	valueAsBytes, err := json.Marshal(user)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to stringfy user as bytes:%v", err))
	}

	return valueAsBytes, nil
}

// ============================================================================================================================
// ReadAllCargoAgents - Get All Cargo Agents
// Inputs - none
// Returns - Array of User
// ============================================================================================================================
func ReadAllCargoAgents(stub shim.ChaincodeStubInterface)([]byte, error){
	var Companies []User
	var err error

	Companies, err = getUsersByType(stub, ROLE_CARGO_AGENT)
	if err != nil {
		return nil, err
	}

	fmt.Println("cargo agent array - ", Companies)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(Companies) //convert to array of bytes
	return everythingAsBytes, nil
}

// ============================================================================================================================
// ReadAllShippers - Get All Shipping Companies
// Inputs - none
// Returns - Array of User
// ============================================================================================================================
func ReadAllShippers(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var Companies []User
	var err error

	Companies, err = getUsersByType(stub, ROLE_SHIPPER)
	if err != nil {
		return nil, err
	}

	fmt.Println("shipping company array - ", Companies)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(Companies) //convert to array of bytes
	return everythingAsBytes, nil
}

// ============================================================================================================================
// ReadAllCarriers - Get All Car Companies
// inputs - none
// outputs -
// {
//	Companies: [ {
//		"id": "",
//	}]
// }
// ============================================================================================================================
func ReadAllCarriers(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var Companies []User
	var err error

	Companies, err = getUsersByType(stub, ROLE_CARRIER)
	if err != nil {
		return nil, err
	}

	fmt.Println("car company array - ", Companies)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(Companies) //convert to array of bytes
	return everythingAsBytes, nil
}

// ============================================================================================================================
// RegisterUser - Register a new user into the ledger
// Inputs - Array of strings
// 	0
// 	User
// Returns -
// ============================================================================================================================
func RegisterUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting registerUser")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	err = common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	var user *User = new(User)
	err = json.Unmarshal([]byte(args[0]), user)
	if err != nil {
		return nil, errors.New("Invalid type of arguments. Expecting user")
	}
	fmt.Printf("register user: %v", *user)
	if len(user.Id) == 0 {
		return nil, errors.New("Error arguments, user.id should not be empty")
	}
	user.SetObjectType()
	if err = user.PutUser(stub); err != nil {
		return nil, err
	}

	//  ==== Index the user id to enable id range queries, e.g. return all blue marbles ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := common.INDEX_USER
	values, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}
	values = append(values, user.Id)
	err = common.PutStringList(stub, indexName, values)
	if err != nil {
		return nil, err
	}

	fmt.Println("- end registerUser")
	return nil, nil
}

// ============================================================================================================================
// UpdateUser - update user info
// Inputs - Array of strings
// 	0,
// 	User
// ============================================================================================================================
func UpdateUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting updateUser")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//input sanitation
	if err = common.SanitizeArguments(args); err != nil {
		return nil, err
	}

	var user *User = new(User)
	if err = json.Unmarshal([]byte(args[0]), user); err != nil {
		return nil, errors.New("Invalid type of arguments. Expecting User")
	}

	// check whether user exist
	var old *User = new(User)
	if err = old.GetUser(stub, user.Id); err != nil {
		fmt.Println(user)
		return nil, errors.New("This user does not exist - " + user.Id) //all stop a user by this id not exists
	}

	if err = user.PutUser(stub); err != nil {
		fmt.Println(user)
		return nil, err
	}

	fmt.Println("- end updateUser")
	return nil, nil
}

// ============================================================================================================================
// DeleteUser - delete a user from the ledger
// Inputs - Array of strings
// 	0
// 	userId
// ============================================================================================================================
func DeleteUesr(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	fmt.Println("starting deleteUser")

	if len(args) != 1 {
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}

	// input sanitation
	err := common.SanitizeArguments(args)
	if err != nil {
		return nil, err
	}

	userId := args[0]

	// get the user
	var user *User = new(User)
	if !IsUserExist(stub, userId, user) {
		fmt.Println("Failed to find user by id " + userId)
		return nil, errors.New("Error user " + userId + "  not exist")
	}

	// remove the user
	err = stub.DelState(user.BuildKey(userId)) //remove the key from chaincode state
	if err != nil {
		return nil, err
	}

	fmt.Println("- end deletUser")
	return nil, nil
}

// ============================================================================================================================
// Get User - get a user asset from ledger
// ============================================================================================================================
func (t *User) GetUser(stub shim.ChaincodeStubInterface, id string) error {
	valueAsBytes, err := stub.GetState(t.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                    //this seems to always succeed, even if key didn't exist
		return errors.New("Failed to find user - " + id)
	}
	json.Unmarshal(valueAsBytes, t) //un stringify it aka JSON.parse()

	if t.Id != id {
		return errors.New("User does not exist - " + id)
	}

	return nil
}

// ============================================================================================================================
// Put User - put a user asset into ledger
// ============================================================================================================================
func (t *User) PutUser(stub shim.ChaincodeStubInterface) error {
	valueAsBytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("User: %#v", t)
		return errors.New("Failed to marshal this user: " + err.Error())
	}

	err = stub.PutState(t.BuildKey(t.Id), valueAsBytes)
	if err != nil {
		return errors.New("Failed to put this user into state: " + err.Error())
	}

	return nil
}

// ============================================================================================================================
// Build Key for a given user Id
// ============================================================================================================================
func (t *User) BuildKey(id string) (key string) {
	return common.UserKeyPrefix + id
}

func (t *User) IsRegularClient() bool {
	return t.Role == ROLE_REGULAR_CLIENT
}

func (t *User) IsCargoAgent() bool {
	return t.Role == ROLE_CARGO_AGENT
}

func (t *User) IsCarrier() bool {
	return t.Role == ROLE_CARRIER
}

func (t *User) IsShipper() bool {
	return t.Role == ROLE_SHIPPER
}

func (t *User) SetObjectType() {
	t.ObjectType = common.ObjecTTypeUser
}

func IsUserExist(stub shim.ChaincodeStubInterface, userId string, user *User) bool {
	if err := user.GetUser(stub, userId); err != nil {
		return false
	}
	return true
}
