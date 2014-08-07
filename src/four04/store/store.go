package store

import (
  "bytes"
  "code.google.com/p/goauth2/oauth"
  "compress/gzip"
  "database/sql"
  "encoding/json"
  "fmt"
  "four04/config"
  "four04/secure"
  _ "github.com/go-sql-driver/mysql"
  "time"
)

const (
  userTableCreate = `
  CREATE TABLE IF NOT EXISTS user (
    Id BIGINT PRIMARY KEY,
    Email VARCHAR(255) UNIQUE NOT NULL,
    Data VARBINARY(32000) NOT NULL
  ) ENGINE=InnoDB
  `

  sessionTableCreate = `
    CREATE TABLE IF NOT EXISTS session (
      Id VARBINARY(16) PRIMARY KEY,
      UserId BIGINT NOT NULL,
      CreatedAt DATETIME NOT NULL,
      ExpiresAt DATETIME NOT NULL
    ) ENGINE=InnoDB
  `
)

type User struct {
  Id        int
  Name      string
  Email     string
  Company   string
  Location  string
  Blog      string
  CreatedAt time.Time
  UpdatedAt time.Time
  Token     *oauth.Token
}

func (u *User) Save(db *sql.DB) error {
  var buf bytes.Buffer

  w := gzip.NewWriter(&buf)

  if err := json.NewEncoder(w).Encode(u); err != nil {
    return err
  }

  if err := w.Close(); err != nil {
    return err
  }

  q := `INSERT INTO user (Id, Email, Data) VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE Email=?, Data=?`
  _, err := db.Exec(q, u.Id, u.Email, buf.Bytes(), u.Email, buf.Bytes())
  return err
}

type Session struct {
  Key       []byte
  UserId    int
  CreatedAt time.Time
  ExpiresAt time.Time
}

func (s *Session) Save(db *sql.DB) error {
  q := `INSERT INTO session (Id, UserId, CreatedAt, ExpiresAt)
        VALUES (?, ?, ?, ?)`
  _, err := db.Exec(q, s.Key, s.UserId, s.CreatedAt, s.ExpiresAt)
  return err
}

func NewSession(userId int) (*Session, error) {
  var key [16]byte
  if err := secure.FillStrongKey(key[:]); err != nil {
    return nil, err
  }

  now := time.Now().UTC()
  return &Session{
    Key:       key[:],
    UserId:    userId,
    CreatedAt: now,
    ExpiresAt: now.Add(24 * time.Hour),
  }, nil
}

func Init(cfg *config.Config) error {
  db, err := Open(cfg)
  if err != nil {
    return err
  }
  defer db.Close()

  cmds := []string{
    userTableCreate,
    sessionTableCreate,
  }

  for _, cmd := range cmds {
    if _, err := db.Exec(cmd); err != nil {
      return err
    }
  }

  return nil
}

func Open(cfg *config.Config) (*sql.DB, error) {
  return sql.Open("mysql",
    fmt.Sprintf("%s@%s/four04", cfg.Mysql.User, cfg.Mysql.Host))
}
