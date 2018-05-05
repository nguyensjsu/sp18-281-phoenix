/**

Starbucks, Inc.
Version 1.0
endpoints will be different for different VPCs
**/

var machine = "http://ec2-13-57-59-79.us-west-1.compute.amazonaws.com:3000/";
var cartendpoint = "http://ec2-13-57-59-79.us-west-1.compute.amazonaws.com:3000/";
var productendpoint = "http://ec2-54-241-198-25.us-west-1.compute.amazonaws.com:3000/products"

var orderEndpoint = "http://ec2-user@ec2-52-52-199-60.us-west-1.compute.amazonaws.com:3000/";

var fs = require('fs');
var express = require('express');
var Client = require('node-rest-client').Client;
var session = require('client-sessions');
var tailor = require('./tailor.js');

var app = express();
app.use(express.bodyParser());
app.use("/images", express.static(__dirname + '/images'));
app.use("/css", express.static(__dirname + '/css'));
app.use("/js", express.static(__dirname + '/js'));
app.use("/bootstrap-3.3.7", express.static(__dirname + '/bootstrap-3.3.7'));

var page = function( req, res, location ) {
    var url = req.originalUrl
    if (req.originalUrl == '/') {
      url = '/index.html'
    }
    console.log(url);
    res.setHeader('Content-Type', 'text/html');
    var body = 'not found'
    try {
      body = tailor.getPage(url);
      res.writeHead(200);
    } catch(err) {
      res.writeHead(404);
    }
    res.end( body );
}

app.use(session({
  cookieName: 'session',
  secret: 'myStore',
  duration: 30 * 60 * 1000,
  activeDuration: 5 * 60 * 1000,
}));

app.post('/setStore', function(req, res) {
  req.session.store = req.body.store;
  res.cookie('storeLocation',req.body.store, { maxAge: 900000, httpOnly: true });
  console.log('cookie created successfully');
  res.redirect('/locations.html');
});

app.post('/cart', function(req, res) {
  var client = new Client();

  var args = {
    data: {"Items": [{
    "Drink": req.body.Drink,
    "Size": req.body.Size,
    "Options": req.body.Options,
    "Price": parseFloat(req.body.Price),
    "Quantity":parseInt(req.body.Quantity, 10)}]
  },
    headers: { "Content-Type": "application/json" }};
  client.post(cartendpoint+"/"+req.body.StoreLocation+"/cart"+, args,
    function(data, response_raw) {
      console.log(data);
      req.session.cartId = data.Id;
      console.log('cookie created successfully');
      res.send(data);
      });
});

app.put('/cart/:cartId', function(req, res) {
  var client = new Client();
  client.get(machine+"/"+req.params.cartId, function (data, response) {
    // parsed response body as js object
    var cart = JSON.parse(data);
    console.log(cart.Items);
    cart.Items.push({
    "Drink": req.body.Drink,
    "Size": req.body.Size,
    "Options": req.body.Options,
    "Price": parseFloat(req.body.Price),
    "Quantity":parseInt(req.body.Quantity, 10)});


    var args = {
    data: {"Items": cart.Items},
    headers: { "Content-Type": "application/json" }};
  client.put( cartendpoint+"/"+req.body.StoreLocation+"/cart/"+req.params.cartId, args,
    function(data, response_raw) {
      console.log(data);
      res.send(data);
      });

  });

});

app.post('/getCart/:cartId', function(req, res) {
  var client = new Client();
  client.get(cartendpoint+"/"+req.body.StoreLocation+"/cart/"+req.params.cartId, function (data, response) {
    // parsed response body as js object
    var cart = JSON.parse(data);
    console.log(cart.Items);
    res.send(cart);
  });
});

var handle_get = function (req, res) {
    page( req, res, "San Jose" ) ;
}

app.delete('/order/:orderId', function(req, res) {
  var client = new Client();
  client.delete(orderEndpoint + "order/" + req.params.orderId, function (data, response) {
    console.log(data);
    res.send(data)
  });
});

app.delete('/cart/:cartId', function(req, res) {
  var client = new Client();
  client.delete(cartendpoint+"/"+req.body.StoreLocation+"/cart/"+req.params.cartId, function (data, response) {
    // parsed response body as js object
    var cart = JSON.parse(data);
    console.log(cart);
    res.send(cart);
  });
});




app.post('/products', function(req, res) {
  var client = new Client();
  client.get(productendpoint+"/"+req.body.StoreLocation+"/products/", function (data, response) {
    // parsed response body as js object
    var cart = JSON.parse(data);
    console.log(data);
    res.send(data);
  });
});

app.set('port', (process.env.PORT || 3000));

// app.post("*", handle_post );
app.get( "*", handle_get ) ;

app.listen(app.get('port'), function() {
  console.log('running on port', app.get('port'));
});
