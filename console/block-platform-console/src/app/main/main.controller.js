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

    angular.module('airs').controller('MainController', MainController);

    /** @ngInject */
    function MainController($timeout,$translate,$location,ApiServer, $state, toastr,$scope) {
       /* jshint validthis: true */
       var vm = this;
       var url = $location.absUrl();
       var theme_11 = {
          themeID: 11,
          navbarHeaderColor: 'bg-info',
          navbarCollapseColor: 'bg-white',
          asideColor: 'bg-dark b-r',
          headerFixed: true,
          asideFixed: false,
          asideFolded: false,
          asideDock: false,
          container: false
      };

      // config
        $scope.app = {
        name: 'air cc',
        version: '0.0.1',
        // for chart colors
        color: {
          primary: '#7266ba',
          info:    '#23b7e5',
          success: '#27c24c',
          warning: '#fad733',
          danger:  '#f05050',
          light:   '#e8eff0',
          dark:    '#3a3f51',
          black:   '#1c2b36'
        },
        settings: theme_11
      };


      // angular translate
        $scope.lang = { isopen: false };
        $scope.langs = {'en-us':'English', 'zh-cn':'中文'};
        $scope.selectLang = $scope.langs[$translate.proposedLanguage()] || "中文";
        $scope.setLang = function(langKey, $event) {
        // set the current lang
            $scope.selectLang = $scope.langs[langKey];
        // You can change the language during runtime
        $translate.use(langKey);
        $scope.lang.isopen = !$scope.lang.isopen;
      };

        vm.awesomeThings = [];
        vm.classAnimation = '';
        vm.creationDate = 1452231070467;
        vm.showToastr = showToastr;


        function showToastr() {
            toastr.info('Fork <a href="https://github.com/Swiip/generator-gulp-angular" target="_blank"><b>generator-gulp-angular</b></a>');
            vm.classAnimation = '';
        }

        if (ApiServer.isAuthed()){

            if (url.indexOf('#') === -1 || url.indexOf('access') !== -1){
                $timeout(function () {
                    var roleType = ApiServer.roleType();
                    // if (roleType === 'cargoagent'){
                    //     $state.go('app.goodorder');
                    // }else if (roleType === 'carrier'){
                    //     $state.go('app.carorder');
                    // }else if (roleType === 'shipper'){
                    //     $state.go('app.shiporder');
                    // }else{
                        $state.go('app.dashboard');
                    // }
                },10);
            }

        }else{
            $timeout(function () {
                $state.go('access.signin');
            },10);
        }


    }
})();
