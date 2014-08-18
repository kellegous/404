/// <reference path="lib/jquery.d.ts" />
/// <reference path="lib/socketio.d.ts" />
module four04 {

var s = io.connect("http://localhost:8080", {path: '/api/sock'});
s.on('connect', function() {
  console.log('connected');
});

}