package main

import (
  "encoding/json"
  "flag"
  "four04/auth"
  "four04/store"
  "github.com/kellegous/pork"
  "net/http"
  "os"
)

type Config struct {
  OAuth struct {
    ClientId     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
  }

  Mysql struct {
    Host string `json:"host"`
    User string `json:"user"`
  }
}

func (c *Config) ClientId() string {
  return c.OAuth.ClientId
}

func (c *Config) ClientSecret() string {
  return c.OAuth.ClientSecret
}

func (c *Config) MysqlHost() string {
  return c.Mysql.Host
}

func (c *Config) MysqlUser() string {
  return c.Mysql.User
}

func (c *Config) Load(filename string) error {
  r, err := os.Open(filename)
  if err != nil {
    return err
  }
  defer r.Close()
  return json.NewDecoder(r).Decode(c)
}

func main() {
  flagAddr := flag.String("addr", ":8080", "")
  flagConf := flag.String("conf", "config.json", "")
  flag.Parse()

  var cfg Config
  if err := cfg.Load(*flagConf); err != nil {
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
