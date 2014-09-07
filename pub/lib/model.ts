/// <reference path="sockjs.d.ts" />
/// <reference path="jquery.d.ts" />
module four04 {

export class Model {
  socketDidConnect = new Signal;

  socketDidDisconnect = new Signal;

  messageDidArrive = new Signal;

  private socket : SockJS;

  constructor(private sockPath : string, private authPath : string) {
  }

  private auth(socket : SockJS) {
    socket.onmessage = (event : SJSMessageEvent) => {
      var msg = JSON.parse(event.data);
      if (msg.Type != 'connect') {
        socket.close();
      }

      this.socket = socket;
      this.socketDidConnect.raise(this);
    };

    socket.onclose = (event : SJSCloseEvent) => {
      this.socket = null;
    };

    $.ajax({
      url: this.authPath,
      dataType: 'text',
      success: (data : string) => {
        socket.send(JSON.stringify({
          Type: 'connect',
          Token: data
        }));
      },
      error: (xhr, status, error) => {
        console.error(error.toString());
      }
    });
  }

  /**
   *
   */
  connect() {
    if (this.socket) {
      return;
    }

    var socket = new SockJS(this.sockPath, null);
    socket.onopen = (e) => {
      $.ajax({
        url: this.authPath,
        dataType: 'text',
        success: (data : string) => {
          socket.send(JSON.stringify({
            Type: 'connect',
            Token: data
          }));
        },
        error: (xhr, status, error) => {
          console.error(error.toString());
        }
      });


      socket.onclose = (event : SJSCloseEvent) => {
        this.socket = null;
      };

      socket.onmessage = (event : SJSMessageEvent) => {
        var msg = JSON.parse(event.data);
        if (msg.Type != 'connect') {
          socket.close();
        }

        socket.onmessage = (event : SJSMessageEvent) => {
          console.log(event);
        };

        socket.onclose = (event : SJSCloseEvent) => {
          this.socket = null;
          this.socketDidDisconnect.raise(this);
        };

        this.socket = socket;
        this.socketDidConnect.raise(this);
      };
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
  send(msg : any) {
    if (this.socket) {
      this.socket.send(JSON.stringify({
        To: 0,
        Type: "m",
        Msg: msg
      }));
    }
  }
}

}