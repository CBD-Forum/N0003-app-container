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

    angular.module('airs').controller('orderCarPackingController', orderCarPackingController);

    /** @ngInject */
    function orderCarPackingController($log,toastr,$state,ApiServer,$stateParams,$uibModal,$filter) {
        /* jshint validthis: true */
        var vm = this;

        var info = ApiServer.info();
        var orderId = $stateParams.orderId;

        vm.containers = [];
        vm.isReadOnly = false;
        vm.packingInfo = {carrierid:info.id,orderid:orderId,packinglist:{items:[],dateforpackinggoods:''}};

        vm.submitAction = submitAction;


        function submitAction(type) {

            if (type === 'later'){
                $state.go('app.carorderadd',{orderId:orderId});
                return;
            }

            //判断是否添加了货物
            if (vm.packingInfo.packinglist.items.length == 0){
                toastr.error('请添加货物');
                return;
            }

            vm.packingInfo.packinglist.dateforpackinggoods = $filter('date')(new Date(),'yyyy-MM-dd HH:mm:ss');
            if (type === 'confirm'){
                console.log(vm.packingInfo);
                ApiServer.carOrderPackgoods(vm.packingInfo,function (res) {
                    toastr.success('操作成功');
                    $state.go('app.carorderadd',{orderId:orderId});
                },function (err) {
                    var errInfo = '操作失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                })
            }else{
                history.back();
            }
        }

        function getOrder() {
            ApiServer.orderGet(orderId,function (res) {
                vm.containers = res.data.bookingform.containers;

               console.log(res.data.bookingform.containers);//ContainerNo

            },function (err) {
                var errInfo = '获取集装箱数据失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            });
        }
        getOrder();



        //Model
        vm.addGoodAction = function (size) {
            var modalInstance = $uibModal.open({
                templateUrl: 'orderModalContent.html',
                size: size,
                controller:'ModalOrder4ClientInstanceCtrl',
                resolve: {
                    good: function () {
                        return {name:'',type:'',measurement:'',grossweight:'',containers:vm.containers};
                    }
                }
            });
            modalInstance.result.then(function (param) {
                vm.packingInfo.packinglist.items.push(param);
            }, function () {
                $log.info('Modal dismissed at: ' + new Date());
            });
        };


    }

})();
