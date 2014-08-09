package main

import (
  "flag"
  "fmt"
  "four04/auth"
  "four04/config"
  "four04/context"
  "four04/store"
  "github.com/kellegous/pork"
  "net/http"
)

func main() {
  flagAddr := flag.String("addr", ":8080", "")
  flagConf := flag.String("conf", "config.json", "")
  flag.Parse()

  var cfg config.Config
  if err := cfg.LoadFromFile(*flagConf); err != nil {
    panic(err)
  }

  if err := store.Init(&cfg); err != nil {
    panic(err)
  }

  r := pork.NewRouter(nil, nil, nil)

  auth.Setup(r, &cfg)

  r.RespondWithFunc("/info", func(w pork.ResponseWriter, r *http.Request) {
    ctx, err := context.FromRequest(r, &cfg)
    if err != nil {
      panic(err)
    }
    defer ctx.Close()

    sess, err := auth.SessionFromRequest(ctx, r)
    if err != nil {
      panic(err)
    }

    fmt.Fprintf(w, "%v", sess)

    if sess == nil {
      return
    }

    user, err := sess.User(ctx)
    if err != nil {
      panic(err)
    }

    fmt.Fprintf(w, "%v", user)
  })

  if err := http.ListenAndServe(*flagAddr, r); err != nil {
    panic(err)
  }
}
