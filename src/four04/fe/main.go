package main

import (
  "errors"
  "flag"
  "fmt"
  "four04/auth"
  "four04/config"
  "four04/context"
  "four04/store"
  "github.com/kellegous/pork"
  "net/http"
  "os"
  "path/filepath"
  "runtime"
)

func findRoot() (string, error) {
  dir, err := filepath.Abs(
    filepath.Join(filepath.Dir(os.Args[0]), "../pub"))
  if err == nil {
    if _, err := os.Stat(dir); err == nil {
      return dir, nil
    }
  }

  _, file, _, _ := runtime.Caller(0)
  dir, err = filepath.Abs(
    filepath.Join(filepath.Dir(file), "../../../pub"))
  if err == nil {
    if _, err := os.Stat(dir); err == nil {
      return dir, nil
    }
  }

  return "", errors.New("cannot locate pub directory")
}

func setup(r pork.Router, root string, cfg *config.Config) {
  c := pork.Content(pork.NewConfig(pork.None),
    http.Dir(root))

  // serves up static content
  r.RespondWithFunc("/", func(w pork.ResponseWriter, r *http.Request) {
    ctx := context.MustOpen(cfg)
    defer ctx.Close()

    _, err := auth.SessionFromRequest(ctx, r)
    if err == store.ErrNotFound {
      http.Redirect(w, r, "/auth/a", http.StatusTemporaryRedirect)
      return
    } else if err != nil {
      panic(err)
    }

    c.ServePork(w, r)
  })

  // some debugging handlers
  r.RespondWithFunc("/info", func(w pork.ResponseWriter, r *http.Request) {
    ctx := context.MustOpen(cfg)
    defer ctx.Close()

    sess, user, err := auth.UserFromRequest(ctx, r)
    if err != nil {
      panic(err)
    }

    fmt.Fprintln(w, sess, user)
  })
}

func main() {
  flagAddr := flag.String("addr", ":8080", "")
  flagConf := flag.String("conf", "config.json", "")
  flag.Parse()

  root, err := findRoot()
  if err != nil {
    panic(err)
  }

  var cfg config.Config
  if err := cfg.LoadFromFile(*flagConf); err != nil {
    panic(err)
  }

  if err := store.Init(&cfg); err != nil {
    panic(err)
  }

  r := pork.NewRouter(nil, nil, nil)

  auth.Setup(r, &cfg)
  setup(r, root, &cfg)

  if err := http.ListenAndServe(*flagAddr, r); err != nil {
    panic(err)
  }
}
