package store

import (
  "bytes"
  "code.google.com/p/goauth2/oauth"
  "compress/gzip"
  "database/sql"
  "encoding/json"
  "errors"
  "four04/config"
  "four04/context"
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

var (
  ErrNotFound = errors.New("record not found")
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

func (u *User) Save(ctx *context.Context) error {
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
  _, err := ctx.Db.Exec(q, u.Id, u.Email, buf.Bytes(), u.Email, buf.Bytes())
  return err
}

func (u *User) fromRow(r *sql.Rows) error {
  var data []byte

  if err := r.Scan(&data); err != nil {
    return err
  }

  gr, err := gzip.NewReader(bytes.NewReader(data))
  if err != nil {
    return err
  }

  return json.NewDecoder(gr).Decode(u)
}

func FindUser(ctx *context.Context, id int) (*User, error) {
  q := "SELECT Data FROM user WHERE Id=?"
  r, err := ctx.Db.Query(q, id)
  if err != nil {
    return nil, err
  }
  defer r.Close()

  if !r.Next() {
    return nil, ErrNotFound
  }

  user := &User{}
  if err := user.fromRow(r); err != nil {
    return nil, err
  }

  return user, nil
}

type Session struct {
  Key       []byte
  UserId    int
  CreatedAt time.Time
  ExpiresAt time.Time
}

func (s *Session) Save(ctx *context.Context) error {
  q := `INSERT INTO session (Id, UserId, CreatedAt, ExpiresAt)
        VALUES (?, ?, ?, ?)`
  _, err := ctx.Db.Exec(q, s.Key, s.UserId, s.CreatedAt, s.ExpiresAt)
  return err
}

func (s *Session) User(ctx *context.Context) (*User, error) {
  return FindUser(ctx, s.UserId)
}

func (s *Session) fromRow(r *sql.Rows) error {
  return r.Scan(&s.Key, &s.UserId, &s.CreatedAt, &s.ExpiresAt)
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

func FindSession(ctx *context.Context, id []byte) (*Session, error) {
  q := "SELECT Id, UserId, CreatedAt, ExpiresAt FROM session WHERE ID=?"
  r, err := ctx.Db.Query(q, id)
  if err != nil {
    return nil, err
  }
  defer r.Close()

  if !r.Next() {
    return nil, ErrNotFound
  }

  sess := &Session{}
  if err := sess.fromRow(r); err != nil {
    return nil, err
  }

  return sess, nil
}

func DeleteSession(ctx *context.Context, id []byte) error {
  q := "DELETE FROM session WHERE id=?"
  _, err := ctx.Db.Exec(q, id)
  return err
}

func Init(cfg *config.Config) error {
  ctx, err := context.Open(cfg)
  if err != nil {
    return err
  }
  defer ctx.Close()

  cmds := []string{
    userTableCreate,
    sessionTableCreate,
  }

  for _, cmd := range cmds {
    if _, err := ctx.Db.Exec(cmd); err != nil {
      return err
    }
  }

  return nil
}
