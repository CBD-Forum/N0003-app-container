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
(function () {
    'use strict';

    /** @ngInject */
    angular
        .module('airs')
        .config(routeConfig);

    function routeConfig($stateProvider, $urlRouterProvider) {

        // $locationProvider.html5Mode(true);
        $urlRouterProvider
            .otherwise('dashboard');
        $stateProvider


            .state('app.dashboard',{
                url: 'dashboard',
                templateUrl: 'app/components/dashboard/dashboard.html'
            })

            /** LOGIN **/
            .state('access', {
                url: '/access',
                templateUrl: 'signin.html',
                controller: 'SigninController',
                controllerAs: 'vm'
            })
            .state('access.signin', {
                url: '/signin',
                templateUrl: 'app/components/signin/signin.html'
            })
            .state('access.signup',{
                url: '/signup',
                templateUrl: 'app/components/signin/signup.html'
            })


            .state('app.profile', {
                url: 'user/profile',
                templateUrl: 'app/components/profile/profile.html'
            })
            .state('app.company', {
                url: 'regular/company?type',
                templateUrl: 'app/components/company/company.html'
            })

            /** ACCOUNT **/

            //////用户
            .state('app.userorder', {
                url: 'regular/order',
                templateUrl: 'app/components/order/user.order.html'
            })
            .state('app.userorderadd', {
                url: 'regular/order/:orderId',
                templateUrl: 'app/components/order/user.order.add.html'
            })


            //////货代公司
            .state('app.goodorder', {
                url: 'good/order',
                templateUrl: 'app/components/order/good.order.html'
            })
            .state('app.goodorderadd', {
                url: 'good/order/:orderId',
                templateUrl: 'app/components/order/good.order.add.html'
            })
            .state('app.goodordercar', {
                url: 'good/order/:orderId/car',
                templateUrl: 'app/components/order/good.order.car.html'
            })
            .state('app.goodorderspace', {
                url: 'good/order/:orderId/space',
                templateUrl: 'app/components/order/good.order.space.html'
            })
            .state('app.goodstorages', {
                url: 'good/storages',
                templateUrl: 'app/components/storage/storages.html'
            })
            .state('app.goodstorage', {
                url: 'good/storage',
                templateUrl: 'app/components/storage/storage.html'
            })

            //////拖车公司
            .state('app.carorder', {
                url: 'car/order',
                templateUrl: 'app/components/order/car.order.html'
            })
            .state('app.carorderadd', {
                url: 'car/order/:orderId',
                templateUrl: 'app/components/order/car.order.add.html'
            })
            .state('app.carorderaddpacking', {
                url: 'car/order/packing/:orderId',
                templateUrl: 'app/components/order/car.order.packing.html'
            })
            .state('app.vehicles', {
                url: 'car/vehicles',
                templateUrl: 'app/components/vehicle/vehicles.html'
            })
            .state('app.vehicle', {
                url: 'car/vehicle/:vehicleId',
                templateUrl: 'app/components/vehicle/vehicle.html'
            })
            .state('app.transporttask', {
                url: 'car/transporttask',
                templateUrl: 'app/components/transporttask/transporttask.html'
            })

            //////船运公司
            .state('app.shiporder', {
                url: 'ship/order',
                templateUrl: 'app/components/order/ship.order.html'
            })
            .state('app.shiporderadd', {
                url: 'ship/order/:orderId',
                templateUrl: 'app/components/order/ship.order.add.html'
            })
            .state('app.containers', {
                url: 'ship/containers',
                templateUrl: 'app/components/container/containers.html'
            })
            .state('app.container', {
                url: 'ship/container/:containerId',
                templateUrl: 'app/components/container/container.html'
            })
            .state('app.shippingschedule', {
                url: 'ship/shippingschedule',
                templateUrl: 'app/components/shippingschedule/shippingschedule.html'
            })
            .state('app.shippingscheduleadd', {
                url: 'ship/shippingschedule/:shippingScheduleId',
                templateUrl: 'app/components/shippingschedule/shippingschedule.add.html'
            })


            .state('app', {
                //abstract: true,
                url: '/',
                controller: 'MainController',
                controllerAs: 'main',
                templateUrl: 'app/main/main.html'
            })

        ;
    }

})();
