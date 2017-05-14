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

    angular.module('airs').controller('TransportTaskController', TransportTaskController);

    /** @ngInject */
    function TransportTaskController($state,toastr,ApiServer) {
        /* jshint validthis: true */
        var vm = this;

        var info = ApiServer.info();

        vm.titles = ['车辆编号','驾驶人','联系方式','驾驶证号','起始','到达'];

        vm.items = [];

        function getData() {

            ApiServer.transportTaskGetByOwner(info.id,function (res) {
                console.log(res);
                vm.items = res.data;
            },function (err) {
                var errInfo = '获取信息失败：' + err.statusText + ' (' + err.status +')';
                toastr.error(errInfo);
            });
        }

        getData();

    }

})();
