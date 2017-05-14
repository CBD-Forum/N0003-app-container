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

    angular.module('airs').controller('orderGoodAddController', orderGoodAddController);

    /** @ngInject */
    function orderGoodAddController($log,$uibModal,$state,ApiServer,$stateParams,toastr,$filter) {
        /* jshint validthis: true */
        var vm = this;

        var orderId = $stateParams.orderId;
        
        

        var info = ApiServer.info();
        vm.orderInfo = {consigningform:{goodslist:[]}};

        vm.acceptAction = acceptAction;


        function acceptAction(type) {
            if (type === 'accept'){

                vm.tipsInfo = {title:'接受',content:'你将接受此订单'};
                vm.openAlert('','accept');

            }else if (type === 'refuse'){

                vm.tipsInfo = {title:'拒绝',content:'确定拒绝接受此订单吗'};
                vm.openAlert('','refuse');

            }else if (type === 'car'){

                $state.go('app.goodordercar',{orderId:orderId});

            }else if (type === 'space'){

                $state.go('app.goodorderspace',{orderId:orderId});

            }else if (type === 'finish'){

                vm.tipsInfo = {title:'完成',content:'确定此订单已经完成了吗'};
                vm.openAlert('','finish');

            }else{

                $state.go('app.goodorder');

            }
        }

        function dealOrderAction(accept) {
            if (accept === 'finish'){
                var now = $filter('date')(new Date(),'yyyy-MM-dd HH:mm:ss');
                var param = {"cargoagentid": info.id, "orderid": orderId, "dateforfinish": now};
                console.log(param);
                ApiServer.goodOrderFinish(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_finished';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else if (accept === 'refuse'){
                var param = {cargoagentid:info.id,orderid:orderId,isorderaccept:false,remark:''};
                ApiServer.goodOrderAccept(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_failed';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else if (accept === 'accept'){
                var param = {cargoagentid:info.id,orderid:orderId,isorderaccept:true,remark:''};
                ApiServer.goodOrderAccept(param,function (res) {
                    toastr.success('操作成功');
                    vm.orderInfo.state = 'order_checked';
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else{

            }
        }

        function getOrderDatas() {
            if (orderId === 'new'){
                return;
            }
            ApiServer.orderGet(orderId,function (res) {
                vm.orderInfo = res.data;
                console.log(vm.orderInfo);
            },function (err) {
                var errInfo = '获取数据失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            });
        }


        getOrderDatas();




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
