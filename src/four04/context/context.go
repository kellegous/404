package context

import (
  "database/sql"
  "fmt"
  "four04/config"
)

var DefaultConfig *config.Config

type Context struct {
  Db  *sql.DB
  Cfg *config.Config
}

func (c *Context) Close() error {
  return c.Db.Close()
}

func MustOpen(cfg *config.Config) *Context {
  ctx, err := Open(cfg)
  if err != nil {
    panic(err)
  }
  return ctx
}

func Open(cfg *config.Config) (*Context, error) {
  db, err := sql.Open("mysql",
    fmt.Sprintf("%s@%s/four04?parseTime=true", cfg.Mysql.User, cfg.Mysql.Host))
  if err != nil {
    return nil, err
  }

  if cfg == nil {
    cfg = DefaultConfig
  }

  return &Context{
    Db:  db,
    Cfg: cfg,
  }, nil
}
