/// <reference path="lib/jquery.d.ts" />
/// <reference path="lib/socketio.d.ts" />
/// <reference path="lib/signal.ts" />
/// <reference path="lib/model.ts" />
module four04 {

var model = Model.fromLocation();

model.socketDidConnect.tap((model? : Model) => {
  console.log('connect');
});

model.socketDidDisconnect.tap((model? : Model) => {
  console.log('disconnect');
});

model.connect();

}