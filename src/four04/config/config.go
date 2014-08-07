package config

import (
  "encoding/json"
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

func (c *Config) LoadFromFile(filename string) error {
  r, err := os.Open(filename)
  if err != nil {
    return err
  }
  defer r.Close()

  return json.NewDecoder(r).Decode(c)
}
