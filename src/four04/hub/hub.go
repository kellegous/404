package hub

import (
  "fmt"
  "four04/auth"
  "four04/context"
  "github.com/googollee/go-socket.io"
  "github.com/kellegous/pork"
  "log"
)

const (
  channelChat = "chat"
)

func Setup(r pork.Router, ctx *context.Context) error {
  srv, err := socketio.NewServer(nil)
  if err != nil {
    return err
  }

  srv.On("connection", func(s socketio.Socket) {
    log.Printf("%v", s)
    _, _, err := auth.UserFromRequest(ctx, s.Request())
    if err != nil {
      s.Emit("disconnect")
      // TODO(knorton): Access denied. There is no API to close a socket.
      return
    }

    s.Join(channelChat)

    s.On("msg", func(msg string) {
      s.BroadcastTo(channelChat, "msg", msg)
    })

    s.On(fmt.Sprintf("%s msg", channelChat), func(msg string) {
      s.Emit("msg", msg)
    })

    s.On("disconnection", func() {
      log.Println("disconnect")
    })
  })

  srv.On("error", func(s socketio.Socket) {
    log.Printf("error: %v", s)
  })

  r.RespondWith("/api/sock/", pork.ResponderFor(srv))

  return nil
}
