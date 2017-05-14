// Copyright [2016] [Lele Guo]
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
/**
 * Created by Otherplayer on 2016/12/12.
 */

(function () {
    'use strict';


    /**
     *
     * facotry是一个单例,它返回一个包含service成员的对象。
     * 注：所有的Angular services都是单例，这意味着每个injector都只有一个实例化的service。
     *
     */
    angular
        .module('airs')
        .factory('ApiServer', ApiServer);

    /** @ngInject */
    function ApiServer(constdata,NetworkService,StorageService,iotUtil,$timeout,$state,toastr) {


        var service = {

            isAuthed: isAuthed,
            roleType: roleType,
            info: info,
            logoutAction: logoutAction,

            userRegister: userRegister,
            userUpdate: userUpdate,
            userGet: userGet,
            userGetByRoleType: userGetByRoleType,

            userLogin: userLogin,
            userRefresh: userRefresh,
            userLogout: userLogout,

            messageAdd: messageAdd,
            messageGet: messageGet,
            messageGetByUserId: messageGetByUserId,
            messageDelete: messageDelete,

            vehicleAdd: vehicleAdd,
            vehicleUpdate: vehicleUpdate,
            vehicleGet: vehicleGet,
            vehicleGetByOwner: vehicleGetByOwner,
            vehicleDelete: vehicleDelete,

            shippingScheduleAdd: shippingScheduleAdd,
            shippingScheduleUpdate: shippingScheduleUpdate,
            shippingScheduleGet: shippingScheduleGet,
            shippingScheduleGetByOwner: shippingScheduleGetByOwner,
            shippingScheduleDelete: shippingScheduleDelete,

            containerAdd: containerAdd,
            containerUpdate: containerUpdate,
            containerGet: containerGet,
            containerGetByOwner: containerGetByOwner,
            containerDelete: containerDelete,
            containerTrack: containerTrack,

            transportTaskGet: transportTaskGet,
            transportTaskGetByOwner: transportTaskGetByOwner,
            transportTaskDelete: transportTaskDelete,

            orderGet:orderGet,
            orderGetByOwner:orderGetByOwner,
            orderDelete:orderDelete,
            clientOrderAdd: clientOrderAdd,
            clientOrderConfirmReceipt: clientOrderConfirmReceipt,
            goodOrderAccept: goodOrderAccept,
            goodOrderCar: goodOrderCar,
            goodOrderSpace: goodOrderSpace,
            goodOrderFinish: goodOrderFinish,
            carOrderFetchEmptyContainers: carOrderFetchEmptyContainers,
            carOrderPackgoods: carOrderPackgoods,
            carOrderArriveyard: carOrderArriveyard,
            shipOrderLoadgoods: shipOrderLoadgoods,
            shipOrderDepartrue: shipOrderDepartrue,
            shipOrderArriveDestinationPort: shipOrderArriveDestinationPort,
            shipOrderDeliverGoods: shipOrderDeliverGoods

        };

        return service;

        ////////////信息

        function isAuthed() {
            var token = StorageService.get(constdata.token);
            if (token && token !== 'undefined'){
                return true;
            }
            return false;
        }
        function roleType() {
            // "regularclient"; "cargoagent"; "carrier"; "shipper";
            var information = info();
            return information.role;
        }
        function info() {
            var information = StorageService.get(constdata.informationKey);
            if (!information || information === 'undefined'){
                // toastr.error('服务器出错了，请稍后重试');
                logoutAction();
            }else{
                return information;
            }
        }
        function logoutAction() {
            var authorizationKey = constdata.token;
            var userInfo = constdata.informationKey;
            $timeout(function () {
                StorageService.clear(authorizationKey);
                StorageService.clear(userInfo);
                StorageService.clear(constdata.token);
            },60);
            $state.go('access.signin');
        }

        ////////////用户

        function userRegister(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.user,param,successHandler,failedHandler);
        }
        function userUpdate(param,successHandler,failedHandler) {
            NetworkService.put(constdata.api.user,param,successHandler,failedHandler);
        }
        function userGet(userId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.user + '/' + userId,null,successHandler,failedHandler);
        }
        function userGetByRoleType(type,successHandler,failedHandler) {
            NetworkService.get(constdata.api.user + '/findByUserRoleType?roletype=' + type,null,successHandler,failedHandler);
        }

        function userLogin(name,pwd,successHandler,failedHandler) {
            NetworkService.get(constdata.api.user + '/session/login?username=' + name + '&password=' + iotUtil.sha256(pwd),null,successHandler,failedHandler);
        }
        function userRefresh(successHandler,failedHandler) {
            NetworkService.put(constdata.api.user + '/session/refresh',null,successHandler,failedHandler);
        }
        function userLogout(successHandler,failedHandler) {
            NetworkService.put(constdata.api.user + '/session/logout',null,successHandler,failedHandler);
        }

        ////////////Message

        function messageAdd(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.message,param,successHandler,failedHandler);
        }
        function messageGet(messageId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.message + '/' + messageId,null,successHandler,failedHandler);
        }
        function messageGetByUserId(userId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.message + '/findByUserId?userId=' + userId,null,successHandler,failedHandler);
        }
        function messageDelete(messageId,successHandler,failedHandler) {
            NetworkService.delete(constdata.api.message + '/' + messageId,null,successHandler,failedHandler);
        }


        ////////////Vehicle 拖车

        function vehicleAdd(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.resource.vehicle,param,successHandler,failedHandler);
        }
        function vehicleUpdate(param,successHandler,failedHandler) {
            NetworkService.put(constdata.api.resource.vehicle,param,successHandler,failedHandler);
        }
        function vehicleGet(vehicleId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.resource.vehicle + '/' + vehicleId,null,successHandler,failedHandler);
        }
        function vehicleGetByOwner(ownerId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.resource.vehicle + '/findByOwnerId?ownerid=' + ownerId,null,successHandler,failedHandler);
        }
        function vehicleDelete(vehicleId,successHandler,failedHandler) {
            NetworkService.delete(constdata.api.resource.vehicle + '/' + vehicleId,null,successHandler,failedHandler);
        }

        ////////////ShippingSchedule

        function shippingScheduleAdd(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.resource.shippingSchedule,param,successHandler,failedHandler);
        }
        function shippingScheduleUpdate(param,successHandler,failedHandler) {
            NetworkService.put(constdata.api.resource.shippingSchedule,param,successHandler,failedHandler);
        }
        function shippingScheduleGet(shippingscheduleId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.resource.shippingSchedule + '/' + shippingscheduleId,null,successHandler,failedHandler);
        }
        function shippingScheduleGetByOwner(ownerId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.resource.shippingSchedule + '/findByOwnerId?ownerid=' + ownerId,null,successHandler,failedHandler);
        }
        function shippingScheduleDelete(shippingscheduleId,successHandler,failedHandler) {
            NetworkService.delete(constdata.api.resource.shippingSchedule + '/' + shippingscheduleId,null,successHandler,failedHandler);
        }

        ////////////Container

        function containerAdd(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.resource.container,param,successHandler,failedHandler);
        }
        function containerUpdate(param,successHandler,failedHandler) {
            NetworkService.put(constdata.api.resource.container,param,successHandler,failedHandler);
        }
        function containerGet(containerId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.resource.container + '/' + containerId,null,successHandler,failedHandler);
        }
        function containerTrack(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.resource.container + '/track',param,successHandler,failedHandler);
        }
        function containerGetByOwner(ownerId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.resource.container + '/findByOwnerId?ownerid=' + ownerId,null,successHandler,failedHandler);
        }
        function containerDelete(containerId,successHandler,failedHandler) {
            NetworkService.delete(constdata.api.resource.container + '/' + containerId,null,successHandler,failedHandler);
        }

        ////////////TransportTask

        function transportTaskGet(transportTaskId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.resource.transportTask + '/' + transportTaskId,null,successHandler,failedHandler);
        }
        function transportTaskGetByOwner(ownerId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.resource.transportTask + '/findByOwnerId?ownerid=' + ownerId,null,successHandler,failedHandler);
        }
        function transportTaskDelete(transportTaskId,successHandler,failedHandler) {
            NetworkService.delete(constdata.api.resource.transportTask + '/' + transportTaskId,null,successHandler,failedHandler);
        }

        ////////////Order
        function orderGet(orderId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.order + '/' + orderId,null,successHandler,failedHandler);
        }
        function orderDelete(containerId,successHandler,failedHandler) {
            NetworkService.delete(constdata.api.order + '/' + containerId,null,successHandler,failedHandler);
        }
        function orderGetByOwner(ownerId,successHandler,failedHandler) {
            NetworkService.get(constdata.api.order + '/findByUserId?ownerid=' + ownerId,null,successHandler,failedHandler);
        }

        function clientOrderAdd(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/client/create',param,successHandler,failedHandler);
        }
        function clientOrderConfirmReceipt(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/client/confirmreceipt',param,successHandler,failedHandler);
        }

        function goodOrderAccept(param,successHandler,failedHandler) {//是否接受订单
            NetworkService.post(constdata.api.order + '/cargoagent/check',param,successHandler,failedHandler);
        }
        function goodOrderSpace(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/cargoagent/bookspace',param,successHandler,failedHandler);
        }
        function goodOrderCar(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/cargoagent/bookvehicle',param,successHandler,failedHandler);
        }
        function goodOrderFinish(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/cargoagent/finish',param,successHandler,failedHandler);
        }

        function carOrderFetchEmptyContainers(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/carrier/fetchemptycontainers',param,successHandler,failedHandler);
        }
        function carOrderPackgoods(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/carrier/packgoods',param,successHandler,failedHandler);
        }
        function carOrderArriveyard(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/carrier/arriveyard',param,successHandler,failedHandler);
        }

        function shipOrderLoadgoods(param,successHandler,failedHandler) {
            NetworkService.post(constdata.api.order + '/shipper/loadgoods',param,successHandler,failedHandler);
        }
        function shipOrderDepartrue(param,successHandler,failedHandler) {//集装箱离港时间
            NetworkService.post(constdata.api.order + '/shipper/departure',param,successHandler,failedHandler);
        }
        function shipOrderArriveDestinationPort(param,successHandler,failedHandler) {//集装箱到目的港口时间
            NetworkService.post(constdata.api.order + '/shipper/arrivedestinationport',param,successHandler,failedHandler);
        }
        function shipOrderDeliverGoods(param,successHandler,failedHandler) {//货物送至客户
            NetworkService.post(constdata.api.order + '/shipper/delivergoods',param,successHandler,failedHandler);
        }



    }

})();
