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

    angular.module('airs').controller('SigninController', SigninController);

    /** @ngInject */
    function SigninController(logger,toastr,StorageService,$timeout,$state,constdata,$rootScope,iotUtil,$translate,ApiServer) {
        /* jshint validthis: true */
        var vm = this;
        var authorizationKey = constdata.token;
        var userInfo = constdata.informationKey;
        //语言
        var langChi = '中文';
        var langEng = 'English';
        var userLanguage = window.localStorage.userLanguage;

        vm.user = {username:'',password:''};
        vm.isLogining = false;


        vm.login = loginAction;
        vm.logout = logoutAction;
        vm.username = username;
        vm.gotoRegisterAction = gotoRegisterAction;
        
        function gotoRegisterAction() {
            $state.go('access.signup');
        }
        function loginAction() {

            if (vm.user.username.length == 0 || vm.user.password.length == 0){
                toastr.error('请输入用户名和密码');
                return;
            }

            vm.isLogining = true;
            StorageService.remove(constdata.token);

            ApiServer.userLogin(vm.user.username,vm.user.password,function (response) {

                console.log('login success');
                var result = response.data;

                var userId = result.userid;
                var sessionId = result.sessionid;
                var token = result.token;

                var sessionInfo = {userid:userId,sessionid:sessionId,token:token};

                StorageService.put(authorizationKey,sessionInfo,24 * 7 * 60 * 60);//3 天过期

                ApiServer.userGet(userId,function (infoResponse) {
                    console.log(infoResponse.data);
                    $rootScope.userInfo = infoResponse.data;
                    StorageService.put(userInfo,infoResponse.data,24 * 3 * 60 * 60);

                    var appGo = 'app.dashboard';
                    var roleType = infoResponse.data.role;
                    // if (roleType === 'cargoagent'){
                    //     appGo = 'app.goodorder';
                    // }else if (roleType === 'carrier'){
                    //     appGo = 'app.carorder';
                    // }else if (roleType === 'shipper'){
                    //     appGo = 'app.shiporder';
                    // }else{
                    //     appGo = 'app.dashboard';
                    // }

                    $rootScope.$on('$locationChangeSuccess',function(){//返回前页时，刷新前页
                        parent.location.reload();
                    });

                    $state.go(appGo);
                },function (infoError) {
                    var errInfo = '登录失败：' + infoError.statusText + ' (' + infoError.status +')';
                    toastr.error(errInfo);
                    vm.isLogining = false;
                });
            },function (err) {
                var errInfo = '登录失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
                vm.isLogining = false;
            });

        }
        function username() {
            var information = StorageService.get(constdata.informationKey);
            return information.username;
        }
        
        function logoutAction() {
            $timeout(function () {
                StorageService.clear(authorizationKey);
                StorageService.clear(userInfo);
                StorageService.clear(constdata.token);
            },60);
            $state.go('access.signin');
        }

        //切换语言
        userLanguage == 'zh-cn' ? vm.langChoosen = langChi : vm.langChoosen = langEng
        userLanguage == 'zh-cn' ? vm.langLeft = langEng : vm.langLeft = langChi
        vm.toggleLang = function(lang) {
            vm.langChoosen = (vm.langChoosen == langChi) ? langEng : langChi
            vm.langLeft = (vm.langLeft == langChi) ? langEng : langChi;
            // console.log(lang);
            lang == langEng ? $translate.use('en-us') : $translate.use('zh-cn');
            lang == langEng ? window.localStorage.userLanguage='en-us' :  window.localStorage.userLanguage='zh-cn'
            // window.location.reload();
        }

    }

})();
