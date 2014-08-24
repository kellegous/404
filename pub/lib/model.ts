/// <reference path="sockjs.d.ts" />
module four04 {

export class Model {
  socketDidConnect = new Signal;

  socketDidDisconnect = new Signal;

  messageDidArrive = new Signal;

  private socket : SockJS;

  constructor(public path : string) {
  }

  /**
   *
   */
  connect() {
    if (this.socket) {
      return;
    }

    var socket = new SockJS(this.path, null);
    socket.onopen = (e) => {
      this.socket = socket;
      this.socketDidConnect.raise(this);
    };

    socket.onclose = (e) => {
      this.socket = null;
      this.socketDidDisconnect.raise(this);
    };

    socket.onmessage = (e) => {
      console.log(e);
    };
  }

  /**
   *
   */
  connected() : boolean {
    return !!this.socket;
  }

  /**
   *
   */
  send(ch : string, msg : any) {
    if (this.socket) {
      // send
    }
  }
}

}