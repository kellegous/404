package hub

import (
  "github.com/googollee/go-socket.io"
  "github.com/kellegous/pork"
  "log"
)

func Setup(r pork.Router) error {
  srv, err := socketio.NewServer(nil)
  if err != nil {
    return err
  }

  srv.On("connection", func(s socketio.Socket) {
    log.Printf("connect: %v", s)
  })

  srv.On("error", func(s socketio.Socket) {
    log.Printf("error: %v", s)
  })

  r.RespondWith("/api/sock/", pork.ResponderFor(srv))

  return nil
}
