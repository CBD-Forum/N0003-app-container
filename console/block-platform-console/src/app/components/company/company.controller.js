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

    angular.module('airs').controller('CompanyController', CompanyController);

    /** @ngInject */
    function CompanyController(toastr,ApiServer,$stateParams,$state) {
        /* jshint validthis: true */
        var vm = this;

        var type = $stateParams.type;

        vm.items = [];
        vm.titles = ['公司名称','公司地址','公司网址','组织机构代码','税务登记号码','统一社会信用代码','商业标识','认证'];

        function getDatas() {
            ApiServer.userGetByRoleType(type,function (res) {
                vm.items = res.data;
                vm.displayedCollection = [].concat(vm.items);
            },function (err) {
                var errInfo = '获取数据失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            })
        }
        

        getDatas();
        
        
    }

})();
