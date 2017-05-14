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
	"cbdforum/app-container/chaincode/src/common"
	"cbdforum/app-container/chaincode/src/resource"
	auth "cbdforum/app-container/chaincode/src/user"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var logger *logging.Logger = logging.MustGetLogger("Order")

type Goods struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Measurement int    `json:"measurement"`
	GrossWeight int    `json:"grossweight"`
}

type Order struct {
	ObjectType     string `json:"docType"`
	Id             string `json:"id"`
	OrderNo        string `json:"orderno"`

	ConsigningForm ConsigningForm `json:"consigningform"`
	BookingForm    BookingForm    `json:"bookingform"`
	CarryingForm   CarryingForm   `json:"carryingform"`

	State          string   `json:"state"`
	DeletedByWho   []string `json:"deletedbywho"`
	Remark         string   `json:"remark"`  // 用来记录订单被拒绝的原因

	CreatedAt      string `json:"createdat"` // 订单创建的时间
}

type ConsigningForm struct {
	ConsigningFormNo string `json:"consigningformno"`
	ClientId         string `json:"clientid"`
	CargoAgentId     string `json:"cargoagentid"`

	GoodsList            []Goods `json:"goodslist"`
	DeliveryAddress      string  `json:"deliveryaddress"`
	ShippingAddress      string  `json:"shippingaddress"`
	Consignee            string  `json:"consignee"`
	ConsigneePhone       string  `json:"consigneePhone"`
	Consignor            string  `json:"consignor"`
	ConsignorPhone       string  `json:"consignorPhone"`
	ExpectedDeliveryDate string  `json:"expecteddeliverydate"`

	DateForConfirmReceipt string `json:"dateforconfirmreceipt"` // 确认签收时间, 由 client 写入
	DateForFinish         string `json:"dateforfinish"`         // 订单最终完成时间，由cargo agent 写入

}

type BookingForm struct {
	BookingFormNo       string                    `json:"bookingformno"`
	ShipperId           string                    `json:"shipperid"`
	Voyage              resource.ShippingSchedule `json:"voyage"`
	BerthNo             []string                  `json:"berthno"`
	Containers          []resource.Container      `json:"containers"`
	DateForLoading      string                    `json:"dateforloading"`      // 货物装船时间
	DateForDeparture    string                    `json:"datefordeparture"`    // 起航时间
	DateForArrival      string                    `json:"dateforarrival"`      // 到达时间
	DateForDeliverGoods string                    `json:"datefordelivergoods"` // 送货时间
}

type CarryingForm struct {
	CarryingFormNo  string             `json:"carryingformno"`
	CarrierId       string             `json:"carrierid"`
	TransportTaskId string             `json:"transporttaskid"`
	Vehicles        []resource.Vehicle `json:"vehicles"`
	PackingList     PackingList        `json:"packinglist"`
	Status          string             `json:"status"`
	DateForReceiver string             `json:"dateforreceiver"` // 设备移入时间
	DateForDeliver  string             `json:"datefordeliver"`  // 设备移出时间
}

type PackingList struct {
	Items               []PackingListItem `json:"items"`
	DateForPackingGoods string            `json:"dateforpackinggoods"` //装箱时间
}

type PackingListItem struct {
	ContainerId string `json:"containerid"`
	Goods       Goods  `json:"goods"`
}

// ============================================================================================================================
// ReadOrder - Read an existing order from the ledger
// Inputs - Array of strings
// 	0
// 	orderId
// ============================================================================================================================
func ReadOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	fmt.Println("starting ReadOrder")

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
	var value *Order = new(Order)
	valueAsBytes, err := stub.GetState(value.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                                 //this seems to always succeed, even if key didn't exist
		fmt.Printf("Failed to find order %s ", id)
		return nil, err
	}


	fmt.Println("- end ReadOrder")
	return valueAsBytes, nil
}

// Notice: 输入只需要一个userid, 需要根据用户角色来判断具体的id
// ============================================================================================================================
// Get all containers by userid (containers)
//
// Inputs - Array of strings
//	1
//	userId
// Returns - Array of Order
//
// ============================================================================================================================
func ReadAllOrdersByUserId(stub shim.ChaincodeStubInterface, args []string)([]byte, error) {
	var orders []Order
	var user *auth.User = new(auth.User)

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}
	userId := args[0]

	if !auth.IsUserExist(stub, userId, user) {
		return nil, errors.New(fmt.Sprintf("User %s not exist", userId))
	}

	// ---- Get All Containers ---- //
	indexName := common.INDEX_ORDER
	idList, err := common.GetStringList(stub, indexName)
	if err != nil {
		return nil, err
	}

	for i, id := range idList {
		if i == 0 {
			continue
		}
		var value Order
		queryValAsBytes, err := stub.GetState(value.BuildKey(id))
		if err != nil {
			//todo: value that have been deleted should update index
			continue
		}
		fmt.Println("on container id - ", id)
		json.Unmarshal(queryValAsBytes, &value) //un stringify it aka JSON.parse()
		if value.ConsigningForm.ClientId == userId || value.ConsigningForm.CargoAgentId == userId || value.BookingForm.ShipperId == userId || value.CarryingForm.CarrierId == userId {
			if (value.IsDeletedByUser(userId)){
				continue
			}
			orders = append(orders, value)  //add this container to the list
		}

	}
	fmt.Println("order array - ", orders)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(orders) //convert to array of bytes
	return everythingAsBytes, nil
}

// Note: Should never delete order, record those who want to delete the order instead
// ============================================================================================================================
// DeleteOrder - record those trying to delete the order
// Inputs - Array of strings
// 	0,		1
// 	orderId,	userId
// ============================================================================================================================
func DeleteOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
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
	userId := args[1]

	// get the container
	var value *Order = new(Order)
	if err := value.GetOrder(stub, id); err != nil {
		fmt.Println("Failed to find contianer by id " + id)
		return nil, err
	}

	if !value.IsOrderFinished() {
		fmt.Printf("Order %v not finished yet, should not be deleted\n", value)
		return nil, errors.New("error to delete unfinished order")
	}

	value.DeletedByWho = append(value.DeletedByWho, userId)

	if err = value.PutOrder(stub); err != nil {
		return nil, NewShimError(ERROR_INTERNAL, err.Error())
	}

	fmt.Println("- end deleteOrder")
	return nil, nil
}

// ============================================================================================================================
// Get Order - get an order asset from ledger
// ============================================================================================================================
func (t *Order) GetOrder(stub shim.ChaincodeStubInterface, id string) error {
	valueAsBytes, err := stub.GetState(t.BuildKey(id)) //getState retreives a key/value from the ledger
	if err != nil {                                    //this seems to always succeed, even if key didn't exist
		return errors.New("Failed to find order - " + id)
	}

	err = json.Unmarshal(valueAsBytes, t) //un stringify it aka JSON.parse()
	if err != nil {
		return err
	}

	return nil
}

// ============================================================================================================================
// Put Order - put an order asset into ledger
// ============================================================================================================================
func (t *Order) PutOrder(stub shim.ChaincodeStubInterface) error {
	valueAsBytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("%#v", t)
		return errors.New("Failed to marshal this order: " + err.Error())
	}

	err = stub.PutState(t.BuildKey(t.Id), valueAsBytes)
	if err != nil {
		return errors.New("Failed to put this order into state: " + err.Error())
	}

	return nil
}

// ============================================================================================================================
// Build Key for a given Container Id
// ============================================================================================================================
func (t *Order) BuildKey(id string) (key string) {
	return common.OrderKeyPrefix + id
}

func (t *Order) SetObjectType() {
	t.ObjectType = common.ObjectTypeOrder
}

func (t *Order) IsOrderFinished() bool {
	return t.State == STATE_FAILED || t.State == STATE_FINISHED
}

func (t *Order) IsDeletedByUser(userId string) bool {
	for _, item := range t.DeletedByWho {
		if item == userId {
			return true
		}
	}
	return false
}

func IsOrderExist(stub shim.ChaincodeStubInterface, orderId string, order *Order) bool {
	if err := order.GetOrder(stub, orderId); err != nil {
		return false
	}
	return true
}

/* Notice: operation not supported by fabric-0.6,
// ============================================================================================================================
// Get history of asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "01490985296352SjAyM"
// ============================================================================================================================
func GetHistoryForOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId  string `json:"txId"`
		Value Order  `json:"value"`
	}
	var history []AuditHistory
	var order Order

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	orderId := args[0]
	fmt.Printf("- start getHistoryForOrder: %s\n", orderId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(order.BuildKey(orderId))
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
		tx.TxId = txID                        //copy transaction id over
		json.Unmarshal(historicValue, &order) //un stringify it aka JSON.parse()
		if historicValue == nil {             //container has been deleted
			var emptyOrder Order
			tx.Value = emptyOrder //copy nil container
		} else {
			json.Unmarshal(historicValue, &order) //un stringify it aka JSON.parse()
			tx.Value = order                      //copy container over
		}
		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistoryForOrder returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
}

*/
