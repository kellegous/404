/// <reference path="jquery.d.ts" />
/// <reference path="model.ts" />
module four04 {

export class ConvoView {
  private msgs : JQuery;
  private text : JQuery;

  constructor(public model : Model) {
    var msgs = $('#messages'),
        text = $('#message');

    text.on('keydown', (e : KeyboardEvent) => {
      if (e.keyCode != 13) {
        return;
      }

      this.model.send('msg', JSON.stringify({
        text : text.val()
      }));

      text.val('');
    });

    model.messageDidArrive.tap((model? : Model, msg? : string) => {
      console.log('message', msg);
    });

    this.msgs = msgs;
    this.text = text;
  }
}

}