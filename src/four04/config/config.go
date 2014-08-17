package config

import (
  "crypto/sha256"
  "encoding/json"
  "fmt"
  "os"
  "path/filepath"
)

type Config struct {
  OAuth struct {
    ClientId     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
  }

  HmacKey []byte `json:"hmac_key"`

  DbPath string `json:"dbpath"`

  RootPath string `json:-`
}

func (c *Config) LoadFromFile(filename string) error {
  r, err := os.Open(filename)
  if err != nil {
    return err
  }
  defer r.Close()

  if err := json.NewDecoder(r).Decode(c); err != nil {
    return err
  }

  if !filepath.IsAbs(c.DbPath) {
    dir, err := filepath.Abs(filepath.Dir(filename))
    if err != nil {
      return err
    }

    c.DbPath = filepath.Join(dir, c.DbPath)
  }

  if len(c.HmacKey) < sha256.BlockSize {
    return fmt.Errorf("HMAC Key is too short: %d < %d", len(c.HmacKey), sha256.BlockSize)
  }

  return nil
}
