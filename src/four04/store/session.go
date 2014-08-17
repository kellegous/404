package store

import (
  "bytes"
  "encoding/gob"
  "four04/context"
  "four04/secure"
  "github.com/DanielMorsing/rocksdb"
  "io"
  "time"
)

type Session struct {
  Key       []byte
  UserId    uint64
  CreatedAt time.Time
  ExpiresAt time.Time
}

func (s *Session) Save(ctx *context.Context) error {
  var buf bytes.Buffer

  if err := s.Encode(&buf); err != nil {
    return err
  }

  wo := rocksdb.NewWriteOptions()
  defer wo.Close()

  return ctx.Db.Put(wo, s.Key, buf.Bytes())
}

func (s *Session) User(ctx *context.Context) (*User, error) {
  return FindUser(ctx, s.UserId)
}

func (u *Session) Encode(w io.Writer) error {
  return gob.NewEncoder(w).Encode(u)
}

func (u *Session) Decode(r io.Reader) error {
  return gob.NewDecoder(r).Decode(u)
}

func NewSession(userId uint64) (*Session, error) {
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
  ro := rocksdb.NewReadOptions()
  defer ro.Close()

  buf, err := ctx.Db.Get(ro, id)
  if err != nil {
    return nil, err
  }

  if buf == nil {
    return nil, ErrNotFound
  }

  s := &Session{}
  if err := s.Decode(bytes.NewBuffer(buf)); err != nil {
    return nil, err
  }

  return s, nil
}

func DeleteSession(ctx *context.Context, id []byte) error {
  wo := rocksdb.NewWriteOptions()
  defer wo.Close()

  return ctx.Db.Delete(wo, id)
}
