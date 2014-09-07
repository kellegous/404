package hub

import (
  "encoding/json"
  _ "expvar"
  "fmt"
  "four04/auth"
  "four04/context"
  "four04/store"
  "github.com/kellegous/pork"
  "gopkg.in/igm/sockjs-go.v2/sockjs"
  "log"
  "net/http"
  _ "net/http/pprof"
)

const (
  typeConnect = "connect"
  typeMessage = "m"
)

type userSession struct {
  user *store.User
  sess map[string]sockjs.Session
}

type hub struct {
  users map[uint64]*userSession
  ch    chan func()
}

func (h *hub) enter(s sockjs.Session, user *store.User) {
  h.ch <- func() {
    us := h.users[user.Id]
    if us == nil {
      us = &userSession{
        user: user,
        sess: map[string]sockjs.Session{},
      }
      h.users[user.Id] = us
    }
    us.sess[s.ID()] = s
    log.Printf("enter: %s", user.Name)
  }
}

func (h *hub) leave(s sockjs.Session, user *store.User) {
  h.ch <- func() {
    us := h.users[user.Id]
    if us == nil {
      return
    }

    delete(us.sess, s.ID())
    log.Printf("exit: %s", user.Name)
  }
}

func (h *hub) sendTo(user uint64, msg string) {
  h.ch <- func() {
    us := h.users[user]
    if us == nil {
      return
    }

    for _, s := range us.sess {
      s.Send(msg)
    }
  }
}

func (h *hub) broadcast(msg string) {
  h.ch <- func() {
    log.Printf("broadcast: %s", msg)
    for _, us := range h.users {
      for _, s := range us.sess {
        if err := s.Send(msg); err != nil {
          log.Printf("send: %s", err)
        }
      }
    }
  }
}

func (h *hub) start() {
  h.users = map[uint64]*userSession{}
  h.ch = make(chan func())
  go func() {
    for fn := range h.ch {
      fn()
    }
  }()
}

func (h *hub) dispatch(s sockjs.Session, user *store.User, msg string) error {
  log.Printf("dispatch: %s (%s)", msg, user.Name)
  var m struct {
    Type string
    To   uint64
  }

  if err := json.Unmarshal([]byte(msg), &m); err != nil {
    return err
  }

  switch m.Type {
  case typeMessage:
    if m.To == 0 {
      h.broadcast(msg)
    } else {
      h.sendTo(m.To, msg)
    }
  }

  return nil
}

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

  if req.Type != typeConnect {
    return nil, fmt.Errorf("hub: expected \"connect\" got \"%s\"", req.Type)
  }

  sess, err := auth.SessionFromToken(ctx, req.Token)
  if err != nil {
    return nil, err
  }

  if err := send(s, map[string]string{
    "Type": typeConnect,
  }); err != nil {
    return nil, err
  }

  return sess, nil
}

func accessDenied(s sockjs.Session) error {
  return s.Close(http.StatusForbidden, http.StatusText(http.StatusForbidden))
}

func Setup(r pork.Router, ctx *context.Context) error {
  var h hub

  h.start()

  hand := sockjs.NewHandler("/api/sock", sockjs.DefaultOptions, func(s sockjs.Session) {
    sess, err := authenticate(s, ctx)
    if err != nil {
      accessDenied(s)
      return
    }

    user, err := sess.User(ctx)
    if err != nil {
      accessDenied(s)
      return
    }

    h.enter(s, user)
    defer h.leave(s, user)

    for {
      msg, err := s.Recv()
      if err != nil {
        log.Printf("recv: %s", err)
        return
      }

      if err := h.dispatch(s, user, msg); err != nil {
        log.Printf("disp: %s", err)
        return
      }
    }
  })

  r.RespondWith("/api/sock/", pork.ResponderFor(hand))
  return nil
}
