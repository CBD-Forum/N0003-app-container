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
 * Created by Otherplayer on 16/7/25.
 */
(function () {
    'use strict';

    angular.module('airs').controller('SignupController', SignupController);

    /** @ngInject */
    function SignupController(toastr,StorageService,$http,$state,constdata,iotUtil) {
        /* jshint validthis: true */
        var vm = this;
        vm.user = {role:'regularclient'};
        vm.roles = [
            {role:'regularclient',name:'普通用户',company:false},
            {role:'cargoagent',name:'货代公司',company:true},
            {role:'carrier',name:'拖车公司',company:true},
            {role:'shipper',name:'船运公司',company:true}
            ];
        vm.gotoLoginAction = gotoLoginAction;
        vm.roleTypeChangedAction = roleTypeChangedAction;
        vm.registerAction = registerAction;
        vm.password = '';


        function registerAction() {
            register(vm.user,vm.password,function (res) {
                toastr.success('注册成功');
                gotoLoginAction();
            });
        }
        function gotoLoginAction() {
            $state.go('access.signin');
        }
        function roleTypeChangedAction() {
            console.log(vm.role);
        }


        function register(params,password,successHandler) {
            params.password = iotUtil.sha256(password);
            // var header = { "password": iotUtil.sha256(password) };
            $http({
                url: getBaseUrl(),
                method: 'POST',
                data: params
                // headers: header
            }).success(function (data) {
                successHandler(data);
            }).error(function (err) {
                toastr.error(err.Error);
            });
        }
        function getBaseUrl() {
            if(constdata.debugMode){
                return constdata.apiHost_OFFLINE + constdata.api.user;
            }else {
                return constdata.apiHost_ONLINE + constdata.api.user;
            }
        }




    }

})();
