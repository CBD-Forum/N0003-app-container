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
 *
 * Run blocks - get executed after the injector is created and are used to kickstart the application.
 * Only instances and constants can be injected into run blocks. This is to prevent further
 * system configuration during application run time.
 *
 */

(function () {

    // 'use strict';

    angular.module('airs').run(runBlock);

    /** @ngInject */
    function runBlock($log) {
        $log.debug('Hold on, Air Start Fly.');
    }


})();
