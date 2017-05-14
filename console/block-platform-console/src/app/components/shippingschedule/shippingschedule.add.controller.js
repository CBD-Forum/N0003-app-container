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

    angular.module('airs').controller('ShippingScheduleAddController', ShippingScheduleAddController);

    /** @ngInject */
    function ShippingScheduleAddController($stateParams,ApiServer,toastr,$scope,$filter,$state,$timeout) {
        /* jshint validthis: true */
        var vm = this;

        var shippingScheduleId = $stateParams.shippingScheduleId;
        var info = ApiServer.info();
        var roleType = info.role;

        vm.isReadOnly = (roleType !== 'shipper');
        vm.statusType = [{id:'free',name:'空闲'},{id:'inuse',name:'使用中'}];
        vm.status = 'free';

        vm.shippingInfo = {ownerid:info.id,vessel:{ownerid:info.id,space:{spaceused:[],spaceunused:[]}}};


        vm.submitAction = submitAction;
        vm.cancelAction = cancelAction;

        function submitAction() {

            if (vm.shippingInfo.vessel.space.restnum > vm.shippingInfo.vessel.space.totalnum){
                toastr.error('未使用舱位数不能大于总舱位数');
                return;
            }else{
                vm.shippingInfo.vessel.space.spaceunused = [];
                vm.shippingInfo.vessel.space.spaceused = [];
                for(var i = 0; i < vm.shippingInfo.vessel.space.totalnum; i++){
                    if (i < vm.shippingInfo.vessel.space.restnum){
                        vm.shippingInfo.vessel.space.spaceunused.push('b'+i);
                    }else{
                        vm.shippingInfo.vessel.space.spaceused.push('a'+i);
                    }
                }
            }
            // console.log(vm.shippingInfo);

            vm.shippingInfo.arrivaldate = $filter('date')($scope.arrivaldate,'yyyy-MM-dd HH:mm:ss');
            vm.shippingInfo.departuredate = $filter('date')($scope.departuredate,'yyyy-MM-dd HH:mm:ss');
            if (shippingScheduleId === 'new'){
                console.log(vm.shippingInfo);
                ApiServer.shippingScheduleAdd(vm.shippingInfo,function (res) {
                    toastr.success('添加成功');
                    $timeout(function () {
                        $state.go('app.shippingschedule');
                    },2000);
                },function (err) {
                    var errInfo = '添加失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }else {
                ApiServer.shippingScheduleUpdate(vm.shippingInfo,function (res) {
                    toastr.success('更新成功');
                },function (err) {
                    var errInfo = '更新失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }

        }
        function cancelAction() {
            $state.go('app.shippingschedule');
        }
        function getData() {
            if (shippingScheduleId !== 'new'){
                ApiServer.shippingScheduleGet(shippingScheduleId,function (res) {
                    vm.shippingInfo = res.data;
                },function (err) {
                    var errInfo = '获取数据失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }
        }

        getData();








        $scope.today = function() {
            $scope.departuredate = new Date();
            $scope.arrivaldate = new Date();
        };
        $scope.today();

        $scope.clear = function () {
            $scope.departuredate = null;
            $scope.arrivaldate = null;
        };

        // Disable weekend selection
        $scope.disabled = function(date, mode) {
            return ( mode === 'day' && ( date.getDay() === 0 || date.getDay() === 6 ) );
        };

        $scope.toggleMin = function() {
            $scope.minDate = $scope.minDate ? null : new Date();
        };
        $scope.toggleMin();

        $scope.open = function($event) {
            $event.preventDefault();
            $event.stopPropagation();

            $scope.opened = true;
        };
        $scope.open2 = function($event) {
            $event.preventDefault();
            $event.stopPropagation();

            $scope.opened2 = true;
        };

        $scope.dateOptions = {
            formatYear: 'yy',
            startingDay: 1,
            class: 'datepicker'
        };

        $scope.formats = ['dd-MMMM-yyyy', 'yyyy/MM/dd', 'dd.MM.yyyy', 'shortDate'];
        $scope.format = $scope.formats[0];

    }

})();
