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
package rest

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/web"
	"net/http"
	"net"
	"golang.org/x/net/netutil"
)


type ServerContainerREST struct {
}


// restResult defines the response payload for a general REST interface request.
type restResult struct {
	OK    string `json:",omitempty"`
	Error string `json:",omitempty"`
}

// SetContainerServer is a middleware function that sets the pointer to the
// underlying ServerOpenchain object and the undeflying Devops object.
func (s *ServerContainerREST) SetContainerServer(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	next(rw, req)
}

// SetResponseType is a middleware function that sets the appropriate response
// headers. Currently, it is setting the "Content-Type" to "application/json" as
// well as the necessary headers in order to enable CORS for Swagger usage.
func (s *ServerContainerREST) SetResponseType(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	rw.Header().Set("Content-Type", "application/json")

	// Enable CORS
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "accept, content-type, userid, sessionid, token")

	next(rw, req)
}



// NotFound returns a custom landing page when a given hyperledger end point
// had not been defined.
func (s *ServerContainerREST) NotFound(rw web.ResponseWriter, r *web.Request) {
	rw.WriteHeader(http.StatusNotFound)
	json.NewEncoder(rw).Encode(restResult{Error: "Container endpoint not found."})
}

func buildServerContainerRESTRouter() *web.Router {
	router := web.New(ServerContainerREST{})

	// Add middleware
	router.Middleware((*ServerContainerREST).SetContainerServer)
	router.Middleware((*ServerContainerREST).SetResponseType)

	// Add routes
	router.Post("/resource/vehicle", (*ServerContainerREST).InsertVehicle)
	router.Put("/resource/vehicle", (*ServerContainerREST).UpdateVehicle)
	router.Get("/resource/vehicle/:id", (*ServerContainerREST).GetVehicleById)
	router.Delete("/resource/vehicle/:id", (*ServerContainerREST).DeleteVehicle)
	router.Get("/resource/vehicle/findByOwnerId", (*ServerContainerREST).FindVehiclesByOwnerId)

	router.Post("/resource/shippingschedule", (*ServerContainerREST).InsertShippingSchedule)
	router.Put("/resource/shippingschedule", (*ServerContainerREST).UpdateShippingSchedule)
	router.Get("/resource/shippingschedule/:id", (*ServerContainerREST).GetShippingScheduleById)
	router.Delete("/resource/shippingschedule/:id", (*ServerContainerREST).DeleteShippingSchedule)
	router.Get("/resource/shippingschedule/findByOwnerId", (*ServerContainerREST).FindShippingSchedulesByOwnerId)

	router.Post("/resource/container", (*ServerContainerREST).InsertContainer)
	router.Put("/resource/container", (*ServerContainerREST).UpdateContainer)
	router.Get("/resource/container/:id", (*ServerContainerREST).GetContainerById)
	router.Delete("/resource/container/:id", (*ServerContainerREST).DeleteContainer)
	router.Get("/resource/container/findByOwnerId", (*ServerContainerREST).FindContainersByOwnerId)
	router.Post("/resource/container/track", (*ServerContainerREST).TrackContainers)

	router.Get("/resource/transporttask/:id", (*ServerContainerREST).GetTransportTaskById)
	router.Delete("/resource/transporttask/:id", (*ServerContainerREST).DeleteTransportTask)
	router.Get("/resource/transporttask/findByOwnerId", (*ServerContainerREST).FindTransportTasksByOwnerId)

	router.Get("/order/:id", (*ServerContainerREST).GetOrderById)
	router.Delete("/order/:id", (*ServerContainerREST).DeleteOrder)
	router.Get("/order/findByUserId", (*ServerContainerREST).FindOrdersByUserId)
	router.Post("/order/client/create", (*ServerContainerREST).CreateOrder)
	router.Post("/order/cargoagent/check", (*ServerContainerREST).CheckOrder)
	router.Post("/order/cargoagent/bookspace", (*ServerContainerREST).BookSpace)
	router.Post("/order/cargoagent/bookvehicle", (*ServerContainerREST).BookVehicle)
	router.Post("/order/carrier/fetchemptycontainers", (*ServerContainerREST).FetchEmptyContainers)
	router.Post("/order/carrier/packgoods", (*ServerContainerREST).PackGoods)
	router.Post("/order/carrier/arriveyard", (*ServerContainerREST).ArriveYard)
	router.Post("/order/shipper/loadgoods", (*ServerContainerREST).LoadGoods)
	router.Post("/order/shipper/departure", (*ServerContainerREST).Departure)
	router.Post("/order/shipper/arrivedestinationport", (*ServerContainerREST).ArriveDestinationPort)
	router.Post("/order/shipper/delivergoods", (*ServerContainerREST).DeliverGoods)
	router.Post("/order/client/confirmreceipt", (*ServerContainerREST).ConfirmReceipt)
	router.Post("/order/cargoagent/finish", (*ServerContainerREST).FinishOrder)

	router.Post("/message", (*ServerContainerREST).InsertMessage)
	router.Put("/message", (*ServerContainerREST).UpdateMessageStatus)
	router.Get("/message/:id", (*ServerContainerREST).GetMessageById)
	router.Delete("/message/:id", (*ServerContainerREST).DeleteMessage)
	router.Get("/message/findByUserId", (*ServerContainerREST).FindMessagesByUserId)

	router.Post("/user", (*ServerContainerREST).RegisterUser)
	router.Put("/user", (*ServerContainerREST).UpdateUser)
	router.Get("/user/:id", (*ServerContainerREST).GetUserById)
	router.Get("/user/findByUserRoleType", (*ServerContainerREST).FindUsersByType)

	router.Get("/user/session/login", (*ServerContainerREST).LogIn)
	router.Get("/user/session/refresh", (*ServerContainerREST).Refresh)
	router.Get("/user/session/logout", (*ServerContainerREST).LogOut)

	// Add not found page
	router.NotFound((*ServerContainerREST).NotFound)

	return router
}

// StartContainerRESTServer initializes the REST service and adds the required middleware and routes.
func StartContainerRESTServer() {
	// Initialize the REST service object
	restLogger.Info("Initializing the REST service on localhost.")

	router := buildServerContainerRESTRouter()

	listener, err := net.Listen("tcp", localServerAddress)
	if err != nil {
		restLogger.Errorf("Failed to listen on port 9090: %v", err)
	}
	defer listener.Close()

	restLogger.Infof("Start to listen and serve on %s, max limited connection is %d", localServerAddress, localServerMaxConnectionLimit)

	listener = netutil.LimitListener(listener, localServerMaxConnectionLimit)
	err = http.Serve(listener, router)
	if err != nil {
		restLogger.Errorf("ListenAndServe: %s", err)
	}

}
