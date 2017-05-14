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
 * Created by Otherplayer on 16/7/27.
 */

(function () {
    'use strict';

    angular
        .module('airs')
        .controller('ModalOrder4ClientInstanceCtrl', ModalOrder4ClientInstanceCtrl);

    /** @ngInject */
    function ModalOrder4ClientInstanceCtrl($uibModalInstance,$scope,good) {
        /* jshint validthis: true */
        // var vm = this;

        $scope.good = good;
        $scope.containers = [];
        $scope.containers = good.containers;
        if ($scope.containers && $scope.containers.length > 0){
            $scope.good.containerno = $scope.containers[0].containerno;
        }


        $scope.ok = function () {
            $uibModalInstance.close(good);
        };

        $scope.cancel = function () {
            $uibModalInstance.dismiss('CANCEL');
        };

    }

})();
