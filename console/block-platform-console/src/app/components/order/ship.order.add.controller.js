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

    angular.module('airs').controller('orderShipAddController', orderShipAddController);

    /** @ngInject */
    function orderShipAddController($scope,$log,$uibModal,ApiServer,$stateParams,toastr,$state,$filter) {
        /* jshint validthis: true */
        var vm = this;


        var orderId = $stateParams.orderId;
        var info = ApiServer.info();
        vm.orderInfo = {consigningform:{goodslist:[]}};
        vm.companys = [];
        vm.company = {};

        vm.acceptAction = acceptAction;


        function acceptAction(type) {
            dealOrderAction(type);
        }

        function dealOrderAction(accept) {
            if (accept === 'loadgoods'){
                var param = {shipperid:info.id,orderid:orderId,dateforloading:$filter('date')(new Date(),'yyyy-MM-dd HH:mm:ss')};
                ApiServer.shipOrderLoadgoods(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_goods_loaded';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else if (accept === 'departrue'){
                var param = {shipperid:info.id,orderid:orderId,datefordeparture:$filter('date')(new Date(),'yyyy-MM-dd HH:mm:ss')};
                ApiServer.shipOrderDepartrue(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_goods_shipping';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else if (accept === 'arrivedestinationport'){
                var param = {shipperid:info.id,orderid:orderId,dateforarrival:$filter('date')(new Date(),'yyyy-MM-dd HH:mm:ss')};
                ApiServer.shipOrderArriveDestinationPort(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_goods_arrived';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else if (accept === 'delivergoods'){
                var param = {shipperid:info.id,orderid:orderId,datefordelivergoods:$filter('date')(new Date(),'yyyy-MM-dd HH:mm:ss')};
                ApiServer.shipOrderDeliverGoods(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_goods_delivered';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else {
                $state.go('app.shiporder');
            }
        }

        function getDatas() {
            ApiServer.orderGet(orderId,function (res) {
                vm.orderInfo = res.data;
                getCompanyById(vm.orderInfo.consigningform.cargoagentid);
            },function (err) {
                var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
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
        vm.openAlert = function (size,index) {
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
