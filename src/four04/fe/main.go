package main

import (
  "errors"
  "flag"
  "fmt"
  "four04/auth"
  "four04/config"
  "four04/context"
  "four04/hub"
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

func setup(r pork.Router, ctx *context.Context) error {
  root, err := findRoot()
  if err != nil {
    return err
  }

  c := pork.Content(pork.NewConfig(pork.None),
    http.Dir(root))

  // serves up static content
  r.RespondWithFunc("/", func(w pork.ResponseWriter, r *http.Request) {
    _, err := auth.SessionFromRequest(ctx, r)
    if err != nil {
      http.Redirect(w, r, "/auth/a", http.StatusTemporaryRedirect)
      return
    }

    c.ServePork(w, r)
  })

  // some debugging handlers
  r.RespondWithFunc("/info", func(w pork.ResponseWriter, r *http.Request) {
    sess, user, err := auth.UserFromRequest(ctx, r)
    if err != nil {
      panic(err)
    }

    fmt.Fprintln(w, sess, user)
  })

  return nil
}

func main() {
  flagAddr := flag.String("addr", ":8080", "")
  flagConf := flag.String("conf", "config.json", "")
  flag.Parse()

  var cfg config.Config
  if err := cfg.LoadFromFile(*flagConf); err != nil {
    panic(err)
  }

  ctx, err := context.Open(&cfg)
  if err != nil {
    panic(err)
  }

  r := pork.NewRouter(nil, nil, nil)

  auth.Setup(r, ctx)

  if err := hub.Setup(r); err != nil {
    panic(err)
  }

  if err := setup(r, ctx); err != nil {
    panic(err)
  }

  if err := http.ListenAndServe(*flagAddr, r); err != nil {
    panic(err)
  }
}
