package store

import (
  "errors"
)

var (
  ErrNotFound = errors.New("record not found")
)

const (
  UserKind    byte = 0x0
  SessionKind byte = 0x1
)
