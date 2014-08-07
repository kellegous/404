package main

import (
  "flag"
  "four04/auth"
  "four04/config"
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

  if err := http.ListenAndServe(*flagAddr, r); err != nil {
    panic(err)
  }
}
