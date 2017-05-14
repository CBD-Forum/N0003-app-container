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
 * Created by Otherplayer on 16/7/21.
 */
(function () {
    'use strict';

    angular.module('airs').controller('orderCarAddController', orderCarAddController);

    /** @ngInject */
    function orderCarAddController($log,$uibModal,$state,ApiServer,$stateParams,toastr,$filter) {
        /* jshint validthis: true */
        var vm = this;

        var orderId = $stateParams.orderId;
        
        var info = ApiServer.info();

        vm.orderInfo = {consigningform:{goodslist:[]}};
        vm.companys = [];
        vm.company = {};

        vm.acceptAction = acceptAction;


        function acceptAction(type) {
            if (type === 'get'){
                var param = {carrierid:info.id,orderid:orderId,datefordelivered:$filter('date')(new Date(),'yyyy-MM-dd HH:mm:ss')};
                ApiServer.carOrderFetchEmptyContainers(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_empty_container_fetched';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else if (type === 'pack'){
                $state.go('app.carorderaddpacking',{orderId:orderId});
            }else if (type === 'arrive'){
                var param = {carrierid:info.id,orderid:orderId,dateforreceiver:$filter('date')(new Date(),'yyyy-MM-dd HH:mm:ss')};
                console.log(param);
                ApiServer.carOrderArriveyard(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_yard_arrived';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else{
                $state.go('app.carorder');
            }
        }

        function dealOrderAction(accept) {
            if (accept === 'finish'){
                var now = new Date();
                var param = {"clientid": info.id, "OrderId": orderId, "DateForFinish": now};
                ApiServer.goodOrderAccept(param,function (res) {

                },function (err) {

                })
            }else if (accept === 'get'){
                var param = {carrierid:info.id,orderid:orderId,datefordelivered:'2017-5-21 18:20:00'};
                ApiServer.carOrderFetchEmptyContainers(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_failed';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else{
                var param = {cargoAgentId:'',orderId:orderId,isOrderAccept:false,remark:''};
                if (accept === 'accept'){
                    param.isOrderAccept = true;
                }
                ApiServer.goodOrderAccept(param,function (res) {

                },function (err) {

                })
            }
        }

        function getDatas() {
            ApiServer.orderGet(orderId,function (res) {
                vm.orderInfo = res.data;
                console.log(vm.orderInfo);
                getCompanyById(vm.orderInfo.consigningform.cargoagentid);
            },function (err) {
                var errInfo = '获取数据失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            })
        }

        function getCompanyDatas() {
            ApiServer.userGetByRoleType('cargoagent',function (res) {
                vm.companys = res.data;
                getDatas();
            },function (err) {
                var errInfo = '获取数据失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            })
        }

        function getCompanyById(cargoagentId) {
            for (var i = 0; i < vm.companys.length; i++){
                if (vm.companys[i].id === cargoagentId){
                    vm.company = vm.companys[i];
                    break;
                }
            }
        }

        getCompanyDatas();


        //Model
        vm.openAlert = function (size,accept) {
            var modalInstance = $uibModal.open({
                templateUrl: 'myModalContent.html',
                size: size,
                controller:'ModalInstanceCtrl',
                resolve: {
                    tipsInfo: function () {
                        return vm.tipsInfo;
                    }
                }
            });
            modalInstance.result.then(function (param) {
                dealOrderAction(accept);
            }, function () {
                $log.info('Modal dismissed at: ' + new Date());
            });
        };


    }

})();
