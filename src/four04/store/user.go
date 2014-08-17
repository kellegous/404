package store

import (
  "bytes"
  "code.google.com/p/goauth2/oauth"
  "encoding/binary"
  "encoding/gob"
  "four04/context"
  "github.com/DanielMorsing/rocksdb"
  "io"
  "time"
)

type User struct {
  Id        uint64
  Name      string
  Email     string
  Company   string
  Location  string
  Blog      string
  CreatedAt time.Time
  UpdatedAt time.Time
  Token     *oauth.Token
}

func userKeyBytes(id uint64) []byte {
  var buf bytes.Buffer
  buf.Write([]byte{UserKind})
  binary.Write(&buf, binary.LittleEndian, id)
  return buf.Bytes()
}

func (u *User) Save(ctx *context.Context) error {
  var buf bytes.Buffer

  if err := u.Encode(&buf); err != nil {
    return err
  }

  wo := rocksdb.NewWriteOptions()
  defer wo.Close()

  return ctx.Db.Put(wo, userKeyBytes(u.Id), buf.Bytes())
}

func (u *User) Encode(r io.Writer) error {
  return gob.NewEncoder(r).Encode(u)
}

func (u *User) Decode(r io.Reader) error {
  return gob.NewDecoder(r).Decode(u)
}

func FindUser(ctx *context.Context, id uint64) (*User, error) {
  ro := rocksdb.NewReadOptions()
  defer ro.Close()

  buf, err := ctx.Db.Get(ro, userKeyBytes(id))
  if err != nil {
    return nil, err
  }

  if buf == nil {
    return nil, ErrNotFound
  }

  u := &User{}
  if err := u.Decode(bytes.NewBuffer(buf)); err != nil {
    return nil, err
  }

  return u, nil
}
