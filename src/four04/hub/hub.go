package hub

import (
  "encoding/json"
  "fmt"
  "four04/auth"
  "four04/context"
  "four04/store"
  "github.com/kellegous/pork"
  "gopkg.in/igm/sockjs-go.v2/sockjs"
  "log"
)

func send(s sockjs.Session, data interface{}) error {
  b, err := json.Marshal(data)
  if err != nil {
    return err
  }

  return s.Send(string(b))
}

func authenticate(s sockjs.Session, ctx *context.Context) (*store.Session, error) {
  val, err := s.Recv()
  if err != nil {
    return nil, err
  }

  var req struct {
    Type  string
    Token string
  }

  if err := json.Unmarshal([]byte(val), &req); err != nil {
    return nil, err
  }

  if req.Type != "connect" {
    return nil, fmt.Errorf("hub: expected \"connect\" got \"%s\"", req.Type)
  }

  sess, err := auth.SessionFromToken(ctx, req.Token)
  if err != nil {
    return nil, err
  }

  if err := send(s, map[string]string{
    "Type": "connect",
  }); err != nil {
    return nil, err
  }

  return sess, nil
}

func Setup(r pork.Router, ctx *context.Context) error {
  hand := sockjs.NewHandler("/api/sock", sockjs.DefaultOptions, func(s sockjs.Session) {
    sess, err := authenticate(s, ctx)
    if err != nil {
      s.Close(403, "Access Denied")
      return
    }
    log.Println(sess)
  })

  r.RespondWith("/api/sock/", pork.ResponderFor(hand))
  return nil
}
