package hub

import (
  "four04/context"
  "github.com/kellegous/pork"
  "gopkg.in/igm/sockjs-go.v2/sockjs"
  "log"
)

const (
  channelChat = "chat"
)

func Setup(r pork.Router, ctx *context.Context) error {
  hand := sockjs.NewHandler("/api/sock", sockjs.DefaultOptions, func(s sockjs.Session) {
    log.Println(s)
  })

  r.RespondWith("/api/sock/", pork.ResponderFor(hand))
  return nil
}
