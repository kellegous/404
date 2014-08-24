/// <reference path="lib/signal.ts" />
/// <reference path="lib/model.ts" />
/// <reference path="lib/convo-view.ts" />
module four04 {

var model = new Model('/api/sock'),
    convo = new ConvoView(model);

model.socketDidConnect.tap((model? : Model) => {
  console.log('connect');
});

model.socketDidDisconnect.tap((model? : Model) => {
  console.log('disconnect');
});

model.connect();

}