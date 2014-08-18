/// <reference path="jquery.d.ts" />
/// <reference path="model.ts" />
module four04 {

export class ConvoView {
  private root : JQuery;

  constructor(public model : Model) {
    this.root = $('#messages');
  }
}

}