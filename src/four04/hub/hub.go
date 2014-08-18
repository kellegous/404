package hub

import (
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
    _, _, err := auth.UserFromRequest(ctx, s.Request())
    if err != nil {
      // TODO(knorton): Access denied. There is no API to close a socket.
      return
    }

    s.Join(channelChat)

    s.On("msg", func(msg string) {
      s.BroadcastTo(channelChat, msg)
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
