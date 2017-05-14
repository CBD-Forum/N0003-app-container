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

    angular.module('airs').controller('DashboardController', DashboardController);

    /** @ngInject */
    function DashboardController($stateParams,ApiServer,toastr,$state,$timeout,$interval) {
        /* jshint validthis: true */
        var vm = this;


        // oh shit ~

        vm.title = '此处是监控信息，';

        /**
         * 此处是控制台信息，集成了公司物联网云平台集装箱芯片信息用于
         * 定位集装箱，航线信息暂不能公开
         *
        var info = ApiServer.info();
        var roleType = info.role;
        vm.title = '物流动态';
        vm.containerlists = [];
        var width = document.body.clientWidth;
        var height = document.body.clientHeight;
        vm.mapSize = {"width":width + 'px',"height":height + 'px'};
        vm.flag = false;
        vm.flag1 = false;

        //创建航线，航线是固定的。。。出于保密原因，航线信息。。


        // 百度地图API功能
        var map = new BMap.Map("goodtrack",{minZoom:3,maxZoom:14});    // 创建Map实例
        map.centerAndZoom("上海",5);      // 初始化地图,用城市名设置地图中心点
        map.enableScrollWheelZoom(false);     //开启鼠标滚轮缩放
        // 创建地址解析器实例
        var myGeo = new BMap.Geocoder();
        //点击地图，获取经纬度坐标
        map.addEventListener("click",function (e) {
            var ll = e.point.lng+","+e.point.lat;
            console.log(ll);
        });



        var top_right_navigation = new BMap.NavigationControl({anchor: BMAP_ANCHOR_TOP_LEFT, type: BMAP_NAVIGATION_CONTROL_SMALL}); //右上角，仅包含平移和缩放按钮
        map.addControl(top_right_navigation);

        // 从集团IoT云平台获取集装箱定位信息


        // 航线是固定的，这里先这样


        function gotoOrder() {
            if (roleType === 'cargoagent'){
                $state.go('app.goodorder');
            }else if (roleType === 'carrier'){
                $state.go('app.carorder');
            }else if (roleType === 'shipper'){
                $state.go('app.shiporder');
            }else{
                $state.go('app.userorder');
            }
        }


        //通过地址获取经纬度
        function getMapPointFromAddress(address,successHandler) {
            if (address.indexOf('-') !== -1){
                var addresses = address.split('-');
                var city = addresses[0];
                var detail = addresses[1];
                getPointFromAddress(city,detail,function (point) {
                    successHandler(point);
                })
            }
        }
        function getPointFromAddress(city,detail,successHandler) {//北京市 北京市海淀区上地10街10号
            myGeo.getPoint(detail, function(point){
                successHandler(point);
            }, city);
        }
        // 编写自定义函数,创建标注
        function addMarker(point){
            var marker = new BMap.Marker(point);
            map.addOverlay(marker);
        }
        function addMainlandRouteCurve(info1,info2,handler,color) {
            if (color){
                addCurve(info1,info2,handler,color,4);
            }else{
                addCurve(info1,info2,handler,"#27c24c",4);
            }
        }
        function addOceanRouteCurve(info1,info2,handler,color) {
            if (color){
                addCurve(info1,info2,handler,color,4);
            }else{
                addCurve(info1,info2,handler,"#23b7e5",4);
            }
        }
        //向地图中添加线函数
        function addLine(points,color) {
            var plPoints;
            if (color){
                plPoints = [{style:"solid",weight:4,color:color,opacity:0.6,points:points}];
            }else{
                plPoints = [{style:"solid",weight:4,color:"#f00",opacity:0.6,points:points}];
            }
            addPolyline(plPoints);
        }
        function addPolyline(plPoints){
            for(var i=0;i<plPoints.length;i++){

                var json = plPoints[i];
                var points = [];
                for(var j=0;j<json.points.length;j++){
                    var p1 = json.points[j].latitude;
                    var p2 = json.points[j].longitude;
                    points.push(new BMap.Point(p1,p2));
                }
                var line = new BMap.Polyline(points,{strokeStyle:json.style,strokeWeight:json.weight,strokeColor:json.color,strokeOpacity:json.opacity});
                map.addOverlay(line);
            }
        }
        function addCurve(info1,info2,handler,color,weight) {
            var p1 = new BMap.Point(info1.latitude,info1.longitude);
            var p2 = new BMap.Point(info2.latitude,info2.longitude);
            var points = [p1, p2];
            var curve = new BMapLib.CurveLine(points, {strokeColor:color, strokeWeight:weight, strokeOpacity:0.7}); //创建弧线对象
            map.addOverlay(curve); //添加到地图中
            curve.addEventListener("click",function () {
                handler();
            });
        }
        function addPoint(info,jump,type) {
            var pt = new BMap.Point(info.latitude, info.longitude);
            var marker;
            if (type){
                var icon;
                if (type === 'b'){
                    icon = 'images/icon-boat.png';
                }else if (type === 'j'){
                    icon = 'images/icon-j.png';
                }else if (type === 'p'){
                    icon = 'images/icon-point.png';
                }else{
                    icon = 'images/icon-point.png';
                }
                var myIcon = new BMap.Icon(icon, new BMap.Size(36,36));
                marker = new BMap.Marker(pt,{icon:myIcon});  // 创建标注
            }else{
                marker = new BMap.Marker(pt);
            }

            map.addOverlay(marker);
            if (jump){
                marker.setAnimation(BMAP_ANIMATION_BOUNCE); //跳动的动画
            }

            var opts = {
                width : 200,     // 信息窗口宽度
                height: 70,     // 信息窗口高度
                title : info.title
            }
            var infoWindow = new BMap.InfoWindow(info.detail, opts);  // 创建信息窗口对象
            marker.addEventListener("click", function(){
                map.openInfoWindow(infoWindow,pt); //开启信息窗口
            });

            if (info.message){
                addTip(info);
            }
        }
        function addTip(info) {
            var pt = new BMap.Point(info.latitude, info.longitude);
            var opts = {
                position : pt,    // 指定文本标注所在的地理位置
                offset   : new BMap.Size(20, -20)    //设置文本偏移量
            }
            var label = new BMap.Label(' ' + info.message, opts);  // 创建文本标注对象
            label.setStyle({
                color : "green",
                fontSize : "12px",
                height : "20px",
                border:"none",
                padding:'0',
                lineHeight : "20px",
                fontFamily:"微软雅黑"
            });
            map.addOverlay(label);
        }

         **/



    }

})();
