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

    angular.module('airs').controller('ContainerController', ContainerController);

    /** @ngInject */
    function ContainerController($stateParams,ApiServer,toastr,$state,$timeout) {
        /* jshint validthis: true */
        var vm = this;

        var info = ApiServer.info();

        vm.containerInfo = {ownerid:info.id,type:'freight',location:'北京'};

        vm.statusType = [{id:'free',name:'空闲'},{id:'inuse',name:'使用中'}];
        vm.status = 'free';

        var containerId = $stateParams.containerId;

        vm.submitAction = submitAction;
        vm.cancelAction = cancelAction;

        function submitAction() {

            if (containerId === 'new'){
                ApiServer.containerAdd(vm.containerInfo,function (res) {
                    toastr.success('添加成功');
                    $timeout(function () {
                        $state.go('app.containers');
                    },2000);

                },function (err) {
                    var errInfo = '添加失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }else {
                ApiServer.containerUpdate(vm.containerInfo,function (res) {
                    toastr.success('更新成功');
                },function (err) {
                    var errInfo = '更新失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }

        }
        function cancelAction() {
            $state.go('app.containers');
        }
        function getData() {
            if (containerId !== 'new'){
                ApiServer.containerGet(containerId,function (res) {
                    console.log(res);
                    vm.containerInfo = res.data;
                },function (err) {
                    var errInfo = '获取数据失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }
        }

        getData();

    }

})();
