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
 * app main controller
 */

(function () {
    'use strict';

    angular.module('airs').controller('AsideController', AsideController);

    /** @ngInject */
    function AsideController(ApiServer,$state) {
        /* jshint validthis: true */
        var vm = this;
        var info = ApiServer.info();

        var height = document.body.clientHeight + 'px';
        vm.navStyle = {'height':height};
        
        vm.clearAllMessageAction = clearAllMessageAction;
        vm.logoutAction = logoutAction;

        vm.title = 'IntelligenceContainerBlock';
        vm.messages = [];
        vm.roleType = ApiServer.roleType();
        if (vm.roleType === 'regularclient'){
            vm.title = '用户系统';
        }else if (vm.roleType === 'cargoagent'){
            vm.title = '货代公司管理系统';
        }else if (vm.roleType === 'carrier'){
            vm.title = '拖车公司管理系统';
        }else if (vm.roleType === 'shipper'){
            vm.title = '船运公司管理系统';
        }

        vm.infomation = ApiServer.info();
        
        function clearAllMessageAction() {
            vm.messages = [];
            for (var i = 0; i < vm.messages.length; i++){
                var msg = vm.messages[i];
                ApiServer.messageDelete(msg.messageid);
            }
        }
        
        ApiServer.messageGetByUserId(info.id,function (res) {
            vm.messages = res.data;
            if (vm.messages.length > 5){
                vm.messages = vm.messages.slice(0,5);
            }
        },function (err) {

        })
        
        function logoutAction() {
            ApiServer.logoutAction();
        }
        
        


    }
})();
