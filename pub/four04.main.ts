/// <reference path="lib/socketio.d.ts" />
/// <reference path="lib/signal.ts" />
/// <reference path="lib/model.ts" />
/// <reference path="lib/convo-view.ts" />
module four04 {

var model = Model.fromLocation(),
    convo = new ConvoView(model);

model.socketDidConnect.tap((model? : Model) => {
  console.log('connect');
});

model.socketDidDisconnect.tap((model? : Model) => {
  console.log('disconnect');
});

model.messageDidArrive.tap((model? : Model, msg? : string) => {
  console.log(msg);
});

model.connect();

}