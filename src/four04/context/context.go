package context

import (
  "four04/config"
  "github.com/DanielMorsing/rocksdb"
  "os"
  "path/filepath"
)

var DefaultConfig *config.Config

type Context struct {
  Db  *rocksdb.DB
  Cfg *config.Config
}

func (c *Context) Close() {
  c.Db.Close()
}

func Open(cfg *config.Config) (*Context, error) {
  p := filepath.Join(cfg.DbPath, "db")
  if _, err := os.Stat(p); err != nil {
    if err := os.MkdirAll(p, os.ModePerm); err != nil {
      return nil, err
    }
  }

  o := rocksdb.NewOptions()
  o.SetCreateIfMissing(true)
  defer o.Close()

  db, err := rocksdb.Open(p, o)
  if err != nil {
    return nil, err
  }

  return &Context{
    Db:  db,
    Cfg: cfg,
  }, nil
}
