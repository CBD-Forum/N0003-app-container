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

    angular.module('airs').controller('ContainersController', ContainersController);

    /** @ngInject */
    function ContainersController($state,$log,$uibModal,ApiServer,toastr) {
        /* jshint validthis: true */
        var vm = this;
        
        var info = ApiServer.info();

        vm.titles = ['集装箱编号','类型','最大重量','毛重','尺寸','操作'];

        vm.items = [];

        vm.displayedCollection = [].concat(vm.items);



        vm.gotoDetail = gotoDetail;

        function gotoDetail(type,index) {
            if (type === 'new'){
                $state.go('app.container',{containerId:'new'});
            } else if (type === 'detail'){
                $state.go('app.container',{containerId:vm.items[index].id});
            }
        }
        function deleteAction(index) {
            ApiServer.vehicleDelete(vm.items[index].id,function (res) {
                toastr.success('删除成功');
            },function (err) {
                var errInfo = '删除失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            })
        }
        function getDatas() {
            ApiServer.containerGetByOwner(info.id,function (res) {
                vm.items = res.data;
                console.log(res);
            },function (err) {
                var errInfo = '获取信息失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            });
        }
        getDatas();


        //Model
        vm.tipsInfo = {title:'删除',content:'确定删除吗？'};
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
                deleteAction(index);
            }, function () {
                $log.info('Modal dismissed at: ' + new Date());
            });
        };




    }

})();
