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

    angular.module('airs').controller('ProfileController', ProfileController);

    /** @ngInject */
    function ProfileController($scope,ApiServer,toastr,constdata,$state) {
        /* jshint validthis: true */
        var vm = this;

        var userInfo = constdata.informationKey;
        var info = ApiServer.info();
        vm.roleType = info.role;
        vm.submitAction = submitAction;
        vm.cancelAction = cancelAction;

        vm.user = {role:info.role};

        function submitAction() {
            ApiServer.userUpdate(vm.user,function (res) {
                //更新本地信息
                // StorageService.put(userInfo,vm.user,24 * 3 * 60 * 60);
                toastr.success('更新成功');

            },function (err) {
                var errInfo = '更新失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            })
        }
        function cancelAction() {
            $state.go('app.dashboard');
        }

        function getDatas() {
            ApiServer.userGet(info.id,function (res) {
                vm.user = res.data;
            },function (err) {
                var errInfo = '获取数据失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            })
        }

        getDatas();


    }

})();
