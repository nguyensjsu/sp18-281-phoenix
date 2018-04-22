/**

Starbucks, Inc.
Version 1.0

**/


var fs = require('fs');
var express = require('express');
var Client = require('node-rest-client').Client;
var session = require('client-sessions');

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
      body = fs.readFileSync('.' + url);
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

var handle_get = function (req, res) {
    page( req, res, "San Jose" ) ;
}

app.set('port', (process.env.PORT || 3000));

// app.post("*", handle_post );
app.get( "*", handle_get ) ;

app.listen(app.get('port'), function() {
  console.log('running on port', app.get('port'));
});
