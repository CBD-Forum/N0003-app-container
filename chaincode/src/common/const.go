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

const (
	MessageServer         string = "messageserver"
	ContainerUIVersion     = "containeruiersion"

	ContainerKeyPrefix        = "container-"
	ShippingScheduleKeyPrefix = "shippingschedule-"
	VehicleKeyPrefix          = "vehicle-"
	TransportTaskKeyPrefix    = "transporttask-"
	UserKeyPrefix             = "user-"
	OrderKeyPrefix            = "order-"

	ResourceStatusFree  = "free"
	ResourceStatusInUse = "inuse"

	StatusInit     = "init"
	StatusFinished = "finished"
	StatusFailed   = "failed"

	ObjecTTypeUser = "user"
	ObjectTypeContainer = "container"
	ObjectTypeShippingSchedule = "shippingschedule"
	ObjectTypeVehicle = "vehicle"
	ObjectTypeTransportTask = "transporttask"
	ObjectTypeOrder = "order"
)

const (
	GOLANG_TIME_FMT_STRING = "2006-01-02 15:04:05"
)

const (
	INDEX_USER = "user-id-list"
	INDEX_CONTAINER = "container-id-list"
	INDEX_VEHICLE = "vehicle-id-list"
	INDEX_TRANSPORT_TASK = "transporttask-id-list"
	INDEX_SHIPPING_SCHEDULE = "shippingschedule-id-list"
	INDEX_ORDER = "order-id-list"
)



