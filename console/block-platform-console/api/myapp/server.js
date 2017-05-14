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
var express = require('express');
var app = express();


//CORS middleware
/*
var allowCrossDomain = function(req, res, next) {
    res.header('Access-Control-Allow-Origin', 'example.com');
    res.header('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE');
    res.header('Access-Control-Allow-Headers', 'Content-Type');

    next();
}

app.configure(function() {
    app.use(express.bodyParser());
    app.use(express.cookieParser());
    app.use(express.session({ secret: 'cool beans' }));
    app.use(express.methodOverride());
    app.use(allowCrossDomain);
    app.use(app.router);
    app.use(express.static(__dirname + '/public'));
});
*/

// Add headers
/*
app.use(function (req, res, next) {

    // Website you wish to allow to connect
    res.setHeader('Access-Control-Allow-Origin', 'http://localhost:8888');

    // Request methods you wish to allow
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');

    // Request headers you wish to allow
    res.setHeader('Access-Control-Allow-Headers', 'X-Requested-With,content-type');

    // Set to true if you need the website to include cookies in the requests sent
    // to the API (e.g. in case you use sessions)
    res.setHeader('Access-Control-Allow-Credentials', true);

    // Pass to next layer of middleware
    next();
});
*/

app.use(function(req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Headers", "Origin,Authorization, X-Requested-With, Content-Type, Accept");
    res.header('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE');
    next();
});

// user
app.post('/iot/api/v1/account/auth', function (req, res) {
  //res.send('Hello World!');
  res.json({ token: 'eyJraWQiOiJJT1RfU0VDVVJFS0VZIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJJT1RfUExBVEZPUk1fSVNTVUVSIiwiYXVkIjoiSU9UX1BMQVRGT1JNX0FVRElFTkNFIiwiZXhwIjoxNzgwMDc0MTU1LCJqdGkiOiJjcTFDMC10RVBXUEk1XzN2Z1FsQ2RBIiwiaWF0IjoxNDY5MDM0MTU1LCJzdWIiOiJhZG1pbiIsImNsYWltLnJvbGVzIjpbIkFETUlOIl19.pzXQDv82gPrpNVas_2DHt8mihoNhqw8mnAMlDwnCC-Jkj5xodi_UBTVG8thOLaNSSLpflOqhJ8eJMstZTEJI9Nsoy1axBIun-U47NGpeZF76GUI9vh7wf_9EpwKVs0UDyK5amAVrzyiO6nQEjtMPPbGX_fWfUasB_JP5H34O2pqTl5cb6irSoJxB-_MB7lxZYJ4V9u0W9XRuFbaQtdG5YSiib7-WHHEhOIQ6X3Xg7y9josfUf41BfD9cOs2U_k3WZjiiosZVajy8DatMxF96BZuGVRh4VxozvczuiThyLAcsXW2TjYen4bgGJcH2AG7ip002NDrPxpaE2STcJwtxBQ' });
});
app.get('/iot/api/v1/account/profile', function (req, res) {
    //res.send('Hello World!');
    res.json({ 'type': 'ADMIN', 'roles':['TENANT','ADMIN'], 'nickName':'nickName-otherplayer',
    'email':'devuser@demoproject.org',
    'creationDate':'2016-06-25','enabled':true,'username':'username-le'});
});
app.post('/iot/api/v1/account/password/change', function (req, res) {
    //res.send('Hello World!');
    res.json();
});
app.put('/iot/api/v1/account/profile', function (req, res){
    console.log(req);
    res.json();
});


// Tenant
app.get('/iot/api/v1/tenants', function (req, res) {
    //res.send('get all');
    res.json({content:[{'email':'1devuser@demoproject.org', 'nickName':'1GE Tenant','username':'1username-le','enabled':true},
        {'email':'1devuser@demoproject.org', 'nickName':'1GE Tenant','username':'1username-le','enabled':true},
        {'email':'1devuser@demoproject.org', 'nickName':'2GE Tenant','username':'1username-le','enabled':true},
        {'email':'1devuser@demoproject.org', 'nickName':'3GE Tenant','username':'1username-le','enabled':true}]});
});
app.get('/iot/api/v1/tenants/id', function (req, res) {
    //res.send('get');
    res.json({'email':'1devuser@demoproject.org', 'nickName':'1GE Tenant','username':'1username-le','enabled':true});
});
app.put('/iot/api/v1/tenants/tenant1', function (req, res) {
    //res.send('update');
    res.json();
});
app.delete('/iot/api/v1/tenants/tenant1', function (req, res) {
    //res.send('delete');
    res.json();
});
app.post('/iot/api/v1/tenants', function (req, res) {
    //res.send('add');
    res.json();
});

// Device Manager
app.get('/iot/api/v1/devices', function (req, res) {
    //res.send('get all');
    res.json([{'sn':'1devuserdemoproject.org','name':'0username-le','status':true},
        {'sn':'2devuserdemoproject.org','name':'1username-le','status':true},
        {'sn':'3devuserdemoproject.org','name':'2username-le','status':true},
        {'sn':'4devuserdemoproject.org','name':'3username-le','status':true}]);
});
app.get('/iot/api/v1/devices/id', function (req, res) {
    //res.send('get');
    res.json({'sn':'1devuserdemoproject.org','name':'1username-le','status':true});
});
app.put('/iot/api/v1/devices/device1', function (req, res) {
    //res.send('update');
    res.json();
});
app.post('/iot/api/v1/devices', function (req, res) {
    //res.send('add');
    res.json();
});

// Application

app.get('/iot/api/v1/applications/application1', function (req, res) {
    res.json({
        "sn": "Windows",
        "name": "Edison",
        "token": "33-v",
        "version": "supply",
        "createdby" : "haihang",
        "date": "2013.12.12"
    });
});





//product
//post
app.post('/iot/api/v1/products', function (req, res) {
    
    res.json();
});
//get all
app.get('/iot/api/v1/products', function (req, res) {
    res.json(
        [
            {
                "productName": "computer1",
                "type": "office1",
                "creatDate": "2013.12.12"
            },
            {
                "productName": "computer2",
                "type": "office2",
                "creatDate": "2013.12.12"
            },
            {
                "productName": "computer3",
                "type": "office3",
                "creatDate": "2013.12.12"
            }

        ]
    );
});
//get one pro by id
app.get('/iot/api/v1/products/id', function (req, res) {
    res.json(
        [
            {"productName": "computer1", "displayName": "My Product1", "allowAutoRegister": false, "active": true, "description": "this is product1!!!!","profiles" : [{'key' : 'name','val' : 'string'},{'key' : 'name','val' : 'number'}]},
            {"productName": "computer2", "displayName": "My Product2", "allowAutoRegister": false, "active": true, "description": "this is product2!!!!","profiles" : [{'key' : 'name','val' : 'string'},{'key' : 'name','val' : 'number'}]},
            {"productName": "computer3", "displayName": "My Product3", "allowAutoRegister": false, "active": true, "description": "this is product3!!!!","profiles" : [{'key' : 'name','val' : 'string'},{'key' : 'name','val' : 'number'}]}
        ]
    );
});
//delete
app.delete('/iot/api/v1/products/product1', function (req, res) {
    
    res.json();
});
//modify
app.put('/iot/api/v1/products/product1', function (req, res) {
    res.json();
});
    
app.delete('/iot/api/v1/devices/device1', function (req, res) {
    //res.send('delete');
    res.json();
});


app.get('/', function(req, res, next) {
  // Handle the get for this route
  res.send('Hello World!');
});

app.post('/', function(req, res, next) {
 // Handle the post for this route
 res.send('Hello World post!');
});


app.listen(8080, function () {
  console.log('API server listening on port 8080!');
});
