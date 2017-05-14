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
(function() {
  'use strict';

  describe('service webDevTec', function() {
    var webDevTec;

    beforeEach(module('airs'));
    beforeEach(inject(function(_webDevTec_) {
      webDevTec = _webDevTec_;
    }));

    it('should be registered', function() {
      expect(webDevTec).not.toEqual(null);
    });

    describe('getTec function', function() {
      it('should exist', function() {
        expect(webDevTec.getTec).not.toEqual(null);
      });

      it('should return array of object', function() {
        var data = webDevTec.getTec();
        expect(data).toEqual(jasmine.any(Array));
        expect(data[0]).toEqual(jasmine.any(Object));
        expect(data.length > 5).toBeTruthy();
      });
    });
  });
})();
