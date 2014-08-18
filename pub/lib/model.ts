
module four04 {

export class Model {
  socketDidConnect = new Signal;

  socketDidDisconnect = new Signal;

  msgDidArrive = new Signal;

  private socket : io.Socket;

  constructor(private host : string) {
  }

  connect() {
    if (this.socket) {
      return;
    }

    var socket = io.connect(this.host, { path: '/api/sock'});
    socket.on('connect', () => {
      this.socket = socket;
      this.socketDidConnect.raise(this);
    });

    socket.on('disconnect', () => {
      this.socket = null;
      this.socketDidDisconnect.raise(this);
    });
  }

  static fromLocation() : Model {
    return new Model(location.protocol + '//' + location.host);
  }
}

}