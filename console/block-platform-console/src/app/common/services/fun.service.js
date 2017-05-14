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
 * Created by Otherplayer on 2017/1/9.
 */
(function() {
    'use strict';

    /**
     * 首字母大写
     *
     */
    angular
        .module('airs')
        .filter('capitalize', capitalize);

    /** @ngInject */
    function capitalize() {
        return function(input){        //input是我们传入的字符串
            if (input) {
                return input[0].toUpperCase() + input.slice(1);
            }
        };
    };


    angular
        .module('airs')
        .filter('orderdesc', orderdesc);

    /** @ngInject */
    function orderdesc() {
        return function(input){        //input是我们传入的字符串

            if (input) {
                if (input === 'order_init'){
                    return '等待确认中';
                }else if (input === 'order_created'){
                    return '订单已接受';
                }else if (input === 'order_space_booked'){
                    return '订舱完成';
                }else if (input === 'order_vehicle_booked'){
                    return '订车完成';
                }else if (input === 'order_empty_container_fetched'){
                    return '集装箱已指派';
                }else if (input === 'order_goods_packed'){
                    return '货物已经装箱';
                }else if (input === 'order_yard_arrived'){
                    return '集装箱抵达堆场';
                }else if (input === 'order_goods_loaded'){
                    return '货物已上船';
                }else if (input === 'order_goods_shipping'){
                    return '海运中';
                }else if (input === 'order_goods_arrived'){
                    return '抵达目的港口';
                }else if (input === 'order_goods_delivered'){
                    return '货物已发出';
                }else if (input === 'order_goods_received'){
                    return '收到货物';
                }else if (input === 'order_finished'){
                    return '订单完成';
                }else if (input === 'order_failed'){
                    return '订单失败';
                }else{
                    return input;
                }
            }
        };
    };


    angular
        .module('airs')
        .filter('descTime', descTime);

        /** @ngInject */
    function descTime() {

        return service;


        function service(input) {

            var result = $filter('date')(input,'yyyy-MM-dd HH:mm:ss');
            return result;
        }

    };






    /**
     * 字符串转换为数字
     *
     */
    angular
        .module('airs')
        .filter('num', num);

    /** @ngInject */
    function num() {
        return function(input){        //input是我们传入的字符串
            if (input) {
                return Number(input);
            }
        };
    };


    angular
        .module('airs')
        .factory('iotUtil', iotUtil);

    /** @ngInject */
    function iotUtil(StorageService) {
        var service = {
            uuid : uuid,
            isNull : isNull,
            pagesize : pagesize,
            htmlToPlaintext : htmlToPlaintext,
            getKeyValueFromURL : getKeyValueFromURL,
            dateDesc: dateDesc,
            sha256: SHA256
        };
        return service;

        ////////////////////

        function dateDesc(input) {
            var result = $filter('date')(input,'yyyy-MM-dd HH:mm:ss');
            return result;
        }
        function uuid() {
            var s = [];
            var hexDigits = "0123456789abcdef";
            for (var i = 0; i < 36; i++) {
                s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
            }
            s[14] = "4";  // bits 12-15 of the time_hi_and_version field to 0010
            s[19] = hexDigits.substr((s[19] & 0x3) | 0x8, 1);  // bits 6-7 of the clock_seq_hi_and_reserved to 01
            s[8] = s[13] = s[18] = s[23] = "";

            var uuid = s.join("");
            return uuid;
        }
        function isNull( str ) {
            if ( !str || str == "" ) return true;
            var regu = "^[ ]+$";
            var re = new RegExp(regu);
            return re.test(str);
        }
        function pagesize() {
            var hgqInfo = constdata.informationKey;
            var infomation = StorageService.get(hgqInfo);
            var pageSizeKey = 'airspc.pagesize.' + infomation.username;
            var tempPageSize = StorageService.get(pageSizeKey);
            if (tempPageSize && tempPageSize != 'undefined'){
                return tempPageSize;
            }
            return 10;
        }
        function htmlToPlaintext(text) {
            return text ? String(text).replace(/<[^>]+>/gm, '') : '';
        }
        function getKeyValueFromURL(key,urlstring) {
            var tempData = urlstring.split(key + '=');
            if (tempData.length == 1){
                return 'undefined';
            }
            var lastPath = tempData.pop();
            tempData = lastPath.split('&');
            var result = tempData[0];
            return result;
        }
        function SHA256(s){

            var chrsz   = 8;
            var hexcase = 0;

            function safe_add (x, y) {
                var lsw = (x & 0xFFFF) + (y & 0xFFFF);
                var msw = (x >> 16) + (y >> 16) + (lsw >> 16);
                return (msw << 16) | (lsw & 0xFFFF);
            }

            function S (X, n) { return ( X >>> n ) | (X << (32 - n)); }
            function R (X, n) { return ( X >>> n ); }
            function Ch(x, y, z) { return ((x & y) ^ ((~x) & z)); }
            function Maj(x, y, z) { return ((x & y) ^ (x & z) ^ (y & z)); }
            function Sigma0256(x) { return (S(x, 2) ^ S(x, 13) ^ S(x, 22)); }
            function Sigma1256(x) { return (S(x, 6) ^ S(x, 11) ^ S(x, 25)); }
            function Gamma0256(x) { return (S(x, 7) ^ S(x, 18) ^ R(x, 3)); }
            function Gamma1256(x) { return (S(x, 17) ^ S(x, 19) ^ R(x, 10)); }

            function core_sha256 (m, l) {
                var K = new Array(0x428A2F98, 0x71374491, 0xB5C0FBCF, 0xE9B5DBA5, 0x3956C25B, 0x59F111F1, 0x923F82A4, 0xAB1C5ED5, 0xD807AA98, 0x12835B01, 0x243185BE, 0x550C7DC3, 0x72BE5D74, 0x80DEB1FE, 0x9BDC06A7, 0xC19BF174, 0xE49B69C1, 0xEFBE4786, 0xFC19DC6, 0x240CA1CC, 0x2DE92C6F, 0x4A7484AA, 0x5CB0A9DC, 0x76F988DA, 0x983E5152, 0xA831C66D, 0xB00327C8, 0xBF597FC7, 0xC6E00BF3, 0xD5A79147, 0x6CA6351, 0x14292967, 0x27B70A85, 0x2E1B2138, 0x4D2C6DFC, 0x53380D13, 0x650A7354, 0x766A0ABB, 0x81C2C92E, 0x92722C85, 0xA2BFE8A1, 0xA81A664B, 0xC24B8B70, 0xC76C51A3, 0xD192E819, 0xD6990624, 0xF40E3585, 0x106AA070, 0x19A4C116, 0x1E376C08, 0x2748774C, 0x34B0BCB5, 0x391C0CB3, 0x4ED8AA4A, 0x5B9CCA4F, 0x682E6FF3, 0x748F82EE, 0x78A5636F, 0x84C87814, 0x8CC70208, 0x90BEFFFA, 0xA4506CEB, 0xBEF9A3F7, 0xC67178F2);
                var HASH = new Array(0x6A09E667, 0xBB67AE85, 0x3C6EF372, 0xA54FF53A, 0x510E527F, 0x9B05688C, 0x1F83D9AB, 0x5BE0CD19);
                var W = new Array(64);
                var a, b, c, d, e, f, g, h, i, j;
                var T1, T2;

                m[l >> 5] |= 0x80 << (24 - l % 32);
                m[((l + 64 >> 9) << 4) + 15] = l;

                for ( var i = 0; i<m.length; i+=16 ) {
                    a = HASH[0];
                    b = HASH[1];
                    c = HASH[2];
                    d = HASH[3];
                    e = HASH[4];
                    f = HASH[5];
                    g = HASH[6];
                    h = HASH[7];

                    for ( var j = 0; j<64; j++) {
                        if (j < 16) W[j] = m[j + i];
                        else W[j] = safe_add(safe_add(safe_add(Gamma1256(W[j - 2]), W[j - 7]), Gamma0256(W[j - 15])), W[j - 16]);

                        T1 = safe_add(safe_add(safe_add(safe_add(h, Sigma1256(e)), Ch(e, f, g)), K[j]), W[j]);
                        T2 = safe_add(Sigma0256(a), Maj(a, b, c));

                        h = g;
                        g = f;
                        f = e;
                        e = safe_add(d, T1);
                        d = c;
                        c = b;
                        b = a;
                        a = safe_add(T1, T2);
                    }

                    HASH[0] = safe_add(a, HASH[0]);
                    HASH[1] = safe_add(b, HASH[1]);
                    HASH[2] = safe_add(c, HASH[2]);
                    HASH[3] = safe_add(d, HASH[3]);
                    HASH[4] = safe_add(e, HASH[4]);
                    HASH[5] = safe_add(f, HASH[5]);
                    HASH[6] = safe_add(g, HASH[6]);
                    HASH[7] = safe_add(h, HASH[7]);
                }
                return HASH;
            }

            function str2binb (str) {
                var bin = Array();
                var mask = (1 << chrsz) - 1;
                for(var i = 0; i < str.length * chrsz; i += chrsz) {
                    bin[i>>5] |= (str.charCodeAt(i / chrsz) & mask) << (24 - i%32);
                }
                return bin;
            }

            function Utf8Encode(string) {
                string = string.replace(/\r\n/g,"\n");
                var utftext = "";

                for (var n = 0; n < string.length; n++) {

                    var c = string.charCodeAt(n);

                    if (c < 128) {
                        utftext += String.fromCharCode(c);
                    }
                    else if((c > 127) && (c < 2048)) {
                        utftext += String.fromCharCode((c >> 6) | 192);
                        utftext += String.fromCharCode((c & 63) | 128);
                    }
                    else {
                        utftext += String.fromCharCode((c >> 12) | 224);
                        utftext += String.fromCharCode(((c >> 6) & 63) | 128);
                        utftext += String.fromCharCode((c & 63) | 128);
                    }

                }

                return utftext;
            }

            function binb2hex (binarray) {
                var hex_tab = hexcase ? "0123456789ABCDEF" : "0123456789abcdef";
                var str = "";
                for(var i = 0; i < binarray.length * 4; i++) {
                    str += hex_tab.charAt((binarray[i>>2] >> ((3 - i%4)*8+4)) & 0xF) +
                        hex_tab.charAt((binarray[i>>2] >> ((3 - i%4)*8  )) & 0xF);
                }
                return str;
            }

            s = Utf8Encode(s);
            return binb2hex(core_sha256(str2binb(s), s.length * chrsz));

        }
    }


    angular
        .module('airs')
        .factory('deepcopy', deepcopy);

    /** @ngInject */
    function deepcopy() {
        var service = {
            copy : copy
        };
        return service;

        ////////////////////

        function copy(source) {
            var result={};
            for (var key in source) {
                result[key] = typeof source[key]==='object'? copy(source[key]): source[key];
            }
            return result;
        }
    }


    // /**
    //  * 过滤器:国际化
    //  *
    //  */
    // angular
    //     .module('airs')
    //     .filter("T", T);
    //
    // /** @ngInject */
    // function T($translate) {
    //     return function(key) {
    //         if(key){
    //             return $translate.instant(key);
    //         }
    //         return key;
    //     };
    // }


    /**
     * 服务:国际化
     *
     */
    angular
        .module('airs')
        .factory("i18n", i18n);

    /** @ngInject */
    function i18n($translate) {
        var service = {
            t: translate,
            value : translate
        };
        return service;

        ////////////////////

        function translate(key) {
            if(key){
                return $translate.instant(key);
            }
            return key;
        }
    }


    /**
     * 日志:logger
     *
     */
    angular
        .module('airs')
        .factory('logger', logger);

    /** @ngInject */
    function logger(constdata) {

        var mod = {
            log: 1,
            debug: 2,
            info: 4,
            warn: 8,
            error: 16
        };
        var debugMode = constdata.debugMode;
        var logLevel = parseInt(constdata.logLevel, 2);


        var service = {
            log: log,
            err: err,
            error: err,
            warn: warn,
            info: info,
            debug: debug
        };

        return service;


        ///////////////////


        function log(args) {
            if (__showLevel('log')) {
                console.log(args);
            }
        }

        function err(args) {
            if (__showLevel('error')) {
                console.error(args);
            }
        }

        function warn(args) {
            if (__showLevel('warn')) {
                console.warn(args);
            }
        }

        function info(args) {
            if (__showLevel('info')) {
                console.info(args);
            }
        }

        function debug(args) {
            if (__showLevel('debug')) {
                console.log(args);
            }
        }

        function __showLevel(name) {
            var isValid = logLevel & mod[name];
            return debugMode && isValid;
        }
    }

})();

Array.prototype.indexOf = function(val) {
    for (var i = 0; i < this.length; i++) {
        if (this[i] == val) return i;
    }
    return -1;
};
Array.prototype.remove = function(val) {
    var index = this.indexOf(val);
    if (index > -1) {
        this.splice(index, 1);
    }
};
Array.prototype.clone = function(){
    return this.slice(0);
};


// /*
//  用途：检查输入字符串是否只由英文字母和数字和文字组成
//  输入：
//  s：字符串
//  返回：
//  如果通过验证返回true,否则返回false
//
//  */
// function noSpecialSymbols(s) {
//     var regu = "^[\a-\z\A-\Z0-9\u4E00-\u9FA5\_\:\.\{\} ]+$";
//     var re = new RegExp(regu);
//     if (re.test(s)) {
//         return true;
//     } else {
//         return false;
//     }
// }
//
// /*
//  用途：检查输入字符串是否只由英文字母和数字和文字组成
//  输入：
//  s：字符串
//  返回：
//  如果通过验证返回true,否则返回false
//
//  */
// function isNumber_Letter(s) {
//     var regu = "^[\a-\z\A-\Z0-9\_]+$";
//     var re = new RegExp(regu);
//     if (re.test(s)) {
//         return true;
//     } else {
//         return false;
//     }
// }
// /*
//  用途：检查输入字符串是否为空或者全部都是空格
//  输入：str
//  返回：
//  如果全是空返回true,否则返回false
//  */
// function isNull( str )
// {
//     if ( !str || str == "" ) return true;
//     var regu = "^[ ]+$";
//     var re = new RegExp(regu);
//     return re.test(str);
// }


(function($) {

    $.fn.extend({
        slimScroll: function(options) {

            var defaults = {

                // width in pixels of the visible scroll area
                width : '240px',

                // height in pixels of the visible scroll area
                height : '150px',

                // width in pixels of the scrollbar and rail
                size : '7px',

                // scrollbar color, accepts any hex/color value
                color: '#000',

                // scrollbar position - left/right
                position : 'right',

                // distance in pixels between the side edge and the scrollbar
                distance : '1px',

                // default scroll position on load - top / bottom / $('selector')
                start : 'top',

                // sets scrollbar opacity
                opacity : .4,

                // enables always-on mode for the scrollbar
                alwaysVisible : false,

                // check if we should hide the scrollbar when user is hovering over
                disableFadeOut : false,

                // sets visibility of the rail
                railVisible : false,

                // sets rail color
                railColor : '#333',

                // sets rail opacity
                railOpacity : .2,

                // whether  we should use jQuery UI Draggable to enable bar dragging
                railDraggable : true,

                // defautlt CSS class of the slimscroll rail
                railClass : 'slimScrollRail',

                // defautlt CSS class of the slimscroll bar
                barClass : 'slimScrollBar',

                // defautlt CSS class of the slimscroll wrapper
                wrapperClass : 'slimScrollDiv',

                // check if mousewheel should scroll the window if we reach top/bottom
                allowPageScroll : false,

                // scroll amount applied to each mouse wheel step
                wheelStep : 20,

                // scroll amount applied when user is using gestures
                touchScrollStep : 200,

                // sets border radius
                borderRadius: '7px',

                // sets border radius of the rail
                railBorderRadius : '7px'
            };

            var o = $.extend(defaults, options);

            // do it for every element that matches selector
            this.each(function(){

                var isOverPanel, isOverBar, isDragg, queueHide, touchDif,
                    barHeight, percentScroll, lastScroll,
                    divS = '<div></div>',
                    minBarHeight = 30,
                    releaseScroll = false;

                // used in event handlers and for better minification
                var me = $(this);

                // ensure we are not binding it again
                if (me.parent().hasClass(o.wrapperClass))
                {
                    // start from last bar position
                    var offset = me.scrollTop();

                    // find bar and rail
                    bar = me.siblings('.' + o.barClass);
                    rail = me.siblings('.' + o.railClass);

                    getBarHeight();

                    // check if we should scroll existing instance
                    if ($.isPlainObject(options))
                    {
                        // Pass height: auto to an existing slimscroll object to force a resize after contents have changed
                        if ( 'height' in options && options.height == 'auto' ) {
                            me.parent().css('height', 'auto');
                            me.css('height', 'auto');
                            var height = me.parent().parent().height();
                            me.parent().css('height', height);
                            me.css('height', height);
                        } else if ('height' in options) {
                            var h = options.height;
                            me.parent().css('height', h);
                            me.css('height', h);
                        }

                        if ('scrollTo' in options)
                        {
                            // jump to a static point
                            offset = parseInt(o.scrollTo);
                        }
                        else if ('scrollBy' in options)
                        {
                            // jump by value pixels
                            offset += parseInt(o.scrollBy);
                        }
                        else if ('destroy' in options)
                        {
                            // remove slimscroll elements
                            bar.remove();
                            rail.remove();
                            me.unwrap();
                            return;
                        }

                        // scroll content by the given offset
                        scrollContent(offset, false, true);
                    }

                    return;
                }
                else if ($.isPlainObject(options))
                {
                    if ('destroy' in options)
                    {
                        return;
                    }
                }

                // optionally set height to the parent's height
                o.height = (o.height == 'auto') ? me.parent().height() : o.height;

                // wrap content
                var wrapper = $(divS)
                    .addClass(o.wrapperClass)
                    .css({
                        position: 'relative',
                        overflow: 'hidden',
                        width: o.width,
                        height: o.height
                    });

                // update style for the div
                me.css({
                    overflow: 'hidden',
                    width: o.width,
                    height: o.height
                });

                // create scrollbar rail
                var rail = $(divS)
                    .addClass(o.railClass)
                    .css({
                        width: o.size,
                        height: '100%',
                        position: 'absolute',
                        top: 0,
                        display: (o.alwaysVisible && o.railVisible) ? 'block' : 'none',
                        'border-radius': o.railBorderRadius,
                        background: o.railColor,
                        opacity: o.railOpacity,
                        zIndex: 90
                    });

                // create scrollbar
                var bar = $(divS)
                    .addClass(o.barClass)
                    .css({
                        background: o.color,
                        width: o.size,
                        position: 'absolute',
                        top: 0,
                        opacity: o.opacity,
                        display: o.alwaysVisible ? 'block' : 'none',
                        'border-radius' : o.borderRadius,
                        BorderRadius: o.borderRadius,
                        MozBorderRadius: o.borderRadius,
                        WebkitBorderRadius: o.borderRadius,
                        zIndex: 99
                    });

                // set position
                var posCss = (o.position == 'right') ? { right: o.distance } : { left: o.distance };
                rail.css(posCss);
                bar.css(posCss);

                // wrap it
                me.wrap(wrapper);

                // append to parent div
                me.parent().append(bar);
                me.parent().append(rail);

                // make it draggable and no longer dependent on the jqueryUI
                if (o.railDraggable){
                    bar.bind("mousedown", function(e) {
                        var $doc = $(document);
                        isDragg = true;
                        t = parseFloat(bar.css('top'));
                        pageY = e.pageY;

                        $doc.bind("mousemove.slimscroll", function(e){
                            currTop = t + e.pageY - pageY;
                            bar.css('top', currTop);
                            scrollContent(0, bar.position().top, false);// scroll content
                        });

                        $doc.bind("mouseup.slimscroll", function(e) {
                            isDragg = false;hideBar();
                            $doc.unbind('.slimscroll');
                        });
                        return false;
                    }).bind("selectstart.slimscroll", function(e){
                        e.stopPropagation();
                        e.preventDefault();
                        return false;
                    });
                }

                // on rail over
                rail.hover(function(){
                    showBar();
                }, function(){
                    hideBar();
                });

                // on bar over
                bar.hover(function(){
                    isOverBar = true;
                }, function(){
                    isOverBar = false;
                });

                // show on parent mouseover
                me.hover(function(){
                    isOverPanel = true;
                    showBar();
                    hideBar();
                }, function(){
                    isOverPanel = false;
                    hideBar();
                });

                // support for mobile
                me.bind('touchstart', function(e,b){
                    if (e.originalEvent.touches.length)
                    {
                        // record where touch started
                        touchDif = e.originalEvent.touches[0].pageY;
                    }
                });

                me.bind('touchmove', function(e){
                    // prevent scrolling the page if necessary
                    if(!releaseScroll)
                    {
                        e.originalEvent.preventDefault();
                    }
                    if (e.originalEvent.touches.length)
                    {
                        // see how far user swiped
                        var diff = (touchDif - e.originalEvent.touches[0].pageY) / o.touchScrollStep;
                        // scroll content
                        scrollContent(diff, true);
                        touchDif = e.originalEvent.touches[0].pageY;
                    }
                });

                // set up initial height
                getBarHeight();

                // check start position
                if (o.start === 'bottom')
                {
                    // scroll content to bottom
                    bar.css({ top: me.outerHeight() - bar.outerHeight() });
                    scrollContent(0, true);
                }
                else if (o.start !== 'top')
                {
                    // assume jQuery selector
                    scrollContent($(o.start).position().top, null, true);

                    // make sure bar stays hidden
                    if (!o.alwaysVisible) { bar.hide(); }
                }

                // attach scroll events
                attachWheel(this);

                function _onWheel(e)
                {
                    // use mouse wheel only when mouse is over
                    if (!isOverPanel) { return; }

                    var e = e || window.event;

                    var delta = 0;
                    if (e.wheelDelta) { delta = -e.wheelDelta/120; }
                    if (e.detail) { delta = e.detail / 3; }

                    var target = e.target || e.srcTarget || e.srcElement;
                    if ($(target).closest('.' + o.wrapperClass).is(me.parent())) {
                        // scroll content
                        scrollContent(delta, true);
                    }

                    // stop window scroll
                    if (e.preventDefault && !releaseScroll) { e.preventDefault(); }
                    if (!releaseScroll) { e.returnValue = false; }
                }

                function scrollContent(y, isWheel, isJump)
                {
                    releaseScroll = false;
                    var delta = y;
                    var maxTop = me.outerHeight() - bar.outerHeight();

                    if (isWheel)
                    {
                        // move bar with mouse wheel
                        delta = parseInt(bar.css('top')) + y * parseInt(o.wheelStep) / 100 * bar.outerHeight();

                        // move bar, make sure it doesn't go out
                        delta = Math.min(Math.max(delta, 0), maxTop);

                        // if scrolling down, make sure a fractional change to the
                        // scroll position isn't rounded away when the scrollbar's CSS is set
                        // this flooring of delta would happened automatically when
                        // bar.css is set below, but we floor here for clarity
                        delta = (y > 0) ? Math.ceil(delta) : Math.floor(delta);

                        // scroll the scrollbar
                        bar.css({ top: delta + 'px' });
                    }

                    // calculate actual scroll amount
                    percentScroll = parseInt(bar.css('top')) / (me.outerHeight() - bar.outerHeight());
                    delta = percentScroll * (me[0].scrollHeight - me.outerHeight());

                    if (isJump)
                    {
                        delta = y;
                        var offsetTop = delta / me[0].scrollHeight * me.outerHeight();
                        offsetTop = Math.min(Math.max(offsetTop, 0), maxTop);
                        bar.css({ top: offsetTop + 'px' });
                    }

                    // scroll content
                    me.scrollTop(delta);

                    // fire scrolling event
                    me.trigger('slimscrolling', ~~delta);

                    // ensure bar is visible
                    showBar();

                    // trigger hide when scroll is stopped
                    hideBar();
                }

                function attachWheel(target)
                {
                    if (window.addEventListener)
                    {
                        target.addEventListener('DOMMouseScroll', _onWheel, false );
                        target.addEventListener('mousewheel', _onWheel, false );
                    }
                    else
                    {
                        document.attachEvent("onmousewheel", _onWheel)
                    }
                }

                function getBarHeight()
                {
                    // calculate scrollbar height and make sure it is not too small
                    barHeight = Math.max((me.outerHeight() / me[0].scrollHeight) * me.outerHeight(), minBarHeight);
                    bar.css({ height: barHeight + 'px' });

                    // hide scrollbar if content is not long enough
                    var display = barHeight == me.outerHeight() ? 'none' : 'block';
                    bar.css({ display: display });
                }

                function showBar()
                {
                    // recalculate bar height
                    getBarHeight();
                    clearTimeout(queueHide);

                    // when bar reached top or bottom
                    if (percentScroll == ~~percentScroll)
                    {
                        //release wheel
                        releaseScroll = o.allowPageScroll;

                        // publish approporiate event
                        if (lastScroll != percentScroll)
                        {
                            var msg = (~~percentScroll == 0) ? 'top' : 'bottom';
                            me.trigger('slimscroll', msg);
                        }
                    }
                    else
                    {
                        releaseScroll = false;
                    }
                    lastScroll = percentScroll;

                    // show only when required
                    if(barHeight >= me.outerHeight()) {
                        //allow window scroll
                        releaseScroll = true;
                        return;
                    }
                    bar.stop(true,true).fadeIn('fast');
                    if (o.railVisible) { rail.stop(true,true).fadeIn('fast'); }
                }

                function hideBar()
                {
                    // only hide when options allow it
                    if (!o.alwaysVisible)
                    {
                        queueHide = setTimeout(function(){
                            if (!(o.disableFadeOut && isOverPanel) && !isOverBar && !isDragg)
                            {
                                bar.fadeOut('slow');
                                rail.fadeOut('slow');
                            }
                        }, 1000);
                    }
                }

            });

            // maintain chainability
            return this;
        }
    });

    $.fn.extend({
        slimscroll: $.fn.slimScroll
    });

})(jQuery);
(function () {
    'use strict';

    angular.module('ui.slimscroll', []).directive('slimscroll', function () {
        'use strict';

        return {
            restrict: 'A',
            link: function ($scope, $elem, $attr) {
                var off = [];
                var option = {};

                var refresh = function () {
                    if ($attr.slimscroll) {
                        option = $scope.$eval($attr.slimscroll);
                    } else if ($attr.slimscrollOption) {
                        option = $scope.$eval($attr.slimscrollOption);
                    }

                    $($elem).slimScroll({ destroy: true });

                    $($elem).slimScroll(option);
                };

                var registerWatch = function () {
                    if ($attr.slimscroll && !option.noWatch) {
                        off.push($scope.$watchCollection($attr.slimscroll, refresh));
                    }

                    if ($attr.slimscrollWatch) {
                        off.push($scope.$watchCollection($attr.slimscrollWatch, refresh));
                    }

                    if ($attr.slimscrolllistento) {
                        off.push($scope.$on($attr.slimscrolllistento, refresh));
                    }
                };

                var destructor = function () {
                    $($elem).slimScroll({ destroy: true });
                    off.forEach(function (unbind) {
                        unbind();
                    });
                    off = null;
                };

                off.push($scope.$on('$destroy', destructor));

                registerWatch();
            }
        };
    });


})();

