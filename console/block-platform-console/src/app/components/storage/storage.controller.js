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

    angular.module('airs').controller('StorageController', StorageController);

    /** @ngInject */
    function StorageController($stateParams,ApiServer,toastr) {
        /* jshint validthis: true */
        var vm = this;

        vm.statusType = [{id:'a',name:'空闲'},{id:'b',name:'使用中'}];
        vm.status = 'a';

        var storageId = $stateParams.storageId;

        vm.submitAction = submitAction;

        function submitAction() {
            var param = {};

            if (storageId === 'new'){
                ApiServer.containerAdd(param,function (res) {
                    toastr.success('添加成功');
                },function (err) {
                    var errInfo = '添加失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }else {
                ApiServer.containerUpdate(null,function (res) {
                    toastr.success('更新成功');
                },function (err) {
                    var errInfo = '更新失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }

        }
        function getData() {
            if (storageId !== 'new'){
                ApiServer.containerGet(storageId,function (res) {
                    console.log(res);
                },function (err) {
                    var errInfo = '获取数据失败：' + err.statusText + ' (' + err.status +')';
                    toastr.error(errInfo);
                });
            }
        }

        getData();

    }

})();
