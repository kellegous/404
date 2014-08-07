package secure

import (
  "bytes"
  "crypto/hmac"
  "crypto/rand"
  "crypto/sha256"
  "encoding/binary"
  "errors"
  "io"
  "time"
)

var (
  ErrInvalidSignature = errors.New("secure: invalid signature")
)

func FillStrongKey(b []byte) error {
  _, err := io.ReadFull(rand.Reader, b)
  return err
}

func NewStrongKey(n int) ([]byte, error) {
  b := make([]byte, n)
  if err := FillStrongKey(b); err != nil {
    return nil, err
  }
  return b, nil
}

func Sign(buf, key []byte) ([]byte, error) {
  var res bytes.Buffer

  // Write the payload
  if _, err := res.Write(buf); err != nil {
    return nil, err
  }

  // Add a timestamp
  if err := binary.Write(
    &res,
    binary.LittleEndian,
    time.Now().UTC().Unix()); err != nil {
    return nil, err
  }

  // Add an hmac signature
  m := hmac.New(sha256.New, key)
  if _, err := m.Write(res.Bytes()); err != nil {
    return nil, err
  }

  if _, err := res.Write(m.Sum(nil)); err != nil {
    return nil, err
  }

  return res.Bytes(), nil
}

func toTime(b []byte) (time.Time, error) {
  var t int64
  buf := bytes.NewBuffer(b)
  if err := binary.Read(buf, binary.LittleEndian, &t); err != nil {
    return time.Time{}, err
  }
  return time.Unix(t, 0), nil
}

func Verify(buf, key []byte) ([]byte, time.Duration, error) {
  m := hmac.New(sha256.New, key)
  bs := len(buf)
  hs := m.Size()

  // buffer must be bigger than sig + time
  if bs <= hs+8 {
    return nil, time.Duration(0), ErrInvalidSignature
  }

  // compute the hmac of the payload
  if _, err := m.Write(buf[:bs-hs]); err != nil {
    return nil, time.Duration(0), ErrInvalidSignature
  }

  // verify that hmac matches
  if !hmac.Equal(m.Sum(nil), buf[bs-hs:]) {
    return nil, time.Duration(0), ErrInvalidSignature
  }

  // all is good, extract the time
  t, err := toTime(buf[bs-hs-8:])
  if err != nil {
    return nil, time.Duration(0), ErrInvalidSignature
  }

  return buf[:bs-hs-8], time.Now().UTC().Sub(t), nil
}
